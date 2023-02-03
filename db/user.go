package db

import (
	"crypto/md5"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	//_ "github.com/xeodou/go-sqlcipher"
	"log"
	"path/filepath"
	"strings"
	"time"
)

//https://support.huaweicloud.com/intl/zh-cn/productdesc-msgsms/phone_numbers.html

const XUSER_SQLITE_DB = "XNebulaUserS3.db"

// 创建用户信息表
var createXUserTab = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
    %s INTEGER PRIMARY KEY,
    %s varchar(128) NOT NULL,
    %s varchar(64) NOT NULL,
    %s varchar(20) NOT NULL ,
    %s varchar(64) NOT NULL ,
    %s varchar(200) NOT NULL,
    %s INTEGER NOT NULL DEFAULT 1, 
    %s INTEGER NOT NULL
    );`,
	UserTab_NAME, UserTab_USER_ID, UserTab_PASSWD, UserTab_NICKNAME, UserTab_MOBILE, UserTab_EMAIL, UserTab_ICON, UserTab_StatueCode, UserTab_CREATE_UTC)

// 创建用户表索引
var sql_CreateXUserTabIndex = [2]string{
	//WHERE %s IS NOT NULL
	fmt.Sprintf(`CREATE UNIQUE INDEX IF NOT EXISTS xuser_mobile ON %s(%s) ;`, UserTab_NAME, UserTab_MOBILE),
	fmt.Sprintf(`CREATE UNIQUE INDEX IF NOT EXISTS xuser_email ON %s(%s);`, UserTab_NAME, UserTab_EMAIL),
}

// 用户token表
var createXUserTokenTab = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
    %s INTEGER PRIMARY KEY,
    %s varchar(128) NOT NULL);`,
	XUserTokenTab_Name, XUserTokenTab_USER_ID, XUserTokenTab_TOKEN)

