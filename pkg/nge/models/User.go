package models

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"fmt"
)

type User struct {
    gorm.Model
    Username       string
    Email          string
    Password       string
    ShownPassword  string
    RoleID         int
	Role 		  Role 	`gorm:"foreignKey:RoleID"`
}


type Claims struct {
	UserID uint
	jwt.StandardClaims
}

var jwtKey = []byte("env.jwtKey")

func GenerateToken(userID uint)(string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	return tokenString, err
}

func ValidateToken(tokenString string) (*Claims, error) {
    claims := &Claims{}

    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })

    if err != nil {
        return nil, err
    }

    if !token.Valid {
        return nil, fmt.Errorf("invalid token")
    }

    return claims, nil
}
func HashPassword(password string)(string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func GetUserByEmail(db *gorm.DB, email string) (*User, error) {
	var user User
	result := db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func Signup(db *gorm.DB, username, email, password string, roleID int) (*User, error) {
    shownPassword := password
    hashedPassword, err := HashPassword(password)
    if err != nil {
        return nil, err
    }

    user := User{
        Username: username,
        Email: email,
        Password: hashedPassword,
        ShownPassword: shownPassword,
        RoleID: roleID,
    }

    result := db.Create(&user)
    if result.Error != nil {
        return nil, result.Error
    }

    err = db.Preload("Role").First(&user, user.ID).Error
    if err != nil {
        return nil, err
    }

    return &user, nil
}


func UpdateUser(db *gorm.DB, user User) (User, error) {
	result := db.Save(user).Error
	return user, result
}

func DeleteUser(db *gorm.DB, id uint) error {
	result := db.Delete(&User{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetUserRole(db *gorm.DB, userID uint) (uint, error) {
    var roleID uint
	fmt.Println("dada")
	fmt.Println(userID)
	fmt.Println("dada")

    query := "SELECT role_id FROM users WHERE id = ?"
    if err := db.Raw(query, userID).Row().Scan(&roleID); err != nil {
        if err == gorm.ErrRecordNotFound {
            return 0, errors.New("user not found")
        }
        return 0, err
    }

	fmt.Println("role", roleID)

    return roleID, nil
}
