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

	params := fmt.Sprintf("user=%s password=%s host=%s port=5432 dbname=%s sslmode=disable",
		pgUser, pgPassword, pgHost, pgDatabase)

	err = orm.RegisterDataBase("default", "postgres", params)
	checkErr(err)
	err = orm.RunSyncdb("default", false, true)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
