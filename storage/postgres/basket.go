package postgres

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"develop/api/models"
	"develop/storage"
)

type basketRepo struct {
	DB *sql.DB
}

func NewBasketRepo(db *sql.DB) storage.IBasket {
	return basketRepo{DB: db}
}

func (b basketRepo) Create(basket models.CreateBasket) (string, error) {

	id := uuid.New()

	 _, err := b.DB.Exec(`insert into baskets(id, customer_id, total_sum)
	values($1, $2, $3)`, id, basket.CustomerID, basket.TotalSum);
	   if err != nil {
		fmt.Println("error is while inserting data", err.Error())
		return "", err
	}
	return id.String(), nil
	
}

func (b basketRepo) GetByID(key models.PrimaryKey) (models.Basket, error) {
	basket := models.Basket{}

	if err := b.DB.QueryRow(`select id, customer_id, total_sum from baskets where id = $1`, key.ID).Scan(&basket.ID,
		&basket.CustomerID, &basket.TotalSum); err != nil {
		fmt.Println("error is while selecting basket", err.Error())
		return models.Basket{}, err
	}

	return basket, nil
}

func (b basketRepo) GetList(request models.GetListRequest) (models.BasketResponse, error) {
	var (
		baskets            = []models.Basket{}
		count             = 0
		countQuery, query string
		page              = request.Page
		offset            = (page - 1) * request.Limit
	)

	countQuery = `
		SELECT count(1) from baskets`


	if err := b.DB.QueryRow(countQuery).Scan(&count); err != nil {
		fmt.Println("error while scanning count of baskets", err.Error())
		return models.BasketResponse{}, err
	}

	query = `
		SELECT  id, customer_id, total_sum from baskets`

	query += ` LIMIT $1 OFFSET $2`

	rows, err := b.DB.Query(query, request.Limit, offset)
	if err != nil {
		fmt.Println("error while query rows", err.Error())
		return models.BasketResponse{}, err
	}

	for rows.Next() {
		basket := models.Basket{}

		err := rows.Scan(
			 &basket.ID,
			 &basket.CustomerID,
			 &basket.TotalSum,
		)
		if err != nil {
			fmt.Println("error while scanning row", err.Error())
			return models.BasketResponse{}, err
		}

		 baskets = append(baskets, basket)
	}

	return models.BasketResponse{
		Baskets: baskets,
		Count: count,
	}, nil
}

func (b basketRepo) Update(basket models.UpdateBasket) (string, error) {
	 _, err := b.DB.Exec(`update baskets set customer_id = $1, total_sum = $2 where id = $3`, &basket.CustomerID, &basket.TotalSum, &basket.ID); 
	 if err != nil {
		return "", err
	}
	return basket.ID, err
}

func (b basketRepo) Delete(key models.PrimaryKey) error {
	if _, err := b.DB.Exec(`delete from baskets where id = $1`, key.ID); err != nil {
		return err
	}
	return nil
}
