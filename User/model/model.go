package model

import (
	"github.com/dgrijalva/jwt-go"
	"gorm.io/datatypes"
)

type Claims struct {
	Data string `json:"username"`
	jwt.StandardClaims
}
type User struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"unique"`
	Password string
}

type Metadata struct {
	ID                     uint `gorm:"primaryKey"`
	Name                   string
	Filename               string
	Chunk_location_details datatypes.JSON
}
