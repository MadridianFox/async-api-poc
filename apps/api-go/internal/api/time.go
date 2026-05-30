package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

type ISOTime struct {
	time.Time
}

func (t ISOTime) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte("null"), nil
	}

	return []byte(`"` + t.Time.Format("2006-01-02T15:04:05-07:00") + `"`), nil
}

func (t *ISOTime) UnmarshalJSON(data []byte) error {
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		t.Time = time.Time{}
		return nil
	}

	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02 15:04:05",
	}

	for _, layout := range layouts {
		parsed, err := time.Parse(layout, value)
		if err == nil {
			t.Time = parsed
			return nil
		}
	}

	return fmt.Errorf("invalid time value %q", value)
}

