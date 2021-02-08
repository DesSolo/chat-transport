package vk

import (
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	apiURL     = "https://api.vk.com/method/"
	apiVersion = "5.52"
)

// CallMethod ...
func (c *Client) callMethod(method string, p RequestParams) ([]byte, error) {
	params, err := p.URLValues()
	if err != nil {
		return nil, err
	}

	params.Set("access_token", c.accessToken)

	if params.Get("v") == "" {
		params.Set("v", apiVersion)
	}

	req, err := http.NewRequest(http.MethodPost, apiURL+method, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(resp.Body)
}
