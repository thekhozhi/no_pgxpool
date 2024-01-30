package models

type Category struct {
	ID 	 string`json:"id"`
	Name string`json:"name"`
}

type CreateCategory struct {
	Name string`json:"name"`
}

type UpdateCategory struct {
	ID 	 string`json:"id"`
	Name string`json:"name"`
}

type CategoryResponse struct {
	Categories []Category`json:"categories"`
	Count int			 `json:"count"`
}
