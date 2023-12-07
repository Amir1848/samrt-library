package users

type UserViewModel struct {
	StudentCode string
	Password    string
}

type GnrUser struct {
	Id          int64
	StudentCode string
	Password    string `gorm:"column:password_c"`
}

type GnrUserRole struct {
	Id     int64
	UserId int64
	Role   UserRole `gorm:"column:role_c"`
}

type UserRole float64

const (
	RoleSysAdmin UserRole = 1
	RoleLibAdmin UserRole = 2
	RoleUser     UserRole = 3
)
