package vk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/net/html/charset"
)

const (
	apiURL     = "https://vk.com/dev"
	apiVersion = "5.130"
)

// CallMethod ...
func (c *Client) callMethod(method, hash string, p RequestParams) ([]byte, error) {
	for param := range p {
		p["param_"+param] = p[param]
		delete(p, param)
	}
	params, err := p.URLValues()
	if err != nil {
		return nil, err
	}

	params.Set("act", "a_run_method")
	params.Set("al", "1")
	params.Set("method", method)
	params.Set("hash", hash)
	params.Set("param_v", apiVersion)

	req, err := http.NewRequest(http.MethodPost, apiURL, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", c.authBasic))
	req.Header.Set("Cookie", fmt.Sprintf("remixsid=%s", c.cookieRemixsID))

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

func unmarshal(b []byte, v interface{}) error {
	var dr map[string]interface{}
	if err := json.Unmarshal(b, &dr); err != nil {
		return err
	}

	payload := dr["payload"].([]interface{})[1].([]interface{})[0].(string)
	if err := json.Unmarshal([]byte(payload), v); err != nil {
		return err
	}

	return nil
}
