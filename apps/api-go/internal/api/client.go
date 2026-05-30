package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type ProductClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewProductClient(baseURL string, timeout time.Duration) *ProductClient {
	return &ProductClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *ProductClient) SearchProducts(ctx context.Context, filter ProductFilter) ProductSearchResult {
	values := url.Values{}

	addInt := func(key string, value *int) {
		if value != nil {
			values.Set(key, strconvFormatInt(*value))
		}
	}

	addString := func(key string, value *string) {
		if value != nil {
			values.Set(key, *value)
		}
	}

	addInt("price_min", filter.MinPrice)
	addInt("price_max", filter.MaxPrice)
	addString("name_like", filter.NameLike)
	addInt("qty_min", filter.MinQty)
	addInt("qty_max", filter.MaxQty)
	addInt("page", filter.Page)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/api/products?"+values.Encode(), nil)
	if err != nil {
		return emptyProductSearchResult()
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return emptyProductSearchResult()
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return emptyProductSearchResult()
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return emptyProductSearchResult()
	}

	var envelope productSearchEnvelope
	if err := json.Unmarshal(body, &envelope); err != nil {
		return emptyProductSearchResult()
	}

	return envelope.toResult()
}

type productSearchEnvelope struct {
	Data []Product    `json:"data"`
	Meta productMeta  `json:"meta"`
}

type productMeta struct {
	Total       int `json:"total"`
	PerPage     int `json:"per_page"`
	CurrentPage int `json:"current_page"`
}

func (e productSearchEnvelope) toResult() ProductSearchResult {
	products := e.Data
	if products == nil {
		products = []Product{}
	}

	return ProductSearchResult{
		Products:    products,
		Total:       e.Meta.Total,
		PerPage:     e.Meta.PerPage,
		CurrentPage: e.Meta.CurrentPage,
	}
}

func emptyProductSearchResult() ProductSearchResult {
	return ProductSearchResult{
		Products: []Product{},
	}
}

type BasketClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewBasketClient(baseURL string, timeout time.Duration) *BasketClient {
	return &BasketClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *BasketClient) SetItem(ctx context.Context, userID, productID, qty int, withItems bool) (Basket, error) {
	values := url.Values{}
	values.Set("user_id", strconvFormatInt(userID))
	values.Set("product_id", strconvFormatInt(productID))
	values.Set("qty", strconvFormatInt(qty))
	values.Set("with_items", boolQueryValue(withItems))

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/baskets/current/set-item?"+values.Encode(), nil)
	if err != nil {
		return Basket{}, fmt.Errorf("build request: %w", err)
	}

	return c.doBasketRequest(request)
}

func (c *BasketClient) GetCurrentBasket(ctx context.Context, userID int, withItems bool) (Basket, error) {
	values := url.Values{}
	values.Set("user_id", strconvFormatInt(userID))
	values.Set("with_items", boolQueryValue(withItems))

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/api/baskets/current?"+values.Encode(), nil)
	if err != nil {
		return Basket{}, fmt.Errorf("build request: %w", err)
	}

	return c.doBasketRequest(request)
}

func (c *BasketClient) doBasketRequest(request *http.Request) (Basket, error) {
	request.Header.Set("Accept", "application/json")

	response, err := c.httpClient.Do(request)
	if err != nil {
		return Basket{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return Basket{}, fmt.Errorf("unexpected status %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return Basket{}, err
	}

	var envelope basketEnvelope
	if err := json.Unmarshal(body, &envelope); err != nil {
		return Basket{}, err
	}

	if envelope.Data == nil {
		return Basket{}, fmt.Errorf("missing data field")
	}

	if envelope.Data.Items == nil {
		envelope.Data.Items = []BasketItem{}
	}

	return *envelope.Data, nil
}

type basketEnvelope struct {
	Data *Basket `json:"data"`
}

func strconvFormatInt(value int) string {
	return strconv.FormatInt(int64(value), 10)
}

func boolQueryValue(value bool) string {
	if value {
		return "1"
	}

	return "0"
}
