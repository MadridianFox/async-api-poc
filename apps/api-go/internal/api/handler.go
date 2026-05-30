package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type Server struct {
	productClient *ProductClient
	basketClient  *BasketClient
	logger        *log.Logger
}

func NewServer(cfg Config) http.Handler {
	return &Server{
		productClient: NewProductClient(cfg.ProductHost, cfg.RequestTimeout),
		basketClient:  NewBasketClient(cfg.BasketHost, cfg.RequestTimeout),
		logger:        log.New(os.Stdout, "", log.LstdFlags),
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		if r.Method != http.MethodGet {
			s.methodNotAllowed(w, http.MethodGet)
			return
		}
		s.handleRoot(w, r)
	case "/api/catalog/search":
		if r.Method != http.MethodPost {
			s.methodNotAllowed(w, http.MethodPost)
			return
		}
		s.handleCatalogSearch(w, r)
	case "/api/basket/current":
		if r.Method != http.MethodGet {
			s.methodNotAllowed(w, http.MethodGet)
			return
		}
		s.handleCurrentBasket(w, r)
	case "/api/basket/current/set-item":
		if r.Method != http.MethodPost {
			s.methodNotAllowed(w, http.MethodPost)
			return
		}
		s.handleSetItem(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (s *Server) handleRoot(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"service": "api-go",
		"status":  "ok",
	})
}

func (s *Server) handleCatalogSearch(w http.ResponseWriter, r *http.Request) {
	input, err := readRequestData(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	filter, validation := parseProductFilter(input)
	if validation != nil {
		writeValidationError(w, validation)
		return
	}

	result := s.productClient.SearchProducts(r.Context(), filter)
	writeJSON(w, http.StatusOK, searchResponse{
		Data: result.Products,
		Meta: searchMeta{
			Total:       result.Total,
			PerPage:     result.PerPage,
			CurrentPage: result.CurrentPage,
		},
	})
}

func (s *Server) handleCurrentBasket(w http.ResponseWriter, r *http.Request) {
	input, err := readRequestData(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	request, validation := parseBasketCurrentRequest(input)
	if validation != nil {
		writeValidationError(w, validation)
		return
	}

	basket, err := s.basketClient.GetCurrentBasket(r.Context(), request.UserID, request.WithItems)
	if err != nil {
		s.logger.Printf("basket current upstream error: %v", err)
		writeError(w, http.StatusInternalServerError, "Internal Error")
		return
	}

	writeJSON(w, http.StatusOK, basketResponse{Data: basket})
}

func (s *Server) handleSetItem(w http.ResponseWriter, r *http.Request) {
	input, err := readRequestData(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	request, validation := parseSetItemRequest(input)
	if validation != nil {
		writeValidationError(w, validation)
		return
	}

	basket, err := s.basketClient.SetItem(r.Context(), request.UserID, request.ProductID, request.Qty, request.WithItems)
	if err != nil {
		s.logger.Printf("basket set-item upstream error: %v", err)
		writeError(w, http.StatusInternalServerError, "Internal Error")
		return
	}

	writeJSON(w, http.StatusOK, basketResponse{Data: basket})
}

func (s *Server) methodNotAllowed(w http.ResponseWriter, allow string) {
	w.Header().Set("Allow", allow)
	writeError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
}

type basketCurrentRequest struct {
	UserID    int
	WithItems bool
}

type setItemRequest struct {
	UserID    int
	ProductID int
	Qty       int
	WithItems bool
}

func parseProductFilter(input *requestData) (ProductFilter, *validationErrors) {
	filter := ProductFilter{}
	errors := validationErrors{}

	if value, ok := input.lookup("qty_min"); ok {
		parsed, valid, err := toInt(value)
		if err != nil {
			errors.add("qty_min", "The qty_min field must be an integer.")
		} else if parsed < 0 {
			errors.add("qty_min", "The qty_min field must be at least 0.")
		} else if valid {
			filter.MinQty = &parsed
		}
	}

	if value, ok := input.lookup("qty_max"); ok {
		parsed, valid, err := toInt(value)
		if err != nil {
			errors.add("qty_max", "The qty_max field must be an integer.")
		} else if parsed < 0 {
			errors.add("qty_max", "The qty_max field must be at least 0.")
		} else if valid {
			filter.MaxQty = &parsed
		}
	}

	if value, ok := input.lookup("price_min"); ok {
		parsed, valid, err := toInt(value)
		if err != nil {
			errors.add("price_min", "The price_min field must be an integer.")
		} else if parsed < 0 {
			errors.add("price_min", "The price_min field must be at least 0.")
		} else if valid {
			filter.MinPrice = &parsed
		}
	}

	if value, ok := input.lookup("price_max"); ok {
		parsed, valid, err := toInt(value)
		if err != nil {
			errors.add("price_max", "The price_max field must be an integer.")
		} else if parsed < 0 {
			errors.add("price_max", "The price_max field must be at least 0.")
		} else if valid {
			filter.MaxPrice = &parsed
		}
	}

	if value, ok := input.lookup("name_like"); ok {
		parsed, valid := toString(value)
		if !valid {
			errors.add("name_like", "The name_like field must be a string.")
		} else if len(parsed) < 1 {
			errors.add("name_like", "The name_like field must be at least 1 character.")
		} else if len(parsed) > 255 {
			errors.add("name_like", "The name_like field must not exceed 255 characters.")
		} else {
			filter.NameLike = &parsed
		}
	}

	if value, ok := input.lookup("page"); ok {
		parsed, valid, err := toInt(value)
		if err != nil {
			errors.add("page", "The page field must be an integer.")
		} else if parsed < 1 {
			errors.add("page", "The page field must be at least 1.")
		} else if valid {
			filter.Page = &parsed
		}
	}

	if errors.empty() {
		return filter, nil
	}

	return ProductFilter{}, &errors
}

func parseBasketCurrentRequest(input *requestData) (basketCurrentRequest, *validationErrors) {
	request := basketCurrentRequest{}
	errors := validationErrors{}

	userID, err := requiredInt(input, "user_id", "The user_id field is required.", "The user_id field must be an integer.", 1)
	if err != nil {
		errors.add("user_id", err.Error())
	} else {
		request.UserID = *userID
	}

	withItems, err := optionalBool(input, "with_items")
	if err != nil {
		errors.add("with_items", "The with_items field must be a boolean.")
	} else {
		request.WithItems = withItems
	}

	if errors.empty() {
		return request, nil
	}

	return basketCurrentRequest{}, &errors
}

func parseSetItemRequest(input *requestData) (setItemRequest, *validationErrors) {
	request := setItemRequest{}
	errors := validationErrors{}

	userID, err := requiredInt(input, "user_id", "The user_id field is required.", "The user_id field must be an integer.", 1)
	if err != nil {
		errors.add("user_id", err.Error())
	} else {
		request.UserID = *userID
	}

	productID, err := requiredInt(input, "product_id", "The product_id field is required.", "The product_id field must be an integer.", 1)
	if err != nil {
		errors.add("product_id", err.Error())
	} else {
		request.ProductID = *productID
	}

	qty, err := requiredInt(input, "qty", "The qty field is required.", "The qty field must be an integer.", 1)
	if err != nil {
		errors.add("qty", err.Error())
	} else {
		request.Qty = *qty
	}

	withItems, err := optionalBool(input, "with_items")
	if err != nil {
		errors.add("with_items", "The with_items field must be a boolean.")
	} else {
		request.WithItems = withItems
	}

	if errors.empty() {
		return request, nil
	}

	return setItemRequest{}, &errors
}

func requiredInt(input *requestData, key, requiredMessage, integerMessage string, min int) (*int, error) {
	value, ok := input.lookup(key)
	if !ok {
		return nil, errString(requiredMessage)
	}

	parsed, valid, err := toInt(value)
	if err != nil || !valid {
		return nil, errString(integerMessage)
	}

	if parsed < min {
		return nil, errString(integerMessage)
	}

	return &parsed, nil
}

func optionalBool(input *requestData, key string) (bool, error) {
	value, ok := input.lookup(key)
	if !ok {
		return false, nil
	}

	parsed, valid, err := toBool(value)
	if err != nil {
		return false, errString("invalid boolean")
	}
	if !valid {
		return false, nil
	}

	return parsed, nil
}

type validationErrors map[string][]string

func (e validationErrors) add(field, message string) {
	e[field] = append(e[field], message)
}

func (e validationErrors) empty() bool {
	return len(e) == 0
}

func errString(message string) error {
	return &validationMessage{message: message}
}

type validationMessage struct {
	message string
}

func (e *validationMessage) Error() string {
	return e.message
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	_ = encoder.Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]any{
		"message": message,
	})
}

func writeValidationError(w http.ResponseWriter, errors *validationErrors) {
	writeJSON(w, http.StatusUnprocessableEntity, map[string]any{
		"message": "The given data was invalid.",
		"errors":  errors,
	})
}
