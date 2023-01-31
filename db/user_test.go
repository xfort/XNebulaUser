package db

import (
	"crypto/md5"
	"fmt"
	"testing"
	"time"
)

func TestXUserS3DB(t *testing.T) {
	xuserDB, err := OpenXUserDB("C:\\WORK\\GoCode\\github.com\\xfort\\XNebulaUser")
	if err != nil {
		t.Fatal(err)
	}

	for index := 0; index < 5; index++ {
		userID := time.Now().UnixMicro()
		xuserDAO := XUserDAO{
			UserID:   userID,
			Nickname: "name",
			Passwd:   fmt.Sprintf("%x", md5.Sum([]byte(time.Now().String()))),
		}
		_, err := xuserDB.Add(&xuserDAO)
		if err != nil {
			t.Fatal(err)
		}
		time.Sleep(2 * time.Second)
	}
	defer xuserDB.Destroy()
}
