package ucounter

import (
	"fmt"
	"git.umu.work/be/goframework/accelerator/cache"
	"git.umu.work/be/goframework/config"
	"git.umu.work/be/goframework/store/gorm"
	"os"
	"path"
	"testing"
)

type commandType string

const (
	incrCommand commandType = "incr"
	decrCommand commandType = "decr"
	getCommand  commandType = "get"
)

func TestMain(m *testing.M) {
	currentPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Println(currentPath)
	config.Init(path.Join(currentPath, "conf"))
	fmt.Printf("config init %+v", config.GetConfig())
	cache.Init(config.GetConfig())
	gorm.Init(config.GetConfig())
	exitCode := m.Run()
	os.Exit(exitCode)
}
