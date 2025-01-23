package category

type UpdateRequest struct {
	Name  *string `json:"name" binding:"omitempty,max=32"`
	Color *string `json:"color" binding:"omitempty,max=7"`
}

type UpdateResponse struct {
	Message string `json:"message"`
}
