package apps

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"servermanager/internal/app/apps"
	"servermanager/internal/app/metrics"
	"servermanager/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Controller struct {
	service *apps.Service
}

func NewController(d apps.Dependencies) *Controller {
	return &Controller{service: apps.NewService(d)}
}

func (c *Controller) Run(ginc *gin.Context) {
	appID, ok := appID(ginc)
	if !ok {
		return
	}

	var req RunRequest
	if err := ginc.ShouldBindJSON(&req); err != nil {
		respondErr(ginc, http.StatusBadRequest, "invalid request body")
		return
	}

	containerID, err := c.service.Run.Execute(ginc.Request.Context(), apps.RunInput{
		AppID:          appID,
		Image:          req.Image,
		Env:            req.Env,
		Port:           req.Port,
		StartCommand:   req.StartCommand,
		MemoryLimit:    req.MemoryLimit,
		CpuLimit:       req.CpuLimit,
		Runtime:        req.Runtime,
		Volumes:        req.Volumes,
		ReadonlyRootfs: req.ReadonlyRootfs,
		DatabaseID:     req.DatabaseID,
	})
	if err != nil {
		respondServiceErr(ginc, "run", appID, err)
		return
	}

	ginc.JSON(http.StatusCreated, RunResponse{
		ContainerID: containerID,
		Name:        domain.AppContainerName(appID),
	})
}

func (c *Controller) Start(ginc *gin.Context) {
	appID, ok := appID(ginc)
	if !ok {
		return
	}
	if err := c.service.Start.Execute(ginc.Request.Context(), appID); err != nil {
		respondServiceErr(ginc, "start", appID, err)
		return
	}
	ginc.Status(http.StatusOK)
}

func (c *Controller) Stop(ginc *gin.Context) {
	appID, ok := appID(ginc)
	if !ok {
		return
	}
	if err := c.service.Stop.Execute(ginc.Request.Context(), appID); err != nil {
		respondServiceErr(ginc, "stop", appID, err)
		return
	}
	ginc.Status(http.StatusOK)
}

func (c *Controller) Remove(ginc *gin.Context) {
	appID, ok := appID(ginc)
	if !ok {
		return
	}
	if err := c.service.Remove.Execute(ginc.Request.Context(), appID); err != nil {
		respondServiceErr(ginc, "remove", appID, err)
		return
	}
	ginc.Status(http.StatusOK)
}

func (c *Controller) Status(ginc *gin.Context) {
	appID, ok := appID(ginc)
	if !ok {
		return
	}
	state, err := c.service.Status.Execute(ginc.Request.Context(), appID)
	if err != nil {
		respondServiceErr(ginc, "status", appID, err)
		return
	}

	resp := ContainerStateResponse{
		Exists:       state.Exists,
		Running:      state.Running,
		RestartCount: state.RestartCount,
		Status:       state.Status,
	}
	if !state.StartedAt.IsZero() {
		resp.StartedAt = &state.StartedAt
	}
	ginc.JSON(http.StatusOK, resp)
}

func (c *Controller) Logs(ginc *gin.Context) {
	appID, ok := appID(ginc)
	if !ok {
		return
	}
	opts, ok := parseLogQuery(ginc)
	if !ok {
		return
	}

	lines, err := c.service.Logs.Execute(ginc.Request.Context(), appID, opts)
	if err != nil {
		respondServiceErr(ginc, "logs", appID, err)
		return
	}
	ginc.JSON(http.StatusOK, toLogsResponse(lines))
}

// StreamLogs serves a live NDJSON log tail (one LogEntryResponse per line,
// flushed as produced): recent history per ?tail=, then new lines as the
// container writes them. The api bridges this stream onto a browser
// WebSocket. Ends when the client disconnects or the container's log
// stream closes.
func (c *Controller) StreamLogs(ginc *gin.Context) {
	appID, ok := appID(ginc)
	if !ok {
		return
	}
	opts, ok := parseLogQuery(ginc)
	if !ok {
		return
	}

	lines, err := c.service.Follow.Execute(ginc.Request.Context(), appID, opts)
	if err != nil {
		respondServiceErr(ginc, "follow logs", appID, err)
		return
	}

	ginc.Header("Content-Type", "application/x-ndjson")
	ginc.Header("Cache-Control", "no-cache")
	ginc.Status(http.StatusOK)
	ginc.Writer.Flush()

	encoder := json.NewEncoder(ginc.Writer)
	index := 0
	for line := range lines {
		// Encode appends the newline NDJSON needs.
		if err := encoder.Encode(toLogEntryResponse(line, index)); err != nil {
			return // client went away; ctx cancellation drains the channel
		}
		ginc.Writer.Flush()
		index++
	}
}

// Metrics serves the collector-fed series for the app. A valid UUID with no
// samples is a 200 with empty series — the store is the local truth, "not
// deployed" is the api's judgement call.
func (c *Controller) Metrics(ginc *gin.Context) {
	appID, ok := appID(ginc)
	if !ok {
		return
	}
	rng, err := metrics.ParseRange(ginc.Query("range"))
	if err != nil {
		respondErr(ginc, http.StatusBadRequest, "range must be 1h or 24h")
		return
	}
	ginc.JSON(http.StatusOK, toMetricsResponse(rng, c.service.Metrics.Execute(appID, rng)))
}

const (
	defaultLogTail = 200
	maxLogTail     = 1000
)

// parseLogQuery validates ?tail= and ?since=. tail is always concrete so the
// adapter never does an unbounded read; values above the cap are clamped, not
// rejected — the caller still gets the freshest lines.
func parseLogQuery(ginc *gin.Context) (domain.LogOptions, bool) {
	opts := domain.LogOptions{Tail: defaultLogTail}

	if raw := ginc.Query("tail"); raw != "" {
		tail, err := strconv.Atoi(raw)
		if err != nil || tail < 1 {
			respondErr(ginc, http.StatusBadRequest, "tail must be a positive integer")
			return domain.LogOptions{}, false
		}
		opts.Tail = min(tail, maxLogTail)
	}

	if raw := ginc.Query("since"); raw != "" {
		since, err := time.Parse(time.RFC3339, raw)
		if err != nil {
			respondErr(ginc, http.StatusBadRequest, "since must be an RFC3339 timestamp")
			return domain.LogOptions{}, false
		}
		opts.Since = since
	}

	return opts, true
}

// appID validates the :id path param. Container and network names are
// derived from it, so it must be a well-formed UUID — never free text.
func appID(ginc *gin.Context) (string, bool) {
	id, err := uuid.Parse(ginc.Param("id"))
	if err != nil {
		respondErr(ginc, http.StatusBadRequest, "id must be a UUID")
		return "", false
	}
	return id.String(), true
}

func respondErr(ginc *gin.Context, code int, msg string) {
	ginc.JSON(code, gin.H{"error": msg})
}

// respondServiceErr maps domain sentinels onto HTTP statuses. Unknown errors
// become a generic 500 — Docker is a local dependency, not an upstream hop —
// with the detail kept in the log, not the response.
func respondServiceErr(ginc *gin.Context, op, appID string, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalid), errors.Is(err, domain.ErrImageMissing):
		respondErr(ginc, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrNotFound):
		respondErr(ginc, http.StatusNotFound, "no container for this application")
	case errors.Is(err, domain.ErrConflict):
		respondErr(ginc, http.StatusConflict, "container already exists; remove it first")
	default:
		slog.Error("app operation failed", "op", op, "app_id", appID, "error", err)
		respondErr(ginc, http.StatusInternalServerError, "internal error")
	}
}
