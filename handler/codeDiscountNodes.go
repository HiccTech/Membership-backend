package handler

import (
	"encoding/json"
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

type CodeDiscountNodeResponse struct {
	CodeDiscountNodeByCode struct {
		CodeDiscount struct {
			Typename   string `json:"__typename"`
			CodesCount struct {
				Count int `json:"count"`
			} `json:"codesCount"`
			ShortSummary string `json:"shortSummary"`
			Title        string `json:"title"`
		} `json:"codeDiscount"`
	} `json:"codeDiscountNodeByCode"`
}

func GetCodeDiscountNodeByCode(code string) (CodeDiscountNodeResponse, error) {
	query := `#graphql
			query codeDiscountNodeByCode($code: String!) {
				codeDiscountNodeByCode(code: $code) {
				codeDiscount {
						... on DiscountCodeBasic {
						title
						}
					}
				}
			}`

	resp, err := utils.CallShopifyGraphQL(query, map[string]interface{}{
		"code": code,
	}, "")
	if err != nil {
		return CodeDiscountNodeResponse{}, err
	}

	result := CodeDiscountNodeResponse{}
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return CodeDiscountNodeResponse{}, err
	}
	return result, nil
}
