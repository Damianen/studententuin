package dtos

type CreateApplicationRequest struct {
	Name string `json:"name" binding:"required"`
	Type string `json:"type" binding:"required"`
	RepoUrl string `json:"repo_url" binding:"required"`
	Branch string `json:"branch" binding:"required"`
	BuildCommand string `json:"build_command" binding:"required"`
	StartCommand string `json:"start_command" binding:"required"`
}

type UpdateApplicationRequest struct {
	Name *string `json:"name,omitempty" binding:"omitempty"`
	Type *string `json:"type,omitempty" binding:"omitempty"`
	RepoUrl *string `json:"repo_url,omitempty" binding:"omitempty"`
	Branch *string `json:"branch,omitempty" binding:"omitempty"`
	BuildCommand *string `json:"build_command,omitempty" binding:"omitempty"`
	StartCommand *string `json:"start_command,omitempty" binding:"omitempty"`
}

type ApplicationListResponse struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Status string `json:"status"`
	RepoUrl string `json:"repo_url"`
	Branch string `json:"branch"`
}
