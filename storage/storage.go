package storage

import (
	"develop/api/models"
)

type IStorage interface {
	Close()
	User() IUserStorage
	Basket() IBasket
	Product() IProduct
	Category() ICategory
	BasketProduct() IBasketProduct
}


type IUserStorage interface {
	Create(models.CreateUser) (string, error)
	GetByID(models.PrimaryKey) (models.User, error)
	GetList(models.GetListRequest) (models.UsersResponse, error)
	Update(models.UpdateUser) (string, error)
	Delete(models.PrimaryKey) error
	GetPassword(id string) (string, error)
	UpdatePassword(models.UpdateUserPassword) error
}

type IBasket interface {
	Create(models.CreateBasket) (string, error)
	GetByID(models.PrimaryKey) (models.Basket, error)
	GetList(models.GetListRequest) (models.BasketResponse, error)
	Update(models.UpdateBasket) (string, error)
	Delete(models.PrimaryKey) error
}

type ICategory interface {
	Create(models.Category)(string,error)
	GetByID(models.PrimaryKey)(models.Category,error)
	GetList(models.GetListRequest)(models.CategoryResponse,error)
	Update(models.Category)(string, error)
	Delete(models.PrimaryKey)(error)
}

type IProduct interface {
	Create(models.Product)(string,error)
	GetByID(models.PrimaryKey)(models.Product,error)
	GetList(models.GetListRequest)(models.ProductResponse,error)
	Update(models.Product)(string, error)
	Delete(models.PrimaryKey)(error)
	Search(map[string]int)(map[string]int, map[string]int, error)
	TakeProducts(map[string]int) error
}

type IBasketProduct interface {
	Create(models.BasketProduct)(string,error)
	GetByID(models.PrimaryKey)(models.BasketProduct,error)
	GetList(models.GetListRequest)(models.BasketProductResponse,error)
	Update(models.BasketProduct)(string, error)
	Delete(models.PrimaryKey)(error)
	AddProducts(string, map[string]int) error
}
 