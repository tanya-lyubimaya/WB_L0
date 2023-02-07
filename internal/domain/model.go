package domain

import "time"

type Order struct {
	OrderUid    string `json:"order_uid" validate:"required,min=19,max=19"`
	TrackNumber string `json:"track_number" validate:"required,min=14,max=14"`
	Entry       string `json:"entry" validate:"required,max=10"`
	Delivery    struct {
		Name    string `json:"name" validate:"required,max=30"`
		Phone   string `json:"phone" validate:"required,min=11,max=14"`
		Zip     string `json:"zip" validate:"required,min=5,max=7"`
		City    string `json:"city" validate:"required,min=2,max=30"`
		Address string `json:"address" validate:"required,max=30"`
		Region  string `json:"region" validate:"required,max=20"`
		Email   string `json:"email" validate:"required,max=30"`
	} `json:"delivery" validate:"required"`
	Payment struct {
		Transaction  string `json:"transaction" validate:"required"`
		RequestId    string `json:"request_id"`
		Currency     string `json:"currency" validate:"required"`
		Provider     string `json:"provider" validate:"required"`
		Amount       int    `json:"amount" validate:"required"`
		PaymentDt    int    `json:"payment_dt" validate:"required"`
		Bank         string `json:"bank" validate:"required"`
		DeliveryCost int    `json:"delivery_cost" validate:"required"`
		GoodsTotal   int    `json:"goods_total" validate:"required"`
		CustomFee    int    `json:"custom_fee" validate:"gte=0"`
	} `json:"payment" validate:"required"`
	Items []struct {
		ChrtId      int    `json:"chrt_id" validate:"required"`
		TrackNumber string `json:"track_number" validate:"required,min=14,max=14"`
		Price       int    `json:"price" validate:"required"`
		Rid         string `json:"rid" validate:"required"`
		Name        string `json:"name" validate:"required"`
		Sale        int    `json:"sale" validate:"required"`
		Size        string `json:"size" validate:"required"`
		TotalPrice  int    `json:"total_price" validate:"required"`
		NmId        int    `json:"nm_id" validate:"required"`
		Brand       string `json:"brand" validate:"required"`
		Status      int    `json:"status" validate:"required"`
	} `json:"items" validate:"required"`
	Locale            string    `json:"locale" validate:"oneof=ru en"`
	InternalSignature string    `json:"internal_signature" validate:"max=20"`
	CustomerId        string    `json:"customer_id" validate:"required,max=20"`
	DeliveryService   string    `json:"delivery_service" validate:"required,max=20"`
	Shardkey          string    `json:"shardkey" validate:"required,max=10"`
	SmId              int       `json:"sm_id" validate:"required,gte=0,lte=10000"`
	DateCreated       time.Time `json:"date_created" format:"2006-01-02T06:22:19Z" validate:"required"`
	OofShard          string    `json:"oof_shard" validate:"required,max=10"`
}
