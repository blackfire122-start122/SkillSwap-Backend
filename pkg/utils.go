package pkg

import (
	"crypto/rand"
	"errors"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"math/big"
	"net/http"
	"os"
	"strconv"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("SECRET_KEY")))

type UserLogin struct {
	Username string
	Password string
}

func Login(w http.ResponseWriter, r *http.Request, userLogin *UserLogin) error {
	session, _ := store.Get(r, "session-name")

	var user User
	if err := DB.First(&user, "Username = ?", userLogin.Username).Error; err != nil {
		return err
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userLogin.Password))
	if err == nil {
		session.Values["id"] = user.Id
		session.Values["password"] = user.Password
		err = session.Save(r, w)
		if err != nil {
			return err
		}
	} else {
		return err
	}

	return nil
}

func Sign(user *UserLogin) error {
	var users []User

	if err := DB.Where("Username = ?", user.Username).Find(&users).Error; err != nil {
		return err
	}
	if len(users) > 0 {
		return errors.New("user with the same username already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	DB.Create(&User{Username: user.Username, Password: string(hashedPassword)})
	return err
}

func CheckSessionUser(r *http.Request) (bool, User) {
	session, _ := store.Get(r, "session-name")

	var user User

	if session.IsNew {
		return false, user
	}

	if err := DB.First(&user, "Id = ?", session.Values["id"]).Error; err != nil {
		return false, user
	}

	if session.Values["password"] != user.Password {
		return false, user
	}
	return true, user
}

//func CheckAdmin(user User) bool {
//	var admin Admin
//	if err := DB.Where("user_id=?", user.Id).Find(&admin).Error; err != nil {
//		return false
//	}
//
//	return admin.UserId == user.Id
//}

func Logout(w http.ResponseWriter, r *http.Request) bool {
	session, _ := store.Get(r, "session-name")

	session.Values["id"] = nil
	session.Values["password"] = nil

	session.Options.MaxAge = -1

	err := session.Save(r, w)
	if err != nil {
		return false
	}
	return true
}

type ChangedUser struct {
	Username string
	Email    string
	Phone    string
	Skills   []string
}

func ChangeUser(changedUser *ChangedUser, user User) error {
	var users []User

	if err := DB.Where("Username = ?", changedUser.Username).Find(&users).Error; err != nil {
		return err
	}

	if len(users) > 0 {
		if users[0].Id != user.Id {
			return errors.New("user with the same username already exists")
		}
	}

	var skillIds []uint64

	for _, skill := range changedUser.Skills {
		skillId, err := strconv.ParseUint(skill, 10, 64)
		if err != nil {
			return err
		}
		skillIds = append(skillIds, skillId)
	}

	var skills []Skill

	if err := DB.Where("Id in ?", skillIds).Find(&skills).Error; err != nil {
		return err
	}

	return DB.Transaction(func(tx *gorm.DB) error {
		user.Username = changedUser.Username
		user.Email = changedUser.Email // ToDo need check on unique
		user.Phone = changedUser.Phone // ToDo need check on unique

		if err := tx.Model(&user).Association("Skills").Clear(); err != nil {
			return err
		}
		if err := tx.Model(&user).Association("Skills").Append(skills); err != nil {
			return err
		}

		if err := tx.Save(&user).Error; err != nil {
			return err
		}

		return nil
	})
}

func GenerateSalt() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	randomString := make([]byte, 10)
	for i := range randomString {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			panic(err)
		}
		randomString[i] = charset[num.Int64()]
	}
	return string(randomString)
}
