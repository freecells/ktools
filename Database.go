/*
 * @Author: Feng
 * @version: v1.0.0
 * @Date: 2020-07-03 13:50:49
 * @LastEditors: Keven
 * @LastEditTime: 2021-09-28 13:12:58
 */
package main

import (
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"

	"gorm.io/driver/sqlserver"

	"gorm.io/gorm"
)

var MainDB *gorm.DB

type DBConfig struct {
	Type     string
	User     string
	Pwd      string
	Database string
	IP       string
	Port     string
}

// func init() {
// 	MainDB = DB()
// }

type lw struct {
}

func (l lw) Println(vals ...interface{}) {
	log.Println("glog:", vals)
}

/*
初始化 db gorm
初始化 gin
*/
func DB(dbc DBConfig) (db *gorm.DB) {

	db, err := getConn(dbc)

	if err != nil {
		log.Println(err.Error())
		return
	}

	return
}

func getConn(dbc DBConfig) (db *gorm.DB, err error) {

	switch dbc.Type {

	case "mysql":
		//注意这里使用的字符集============！！！！！！！！！！！！！！！！！
		// user:password@tcp(localhost:5555)/dbname?tls=skip-verify&autocommit=true
		conStr := dbc.User + ":" + dbc.Pwd + "@tcp(" + dbc.IP + ":" + dbc.Port + ")" + "/" + dbc.Database + "?charset=utf8mb4&parseTime=True&loc=Local"

		db, err = gorm.Open(mysql.Open(conStr), &gorm.Config{
			SkipDefaultTransaction: true,
		})

		if err != nil {
			log.Println(err.Error())
			return
		}

	case "pg":
		// host=myhost port=myport user=gorm dbname=gorm password=mypassword
		conStr := "host=" + dbc.IP + " port=" + dbc.Port + " user=" + dbc.User + " dbname=" + dbc.Database + " password=" + dbc.Pwd
		db, err = gorm.Open(postgres.Open(conStr), &gorm.Config{
			SkipDefaultTransaction: true,
		})

		if err != nil {
			log.Println(err.Error())
			return
		}

	case "sqlite":
		// /tmp/gorm.db
		conStr := "storage/gorm.db"

		db, err = gorm.Open(sqlite.Open(conStr), &gorm.Config{
			SkipDefaultTransaction: true,
		})

		if err != nil {
			log.Println(err.Error())
			return
		}

	case "mssql":

		connFmt := "sqlserver://%s:%s@%s:%s?database=%s&encrypt=disable" // sqlserver2005 需要禁用加密传输

		conStr := fmt.Sprintf(connFmt, dbc.User, dbc.Pwd, dbc.IP, dbc.Port, dbc.Database)

		db, err = gorm.Open(sqlserver.Open(conStr), &gorm.Config{
			SkipDefaultTransaction: true,
		})

		if err != nil {
			log.Println(err.Error())
			return
		}

	}

	return
}

func Rdb() (rdb *redis.Client) {

	rdb = redis.NewClient(&redis.Options{
		Addr:     ":6379",
		DB:       0,
		Password: "",
	})

	return
}
