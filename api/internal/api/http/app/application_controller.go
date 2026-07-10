package app

import (
	"api/internal/api/dtos"
	"api/internal/api/middlewares"
	"api/internal/app/app"
	"api/internal/app/ports"
	"api/internal/app/subdomain"
	"api/internal/domain"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/coder/websocket"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Controller struct {
	appService       *app.Service
	subdomainService *subdomain.Service
}

func NewController(d app.Dependencies, sd subdomain.Dependencies) *Controller {
	return &Controller{
		appService:       app.NewService(d),
		subdomainService: subdomain.NewService(sd),
	}
}

func validType(t domain.ApplicationType) bool {
	switch t {
	case domain.ApplicationTypeNodejs:
		return true
	default:
		return false
	}
}

func normalizeApplicationType(raw string) (domain.ApplicationType, bool) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "nodejs", "node.js", "node":
		return domain.ApplicationTypeNodejs, true
	default:
		return "", false
	}
}

func (c *Controller) Create(ginc *gin.Context) {
	var req dtos.CreateApplicationRequest
	err := ginc.ShouldBindBodyWithJSON(&req)
	appType, ok := normalizeApplicationType(req.Type)
	if err != nil {
		middlewares.Respond(ginc, http.StatusBadRequest, "Invalid JSON or missing value", nil)
		return
	}
	if !ok || !validType(appType) {
		middlewares.Respond(ginc, http.StatusBadRequest, "Unsupported application type", nil)
		return
	}

	userID := ginc.GetString("userID")
	subdomainId := ginc.Param("id")

	if !middlewares.CheckOwnership(ginc, c.subdomainService.CheckUser.Execute, userID, subdomainId, "subdomain") {
		return
	}

	appInput := app.ApplicationInput{
		SubdomainID:  subdomainId,
		Name:         req.Name,
		Type:         appType,
		RepoUrl:      &req.RepoUrl,
		Branch:       &req.Branch,
		StartCommand: &req.StartCommand,
		BuildCommand: &req.BuildCommand,
		Status:       domain.ApplicationStatusPending,
	}

	err = c.appService.Create.Execute(ginc.Request.Context(), appInput)
	if err != nil {
		fmt.Println(err.Error())
		middlewares.Respond(ginc, http.StatusInternalServerError, "failed to create subdomain", nil)
		return
	}

	middlewares.Respond(ginc, http.StatusCreated, "success", nil)
}

func (c *Controller) Updates(ginc *gin.Context) {
	var req dtos.UpdateApplicationRequest
	if !middlewares.BindJSON(ginc, &req) {
		return
	}

	if req.Type != nil {
		appType, ok := normalizeApplicationType(*req.Type)
		if !ok || !validType(appType) {
			middlewares.Respond(ginc, http.StatusBadRequest, "Unsupported application type", nil)
			return
		}
		normalizedType := string(appType)
		req.Type = &normalizedType
	}

	userID := ginc.GetString("userID")
	subdomainId := ginc.Param("id")
	appId := ginc.Param("appId")

	if !middlewares.CheckOwnership(ginc, c.subdomainService.CheckUser.Execute, userID, subdomainId, "subdomain") {
		return
	}

	appInput := app.ApplicationUpdateInput{
		ID:           appId,
		Name:         req.Name,
		Branch:       req.Branch,
		RepoUrl:      req.RepoUrl,
		Type:         (*domain.ApplicationType)(req.Type),
		BuildCommand: req.BuildCommand,
		StartCommand: req.StartCommand,
	}

	err := c.appService.Update.Execute(ginc.Request.Context(), appInput)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		middlewares.Respond(ginc, http.StatusNotFound, "application not found", nil)
		return
	}
	if err != nil {
		fmt.Println(err.Error())
		middlewares.Respond(ginc, http.StatusInternalServerError, "failed to update application", nil)
		return
	}

	middlewares.Respond(ginc, http.StatusOK, "success", nil)
}

func (c *Controller) Delete(ginc *gin.Context) {
	userID := ginc.GetString("userID")
	subdomainId := ginc.Param("id")
	appId := ginc.Param("appId")

	if !middlewares.CheckOwnership(ginc, c.subdomainService.CheckUser.Execute, userID, subdomainId, "subdomain") {
		return
	}

	err := c.appService.Delete.Execute(ginc.Request.Context(), appId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		middlewares.Respond(ginc, http.StatusNotFound, "application not found", nil)
		return
	}
	if err != nil {
		fmt.Println(err.Error())
		middlewares.Respond(ginc, http.StatusInternalServerError, "failed to delete application", nil)
		return
	}

	middlewares.Respond(ginc, http.StatusOK, "success", nil)
}

func (c *Controller) Get(ginc *gin.Context) {
	userID := ginc.GetString("userID")
	subdomainId := ginc.Param("id")
	appId := ginc.Param("appId")

	if !middlewares.CheckOwnership(ginc, c.subdomainService.CheckUser.Execute, userID, subdomainId, "subdomain") {
		return
	}

	application, err := c.appService.Get.Execute(ginc.Request.Context(), appId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		middlewares.Respond(ginc, http.StatusNotFound, "application not found", nil)
		return
	}
	if err != nil {
		fmt.Println(err.Error())
		middlewares.Respond(ginc, http.StatusInternalServerError, "failed to get application", nil)
		return
	}

	applicationResponse := dtos.ApplicationListResponse{
		ID:      application.ID.String(),
		Name:    *application.Name,
		Type:    string(application.Type),
		Status:  string(application.Status),
		RepoUrl: *application.RepositoryURL,
		Branch:  *application.Branch,
	}

	middlewares.Respond(ginc, http.StatusOK, "success", applicationResponse)
}

