package mysql

import (
	"fmt"
	"os"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db     *gorm.DB
	dbOnce sync.Once
	//dsn = "username:password@(ip:port)/database?timeout=5000ms&readTimeout=5000ms&writeTimeout=5000ms&charset=utf8mb4&parseTime=true&loc=Local"
	dsn     string
	pswPath = "..\\..\\..\\..\\Key\\mysql.key"
)

func getDB() (*gorm.DB, error) {
	var err error
	psw := readKey()
	dsn = "root:" + psw + "@47.105.123.226/database?timeout=5000ms&readTimeout=5000ms&writeTimeout=5000ms&charset=utf8mb4&parseTime=true&loc=Local"
	dbOnce.Do(func() {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	})
	return db, err
}
func readKey() string {
	data, err := os.ReadFile(pswPath)
	str := ""
	if err != nil {
		fmt.Println(err)
		return str
	}
	str = string(data)
	return str
}
