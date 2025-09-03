package handler

import (
	"net/http"
	"time"

	"hiccpet/service/model"

	"hiccpet/service/response"

	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AddPet(c *gin.Context, db *gorm.DB) {
	var req struct {
		ShopifyCustomerId     string `json:"shopifyCustomerId" binding:"required"`
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
	if err := db.First(&customer, "shopify_customer_id = ?", req.ShopifyCustomerId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果不存在，则创建新 Customer
			customer = model.Customer{
				ShopifyCustomerId: req.ShopifyCustomerId,
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
		ShopifyCustomerId:     req.ShopifyCustomerId,
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

	shopifyCustomerId := c.Query("shopifyCustomerId")
	if shopifyCustomerId == "" {
		response.Error(c, http.StatusBadRequest, "shopifyCustomerId is required")
		return
	}

	var pets []model.Pet

	if err := db.Where("shopify_customer_id = ?", shopifyCustomerId).Find(&pets).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch pets")
		return
	}

	response.Success(c, pets)
}

func DeletePetById(c *gin.Context, db *gorm.DB) {
	var req struct {
		Id                int    `json:"id" binding:"required"`
		ShopifyCustomerId string `json:"shopifyCustomerId" binding:"required"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查是否存在该客户的宠物
	var pet model.Pet
	if err := db.Where("id = ? AND shopify_customer_id = ?", req.Id, req.ShopifyCustomerId).First(&pet).Error; err != nil {
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

func UpdatePetById(c *gin.Context, db *gorm.DB) {
	type UpdatePetRequest struct {
		ShopifyCustomerId string `json:"shopifyCustomerId" binding:"required"`
		PetId             uint   `json:"petId" binding:"required"`
		Name              string `json:"name"`
		Type              string `json:"type"`
		Birthday          string `json:"birthday"`
	}

	var req UpdatePetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	// 先查找该客户下的 pet
	var pet model.Pet
	if err := db.Where("id = ? AND shopify_customer_id = ?", req.PetId, req.ShopifyCustomerId).
		First(&pet).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, http.StatusNotFound, "pet not found for this customer")
		} else {
			response.Error(c, http.StatusInternalServerError, "failed to query pet")
		}
		return
	}

	// 更新字段
	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Type != "" {
		updates["type"] = req.Type
	}
	if req.Birthday != "" {
		if t, err := time.Parse("2006-01-02", req.Birthday); err == nil {
			updates["birthday"] = t
		} else {
			response.Error(c, http.StatusBadRequest, "invalid birthday format, expected YYYY-MM-DD")
			return
		}
	}

	if len(updates) == 0 {
		response.Error(c, http.StatusBadRequest, "no fields to update")
		return
	}

	if err := db.Model(&pet).Updates(updates).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to update pet")
		return
	}

	response.Success(c, pet)
}
