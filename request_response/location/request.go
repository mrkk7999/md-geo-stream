package location

import "time"

type LocationReq struct {
	ID        int       `json:"id"`
	TenantID  string    `json:"tenant_id"`
	UserID    string    `json:"user_id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Timestamp time.Time `json:"timestamp"`
}

func (LocationReq) TableName() string {
	return "location_data"
}
