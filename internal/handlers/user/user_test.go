package user_handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	dbmodels "github.com/eFournierRobert/themedia/internal/models/db"
	jsonmodels "github.com/eFournierRobert/themedia/internal/models/json"
	"github.com/eFournierRobert/themedia/internal/tools"
	ban_tools "github.com/eFournierRobert/themedia/internal/tools/ban"
	init_tools "github.com/eFournierRobert/themedia/internal/tools/init"
	user_tools "github.com/eFournierRobert/themedia/internal/tools/user"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPostValidUser(t *testing.T) {
	assert := assert.New(t)
	teardownSuite := init_tools.SetupDatabase(t)
	defer teardownSuite(t)

	router, recorder := setupRouterAndRecorder()

	user := jsonmodels.UserPost{
		Username: "username",
		Password: "password",
		Role:     "user",
	}

	jsonUser, err := json.Marshal(&user)
	if err != nil {
		t.Errorf("Coudln't encode user to json")
	}

	router.ServeHTTP(
		recorder,
		httptest.NewRequest("POST", "/u", strings.NewReader(string(jsonUser))),
	)
	assert.Equal(http.StatusCreated, recorder.Code, "HTTP code should be created")

	var userCheck dbmodels.User
	tools.DB.Where("username = ?", user.Username).First(&userCheck)
	assert.NotEqual(0, userCheck.ID, "User was not found in the database")
}

func TestPostUserWithEmptyPassword(t *testing.T) {
	assert := assert.New(t)
	teardownSuite := init_tools.SetupDatabase(t)
	defer teardownSuite(t)

	router, recorder := setupRouterAndRecorder()

	user := jsonmodels.UserPost{
		Username: "username",
		Role:     "user",
	}

	jsonUser, err := json.Marshal(&user)
	if err != nil {
		t.Errorf("Coudln't encode user to json")
	}

	router.ServeHTTP(
		recorder,
		httptest.NewRequest("POST", "/u", strings.NewReader(string(jsonUser))),
	)
	assert.Equal(http.StatusBadRequest, recorder.Code, "HTTP code should be bad request")
}

func TestPostUserWithEmptyUsername(t *testing.T) {
	assert := assert.New(t)
	teardownSuite := init_tools.SetupDatabase(t)
	defer teardownSuite(t)

	router, recorder := setupRouterAndRecorder()

	user := jsonmodels.UserPost{
		Password: "password",
		Role:     "user",
	}

	jsonUser, err := json.Marshal(&user)
	if err != nil {
		t.Errorf("Coudln't encode user to json")
	}

	router.ServeHTTP(
		recorder,
		httptest.NewRequest("POST", "/u", strings.NewReader(string(jsonUser))),
	)
	assert.Equal(http.StatusBadRequest, recorder.Code, "HTTP code should be bad request")
}

func TestGetValidUserWithUUID(t *testing.T) {
	assert := assert.New(t)
	teardownSuite := init_tools.SetupDatabase(t)
	defer teardownSuite(t)

	router, recorder := setupRouterAndRecorder()

	uuid := "de0c8142-5973-478b-9287-37ff25e4e332"
	router.ServeHTTP(
		recorder,
		httptest.NewRequest("GET", "/u/"+uuid, nil),
	)
	assert.Equal(http.StatusFound, recorder.Code, "HTTP code should be found")

	dbUser, err := user_tools.FindFullUserByUUID(&uuid)
	if err != nil {
		t.Errorf("Got error %s", err.Error())
	}

	var user jsonmodels.UserResponse
	json.Unmarshal(recorder.Body.Bytes(), &user)

	assert.Equal(dbUser.UserUUID, user.UUID, "UUIDs should be equals")
	assert.Equal(dbUser.Username, user.Username, "Usernames should be equals")
	assert.Equal(dbUser.RoleUUID, user.Role.UUID, "Roles UUID should be equals")
	assert.Equal(dbUser.Name, user.Role.Name, "Roles name should be equals")
}

