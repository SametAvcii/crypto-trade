package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type API struct {
	EndPoint string `json:"end_point"`
}

func NewAPI(endpoint string) *API {
	return &API{
		EndPoint: endpoint,
	}
}

func (t *API) Get(path string, payload interface{}, response interface{}, headers ...map[string]string) error {
	url := t.EndPoint + path
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if len(headers) > 0 {
		for k, v := range headers[0] {
			req.Header.Set(k, v)
		}
	}
	return t.do(req, response)
}

func (t *API) do(req *http.Request, response interface{}) error {
	log.Println("Request URL:", req.URL.String()) // <--- ekle
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("BODY:", string(body)) // <--- ekle

	decode := json.NewDecoder(bytes.NewReader(body))
	decode.DisallowUnknownFields()
	decode.UseNumber()
	if err := decode.Decode(response); err != nil {
		return err
	}

	decode = json.NewDecoder(bytes.NewReader(body))
	decode.DisallowUnknownFields()
	decode.UseNumber()
	if err := decode.Decode(response); err != nil {
		return err
	}
	return nil
}
