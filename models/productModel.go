package models

type PriceChange struct {
	RequestedPrice float32 `json:"requestedPrice"`
	Timestamp      int64   `json:"ts"`
}
type ProductModel struct {
	Namespace           string        `json:"ns" binding:"required"`
	SequenceNum         int64         `json:"sequenceNum" binding:"required"`
	SKU                 string        `json:"sku" binding:"required"`
	Title               string        `json:"title"`
	Description         string        `json:"description"`
	Images              []string      `json:"images"`
	PrimaryImgIdx       int           `json:"primaryImgIdx"`
	Url                 string        `json:"url"`
	IsContraband        bool          `json:"is_contraband"`
	IsActive            bool          `json:"is_active"`
	HeadCheckOk         bool          `json:"headCheckOK"`
	LastHeadCheck       int64         `json:"lastHeadCheck"`
	Price               float32       `json:"price"`
	PriceChangeRequests []PriceChange `json:"priceChanges"`
}

var SampleProduct = ProductModel{
	Namespace: "Nike",
	SKU:       "SHOE001",
	Title:     "Jordan Delta Breathe",
	Description: `Inspired by high-tech functionality and handmade craftsmanship, 
the Jordan Delta Breathe combines natural and synthetic materials.`,
	Images: []string{
		"https://static.nike.com/a/images/t_PDP_864_v1/f_auto,b_rgb:f5f5f5,q_80/91e12f77-89de-4ed1-8d26-81a5f80c508c/jordan-delta-breathe-mens-shoe-2ggX3h.jpg",
		"https://static.nike.com/a/images/t_PDP_864_v1/f_auto,b_rgb:f5f5f5,q_80/a0a287a1-3518-4230-9bd0-996b567c9019/jordan-delta-breathe-mens-shoe-2ggX3h.jpg",
	},
	PrimaryImgIdx: 0,
	SequenceNum:   0, // ignored when sending create command
	IsContraband:  false,
	IsActive:      true,
	Url:           "https://www.nike.com/t/jordan-delta-breathe-mens-shoe-2ggX3h/CW0783-901",
}
