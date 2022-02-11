package model

type UserDetails struct {
	//用户标识
	UserId   int64
	Username string
	Password string
	//用户具有的权限
	Authorities []string
}