func TestGetUserWithInvalidUUID(t *testing.T) {
	assert := assert.New(t)
	teardownSuite := init_tools.SetupDatabase(t)
	defer teardownSuite(t)

	router, recorder := setupRouterAndRecorder()

	uuid := "de0c8142-5973-478b"
	router.ServeHTTP(
		recorder,
		httptest.NewRequest("GET", "/u/"+uuid, nil),
	)
	assert.Equal(http.StatusBadRequest, recorder.Code, "HTTP code should be bad request")
}

func TestGetInvalidUserWithUUID(t *testing.T) {
	assert := assert.New(t)
	teardownSuite := init_tools.SetupDatabase(t)
	defer teardownSuite(t)

	router, recorder := setupRouterAndRecorder()

	uuid := "90e51714-df99-4dd5-8408-3bb65cb0da00"
	router.ServeHTTP(
		recorder,
		httptest.NewRequest("GET", "/u/"+uuid, nil),
	)
	assert.Equal(http.StatusNotFound, recorder.Code, "HTTP code should be not found")
}

func TestPostValidLogin(t *testing.T) {
	assert := assert.New(t)
	teardownSuite := init_tools.SetupDatabase(t)
	defer teardownSuite(t)

	router, recorder := setupRouterAndRecorder()
	dbNewUser := createTestUserInDatabase()
	user := jsonmodels.UserPost{
		UUID:     dbNewUser.UUID,
		Username: dbNewUser.Username,
		Password: "password",
	}
	jsonUser, err := json.Marshal(user)
	if err != nil {
		t.Errorf("Couldn't turn the user into json. Got %s", err.Error())
	}

	router.ServeHTTP(
		recorder,
		httptest.NewRequest("POST", "/u/login", strings.NewReader(string(jsonUser))),
	)

	assert.Equal(http.StatusOK, recorder.Code, "HTTP code should be OK")
	assert.NotNil(recorder.Result().Cookies()[0])
}

func TestPostLoginWithBannedUser(t *testing.T) {
	assert := assert.New(t)
	teardownSuite := init_tools.SetupDatabase(t)
	defer teardownSuite(t)

	router, recorder := setupRouterAndRecorder()
	dbNewUser := createTestUserInDatabase()
	user := jsonmodels.UserPost{
		UUID:     dbNewUser.UUID,
		Username: dbNewUser.Username,
		Password: "password",
	}
	jsonUser, err := json.Marshal(user)
	if err != nil {
		t.Errorf("Couldn't turn the user into json. Got %s", err.Error())
	}

	ban_tools.CreateBan(dbNewUser.UUID, time.Now().Add(time.Hour))

	router.ServeHTTP(
		recorder,
		httptest.NewRequest("POST", "/u/login", strings.NewReader(string(jsonUser))),
	)

	assert.Equal(http.StatusUnauthorized, recorder.Code, "HTTP code should be unauthorized")
}

func TestPostLoginWithDeletedUser(t *testing.T) {
	assert := assert.New(t)
	teardownSuite := init_tools.SetupDatabase(t)
	defer teardownSuite(t)

	router, recorder := setupRouterAndRecorder()
	user := jsonmodels.UserPost{
		Username: "deleted",
	}
	jsonUser, err := json.Marshal(user)
	if err != nil {
		t.Errorf("Couldn't turn the user into json. Got %s", err.Error())
	}

	router.ServeHTTP(
		recorder,
		httptest.NewRequest("POST", "/u/login", strings.NewReader(string(jsonUser))),
	)

	assert.Equal(http.StatusUnauthorized, recorder.Code, "HTTP code should be unauthorized")
}

