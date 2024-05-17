package pkg

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

type User struct {
	gorm.Model
	Id         uint64 `gorm:"primaryKey"`
	Username   string
	Password   string
	Email      string
	Phone      string
	Image      string
	Rating     uint8
	Skills     []Skill      `gorm:"many2many:user_skills;"`
	Categories []Category   `gorm:"many2many:categories_user;"`
	PriceSkill []PriceSkill `gorm:"many2many:price_skills;"`
}

type Category struct {
	gorm.Model
	Id     uint64 `gorm:"primaryKey"`
	Name   string
	Skills []Skill `gorm:"many2many:categories_skills;"`
}

type Skill struct {
	gorm.Model
	Id   uint64 `gorm:"primaryKey"`
	Name string
}

type PriceSkill struct {
	gorm.Model
	Id    uint64 `gorm:"primaryKey"`
	price uint64
}

func init() {
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	DB = db

	if err != nil {
		panic("failed to connect database")
	}

	err = DB.AutoMigrate(&User{}, &Category{}, Skill{})
	if err != nil {
		panic("Error autoMigrate: ")
	}

}
