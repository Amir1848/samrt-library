package users

import (
	"context"

	"github.com/Amir1848/samrt-library/utils/dbutil"
)

func GetUsers(ctx context.Context) error {
	_, err := dbutil.GetDBConnection(ctx)
	if err != nil {
		return err
	}

	return nil
	// db.Table("users")
}
