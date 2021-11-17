# EDP Admin Console

| :heavy_exclamation_mark: Please refer to [Admin Console user guide](https://epam.github.io/edp-install/user-guide/) to get the notion of the main concepts and guidelines. |
| ---|

## Overview
Admin Console is a central management tool in the EDP ecosystem that provides the ability to deploy infrastructure, project resources and new technologies in a simple way.

## Local Development

### Requirements
* GoLang version higher than 1.13;

>_**NOTE**: The GOPATH and GOROOT environment variables should be added in PATH._
>```
>export GOPATH=C:\Users\<<username>>\go
>export GOROOT=C:\Go
>``

* PostgreSQL client version higher than 9.5;
* Configured access to the VCS, for details, refer to the [Gerrit Setup for Developer](https://kb.epam.com/display/EPMDEDP/Gerrit+Setup+for+Developer) page;
* GoLand Intellij IDEA or another IDE.

## Operator Launch
In order to run the operator, follow the steps below:

1. Clone repository;
2. Open folder in GoLand Intellij IDEA, click the ![add_config_button](readme-resource/add_config_button.png "add_config_button") button and select the **Go Build** option:
   ![add_configuration](readme-resource/add_configuration.png "add_configuration")
3. In Configuration tab, fill in the following:

    3.1. In the Field field, indicate the path to the main.go file;

    3.2. In the Working directory field, indicate the path to the operator;

    3.3. In the Environment field, specify the platform name (OpenShift/Kubernetes) and NameSpace;
   ```
   WATCH_NAMESPACE=test-go-env;PLATFORM_TYPE=openshift
   ```
    ![build-config](readme-resource/build_config.png "build-config")
4. Create the PostgreSQL database, schema, and a user for the EDP Admin Console operator:
     * Create database with a user:
   ```yaml
   CREATE DATABASE edp-db WITH ENCODING 'UTF8';
   CREATE USER postgres WITH PASSWORD 'password';
   GRANT ALL PRIVILEGES ON DATABASE 'edp-db' to postgres;
   ```
     * Create a schema:
   ```yaml
    CREATE SCHEMA [IF NOT EXISTS] 'develop';
   ```
   EDP Admin Console operator supports two modes for running: local and prod.
   For local deploy, modify ```<strong>edp-admin-console/conf/app.conf</strong>``` and set the following parameters:
   ```
    runmode=local
    [local]
    dbEnabled=true
    pgHost=localhost
    pgPort=5432
    pgDatabase=edp-db
    pgUser=postgres
    pgPassword=password
    edpName=develop
   ```
5. Run 'go build main.go' (Shift+F10);
6. After the successful setup, follow the [http://localhost:8080](http://localhost:8080) URL address to check the result:

![check-deploy](readme-resource/check_deploy.png "check-deploy")

## Exceptional Cases
After starting the Go build process, the following error will appear:
```
go: finding github.com/openshift/api v3.9.0
go: finding github.com/openshift/client-go v3.9.0
go: errors parsing go.mod:
C:\Users\<<username>>\Desktop\EDP\edp-admin-console\go.mod:36: require github.com/openshift/api: version "v3.9.0" invalid: unknown revision v3.9.0

Compilation finished with exit code 1
```
To resolve the issue, update the go dependency by applying the Golang command:

```
go get github.com/openshift/api@v3.9.0
```
