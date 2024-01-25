package library

import (
	"context"

	authutil "github.com/Amir1848/samrt-library/routes/authUtil"
	"github.com/Amir1848/samrt-library/utils/dbutil"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func Insert(ctx context.Context, model *GnrLibrary, items []*GnrLibraryItem) (int64, error) {
	db, err := dbutil.GetDBConnection(ctx)
	if err != nil {
		return 0, err
	}

	userId := authutil.GetUserId(ctx)
	model.UserId = userId
	model.Token = uuid.New().String()

	err = dbutil.CreateDatabaseTransaction(db, func(tx *gorm.DB) error {
		err := tx.Table("gnr_library").Create(model).Error
		if err != nil {
			return err
		}

		for _, item := range items {
			item.LibraryId = model.Id
		}

		err = tx.Table("gnr_library_item").Create(&items).Error
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

func GetLibAdminLibrary(ctx context.Context) (result *GnrLibrary, err error) {
	db, err := dbutil.GetDBConnection(ctx)
	if err != nil {
		return nil, err
	}

	userId := authutil.GetUserId(ctx)

	lib := GnrLibrary{}
	err = db.Table("gnr_library").
		Where("user_id = ?", userId).
		Limit(1).
		Select("*").Scan(&lib).Error
	if err != nil {
		return nil, err
	}

	return &lib, nil
}

func GetLibraryWithToken(ctx context.Context, tx *gorm.DB, token string) (lib *GnrLibrary, found bool, err error) {
	lib = &GnrLibrary{}
	fetchResult := tx.Table("gnr_library").Where("token = ?", token).Scan(lib)
	if fetchResult.Error != nil {
		return nil, false, err
	}

	return lib, fetchResult.RowsAffected > 0, nil
}

func SetLibraryStatus(ctx context.Context, tx *gorm.DB, libID int64, isOnline bool) error {
	return tx.Table("gnr_library").Where("id = ?", libID).Update("is_online", isOnline).Error
}

func SetLibraryItemsAsUnknown(ctx context.Context, tx *gorm.DB, tableId int64) error {
	err := tx.Table("gnr_library_item").
		Where("library_id = ?", tableId).
		Update("status", 0).Error
	if err != nil {
		return err
	}

	return nil
}

func SetLibItemStatus(ctx context.Context, tx *gorm.DB, libID int64, itemName string, status int) error {
	err := tx.Table("gnr_library_item").
		Where("library_id = ?", libID).
		Update("status", status).Error
	if err != nil {
		return err
	}

	return nil
}

func GetAllLibraries(ctx context.Context) ([]*LibraryForView, error) {
	db, err := dbutil.GetDBConnection(ctx)
	if err != nil {
		return nil, err
	}

	libs := []*LibraryForView{}
	err = db.Table("gnr_library").Select("*").Scan(&libs).Error
	if err != nil {
		return nil, err
	}

	return libs, nil
}

func GetLibraryWithItems(ctx context.Context, id int64) (result *LibraryInfoWithItems, found bool, err error) {
	db, err := dbutil.GetDBConnection(ctx)
	if err != nil {
		return
	}

	lib := LibraryForView{}
	fetchResult := db.Table("gnr_library").Where("id = ?", id).Select("*").Scan(&lib)
	if fetchResult.Error != nil {
		err = fetchResult.Error
		return
	}

	if fetchResult.RowsAffected == 0 {
		found = false
		return
	}
	found = true

	libItems := []*GnrLibraryItem{}
	err = db.Table("gnr_library_item").Where("library_id = ?", id).Scan(&libItems).Error
	if err != nil {
		return
	}

	result = &LibraryInfoWithItems{
		Library:      &lib,
		LibraryItems: libItems,
	}
	return result, true, nil
}
