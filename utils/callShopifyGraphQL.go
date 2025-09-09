package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hiccpet/service/config"
	"io"
	"net/http"
)

type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

type GraphQLResponse struct {
	Data       json.RawMessage          `json:"data"`
	Errors     []map[string]interface{} `json:"errors,omitempty"`
	Extensions map[string]interface{}   `json:"extensions,omitempty"`
}

func (r *GraphQLResponse) String() string {
	pretty, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return fmt.Sprintf("GraphQLResponse(marshal error: %v)", err)
	}
	return string(pretty)
}

// 提供一个通用的方法，把 Data 解析成目标结构体
func (r *GraphQLResponse) UnmarshalData(v interface{}) error {
	return json.Unmarshal(r.Data, v)
}

// CallShopifyGraphQL 通用方法
func CallShopifyGraphQL(query string, variables map[string]interface{}, apiVersion string) (*GraphQLResponse, error) {
	reqBody := GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	bodyBytes, _ := json.Marshal(reqBody)

	if apiVersion == "" {
		apiVersion = "2025-07"

	}

	url := "https://" + config.Cfg.StoreDomain + "/admin/api/" + apiVersion + "/graphql.json"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Shopify-Access-Token", config.Cfg.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	fmt.Println("Raw Response:", string(respBody))

	var graphResp GraphQLResponse
	if err := json.Unmarshal(respBody, &graphResp); err != nil {
		return nil, err
	}

	return &graphResp, nil
}
