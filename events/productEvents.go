package events

import "github.com/efvincent/archex5/models"

type Event struct {
	SeqNum    int64  `json:"seqNum"`
	Namespace string `json:"ns"`
	Timestamp int    `json:"ts"`
	EventType string `json:"string"`
}

type ProductCreated struct {
	Event
	Source  string              `json:"source"`
	Product models.ProductModel `json:"product"`
}

type AttribsUpdated struct {
	Event
	Title       string `json:"title"`
	Description string `json:"description"`
	Url         string `json:"url"`
}

type ImagesUpdated struct {
	Event
	Images        []string `json:"images"`
	PrimaryImgIdx int      `json:"primatyImgIdx"`
}

type PriceUpdated struct {
	Event
	Price float32 `json:"price"`
}

type HeadCheckPerformed struct {
	Event
	Reason  string `json:"reason"`
	Success bool   `json:"success"`
	Info    string `json:"info"`
}
