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
	CityID       uint64
	City         City         `gorm:"foreignKey:CityID;"`
	Skills       []Skill      `gorm:"many2many:user_skills;"`
	Categories   []Category   `gorm:"many2many:categories_user;"`
	PricesSkills []PriceSkill `gorm:"foreignKey:UserID;"`
	Reviews      []Review     `gorm:"foreignKey:UserID;"`
	//Messages     		 []UserMessage `gorm:"foreignKey:UserID;"`
	CustomerSkillChats  []SkillChat `gorm:"foreignKey:CustomerID;"`
	PerformerSkillChats []SkillChat `gorm:"foreignKey:PerformerID;"`
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

//type UserMessage struct {
//	gorm.Model
//	UserID  uint64
//	Message string
//	Read    bool
//}

type Message struct {
	gorm.Model
	Message     string
	Read        bool
	UserID      uint64
	User        User `gorm:"foreignKey:UserID;"`
	SkillChatID uint64
}

type SkillChat struct {
	gorm.Model
	Messages    []Message `gorm:"foreignKey:SkillChatID;"`
	SkillID     uint64
	Skill       Skill `gorm:"foreignKey:SkillID;"`
	CustomerID  uint64
	PerformerID uint64
}

type City struct {
	gorm.Model
	Name string
}

func init() {
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{}) //ToDo change on postgres
	DB = db

	if err != nil {
		panic("failed to connect database")
	}

	err = DB.AutoMigrate(
		User{},
		Category{},
		Skill{},
		PriceSkill{},
		Review{},
		SkillChat{},
		Message{},
		City{},
	)

	if err != nil {
		panic("Error autoMigrate: ")
	}
}
