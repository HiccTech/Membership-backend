package handler

import (
	"net/http"

	"hiccpet/service/model"

	"hiccpet/service/response"

	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AddPet(c *gin.Context, db *gorm.DB) {
	var req struct {
		ShopifyCustomerID     string `json:"shopifyCustomerID" binding:"required"`
		Phone                 string `json:"phone"`
		PetType               string `json:"petType"`
		Breed                 string `json:"breed"`
		PetIns                string `json:"petIns"`
		Birthday              string `json:"birthday"`
		Gender                string `json:"gender"`
		AdditionalInformation string `json:"additionalInformation"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var customer model.Customer
	if err := db.First(&customer, "shopify_customer_id = ?", req.ShopifyCustomerID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果不存在，则创建新 Customer
			customer = model.Customer{
				ShopifyCustomerID: req.ShopifyCustomerID,
			}
			if err := db.Create(&customer).Error; err != nil {
				response.Error(c, http.StatusInternalServerError, "failed to create customer")
				return
			}
		} else {
			response.Error(c, http.StatusInternalServerError, "database error")
			return
		}
	}

	pet := model.Pet{
		ShopifyCustomerID:     req.ShopifyCustomerID,
		Phone:                 req.Phone,
		PetType:               req.PetType,
		Breed:                 req.Breed,
		PetIns:                req.PetIns,
		Birthday:              req.Birthday,
		Gender:                req.Gender,
		AdditionalInformation: req.AdditionalInformation,
	}

	if err := db.Create(&pet).Error; err != nil {
		response.Error(c, http.StatusBadRequest, "Failed to add pet")
		return
	}

	response.Success(c, pet)
}

func GetPetsByShopifyCustomerID(c *gin.Context, db *gorm.DB) {

	shopifyCustomerID := c.Query("shopifyCustomerID")
	if shopifyCustomerID == "" {
		response.Error(c, http.StatusBadRequest, "shopifyCustomerID is required")
		return
	}

	var pets []model.Pet

	if err := db.Where("shopify_customer_id = ?", shopifyCustomerID).Find(&pets).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch pets")
		return
	}

	response.Success(c, pets)
}

func DeletePetById(c *gin.Context, db *gorm.DB) {
	var req struct {
		Id                int    `json:"id" binding:"required"`
		ShopifyCustomerID string `json:"shopifyCustomerID" binding:"required"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查是否存在该客户的宠物
	var pet model.Pet
	if err := db.Where("id = ? AND shopify_customer_id = ?", req.Id, req.ShopifyCustomerID).First(&pet).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, http.StatusNotFound, "pet not found for this customer")
		} else {
			response.Error(c, http.StatusInternalServerError, "failed to query pet")
		}
		return
	}

	// 删除记录
	if err := db.Delete(&pet).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to delete pet")
		return
	}

	response.Success(c, "pet deleted successfully")
}
