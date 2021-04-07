package models

type ProductModel struct {
	Namespace     string   `json:"ns" binding:"required"`
	SequenceNum   int64    `json:"sequenceNum" binding:"required"`
	SKU           string   `json:"sku" binding:"required"`
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	Images        []string `json:"images"`
	PrimaryImgIdx int      `json:"primaryImgIdx"`
	Url           string   `json:"url"`
	IsContraband  bool     `json:"is_contraband"`
	Price         float32  `json:"price"`
	LastPrice     float32  `json:"lastPrice"`
}

func MakeProduct() ProductModel {
	return ProductModel{
		SequenceNum:  0,
		SKU:          "UNK",
		Title:        "Product-0",
		IsContraband: false,
		Price:        0.0,
		LastPrice:    0.0,
	}
}
