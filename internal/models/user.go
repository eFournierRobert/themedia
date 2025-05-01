package models

type UserPost struct {
	Username string `json:"username"`
	Password string `json:"password"`
	RoleUUID string `json:"roleUuid"`
}

type UserResponse struct {
	UUID string `json:"uuid"`
	Username string `json:"username"`
	Role RoleResponse `json:"role"`
}

type RoleResponse struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}
