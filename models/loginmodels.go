package models

type RefreshToken struct {
	ID         string `json:"id"`
	UserID     string `json:"user_id"`
	Token      string `json:"token"`
	ExpiresAt  int64  `json:"expires_at"`
	CreatedAt  int64  `json:"created_at"`
	Revoked    bool   `json:"revoked"`
	DeviceInfo string `json:"device_info"`
}
