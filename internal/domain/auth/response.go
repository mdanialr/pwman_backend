package auth

import "time"

// Response standard response that may be used in auth domain.
type Response struct {
	AccessToken string    `json:"access_token"`
	ExpiredAt   time.Time `json:"expired_at"`
}
