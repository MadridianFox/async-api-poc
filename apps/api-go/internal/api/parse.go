package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type requestData struct {
	form url.Values
	json map[string]any
}

func readRequestData(r *http.Request) (*requestData, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	data := &requestData{
		form: r.Form,
	}

	contentType := strings.ToLower(r.Header.Get("Content-Type"))
	if !strings.Contains(contentType, "json") {
		return data, nil
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	if len(bytes.TrimSpace(body)) == 0 {
		return data, nil
	}

	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.UseNumber()

	var raw map[string]any
	if err := decoder.Decode(&raw); err != nil {
		return nil, err
	}

	data.json = raw

	return data, nil
}

func (d *requestData) lookup(key string) (any, bool) {
	if d.json != nil {
		if value, ok := d.json[key]; ok {
			if value == nil {
				return nil, false
			}
			return value, true
		}
	}

	if values, ok := d.form[key]; ok && len(values) > 0 {
		return values[0], true
	}

	return nil, false
}

func toString(value any) (string, bool) {
	switch typed := value.(type) {
	case string:
		return typed, true
	case fmt.Stringer:
		return typed.String(), true
	case json.Number:
		return typed.String(), true
	default:
		return "", false
	}
}

func toInt(value any) (int, bool, error) {
	switch typed := value.(type) {
	case int:
		return typed, true, nil
	case int8:
		return int(typed), true, nil
	case int16:
		return int(typed), true, nil
	case int32:
		return int(typed), true, nil
	case int64:
		return int(typed), true, nil
	case uint:
		return int(typed), true, nil
	case uint8:
		return int(typed), true, nil
	case uint16:
		return int(typed), true, nil
	case uint32:
		return int(typed), true, nil
	case uint64:
		return int(typed), true, nil
	case float32:
		if math.Trunc(float64(typed)) != float64(typed) {
			return 0, false, fmt.Errorf("not an integer")
		}
		return int(typed), true, nil
	case float64:
		if math.Trunc(typed) != typed {
			return 0, false, fmt.Errorf("not an integer")
		}
		return int(typed), true, nil
	case json.Number:
		parsed, err := typed.Int64()
		if err != nil {
			return 0, false, err
		}
		return int(parsed), true, nil
	case string:
		if strings.TrimSpace(typed) == "" {
			return 0, false, nil
		}
		parsed, err := strconv.Atoi(typed)
		if err != nil {
			return 0, false, err
		}
		return parsed, true, nil
	default:
		return 0, false, fmt.Errorf("unsupported integer type %T", value)
	}
}

func toBool(value any) (bool, bool, error) {
	switch typed := value.(type) {
	case bool:
		return typed, true, nil
	case int:
		return typed != 0, true, nil
	case int8:
		return typed != 0, true, nil
	case int16:
		return typed != 0, true, nil
	case int32:
		return typed != 0, true, nil
	case int64:
		return typed != 0, true, nil
	case float32:
		return typed != 0, true, nil
	case float64:
		return typed != 0, true, nil
	case json.Number:
		parsed, err := typed.Int64()
		if err != nil {
			return false, false, err
		}
		return parsed != 0, true, nil
	case string:
		if strings.TrimSpace(typed) == "" {
			return false, false, nil
		}

		switch typed {
		case "1":
			return true, true, nil
		case "0":
			return false, true, nil
		default:
			parsed, err := strconv.ParseBool(typed)
			if err != nil {
				return false, false, err
			}
			return parsed, true, nil
		}
	default:
		return false, false, fmt.Errorf("unsupported boolean type %T", value)
	}
}
