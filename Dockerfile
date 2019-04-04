FROM golang:1.10.3-alpine3.8

WORKDIR /go/bin

COPY deployments deployments
COPY edp-admin-console .
COPY static static
COPY views views
COPY conf conf
COPY db db

ENTRYPOINT ["edp-admin-console"]