package commands

import (
	"encoding/json"

	models "github.com/efvincent/archex5/models"
)

func UnmarshalAsTypedCommand(cmdTypeKey string, rawJson []byte) interface{} {
	switch cmdTypeKey {
	case "create-product":
		cmd := &CreateProductCmd{}
		json.Unmarshal(rawJson, cmd)
		return cmd
	case "update-product-attribs":
		cmd := &UpdateProductAttributesCmd{}
		json.Unmarshal(rawJson, cmd)
		return cmd
	case "update-product-images":
		cmd := &UpdateProductImagesCmd{}
		json.Unmarshal(rawJson, cmd)
		return cmd
	case "update-product-price":
		cmd := &UpdatePriceCmd{}
		json.Unmarshal(rawJson, cmd)
		return cmd
	case "product-headcheck":
		cmd := &HeadCheckCmd{}
		json.Unmarshal(rawJson, cmd)
		return cmd
	case "product-set-active":
		cmd := &SetActiveCmd{}
		json.Unmarshal(rawJson, cmd)
		return cmd
	}
	return nil
}

type ProductCmd struct {
	Namespace string `json:"ns" binding:"required"`
	Timestamp int    `json:"ts" binding:"required"`
	SKU       string `json:"sku" binding:"required"`
}

// A request to create a new product (Namespace + SKU) that explicitly does not exist -
// ie if the product exist this command fails. For product updates there are specific
// commands for the types of updates, see below
type CreateProductCmd struct {
	ProductCmd
	Source  string              `json:"source"`
	Product models.ProductModel `json:"product"`
}

// Used to update attributes on the product that do not require special
// handling or verification
type UpdateProductAttributesCmd struct {
	ProductCmd
	Title       string `json:"title"`
	Description string `json:"description"`
	Url         string `json:"url"`
}

type UpdateProductImagesCmd struct {
	ProductCmd
	Images        []string `json:"images"`
	PrimaryImgIdx int      `json:"primaryImgIdx"`
}

type UpdatePriceCmd struct {
	ProductCmd
	Version int64   `json:"version"`
	Price   float32 `json:"price"`
}

type HeadCheckCmd struct {
	ProductCmd
	Reason string `json:"reason"`
}

type SetActiveCmd struct {
	ProductCmd
	Active bool `json:"active"`
}
