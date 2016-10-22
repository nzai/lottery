package conn

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nzai/lottery/config"
)

const (
	defaultConnectionString = "lottery.db"
)

//	获取数据库连接
func GetConn() (*sql.DB, error) {

	//	连接字符串
	connectionString := config.String("database", "conn", defaultConnectionString)

	//	创建连接
	db, err := sql.Open("sqlite3", connectionString)
	if err != nil {
		log.Print("创建数据库连接失败: ", err)
		return nil, err
	}

	//	测试一下连接
	err = db.Ping()
	if err != nil {
		log.Print("连接数据库失败: ", err)
		return nil, err
	}

	return db, err
}
