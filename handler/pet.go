package handler

import (
	"fmt"
	"net/http"
	"path/filepath"

	"hiccpet/service/middleware"
	"hiccpet/service/model"
	"hiccpet/service/service"

	"hiccpet/service/response"

	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type PetReq struct {
	Phone                 string         `json:"phone"`
	PetAvatarUrl          string         `json:"petAvatarUrl"`
	PetName               string         `json:"petName"`
	PetType               string         `json:"petType"`
	Breed                 string         `json:"breed"`
	Birthday              string         `json:"birthday"`
	Weight                string         `json:"weight"`
	Gender                string         `json:"gender"`
	VaccinationRecords    int            `json:"vaccinationRecords"`
	SterilizationStatus   int            `json:"sterilizationStatus"`
	HasMedicalCondition   bool           `json:"hasMedicalCondition"`
	MedicalConditionMap   datatypes.JSON `json:"medicalConditionMap"`
	MedicalConditionOther string         `json:"medicalConditionOther"`
	CoatType              string         `json:"coatType"`
	GroomingFrequency     string         `json:"groomingFrequency"`
}

func AddPet(c *gin.Context, db *gorm.DB) {
	var req PetReq

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
		Birthday:              req.Birthday,
		Weight:                req.Weight,
		Gender:                req.Gender,
		VaccinationRecords:    &req.VaccinationRecords,
		SterilizationStatus:   &req.SterilizationStatus,
		HasMedicalCondition:   &req.HasMedicalCondition,
		MedicalConditionMap:   req.MedicalConditionMap,
		MedicalConditionOther: req.MedicalConditionOther,
		CoatType:              req.CoatType,
		GroomingFrequency:     req.GroomingFrequency,
	}

	if err := db.Create(&pet).Error; err != nil {
		response.Error(c, http.StatusBadRequest, "Failed to add pet")
		return
	}

	var count int64
	if err := db.Model(&model.Pet{}).
		Where("shopify_customer_id = ?", shopifyCustomerId).
		Count(&count).Error; err != nil {
		fmt.Println("Query error:", err)
	}

	if count == 1 {
		fmt.Println(count, " -----count")

		// 发放权益
		go func(customer model.Customer, pet model.Pet) {
			_ = service.GrantPetBenefit(shopifyCustomerId, db, &customer, &pet)
		}(customer, pet)
	}

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

type PetUpdateReq struct {
	Id uint `json:"id" binding:"required"`
	PetReq
}

func UpdatePetById(c *gin.Context, db *gorm.DB) {
	var req PetUpdateReq
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
		Birthday:              req.Birthday,
		Weight:                req.Weight,
		Gender:                req.Gender,
		VaccinationRecords:    &req.VaccinationRecords,
		SterilizationStatus:   &req.SterilizationStatus,
		HasMedicalCondition:   &req.HasMedicalCondition,
		MedicalConditionMap:   req.MedicalConditionMap,
		MedicalConditionOther: req.MedicalConditionOther,
		CoatType:              req.CoatType,
		GroomingFrequency:     req.GroomingFrequency,
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
