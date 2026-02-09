package subdomain

import (
	"api/internal/api/dtos"
	"api/internal/api/middlewares"
	"api/internal/app/subdomain"
	"api/internal/infra/postgres"
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


type Controller struct {
	service *subdomain.Service
}

func NewController(d subdomain.Dependencies) *Controller {
	return &Controller{
		service: subdomain.NewService(d),
	}
}

func (c *Controller) Create(ginc *gin.Context) {
	context := ginc.Request.Context()

	var req dtos.CreateSubdomainRequest
	err := ginc.ShouldBindBodyWithJSON(&req)
	if err != nil {
		fmt.Println(err.Error())
		middlewares.Respond(ginc, 400, "Invalid JSON or missing value", nil)
		return
	}

	userID := ginc.GetString("userID")

	subdomainInput := subdomain.SubdomainInput{
		Name: req.Name,
		FullDomain: req.FullDomain,
		UserID: userID,
	}

	err = c.service.Create.Execute(context, subdomainInput)
	if errors.Is(err, postgres.ErrFullDomainAlreadyInUse) {
    	middlewares.Respond(ginc, 409, "domain already in use", nil)
    	return
	}
	if err != nil {
		fmt.Println(err.Error())
		middlewares.Respond(ginc, 500, "failed to create subdomain", nil)
		return
	}

	middlewares.Respond(ginc, 201, "success", nil)
}

func (c *Controller) Update(ginc *gin.Context) {
	context := ginc.Request.Context()

	var req dtos.UpdateSubdomainRequest
	err := ginc.ShouldBindBodyWithJSON(&req)
	if err != nil {
		fmt.Println(err.Error())
		middlewares.Respond(ginc, 400, "Invalid JSON or missing value", nil)
		return
	}

	id := ginc.Param("id")
	userID := ginc.GetString("userID")

	allowed, err := c.service.CheckUser.Execute(context, userID, id)
	if errors.Is(err, gorm.ErrRecordNotFound) || (err != nil && strings.Contains(err.Error(), "invalid input syntax")) {
		middlewares.Respond(ginc, 404, "subdomain not found", nil)
		return
	}
	if err != nil {
		fmt.Println(err.Error())
		middlewares.Respond(ginc, 500, "Internal server error", nil)
		return
	}

	if !allowed {
		middlewares.Respond(ginc, 403, "Unauthorized", nil)
		return
	}


	subdomainInput := subdomain.SubdomainUpdateInput{
		ID: id,
		Name: req.Name,
		FullDomain: req.FullDomain,
	}

	err = c.service.Update.Execute(context, subdomainInput)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		middlewares.Respond(ginc, 404, "subdomain not found", nil)
		return
	}
	if err != nil {
		fmt.Println(err.Error())
		middlewares.Respond(ginc, 500, "failed to update subdomain", nil)
		return
	}

	ginc.Status(204)
}

func (c *Controller) Delete(ginc *gin.Context) {
	context := ginc.Request.Context()

	id := ginc.Param("id")
	userID := ginc.GetString("userID")

	allowed, err := c.service.CheckUser.Execute(context, userID, id)
	if errors.Is(err, gorm.ErrRecordNotFound) || (err != nil && strings.Contains(err.Error(), "invalid input syntax")) {
		middlewares.Respond(ginc, 404, "subdomain not found", nil)
		return
	}
	if err != nil {
		fmt.Println(err.Error())
		middlewares.Respond(ginc, 500, "Internal server error", nil)
		return
	}

	if !allowed {
		middlewares.Respond(ginc, 403, "Unauthorized", nil)
		return
	}

	err = c.service.Delete.Execute(context, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		middlewares.Respond(ginc, 404, "subdomain not found", nil)
		return
	}
	if err != nil {
		fmt.Println(err.Error())
		middlewares.Respond(ginc, 500, "failed to delete subdomain", nil)
		return
	}

	ginc.Status(204)
}

func (c *Controller) Get(ginc *gin.Context) {
	context := ginc.Request.Context()

	id := ginc.Param("id")
	userID := ginc.GetString("userID")

	allowed, err := c.service.CheckUser.Execute(context, userID, id)
	if errors.Is(err, gorm.ErrRecordNotFound) || (err != nil && strings.Contains(err.Error(), "invalid input syntax")) {
		middlewares.Respond(ginc, 404, "subdomain not found", nil)
		return
	}
	if err != nil {
		fmt.Println(err.Error())
		middlewares.Respond(ginc, 500, "Internal server error", nil)
		return
	}

	if !allowed {
		middlewares.Respond(ginc, 403, "Unauthorized", nil)
		return
	}

	subdomain, err := c.service.Get.Execute(context, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		middlewares.Respond(ginc, 404, "subdomain not found", nil)
		return
	}
	if err != nil {
		fmt.Println(err.Error())
		middlewares.Respond(ginc, 500, "failed to get subdomain", nil)
		return
	}

	subdomainResponse := dtos.SubdomainResponse{
		ID: subdomain.ID.String(),
		Name: subdomain.Name,
		FullDomain: subdomain.FullDomain,
		IsActive: subdomain.IsActive,
	}

	middlewares.Respond(ginc, 200, "success", subdomainResponse)
}
