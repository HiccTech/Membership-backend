package service

import (
	"encoding/json"
	"fmt"
	"hiccpet/service/model"
	"hiccpet/service/response"
	"hiccpet/service/utils"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

type Perk struct {
	StoreCredit  string
	DiscountCode DiscountCode
}

type DiscountCode struct {
	ShopifyDiscountCodeNodeId   string `json:"shopifyDiscountCodeNodeId"`
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

// GetTodayAndNextYear 返回新加坡时区今天零点和一年后的今天23:59:59 (RFC3339格式)
func GetTodayAndNextYear() (string, string) {
	// 加载新加坡时区
	loc, _ := time.LoadLocation("Asia/Singapore")

	now := time.Now().In(loc)

	// 今天零点
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)

	// 一年后的今天 23:59:59
	end := start.AddDate(1, 0, 0).Add(-time.Second)

	return start.Format(time.RFC3339), end.Format(time.RFC3339)
}

// 用独立随机源，避免全局 rand.Seed
var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// 生成唯一 11 位数字
func generate11Digits() string {
	return fmt.Sprintf("%011d", rng.Int63n(1e11))
}

// 生成折扣码
func generateDiscountCode(prefix string) string {
	return prefix + generate11Digits()
}

func GrantPetBenefit(shopifyCustomerId string, db *gorm.DB, customer *model.Customer, pet *model.Pet) error {
	fmt.Println("Shopify customer ID:", shopifyCustomerId)

	// TopupStoreCredit(id, "50.00")

	start, end := GetTodayAndNextYear()
	discountCodes :=
		[]DiscountCode{
			{Title: "Birthday", Code: generateDiscountCode("L"), CustomerGetsValuePercentage: 1, CustomerGetsProductId: "gid://shopify/Product/10227744014517", StartsAt: start, EndsAt: end, UsageLimit: 10},
			{Title: "1V1 Personalized Grooming Class", Code: generateDiscountCode("L"), CustomerGetsValuePercentage: 1, CustomerGetsProductId: "gid://shopify/Product/10227742015669", StartsAt: start, EndsAt: end, UsageLimit: 1},
			// {Title: "Free Massage", Code: generateDiscountCode("L"), CustomerGetsValuePercentage: 1, CustomerGetsProductId: "gid://shopify/Product/10227740934325", StartsAt: start, EndsAt: end, UsageLimit: 3},
			// {Title: "Free Aromatherapyor Grass Mud Spa", Code: generateDiscountCode("L"), CustomerGetsValuePercentage: 1, CustomerGetsProductId: "gid://shopify/Product/10227739754677", StartsAt: start, EndsAt: end, UsageLimit: 3},
			{Title: "Sign-up Gift", Code: generateDiscountCode("L"), CustomerGetsValuePercentage: 1, CustomerGetsProductId: "gid://shopify/Product/10227738411189", StartsAt: start, EndsAt: end, UsageLimit: 3},
			// {Title: "20% off Pet Party Venue Rental", Code: generateDiscountCode("L"), CustomerGetsValuePercentage: 1, CustomerGetsProductId: "gid://shopify/Product/10227725467829", StartsAt: start, EndsAt: end, UsageLimit: 3},
		}

	CreateDiscountCode(shopifyCustomerId, &discountCodes)

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
	for i := range *discountCodes {
		// fmt.Println(i, d)
		d := &(*discountCodes)[i] // 取元素指针，确保修改写回切片

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

		type DiscountResp struct {
			DiscountCodeBasicCreate struct {
				CodeDiscountNode struct {
					ID string `json:"id"`
				} `json:"codeDiscountNode"`
				UserErrors []interface{} `json:"userErrors"`
			} `json:"discountCodeBasicCreate"`
		}

		var data DiscountResp
		if err := json.Unmarshal(resp.Data, &data); err != nil {
			fmt.Println("Error parse:", err)
		}

		fmt.Println("ID:", data.DiscountCodeBasicCreate.CodeDiscountNode.ID)
		d.ShopifyDiscountCodeNodeId = data.DiscountCodeBasicCreate.CodeDiscountNode.ID

	}

	UpdateCustomerMetafield(shopifyCustomerId, discountCodes)

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

func UpdateCustomerMetafield(shopifyCustomerId string, value *[]DiscountCode) {

	sseApp := NewSSEServer()
	sseApp.PushToClient(shopifyCustomerId, `{"code":0,"message":"pending"}`)

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
				JsonValue [][]DiscountCode `json:"jsonValue"` // 直接数组
			} `json:"discountcodejson"`
		} `json:"customer"`
	}

	if err := resp1.UnmarshalData(&discountResp); err != nil {
		panic(err)
	}

	fmt.Println("Discounts:", discountResp.Customer.DiscountCodeJson.JsonValue)
	fmt.Println(resp1, "查询成功")

	// 合并已有折扣和新折扣
	result := append([][]DiscountCode{*value}, discountResp.Customer.DiscountCodeJson.JsonValue...)

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

	createdValueJson, _ := json.Marshal(response.Response{Data: value, Message: "succcess", Code: 0})
	fmt.Println(resp, "更新成功")
	sseApp.PushToClient(shopifyCustomerId, string(createdValueJson))
}
