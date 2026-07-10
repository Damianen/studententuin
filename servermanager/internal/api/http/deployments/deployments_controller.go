package deployments

import (
	"errors"
	"log/slog"
	"net/http"

	"servermanager/internal/app/deployments"
	"servermanager/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Controller struct {
	service *deployments.Service
}

func NewController(d deployments.Dependencies) *Controller {
	return &Controller{service: deployments.NewService(d)}
}

// Deploy accepts the async job: everything validatable without touching git
// or Docker is rejected here (400/409); the pipeline reports through the
// deployment resource afterwards.
func (c *Controller) Deploy(ginc *gin.Context) {
	appID, ok := pathUUID(ginc, "id")
	if !ok {
		return
	}

	var req DeployRequest
	if err := ginc.ShouldBindJSON(&req); err != nil {
		respondErr(ginc, http.StatusBadRequest, "invalid request body")
		return
	}

	deploymentID, err := c.service.Deploy.Execute(ginc.Request.Context(), deployments.DeployInput{
		AppID:         appID,
		RepositoryURL: req.RepositoryURL,
		Branch:        req.Branch,
		Type:          req.Type,
		BuildCommand:  req.BuildCommand,
		StartCommand:  req.StartCommand,
		Env:           req.Env,
		Port:          req.Port,
		MemoryLimit:   req.MemoryLimit,
		CpuLimit:      req.CpuLimit,
		Runtime:       req.Runtime,
	})
	if err != nil {
		respondServiceErr(ginc, "deploy", appID, err)
		return
	}

	ginc.JSON(http.StatusAccepted, DeployResponse{DeploymentID: deploymentID})
}

func (c *Controller) Get(ginc *gin.Context) {
	deploymentID, ok := pathUUID(ginc, "id")
	if !ok {
		return
	}

	job, err := c.service.Get.Execute(deploymentID)
	if err != nil {
		respondServiceErr(ginc, "get deployment", deploymentID, err)
		return
	}

	ginc.JSON(http.StatusOK, DeploymentResponse{
		ID:            job.ID,
		AppID:         job.AppID,
		Status:        string(job.Status),
		Error:         job.Error,
		CommitSHA:     job.CommitSHA,
		Image:         job.Image,
		ContainerID:   job.ContainerID,
		ContainerName: job.ContainerName,
		BuildLog:      job.BuildLog,
		CreatedAt:     job.CreatedAt,
		UpdatedAt:     job.UpdatedAt,
	})
}

// pathUUID validates a path param the same way the apps package does —
// derived names must never come from free text.
func pathUUID(ginc *gin.Context, name string) (string, bool) {
	id, err := uuid.Parse(ginc.Param(name))
	if err != nil {
		respondErr(ginc, http.StatusBadRequest, name+" must be a UUID")
		return "", false
	}
	return id.String(), true
}

func respondErr(ginc *gin.Context, code int, msg string) {
	ginc.JSON(code, gin.H{"error": msg})
}

func respondServiceErr(ginc *gin.Context, op, id string, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalid):
		respondErr(ginc, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrNotFound):
		respondErr(ginc, http.StatusNotFound, "no such deployment")
	case errors.Is(err, domain.ErrConflict):
		respondErr(ginc, http.StatusConflict, "a deployment is already in flight for this application")
	default:
		slog.Error("deployment operation failed", "op", op, "id", id, "error", err)
		respondErr(ginc, http.StatusInternalServerError, "internal error")
	}
}