const (
	defaultLogTail = 200
	maxLogTail     = 1000
)

func (c *Controller) GetLogs(ginc *gin.Context) {
	userID := ginc.GetString("userID")
	subdomainId := ginc.Param("id")
	appId := ginc.Param("appId")

	if !middlewares.CheckOwnership(ginc, c.subdomainService.CheckUser.Execute, userID, subdomainId, "subdomain") {
		return
	}

	opts := ports.LogOptions{Tail: defaultLogTail}
	if raw := ginc.Query("tail"); raw != "" {
		tail, err := strconv.Atoi(raw)
		if err != nil || tail < 1 {
			middlewares.Respond(ginc, http.StatusBadRequest, "tail must be a positive integer", nil)
			return
		}
		opts.Tail = min(tail, maxLogTail)
	}
	if raw := ginc.Query("since"); raw != "" {
		since, err := time.Parse(time.RFC3339, raw)
		if err != nil {
			middlewares.Respond(ginc, http.StatusBadRequest, "since must be an RFC3339 timestamp", nil)
			return
		}
		opts.Since = since
	}

	entries, err := c.appService.Logs.Execute(ginc.Request.Context(), subdomainId, appId, opts)
	if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, app.ErrNotInSubdomain) {
		middlewares.Respond(ginc, http.StatusNotFound, "application not found", nil)
		return
	}
	if err != nil {
		// The manager is an upstream hop, unlike the database — surface it as
		// a gateway problem and keep the detail in the server log.
		fmt.Println(err.Error())
		middlewares.Respond(ginc, http.StatusBadGateway, "log service unavailable", nil)
		return
	}

	logsResponse := make([]dtos.LogEntryResponse, 0, len(entries))
	for _, entry := range entries {
		logsResponse = append(logsResponse, toLogEntryResponse(entry))
	}

	middlewares.Respond(ginc, http.StatusOK, "success", logsResponse)
}

func toLogEntryResponse(entry ports.LogEntry) dtos.LogEntryResponse {
	return dtos.LogEntryResponse{
		ID:        entry.ID,
		Timestamp: entry.Timestamp.UTC().Format(time.RFC3339Nano),
		Level:     entry.Level,
		Message:   entry.Message,
	}
}

// StreamLogs bridges the manager's live log stream onto a browser WebSocket.
// Auth and ownership run on the handshake request (the JWT cookie rides
// along on same-origin WebSocket connects), so by the time the connection
// upgrades the caller is known to own the app.
func (c *Controller) StreamLogs(ginc *gin.Context) {
	userID := ginc.GetString("userID")
	subdomainId := ginc.Param("id")
	appId := ginc.Param("appId")

	if !middlewares.CheckOwnership(ginc, c.subdomainService.CheckUser.Execute, userID, subdomainId, "subdomain") {
		return
	}

	opts := ports.LogOptions{Tail: defaultLogTail}
	if raw := ginc.Query("tail"); raw != "" {
		tail, err := strconv.Atoi(raw)
		if err != nil || tail < 1 {
			middlewares.Respond(ginc, http.StatusBadRequest, "tail must be a positive integer", nil)
			return
		}
		opts.Tail = min(tail, maxLogTail)
	}

	// The stream must die with the socket: derive a cancelable context so the
	// manager-side follow unwinds when the browser disconnects.
	ctx, cancel := context.WithCancel(ginc.Request.Context())
	defer cancel()

	entries, err := c.appService.StreamLogs.Execute(ctx, subdomainId, appId, opts)
	if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, app.ErrNotInSubdomain) {
		middlewares.Respond(ginc, http.StatusNotFound, "application not found", nil)
		return
	}
	if errors.Is(err, ports.ErrAppNotDeployed) {
		// No container to follow — let the client fall back to polling,
		// which renders the "no logs yet" empty state.
		middlewares.Respond(ginc, http.StatusNotFound, "application not deployed", nil)
		return
	}
	if err != nil {
		fmt.Println(err.Error())
		middlewares.Respond(ginc, http.StatusBadGateway, "log service unavailable", nil)
		return
	}

	// Requires gin >= 1.12: 1.11's Hijack guard rejects coder/websocket's
	// WriteHeaderNow workaround ("gin: response already written").
	conn, err := websocket.Accept(ginc.Writer, ginc.Request, &websocket.AcceptOptions{
		// Same-origin in prod (reverse proxy); the vite dev server runs on a
		// different localhost port, so allow localhost explicitly.
		OriginPatterns: []string{"localhost:*", "127.0.0.1:*"},
	})
	if err != nil {
		// Accept already wrote the HTTP error response.
		return
	}
	defer conn.Close(websocket.StatusInternalError, "stream aborted")

	// Reads are only needed to surface client close/pings promptly.
	readCtx := conn.CloseRead(ctx)

	for {
		select {
		case entry, open := <-entries:
			if !open {
				_ = conn.Close(websocket.StatusNormalClosure, "log stream ended")
				return
			}
			payload, err := json.Marshal(toLogEntryResponse(entry))
			if err != nil {
				continue
			}
			if err := conn.Write(readCtx, websocket.MessageText, payload); err != nil {
				return // client gone; deferred cancel unwinds the follow
			}
		case <-readCtx.Done():
			_ = conn.Close(websocket.StatusNormalClosure, "client closed")
			return
		}
	}
}
