package handler

import (
	"errors"
	"net/http"
	"strconv"
	"develop/api/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

 

var w *gin.Context

func (h Handler) CreateBasket(c *gin.Context) {
	createBasket := models.CreateBasket{}

	err := c.ShouldBindJSON(&createBasket)
	if err != nil{
		handleResponse(c,"Error while reading data body from client!",http.StatusInternalServerError,err)
	}

	pKey, err := h.storage.Basket().Create(createBasket)
	if err != nil {
		handleResponse(c, "error is while creating basket",http.StatusInternalServerError, err)
		return
	}

	id := models.PrimaryKey{ID: pKey}
	res, err := h.storage.Basket().GetByID(id)
	if err != nil {
		handleResponse(c,"error is while getting by id",http.StatusInternalServerError,err)
		return
	}
	handleResponse(w,"",http.StatusCreated, res)
}

func (h Handler) GetBasketByID(c *gin.Context) {
	 uid := c.Param("id")

	 basket, err := h.storage.Basket().GetByID(models.PrimaryKey{ID: uid})
	 if err != nil{
		handleResponse(c,"error while getting basket by id!",http.StatusInternalServerError, err)
	 }

	 handleResponse(c,"",http.StatusOK,basket)
}

func (h Handler) GetBasketList(c *gin.Context) {
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

func (h Handler) UpdateBasket(c *gin.Context) {
	 updatebasket := models.UpdateBasket{}

	 uid := c.Param("id")
	 if uid == ""{
		handleResponse(c, "invalid uuid!", http.StatusBadRequest,errors.New("uuid is not valid"))
		return
	 }

	 updatebasket.ID = uid

	 err := c.ShouldBindJSON(&updatebasket)
	 if err != nil{
		handleResponse(c, "error while reading body", http.StatusBadRequest,err)
		return
	 }

	 pKey, err := h.storage.Basket().Update(updatebasket)
	 if err != nil{
		handleResponse(c, "error while updating basket!", http.StatusInternalServerError,err)
		return
	 }

	 basket, err := h.storage.Basket().GetByID(models.PrimaryKey{
		ID: pKey,
	 })
	 if err != nil{
		handleResponse(c, "erro while getting basl=ket by id", http.StatusInternalServerError,err)
		return
	 }

	 handleResponse(c,"", http.StatusOK,basket)
}

func (h Handler) DeleteBasket(c *gin.Context) {
	 uid := c.Param("id")
	 id, err := uuid.Parse(uid)
	 if err != nil{
		handleResponse(c, "uuid is not valid!",http.StatusBadRequest, err)
		return
	 }

	 err = h.storage.Basket().Delete(models.PrimaryKey{
		ID: id.String(),
	 })
	 if err != nil{
		handleResponse(c, "error while deleting basket by id!", http.StatusInternalServerError,err)
		return
	 }

	 handleResponse(c, "",http.StatusOK,err)
}

