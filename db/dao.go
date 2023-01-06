package db

type XNebulaUserDBer interface {
	Add(*XUserDAO) error
	QueryByMobile(mobile string) (*XUserDAO, error)
	QueryByEmailPasswd(email string, passwd string) (*XUserDAO, error)
	QueryByMobilePasswd(mobile string, passwd string) (*XUserDAO, error)
	Logout(userID int64, token string) error
	Destroy()
}

type XNebulaAdminDBer interface {
	SignUp()
	Login()
	Logout()
	Destroy()
}

type XUserDAO struct {
	UserID   int64
	Nickname string
	Passwd   string
	Mobile   *string
	Email    *string
	Icon     *string
	Token    *string
}
