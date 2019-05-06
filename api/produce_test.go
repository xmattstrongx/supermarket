package api

import (
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

func generateTestProduceData() []models.Produce {
	return []models.Produce{
		models.Produce{
			Name:        "Lettuce",
			ProduceCode: "A12T-4GH7-QPL9-3N4M",
			UnitPrice:   3.46,
		},
		models.Produce{
			Name:        "Peach",
			ProduceCode: "E5T6-9UI3-TH15-QR88",
			UnitPrice:   2.99,
		},
		models.Produce{
			Name:        "Green Pepper",
			ProduceCode: "YRT6-72AS-K736-L4AR",
			UnitPrice:   0.79,
		},
	}
}

func produceSortedByNameAsc() []models.Produce {
	return []models.Produce{
		models.Produce{
			Name:        "Green Pepper",
			ProduceCode: "YRT6-72AS-K736-L4AR",
			UnitPrice:   0.79,
		},
		models.Produce{
			Name:        "Lettuce",
			ProduceCode: "A12T-4GH7-QPL9-3N4M",
			UnitPrice:   3.46,
		},
		models.Produce{
			Name:        "Peach",
			ProduceCode: "E5T6-9UI3-TH15-QR88",
			UnitPrice:   2.99,
		},
	}
}

func produceSortedByNameDesc() []models.Produce {
	return []models.Produce{
		models.Produce{
			Name:        "Peach",
			ProduceCode: "E5T6-9UI3-TH15-QR88",
			UnitPrice:   2.99,
		},
		models.Produce{
			Name:        "Lettuce",
			ProduceCode: "A12T-4GH7-QPL9-3N4M",
			UnitPrice:   3.46,
		},
		models.Produce{
			Name:        "Green Pepper",
			ProduceCode: "YRT6-72AS-K736-L4AR",
			UnitPrice:   0.79,
		},
	}
}

func produceSortedByPriceAsc() []models.Produce {
	return []models.Produce{
		models.Produce{
			Name:        "Green Pepper",
			ProduceCode: "YRT6-72AS-K736-L4AR",
			UnitPrice:   0.79,
		},
		models.Produce{
			Name:        "Peach",
			ProduceCode: "E5T6-9UI3-TH15-QR88",
			UnitPrice:   2.99,
		},
		models.Produce{
			Name:        "Lettuce",
			ProduceCode: "A12T-4GH7-QPL9-3N4M",
			UnitPrice:   3.46,
		},
	}
}

func produceSortedByPriceDesc() []models.Produce {
	return []models.Produce{
		models.Produce{
			Name:        "Lettuce",
			ProduceCode: "A12T-4GH7-QPL9-3N4M",
			UnitPrice:   3.46,
		},
		models.Produce{
			Name:        "Peach",
			ProduceCode: "E5T6-9UI3-TH15-QR88",
			UnitPrice:   2.99,
		},
		models.Produce{
			Name:        "Green Pepper",
			ProduceCode: "YRT6-72AS-K736-L4AR",
			UnitPrice:   0.79,
		},
	}
}

func produceSortedByProduceCodeAsc() []models.Produce {
	return []models.Produce{
		models.Produce{
			Name:        "Lettuce",
			ProduceCode: "A12T-4GH7-QPL9-3N4M",
			UnitPrice:   3.46,
		},
		models.Produce{
			Name:        "Peach",
			ProduceCode: "E5T6-9UI3-TH15-QR88",
			UnitPrice:   2.99,
		},
		models.Produce{
			Name:        "Green Pepper",
			ProduceCode: "YRT6-72AS-K736-L4AR",
			UnitPrice:   0.79,
		},
	}
}

func produceSortedByProduceCodeDesc() []models.Produce {
	return []models.Produce{
		models.Produce{
			Name:        "Green Pepper",
			ProduceCode: "YRT6-72AS-K736-L4AR",
			UnitPrice:   0.79,
		},
		models.Produce{
			Name:        "Peach",
			ProduceCode: "E5T6-9UI3-TH15-QR88",
			UnitPrice:   2.99,
		},
		models.Produce{
			Name:        "Lettuce",
			ProduceCode: "A12T-4GH7-QPL9-3N4M",
			UnitPrice:   3.46,
		},
	}
}

func produceSortedByNameLimit2Offset1() []models.Produce {
	return []models.Produce{
		models.Produce{
			Name:        "Lettuce",
			ProduceCode: "A12T-4GH7-QPL9-3N4M",
			UnitPrice:   3.46,
		},
		models.Produce{
			Name:        "Peach",
			ProduceCode: "E5T6-9UI3-TH15-QR88",
			UnitPrice:   2.99,
		},
	}
}

func produceSortedByPriceLimit2Offset0Desc() []models.Produce {
	return []models.Produce{
		models.Produce{
			Name:        "Lettuce",
			ProduceCode: "A12T-4GH7-QPL9-3N4M",
			UnitPrice:   3.46,
		},
		models.Produce{
			Name:        "Peach",
			ProduceCode: "E5T6-9UI3-TH15-QR88",
			UnitPrice:   2.99,
		},
	}
}

func produceSortedByProductCodeOffset1Desc() []models.Produce {
	return []models.Produce{
		models.Produce{
			Name:        "Peach",
			ProduceCode: "E5T6-9UI3-TH15-QR88",
			UnitPrice:   2.99,
		},
		models.Produce{
			Name:        "Lettuce",
			ProduceCode: "A12T-4GH7-QPL9-3N4M",
			UnitPrice:   3.46,
		},
	}
}

func produceOffsetOverLengthLimit2() []models.Produce {
	return []models.Produce{
		models.Produce{
			Name:        "Lettuce",
			ProduceCode: "A12T-4GH7-QPL9-3N4M",
			UnitPrice:   3.46,
		},
		models.Produce{
			Name:        "Peach",
			ProduceCode: "E5T6-9UI3-TH15-QR88",
			UnitPrice:   2.99,
		},
	}
}
