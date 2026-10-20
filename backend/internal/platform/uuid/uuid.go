package uuid

import googleuuid "github.com/google/uuid"

type Uuid interface {
	GenerateAccessToken() string
	GenerateRefreshToken() string
}

type uuidImpl struct{}

func New() Uuid {
	return &uuidImpl{}
}

func (u *uuidImpl) GenerateAccessToken() string {
	return googleuuid.NewString()
}

func (u *uuidImpl) GenerateRefreshToken() string {
	return googleuuid.NewString()
}
