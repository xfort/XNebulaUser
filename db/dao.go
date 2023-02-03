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

	UserTab_NICKNAME UserTable = "nickname"
	UserTab_ICON     UserTable = "icon"
	//用户状态码，是否被禁用、注销
	UserTab_StatueCode UserTable = "statue_code"
	UserTab_CREATE_UTC UserTable = "create_utc"
)

const (
	//用户token表
	XUserTokenTab_Name    = "XUserToken"
	XUserTokenTab_USER_ID = "user_id"
	XUserTokenTab_TOKEN   = "token"
)

const (
	//用户邮箱验证表
	XUserEmailAuthTab_Name    = "XUserEmailAuth"
	XUserEmailAuthTab_USER_ID = "user_id"
	XUserEmailAuthTab_EMAIL   = "email"
	XUserEmailAuthTab_CODE    = "email_auth_code"

	//业务类别，例如注册、登录、找回密码、改绑等
	XUserEmailAuthTab_ACTION = "action_type"
	//单位时间内创建的次数，用于限制 10分钟内发送验证码次数
	XUserEmailAuthTab_CREATE_COUNT = "create_count"
	//复活时间点。限制发送次数
	XUserEmailAuthTab_DEAD_UTC = "dead_utc"
	//最后一次创建时间 UTC
	XUserEmailAuthTab_LAST_CREATE = "last_create_utc"
)

const (
	//用户手机号验证表
	XUserMobileAuthTab_Name    = "XUserMobileAuth"
	XUserMobileAuthTab_USER_ID = "user_id"
	XUserMobileAuthTab_MOBILE  = "mobile"
	XUserMobileAuthTab_CODE    = "auth_code"
	//业务类别，例如注册、登录、找回密码、改绑其它手机号
	XUserMobileAuthTab_ACTION = "action_type"
	//单位时间内创建的次数，用于限制 60分钟内发送验证码次数
	XUserMobileAuthTab_CREATE_COUNT = "create_count"
	//复活时间点。
	XUserMobileAuthTab_DEAD_UTC = "dead_utc"
)

const (
	//用户登录异常、失败表,存储登录失败的时间、次数，用于在连续多次登录失败时暂时锁定账户
	XUserLoginErrTab_Name    = "XUserLoginErr"
	XUserLoginErrTab_USER_ID = "user_id"

	//单位时间内创建的次数，用于60分钟内限制连续登录失败次数
	XUserLoginErrTab_CREATE_COUNT = "err_count"
	//复活时间点。例如在 1：00时登录失败，则复活时间为60分钟后2：00，在这60分钟内任何登录失败都计入 err_count
	XUserLoginErrTab_DEAD_UTC = "dead_utc"
	XUserLoginErrTab_LAST_ERR = "last_error_utc"
)
