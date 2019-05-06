package models

type Produce struct {
	Name        string  `json:"name"`
	ProduceCode string  `json:"produceCode"`
	UnitPrice   float64 `json:"unitPrice"`
}

type CreateProduceReponse struct {
	Created []Produce `json:"created"`
	Invalid []Produce `json:"createFailed"`
}
