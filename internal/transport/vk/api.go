package vk

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/net/html/charset"
)

const (
	apiURL     = "https://api.vk.com/method/"
	apiVersion = "5.130"
)

// CallMethod ...
func (c *Client) callMethod(method string, p RequestParams) ([]byte, error) {
	params, err := p.URLValues()
	if err != nil {
		return nil, err
	}

	params.Set("access_token", c.accessToken)
	params.Set("v", apiVersion)

	req, err := http.NewRequest(http.MethodPost, apiURL+method, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("not valid status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	utf8, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(utf8)
}
