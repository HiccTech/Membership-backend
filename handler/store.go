package handler

import (
	"net/http"

	"hiccpet/service/model"

	"hiccpet/service/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AddStore(c *gin.Context, db *gorm.DB) {
	var req struct {
		StoreName   string `json:"storeName"`
		CountryCode string `json:"countryCode"`
		StoreDomain string `json:"storeDomain"`
		AccessToken string `json:"accessToken"`
		Admin       string `json:"admin"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	store := model.Store{
		StoreName:   req.StoreName,
		CountryCode: req.CountryCode,
		StoreDomain: req.StoreDomain,
		AccessToken: req.AccessToken,
		Admin:       req.Admin,
	}

	if err := db.Create(&store).Error; err != nil {
		// c.JSON(http.StatusBadRequest, gin.H{"error": "CountryCode already exists"})
		response.Error(c, http.StatusBadRequest, "CountryCode already exists")
		return
	}

	response.Success(c, store)
	// c.JSON(http.StatusOK, gin.H{"message": "registration successful"})
}

func GetStores(c *gin.Context, db *gorm.DB) {

	type StoreResponse struct {
		ID          uint   `json:"id"`
		StoreName   string `json:"storeName"`
		CountryCode string `json:"countryCode"`
		StoreDomain string `json:"storeDomain"`
		Admin       string `json:"admin"`
	}

	var stores []model.Store

	if err := db.Find(&stores).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch stores")
		return
	}

	var storeResponses []StoreResponse
	for _, s := range stores {
		storeResponses = append(storeResponses, StoreResponse{
			ID:          s.ID,
			StoreName:   s.StoreName,
			CountryCode: s.CountryCode,
			StoreDomain: s.StoreDomain,
			Admin:       s.Admin,
		})
	}

	if storeResponses == nil {
		storeResponses = []StoreResponse{}
	}

	response.Success(c, storeResponses)
}