// 邮箱验证表，user_id、email都不为空
var createXUserEmailAuthTab = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
	%s INTEGER PRIMARY KEY,
    %s varchar(64) UNIQUE INDEX,
    %s INT NOT NULL DEFAULT 1,
    %s INT NOT NULL DEFAULT 1,
    %s INT NOT NULL DEFAULT 0,
    %s INTEGER NOT NULL DEFAULT 0,
    %s INTEGER NOT NULL DEFAULT 0
    );`,
	XUserEmailAuthTab_Name, XUserEmailAuthTab_USER_ID, XUserEmailAuthTab_EMAIL, XUserEmailAuthTab_CODE, XUserEmailAuthTab_ACTION,
	XUserEmailAuthTab_CREATE_COUNT, XUserEmailAuthTab_DEAD_UTC, XUserEmailAuthTab_LAST_CREATE)

var xUserInfoColumens = [6]UserTable{UserTab_USER_ID, UserTab_NICKNAME, UserTab_MOBILE, UserTab_EMAIL, UserTab_ICON, UserTab_StatueCode}

// 用于查询用户表的前段固定SQL
var sql_QueryXUserInfo = fmt.Sprintf(`SELECT %s,%s,%s,%s,%s FROM %s`, xUserInfoColumens[0], xUserInfoColumens[1], xUserInfoColumens[2], xUserInfoColumens[3], xUserInfoColumens[4], UserTab_NAME)

// xuser 表的版本，用于以后更新表结构
const XUserTab_Ver = 2023011101

type XUserSQLiteDB struct {
	userDB *sql.DB
}

func OpenXUserDB(dbDir string) (*XUserSQLiteDB, error) {
	userDBPath := filepath.Join(dbDir, XUSER_SQLITE_DB)
	adminPass := "XNebulaUser"
	authTEMP := md5.Sum([]byte(XUSER_SQLITE_DB))
	for index := 0; index < len(authTEMP); index++ {
		if index%2 == 0 {
			adminPass = adminPass + fmt.Sprintf("%x", authTEMP[index])
		}
	}

	dbURI := fmt.Sprintf("file:%s?_auth&_auth_user=XNebulaUserAdmin&_auth_pass=%s&_key=%s", userDBPath, adminPass, adminPass)

	xUserDB, err := sql.Open("sqlite3", dbURI)
	if err != nil {
		return nil, err
	}
	xUserDB.SetMaxOpenConns(1)
	err = createTable(xUserDB)
	if err != nil {
		xUserDB.Close()
		return nil, fmt.Errorf("创建用户表失败,%s,%s", createXUserTab, err.Error())
	}

	return &XUserSQLiteDB{userDB: xUserDB}, nil
}

func createTable(userDB *sql.DB) error {
	_, err := userDB.Exec(createXUserTab)
	if err != nil {
		return fmt.Errorf("创建用户表错误：%s", err.Error())
	}
	for _, itemSQL := range sql_CreateXUserTabIndex {
		_, err := userDB.Exec(itemSQL)
		if err != nil {
			return fmt.Errorf("创建用户表索引错误%s,%s", itemSQL, err.Error())
		}
	}
	_, err = userDB.Exec(createXUserTokenTab)
	return err
}

// 新增用户。先检查 手机号、邮箱是否已使用
func (xuserDB *XUserSQLiteDB) Add(user *XUserDAO) (*XUserDAO, error) {
	if user.UserID <= 0 || len(user.Passwd) < 16 {
		return user, errors.New("userID，passwd 参数异常")
	}
	//if (len(user.Mobile) < 8) || (len(user.Email) < 4) {
	//	return user, errors.New("Mobile、Email参数异常_")
	//}
	sql := fmt.Sprintf(`INSERT INTO %s(%s,%s,%s,%s,%s,%s,%s) VALUES(?,?,?,?,?,?,?);`,
		UserTab_NAME, UserTab_USER_ID, UserTab_PASSWD, UserTab_NICKNAME, UserTab_MOBILE, UserTab_EMAIL, UserTab_ICON, UserTab_CREATE_UTC)

	_, err := xuserDB.userDB.Exec(sql, user.UserID, user.Nickname, user.Passwd, user.Mobile, user.Email, user.Icon, time.Now().Unix())
	return user, err
}

// 手机+短信验证
func (xuserDB *XUserSQLiteDB) QueryByMobile(mobile string, authCode int32) (*XUserDAO, error) {
	if !strings.Contains(mobile, "-") || len(mobile) < 8 {
		return nil, errors.New("手机号格式异常，格式必须是地区码-手机号")
	}
	sql := sql_QueryXUserInfo + fmt.Sprintf(` WHERE %s=? LIMIT 1;`, UserTab_MOBILE)
	resRow, err := xuserDB.userDB.Query(sql, mobile)
	if err != nil {
		return nil, err
	}
	defer resRow.Close()

	xuserDAO := XUserDAO{}
	err = resRow.Scan(&xuserDAO.UserID, &xuserDAO.Nickname, xuserDAO.Mobile, xuserDAO.Email, xuserDAO.Icon)
	if err != nil {
		return nil, err
	}
	return &xuserDAO, nil
}

func (xuserDB *XUserSQLiteDB) QueryByEmailPasswd(email string, passwd string) (*XUserDAO, error) {
	if len(email) <= 4 || len(passwd) != 128 {
		return nil, errors.New("email、密码参数异常_" + email + "_" + passwd)
	}
	sql := sql_QueryXUserInfo + fmt.Sprintf(` WHERE %s=? AND %s=? LIMIT 1;`, UserTab_EMAIL, UserTab_PASSWD)
	resRow, err := xuserDB.userDB.Query(sql, email, passwd)
	if err != nil {
		return nil, err
	}
	defer resRow.Close()
	xuserDAO := XUserDAO{}
	err = resRow.Scan(&xuserDAO.UserID, &xuserDAO.Nickname, xuserDAO.Mobile, xuserDAO.Email, xuserDAO.Icon)
	if err != nil {
		return nil, err
	}
	return &xuserDAO, nil
}

func (xuserDB *XUserSQLiteDB) QueryByMobilePasswd(mobile string, passwd string) (*XUserDAO, error) {
	if !strings.Contains(mobile, "-") || len(mobile) < 8 || len(passwd) != 128 {
		return nil, errors.New("手机号异常(格式必须是地区码-手机号)/密码异常")
	}
	sql := sql_QueryXUserInfo + fmt.Sprintf(` WHERE %s=? AND %s=? LIMIT 1;`, UserTab_MOBILE, UserTab_PASSWD)
	resRow, err := xuserDB.userDB.Query(sql, mobile, passwd)
	if err != nil {
		return nil, err
	}
	defer resRow.Close()
	xuserDAO := XUserDAO{}
	err = resRow.Scan(&xuserDAO.UserID, &xuserDAO.Nickname, xuserDAO.Mobile, xuserDAO.Email, xuserDAO.Icon)
	if err != nil {
		return nil, err
	}
	return &xuserDAO, nil
}

func (xuserDB *XUserSQLiteDB) Logout(userID int64, token string) error {
	//TODO
	return nil
}

func (xuserDB *XUserSQLiteDB) Destroy() {
	if xuserDB.userDB != nil {
		err := xuserDB.userDB.Close()
		if err != nil {
			log.Println(err)
		}
		xuserDB.userDB = nil
	}
}
