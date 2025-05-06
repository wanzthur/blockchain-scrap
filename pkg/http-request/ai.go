// package ai

package httprequest

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"os"
// )

// type AiDocClient struct {
// 	baseURL    string
// 	httpClient *http.Client
// }

// type APIRequest struct {
// 	Query        string `json:"query"`
// 	StartingDate string `json:"starting_date"`
// 	EndDate      string `json:"end_date"`
// }

// type APIResponse struct {
// 	Success               bool     `json:"success"`
// 	LinkWebsiteCollection []string `json:"link_website_collection"`
// 	ContentSummary        string   `json:"content_summary"`
// }

// func New() *AiDocClient {
// 	return &AiDocClient{
// 		baseURL:    os.Getenv("URL_PROCESS"),
// 		httpClient: &http.Client{},
// 	}
// }

// func (a *AiDocClient) doRequest(endpoint, method string) ([]byte, errs.MessageErr) {
// 	url := a.baseURL + endpoint

// 	var buf bytes.Buffer
// 	req, err := http.NewRequestWithContext(context.Background(), method, url, &buf)
// 	if err != nil {
// 		return nil, errs.NewInternalServerError(fmt.Sprintf("AiDoc - http.NewRequest - err: %v", err))
// 	}

// 	req.Header.Set("Content-Type", "application/json")
// 	resp, err := a.httpClient.Do(req)
// 	if err != nil {
// 		return nil, errs.NewInternalServerError(fmt.Sprintf("AiDoc - httpClient.Do - err: %v", err))
// 	}
// 	defer resp.Body.Close()

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, errs.NewInternalServerError(fmt.Sprintf("AiDoc - io.ReadAll - err: %v", err))
// 	}

// 	switch resp.StatusCode {
// 	case http.StatusOK:
// 		return body, nil
// 	case http.StatusUnprocessableEntity:
// 		return nil, errs.NewUnprocessableEntity("Invalid request body")
// 	case http.StatusRequestTimeout:
// 		return nil, errs.NewRequestTimeout("Request timeout")
// 	default:
// 		return nil, errs.NewInternalServerError("API request failed")
// 	}
// }

// func (a *AiDocClient) SearchContent(request APIRequest) (APIResponse, errs.MessageErr) {
// 	data, err := a.doRequest("/scraping_by_query_search_v3", http.MethodPost, request)
// 	if err != nil {
// 		return APIResponse{}, err
// 	}

// 	var response APIResponse
// 	if err := json.Unmarshal(data, &response); err != nil {
// 		return APIResponse{}, errs.NewInternalServerError(fmt.Sprintf("AiDoc - SearchContent - json.Unmarshal - err: %v", err))
// 	}

// 	return response, nil
// }
