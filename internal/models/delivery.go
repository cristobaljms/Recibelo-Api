package models

import "gopkg.in/mgo.v2/bson"

const (
	DELIVERY_STATE_INVALID = iota
	DELIVERY_STATE_SENT
	DELIVERY_STATE_IN_TRANSIT
	DELIVERY_STATE_RECEIVED
)

type DeliveryState int64

type Delivery struct {
	ID               bson.ObjectId     `json:"id"                 bson:"_id,omitempty" `
	TrackerID        string            `json:"tracker_id"         bson:"tracker_id"`
	Rated            bool              `json:"rated"              bson:"rated"`
	UserID           bson.ObjectId     `json:"user_id"            bson:"user_id,omitempty"`
	Rating           int64             `json:"rating"             bson:"rating"`
	Courier          Courier           `json:"courier"            bson:"courier"`
	Description      string            `json:"description"        bson:"description"`
	DeliveryState    DeliveryState     `json:"delivery_state"     bson:"delivery_state"`
	DeliveryHistory  []DeliveryHistory `json:"delivery_history"   bson:"delivery_history"`
	EstimatedDueDate int64             `json:"estimated_due_date" bson:"estimated_due_date"`
}

type TeamDelivery struct {
	UserID    bson.ObjectId `json:"user_id"    bson:"user_id,omitempty"`
	TrackerID string        `json:"tracker_id" bson:"tracker_id"`
}

type DeliveryHistory struct {
	Timestamp     int64         `json:"timestamp"      bson:"timestamp"`
	Description   string        `json:"description"    bson:"description"`
	DeliveryState DeliveryState `json:"delivery_state" bson:"delivery_state"`
}

type Courier struct {
	Name    string `json:"name"     bson:"name"`
	IconURL string `json:"icon_url" bson:"icon_url"`
}
