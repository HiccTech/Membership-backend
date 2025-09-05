package service

import (
	"hiccpet/service/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GrantPetBenefit(c *gin.Context, db *gorm.DB, customer *model.Customer, pet *model.Pet) error {

	return nil
}
