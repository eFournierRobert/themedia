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
