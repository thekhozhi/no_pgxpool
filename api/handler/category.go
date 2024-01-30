package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"develop/api/models"
)

func (h Handler) CreateCategory(c *gin.Context) {
	createCategory := models.Category{}

	if err := c.ShouldBindJSON(&createCategory); err != nil {
		handleResponse(c, "error while reading body from client", http.StatusBadRequest, err)
		return
	}

	pKey, err := h.storage.Category().Create(createCategory)
	if err != nil {
		handleResponse(c, "error while creating category", http.StatusInternalServerError, err)
		return
	}

	category, err := h.storage.Category().GetByID(models.PrimaryKey{
		ID: pKey,
	})
	if err != nil {
		handleResponse(c, "error while getting category by id", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusCreated, category)
}

func (h Handler) GetCategory(c *gin.Context) {
	var err error

	uid := c.Param("id")

	category, err := h.storage.Category().GetByID(models.PrimaryKey{
		ID: uid,
	})
	if err != nil {
		handleResponse(c, "error while getting category by id", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, category)
}

func (h Handler) GetCategoryList(c *gin.Context) {
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

func (h Handler) UpdateCategory(c *gin.Context) {
	updateCategory := models.Category{}

	uid := c.Param("id")
	if uid == "" {
		handleResponse(c, "invalid uuid", http.StatusBadRequest, errors.New("uuid is not valid"))
		return
	}

	updateCategory.ID = uid

	if err := c.ShouldBindJSON(&updateCategory); err != nil {
		handleResponse(c, "error while reading body", http.StatusBadRequest, err.Error())
		return
	}

	pKey, err := h.storage.Category().Update(updateCategory)
	if err != nil {
		handleResponse(c, "error while updating category", http.StatusInternalServerError, err.Error())
		return
	}

	category, err := h.storage.Category().GetByID(models.PrimaryKey{
		ID: pKey,
	})
	if err != nil {
		handleResponse(c, "error while getting category by id", http.StatusInternalServerError, err)
		return
	}

	handleResponse(c, "", http.StatusOK, category)
}

func (h Handler) DeleteCategory(c *gin.Context) {
	uid := c.Param("id")
	id, err := uuid.Parse(uid)
	if err != nil {
		handleResponse(c, "uuid is not valid", http.StatusBadRequest, err.Error())
		return
	}

	if err = h.storage.Category().Delete(models.PrimaryKey{
		ID: id.String(),
	}); err != nil {
		handleResponse(c, "error while deleting category by id", http.StatusInternalServerError, err.Error())
		return
	}

	handleResponse(c, "", http.StatusOK, "data successfully deleted")
}

 