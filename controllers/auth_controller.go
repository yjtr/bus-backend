package controllers

import (
	"TapTransit-backend/models"
	"TapTransit-backend/utils"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthController struct{}

func NewAuthController() *AuthController {
	return &AuthController{}
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginUser struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Role     string `json:"role"`
}

type loginResponse struct {
	Token string    `json:"token"`
	User  loginUser `json:"user"`
}

// Login 简单登录（开发用，明文密码）
func (a *AuthController) Login(ctx *gin.Context) {
	var req loginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "缺少用户名或密码")
		return
	}

	var user models.User
	if err := utils.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		utils.Unauthorized(ctx, "用户名或密码错误")
		return
	}

	if user.Password != req.Password {
		utils.Unauthorized(ctx, "用户名或密码错误")
		return
	}

	resp := loginResponse{
		Token: fmt.Sprintf("dev_token_%d", time.Now().Unix()),
		User: loginUser{
			ID:       user.ID,
			Username: user.Username,
			Name:     user.RealName,
			Role:     user.Role,
		},
	}
	utils.Success(ctx, resp)
}

// Logout 登出（占位）
func (a *AuthController) Logout(ctx *gin.Context) {
	utils.Success(ctx, gin.H{"ok": true})
}
