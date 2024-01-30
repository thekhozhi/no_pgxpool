package postgres

import (
	"database/sql"
	"fmt"
	"develop/api/models"
	"develop/storage"
	"github.com/google/uuid"
)

type categoryRepo struct {
	DB *sql.DB
}

func NewCategoryRepo(db *sql.DB)storage.ICategory{
	return categoryRepo{DB: db}
}

func (c categoryRepo) Create(cat models.Category) (string, error) {

	id := uuid.New()

	 _, err := c.DB.Exec(`INSERT INTO categories(id, name)
	values($1, $2)`, id, cat.Name);
	   if err != nil {
		fmt.Println("error is while inserting data", err.Error())
		return "", err
	}
	return id.String(), nil
	
}

func  (c categoryRepo) GetByID(key models.PrimaryKey) (models.Category, error) {
	 cat := models.Category{}

		err := c.DB.QueryRow(`select id,name from categories where id = $1`, key.ID).Scan(
			&cat.ID,
			&cat.Name,
		)
		 
		if err != nil {
		fmt.Println("error is while selecting category", err.Error())
		return models.Category{}, err
	}

	return cat, nil
}

func (c categoryRepo) GetList(request models.GetListRequest) (models.CategoryResponse, error) {
	var (
		categories            = []models.Category{}
		count             = 0
		countQuery, query string
		page              = request.Page
		offset            = (page - 1) * request.Limit
	)

	countQuery = `
		SELECT count(1) from categories`


	if err := c.DB.QueryRow(countQuery).Scan(&count); err != nil {
		fmt.Println("error while scanning count of categories", err.Error())
		return models.CategoryResponse{}, err
	}

	query = `
		SELECT  id, name from categories`

	query += ` LIMIT $1 OFFSET $2`

	rows, err := c.DB.Query(query, request.Limit, offset)
	if err != nil {
		fmt.Println("error while query rows", err.Error())
		return models.CategoryResponse{}, err
	}

	for rows.Next() {
		cat := models.Category{}

		err := rows.Scan(
			&cat.ID,
			&cat.Name,
		)
		if err != nil {
			fmt.Println("error while scanning row", err.Error())
			return models.CategoryResponse{}, err
		}

		categories = append(categories, cat)
	}

	return models.CategoryResponse{
		Categories: categories,
		Count: count,
	}, nil
}

func (c categoryRepo) Update(cat models.Category) (string, error) {

	 _, err := c.DB.Exec(`UPDATE categories SET name = $1 where id = $2`,&cat.Name, &cat.ID); 
	 if err != nil {
		return "", err
	}
	return cat.ID, err
}

func (c categoryRepo) Delete(key models.PrimaryKey) error {
	if _, err := c.DB.Exec(`DELETE from categories WHERE id = $1`, key.ID); err != nil {
		return err
	}
	return nil
}
