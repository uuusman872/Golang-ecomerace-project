package server

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/uuusman872/learning-go-shop/internal/dto"
	"github.com/uuusman872/learning-go-shop/internal/utils"
)

// createCategory godoc
// @Summary      Create category
// @Description  Admin only
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      dto.CreateCategoryRequest  true  "Category payload"
// @Success      201   {object}  utils.Response
// @Failure      400   {object}  utils.Response
// @Failure      401   {object}  utils.Response
// @Failure      403   {object}  utils.Response
// @Failure      500   {object}  utils.Response
// @Router       /api/v1/categories/ [post]
func (s *Server) createCategory(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	category, err := s.productService.CreateCategory(&req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to create category", err)
		return
	}
	utils.CreateResponse(c, "Category created successfully", category)
}

// getCategory godoc
// @Summary      List categories
// @Tags         Categories
// @Produce      json
// @Success      200  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /api/v1/categories [get]
func (s *Server) getCategory(c *gin.Context) {
	category, err := s.productService.GetCategory()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to get category", err)
		return
	}
	utils.SuccessResponse(c, "Categories retrieved successfully", category)
}

// updateCategory godoc
// @Summary      Update category
// @Description  Admin only
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      int                       true  "Category ID"
// @Param        body  body      dto.UpdateCategoryRequest  true  "Category payload"
// @Success      200   {object}  utils.Response
// @Failure      400   {object}  utils.Response
// @Failure      401   {object}  utils.Response
// @Failure      403   {object}  utils.Response
// @Failure      500   {object}  utils.Response
// @Router       /api/v1/categories/{id} [put]
func (s *Server) updateCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid category ID", err)
		return
	}
	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	category, err := s.productService.UpdateCategory(uint(id), &req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update category", err)
		return
	}
	utils.SuccessResponse(c, "Category updated successfully", category)
}

// deleteCategory godoc
// @Summary      Delete category
// @Description  Admin only
// @Tags         Categories
// @Produce      json
// @Security     BearerAuth
// @Param        id   path  int  true  "Category ID"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      401  {object}  utils.Response
// @Failure      403  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /api/v1/categories/{id} [delete]
func (s *Server) deleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid category ID", err)
		return
	}
	if err := s.productService.DeleteCategory(uint(id)); err != nil {
		utils.InternalServerErrorResponse(c, "Failed to delete category", err)
		return
	}
}

// createProduct godoc
// @Summary      Create product
// @Description  Admin only
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      dto.CreateProductRequest  true  "Product payload"
// @Success      201   {object}  utils.Response
// @Failure      400   {object}  utils.Response
// @Failure      401   {object}  utils.Response
// @Failure      403   {object}  utils.Response
// @Failure      500   {object}  utils.Response
// @Router       /api/v1/products/ [post]
func (s *Server) createProduct(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	product, err := s.productService.CreateProduct(&req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to create product", err)
		return
	}
	utils.CreateResponse(c, "Product created successfully", product)
}

// getProducts godoc
// @Summary      List products
// @Tags         Products
// @Produce      json
// @Param        page   query     int  false  "Page number"  default(1)
// @Param        limit  query     int  false  "Page size"    default(10)
// @Success      200    {object}  utils.PaginationResponse
// @Failure      500    {object}  utils.Response
// @Router       /api/v1/products [get]
func (s *Server) getProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	products, meta, err := s.productService.GetProducts(page, limit)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to fetch products", err)
		return
	}
	utils.PaginationSuccessResponse(c, "Products retrieved successfully", products, *meta)
}

// getProduct godoc
// @Summary      Get product by ID
// @Tags         Products
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      404  {object}  utils.Response
// @Router       /api/v1/products/{id} [get]
func (s *Server) getProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid product ID", err)
		return
	}
	product, err := s.productService.GetProduct(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Product not found")
		return
	}
	utils.SuccessResponse(c, "Product retrieved successfully", product)
}

// updateProduct godoc
// @Summary      Update product
// @Description  Admin only
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      int                      true  "Product ID"
// @Param        body  body      dto.UpdateProductRequest  true  "Product payload"
// @Success      200   {object}  utils.Response
// @Failure      400   {object}  utils.Response
// @Failure      401   {object}  utils.Response
// @Failure      403   {object}  utils.Response
// @Failure      500   {object}  utils.Response
// @Router       /api/v1/products/{id} [put]
func (s *Server) updateProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid product ID", err)
		return
	}
	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	product, err := s.productService.UpdateProduct(uint(id), &req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update product", err)
		return
	}
	utils.SuccessResponse(c, "Product Updated Successfully", product)
}

// deleteProduct godoc
// @Summary      Delete product
// @Description  Admin only
// @Tags         Products
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      401  {object}  utils.Response
// @Failure      403  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /api/v1/products/{id} [delete]
func (s *Server) deleteProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid request ID", err)
		return
	}
	if err := s.productService.DeleteProduct(uint(id)); err != nil {
		utils.InternalServerErrorResponse(c, "Failed to delete product", err)
		return
	}
	utils.SuccessResponse(c, "Product deleted successfully", nil)
}

// uploadProductImage godoc
// @Summary      Upload product image
// @Description  Admin only. Multipart form field name must be `image`.
// @Tags         Products
// @Accept       mpfd
// @Produce      json
// @Security     BearerAuth
// @Param        id     path      int   true  "Product ID"
// @Param        image  formData  file  true  "Image file (.jpg .jpeg .png .gif .webp)"
// @Success      200    {object}  utils.Response
// @Failure      400    {object}  utils.Response
// @Failure      401    {object}  utils.Response
// @Failure      403    {object}  utils.Response
// @Failure      500    {object}  utils.Response
// @Router       /api/v1/products/{id}/images [post]
func (s *Server) uploadProductImage(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalidproduct ID", err)
		return
	}
	file, err := c.FormFile("image")
	if err != nil {
		utils.BadRequestResponse(c, "Invalid image file", err)
		return
	}
	url, err := s.uploadService.UploadProductImage(uint(id), file)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to upload image", err)
		return
	}
	if err := s.productService.AddProductImage(uint(id), url, file.Filename); err != nil {
		utils.InternalServerErrorResponse(c, "Failed to save image record", err)
		return
	}
	utils.SuccessResponse(c, "Image uploaded successfully", map[string]string{"url": url})
}
