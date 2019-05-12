package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/xmattstrongx/supermarket/models"
)

func Test_isValidProduceCode(t *testing.T) {
	type args struct {
		produceCode string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid product code",
			args: args{
				produceCode: "A12T-4GH7-QPL9-3N4M",
			},
			want: true,
		},
		{
			name: "valid product code is case insensitive",
			args: args{
				produceCode: "a12T-4Gh7-QpL9-3n4m",
			},
			want: true,
		},
		{
			name: "invalid product code only three sets",
			args: args{
				produceCode: "a12T-4Gh7-QpL9-",
			},
			want: false,
		},
		{
			name: "invalid product code special character !",
			args: args{
				produceCode: "a12T-4Gh7-QpL9-ALL!",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidProduceCode(tt.args.produceCode); got != tt.want {
				t.Errorf("isValidProduceCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sortProduce(t *testing.T) {
	type args struct {
		queryParams queryParameters
	}
	tests := []struct {
		name string
		args args
		want []models.Produce
	}{
		{
			name: "sort by name",
			args: args{
				queryParameters{
					sortBy: "name",
				},
			},
			want: produceSortedByNameAsc(),
		},
		{
			name: "sort by name descending",
			args: args{
				queryParameters{
					sortBy: "name",
					order:  "desc",
				},
			},
			want: produceSortedByNameDesc(),
		},
		{
			name: "sort by price",
			args: args{
				queryParameters{
					sortBy: "unitPrice",
				},
			},
			want: produceSortedByPriceAsc(),
		},
		{
			name: "sort by price descending",
			args: args{
				queryParameters{
					sortBy: "unitPrice",
					order:  "desc",
				},
			},
			want: produceSortedByPriceDesc(),
		},
		{
			name: "sort by produceCode",
			args: args{
				queryParameters{
					sortBy: "produceCode",
				},
			},
			want: produceSortedByProduceCodeAsc(),
		},
		{
			name: "sort by produceCode descending",
			args: args{
				queryParameters{
					sortBy: "produceCode",
					order:  "desc",
				},
			},
			want: produceSortedByProduceCodeDesc(),
		},
		{
			name: "sort by name limit 2 offset 1",
			args: args{
				queryParameters{
					sortBy: "name",
					offset: "1",
					limit:  "2",
				},
			},
			want: produceSortedByNameLimit2Offset1(),
		},
		{
			name: "sort by price desc limit 2 offset 0",
			args: args{
				queryParameters{
					sortBy: "unitPrice",
					order:  "desc",
					offset: "0",
					limit:  "2",
				},
			},
			want: produceSortedByPriceLimit2Offset0Desc(),
		},
		{
			name: "sort by produceCode desc offset 1",
			args: args{
				queryParameters{
					sortBy: "produceCode",
					order:  "desc",
					offset: "1",
				},
			},
			want: produceSortedByProductCodeOffset1Desc(),
		},
		{
			name: "limit greater than length",
			args: args{
				queryParameters{
					limit: string(len(generateTestProduceData()) + 1),
				},
			},
			want: generateTestProduceData(),
		},
		{
			name: "offset greater than length",
			args: args{
				queryParameters{
					offset: string(len(generateTestProduceData())),
				},
			},
			want: generateTestProduceData(),
		},
		{
			name: "offset greater than length limit 2",
			args: args{
				queryParameters{
					offset: string(len(generateTestProduceData())),
					limit:  "2",
				},
			},
			want: produceOffsetOverLengthLimit2(),
		},
		{
			name: "offset is negative",
			args: args{
				queryParameters{
					offset: "-1",
				},
			},
			want: generateTestProduceData(),
		},
		{
			name: "limit is negative",
			args: args{
				queryParameters{
					limit: "-1",
				},
			},
			want: generateTestProduceData(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sortProduce(generateTestProduceData(), tt.args.queryParams); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sortProduce() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestListProduce(t *testing.T) {
	s := NewServer()
	req, err := http.NewRequest(http.MethodGet, "/produce", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.ListProduce)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	produce := []models.Produce{}
	if err := json.Unmarshal([]byte(rr.Body.String()), &produce); err != nil {
		t.Errorf("Failed to unmarshal response: %s", err)
	}

	if len(produce) != 4 {
		t.Errorf("response array has unexpected length got %d want %d", len(produce), 4)
	}

	contains := func(list []models.Produce, find models.Produce) bool {
		for _, val := range list {
			if val == find {
				return true
			}
		}
		return false
	}

	if ok := contains(
		produce,
		models.Produce{
			Name:        "Peach",
			ProduceCode: "E5T6-9UI3-TH15-QR88",
			UnitPrice:   2.99,
		},
	); !ok {
		t.Errorf("response array does not contain expected field")
	}

	if ok := contains(
		produce,
		models.Produce{
			Name:        "Green Pepper",
			ProduceCode: "YRT6-72AS-K736-L4AR",
			UnitPrice:   0.79,
		},
	); !ok {
		t.Errorf("response array does not contain expected field")
	}

	if ok := contains(
		produce,
		models.Produce{
			Name:        "Gala Apple",
			ProduceCode: "TQ4C-VV6T-75ZX-1RMR",
			UnitPrice:   3.59,
		},
	); !ok {
		t.Errorf("response array does not contain expected field")
	}

	if ok := contains(
		produce,
		models.Produce{
			Name:        "Lettuce",
			ProduceCode: "A12T-4GH7-QPL9-3N4M",
			UnitPrice:   3.46,
		},
	); !ok {
		t.Errorf("response array does not contain expected field")
	}

}

func TestListProduceSortByName(t *testing.T) {
	s := NewServer()
	req, err := http.NewRequest(http.MethodGet, "/produce?sort_by=name", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.ListProduce)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `[{"name":"Gala Apple","produceCode":"TQ4C-VV6T-75ZX-1RMR","unitPrice":3.59},{"name":"Green Pepper","produceCode":"YRT6-72AS-K736-L4AR","unitPrice":0.79},{"name":"Lettuce","produceCode":"A12T-4GH7-QPL9-3N4M","unitPrice":3.46},{"name":"Peach","produceCode":"E5T6-9UI3-TH15-QR88","unitPrice":2.99}]`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestListProduceSortByProduceCode(t *testing.T) {
	s := NewServer()

	req, err := http.NewRequest(http.MethodGet, "/produce?sort_by=produceCode", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.ListProduce)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `[{"name":"Lettuce","produceCode":"A12T-4GH7-QPL9-3N4M","unitPrice":3.46},{"name":"Peach","produceCode":"E5T6-9UI3-TH15-QR88","unitPrice":2.99},{"name":"Gala Apple","produceCode":"TQ4C-VV6T-75ZX-1RMR","unitPrice":3.59},{"name":"Green Pepper","produceCode":"YRT6-72AS-K736-L4AR","unitPrice":0.79}]`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestListProduceSortByUnitPrice(t *testing.T) {
	s := NewServer()

	req, err := http.NewRequest(http.MethodGet, "/produce?sort_by=unitPrice", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.ListProduce)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `[{"name":"Green Pepper","produceCode":"YRT6-72AS-K736-L4AR","unitPrice":0.79},{"name":"Peach","produceCode":"E5T6-9UI3-TH15-QR88","unitPrice":2.99},{"name":"Lettuce","produceCode":"A12T-4GH7-QPL9-3N4M","unitPrice":3.46},{"name":"Gala Apple","produceCode":"TQ4C-VV6T-75ZX-1RMR","unitPrice":3.59}]`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestCreateProduce(t *testing.T) {
	s := NewServer()

	req, err := http.NewRequest(http.MethodPost, "/produce", bytes.NewBuffer([]byte(`[{"name":"fumanchu","produceCode":"XX1X-4GH7-QPL9-3N4M","unitPrice": 1.13333}]`)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.CreateProduce)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	// Check the response body is what we expect.
	expected := `{"created":[{"name":"fumanchu","produceCode":"XX1X-4GH7-QPL9-3N4M","unitPrice":1.13}],"createFailed":[]}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestCreateProduceBadRequest(t *testing.T) {
	s := NewServer()

	req, err := http.NewRequest(http.MethodPost, "/produce", bytes.NewBuffer([]byte(`[{"name":"fumanchu","produceCode":"invalid","unitPrice": 1.13333}]`)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.CreateProduce)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	// Check the response body is what we expect.
	expected := `{"created":[],"createFailed":[{"name":"fumanchu","produceCode":"invalid","unitPrice":1.13333}]}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestCreateProduceValidAndInvalidRequests(t *testing.T) {
	s := NewServer()

	req, err := http.NewRequest(http.MethodPost, "/produce", bytes.NewBuffer([]byte(`[{"name":"fumanchu","produceCode":"XX1X-4GH7-QPL9-3N4M","unitPrice": 1.13333},{"name":"fumanchu","produceCode":"XX1X-4GH7-QPL9-","unitPrice": 1.00}]`)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.CreateProduce)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusMultiStatus {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	// Check the response body is what we expect.
	expected := `{"created":[{"name":"fumanchu","produceCode":"XX1X-4GH7-QPL9-3N4M","unitPrice":1.13}],"createFailed":[{"name":"fumanchu","produceCode":"XX1X-4GH7-QPL9-","unitPrice":1}]}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestDeleteProduce(t *testing.T) {
	s := NewServer()

	req, err := http.NewRequest(http.MethodDelete, "/produce/A12T-4GH7-QPL9-3N4M", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.DeleteProduce)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	// Check the response body is what we expect.
	expected := ``
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
