// antmongo
package antmongo

import (
	_f "antlinker.com/antlog"
	_t "antlinker.com/tools"
	"gopkg.in/mgo.v2"
	"os"
	"strconv"
)

var log = _f.AntLogger
var (
	session   *mgo.Session
	dburl     string = "mongodb://goServer:goServer@localhost:27017,goServer:goServer@localhost:27017"
	dbname    string = "test"
	poolLimit int    = 100
)

// 读取配置文件,初始化mongodb连接配置参数
func init() {
	conf := _t.GetCurrentPath() + "/config/antmongo.conf"
	log.Debug("MongoDB初始化读取配置文件:%s", conf)
	p, err := _t.ReadConfig(conf)
	if err != nil {
		log.Error("读取MongoDB配置文件错误:%s", err.Error())
		os.Exit(1)
	}
	dburl = p["MGO_URL"]
	dbname = p["MGO_DB_NAME"]
	i, e := strconv.Atoi(p["MGO_POOL_LIMIT"])
	if e != nil {
		log.Error("读取MongoDB配置文件错误:%s", e.Error())
	}
	poolLimit = i
	log.Debug("MGO_URL = %s", dburl)
	log.Debug("MGO_DB_NAME = %s", dbname)
	log.Debug("MGO_POOL_LIMIT = %d", i)
}

// 获取Mongodb的session，如果为空则重新创建，否则返回原有值
func GetSession() (*mgo.Session, error) {
	if session == nil {
		var err error
		session, err = mgo.Dial(dburl)
		if err != nil {
			log.Error("连接mongodb错误:%s", err.Error())
			return nil, err
		}
		session.SetPoolLimit(poolLimit)
	}
	return session.Clone(), nil
}

// 定义执行Mongodb操作的函数操作完成后关闭session
// 传入collection的名字和操作函数即可
func OPC(collection string, f func(*mgo.Collection)) {
	session, err := GetSession()
	if err != nil {
		log.Error("获取mongodb连接错误:%s", err.Error())
	}
	defer func() {
		session.Close()
		if err := recover(); err != nil {
			log.Error("执行mongodb操作错误:%s", err)
		}
	}()
	c := session.DB(dbname).C(collection)
	f(c)
}

// 定义执行Mongodb操作的函数操作完成后关闭session
// 传入数据库名字和collection的名字，外加操作函数即可
func OPDC(newdb string, collection string, f func(*mgo.Collection)) {
	session, err := GetSession()
	if err != nil {
		log.Error("获取mongodb连接错误:%s", err.Error())
	}
	defer func() {
		session.Close()
		if err := recover(); err != nil {
			log.Error("执行mongodb操作错误:%s", err)
		}
	}()
	c := session.DB(newdb).C(collection)
	f(c)
}
