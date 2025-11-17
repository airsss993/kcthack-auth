package v1

import "time"

type RegisterReq struct {
	FirstName string
	LastName  string
	Email     string
}

type LoginReq struct {
	Email    string
	Password string
}

type LoginResp struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    time.Time
}