func TestPostLogout(t *testing.T) {
	assert := assert.New(t)
	teardownSuite := init_tools.SetupDatabase(t)
	defer teardownSuite(t)

	router, recorder := setupRouterAndRecorder()

	router.ServeHTTP(
		recorder,
		httptest.NewRequest("POST", "/u/logout", nil),
	)

	assert.Equal(http.StatusOK, recorder.Code, "HTTP code should be OK")
	assert.Equal("Authorization", recorder.Result().Cookies()[0].Name)
	assert.Empty(recorder.Result().Cookies()[0].Value)
	assert.Greater(time.Now(), recorder.Result().Cookies()[0].Expires)
}

func TestDeleteUserWithValidUser(t *testing.T) {
	assert := assert.New(t)
	teardownSuite := init_tools.SetupDatabase(t)
	defer teardownSuite(t)

	router, recorder := setupRouterAndRecorder()
	adminCookie := getAdminAuthCookie(router)

	uuidToDelete := "35ad671e-0fa0-4829-ae8e-37043d95fc33"
	req := httptest.NewRequest("DELETE", "/u/"+uuidToDelete, nil)
	req.AddCookie(adminCookie)
	router.ServeHTTP(
		recorder,
		req,
	)

	assert.Equal(http.StatusOK, recorder.Code, "HTTP code should be OK")

	var deletedUser dbmodels.User
	tools.DB.Where("uuid = ?", uuidToDelete).First(&deletedUser)

	assert.Zero(deletedUser.ID, "Deleted user should be deleted, but was found in the database")
}

func TestDeleteWithInvalidUser(t *testing.T) {
	assert := assert.New(t)
	teardownSuite := init_tools.SetupDatabase(t)
	defer teardownSuite(t)

	router, recorder := setupRouterAndRecorder()
	userCookie := getUserAuthCookie(router)

	uuidToDelete := "35ad671e-0fa0-4829-ae8e-37043d95fc33"
	req := httptest.NewRequest("DELETE", "/u/"+uuidToDelete, nil)
	req.AddCookie(userCookie)
	router.ServeHTTP(
		recorder,
		req,
	)

	assert.Equal(http.StatusUnauthorized, recorder.Code, "HTTP code should be unauthorized")

	var user dbmodels.User
	tools.DB.Where("uuid = ?", uuidToDelete).First(&user)

	assert.NotZero(user.ID, "User shouldn't have been deleted from database")
}

func TestPutUserWithValidUser(t *testing.T) {
	assert := assert.New(t)
	teardownSuite := init_tools.SetupDatabase(t)
	defer teardownSuite(t)

	router, recorder := setupRouterAndRecorder()
	userCookie := getUserAuthCookie(router)
	var user dbmodels.User
	tools.DB.Last(&user)

	modifiedJsonUser := jsonmodels.UserPost{
		Username: "New username",
		Bio:      user.Bio,
	}
	jsonUser, _ := json.Marshal(modifiedJsonUser)

	uuidToModify := user.UUID
	req := httptest.NewRequest("PUT", "/u/"+uuidToModify, strings.NewReader(string(jsonUser)))
	req.AddCookie(userCookie)
	router.ServeHTTP(
		recorder,
		req,
	)

	assert.Equal(http.StatusOK, recorder.Code, "HTTP code should be OK")

	tools.DB.Where("uuid = ?", user.UUID).First(&user)
	assert.Equal(modifiedJsonUser.Username, user.Username, "New username has not been set")
	assert.Equal(modifiedJsonUser.Bio, user.Bio, "Bio shouldn't have changed")
}

func TestPutUserWithNewRoleAndInvalidUser(t *testing.T) {
	assert := assert.New(t)
	teardownSuite := init_tools.SetupDatabase(t)
	defer teardownSuite(t)

	router, recorder := setupRouterAndRecorder()
	userCookie := getUserAuthCookie(router)
	var user dbmodels.User
	tools.DB.Last(&user)
	adminRole := getAdminRole()

	modifiedJsonUser := jsonmodels.UserPost{
		Role: adminRole.Name,
	}
	jsonUser, _ := json.Marshal(modifiedJsonUser)

	uuidToModify := user.UUID
	req := httptest.NewRequest("PUT", "/u/"+uuidToModify, strings.NewReader(string(jsonUser)))
	req.AddCookie(userCookie)
	router.ServeHTTP(
		recorder,
		req,
	)

	assert.Equal(http.StatusForbidden, recorder.Code, "HTTP code should be forbidden")

	tools.DB.Where("uuid = ?", user.UUID).First(&user)
	assert.NotEqual(adminRole.ID, user.RoleID, "Role should not have been updated in the database")
}

