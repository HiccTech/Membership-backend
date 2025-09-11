package handler

import (
	"encoding/json"
	"fmt"
	"hiccpet/service/response"
	"hiccpet/service/service"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Order struct {
	ID        int64      `json:"id"`
	LineItems []LineItem `json:"line_items"`
	Customer  Customer   `json:"customer"`
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
			service.TopupStoreCredit(customerId, "1000")
			service.CreateDiscountCode(customerId, &[]service.DiscountCode{
				{Title: "Free Massage 10 sessions", Code: service.GenerateDiscountCode("C"), CustomerGetsValuePercentage: 1, CustomerGetsProductId: "gid://shopify/Product/10227740934325", StartsAt: start, EndsAt: end, UsageLimit: 10},
				{Title: "Free Aromatherapyor Grass Mud Spa", Code: service.GenerateDiscountCode("C"), CustomerGetsValuePercentage: 1, CustomerGetsProductId: "gid://shopify/Product/10227739754677", StartsAt: start, EndsAt: end, UsageLimit: 1},
				{Title: "Pet Party Venue Rental Free 1h", Code: service.GenerateDiscountCode("C"), CustomerGetsValuePercentage: 1, CustomerGetsProductId: "gid://shopify/Product/10227792937141", StartsAt: start, EndsAt: end, UsageLimit: 1},
			})
			break forLoop
		case 10228688453813:
			// 充值2000
			println("充值2000")
			service.TopupStoreCredit(customerId, "2000")
			service.CreateDiscountCode(customerId, &[]service.DiscountCode{
				{Title: "Free Massage 20 sessions", Code: service.GenerateDiscountCode("P"), CustomerGetsValuePercentage: 1, CustomerGetsProductId: "gid://shopify/Product/10227793035445", StartsAt: start, EndsAt: end, UsageLimit: 20},
				{Title: "Free Aromatherapyor Grass Mud Spa", Code: service.GenerateDiscountCode("P"), CustomerGetsValuePercentage: 1, CustomerGetsProductId: "gid://shopify/Product/10227739754677", StartsAt: start, EndsAt: end, UsageLimit: 1},
				{Title: "Pet Party Venue Rental Free 3h", Code: service.GenerateDiscountCode("P"), CustomerGetsValuePercentage: 1, CustomerGetsProductId: "gid://shopify/Product/10227792969909", StartsAt: start, EndsAt: end, UsageLimit: 1},
			})
			break forLoop
		default:

		}
	}

	response.Success(c, "充值成功")
}
