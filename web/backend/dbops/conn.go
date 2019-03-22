/**
* @Author: yanKoo
* @Date: 2019/3/11 11:17
* @Description: 连接数据库
 */
package dbops

import (
	"database/sql"
	"github.com/go-ini/ini"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gomodule/redigo/redis"
	"log"
	"os"
)

var (
	dbConn  *sql.DB
	rdsConn redis.Conn
	err     error
)

/**
 * 连接数据库
 */
func init() {
	cfg, err := ini.Load("../../common/conf/db.ini")  // 编译之后的执行文件所在位置的相对位置
	if err != nil {
		log.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	driver := cfg.Section("db").Key("driver").String()
	user := cfg.Section("db").Key("user").String()
	pwd := cfg.Section("db").Key("password").String()
	host := cfg.Section("db").Key("host").String()
	port := cfg.Section("db").Key("port").String()
	db := cfg.Section("db").Key("db").String()
	maxConn, _ := cfg.Section("db").Key("max_conn").Int()

	dataSource := user + ":" + pwd + "@tcp(" + host + ":" + port + ")/" + db

	dbConn, err = sql.Open(driver, dataSource)
	dbConn.SetMaxOpenConns(maxConn)
	if err != nil {
		panic(err.Error())
	}

	// redis Temporary
	rdsConn, err = redis.Dial("tcp", "172.16.0.74.:6379")
	if err != nil {
		log.Println("connect redis error :", err)
		return
	}
}
