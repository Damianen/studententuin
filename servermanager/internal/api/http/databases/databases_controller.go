package databases

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"servermanager/internal/app/databases"
	"servermanager/internal/app/metrics"
	"servermanager/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Controller struct {
	service *databases.Service
}

func NewController(d databases.Dependencies) *Controller {
	return &Controller{service: databases.NewService(d)}
}

func (c *Controller) Provision(ginc *gin.Context) {
	dbID, ok := dbID(ginc)
	if !ok {
		return
	}

	var req ProvisionRequest
	if err := ginc.ShouldBindJSON(&req); err != nil {
		respondErr(ginc, http.StatusBadRequest, "invalid request body")
		return
	}

	provisionID, err := c.service.Provision.Execute(ginc.Request.Context(), databases.ProvisionInput{
		DBID:        dbID,
		Type:        req.Type,
		Version:     req.Version,
		DBName:      req.DBName,
		DBUser:      req.DBUser,
		DBPassword:  req.DBPassword,
		MemoryLimit: req.MemoryLimit,
		CpuLimit:    req.CpuLimit,
		Runtime:     req.Runtime,
	})
	if err != nil {
		respondServiceErr(ginc, "provision", dbID, err)
		return
	}

	ginc.JSON(http.StatusAccepted, ProvisionResponse{ProvisionID: provisionID})
}

func (c *Controller) Status(ginc *gin.Context) {
	dbID, ok := dbID(ginc)
	if !ok {
		return
	}
	status, err := c.service.Status.Execute(ginc.Request.Context(), dbID)
	if err != nil {
		respondServiceErr(ginc, "status", dbID, err)
		return
	}

	resp := DatabaseStatusResponse{
		Exists:       status.State.Exists,
		Running:      status.State.Running,
		RestartCount: status.State.RestartCount,
		Status:       status.State.Status,
		Health:       status.State.Health,
	}
	if !status.State.StartedAt.IsZero() {
		resp.StartedAt = &status.State.StartedAt
	}
	if job := status.Provision; job != nil {
		resp.Provision = &ProvisionStateResponse{
			ID:            job.ID,
			Status:        string(job.Status),
			Error:         job.Error,
			Image:         job.Image,
			ContainerID:   job.ContainerID,
			ContainerName: job.ContainerName,
			CreatedAt:     job.CreatedAt,
			UpdatedAt:     job.UpdatedAt,
		}
	}
	ginc.JSON(http.StatusOK, resp)
}

func (c *Controller) Delete(ginc *gin.Context) {
	dbID, ok := dbID(ginc)
	if !ok {
		return
	}
	if err := c.service.Delete.Execute(ginc.Request.Context(), dbID); err != nil {
		respondServiceErr(ginc, "delete", dbID, err)
		return
	}
	ginc.Status(http.StatusOK)
}

func (c *Controller) Logs(ginc *gin.Context) {
	dbID, ok := dbID(ginc)
	if !ok {
		return
	}
	opts, ok := parseLogQuery(ginc)
	if !ok {
		return
	}

	lines, err := c.service.Logs.Execute(ginc.Request.Context(), dbID, opts)
	if err != nil {
		respondServiceErr(ginc, "logs", dbID, err)
		return
	}
	ginc.JSON(http.StatusOK, toLogsResponse(lines))
}

// StreamLogs serves a live NDJSON log tail, mirroring the apps endpoint; the
// api bridges it onto a browser WebSocket.
func (c *Controller) StreamLogs(ginc *gin.Context) {
	dbID, ok := dbID(ginc)
	if !ok {
		return
	}
	opts, ok := parseLogQuery(ginc)
	if !ok {
		return
	}

	lines, err := c.service.Follow.Execute(ginc.Request.Context(), dbID, opts)
	if err != nil {
		respondServiceErr(ginc, "follow logs", dbID, err)
		return
	}

	ginc.Header("Content-Type", "application/x-ndjson")
	ginc.Header("Cache-Control", "no-cache")
	ginc.Status(http.StatusOK)
	ginc.Writer.Flush()

	encoder := json.NewEncoder(ginc.Writer)
	index := 0
	for line := range lines {
		if err := encoder.Encode(toLogEntryResponse(line, index)); err != nil {
			return // client went away; ctx cancellation drains the channel
		}
		ginc.Writer.Flush()
		index++
	}
}

// Metrics serves the collector-fed series for the database, mirroring the
// apps endpoint: a valid UUID with no samples is a 200 with empty series.
func (c *Controller) Metrics(ginc *gin.Context) {
	dbID, ok := dbID(ginc)
	if !ok {
		return
	}
	rng, err := metrics.ParseRange(ginc.Query("range"))
	if err != nil {
		respondErr(ginc, http.StatusBadRequest, "range must be 1h or 24h")
		return
	}
	ginc.JSON(http.StatusOK, toMetricsResponse(rng, c.service.Metrics.Execute(dbID, rng)))
}

const (
	defaultLogTail = 200
	maxLogTail     = 1000
)

// parseLogQuery validates ?tail= and ?since=, mirroring the apps endpoint.
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

// dbID validates the :id path param. Container, network, and volume names are
// derived from it, so it must be a well-formed UUID — never free text.
func dbID(ginc *gin.Context) (string, bool) {
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

// respondServiceErr maps domain sentinels onto HTTP statuses, mirroring the
// apps controller. Provision requests never log their body — the password
// travels in it.
func respondServiceErr(ginc *gin.Context, op, dbID string, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalid), errors.Is(err, domain.ErrImageMissing):
		respondErr(ginc, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrNotFound):
		respondErr(ginc, http.StatusNotFound, "no container for this database")
	case errors.Is(err, domain.ErrConflict):
		respondErr(ginc, http.StatusConflict, err.Error())
	default:
		slog.Error("database operation failed", "op", op, "db_id", dbID, "error", err)
		respondErr(ginc, http.StatusInternalServerError, "internal error")
	}
}
