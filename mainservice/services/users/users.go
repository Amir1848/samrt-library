package users

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/Amir1848/samrt-library/utils/dbutil"
	jwt "github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func GetUsers(ctx context.Context) error {
	_, err := dbutil.GetDBConnection(ctx)
	if err != nil {
		return err
	}

	return nil
	// db.Table("users")
}

func RegisterUser(ctx context.Context, model *UserViewModel) error {
	dbUser := &GnrUser{
		StudentCode: model.StudentCode,
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(model.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	dbUser.Password = string(hash)

	db, err := dbutil.GetDBConnection(ctx)
	if err != nil {
		return err
	}

	err = dbutil.CreateDatabaseTransaction(db, func(tx *gorm.DB) error {
		err = tx.Table("gnr_user").Create(dbUser).Error
		if err != nil {
			return err
		}

		userRole := GnrUserRole{
			UserId: dbUser.Id,
			Role:   RoleUser,
		}

		err = tx.Table("gnr_user_role").Create(&userRole).Error
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func Login(ctx context.Context, loginReq *UserViewModel) (jwtToken string, result bool, err error) {
	//todo: sleep n senconds (n is random)
	db, err := dbutil.GetDBConnection(ctx)
	if err != nil {
		return "", false, err
	}

	userRoles := []UserRole{}
	var userId int64

	err = dbutil.CreateDatabaseTransaction(db, func(tx *gorm.DB) error {
		user := &GnrUser{}
		fetchResult := tx.Table("gnr_user").Where("student_code = ?", loginReq.StudentCode).Scan(user)
		if fetchResult.Error != nil {
			return fetchResult.Error
		}

		userId = user.Id

		if fetchResult.RowsAffected == 0 {
			return errors.New("username or password is wrong")
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password))
		if err != nil {
			return errors.New("username or password is wrong")
		}

		err = tx.Table("gnr_user_role").Where("user_id = ?", user.Id).Pluck("role_c", &userRoles).Error
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return "", false, err
	}

	userToken, err := createJwtToken(loginReq.StudentCode, userRoles, userId)
	if err != nil {
		return "", true, err
	}

	return userToken, true, nil
}

func createJwtToken(studentCode string, roles []UserRole, userId int64) (string, error) {
	key := []byte(os.Getenv("JwtKey"))
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"studentCode": studentCode,
		"roles":       roles,
		"userId":      userId,
		"exp":         time.Now().Add(3 * time.Hour).Unix(),
	})
	s, err := t.SignedString(key)

	return s, err
}
