package dtos

// CreateDatabaseRequest deliberately has no password field: credentials are
// generated server-side (§6 phase 5) and surface only via the connection
// string on the record.
type CreateDatabaseRequest struct {
	Name string `json:"name" binding:"required"`
	Type string `json:"type" binding:"required"`
	Version string `json:"version" binding:"required"`
	DbName string `json:"db_name" binding:"required"`
}

type UpdateDatabaseRequest struct {
	Name *string `json:"name,omitempty" binding:"omitempty"`
	Type *string `json:"type,omitempty" binding:"omitempty"`
	Version *string `json:"version,omitempty" binding:"omitempty"`
}

type DatabaseListResponse struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Version string `json:"version"`
	Status string `json:"status"`
	ConnectionString string `json:"connection_string"`
}
