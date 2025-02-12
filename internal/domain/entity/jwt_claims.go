package entity

import "time"

// AuthClaims Adapter for jwt pkg
type AuthClaims struct {
	Role       string    `json:"role"`
	DeviceID   string    `json:"device_id"`
	Subject    string    `json:"subject"`
	JTIDY      string    `json:"jwt_id"`
	ExpireTime time.Time `json:"expire_time"`
}
