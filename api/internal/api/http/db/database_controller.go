package db

import (
	"api/internal/api/dtos"
	"api/internal/api/middlewares"
	"api/internal/app/db"
	"api/internal/app/subdomain"
	"api/internal/domain"
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Controller struct {
	dbService *db.Service
	subdomainService *subdomain.Service
}

func NewController(d db.Dependencies, sd subdomain.Dependencies) *Controller {
	return &Controller{
		dbService: db.NewService(d),
		subdomainService: subdomain.NewService(sd),
	}
}

func validType(t domain.DatabaseType) bool {
    switch t {
    case domain.DatabaseTypePostgres, domain.DatabaseTypeMySQL, domain.DatabaseTypeMongoDB:
        return true
    default:
        return false
    }
}

func (c *Controller) Create(ginc *gin.Context) {
	context := ginc.Request.Context()

	var req dtos.CreateDatabaseRequest
	err := ginc.ShouldBindBodyWithJSON(&req)
	if err != nil || validType(domain.DatabaseType(req.Type)) {
		middlewares.Respond(ginc, 400, "Invalid JSON or missing value", nil)
		return
	}

	userID := ginc.GetString("userID")
	subdomainId := ginc.Param("id")

	allowed, err := c.subdomainService.CheckUser.Execute(context, userID, subdomainId)
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
		middlewares.Respond(ginc, 403, "unauthorized", nil)
		return
	}

	dbInput := db.DatabaseInput{
		SubdomainID: subdomainId,
		Name: req.Name,
		Type: domain.DatabaseType(req.Type),
		Status: domain.DatabaseStatusProvisioning,
		Version: req.Version,
	}

	err = c.dbService.Create.Execute(context, dbInput)
	if err != nil {
		fmt.Println(err.Error())
		middlewares.Respond(ginc, 500, "failed to create subdomain", nil)
		return
	}

	middlewares.Respond(ginc, 201, "success", nil)
}

func (c *Controller) Update(ginc *gin.Context) {
	context := ginc.Request.Context()

	var req dtos.UpdateDatabaseRequest
	err := ginc.ShouldBindBodyWithJSON(&req)
	if err != nil {
		fmt.Println(err.Error())
		middlewares.Respond(ginc, 400, "invalid JSON or missing value", nil)
		return
	}

	if req.Type != nil {
		if validType(domain.DatabaseType(*req.Type)) {
			middlewares.Respond(ginc, 400, "invalid JSON or missing value", nil)
			return
		}
	}

	userID := ginc.GetString("userID")
	subdomainId := ginc.Param("id")
	dbId := ginc.Param("dbId")

	allowed, err := c.subdomainService.CheckUser.Execute(context, userID, subdomainId)
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

	databaseInput := db.DatabaseUpdateInput{
		ID: dbId,
		Name: req.Name,
		Type: (*domain.DatabaseType)(req.Type),
		Version: req.Version,
	}

	err = c.dbService.Update.Execute(context, databaseInput)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		middlewares.Respond(ginc, 404, "database not found", nil)
		return
	}
	if err != nil {
		fmt.Println(err.Error())
		middlewares.Respond(ginc, 500, "failed to update the database", nil)
		return
	}

	ginc.Status(204)
}

func (c *Controller) Delete(ginc *gin.Context) {
	context := ginc.Request.Context()

	userID := ginc.GetString("userID")
	subdomainId := ginc.Param("id")
	id := ginc.Param("dbId")

	allowed, err := c.subdomainService.CheckUser.Execute(context, userID, subdomainId)
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

	err = c.dbService.Delete.Execute(context, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		middlewares.Respond(ginc, 404, "database not found", nil)
		return
	}
	if err != nil {
		fmt.Println(err.Error())
		middlewares.Respond(ginc, 500, "failed to delete the database", nil)
		return
	}

	ginc.Status(204)
}

func (c *Controller) Get(ginc *gin.Context) {
	context := ginc.Request.Context()

	userID := ginc.GetString("userID")
	subdomainId := ginc.Param("id")
	id := ginc.Param("dbId")

	allowed, err := c.subdomainService.CheckUser.Execute(context, userID, subdomainId)
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

	database, err := c.dbService.Get.Execute(context, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		middlewares.Respond(ginc, 404, "database not found", nil)
		return
	}

	if err != nil {
		fmt.Println(err.Error())
		middlewares.Respond(ginc, 500, "failed to get database", nil)
	}

	databaseResponse := dtos.DatabaseListResponse{
		ID: database.ID.String(),
		Name: database.Name,
		Type: string(database.Type),
		Version: database.Version,
		Status: string(database.Status),
	}

	middlewares.Respond(ginc, 200, "success", databaseResponse)
}
