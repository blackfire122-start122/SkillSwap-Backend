package pkg

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

type User struct {
	gorm.Model
	Id           uint64 `gorm:"primaryKey"`
	Username     string
	Password     string
	Email        string
	Phone        string
	Image        string
	Rating       uint8
	Skills       []Skill      `gorm:"many2many:user_skills;"`
	Categories   []Category   `gorm:"many2many:categories_user;"`
	PricesSkills []PriceSkill `gorm:"foreignKey:UserID;"`
	Reviews      []Review     `gorm:"foreignKey:UserID;"`
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

type Review struct {
	gorm.Model
	Id         uint64 `gorm:"primaryKey"`
	UserID     uint64
	ReviewerID uint64
	Reviewer   User `gorm:"foreignKey:ReviewerID;"`
	Review     string
	Rating     uint8
}

type PriceSkill struct {
	gorm.Model
	UserID  uint64
	SkillId uint64
	Skill   Skill `gorm:"foreignKey:SkillId;"`
	Price   uint64
}

func init() {
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{}) //ToDo change on postgres
	DB = db

	if err != nil {
		panic("failed to connect database")
	}

	err = DB.AutoMigrate(&User{}, &Category{}, Skill{}, PriceSkill{}, Review{})
	if err != nil {
		panic("Error autoMigrate: ")
	}
}
