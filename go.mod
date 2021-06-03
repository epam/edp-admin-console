module edp-admin-console

go 1.14

replace (
	github.com/kubernetes-incubator/reference-docs => github.com/kubernetes-sigs/reference-docs v0.0.0-20170929004150-fcf65347b256
	github.com/markbates/inflect => github.com/markbates/inflect v1.0.4
	k8s.io/api => k8s.io/api v0.20.2
	k8s.io/apimachinery => k8s.io/apimachinery v0.20.2
)

require (
	github.com/astaxie/beego v1.12.0
	github.com/blang/semver v3.5.1+incompatible
	github.com/coreos/go-oidc v2.1.0+incompatible
	github.com/docker/docker v1.13.1 // indirect
	github.com/epam/edp-cd-pipeline-operator/v2 v2.3.0-58.0.20210603104955-7b2bc4604c1d
	github.com/epam/edp-codebase-operator/v2 v2.3.0-95.0.20210531100750-614060915d79
	github.com/golang-migrate/migrate v3.5.4+incompatible
	github.com/lib/pq v1.8.0
	github.com/openshift/api v3.9.0+incompatible
	github.com/pkg/errors v0.9.1
	github.com/pquerna/cachecontrol v0.0.0-20180517163645-1555304b9b35 // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/shiena/ansicolor v0.0.0-20151119151921-a422bbe96644 // indirect
	github.com/stretchr/testify v1.6.1
	github.com/xanzy/ssh-agent v0.2.1 // indirect
	go.uber.org/zap v1.15.0
	golang.org/x/net v0.0.0-20210421230115-4e50805a0758
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	gopkg.in/square/go-jose.v2 v2.3.0 // indirect
	gopkg.in/src-d/go-billy.v4 v4.3.0 // indirect
	gopkg.in/src-d/go-git.v4 v4.10.0
	k8s.io/api v0.21.0-rc.0
	k8s.io/apimachinery v0.21.0-rc.0
	k8s.io/client-go v0.20.2
)
