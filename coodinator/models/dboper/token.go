package dboper

import (
	"github.com/kigland/HPC-Scheduler/coodinator/models/dbmod"
	"github.com/kigland/HPC-Scheduler/coodinator/shared"
)

func GetTokenByToken(token string) (dbmod.Token, error) {
	var tk dbmod.Token
	err := shared.DB.Where("token = ?", token).First(&tk).Error
	return tk, err
}

func CreateToken(token string, userID string) error {
	return shared.DB.Create(&dbmod.Token{Token: token, UserId: userID}).Error
}
