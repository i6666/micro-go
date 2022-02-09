package service

import (
	"context"
	"github.com/i6666/micro-go/user/dao"
	"github.com/i6666/micro-go/user/redis"
	"testing"
)

func TestUserServiceImpl_Register(t *testing.T) {

	err := dao.InitMysql("180.76.118.233", "3306", "root", "strong", "user")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	err = redis.InitRedis("jlili.cn", "6379", "strong")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	userService := &UserServiceImpl{
		&dao.UserDaoImpl{},
	}

	user, err := userService.Register(context.Background(), &RegisterUserVo{
		Username: "strong23",
		Email:    "strong33@mmail.cc",
		Password: "strong",
	})

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Logf("user id is %d", user.ID)

}

func TestUserServiceImpl_Login(t *testing.T) {

	err := dao.InitMysql("180.76.118.233", "3306", "root", "strong", "user")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	err = redis.InitRedis("jlili.cn", "6379", "strong")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	userService := &UserServiceImpl{
		&dao.UserDaoImpl{},
	}

	user, err := userService.Login(context.Background(), "cc@mmail.cc", "strong")

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Logf("user id is %d", user.ID)

}
