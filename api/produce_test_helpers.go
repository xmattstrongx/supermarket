package api

import "github.com/xmattstrongx/supermarket/models"

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
