package handler

import (
	"net/http"

	"hiccpet/service/model"
	"hiccpet/service/utils"

	"hiccpet/service/response"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"hiccpet/service/config"
)

// 注册
func Register(c *gin.Context, db *gorm.DB) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user := model.User{
		Username: req.Username,
		Password: string(hashed),
	}

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "registration successful"})
}

// 登录
func Login(c *gin.Context, db *gorm.DB) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user model.User
	if err := db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		// c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		response.Error(c, http.StatusUnauthorized, "invalid credentials")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		// c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		response.Error(c, http.StatusUnauthorized, "invalid credentials")
		return
	}

	token, _ := utils.GenerateToken(user.Username)
	// c.JSON(http.StatusOK, gin.H{"token": token})

	response.Success(c, gin.H{"token": token, "storeName": config.Cfg.StoreName})
}
