package events

import "github.com/efvincent/archex5/models"

// the constants in this format are used when deserializing (unmarshaling) events
// from the event store so we know what type of event is being read.
const ProductCreatedT = "prodCreated-1"

type ProductCreated struct {
	Namespace string               `json:"ns" binding:"required"`
	SKU       string               `json:"sku" binding:"required"`
	Source    string               `json:"source"`
	Product   *models.ProductModel `json:"product"`
}

const AttribsUpdatedT = "attrUpd-1"

type AttribsUpdated struct {
	Namespace   string `json:"ns" binding:"required"`
	SKU         string `json:"sku" binding:"required"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Url         string `json:"url"`
}

const ImagesUpdatedT = "imgUpd-1"

type ImagesUpdated struct {
	Namespace     string   `json:"ns" binding:"required"`
	SKU           string   `json:"sku" binding:"required"`
	Images        []string `json:"images"`
	PrimaryImgIdx int      `json:"primatyImgIdx"`
}

const PriceUpdatedT = "priceUpd-1"

type PriceUpdated struct {
	Namespace string  `json:"ns" binding:"required"`
	SKU       string  `json:"sku" binding:"required"`
	Price     float32 `json:"price"`
}

const HeadCheckPerformedT = "headcheck-1"

type HeadCheckPerformed struct {
	Namespace string `json:"ns" binding:"required"`
	SKU       string `json:"sku" binding:"required"`
	Reason    string `json:"reason"`
	Success   bool   `json:"success"`
	Info      string `json:"info"`
}

const ActiveStateSetT = "setactivestate-1"

type ActiveStateSet struct {
	Namespace string `json:"ns" binding:"required"`
	SKU       string `json:"sku" binding:"required"`
	Active    bool   `json:"active" binding:"required"`
}
