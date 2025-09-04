package utils

import (
	"bytes"
	"encoding/json"
	"hiccpet/service/config"
	"io"
	"net/http"
)

type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

type GraphQLResponse struct {
	Data   json.RawMessage          `json:"data"`
	Errors []map[string]interface{} `json:"errors,omitempty"`
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

	var graphResp GraphQLResponse
	if err := json.Unmarshal(respBody, &graphResp); err != nil {
		return nil, err
	}

	return &graphResp, nil
}
