package server

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/uuusman872/learning-go-shop/internal/utils"
)

// createOrder godoc
// @Summary      Create order from cart
// @Description  Creates an order from the authenticated user's cart items and clears the cart
// @Tags         Orders
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      401  {object}  utils.Response
// @Failure      403  {object}  utils.Response
// @Router       /api/v1/orders/ [post]
func (s *Server) createOrder(c *gin.Context) {
	userID := c.GetUint("user_id")
	order, err := s.orderService.CreateOrder(userID)
	if err != nil {
		utils.BadRequestResponse(c, "Failed to create order", err)
		return
	}
	utils.SuccessResponse(c, "Order created successfully", order)
}

// getOrders godoc
// @Summary      List orders
// @Description  Returns a paginated list of the authenticated user's orders
// @Tags         Orders
// @Produce      json
// @Security     BearerAuth
// @Param        page   query     int  false  "Page number"  default(1)
// @Param        limit  query     int  false  "Page size"    default(10)
// @Success      200    {object}  utils.PaginationResponse
// @Failure      401    {object}  utils.Response
// @Failure      403    {object}  utils.Response
// @Failure      500    {object}  utils.Response
// @Router       /api/v1/orders/ [get]
func (s *Server) getOrders(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	orders, meta, err := s.orderService.GetOrders(userID, page, limit)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to get orders", err)
		return
	}
	utils.PaginationSuccessResponse(c, "Orders", orders, *meta)
}

// getOrder godoc
// @Summary      Get order by ID
// @Description  Returns a single order belonging to the authenticated user
// @Tags         Orders
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Order ID"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      401  {object}  utils.Response
// @Failure      403  {object}  utils.Response
// @Failure      404  {object}  utils.Response
// @Router       /api/v1/orders/{id} [get]
func (s *Server) getOrder(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid order id", err)
		return
	}
	order, err := s.orderService.GetOrder(userID, uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "Order not found")
		return
	}
	utils.SuccessResponse(c, "Order retrieved successfully", order)
}
