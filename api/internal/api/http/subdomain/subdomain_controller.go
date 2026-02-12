package subdomain

import (
	"api/internal/api/dtos"
	"api/internal/api/middlewares"
	"api/internal/app/subdomain"
	"api/internal/infra/postgres"
	"errors"
	"fmt"
	"net/http"

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
	var req dtos.CreateSubdomainRequest
	if !middlewares.BindJSON(ginc, &req) {
		return
	}

	userID := ginc.GetString("userID")
	subdomainInput := subdomain.SubdomainInput{
		Name:       req.Name,
		FullDomain: req.FullDomain,
		UserID:     userID,
	}

	err := c.service.Create.Execute(ginc.Request.Context(), subdomainInput)
	if errors.Is(err, postgres.ErrFullDomainAlreadyInUse) {
		middlewares.Respond(ginc, http.StatusConflict, "domain already in use", nil)
		return
	}
	if err != nil {
		fmt.Println(err.Error())
		middlewares.Respond(ginc, http.StatusInternalServerError, "failed to create subdomain", nil)
		return
	}

	middlewares.Respond(ginc, http.StatusCreated, "success", nil)
}

func (c *Controller) Update(ginc *gin.Context) {
	var req dtos.UpdateSubdomainRequest
	if !middlewares.BindJSON(ginc, &req) {
		return
	}

	id := ginc.Param("id")
	userID := ginc.GetString("userID")

	if !middlewares.CheckOwnership(ginc, c.service.CheckUser.Execute, userID, id, "subdomain") {
		return
	}

	subdomainInput := subdomain.SubdomainUpdateInput{
		ID:         id,
		Name:       req.Name,
		FullDomain: req.FullDomain,
	}

	err := c.service.Update.Execute(ginc.Request.Context(), subdomainInput)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		middlewares.Respond(ginc, http.StatusNotFound, "subdomain not found", nil)
		return
	}
	if err != nil {
		fmt.Println(err.Error())
		middlewares.Respond(ginc, http.StatusInternalServerError, "failed to update subdomain", nil)
		return
	}

	middlewares.Respond(ginc, http.StatusOK, "success", nil)
}

func (c *Controller) Delete(ginc *gin.Context) {
	id := ginc.Param("id")
	userID := ginc.GetString("userID")

	if !middlewares.CheckOwnership(ginc, c.service.CheckUser.Execute, userID, id, "subdomain") {
		return
	}

	err := c.service.Delete.Execute(ginc.Request.Context(), id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		middlewares.Respond(ginc, http.StatusNotFound, "subdomain not found", nil)
		return
	}
	if err != nil {
		fmt.Println(err.Error())
		middlewares.Respond(ginc, http.StatusInternalServerError, "failed to delete subdomain", nil)
		return
	}

	middlewares.Respond(ginc, http.StatusOK, "success", nil)
}

func (c *Controller) Get(ginc *gin.Context) {
	id := ginc.Param("id")
	userID := ginc.GetString("userID")

	if !middlewares.CheckOwnership(ginc, c.service.CheckUser.Execute, userID, id, "subdomain") {
		return
	}

	subdomain, err := c.service.Get.Execute(ginc.Request.Context(), id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		middlewares.Respond(ginc, http.StatusNotFound, "subdomain not found", nil)
		return
	}
	if err != nil {
		fmt.Println(err.Error())
		middlewares.Respond(ginc, http.StatusInternalServerError, "failed to get subdomain", nil)
		return
	}

	subdomainResponse := dtos.SubdomainResponse{
		ID:         subdomain.ID.String(),
		Name:       subdomain.Name,
		FullDomain: subdomain.FullDomain,
		IsActive:   subdomain.IsActive,
	}

	middlewares.Respond(ginc, http.StatusOK, "success", subdomainResponse)
}

func (c *Controller) GetAll(ginc *gin.Context) {
	userID := ginc.GetString("userID")

	subdomains, err := c.service.GetAll.Execute(ginc.Request.Context(), userID)
	if err != nil {
		fmt.Println(err.Error())
		middlewares.Respond(ginc, http.StatusInternalServerError, "failed to get subdomains", nil)
		return
	}

	var response []dtos.SubdomainListItemResponse
	for _, s := range subdomains {
		item := dtos.SubdomainListItemResponse{
			ID:         s.ID.String(),
			Name:       s.Name,
			FullDomain: s.FullDomain,
			IsActive:   s.IsActive,
		}

		if s.Database != nil {
			item.Database = &dtos.DatabaseListResponse{
				ID:      s.Database.ID.String(),
				Name:    s.Database.Name,
				Type:    string(s.Database.Type),
				Version: s.Database.Version,
				Status:  string(s.Database.Status),
			}
		}

		if s.Application != nil {
			item.Application = &dtos.ApplicationListResponse{
				ID:      s.Application.ID.String(),
				Name:    *s.Application.Name,
				Type:    string(s.Application.Type),
				Status:  string(s.Application.Status),
				RepoUrl: *s.Application.RepositoryURL,
				Branch:  *s.Application.Branch,
			}
		}

		response = append(response, item)
	}

	middlewares.Respond(ginc, http.StatusOK, "success", response)
}
