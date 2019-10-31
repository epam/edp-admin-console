/*
 * Copyright 2019 EPAM Systems.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package context

import (
	"edp-admin-console/models/query"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"log"
)

func InitDb() {
	err := orm.RegisterDriver("postgres", orm.DRPostgres)
	checkErr(err, fmt.Sprintf("couldn't register %v driver", "postgres"))

	pgUser := beego.AppConfig.String("pgUser")
	pgPassword := beego.AppConfig.String("pgPassword")
	pgHost := beego.AppConfig.String("pgHost")
	pgDatabase := beego.AppConfig.String("pgDatabase")
	pgPort := beego.AppConfig.String("pgPort")
	pgSchema := Tenant

	params := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s search_path=%s sslmode=disable",
		pgUser, pgPassword, pgHost, pgPort, pgDatabase, pgSchema)

	err = orm.RegisterDataBase("default", "postgres", params)
	checkErr(err, "couldn't register database")
	log.Printf("Connection to %s:%s database is established.", pgHost, pgPort)

	db, err := orm.GetDB("default")
	checkErr(err, "couldnt' get db")
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		pgDatabase, driver)
	checkErr(err, "couldn't initialize migration instance")
	err = m.Up()
	checkErr(err, "db migration is failed")
	debug, err := beego.AppConfig.Bool("ormDebug")
	if err != nil {
		log.Printf("Cannot read orm debug config. Set to false %v", err)
		debug = false
	}
	orm.Debug = debug
	orm.RegisterModel(new(query.Codebase), new(query.ActionLog), new(query.CodebaseBranch), new(query.ThirdPartyService),
		new(query.CDPipeline), new(query.Stage), new(query.QualityGate), new(query.ApplicationsToPromote),
		new(query.CodebaseDockerStream), new(query.GitServer), new(query.JenkinsSlave), new(query.JobProvisioning),
		new(query.EDPComponent))
}

func checkErr(err error, msg string) {
	if err != nil {
		panic(errors.Wrap(err, msg))
	}
}
