package dtos

type CreateSubdomainRequest struct {
	Name       string `json:"name" binding:"required"`
	FullDomain string `json:"fullDomain" binding:"required"`
}

type UpdateSubdomainRequest struct {
	Name       *string `json:"name,omitempty" binding:"omitempty"`
	FullDomain *string `json:"fullDomain,omitempty" binding:"omitempty"`
}

type SubdomainResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	FullDomain string `json:"fullDomain"`
	IsActive   bool   `json:"isActive"`
}

type SubdomainListItemResponse struct {
	ID          string                   `json:"id"`
	Name        string                   `json:"name"`
	FullDomain  string                   `json:"fullDomain"`
	IsActive    bool                     `json:"isActive"`
	Database    *DatabaseListResponse    `json:"database,omitempty"`
	Application *ApplicationListResponse `json:"application,omitempty"`
}
