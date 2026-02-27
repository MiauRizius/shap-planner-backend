package models

type RefreshToken struct {
	ID         string `json:id`
	UserID     string `json:userid`
	Token      string `json:token`
	ExpiresAt  int64  `json:expiresat`
	CreatedAt  int64  `json:createdat`
	Revoked    bool   `json:revoked`
	DeviceInfo string `json:deviceinfo`
}
