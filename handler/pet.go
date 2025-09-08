package handler

import (
	"net/http"
	"path/filepath"

	"hiccpet/service/middleware"
	"hiccpet/service/model"
	"hiccpet/service/service"

	"hiccpet/service/response"

	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func AddPet(c *gin.Context, db *gorm.DB) {
	var req struct {
		Phone                 string `json:"phone"`
		PetAvatarUrl          string `json:"petAvatarUrl"`
		PetName               string `json:"petName"`
		PetType               string `json:"petType"`
		Breed                 string `json:"breed"`
		PetIns                string `json:"petIns"`
		Birthday              string `json:"birthday"`
		Weight                string `json:"weight"`
		Gender                string `json:"gender"`
		AdditionalInformation string `json:"additionalInformation"`
	}

	shopifyCustomerId := c.MustGet("shopifyClaims").(*middleware.ShopifyClaims).Sub

	if err := c.BindJSON(&req); err != nil {
		// c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var customer model.Customer
	if err := db.First(&customer, "shopify_customer_id = ?", shopifyCustomerId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果不存在，则创建新 Customer
			customer = model.Customer{
				ShopifyCustomerId: shopifyCustomerId,
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
		ShopifyCustomerId:     shopifyCustomerId,
		Phone:                 req.Phone,
		PetAvatarUrl:          req.PetAvatarUrl,
		PetName:               req.PetName,
		PetType:               req.PetType,
		Breed:                 req.Breed,
		PetIns:                req.PetIns,
		Birthday:              req.Birthday,
		Weight:                req.Weight,
		Gender:                req.Gender,
		AdditionalInformation: req.AdditionalInformation,
	}

	if err := db.Create(&pet).Error; err != nil {
		response.Error(c, http.StatusBadRequest, "Failed to add pet")
		return
	}

	// 发放权益
	go func(customer model.Customer, pet model.Pet) {
		_ = service.GrantPetBenefit(c, db, &customer, &pet)
	}(customer, pet)

	response.Success(c, pet)
}

func GetPetsByShopifyCustomerID(c *gin.Context, db *gorm.DB) {

	shopifyCustomerId := c.MustGet("shopifyClaims").(*middleware.ShopifyClaims).Sub

	if shopifyCustomerId == "" {
		// response.Error(c, http.StatusBadRequest, "shopifyCustomerId is required")
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
		Id int `json:"id" binding:"required"`
	}
	if err := c.BindJSON(&req); err != nil {
		// c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	shopifyCustomerId := c.MustGet("shopifyClaims").(*middleware.ShopifyClaims).Sub

	// 检查是否存在该客户的宠物
	var pet model.Pet
	if err := db.Where("id = ? AND shopify_customer_id = ?", req.Id, shopifyCustomerId).First(&pet).Error; err != nil {
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
		Id                    uint   `json:"id" binding:"required"`
		Phone                 string `json:"phone"`
		PetAvatarUrl          string `json:"petAvatarUrl"`
		PetName               string `json:"petName"`
		PetType               string `json:"petType"`
		Breed                 string `json:"breed"`
		PetIns                string `json:"petIns"`
		Birthday              string `json:"birthday"`
		Weight                string `json:"weight"`
		Gender                string `json:"gender"`
		AdditionalInformation string `json:"additionalInformation"`
	}

	var req UpdatePetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	shopifyCustomerId := c.MustGet("shopifyClaims").(*middleware.ShopifyClaims).Sub

	// 先查找该客户下的 pet
	var pet model.Pet
	if err := db.Where("id = ? AND shopify_customer_id = ?", req.Id, shopifyCustomerId).
		First(&pet).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, http.StatusNotFound, "pet not found for this customer")
		} else {
			response.Error(c, http.StatusInternalServerError, "failed to query pet")
		}
		return
	}

	// 更新字段
	updates := model.Pet{
		Phone:                 req.Phone,
		PetAvatarUrl:          req.PetAvatarUrl,
		PetName:               req.PetName,
		PetType:               req.PetType,
		Breed:                 req.Breed,
		PetIns:                req.PetIns,
		Birthday:              req.Birthday,
		Weight:                req.Weight,
		Gender:                req.Gender,
		AdditionalInformation: req.AdditionalInformation,
	}

	if err := db.Model(&pet).Updates(updates).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to update pet")
		return
	}

	// 返回最新数据
	if err := db.First(&pet, pet.ID).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to fetch updated pet")
		return
	}

	response.Success(c, pet)
}

func UploadPetAvatar(c *gin.Context, db *gorm.DB) {
	file, err := c.FormFile("image")
	if err != nil {
		// c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	filename := uuid.New().String() + filepath.Ext(file.Filename)
	savePath := "uploadPetImgs/" + filename

	// 保存文件
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{"url": "/" + filename})
}
