package dboper

import (
	"github.com/kigland/OpenHPC/coodinator/models/dbmod"
	"github.com/kigland/OpenHPC/coodinator/shared"
)

func GetUserByID(id string) (*dbmod.User, error) {
	var user dbmod.User
	if err := shared.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
