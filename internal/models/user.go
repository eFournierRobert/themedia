package models

type UserPost struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserPostResponse struct {
	UUID string `json:"uuid"`
	Username string `json:"username"`
}
