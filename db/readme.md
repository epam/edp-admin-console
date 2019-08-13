For local development you should follow instruction, how to setup migrate cli:

https://github.com/golang-migrate/migrate/tree/master/cmd/migrate

Then you can use this command to verify your migration scripts:

migrate -path=go/src/edp-admin-console/db/migrations -database="postgres://user:password@localhost:5432/edp?sslmode=disable&search_path=main" goto 71