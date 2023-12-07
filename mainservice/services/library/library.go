package library

import (
	"context"

	authutil "github.com/Amir1848/samrt-library/routes/authUtil"
	"github.com/Amir1848/samrt-library/utils/dbutil"
	"gorm.io/gorm"
)

func Insert(ctx context.Context, model *GnrLibrary) (int64, error) {
	db, err := dbutil.GetDBConnection(ctx)
	if err != nil {
		return 0, err
	}

	userId, _ := authutil.GetUserId(ctx)
	model.UserId = userId

	err = dbutil.CreateDatabaseTransaction(db, func(tx *gorm.DB) error {
		err := tx.Table("gnr_library").Create(model).Error
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return model.Id, err
}

func GetLibraries(ctx context.Context) (result []*GnrLibrary, err error) {
	db, err := dbutil.GetDBConnection(ctx)
	if err != nil {
		return nil, err
	}

	libs := []*GnrLibrary{}
	err = db.Table("gnr_library").Select("*").Scan(&libs).Error
	if err != nil {
		return nil, err
	}

	return libs, nil
}
