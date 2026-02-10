package app

import (
	"api/internal/api/dtos"
	"api/internal/api/middlewares"
	"api/internal/app/app"
	"api/internal/app/subdomain"
	"api/internal/domain"
	"errors"
	"fmt"
	"net/http"

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

func (c *Controller) Create(ginc *gin.Context) {
	var req dtos.CreateApplicationRequest
	err := ginc.ShouldBindBodyWithJSON(&req)
	if err != nil || !validType(domain.ApplicationType(req.Type)) {
		middlewares.Respond(ginc, http.StatusBadRequest, "Invalid JSON or missing value", nil)
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
		Type:         domain.ApplicationType(req.Type),
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
		if validType(domain.ApplicationType(*req.Type)) {
			middlewares.Respond(ginc, http.StatusBadRequest, "Invalid JSON or missing value", nil)
			return
		}
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

	middlewares.Respond(ginc, http.StatusNoContent, "success", nil)
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

	middlewares.Respond(ginc, http.StatusNoContent, "success", nil)
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
