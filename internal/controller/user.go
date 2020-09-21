package controller

import (
	"github.com/go-xorm/xorm"
	"github.com/thunes/internal/db"
	"github.com/thunes/internal/model"
)

func NewUserHandler() *User {
	return &User{
		M: &UserModel{
			session: db.Cli().Table(&model.User{}),
		},
	}
}

type User struct {
	M *UserModel
}

///////////////////////////////////////////////////////////////////////////////

// UserModel 提供数据库数据操作功能接口
type UserModel struct {
	session *xorm.Session
}

func (M *UserModel) Get(uid model.UID) (user model.User, err error) {
	buf := []model.User{}

	if err = M.session.Where("id = ?", uid).Find(&buf); err != nil {
		return
	}

	if len(buf) == 0 {
		err = ErrAccountNotFound
		return
	}

	return buf[0], nil
}

func (M *UserModel) FindByEmail(email string) (user model.User, err error) {
	buf := []model.User{}

	if err = M.session.Where("email = ?", email).Find(&buf); err != nil {
		return
	}

	if len(buf) == 0 {
		err = ErrAccountNotFound
		return
	}

	return buf[0], nil
}

func (M UserModel) EncryptPWD(user model.User) (outer string) {
	//TODO: 依据user中信息生成一个非可逆的sha256值
	return user.PWD
}

func (M UserModel) VerifyPWD(user model.User, pwd string) (ok bool) {
	buf := user
	buf.PWD = pwd

	if M.EncryptPWD(buf) != user.PWD {
		return false
	}

	return true
}
