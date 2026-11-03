package server

import (
	"github.com/gin-gonic/gin"
	"github.com/uuusman872/learning-go-shop/internal/dto"
	"github.com/uuusman872/learning-go-shop/internal/utils"
)

// register godoc
// @Summary      Register a new user
// @Description  Create a customer account
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.RegisterRequest  true  "Register payload"
// @Success      201   {object}  utils.Response
// @Failure      400   {object}  utils.Response
// @Router       /api/v1/auth/register [post]
func (s *Server) register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "invalid request data", err)
		return
	}

	response, err := s.authServices.Register(&req)
	if err != nil {
		utils.BadRequestResponse(c, "Registration failed", err)
		return
	}
	utils.CreateResponse(c, "User registered successfully", response)

}

// login godoc
// @Summary      Login
// @Description  Authenticate with email and password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.LoginRequest  true  "Login payload"
// @Success      200   {object}  utils.Response
// @Failure      400   {object}  utils.Response
// @Failure      401   {object}  utils.Response
// @Router       /api/v1/auth/login [post]
func (s *Server) login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	response, err := s.authServices.Login(&req)
	if err != nil {
		utils.UnauthorizedResponse(c, "Login Failed")
		return
	}

	utils.SuccessResponse(c, "Login successful", response)
}

// refreshToken godoc
// @Summary      Refresh access token
// @Description  Exchange a refresh token for a new token pair
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.RefreshTokenRequest  true  "Refresh token payload"
// @Success      200   {object}  utils.Response
// @Failure      400   {object}  utils.Response
// @Failure      401   {object}  utils.Response
// @Router       /api/v1/auth/refresh [post]
func (s *Server) refreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	response, err := s.authServices.RefreshToken(&req)
	if err != nil {
		utils.UnauthorizedResponse(c, "Token refresh failed")
		return
	}
	utils.SuccessResponse(c, "Token refreshed successfully", response)
}

// logout godoc
// @Summary      Logout
// @Description  Invalidate the refresh token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.RefreshTokenRequest  true  "Refresh token payload"
// @Success      200   {object}  utils.Response
// @Failure      400   {object}  utils.Response
// @Failure      500   {object}  utils.Response
// @Router       /api/v1/auth/logout [post]
func (s *Server) logout(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	if err := s.authServices.Logout(req.RefreshToken); err != nil {
		utils.InternalServerErrorResponse(c, "Logout failed", err)
		return
	}
	utils.SuccessResponse(c, "Logout successful", nil)
}

// getProfile godoc
// @Summary      Get current user profile
// @Tags         Users
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  utils.Response
// @Failure      401  {object}  utils.Response
// @Failure      404  {object}  utils.Response
// @Router       /api/v1/users/profile [get]
func (s *Server) getProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	profile, err := s.userService.GetProfile(userID)
	if err != nil {
		utils.NotFoundResponse(c, "User not found")
		return
	}
	utils.SuccessResponse(c, "Profile retrieved successfully", profile)
}

// updateProfile godoc
// @Summary      Update current user profile
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      dto.UpdateProfileRequest  true  "Profile payload"
// @Success      200   {object}  utils.Response
// @Failure      400   {object}  utils.Response
// @Failure      401   {object}  utils.Response
// @Failure      500   {object}  utils.Response
// @Router       /api/v1/users/profile [put]
func (s *Server) updateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	profile, err := s.userService.UpdateProfile(userID, &req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to Update profile", err)
		return
	}
	utils.SuccessResponse(c, "Profile update successfully", profile)
}
