package handler

import (
	"develop/api/models"
	"errors"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h Handler) CreateProduct(c *gin.Context) {
	createProd := models.Product{}

	err := c.ShouldBindJSON(&createProd)
	if err != nil {
		handleResponse(c, "error while reading body from client", http.StatusBadRequest, err)
		return
	}

	pKey, err := h.storage.Product().Create(createProd)
	if err != nil {
		handleResponse(c, "error while creating user", http.StatusInternalServerError, err)
		return
	}

	product, err := h.storage.Product().GetByID(models.PrimaryKey{
		ID: pKey,
	})
	if err != nil {
		handleResponse(c, "error while getting user by id", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusCreated, product)
}

func (h Handler) GetProduct(c *gin.Context) {
	var err error

	uid := c.Param("id")

	product, err := h.storage.Product().GetByID(models.PrimaryKey{
		ID: uid,
	})
	if err != nil {
		handleResponse(c, "error while getting product by id", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, product)
}

func (h Handler) GetProductList(c *gin.Context) {
	var (
		page, limit int
	)

	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		handleResponse(c, "error while parsing page", http.StatusBadRequest, err.Error())
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err = strconv.Atoi(limitStr)
	if err != nil {
		handleResponse(c, "error while parsing limit", http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.storage.User().GetList(models.GetListRequest{
		Page:   page,
		Limit:  limit,
	})
	if err != nil {
		handleResponse(c, "error while getting users", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, resp)
}

func (h Handler) UpdateProduct(c *gin.Context) {
	updateProd := models.Product{}

	uid := c.Param("id")
	if uid == "" {
		handleResponse(c, "invalid uuid", http.StatusBadRequest, errors.New("uuid is not valid"))
		return
	}

	updateProd.ID = uid

	if err := c.ShouldBindJSON(&updateProd); err != nil {
		handleResponse(c, "error while reading body", http.StatusBadRequest, err.Error())
		return
	}

	pKey, err := h.storage.Product().Update(updateProd)
	if err != nil {
		handleResponse(c, "error while updating user", http.StatusInternalServerError, err.Error())
		return
	}

	product, err := h.storage.Product().GetByID(models.PrimaryKey{
		ID: pKey,
	})
	if err != nil {
		handleResponse(c, "error while getting product by id", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, product)
}

func (h Handler) DeleteProduct(c *gin.Context) {
	uid := c.Param("id")
	id, err := uuid.Parse(uid)
	if err != nil {
		handleResponse(c, "uuid is not valid", http.StatusBadRequest, err.Error())
		return
	}

	if err = h.storage.Product().Delete(models.PrimaryKey{
		ID: id.String(),
	}); err != nil {
		handleResponse(c, "error while deleting product by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, "data successfully deleted")
}

func (h Handler) StartSellNew(c * gin.Context){
	request := models.SellRequest{}

	err := c.ShouldBindJSON(&request)
	if err != nil{
		handleResponse(c, "error while reading from json!", http.StatusBadRequest, err.Error())
	}

	selectedProducts, productPrices, err := h.storage.Product().Search(request.Products)
	if err != nil{
		handleResponse(c, "error while searching products", http.StatusInternalServerError, err.Error())
		return
	}

	basket, err := h.storage.Basket().GetByID(models.PrimaryKey{
		ID: request.BasketID,
	})

	if err != nil{
		handleResponse(c, "error while searching products", http.StatusInternalServerError, err.Error())
		return
	}

	customer, err := h.storage.User().GetByID(models.PrimaryKey{
		ID: basket.CustomerID,
	})

	if err != nil{
		handleResponse(c, "error while getting customer by id", http.StatusInternalServerError,err.Error())
		return
	}

	totalSum, profit := 0, 0
	basketProducts := map[string]int{}

	for productID, price := range selectedProducts{
		customerQuantity := request.Products[productID]
		totalSum += price * customerQuantity

		profit += customerQuantity * (price - productPrices[productID])
		basketProducts[productID] = customerQuantity
	}

	if customer.Cash < uint(totalSum){
		handleResponse(c, "Not enough cash!", http.StatusBadRequest, err.Error())
		return
	}

	err = h.storage.Product().TakeProducts(basketProducts)
	if err != nil{
		handleResponse(c, "error while minusing products from database!", http.StatusInternalServerError,err.Error())
		return
	}

	err = h.storage.BasketProduct().AddProducts(basket.ID, basketProducts)
	if err != nil{
		handleResponse(c, "error while adding products!",http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "Succesfully finished the purchase!",http.StatusOK, profit)
}