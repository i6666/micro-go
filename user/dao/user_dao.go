package dao

import "time"

type UserEntity struct {
	ID        int64
	Username  string
	Password  string
	Email     string
	CreatedAt time.Time
}

// 方法就是一类带特殊 接收者参数的函数，方法接收者位于func 关键字和方法名之间
func (UserEntity) TableName() string {
	return "user"
}

type UserDao interface {
	SelectByEmail(email string) (*UserEntity, error)
	Save(user *UserEntity) error
}

type UserDaoImpl struct {
}

func (userDAO *UserDaoImpl) SelectByEmail(email string) (*UserEntity, error) {
	user := &UserEntity{}
	err := db.Where("email = ?", email).First(user).Error
	return user, err
}

func (userDAO *UserDaoImpl) Save(user *UserEntity) error {
	return db.Create(user).Error
}
