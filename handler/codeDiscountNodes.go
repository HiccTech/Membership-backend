package handler

import (
	"hiccpet/service/response"
	"hiccpet/service/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetCodeDiscountNodes(c *gin.Context, db *gorm.DB) {
	var req struct {
		Query string `json:"query"`
	}

	if err := c.BindJSON(&req); err != nil {
		// c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	query := `#graphql
				query GetCodeDiscountNodes($query:String){
					codeDiscountNodes(first: 250,query:$query,reverse:true) {
						nodes {
							id
							codeDiscount {
								... on DiscountCodeBasic {
									asyncUsageCount
									title
									startsAt
									endsAt
									status
									usageLimit
									codes(first:1){
									nodes{
										code
								}
							}
						}
					}
				}
			}
		}`

	resp, err := utils.CallShopifyGraphQL(query, map[string]interface{}{
		"query": req.Query,
	}, "")
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, resp.Data)
}