func TestPutUserWithNewRoleAndValidUser(t *testing.T) {
	assert := assert.New(t)
	teardownSuite := init_tools.SetupDatabase(t)
	defer teardownSuite(t)

	router, recorder := setupRouterAndRecorder()
	adminCookie := getAdminAuthCookie(router)
	user := createTestUserInDatabase()
	adminRole := getAdminRole()

	modifiedJsonUser := jsonmodels.UserPost{
		Role: adminRole.Name,
	}
	jsonUser, _ := json.Marshal(modifiedJsonUser)

	uuidToModify := user.UUID
	req := httptest.NewRequest("PUT", "/u/"+uuidToModify, strings.NewReader(string(jsonUser)))
	req.AddCookie(adminCookie)
	router.ServeHTTP(
		recorder,
		req,
	)

	assert.Equal(http.StatusOK, recorder.Code, "HTTP code should be OK")

	tools.DB.Where("uuid = ?", user.UUID).First(&user)
	assert.Equal(adminRole.ID, user.RoleID, "Role have been updated in the database")
}

func TestPutUserWithNewPasswordAndValidUser(t *testing.T) {
	assert := assert.New(t)
	teardownSuite := init_tools.SetupDatabase(t)
	defer teardownSuite(t)

	router, recorder := setupRouterAndRecorder()
	userCookie := getUserAuthCookie(router)
	var user dbmodels.User
	tools.DB.Last(&user)

	modifiedJsonUser := jsonmodels.UserPost{
		Password: "new password",
	}
	jsonUser, _ := json.Marshal(modifiedJsonUser)

	uuidToModify := user.UUID
	req := httptest.NewRequest("PUT", "/u/"+uuidToModify, strings.NewReader(string(jsonUser)))
	req.AddCookie(userCookie)
	router.ServeHTTP(
		recorder,
		req,
	)

	assert.Equal(http.StatusOK, recorder.Code, "HTTP code should be OK")

	var updatedUser dbmodels.User
	tools.DB.Where("uuid = ?", user.UUID).First(&updatedUser)
	assert.NotEqual(user.PasswordHash, updatedUser.PasswordHash, "Password should have been updated in the database")
}

func TestPutUserWithAnotherUser(t *testing.T) {
	assert := assert.New(t)
	teardownSuite := init_tools.SetupDatabase(t)
	defer teardownSuite(t)

	router, recorder := setupRouterAndRecorder()
	userCookie := getUserAuthCookie(router)
	var user dbmodels.User
	tools.DB.Last(&user)
	otherUser := createTestUserInDatabase()

	modifiedJsonUser := jsonmodels.UserPost{
		Bio: "new bio",
	}
	jsonUser, _ := json.Marshal(modifiedJsonUser)

	uuidToModify := otherUser.UUID
	req := httptest.NewRequest("PUT", "/u/"+uuidToModify, strings.NewReader(string(jsonUser)))
	req.AddCookie(userCookie)
	router.ServeHTTP(
		recorder,
		req,
	)

	assert.Equal(http.StatusUnauthorized, recorder.Code, "HTTP code should be unauthorized")

	tools.DB.Where("uuid = ?", otherUser.UUID).First(&otherUser)
	assert.NotEqual(otherUser.Bio, modifiedJsonUser.Bio, "Bio should not have been updated in the database")
}

