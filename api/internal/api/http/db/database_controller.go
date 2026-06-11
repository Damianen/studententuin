package db

import (
	"api/internal/api/dtos"
	"api/internal/api/middlewares"
	"api/internal/app/db"
	"api/internal/app/subdomain"
	"api/internal/domain"
	"errors"
	"fmt"
	"net/http"

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

func validType(t domain.DatabaseType) bool {
	switch t {
	case domain.DatabaseTypePostgres, domain.DatabaseTypeMySQL, domain.DatabaseTypeMongoDB:
		return true
	default:
		return false
	}
}

func (c *Controller) Create(ginc *gin.Context) {
	var req dtos.CreateDatabaseRequest
	err := ginc.ShouldBindBodyWithJSON(&req)
	if err != nil || !validType(domain.DatabaseType(req.Type)) {
		middlewares.Respond(ginc, http.StatusBadRequest, "Invalid JSON or missing value", nil)
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
		Status:      domain.DatabaseStatusProvisioning,
		Version:     req.Version,
	}

	err = c.dbService.Create.Execute(ginc.Request.Context(), dbInput)
	if err != nil {
		fmt.Println(err.Error())
		middlewares.Respond(ginc, http.StatusInternalServerError, "failed to create subdomain", nil)
		return
	}

	middlewares.Respond(ginc, http.StatusCreated, "success", nil)
}

func (c *Controller) Update(ginc *gin.Context) {
	var req dtos.UpdateDatabaseRequest
	if !middlewares.BindJSON(ginc, &req) {
		return
	}

	if req.Type != nil {
		if !validType(domain.DatabaseType(*req.Type)) {
			middlewares.Respond(ginc, http.StatusBadRequest, "invalid JSON or missing value", nil)
			return
		}
	}

	userID := ginc.GetString("userID")
	subdomainId := ginc.Param("id")
	dbId := ginc.Param("dbId")

	if !middlewares.CheckOwnership(ginc, c.subdomainService.CheckUser.Execute, userID, subdomainId, "subdomain") {
		return
	}

	databaseInput := db.DatabaseUpdateInput{
		ID:      dbId,
		Name:    req.Name,
		Type:    (*domain.DatabaseType)(req.Type),
		Version: req.Version,
	}

	err := c.dbService.Update.Execute(ginc.Request.Context(), databaseInput)
	if errors.Is(err, gorm.ErrRecordNotFound) {
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

	err := c.dbService.Delete.Execute(ginc.Request.Context(), id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		middlewares.Respond(ginc, http.StatusNotFound, "database not found", nil)
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

	database, err := c.dbService.Get.Execute(ginc.Request.Context(), id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
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
