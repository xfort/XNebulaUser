package db

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"path/filepath"
	"time"
)

//https://support.huaweicloud.com/intl/zh-cn/productdesc-msgsms/phone_numbers.html

const XUSER_DB = "XNebulaUser_2.s3db"

type XNebulaUserDB struct {
	userDB *sql.DB
}

func OpenXUserDB(dbDir string) (*XNebulaUserDB, error) {
	//execABPath, err := os.Executable()
	//if err != nil {
	//	return nil, err
	//}
	//execDirPath := filepath.Dir(execABPath)
	userDBPath := filepath.Join(dbDir, XUSER_DB)
	adminPass := "XNebulaUser"
	authTEMP := md5.Sum([]byte(XUSER_DB))

	for index := 0; index < len(authTEMP); index++ {
		if index%2 == 0 {
			adminPass = adminPass + fmt.Sprintf("%x", authTEMP[index])
		}
	}

	dbURI := fmt.Sprintf("file:%s?_auth&_auth_user=XNebulaUserAdmin&_auth_pass=%s&cache=private", userDBPath, adminPass)
	log.Println("dbURI", dbURI)
	xUserDB, err := sql.Open("sqlite3", dbURI)
	if err != nil {
		return nil, err
	}
	xUserDB.SetMaxOpenConns(1)
	err = createTable(xUserDB)
	if err != nil {
		return nil, err
	}
	return &XNebulaUserDB{userDB: xUserDB}, nil
}

func createTable(userDB *sql.DB) error {
	_, err := userDB.Exec(`CREATE TABLE IF NOT EXISTS  user (
    userID INTEGER PRIMARY KEY, 
	nickname varchar(64) NOT NULL,
    passwd varchar(128) NOT NULL,
    mobile varchar(20),
    email varchar(64),
    icon varchar(255),
    signupUTC INTEGER NOT NULL
     );`)
	return err
}

func (xuserDB *XNebulaUserDB) Add(user *XUserDAO) error {
	_, err := xuserDB.userDB.Exec(`INSERT INTO user(userID,nickname,passwd,mobile,email,icon,signupUTC)
						VALUES(?,?,?,?,?,?,?);`,
		user.UserID, user.Nickname, user.Passwd, user.Mobile, user.Email, user.Icon, time.Now().Unix())
	return err
}

func (xuserDB *XNebulaUserDB) QueryByMobile(mobile string) (*XUserDAO, error) {

	return nil, nil
}

func (xuserDB *XNebulaUserDB) QueryByEmailPasswd(email string, passwd string) (*XUserDAO, error) {

	return nil, nil
}

func (xuserDB *XNebulaUserDB) QueryByMobilePasswd(mobile string, passwd string) (*XUserDAO, error) {

	return nil, nil
}

func (xuserDB *XNebulaUserDB) Logout(userID int64, token string) error {
	return nil
}

func (xuserDB *XNebulaUserDB) Destroy() {
	if xuserDB.userDB != nil {
		err := xuserDB.userDB.Close()
		if err != nil {
			log.Println(err)
		}
		xuserDB.userDB = nil
	}
}
