package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/xmattstrongx/supermarket/models"
)

type produceClient struct {
	client   *http.Client
	endpoint string
}

func newClient(opts ...func(*produceClient) error) (*produceClient, error) {
	produceClient := &produceClient{
		client: &http.Client{},
	}

	for _, opt := range opts {
		if err := opt(produceClient); err != nil {
			return nil, err
		}
	}

	return produceClient, nil
}

func withTimeout(timeout string) func(*produceClient) error {
	return func(p *produceClient) error {
		t, err := time.ParseDuration("10s")
		if err != nil {
			return err
		}
		p.client.Timeout = t
		return nil
	}
}

func withEndpoint(endpoint string) func(*produceClient) error {
	return func(p *produceClient) error {
		p.endpoint = endpoint
		return nil
	}
}

func (p *produceClient) listProduce(sortBy, order, limit, offset string) ([]models.Produce, int, error) {
	listProduceResponse := &[]models.Produce{}
	u, err := url.Parse(fmt.Sprintf("%s/api/v1/produce", p.endpoint))
	if err != nil {
		return nil, 0, err
	}
	q := u.Query()

	if sortBy != "" {
		q.Set("sort_by", sortBy)
	}

	if order != "" {
		q.Set("order", order)
	}

	if limit != "" {
		q.Set("limit", limit)
	}

	if offset != "" {
		q.Set("offset", offset)
	}

	u.RawQuery = q.Encode()
	fmt.Println(u)

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, 0, err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	if err := json.Unmarshal(bodyBytes, listProduceResponse); err != nil {
		return nil, 0, err
	}

	return *listProduceResponse, resp.StatusCode, nil
}

func (p *produceClient) deleteProduce(produceCode string) ([]byte, int, error) {
	url := fmt.Sprintf("%s/api/v1/produce/%s", p.endpoint, produceCode)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, 0, err
	}

	fmt.Println(url)

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, 0, err
	}

	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	return bodyBytes, resp.StatusCode, nil
}
