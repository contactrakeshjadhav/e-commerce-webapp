package model

import "time"

const (
	StandartISOFormat = "2006-01-02T15:04:05-0700"
	SERVICE_NAME      = "product"
)

var (
	TokenCTXKey = CTXKey("Token")
)

type Base struct {
	ID ID

	CreatedAt time.Time
	UpdatedAt time.Time
}

type ID string
type CTXKey string
