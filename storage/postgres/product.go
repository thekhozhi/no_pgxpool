package postgres

import (
	"database/sql"
	"develop/api/models"
	"develop/storage"
	"fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type productRepo struct {
	DB *sql.DB
}

func NewProductRepo(db *sql.DB) storage.IProduct {
	return productRepo{DB: db}
}

func (p productRepo) Create(product models.Product) (string, error) {

	id := uuid.New()

	 _, err := p.DB.Exec(`INSERT INTO products(id ,name, price, original_price, quantity, category_id)
	values($1, $2, $3, $4, $5, $6)`, 
	id, 
	product.Name, 
	product.Price, 
	product.OriginalPrice,
	product.Quantity, 
	product.CategoryID);
	   if err != nil {
		fmt.Println("error is while inserting data", err.Error())
		return "", err
	}
	return id.String(), nil
	
}

func (p productRepo) GetByID(key models.PrimaryKey) (models.Product, error) {
	prod := models.Product{}

	 err := p.DB.QueryRow(`SELECT id, name, price, original_price, quantity, category_id from products where id = $1`, key.ID).Scan(
		 &prod.ID,
		 &prod.Name,
		 &prod.Price,
		 &prod.OriginalPrice,
		 &prod.Quantity,
		 &prod.CategoryID,
	 )
		if err != nil {
		fmt.Println("error is while selecting product", err.Error())
		return models.Product{}, err
	}

	return prod, nil
}

func (p productRepo) GetList(request models.GetListRequest) (models.ProductResponse, error) {
	var (
		products            = []models.Product{}
		count             = 0
		countQuery, query string
		page              = request.Page
		offset            = (page - 1) * request.Limit
	)

	countQuery = `
		SELECT count(1) from products`


	err := p.DB.QueryRow(countQuery).Scan(&count) 
		if err != nil {
		fmt.Println("error while scanning count of products", err.Error())
		return models.ProductResponse{}, err
	}

	query = `
		SELECT  id, name, price, original_price, quantity, category_id from products`

	query += ` LIMIT $1 OFFSET $2`

	rows, err := p.DB.Query(query, request.Limit, offset)
	if err != nil {
		fmt.Println("error while query rows", err.Error())
		return models.ProductResponse{}, err
	}

	for rows.Next() {
		 prod := models.Product{}

		err := rows.Scan(
			&prod.ID,
			&prod.Name,
			&prod.Price,
			&prod.OriginalPrice,
			&prod.Quantity,
			&prod.CategoryID,
		)
		if err != nil {
			fmt.Println("error while scanning row", err.Error())
			return models.ProductResponse{}, err
		}

		  products = append(products, prod)
	}

	return models.ProductResponse{
		Products: products,
		Count: count,
	}, nil
}

func (p productRepo) Update(prod models.Product) (string, error) {
	 _, err := p.DB.Exec(`UPDATE products SET name = $1, price = $2, original_price = $3, quantity = 
	 $4, category_id = $5 WHERE id = $6`, 
	 prod.Name, 
	 prod.Price,
	 prod.OriginalPrice,
	 prod.Quantity,
	 prod.CategoryID,
	 prod.ID); 
	 if err != nil {
		return "", err
	}
	return prod.ID, err
}

func (p productRepo) Delete(key models.PrimaryKey) error {
	_, err := p.DB.Exec(`DELETE from products WHERE id = $1`, key.ID)
	if err != nil {
		return err
	}
	return nil
}


func (p productRepo) Search(customerProductIDs map[string]int) (map[string]int, map[string]int, error) {
	var (
		selectedProducts = models.SellRequest{
			Products: map[string]int{},
		}
		arrayProductIDs = make([]string, len(customerProductIDs))
		productPrices = make(map[string]int, 0)
	)

	for key := range customerProductIDs{
		arrayProductIDs = append(arrayProductIDs, key)
	}

	query := `
	SELECT id, quantity, price, original_price from products where id::varchar = ANY($1)`

	rows, err := p.DB.Query(query, pq.Array(arrayProductIDs))
		if err != nil{
			fmt.Println("error while getting products by product ids!", err.Error())
			return nil,nil, err
		}

	for rows.Next(){
		var(
			quantity,price,originalPrice int
			productID 				  string
		)
		
		err := rows.Scan(
			&productID,
			&quantity,
			&price,
			&originalPrice,
		)
		if err != nil{
			fmt.Println("Error while scanning rows one by one", err.Error())
			return nil, nil, err
		}

		if customerProductIDs[productID] <= quantity{
			selectedProducts.Products[productID] = price
			productPrices[productID] = originalPrice
		}
	}
	return selectedProducts.Products, productPrices, nil
}

func (p productRepo) TakeProducts(products map[string]int) error {
	query := `
	UPDATE products SET quantity = quantity - $1 where id = $2`

	stmt, err := p.DB.Prepare(query)
	if err != nil{
		fmt.Println("Error while preparing query",err.Error())
		return err
	}
	
	for productID, quantity := range products{
		_, err := stmt.Exec(query, quantity, productID)
		if err != nil{
			fmt.Println("Error while updating product quantity",err.Error())
			return err
		}
	}
	return nil
}