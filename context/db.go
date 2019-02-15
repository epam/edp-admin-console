package context

import (
	"edp-admin-console/models"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/lib/pq"
	"log"
)

func InitDb() {
	orm.RegisterModel(new(models.EDPTenant))
	err := orm.RegisterDriver("postgres", orm.DRPostgres)
	checkErr(err)

	pgUser := beego.AppConfig.String("pgUser")
	pgPassword := beego.AppConfig.String("pgPassword")
	pgHost := beego.AppConfig.String("pgHost")
	pgDatabase := beego.AppConfig.String("pgDatabase")
	pgPort := beego.AppConfig.String("pgPort")

	params := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		pgUser, pgPassword, pgHost, pgPort, pgDatabase)

	err = orm.RegisterDataBase("default", "postgres", params)
	checkErr(err)
	log.Printf("Connection to %s:%s database is established.", pgHost, pgPort)
	err = orm.RunSyncdb("default", false, true)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
