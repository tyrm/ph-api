#!/bin/bash

go get github.com/jinzhu/gorm
go get github.com/juju/loggo
go get gopkg.in/go-oauth2/redis.v1
go get gopkg.in/oauth2.v3
go get gopkg.in/session.v1

CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o puphaus-api .