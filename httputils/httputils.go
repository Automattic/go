package httputils

import "net/url"

func GetFormValueByIndex(form url.Values, key string, index int, defaultValue string) string {
	if index < 0 {
		return defaultValue
	}

	if len(form[key]) <= index {
		return defaultValue
	}

	return form[key][index]
}
