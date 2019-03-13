package context

import (
	"github.com/astaxie/beego/orm"
	"io/ioutil"
	"log"
)

func TryToCreateTables() {
	log.Println("Try to create tables...")
	bytes, err := ioutil.ReadFile("deployments/init.sql")
	checkErr(err)

	o := orm.NewOrm()
	_, err = o.Raw(string(bytes)).Exec()
	checkErr(err)
}
