package vk

import (
	"fmt"
	"net/url"
)

type responseMessage struct {
	Response struct {
		Count int
		Items []struct {
			ID        int
			Date      int
			Out       int
			UserID    int
			ReadState int
			Title     string
			Body      string
		}
	}
}

// RequestParams ...
type RequestParams map[string]interface{}

// URLValues ...
func (p RequestParams) URLValues() (url.Values, error) {
	values := url.Values{}

	for k, v := range p {
		values.Add(k, fmt.Sprint(v))
	}

	return values, nil
}
