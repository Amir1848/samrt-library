package notification

import (
	"context"
	"time"

	authutil "github.com/Amir1848/samrt-library/routes/authUtil"
	"github.com/Amir1848/samrt-library/services/users"
	"github.com/Amir1848/samrt-library/utils/dbutil"
	"gorm.io/gorm"
)

func NotifyUser(ctx context.Context, studentCode string, messageType int) error {
	db, err := dbutil.GetDBConnection(ctx)
	if err != nil {
		return err
	}

	err = dbutil.CreateDatabaseTransaction(db, func(tx *gorm.DB) error {
		userId, err := users.GetUserIdByStudentCode(ctx, tx, studentCode)
		if err != nil {
			return err
		}

		err = tx.Table("gnr_notification").Create(
			&GnrNotification{
				UserRef: userId,
				Type:    messageType,
				Date:    time.Now(),
			},
		).Error
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

func GetUsersNotifications(ctx context.Context) ([]*GnrNotification, error) {
	db, err := dbutil.GetDBConnection(ctx)
	if err != nil {
		return nil, err
	}

	userId := authutil.GetUserId(ctx)

	var result = []*GnrNotification{}
	err = db.Table("gnr_notification").Where("user_ref = ?", userId).Scan(&result).Error
	if err != nil {
		return nil, err
	}

	var messageTypeMapper = map[int]string{
		1: "داش بیا سر جات",
	}

	for _, item := range result {
		item.Title = messageTypeMapper[item.Type]
	}

	return result, nil
}
