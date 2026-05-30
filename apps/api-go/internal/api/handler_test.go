package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestCatalogSearchReturnsUpstreamData(t *testing.T) {
	t.Parallel()

	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method %s", r.Method)
		}

		if got := r.URL.Path; got != "/api/products" {
			t.Fatalf("unexpected path %s", got)
		}

		assertQueryValue(t, r.URL.Query(), "price_min", "10")
		assertQueryValue(t, r.URL.Query(), "price_max", "100")
		assertQueryValue(t, r.URL.Query(), "name_like", "apple")
		assertQueryValue(t, r.URL.Query(), "qty_min", "2")
		assertQueryValue(t, r.URL.Query(), "qty_max", "20")
		assertQueryValue(t, r.URL.Query(), "page", "3")

		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{
					"id":         1,
					"name":       "apple",
					"price":      50,
					"qty":        10,
					"created_at": "2026-05-30T10:00:00+00:00",
					"updated_at": "2026-05-30T11:00:00+00:00",
				},
			},
			"meta": map[string]any{
				"total":        1,
				"per_page":     20,
				"current_page": 3,
			},
		})
	}))
	defer upstream.Close()

	app := newTestServer(t, upstream.URL, upstream.URL)

	request := httptest.NewRequest(http.MethodPost, "/api/catalog/search", strings.NewReader(`{"price_min":10,"price_max":100,"name_like":"apple","qty_min":2,"qty_max":20,"page":3}`))
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	app.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("unexpected status %d", response.Code)
	}

	var payload struct {
		Data []Product  `json:"data"`
		Meta searchMeta `json:"meta"`
	}
	decodeBody(t, response.Body.Bytes(), &payload)

	if len(payload.Data) != 1 {
		t.Fatalf("unexpected data length %d", len(payload.Data))
	}

	if payload.Meta.Total != 1 || payload.Meta.PerPage != 20 || payload.Meta.CurrentPage != 3 {
		t.Fatalf("unexpected meta: %+v", payload.Meta)
	}

	if got := payload.Data[0].CreatedAt.Time.Format("2006-01-02T15:04:05-07:00"); got != "2026-05-30T10:00:00+00:00" {
		t.Fatalf("unexpected created_at %s", got)
	}
}

func TestCatalogSearchReturnsEmptyOnBadUpstreamStatus(t *testing.T) {
	t.Parallel()

	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer upstream.Close()

	app := newTestServer(t, upstream.URL, upstream.URL)

	request := httptest.NewRequest(http.MethodPost, "/api/catalog/search", strings.NewReader(`{"page":1}`))
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	app.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("unexpected status %d", response.Code)
	}

	var payload struct {
		Data []Product  `json:"data"`
		Meta searchMeta `json:"meta"`
	}
	decodeBody(t, response.Body.Bytes(), &payload)

	if len(payload.Data) != 0 {
		t.Fatalf("unexpected products length %d", len(payload.Data))
	}

	if payload.Meta.Total != 0 || payload.Meta.PerPage != 0 || payload.Meta.CurrentPage != 0 {
		t.Fatalf("unexpected meta: %+v", payload.Meta)
	}
}

func TestCurrentBasketReturnsData(t *testing.T) {
	t.Parallel()

	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method %s", r.Method)
		}

		if got := r.URL.Path; got != "/api/baskets/current" {
			t.Fatalf("unexpected path %s", got)
		}

		assertQueryValue(t, r.URL.Query(), "user_id", "77")
		assertQueryValue(t, r.URL.Query(), "with_items", "1")

		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"id":      nil,
				"user_id": 77,
				"price":   250,
				"status":  0,
				"items": []map[string]any{
					{
						"id":           9,
						"product_id":   5,
						"price_by_one": 25,
						"qty":          10,
					},
				},
			},
		})
	}))
	defer upstream.Close()

	app := newTestServer(t, upstream.URL, upstream.URL)

	request := httptest.NewRequest(http.MethodGet, "/api/basket/current?user_id=77&with_items=1", nil)
	response := httptest.NewRecorder()
	app.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("unexpected status %d", response.Code)
	}

	var payload basketResponse
	decodeBody(t, response.Body.Bytes(), &payload)

	if payload.Data.ID != nil {
		t.Fatalf("expected nil id, got %v", *payload.Data.ID)
	}

	if payload.Data.UserID != 77 || payload.Data.Price != 250 || payload.Data.Status != 0 {
		t.Fatalf("unexpected basket: %+v", payload.Data)
	}

	if len(payload.Data.Items) != 1 {
		t.Fatalf("unexpected items length %d", len(payload.Data.Items))
	}
}

func TestSetItemValidationFailure(t *testing.T) {
	t.Parallel()

	app := newTestServer(t, "http://127.0.0.1/", "http://127.0.0.1/")

	request := httptest.NewRequest(http.MethodPost, "/api/basket/current/set-item?product_id=5&qty=2", nil)
	response := httptest.NewRecorder()
	app.ServeHTTP(response, request)

	if response.Code != http.StatusUnprocessableEntity {
		t.Fatalf("unexpected status %d", response.Code)
	}
}

func TestBasketUpstreamFailureBecomes500(t *testing.T) {
	t.Parallel()

	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"code":    "BasketItemNotAdded",
			"message": "Basket Item Not Added",
		})
	}))
	defer upstream.Close()

	app := newTestServer(t, upstream.URL, upstream.URL)

	request := httptest.NewRequest(http.MethodPost, "/api/basket/current/set-item?user_id=1&product_id=2&qty=3", nil)
	response := httptest.NewRecorder()
	app.ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("unexpected status %d", response.Code)
	}
}

func newTestServer(t *testing.T, basketHost, productHost string) http.Handler {
	t.Helper()

	return NewServer(Config{
		ListenAddr:    ":0",
		ProductHost:   productHost,
		BasketHost:    basketHost,
		RequestTimeout: 2 * time.Second,
	})
}

func decodeBody[T any](t *testing.T, body []byte, value *T) {
	t.Helper()

	if err := json.Unmarshal(body, value); err != nil {
		t.Fatalf("decode body: %v", err)
	}
}

func assertQueryValue(t *testing.T, query url.Values, key, want string) {
	t.Helper()

	if got := query.Get(key); got != want {
		t.Fatalf("unexpected query %s: got %q want %q", key, got, want)
	}
}

