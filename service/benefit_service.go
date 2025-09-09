package service

import (
	"encoding/json"
	"fmt"
	"hiccpet/service/model"
	"hiccpet/service/utils"

	"gorm.io/gorm"
)

type Perk struct {
	StoreCredit  string
	DiscountCode DiscountCode
}

type DiscountCode struct {
	Title                       string `json:"title"`
	Code                        string `json:"code"`
	CustomerGetsValuePercentage int    `json:"customerGetsValuePercentage"`
	CustomerGetsProductId       string `json:"customerGetsProductId"`
	StartsAt                    string `json:"startsAt"`
	EndsAt                      string `json:"endsAt"`
	UsageLimit                  int    `json:"usageLimit"`
}

type Discount struct {
	Title                       string `json:"title"`
	Code                        string `json:"code"`
	CustomerGetsValuePercentage int    `json:"customerGetsValuePercentage"`
	CustomerGetsProductId       string `json:"customerGetsProductId"`
	StartsAt                    string `json:"startsAt"`
	EndsAt                      string `json:"endsAt"`
	UsageLimit                  int    `json:"usageLimit"`
}

func GrantPetBenefit(shopifyCustomerId string, db *gorm.DB, customer *model.Customer, pet *model.Pet) error {

	fmt.Println("Shopify customer ID:", shopifyCustomerId)

	// discountCodes := []DiscountCode{
	// 	{Title: "Pet Party Venue Rental1", Code: "DISCOUNT20251", CustomerGetsValuePercentage: 1, CustomerGetsProductId: "gid://shopify/Product/10217653829813", StartsAt: "2025-09-01T00:00:00Z", EndsAt: "2025-12-31T23:59:59Z", UsageLimit: 10},
	// 	{Title: "Pet Party Venue Rental2", Code: "DISCOUNT20252", CustomerGetsValuePercentage: 1, CustomerGetsProductId: "gid://shopify/Product/10217653829813", StartsAt: "2025-09-01T00:00:00Z", EndsAt: "2025-12-31T23:59:59Z", UsageLimit: 1},
	// 	{Title: "Pet Party Venue Rental3", Code: "DISCOUNT20253", CustomerGetsValuePercentage: 1, CustomerGetsProductId: "gid://shopify/Product/10217653829813", StartsAt: "2025-09-01T00:00:00Z", EndsAt: "2025-12-31T23:59:59Z", UsageLimit: 3},
	// }
	// CreateDiscountCode(id, &discountCodes)

	// TopupStoreCredit(id, "50.00")

	// jsonValue := map[string]interface{}{
	// 	"discounts": []map[string]interface{}{
	// 		{"id": 1, "code": "DISCOUNT20251"},
	// 		{"id": 2, "code": "DISCOUNT20252"},
	// 		{"id": 3, "code": "DISCOUNT20253"},
	// 	},
	// }

	discountCodes :=
		[]Discount{
			{Title: "Pet Party Venue Rental1", Code: "DISCOUNT20251", CustomerGetsValuePercentage: 1, CustomerGetsProductId: "gid://shopify/Product/10217653829813", StartsAt: "2025-09-01T00:00:00Z", EndsAt: "2025-12-31T23:59:59Z", UsageLimit: 10},
			{Title: "Pet Party Venue Rental2", Code: "DISCOUNT20252", CustomerGetsValuePercentage: 1, CustomerGetsProductId: "gid://shopify/Product/10217653829813", StartsAt: "2025-09-01T00:00:00Z", EndsAt: "2025-12-31T23:59:59Z", UsageLimit: 1},
			{Title: "Pet Party Venue Rental3", Code: "DISCOUNT20253", CustomerGetsValuePercentage: 1, CustomerGetsProductId: "gid://shopify/Product/10217653829813", StartsAt: "2025-09-01T00:00:00Z", EndsAt: "2025-12-31T23:59:59Z", UsageLimit: 3},
		}

	UpdateCustomerMetafield(shopifyCustomerId, discountCodes)

	return nil
}

