package handler

import (
	"fmt"
	"hiccpet/service/middleware"
	"hiccpet/service/response"
	"hiccpet/service/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetStoreCreditBalance(c *gin.Context, db *gorm.DB) {
	id := c.MustGet("shopifyClaims").(*middleware.ShopifyClaims).Sub
	fmt.Println("Shopify customer ID:", id)
	query := `#graphql
		query GetCustomer($id: ID!){
			customer(id: $id) {
			id
			tags
			storeCreditAccounts(first:3){
				nodes {
						balance {
							amount
							currencyCode
						}
						id
					}
				}
			}
		}`

	resp, err := utils.CallShopifyGraphQL(query, map[string]interface{}{
		"id": id,
	}, "")
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, resp.Data)
}
