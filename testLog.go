package main

//import (
//	antlog "antlinker.com/antlog"
//	db "antlinker.com/antmongo"
//	_t "antlinker.com/tools"
//	mgo "gopkg.in/mgo.v2"
//	"io/ioutil"
//)

//var ALOG = antlog.AntLogger

////uid uphoto
//type Model_uphoto struct {
//	UID    string `bson:"UID"`
//	UPHOTO string `bson:"UPHOTO"`
//}

//func main() {
//	var usersAll []Model_uphoto
//	ALOG.Debug("开始取数据")
//	usersAll = QueryAllUsers()
//	for _, user := range usersAll {
//		uid := user.UID
//		uphoto := user.UPHOTO
//		ALOG.Debug("UID:", uid)
//		filecontent, err := _t.Base64Decode(uphoto)
//		if err != nil {
//			ALOG.Error("Base64Decode err", err.Error())
//			panic(err)
//		}
//		filename := "D:/uphoto/" + uid + ".png"
//		ioutil.WriteFile(filename, filecontent, 0x644)
//	}
//}

//func QueryAllUsers() []Model_uphoto {
//	ALOG.Debug("开始QueryAllUsers")
//	var users []Model_uphoto
//	db.OPC("userinfo", func(d *mgo.Collection) {
//		err := d.Find(nil).All(&users)
//		if err != nil {
//			ALOG.Error("insert err", err.Error())
//			panic(err)
//		}
//	})
//	return users
//}