func CreateDiscountCode(shopifyCustomerId string, discountCodes *[]DiscountCode) error {
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
	for i, d := range *discountCodes {
		fmt.Println(i, d)

		resp, err := utils.CallShopifyGraphQL(query, map[string]interface{}{
			"basicCodeDiscount": map[string]interface{}{
				"title": d.Title,
				"code":  d.Code,
				"customerSelection": map[string]interface{}{
					"customers": map[string]interface{}{
						"add": []string{shopifyCustomerId},
					},
				},
				"customerGets": map[string]interface{}{
					"value": map[string]interface{}{
						"percentage": d.CustomerGetsValuePercentage,
					},
					"items": map[string]interface{}{
						"all": false,
						"products": map[string]interface{}{
							"productsToAdd": []string{d.CustomerGetsProductId},
						},
					},
				},
				"startsAt":   d.StartsAt,
				"endsAt":     d.EndsAt,
				"usageLimit": d.UsageLimit,
			},
		}, "")
		if err != nil {
			fmt.Println("Error creating discount code:", err)
		}
		print(resp)

	}
	return nil

}

func TopupStoreCredit(shopifyCustomerId string, amount string) {
	query := `#graphql
			mutation storeCreditAccountCredit($id: ID!, $creditInput: StoreCreditAccountCreditInput!) {
				storeCreditAccountCredit(id: $id, creditInput: $creditInput) {
				storeCreditAccountTransaction {
					amount {
						amount
						currencyCode
					}
					account {
						id
						balance {
							amount
							currencyCode
						}
					}
				}
				userErrors {
					message
					field
				}
				}
			}`

	resp, err := utils.CallShopifyGraphQL(query, map[string]interface{}{
		"id": shopifyCustomerId,
		"creditInput": map[string]interface{}{
			"creditAmount": map[string]interface{}{
				"amount":       amount,
				"currencyCode": "SGD",
			},
		},
	}, "")
	if err != nil {
		fmt.Println("Error top up:", err)
	}
	print(resp)
}

func UpdateCustomerMetafield(shopifyCustomerId string, value []Discount) {
	// 查询已有折扣
	queryMetafieldByCustomer := `#graphql
		query ($id:ID!){
			customer(id: $id) {
				id
				email
				discountcodejson:metafield(namespace:"custom",key:"discountcodejson"){
					jsonValue
				}
			}
		}`

	resp1, err := utils.CallShopifyGraphQL(queryMetafieldByCustomer, map[string]interface{}{
		"id": shopifyCustomerId,
	}, "")
	if err != nil {
		fmt.Println("Error query:", err)
		return
	}

	var discountResp struct {
		Customer struct {
			ID               string `json:"id"`
			Email            string `json:"email"`
			DiscountCodeJson struct {
				JsonValue [][]Discount `json:"jsonValue"` // 直接数组
			} `json:"discountcodejson"`
		} `json:"customer"`
	}

	if err := resp1.UnmarshalData(&discountResp); err != nil {
		panic(err)
	}

	fmt.Println("Discounts:", discountResp.Customer.DiscountCodeJson.JsonValue)
	fmt.Println(resp1, "查询成功")

	// 合并已有折扣和新折扣
	result := append([][]Discount{value}, discountResp.Customer.DiscountCodeJson.JsonValue...)

	// Marshal 为 JSON 字符串
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Println("Error marshal JSON:", err)
		return
	}

	// 更新 Shopify metafield
	query := `#graphql
		mutation MetafieldsSet($metafields: [MetafieldsSetInput!]!) {
			metafieldsSet(metafields: $metafields) {
				metafields {
					key
					namespace
					value
					createdAt
					updatedAt
				}
				userErrors {
					field
					message
					code
				}
			}
		}`

	resp, err := utils.CallShopifyGraphQL(query, map[string]interface{}{
		"metafields": []map[string]interface{}{
			{
				"key":       "discountcodejson",
				"namespace": "custom",
				"ownerId":   shopifyCustomerId,
				"type":      "json",
				"value":     string(data),
			},
		},
	}, "")
	if err != nil {
		fmt.Println("Error update metafield:", err)
		return
	}

	fmt.Println(resp, "更新成功")
}
