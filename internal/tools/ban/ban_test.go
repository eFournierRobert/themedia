package ban_tools

import (
	"testing"
	"time"

	dbmodels "github.com/eFournierRobert/themedia/internal/models/db"
	"github.com/eFournierRobert/themedia/internal/tools"
	init_tools "github.com/eFournierRobert/themedia/internal/tools/init"
)

func TestCreateBanWithValidUser(t *testing.T) {
	teardownSuite := init_tools.SetupDatabase(*t)
	defer teardownSuite(*t)

	err := CreateBan("35ad671e-0fa0-4829-ae8e-37043d95fc33", time.Now())
	if err != nil {
		t.Errorf("Couldn't ban user. Got error %s", err.Error())
	}
}

func TestCreateBanWithInvalidUser(t *testing.T) {
	teardownSuite := init_tools.SetupDatabase(*t)
	defer teardownSuite(*t)

	err := CreateBan("Fraise", time.Now())
	if err == nil {
		t.Errorf("Could create ban with invalid user UUID Fraise")
	}
}

func TestIsUserBanned(t *testing.T) {
	teardownSuite := init_tools.SetupDatabase(*t)
	defer teardownSuite(*t)

	b := IsUserBanned(*getDeletedUserUUID())

	if !b {
		t.Errorf("Banned user %s is not seen as banned", *getDeletedUserUUID())
	}
}

func getDeletedUserUUID() *string {
	var usr dbmodels.User
	tools.DB.Where("username = ?", "deleted").First(&usr)

	return &usr.UUID
}
