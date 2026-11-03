package server

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/uuusman872/learning-go-shop/internal/dto"
	"github.com/uuusman872/learning-go-shop/internal/utils"
)

// getCart godoc
// @Summary      Get cart
// @Description  Requires authentication
// @Tags         Cart
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  utils.Response
// @Failure      401  {object}  utils.Response
// @Failure      403  {object}  utils.Response
// @Failure      404  {object}  utils.Response
// @Router       /api/v1/cart/ [get]
func (s *Server) getCart(c *gin.Context) {
	userID := c.GetUint("user_id")
	cart, err := s.cartService.GetCart(userID)
	if err != nil {
		utils.NotFoundResponse(c, "Cart not found")
		return
	}
	utils.SuccessResponse(c, "Cart retrieved successfully", cart)
}

// addToCart godoc
// @Summary      Add item to cart
// @Tags         Cart
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      dto.AddToCartRequest  true  "Cart item payload"
// @Success      200   {object}  utils.Response
// @Failure      400   {object}  utils.Response
// @Failure      401   {object}  utils.Response
// @Failure      403   {object}  utils.Response
// @Router       /api/v1/cart/items [post]
func (s *Server) addToCart(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req dto.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	cart, err := s.cartService.AddToCart(userID, &req)
	if err != nil {
		utils.BadRequestResponse(c, "Failed to add item to cart", err)
		return
	}
	utils.SuccessResponse(c, "Item added to cart successfully", cart)
}

// updateCartItem godoc
// @Summary      Update cart item quantity
// @Tags         Cart
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      int                         true  "Cart item ID"
// @Param        body  body      dto.UpdateCartItemRequest  true  "Quantity payload"
// @Success      200   {object}  utils.Response
// @Failure      400   {object}  utils.Response
// @Failure      401   {object}  utils.Response
// @Failure      403   {object}  utils.Response
// @Router       /api/v1/cart/items/{id} [put]
func (s *Server) updateCartItem(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid id", err)
		return
	}
	var req dto.UpdateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	cart, err := s.cartService.UpdateCartItem(userID, uint(id), &req)
	if err != nil {
		utils.BadRequestResponse(c, "Failed to update item to cart", err)
		return
	}
	utils.SuccessResponse(c, "Item updated successfully", cart)
}

// removeFromCart godoc
// @Summary      Remove item from cart
// @Tags         Cart
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Cart item ID"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      401  {object}  utils.Response
// @Failure      403  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /api/v1/cart/items/{id} [delete]
func (s *Server) removeFromCart(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid id", err)
	}
	if err := s.cartService.RemoveFromCart(userID, uint(id)); err != nil {
		utils.InternalServerErrorResponse(c, "Failed to remove item from cart", err)
		return
	}
	utils.SuccessResponse(c, "Item removed from cart successfully", nil)
}
