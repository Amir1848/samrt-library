package library

import (
	"context"

	"github.com/Amir1848/samrt-library/utils/dbutil"
	"gorm.io/gorm"
)

func Insert(ctx context.Context, model *GnrLibrary) (int64, error) {
	db, err := dbutil.GetDBConnection(ctx)
	if err != nil {
		return 0, err
	}

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
