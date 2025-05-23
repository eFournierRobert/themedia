// Package models is the package that contains
// all the structs that will be returned in JSON.
package jsonmodels

// UserPost is the struct used to deserialize
// the request body for a new user.
type UserPost struct {
	UUID     string `json:"uuid"`
	Username string `json:"username"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
	Role     string `json:"role"`
}

// UserResponse is the struct used to serialize
// the user we want to return.
type UserResponse struct {
	UUID     string       `json:"uuid"`
	Username string       `json:"username"`
	Bio      string       `json:"bio"`
	Role     RoleResponse `json:"role"`
}

// RoleResponse is the struct used to serialize
// the role we want to return.
type RoleResponse struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}
