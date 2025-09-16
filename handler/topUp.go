package handler

import (
	"encoding/json"
	"fmt"
	"hiccpet/service/email"
	"hiccpet/service/middleware"
	"hiccpet/service/model"
	"hiccpet/service/response"
	"hiccpet/service/service"
	"hiccpet/service/utils"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Order struct {
	ID                   int64                 `json:"id"`
	LineItems            []LineItem            `json:"line_items"`
	Customer             Customer              `json:"customer"`
	DiscountApplications []DiscountApplication `json:"discount_applications"`
}

type LineItem struct {
	ID                int64    `json:"id"`
	AdminGraphqlApiID string   `json:"admin_graphql_api_id"`
	CurrentQuantity   int      `json:"current_quantity"`
	Name              string   `json:"name"`
	Price             string   `json:"price"`
	PriceSet          PriceSet `json:"price_set"`
	ProductID         int64    `json:"product_id"`
	Quantity          int      `json:"quantity"`
	SKU               string   `json:"sku"`
	Title             string   `json:"title"`
	VariantID         int64    `json:"variant_id"`
}

type Customer struct {
	AdminGraphqlApiId string `json:"admin_graphql_api_id"`
	Email             string `json:"email"`
}

type PriceSet struct {
	ShopMoney        Money `json:"shop_money"`
	PresentmentMoney Money `json:"presentment_money"`
}

type Money struct {
	Amount       string `json:"amount"`
	CurrencyCode string `json:"currency_code"`
}

type DiscountApplication struct {
	Type string `json:"type"`
	Code string `json:"code"`
}

func HandleTopUp(c *gin.Context, db *gorm.DB) {
	b, _ := io.ReadAll(c.Request.Body)
	fmt.Println(string(b), " -----------------------")

	var order Order
	if err := json.Unmarshal(b, &order); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Printf("订单ID: %d\n", order.ID)

	customerId := order.Customer.AdminGraphqlApiId
	start, end := service.GetTodayAndNextYear()

forLoop:
	for _, item := range order.LineItems {
		fmt.Printf("商品: %s, 价格: %s %s , 产品id：%d\n",
			item.Title,
			item.PriceSet.ShopMoney.Amount,
			item.PriceSet.ShopMoney.CurrencyCode,
			item.ProductID,
		)

		switch item.ProductID {
		case 10228688158901:
			// 充值1000
			println("充值1000")
			storeTopup(c, db, 1, customerId, &order)
			service.TopupStoreCredit(customerId, "1000", end)

			discountCodes := []service.DiscountCode{
				{Title: "Free Massage 10 sessions", Code: service.GenerateDiscountCode("C"), CustomerGetsValuePercentage: 1, CustomerGetsProductId: "gid://shopify/Product/10227740934325", StartsAt: start, EndsAt: end, UsageLimit: 10},
				{Title: "Free Aromatherapyor Grass Mud Spa", Code: service.GenerateDiscountCode("C"), CustomerGetsValuePercentage: 1, CustomerGetsProductId: "gid://shopify/Product/10227739754677", StartsAt: start, EndsAt: end, UsageLimit: 1},
				{Title: "Pet Party Venue Rental Free 1h", Code: service.GenerateDiscountCode("C"), CustomerGetsValuePercentage: 1, CustomerGetsProductId: "gid://shopify/Product/10227792937141", StartsAt: start, EndsAt: end, UsageLimit: 1},
			}
			service.CreateDiscountCode(customerId, &discountCodes)
			service.AddTagsToCustomer(customerId, "Club 1000")
			expiredAt, err := utils.FormatDate(end)
			if err != nil {
				fmt.Println("Error formatting date:", err)
			}
			service.SendEmail(
				service.SendEmailData{ShopifyCustomerId: customerId, CustomerEmail: order.Customer.Email, Template: "email/clubEmailWithTopup.tmpl", Subject: "Thank you for your top-up of $1000", DiscountCodes: &discountCodes,
					StoreCredit: &email.StoreCredit{Amount: 1000, Currency: "$", ExpiredAt: expiredAt}},
			)
			break forLoop
		case 10228688453813:
			// 充值2000
			println("充值2000")
			storeTopup(c, db, 2, customerId, &order)
			service.TopupStoreCredit(customerId, "2000", end)
			discountCodes := []service.DiscountCode{
				{Title: "Free Massage 20 sessions", Code: service.GenerateDiscountCode("P"), CustomerGetsValuePercentage: 1, CustomerGetsProductId: "gid://shopify/Product/10227793035445", StartsAt: start, EndsAt: end, UsageLimit: 20},
				{Title: "Free Aromatherapyor Grass Mud Spa", Code: service.GenerateDiscountCode("P"), CustomerGetsValuePercentage: 1, CustomerGetsProductId: "gid://shopify/Product/10227739754677", StartsAt: start, EndsAt: end, UsageLimit: 1},
				{Title: "Pet Party Venue Rental Free 3h", Code: service.GenerateDiscountCode("P"), CustomerGetsValuePercentage: 1, CustomerGetsProductId: "gid://shopify/Product/10227792969909", StartsAt: start, EndsAt: end, UsageLimit: 1},
			}
			service.CreateDiscountCode(customerId, &discountCodes)
			service.AddTagsToCustomer(customerId, "Club 2000")

			expiredAt, err := utils.FormatDate(end)
			if err != nil {
				fmt.Println("Error formatting date:", err)
			}
			service.SendEmail(
				service.SendEmailData{ShopifyCustomerId: customerId, CustomerEmail: order.Customer.Email, Template: "email/clubEmailWithTopup.tmpl", Subject: "Thank you for your top-up of $2000", DiscountCodes: &discountCodes,
					StoreCredit: &email.StoreCredit{Amount: 2000, Currency: "$", ExpiredAt: expiredAt}},
			)
			break forLoop
		case 10227739754677, 10227740934325, 10227792937141:
			if len(order.DiscountApplications) != 0 {
				println("权益消费 ", order.DiscountApplications[0].Code)
				code := order.DiscountApplications[0].Code
				if code != "" {
					if res, err := GetCodeDiscountNodeByCode(code); err == nil {
						fmt.Println("使用权益: ", res.CodeDiscountNodeByCode.CodeDiscount.Title)
					} else {
						fmt.Println("查询权益失败: ", err)
					}
				}
			}
		default:

		}
	}

	response.Success(c, "充值成功")
}

func storeTopup(c *gin.Context, db *gorm.DB, topupType int, shopifyCustomerId string, order *Order) {
	topup := model.Topup{
		OrderId:           order.ID,
		Type:              topupType,
		ShopifyCustomerId: shopifyCustomerId,
		Email:             order.Customer.Email,
	}

	if err := db.Create(&topup).Error; err != nil {
		response.Error(c, http.StatusBadRequest, "Failed to add topup")
		return
	}
}

func TopupCount(c *gin.Context, db *gorm.DB) {

	shopifyCustomerId := c.MustGet("shopifyClaims").(*middleware.ShopifyClaims).Sub

	if shopifyCustomerId == "" {
		response.Error(c, http.StatusBadRequest, "shopifyCustomerId is required")
		return
	}

	var count1 int64
	if err := db.Model(&model.Topup{}).
		Where("shopify_customer_id = ?", shopifyCustomerId).Where("type = ?", 1).
		Count(&count1).Error; err != nil {
		fmt.Println("Query error:", err)
	}

	var count2 int64
	if err := db.Model(&model.Topup{}).
		Where("shopify_customer_id = ?", shopifyCustomerId).Where("type = ?", 2).
		Count(&count2).Error; err != nil {
		fmt.Println("Query error:", err)
	}

	response.Success(c, gin.H{"count1": count1, "count2": count2})
}
