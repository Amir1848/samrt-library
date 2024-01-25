package library

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
