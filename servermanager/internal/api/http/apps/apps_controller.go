package apps

import (
	"errors"
	"log/slog"
	"net/http"

	"servermanager/internal/app/apps"
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
