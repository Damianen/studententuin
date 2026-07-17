package db

import (
	"api/internal/api/dtos"
	"api/internal/api/middlewares"
	"api/internal/app/db"
	"api/internal/app/ports"
	"api/internal/app/subdomain"
	"api/internal/domain"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/coder/websocket"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Controller struct {
	dbService        *db.Service
	subdomainService *subdomain.Service
}

func NewController(d db.Dependencies, sd subdomain.Dependencies) *Controller {
	return &Controller{
		dbService:        db.NewService(d),
		subdomainService: subdomain.NewService(sd),
	}
}

// supportedEngines is the api-side twin of the manager's image allowlist:
// postgres-only through phase 5 (§10.3), so a disallowed engine fails here
// with a clean 400 instead of a created-then-failed record.
var supportedEngines = map[domain.DatabaseType]map[string]bool{
	domain.DatabaseTypePostgres: {"16": true, "17": true},
}

func validEngine(t domain.DatabaseType, version string) bool {
	return supportedEngines[t][version]
}

// dbNameRe mirrors the manager's identifier rule; the value becomes
// POSTGRES_DB and part of the connection string.
var dbNameRe = regexp.MustCompile(`^[a-z_][a-z0-9_]{0,62}$`)

func (c *Controller) Create(ginc *gin.Context) {
	var req dtos.CreateDatabaseRequest
	err := ginc.ShouldBindBodyWithJSON(&req)
	if err != nil {
		middlewares.Respond(ginc, http.StatusBadRequest, "Invalid JSON or missing value", nil)
		return
	}
	if !validEngine(domain.DatabaseType(req.Type), req.Version) {
		middlewares.Respond(ginc, http.StatusBadRequest, "unsupported database type or version (postgres 16/17 for now)", nil)
		return
	}
	if !dbNameRe.MatchString(req.DbName) {
		middlewares.Respond(ginc, http.StatusBadRequest, "db_name must be lowercase letters, digits or underscores, starting with a letter", nil)
		return
	}

	userID := ginc.GetString("userID")
	subdomainId := ginc.Param("id")

	if !middlewares.CheckOwnership(ginc, c.subdomainService.CheckUser.Execute, userID, subdomainId, "subdomain") {
		return
	}

	dbInput := db.DatabaseInput{
		SubdomainID: subdomainId,
		Name:        req.Name,
		Type:        domain.DatabaseType(req.Type),
		Version:     req.Version,
		DbName:      req.DbName,
	}

	err = c.dbService.Create.Execute(ginc.Request.Context(), dbInput)
	if err != nil {
		fmt.Println(err.Error())
		middlewares.Respond(ginc, http.StatusInternalServerError, "failed to create database", nil)
		return
	}

	middlewares.Respond(ginc, http.StatusCreated, "success", nil)
}

func (c *Controller) Update(ginc *gin.Context) {
	var req dtos.UpdateDatabaseRequest
	if !middlewares.BindJSON(ginc, &req) {
		return
	}

	// Type/version changes would require a re-provision; the settings tab
	// only edits the display name. Reject engine changes outright.
	if req.Type != nil || req.Version != nil {
		middlewares.Respond(ginc, http.StatusBadRequest, "type and version cannot be changed — delete and recreate the database", nil)
		return
	}

	userID := ginc.GetString("userID")
	subdomainId := ginc.Param("id")
	dbId := ginc.Param("dbId")

	if !middlewares.CheckOwnership(ginc, c.subdomainService.CheckUser.Execute, userID, subdomainId, "subdomain") {
		return
	}

	databaseInput := db.DatabaseUpdateInput{
		ID:          dbId,
		SubdomainID: subdomainId,
		Name:        req.Name,
	}

	err := c.dbService.Update.Execute(ginc.Request.Context(), databaseInput)
	if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, db.ErrNotInSubdomain) {
		middlewares.Respond(ginc, http.StatusNotFound, "database not found", nil)
		return
	}
	if err != nil {
		fmt.Println(err.Error())
		middlewares.Respond(ginc, http.StatusInternalServerError, "failed to update the database", nil)
		return
	}

	middlewares.Respond(ginc, http.StatusOK, "success", nil)
}

func (c *Controller) Delete(ginc *gin.Context) {
	userID := ginc.GetString("userID")
	subdomainId := ginc.Param("id")
	id := ginc.Param("dbId")

	if !middlewares.CheckOwnership(ginc, c.subdomainService.CheckUser.Execute, userID, subdomainId, "subdomain") {
		return
	}

	err := c.dbService.Delete.Execute(ginc.Request.Context(), subdomainId, id)
	if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, db.ErrNotInSubdomain) {
		middlewares.Respond(ginc, http.StatusNotFound, "database not found", nil)
		return
	}
	if errors.Is(err, ports.ErrProvisionInFlight) {
		middlewares.Respond(ginc, http.StatusConflict, "database is still provisioning — try again shortly", nil)
		return
	}
	if err != nil {
		fmt.Println(err.Error())
		middlewares.Respond(ginc, http.StatusInternalServerError, "failed to delete the database", nil)
		return
	}

	middlewares.Respond(ginc, http.StatusOK, "success", nil)
}

