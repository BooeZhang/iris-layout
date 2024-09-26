package mysql

import (
	"github.com/kataras/iris/v12"
	"gorm.io/gorm"

	"irir-layout/internal/model"
	"irir-layout/internal/repo"
)

var _ repo.UserRepo = (*UserRepo)(nil)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (ur *UserRepo) GetUserByName(_ iris.Context, account string) (*model.User, error) {
	user := &model.User{}
	err := ur.db.Where("account = ?", account).Find(user).Error
	return user, err
}

func (ur *UserRepo) GetUserById(_ iris.Context, uid int64) (*model.User, error) {
	user := &model.User{}
	err := ur.db.Where("id = ?", uid).Find(user).Error
	return user, err
}

func (ur *UserRepo) GetUserByMobile(_ iris.Context, mobile string) (*model.User, error) {
	user := &model.User{}
	err := ur.db.
		Where("mobile = ?", mobile).
		Where("enabled_status = 1").
		First(user).Error
	return user, err
}
