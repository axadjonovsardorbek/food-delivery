package handler

import (
	ap "auth/genproto/auth"
	"context"
	"fmt"
	"net/http"
	"regexp"

	"github.com/golang-jwt/jwt"
	_ "github.com/swaggo/swag"
	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

const emailRegex = `^[a-zA-Z0-9._]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

func isValidEmail(email string) bool {
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}


// Register godoc
// @Summary Register a new user
// @Description Register a new user
// @Tags user
// @Accept json
// @Produce json
// @Param user body ap.UserCreateReq true "User registration request"
// @Success 201 {object} string "User registered"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 500 {object} string "Server error"
// @Router /user/register [post]
func (h *Handler) UserRegister(c *gin.Context) {
	var req ap.UserCreateReq
	
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !isValidEmail(req.Email){
		c.JSON(http.StatusBadRequest, gin.H{"message": "Incorrect email"})
		return
	}

	req.Role = "user"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}
	req.Password = string(hashedPassword)

	_, err = h.User.Register(context.Background(), &req)

	// input, err := json.Marshal(&req)
	// if err != nil {
	// 	c.JSON(500, gin.H{"error": err.Error()})
	// 	return
	// }
	// err = h.Producer.ProduceMessages("user", input)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered"})
}

// Register godoc
// @Summary Register a new courier
// @Description Register a new courier
// @Tags courier
// @Accept json
// @Produce json
// @Param user body ap.UserCreateReq true "User registration request"
// @Success 201 {object} string "Courier registered"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 500 {object} string "Server error"
// @Router /courier/register [post]
func (h *Handler) CourierRegister(c *gin.Context) {
	var req ap.UserCreateReq

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !isValidEmail(req.Email){
		c.JSON(http.StatusBadRequest, gin.H{"message": "Incorrect email"})
		return
	}

	req.Role = "courier"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}
	req.Password = string(hashedPassword)

	_, err = h.User.Register(context.Background(), &req)

	// input, err := json.Marshal(&req)
	// if err != nil {
	// 	c.JSON(500, gin.H{"error": err.Error()})
	// 	return
	// }
	// err = h.Producer.ProduceMessages("user", input)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Courier registered"})
}

// Login godoc
// @Summary Login a user
// @Description Authenticate user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body ap.UserLoginReq true "User login credentials"
// @Success 200 {object} ap.TokenRes "JWT tokens"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 401 {object} string "Invalid email or password"
// @Router /login [post]
func (h *Handler) Login(c *gin.Context) {
	var req ap.UserLoginReq

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.User.Login(context.Background(), &req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get the profile of the authenticated user
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} ap.UserRes
// @Failure 401 {object} string "Unauthorized"
// @Failure 404 {object} string "User not found"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /profile [get]
func (h *Handler) Profile(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := claims.(jwt.MapClaims)["user_id"].(string)

	user, err := h.User.Profile(context.Background(), &ap.ById{Id: id})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update the profil of the authenticated user
// @Tags user
// @Accept json
// @Produce json
// @Param username query string false "Username"
// @Param email query string false "Email"
// @Success 200 {object} string "User profile updated"
// @Failure 401 {object} string "Unauthorized"
// @Failure 404 {object} string "User settings not found"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /profile/update [put]
func (h *Handler) UpdateProfile(c *gin.Context) {
	var req ap.UserUpdateReq

	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := claims.(jwt.MapClaims)["user_id"].(string)
	name := c.Query("username")
	email := c.Query("email")

	if !isValidEmail(email){
		c.JSON(http.StatusBadRequest, gin.H{"message": "Incorrect email"})
		return
	}

	req.Id = id
	req.Username = name
	req.Email = email

	_, err := h.User.UpdateProfile(context.Background(), &req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User profile updated"})
}

// DeleteProfile godoc
// @Summary Delete user profile
// @Description Delete the profil of the authenticated user
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} string "User profile deleted"
// @Failure 401 {object} string "Unauthorized"
// @Failure 404 {object} string "User not found"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /profile/delete [delete]
func (h *Handler) DeleteProfile(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := claims.(jwt.MapClaims)["user_id"].(string)

	_, err := h.User.DeleteProfile(context.Background(), &ap.ById{Id: id})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User profile deleted"})
}

// ChangePassword godoc
// @Summary ChangePassword
// @Description ChangePassword
// @Tags user
// @Accept json
// @Produce json
// @Param current_password query string false "CurrentPassword"
// @Param new_password query string false "NewPassword"
// @Success 200 {object} string "Changed password"
// @Failure 401 {object} string "Unauthorized"
// @Failure 404 {object} string "Password incorrect"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /change-password [put]
func (h *Handler) ChangePassword(c *gin.Context) {
	var req ap.UsersChangePassword

	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := claims.(jwt.MapClaims)["user_id"].(string)
	cur_pass := c.Query("current_password")
	new_pass := c.Query("new_password")

	if cur_pass == "" || cur_pass == "string" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Password incorrect"})
		return
	}
	if new_pass == "" || new_pass == "string" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Password incorrect"})
		return
	}

	req.Id = id
	req.CurrentPassword = cur_pass
	req.NewPasword = new_pass

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPasword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}

	req.NewPasword = string(hashedPassword)

	_, err = h.User.ChangePassword(context.Background(), &req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Changed password"})
}

// ForgotPassword godoc
// @Summary Send a reset password code to the user's email
// @Description Send a reset password code to the user's email
// @Tags user
// @Accept  json
// @Produce  json
// @Param  email  body  ap.UsersForgotPassword  true  "Email data"
// @Success 200 {object} string "Reset password code sent successfully"
// @Failure 400 {object} string "Invalid input"
// @Failure 500 {object} string "Internal server error"
// @Security BearerAuth
// @Router /forgot-password [post]
func (h *Handler) ForgotPassword(c *gin.Context) {
	req := ap.UsersForgotPassword{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	_, err := h.User.ForgotPassword(context.Background(), &req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(req.Email)

	// input,err := json.Marshal(req)
	// if err != nil {
	// 	c.JSON(500, gin.H{"error": err.Error()})
	// }

	// err = h.Producer.ProduceMessages("forgot_password",input)
	// if err != nil {
	// 	c.JSON(500, gin.H{"error": err.Error()})
	// 	return
	// }

	c.JSON(200, gin.H{"message": "Reset password code sent successfully"})
}

// ResetPassword godoc
// @Summary Reset user password
// @Description Reset user password with the provided reset code and new password
// @Tags user
// @Accept  json
// @Produce  json
// @Param reset_token query string false "ResetToken"
// @Param new_password query string false "NewPassword"
// @Success 200 {object} string "Password reset successfully"
// @Failure 400 {object} string "Invalid input"
// @Failure 500 {object} string "Internal server error"
// @Security BearerAuth
// @Router /reset-password [post]
func (h *Handler) ResetPassword(c *gin.Context) {
	var resetCode ap.UsersResetPassword
	reset_token := c.Query("reset_token")
	new_password := c.Query("new_password")


	if new_password == "" || new_password == "string" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Passwrod is empty"})
	}

	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	email := claims.(jwt.MapClaims)["email"].(string)

	resetCode.NewPassword = new_password
	resetCode.ResetToken = reset_token
	resetCode.Email = email

	_, err := h.User.ResetPassword(context.Background(), &resetCode)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Password reset successfully"})
}

// RefreshToken godoc
// @Summary Get token
// @Description Get the token of the authenticated user
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} ap.TokenRes
// @Failure 401 {object} string "Unauthorized"
// @Failure 404 {object} string "User not found"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /refresh-token [get]
func (h *Handler) RefreshToken(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := claims.(jwt.MapClaims)["user_id"].(string)

	token, err := h.User.RefreshToken(context.Background(), &ap.ById{Id: id})

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, token)
}
