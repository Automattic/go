package httputils

import (
	"net/url"
	"testing"
)

func TestGetFormValueByIndex(t *testing.T) {
	var value string

	form := url.Values{
		"key": {"first", "second"},
	}
	key := "key"
	defaultValue := "default"

	// empty, index 0
	value = GetFormValueByIndex(url.Values{}, key, 0, defaultValue)
	if value != defaultValue {
		t.Errorf("empty form, index 0; got %s; expected %d", value, defaultValue)
	}

	value = GetFormValueByIndex(form, key, -1, defaultValue)
	if value != defaultValue {
		t.Errorf("invalid, index < 0; got %s; expected %d", value, defaultValue)
	}

	value = GetFormValueByIndex(form, key, 2, defaultValue)
	if value != defaultValue {
		t.Errorf("index == len; got %s; expected %d", value, defaultValue)
	}

	value = GetFormValueByIndex(form, key, 3, defaultValue)
	if value != defaultValue {
		t.Errorf("index > len; got %s; expected %d", value, defaultValue)
	}

	value = GetFormValueByIndex(form, "invalid-key", 0, defaultValue)
	if value != defaultValue {
		t.Errorf("invalid key; got %s; expected %d", value, defaultValue)
	}

	value = GetFormValueByIndex(form, key, 0, defaultValue)
	if value != form[key][0] {
		t.Errorf("valid; got %s; expected %d", value, form[key][0])
	}
}
