#!/bin/bash
rm ./ucommon/model/*.gen.go
rm ./ucommon/query/*.gen.go
go run generate.go -dsn="root-pd:dWcErvYVXAxnbhYu@tcp(dev_mysql.umucdn.cn:3306)/uai_common?charset=utf8mb4&parseTime=True&loc=Local"