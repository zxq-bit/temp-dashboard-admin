package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func NewHttpClient(timeout time.Duration) (*http.Client, error) {
	if timeout < 0 {
		return nil, fmt.Errorf("illegal timeout: %v", timeout)
	}
	return &http.Client{
		Timeout: timeout,
	}, nil
}

func doGet(c *http.Client, url string, expectedCode int, result interface{}) error {
	resp, e := c.Get(url)
	if e != nil {
		return e
	}
	defer resp.Body.Close()
	b, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		return e
	}
	if resp.StatusCode != expectedCode {
		return fmt.Errorf("unexpected code, %v != %v, %s", resp.StatusCode, expectedCode, string(b))
	}
	e = json.Unmarshal(b, result)
	if e != nil {
		return fmt.Errorf("unmarshal failed, [%v]'%s', %v", resp.StatusCode, string(b), e)
	}
	return nil
}