func (c *Controller) Get(ginc *gin.Context) {
	userID := ginc.GetString("userID")
	subdomainId := ginc.Param("id")
	id := ginc.Param("dbId")

	if !middlewares.CheckOwnership(ginc, c.subdomainService.CheckUser.Execute, userID, subdomainId, "subdomain") {
		return
	}

	database, err := c.dbService.Get.Execute(ginc.Request.Context(), subdomainId, id)
	if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, db.ErrNotInSubdomain) {
		middlewares.Respond(ginc, http.StatusNotFound, "database not found", nil)
		return
	}
	if err != nil {
		fmt.Println(err.Error())
		middlewares.Respond(ginc, http.StatusInternalServerError, "failed to get database", nil)
		return
	}

	databaseResponse := dtos.DatabaseListResponse{
		ID:      database.ID.String(),
		Name:    database.Name,
		Type:    string(database.Type),
		Version: database.Version,
		Status:  string(database.Status),
	}
	if database.ConnectionString != nil {
		databaseResponse.ConnectionString = *database.ConnectionString
	}

	middlewares.Respond(ginc, http.StatusOK, "success", databaseResponse)
}

const (
	defaultLogTail = 200
	maxLogTail     = 1000
)

func (c *Controller) GetLogs(ginc *gin.Context) {
	userID := ginc.GetString("userID")
	subdomainId := ginc.Param("id")
	dbId := ginc.Param("dbId")

	if !middlewares.CheckOwnership(ginc, c.subdomainService.CheckUser.Execute, userID, subdomainId, "subdomain") {
		return
	}

	opts, ok := parseLogQuery(ginc)
	if !ok {
		return
	}

	entries, err := c.dbService.Logs.Execute(ginc.Request.Context(), subdomainId, dbId, opts)
	if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, db.ErrNotInSubdomain) {
		middlewares.Respond(ginc, http.StatusNotFound, "database not found", nil)
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

// metricsRange is fixed for now: the dashboard always shows the last 24
// hours (the manager also accepts 1h).
const metricsRange = "24h"

func (c *Controller) GetMetrics(ginc *gin.Context) {
	userID := ginc.GetString("userID")
	subdomainId := ginc.Param("id")
	dbId := ginc.Param("dbId")

	if !middlewares.CheckOwnership(ginc, c.subdomainService.CheckUser.Execute, userID, subdomainId, "subdomain") {
		return
	}

	metrics, err := c.dbService.Metrics.Execute(ginc.Request.Context(), subdomainId, dbId, metricsRange)
	if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, db.ErrNotInSubdomain) {
		middlewares.Respond(ginc, http.StatusNotFound, "database not found", nil)
		return
	}
	if err != nil {
		// The manager is an upstream hop, unlike the database — surface it as
		// a gateway problem and keep the detail in the server log.
		fmt.Println(err.Error())
		middlewares.Respond(ginc, http.StatusBadGateway, "metrics service unavailable", nil)
		return
	}

	middlewares.Respond(ginc, http.StatusOK, "success", toMetricsResponse(metrics))
}

func toMetricsResponse(metrics *ports.MetricsResponse) dtos.MetricsResponse {
	series := make(map[string][]dtos.MetricPointResponse, len(metrics.Series))
	for key, points := range metrics.Series {
		mapped := make([]dtos.MetricPointResponse, 0, len(points))
		for _, point := range points {
			mapped = append(mapped, dtos.MetricPointResponse{
				Time:  point.Time.UTC().Format(time.RFC3339),
				Value: point.Value,
			})
		}
		series[key] = mapped
	}
	return dtos.MetricsResponse{Range: metrics.Range, Series: series}
}

// StreamLogs bridges the manager's live database log tail onto a browser
// WebSocket, mirroring the application handler.
func (c *Controller) StreamLogs(ginc *gin.Context) {
	userID := ginc.GetString("userID")
	subdomainId := ginc.Param("id")
	dbId := ginc.Param("dbId")

	if !middlewares.CheckOwnership(ginc, c.subdomainService.CheckUser.Execute, userID, subdomainId, "subdomain") {
		return
	}

	opts, ok := parseLogQuery(ginc)
	if !ok {
		return
	}

	// The stream must die with the socket: derive a cancelable context so the
	// manager-side follow unwinds when the browser disconnects.
	ctx, cancel := context.WithCancel(ginc.Request.Context())
	defer cancel()

	entries, err := c.dbService.StreamLogs.Execute(ctx, subdomainId, dbId, opts)
	if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, db.ErrNotInSubdomain) {
		middlewares.Respond(ginc, http.StatusNotFound, "database not found", nil)
		return
	}
	if errors.Is(err, ports.ErrDatabaseNotProvisioned) {
		// No container to follow — let the client fall back to polling,
		// which renders the "no logs yet" empty state.
		middlewares.Respond(ginc, http.StatusNotFound, "database not provisioned", nil)
		return
	}
	if err != nil {
		fmt.Println(err.Error())
		middlewares.Respond(ginc, http.StatusBadGateway, "log service unavailable", nil)
		return
	}

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

func parseLogQuery(ginc *gin.Context) (ports.LogOptions, bool) {
	opts := ports.LogOptions{Tail: defaultLogTail}
	if raw := ginc.Query("tail"); raw != "" {
		tail, err := strconv.Atoi(raw)
		if err != nil || tail < 1 {
			middlewares.Respond(ginc, http.StatusBadRequest, "tail must be a positive integer", nil)
			return ports.LogOptions{}, false
		}
		opts.Tail = min(tail, maxLogTail)
	}
	if raw := ginc.Query("since"); raw != "" {
		since, err := time.Parse(time.RFC3339, raw)
		if err != nil {
			middlewares.Respond(ginc, http.StatusBadRequest, "since must be an RFC3339 timestamp", nil)
			return ports.LogOptions{}, false
		}
		opts.Since = since
	}
	return opts, true
}

func toLogEntryResponse(entry ports.LogEntry) dtos.LogEntryResponse {
	return dtos.LogEntryResponse{
		ID:        entry.ID,
		Timestamp: entry.Timestamp.UTC().Format(time.RFC3339Nano),
		Level:     entry.Level,
		Message:   entry.Message,
	}
}
