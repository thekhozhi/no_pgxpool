package postgres

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"develop/api/models"
	"develop/storage"
)

type basProductRepo struct {
	DB *sql.DB
}

func NewBasProductRepo(db *sql.DB) storage.IBasketProduct {
	return basProductRepo{DB: db}
}

func (b basProductRepo) Create(basProd models.BasketProduct) (string, error) {

	id := uuid.New()

	 _, err := b.DB.Exec(`insert into basket_products(id, basket_id, product_id, quantity)
	values($1, $2, $3, $4)`, id, basProd.BasketID, basProd.ProductID, basProd.Quantity)
	   if err != nil {
		fmt.Println("error is while inserting data", err.Error())
		return "", err
	}
	return id.String(), nil
	
}

func (b basProductRepo) GetByID(key models.PrimaryKey) (models.BasketProduct, error) {
	bProd := models.BasketProduct{}

	 err := b.DB.QueryRow(`select id, basket_id, product_id, quantity from basket_products where id = $1`, key.ID).Scan(
		 &bProd.ID,
		 &bProd.BasketID,
		 &bProd.ProductID,
		 &bProd.Quantity,
	 )
		if err != nil {
		fmt.Println("error is while selecting basket_product", err.Error())
		return models.BasketProduct{}, err
	}

	return bProd, nil
}

func (b basProductRepo) GetList(request models.GetListRequest) (models.BasketProductResponse, error) {
	var (
		bProducts           = []models.BasketProduct{}
		count             = 0
		countQuery, query string
		page              = request.Page
		offset            = (page - 1) * request.Limit
	)

	countQuery = `
		SELECT count(1) from basket_products`


	if err := b.DB.QueryRow(countQuery).Scan(&count); err != nil {
		fmt.Println("error while scanning count of baskets", err.Error())
		return models.BasketProductResponse{}, err
	}

	query = `
		SELECT  id, basket_id, product_id, quantity from basket_products`

	query += ` LIMIT $1 OFFSET $2`

	rows, err := b.DB.Query(query, request.Limit, offset)
	if err != nil {
		fmt.Println("error while query rows", err.Error())
		return models.BasketProductResponse{}, err
	}

	for rows.Next() {
		 bProd := models.BasketProduct{}

		err := rows.Scan(
			&bProd.ID,
			&bProd.BasketID,
			&bProd.ProductID,
			&bProd.Quantity,
		)
		if err != nil {
			fmt.Println("error while scanning row", err.Error())
			return models.BasketProductResponse{}, err
		}

		bProducts = append(bProducts, bProd)
	}

	return models.BasketProductResponse{
		BasketProducts: bProducts,
		Count: count,
	}, nil
}

func (b basProductRepo) Update(bProd models.BasketProduct) (string, error) {
	 _, err := b.DB.Exec(`update basket_products set basket_id = $1, product_id = $2, quantity = $3 where id = $4`, 
	  bProd.BasketID, bProd.ProductID, bProd.Quantity, bProd.ID); 
	 if err != nil {
		return "", err
	}
	return bProd.ID, err
}

func (b basProductRepo) Delete(key models.PrimaryKey) error {
	if _, err := b.DB.Exec(`delete from basket_products where id = $1`, key.ID); err != nil {
		return err
	}
	return nil
}

func (b basProductRepo) AddProducts(basketID string, products map[string]int) error {
	query := `
		INSERT INTO basket_products
			(id, basket_id, product_id, quantity)
			values ($1, $2, $3, $4)`

	stmt, err := b.DB.Prepare(query)
	if err != nil{
		fmt.Println("Error while preparing statements!",err.Error())
		return err
	}

	for productID, quantity := range products{
		 _, err := stmt.Exec(query, uuid.New(), basketID, productID, quantity)
		 if err != nil{
			fmt.Println("Error while adding product to basket_products table", err.Error())
			return err
		}
	}
	return nil
}