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
	currencyCode, exists := c.GetQuery("currencyCode")
	if !exists {
		currencyCode = "SGD"
	}

	fmt.Println("Shopify customer ID:", id)
	query := `#graphql
		query GetCustomer($id: ID!,$currencyCodeQuery: String!){
			customer(id: $id) {
			storeCreditAccounts(first:1, query: $currencyCodeQuery){
				nodes {
						balance {
							amount
							currencyCode
						}
					}
				}
			}
		}`

	resp, err := utils.CallShopifyGraphQL(query, map[string]interface{}{
		"id":                id,
		"currencyCodeQuery": fmt.Sprintf("currency_code:%s", currencyCode),
	}, "")
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, resp.Data)
}
