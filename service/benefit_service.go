package service

import (
	"fmt"
	"hiccpet/service/middleware"
	"hiccpet/service/model"
	"hiccpet/service/response"
	"hiccpet/service/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GrantPetBenefit(c *gin.Context, db *gorm.DB, customer *model.Customer, pet *model.Pet) error {

	id := c.MustGet("shopifyClaims").(*middleware.ShopifyClaims).Sub
	fmt.Println("Shopify customer ID:", id)
	query := `#graphql
			mutation CreateDiscountCode($basicCodeDiscount: DiscountCodeBasicInput!) {
				discountCodeBasicCreate(basicCodeDiscount: $basicCodeDiscount) {
				codeDiscountNode {
					id
				}
					userErrors {
						field
						message
					}
				}
			}`

	resp, err := utils.CallShopifyGraphQL(query, map[string]interface{}{
		"basicCodeDiscount": map[string]interface{}{
			"title": "10%/ off selected items",
			"code":  "DISCOUNT2024",
			"customerSelection": map[string]interface{}{
				"customers": map[string]interface{}{
					"add": []string{id},
				},
			},
			"customerGets": map[string]interface{}{
				"value": map[string]interface{}{
					"percentage": 1,
				},
				"items": map[string]interface{}{
					"all": true,
				},
			},
			"startsAt":   "2025-09-01T00:00:00Z",
			"endsAt":     "2025-12-31T23:59:59Z",
			"usageLimit": 10,
		},
	}, "")
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
	}
	print(resp)

	return nil
}
