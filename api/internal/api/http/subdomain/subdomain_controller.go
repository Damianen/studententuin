package subdomain

import (
	"api/internal/api/dtos"
	"api/internal/api/middlewares"
	"api/internal/app/subdomain"
	"fmt"

	"github.com/gin-gonic/gin"
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

	userID := ginc.GetString("UserID")

	subdomainInput := subdomain.SubdomainInput{
		Name: req.Name,
		FullDomain: req.FullDomain,
		UserID: userID,
	}

	err = c.service.Create.Execute(context, subdomainInput)
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
	userID := ginc.GetString("UserID")

	correct, err := c.service.CheckUser.Execute(context, userID, id)
	if err != nil {
		fmt.Println(err.Error())
		middlewares.Respond(ginc, 500, "Internal server error", nil)
		return
	}

	if !correct {
		middlewares.Respond(ginc, 401, "Unauthorized", nil)
		return
	}


	subdomainInput := subdomain.SubdomainUpdateInput{
		ID: id,
		Name: req.Name,
		FullDomain: req.FullDomain,
	}

	err = c.service.Update.Execute(context, subdomainInput)
	if err != nil {
		fmt.Println(err.Error())
		middlewares.Respond(ginc, 500, "failed to update subdomain", nil)
		return
	}

	middlewares.Respond(ginc, 204, "success", nil)
}

func (c *Controller) Delete(ginc *gin.Context) {
	context := ginc.Request.Context()

	id := ginc.Param("id")
	userID := ginc.GetString("UserID")

	correct, err := c.service.CheckUser.Execute(context, userID, id)
	if err != nil {
		fmt.Println(err.Error())
		middlewares.Respond(ginc, 500, "Internal server error", nil)
		return
	}

	if !correct {
		middlewares.Respond(ginc, 401, "Unauthorized", nil)
		return
	}

	err = c.service.Delete.Execute(context, id)
	if err != nil {
		fmt.Println(err.Error())
		middlewares.Respond(ginc, 500, "failed to delete subdomain", nil)
		return
	}

	middlewares.Respond(ginc, 204, "success", nil)
}

func (c *Controller) Get(ginc *gin.Context) {
	context := ginc.Request.Context()

	id := ginc.Param("id")
	userID := ginc.GetString("UserID")

	correct, err := c.service.CheckUser.Execute(context, userID, id)
	if err != nil {
		fmt.Println(err.Error())
		middlewares.Respond(ginc, 500, "Internal server error", nil)
		return
	}

	if !correct {
		middlewares.Respond(ginc, 401, "Unauthorized", nil)
		return
	}

	subdomain, err := c.service.Get.Execute(context, id)
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
