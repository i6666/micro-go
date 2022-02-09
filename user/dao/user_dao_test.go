package dao

import "testing"

func TestUserDaoImpl_SelectByEmail(t *testing.T) {
	userDAO := &UserDaoImpl{}
	err := InitMysql("180.76.118.233", "3306", "root", "strong", "user")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	user, err := userDAO.SelectByEmail("cc@mmail.cc")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("result uesrname is %s", user.Username)

}

var userDAO = &UserDaoImpl{}

func TestUserDaoImpl_Save(t *testing.T) {

	userDAO := &UserDaoImpl{}

	err := InitMysql("180.76.118.233", "3306", "root", "strong", "user")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	user := &UserEntity{
		Username: "strong",
		Password: "strong",
		Email:    "cc@mmail.cc",
	}
	err = userDAO.Save(user)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("new User ID is %d", user.ID)
}
