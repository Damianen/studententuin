package app

import (
	"api/internal/api/dtos"
	"api/internal/api/middlewares"
	"api/internal/app/app"
	"api/internal/app/subdomain"
	"api/internal/domain"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Controller struct {
	appService *app.Service
	subdomainService *subdomain.Service
}

func NewController(d app.Dependencies, sd subdomain.Dependencies) *Controller {
	return &Controller{
		appService: app.NewService(d),
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
	context := ginc.Request.Context()

	var req dtos.CreateApplicationRequest
	err := ginc.ShouldBindBodyWithJSON(&req)
	if err != nil || !validType(domain.ApplicationType(req.Type)) {
		middlewares.Respond(ginc, 400, "Invalid JSON or missing value", nil)
		return
	}

	userID := ginc.GetString("userID")
	subdomainId := ginc.Param("id")

	allowed, err := c.subdomainService.CheckUser.Execute(context, userID, subdomainId)
	if err != nil {
		fmt.Println(err.Error())
		middlewares.Respond(ginc, 500, "Internal server error", nil)
		return
	}

	if !allowed {
		middlewares.Respond(ginc, 403, "unauthorized", nil)
		return
	}

	appInput := app.ApplicationInput{
		SubdomainID: subdomainId,
		Name: req.Name,
		Type: domain.ApplicationType(req.Type),
		RepoUrl: &req.RepoUrl,
		Branch: &req.Branch,
		StartCommand: &req.StartCommand,
		BuildCommand: &req.BuildCommand,
		Status: domain.ApplicationStatusPending,
	}

	err = c.appService.Create.Execute(context, appInput)
	if err != nil {
		fmt.Println(err.Error())
		middlewares.Respond(ginc, 500, "failed to create subdomain", nil)
		return
	}

	middlewares.Respond(ginc, 201, "success", nil)
}

func (c *Controller) Updates(ginc *gin.Context) {
	context := ginc.Request.Context()

	var req dtos.UpdateApplicationRequest
	err := ginc.ShouldBindBodyWithJSON(&req)
	if err != nil {
		middlewares.Respond(ginc, 400, "Invalid JSON or missing value", nil)
		return
	}

	if req.Type != nil {
		if validType(domain.ApplicationType(*req.Type)) {
			middlewares.Respond(ginc, 400, "Invalid JSON or missing value", nil)
			return
		}
	}

	userID := ginc.GetString("userID")
	subdomainId := ginc.Param("id")
	appId := ginc.Param("appId")

	allowed, err := c.subdomainService.CheckUser.Execute(context, userID, subdomainId)
	if err != nil {
		fmt.Println(err.Error())
		middlewares.Respond(ginc, 500, "Internal server error", nil)
		return
	}

	if !allowed {
		middlewares.Respond(ginc, 403, "unauthorized", nil)
		return
	}

	appInput := app.ApplicationUpdateInput{
		ID: appId,
		Name: req.Name,
		Branch: req.Branch,
		RepoUrl: req.RepoUrl,
		Type: (*domain.ApplicationType)(req.Type),
		BuildCommand: req.BuildCommand,
		StartCommand: req.StartCommand,
	}

	err = c.appService.Update.Execute(context, appInput)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		middlewares.Respond(ginc, 404, "application not found", nil)
		return
	}
	if err != nil {
		fmt.Println(err.Error())
		middlewares.Respond(ginc, 500, "failed to update application", nil)
		return
	}

	ginc.Status(204)
}

func (c *Controller) Delete(ginc *gin.Context) {
	context := ginc.Request.Context()

	userID := ginc.GetString("userID")
	subdomainId := ginc.Param("id")
	appId := ginc.Param("appId")

	allowed, err := c.subdomainService.CheckUser.Execute(context, userID, subdomainId)
	if err != nil {
		fmt.Println(err.Error())
		middlewares.Respond(ginc, 500, "Internal server error", nil)
		return
	}

	if !allowed {
		middlewares.Respond(ginc, 403, "unauthorized", nil)
		return
	}


	err = c.appService.Delete.Execute(context, appId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		middlewares.Respond(ginc, 404, "application not found", nil)
		return
	}
	if err != nil {
		fmt.Println(err.Error())
		middlewares.Respond(ginc, 500, "failed to delete application", nil)
		return
	}

	ginc.Status(204)
}

func (c *Controller) Get(ginc *gin.Context) {
	context := ginc.Request.Context()

	userID := ginc.GetString("userID")
	subdomainId := ginc.Param("id")
	appId := ginc.Param("appId")

	allowed, err := c.subdomainService.CheckUser.Execute(context, userID, subdomainId)
	if err != nil {
		fmt.Println(err.Error())
		middlewares.Respond(ginc, 500, "Internal server error", nil)
		return
	}

	if !allowed {
		middlewares.Respond(ginc, 403, "unauthorized", nil)
		return
	}


	application, err := c.appService.Get.Execute(context, appId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		middlewares.Respond(ginc, 404, "application not found", nil)
		return
	}
	if err != nil {
		fmt.Println(err.Error())
		middlewares.Respond(ginc, 500, "failed to get application", nil)
		return
	}

	applicationResponse := dtos.ApplicationListResponse{
		ID: application.ID.String(),
		Name: *application.Name,
		Type: string(application.Type),
		Status: string(application.Status),
		RepoUrl: *application.RepositoryURL,
		Branch: *application.Branch,
	}

	middlewares.Respond(ginc, 200, "success", applicationResponse)
}
