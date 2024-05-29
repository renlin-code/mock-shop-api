package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/renlin-code/mock-shop-api/pkg/domain"
)

// @Summary Get Categories
// @Tags Product Сategories
// @Description Get all product categories.
// @ID get-categories
// @Accept json
// @Produce json
// @Param page query string false "Pagination: page number"
// @Param pageSize query string false "Pagination: amount of items per page"
// @Success 200 {object} response
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /api/categories [get]
func (h *Handler) getAllCategories(c *gin.Context) {
	var params domain.PaginationParams
	if err := c.BindQuery(&params); err != nil {
		Fail(c, bindPaginationParamsErrorText, http.StatusBadRequest)
		return
	}
	if err := params.Validate(); err != nil {
		Fail(c, err.Error(), http.StatusBadRequest)
		return
	}

	categories, err := h.services.Category.GetAll(computePaginationParams(params))
	if err != nil {
		FailAndHandleErr(c, err)
		return
	}

	Response(c, categories)
}

// @Summary Get Category By Id
// @Tags Product Сategories
// @Description Get category by id.
// @ID get-category-by-id
// @Accept json
// @Produce json
// @Param id path int true "Category id"
// @Success 200 {object} response
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /api/categories/{id} [get]
func (h *Handler) getCategoryById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		Fail(c, invalidIdErrorText, http.StatusBadRequest)
		return
	}

	category, err := h.services.Category.GetById(id)
	if err != nil {
		FailAndHandleErr(c, err)
		return
	}

	Response(c, category)
}

// @Summary Get Category Products
// @Tags Products
// @Description Get all products in a category.
// @ID get-category-products
// @Accept json
// @Produce json
// @Param id path int true "Category id"
// @Param page query string false "Pagination: page number"
// @Param pageSize query string false "Pagination: amount of items per page"
// @Success 200 {object} response
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /api/categories/{id}/products [get]
func (h *Handler) getCategoryProducts(c *gin.Context) {
	var params domain.PaginationParams
	if err := c.BindQuery(&params); err != nil {
		Fail(c, bindPaginationParamsErrorText, http.StatusBadRequest)
		return
	}
	if err := params.Validate(); err != nil {
		Fail(c, err.Error(), http.StatusBadRequest)
		return
	}

	catId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		Fail(c, invalidIdErrorText, http.StatusBadRequest)
		return
	}

	limit, offset := computePaginationParams(params)
	products, err := h.services.Category.GetProducts(catId, limit, offset)
	if err != nil {
		FailAndHandleErr(c, err)
		return
	}

	Response(c, products)
}

// @Summary Create Category
// @Security ApiKeyAuth
// @Tags Admin
// @Description Create a new category.
// @ID create-category
// @Accept  multipart/form-data
// @Produce json
// @Param name formData string true "Category name"
// @Param description formData string true "Category description"
// @Param available formData boolean true "Category is available"
// @Param image_file formData file true "Category image"
// @Success 200 {object} response
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admin/categories [post]
func (h *Handler) adminCreateCategory(c *gin.Context) {
	r := c.Request
	var input domain.CreateCategoryInput
	input.Name = r.FormValue("name")
	input.Description = r.FormValue("description")

	available := r.FormValue("available")
	input.Available = available == "true"

	file, handler, err := r.FormFile("image_file")
	if err != nil {
		if err == http.ErrMissingFile {
			Fail(c, "image_file: can not be blank", http.StatusBadRequest)
			return
		}
		FailAndHandleErr(c, err)
		return
	}
	defer file.Close()

	input.ImgFile = handler

	if err := input.Validate(); err != nil {
		Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := h.services.Category.CreateCategory(input, file)
	if err != nil {
		FailAndHandleErr(c, err)
		return
	}

	OKId(c, id)
}

// @Summary Update Category
// @Security ApiKeyAuth
// @Tags Admin
// @Description Update category.
// @ID update-category
// @Accept  multipart/form-data
// @Produce json
// @Param id path int true "Category id"
// @Param name formData string false "Category name"
// @Param description formData string false "Category description"
// @Param available formData boolean false "Category is available"
// @Param image_file formData file false "Category image"
// @Success 200 {object} response
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admin/categories/{id} [put]
func (h *Handler) adminUpdateCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		Fail(c, invalidIdErrorText, http.StatusBadRequest)
		return
	}

	r := c.Request
	var input domain.UpdateCategoryInput

	name := r.FormValue("name")
	if name != "" {
		input.Name = &name
	}

	description := r.FormValue("description")
	if description != "" {
		input.Description = &description
	}

	availableField := r.FormValue("available")
	available := availableField != "false"
	input.Available = &available

	file, handler, err := r.FormFile("image_file")
	if err != nil {
		if err != http.ErrMissingFile {
			Fail(c, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		defer file.Close()
	}

	input.ImgFile = handler

	if err := input.Validate(); err != nil {
		Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.services.Category.UpdateCategory(id, input, file)
	if err != nil {
		FailAndHandleErr(c, err)
		return
	}

	OK(c)
}
