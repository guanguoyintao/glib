package uobserver

import (
	"context"
	"git.umu.work/be/goframework/config"
	"git.umu.work/be/goframework/logger"
	"github.com/micro/go-micro/v2/debug/profile"
	"github.com/micro/go-micro/v2/debug/profile/http"
)

// ServerObserver 用于监控和记录应用程序的运行时信息，以便进行性能分析、故障排除、日志记录等操作
func ServerObserver(ctx context.Context, serverName string) error {
	c := config.GetConfig()
	var serviceConfig ServiceConfig
	if err := c.Scan(&serviceConfig); err != nil {
		panic(err)
	}
	if serviceConfig.Service.RunMode == "debug" {
		logger.GetLogger(ctx).Info("server run mode is debug")
		go func() {
			pf := http.NewProfile(
				profile.Name(serverName),
			)
			err := pf.Start()
			if err != nil {
				panic(err)
			}
		}()
	}

	return nil
}
