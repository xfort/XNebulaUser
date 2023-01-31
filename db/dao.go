package db

type XNebulaUserDBer interface {
	Add(*XUserDAO) (*XUserDAO, error)
	//手机+短信验证
	QueryByMobile(mobile string, authCode int32) (*XUserDAO, error)
	//邮箱+密码
	QueryByEmailPasswd(email string, passwd string) (*XUserDAO, error)
	//手机号+密码
	QueryByMobilePasswd(mobile string, passwd string) (*XUserDAO, error)
	//退出登录
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
	Mobile   string
	Email    string
	Icon     string
	Token    string
}

//userID,nickname,passwd,mobile,email,icon,signupUTC

type UserTable string

const (
	UserTab_NAME    UserTable = "xuser"
	UserTab_USER_ID UserTable = "user_id"
	UserTab_PASSWD  UserTable = "passwd"
	//格式是 地区码-手机号，必须包含地区码，例如 86-16602110123
	UserTab_MOBILE UserTable = "mobile"
	UserTab_EMAIL  UserTable = "email"
	UserTab_TOKEN  UserTable = "token"

	UserTab_NICKNAME   UserTable = "nickname"
	UserTab_ICON       UserTable = "icon"
	UserTab_StatueCode UserTable = "statue_code"
	UserTab_CREATE_UTC UserTable = "create_utc"
)
const (
	XUserTokenTab_Name    = "XUserToken"
	XUserTokenTab_USER_ID = "user_id"
	XUserTokenTab_TOKEN   = "token"
)
