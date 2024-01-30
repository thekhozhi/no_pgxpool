package handler

import (
	"errors"
	"net/http"
	"strconv"
	"develop/api/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)


func (h Handler) CreateBasketProduct(c *gin.Context) {
	createBasketProd := models.BasketProduct{}

	err := c.ShouldBindJSON(&createBasketProd)
	if err != nil{
		handleResponse(c,"Error while reading data body from client!",http.StatusInternalServerError,err)
	}

	pKey, err := h.storage.BasketProduct().Create(createBasketProd)
	if err != nil {
		handleResponse(c, "error is while creating basket product",http.StatusInternalServerError, err)
		return
	}

	id := models.PrimaryKey{ID: pKey}
	res, err := h.storage.BasketProduct().GetByID(id)
	if err != nil {
		handleResponse(c,"error is while getting by id",http.StatusInternalServerError,err)
		return
	}
	handleResponse(w,"",http.StatusCreated, res)
}

func (h Handler) GetBasketProductByID(c *gin.Context) {
	 uid := c.Param("id")

	 basketProd, err := h.storage.BasketProduct().GetByID(models.PrimaryKey{ID: uid})
	 if err != nil{
		handleResponse(c,"error while getting basket product by id!",http.StatusInternalServerError, err)
	 }

	 handleResponse(c,"",http.StatusOK,basketProd)
}

func (h Handler) GetBasketProductList(c *gin.Context) {
	var (
		page, limit int
		err         error
	)

	pageStr := c.DefaultQuery("page", "1")
	page, err = strconv.Atoi(pageStr)
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

	resp, err := h.storage.Basket().GetList(models.GetListRequest{
		Page: page,
		Limit: limit,
	})
		 
	if err != nil {
		handleResponse(c, "error while getting baskets", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, resp)
}

func (h Handler) UpdateBasketProduct(c *gin.Context) {
	 updatebasketProd := models.BasketProduct{}

	 uid := c.Param("id")
	 if uid == ""{
		handleResponse(c, "invalid uuid!", http.StatusBadRequest,errors.New("uuid is not valid"))
		return
	 }

	 updatebasketProd.ID = uid

	 err := c.ShouldBindJSON(&updatebasketProd)
	 if err != nil{
		handleResponse(c, "error while reading body", http.StatusBadRequest,err)
		return
	 }

	 pKey, err := h.storage.BasketProduct().Update(updatebasketProd)
	 if err != nil{
		handleResponse(c, "error while updating basket product!", http.StatusInternalServerError,err)
		return
	 }

	 basketProd, err := h.storage.BasketProduct().GetByID(models.PrimaryKey{
		ID: pKey,
	 })
	 if err != nil{
		handleResponse(c, "error while getting basket product by id", http.StatusInternalServerError,err)
		return
	 }

	 handleResponse(c,"", http.StatusOK,basketProd)
}

func (h Handler) DeleteBasketProduct(c *gin.Context) {
	 uid := c.Param("id")
	 id, err := uuid.Parse(uid)
	 if err != nil{
		handleResponse(c, "uuid is not valid!",http.StatusBadRequest, err)
		return
	 }

	 err = h.storage.BasketProduct().Delete(models.PrimaryKey{
		ID: id.String(),
	 })
	 if err != nil{
		handleResponse(c, "error while deleting basket product by id!", http.StatusInternalServerError,err)
		return
	 }

	 handleResponse(c, "",http.StatusOK,err)
}