func TestPostBanWithValidUser(t *testing.T) {
	assert := assert.New(t)
	teardownSuite := init_tools.SetupDatabase(t)
	defer teardownSuite(t)

	router, recorder := setupRouterAndRecorder()
	adminCookie := getAdminAuthCookie(router)
	banBody := jsonmodels.Ban{
		EndDatetime: time.Now().Add(time.Hour * 3),
	}
	jsonBan, _ := json.Marshal(banBody)

	uuidToBan := "35ad671e-0fa0-4829-ae8e-37043d95fc33"
	req := httptest.NewRequest("POST", "/u/"+uuidToBan+"/ban", strings.NewReader(string(jsonBan)))
	req.AddCookie(adminCookie)
	router.ServeHTTP(
		recorder,
		req,
	)

	assert.Equal(http.StatusOK, recorder.Code, "HTTP code should be OK")

	var ban dbmodels.Ban
	tools.DB.Last(&ban)
	assert.Equal(uint(3), ban.UserId, "User should be banned in the database")
}

func TestPostBanWithInvalidUser(t *testing.T) {
	assert := assert.New(t)
	teardownSuite := init_tools.SetupDatabase(t)
	defer teardownSuite(t)

	router, recorder := setupRouterAndRecorder()
	adminCookie := getUserAuthCookie(router)
	banBody := jsonmodels.Ban{
		EndDatetime: time.Now().Add(time.Hour * 3),
	}
	jsonBan, _ := json.Marshal(banBody)

	uuidToBan := "35ad671e-0fa0-4829-ae8e-37043d95fc33"
	req := httptest.NewRequest("POST", "/u/"+uuidToBan+"/ban", strings.NewReader(string(jsonBan)))
	req.AddCookie(adminCookie)
	router.ServeHTTP(
		recorder,
		req,
	)

	assert.Equal(http.StatusUnauthorized, recorder.Code, "HTTP code should be unauthorized")

	var ban dbmodels.Ban
	tools.DB.Last(&ban)
	assert.NotEqual(uint(3), ban.UserId, "User should not be banned in the database")
}

func getAdminAuthCookie(router *gin.Engine) *http.Cookie {
	loginRecorder := httptest.NewRecorder()
	dbNewAdmin := createTestAdminInDatabase()
	user := jsonmodels.UserPost{
		UUID:     dbNewAdmin.UUID,
		Username: dbNewAdmin.Username,
		Password: "password",
	}
	jsonUser, _ := json.Marshal(user)

	router.ServeHTTP(
		loginRecorder,
		httptest.NewRequest("POST", "/u/login", strings.NewReader(string(jsonUser))),
	)

	return loginRecorder.Result().Cookies()[0]
}

func getAdminRole() *dbmodels.Role {
	var role dbmodels.Role
	tools.DB.Where("name = ?", "admin").First(&role)
	return &role
}

func getUserAuthCookie(router *gin.Engine) *http.Cookie {
	loginRecorder := httptest.NewRecorder()
	dbNewUser := createTestUserInDatabase()
	user := jsonmodels.UserPost{
		UUID:     dbNewUser.UUID,
		Username: dbNewUser.Username,
		Password: "password",
	}
	jsonUser, _ := json.Marshal(user)

	router.ServeHTTP(
		loginRecorder,
		httptest.NewRequest("POST", "/u/login", strings.NewReader(string(jsonUser))),
	)

	return loginRecorder.Result().Cookies()[0]
}

func createTestUserInDatabase() *dbmodels.User {
	username := "the user"
	password := "password"
	var role dbmodels.Role
	tools.DB.Where("name = ?", "user").First(&role)
	newUser, _ := user_tools.CreateUser(&username, &password, &role)

	return newUser
}

func createTestAdminInDatabase() *dbmodels.User {
	username := "the user"
	password := "password"
	var role dbmodels.Role
	tools.DB.Where("name = ?", "admin").First(&role)
	newUser, _ := user_tools.CreateUser(&username, &password, &role)

	return newUser
}

func setupRouterAndRecorder() (*gin.Engine, *httptest.ResponseRecorder) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	AddEndpointsToRouter(router)

	return router, httptest.NewRecorder()
}
