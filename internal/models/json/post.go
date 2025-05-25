package jsonmodels

type Post struct {
	UUID string `json:"uuid"`
	Title string `json:"title"`
	Body string `json:"body"`
	UserUUID string `json:"user_uuid"`
}
