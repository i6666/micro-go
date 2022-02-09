package service

import (
	"context"
	"errors"
	"github.com/i6666/micro-go/user/dao"
	"github.com/i6666/micro-go/user/redis"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

var (
	ErrUserExisted = errors.New("user is existed")
	ErrPassword    = errors.New("email and password are not match")
	ErrRegistering = errors.New("email is registering")
)

//用户信息
//StructTag 一般由一个或者多个键值对组成
type UserInfoDto struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type RegisterUserVo struct {
	Username string
	Password string
	Email    string
}

type UserService interface {
	//登录
	Login(ctx context.Context, email, password string) (*UserInfoDto, error)
	//注册
	Register(ctx context.Context, vo *RegisterUserVo) (*UserInfoDto, error)
}

type UserServiceImpl struct {
	useDAO dao.UserDao
}

func (userService *UserServiceImpl) Register(ctx context.Context, vo *RegisterUserVo) (*UserInfoDto, error) {
	lock := redis.GetRedisLock(vo.Email, time.Duration(5)*time.Second)
	err := lock.Lock()
	if err != nil {
		log.Printf("err: %s", err)
		return nil, ErrRegistering
	}
	defer lock.Unlock()

	existUser, err := userService.useDAO.SelectByEmail(vo.Email)
	if (err == nil && existUser == nil) || err == gorm.ErrRecordNotFound {
		newUser := &dao.UserEntity{
			Username: vo.Username,
			Password: vo.Password,
			Email:    vo.Email,
		}
		err = userService.useDAO.Save(newUser)
		if err == nil {
			return &UserInfoDto{
				ID:       newUser.ID,
				Email:    newUser.Email,
				Username: newUser.Username,
			}, nil
		}
	}
	if err != nil {
		err = ErrUserExisted
	}
	return nil, err
}
func MakeUserServiceImpl(userDAO dao.UserDao) UserService {
	return &UserServiceImpl{
		useDAO: userDAO,
	}
}

func (userService *UserServiceImpl) Login(ctx context.Context, email, password string) (*UserInfoDto, error) {
	user, err := userService.useDAO.SelectByEmail(email)
	if err == nil {
		if user.Password == password {
			return &UserInfoDto{
				ID:       user.ID,
				Username: user.Username,
				Email:    user.Email,
			}, nil
		}
	} else {
		log.Printf("err : %s", ErrPassword)
	}
	return nil, err
}
