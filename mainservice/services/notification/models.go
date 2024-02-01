package notification

import "time"

type GnrNotification struct {
	Id      int64
	UserRef int64
	Type    int       `gorm:"column:type_c"`
	Date    time.Time `gorm:"column:date_c"`
}
