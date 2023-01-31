package main

import (
	"crypto/md5"
	"fmt"
	xdb "github.com/xfort/XNebulaUser/db"
	"log"
	"time"
)

func main() {
	xuserDB, err := xdb.OpenXUserDB("C:\\WORK\\GoCode\\github.com\\xfort\\XNebulaUser")
	if err != nil {
		log.Fatalln(err)
	}

	var xuserDBer xdb.XNebulaUserDBer
	xuserDBer = xuserDB
	userID := time.Now().UnixMicro()

	xuserDAO := xdb.XUserDAO{
		UserID:   userID,
		Nickname: "name",
		Passwd:   fmt.Sprintf("%x", md5.Sum([]byte(time.Now().String()))),
	}
	xuserDBer.Add(&xuserDAO)
	xuserDBer.Destroy()
}
