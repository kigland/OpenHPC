package dboper

import (
	"github.com/kigland/OpenHPC/coordinator/models/dbmod"
	"github.com/kigland/OpenHPC/coordinator/shared"
)

func GetUserByID(id string) (*dbmod.User, error) {
	var user dbmod.User
	if err := shared.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
