package main

import (
	"flag"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
	"strings"
)

func dataTypeMapping() map[string]func(detailType string) (dataType string) {
	dataMap := map[string]func(detailType string) (dataType string){
		"int": func(detailType string) (dataType string) {
			if strings.HasSuffix(detailType, "unsigned") {
				return "uint32"
			} else {
				return "int32"
			}
		},
		"bigint": func(detailType string) (dataType string) {
			if strings.HasSuffix(detailType, "unsigned") {
				return "uint64"
			} else {
				return "int64"
			}
		},
		"smallint": func(detailType string) (dataType string) {
			if strings.HasSuffix(detailType, "unsigned") {
				return "uint16"
			} else {
				return "int16"
			}
		},
		"tinyint": func(detailType string) (dataType string) {
			if strings.HasSuffix(detailType, "unsigned") {
				return "uint8"
			} else {
				return "int8"
			}
		},
	}
	return dataMap
}

type CounterMethod interface {
	// SelectForUpdateByCounterUniqueKey
	//
	// sql(select * from counter where counter_key=@key and msg_id=@msgID for update)
	SelectForUpdateByCounterUniqueKey(key string, msgID uint64) (gen.T, error)
}

func main() {
	var dsn string
	flag.StringVar(&dsn, "dsn", "", "the DSN of db")
	flag.Parse()
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(err)
	}
	g := gen.NewGenerator(gen.Config{
		OutPath: "./ucommon/query",
		// OutFile: "",
		// ModelPkgPath:      "",
		WithUnitTest:      false,
		FieldNullable:     true,
		FieldCoverable:    false,
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,
		Mode:              0,
	})
	g.UseDB(db)
	// 类型映射
	g.WithDataTypeMap(dataTypeMapping())

	// 指定生成表
	counter := g.GenerateModel("counter")
	dConfig := g.GenerateModel("dconfig")
	dconfigEnterprise := g.GenerateModel("dconfig_enterprise")
	g.ApplyBasic(counter)
	g.ApplyBasic(dConfig)
	g.ApplyBasic(dconfigEnterprise)

	// 自定义SQL
	g.ApplyInterface(func(method CounterMethod) {}, counter)

	g.Execute()
}
