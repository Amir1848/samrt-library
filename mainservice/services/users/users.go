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

func RegisterUser(ctx context.Context, model *UserViewModel, roles []UserRole) error {
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

		userRoles := []*GnrUserRole{}
		for _, item := range roles {
			userRoles = append(userRoles, &GnrUserRole{
				UserId: dbUser.Id,
				Role:   item,
			})
		}

		err = tx.Table("gnr_user_role").Create(&userRoles).Error
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

func GetUserIdByStudentCode(ctx context.Context, tx *gorm.DB, studentCode string) (int64, error) {

	var result int64 = 0
	err := tx.Table("gnr_user").
		Where("student_code = ?", studentCode).
		Select("id").Scan(&result).Error
	if err != nil {
		return 0, err
	}

	return result, nil
}

func GetAllUserStudentCodes(ctx context.Context) ([]string, error) {
	var allStudents = []string{}

	db, err := dbutil.GetDBConnection(ctx)
	if err != nil {
		return nil, err
	}

	err = db.Table("gnr_user").Select("student_code").Scan(&allStudents).Error
	if err != nil {
		return nil, err
	}

	return allStudents, nil
}
