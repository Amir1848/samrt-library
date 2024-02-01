package library

import "time"

type GnrLibrary struct {
	Id       int64
	Title    string
	UserId   int64
	Token    string
	IsOnline bool
}

type GnrLibraryItem struct {
	Title     string
	Status    LibraryStatus
	LibraryId int64
}

type LibraryStatus int

type LibraryForView struct {
	Id       int64
	Title    string
	IsOnline bool
}

type LibraryInfoWithItems struct {
	Library      *LibraryForView
	LibraryItems []*GnrLibraryItem
}

type GnrStudentHistory struct {
	Id             int64
	UserRef        int64
	Date           time.Time `gorm:"column:date_c"`
	LibraryItemRef int64
}
