package uhttp

import (
	"context"
	"git.umu.work/be/goframework/config"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

func NewRestyClient(ctx context.Context) *resty.Client {
	var err error
	var httpConf HttpConfig
	c := config.GetConfig()
	if err = c.Scan(&httpConf); err != nil {
		panic(err)
	}
	if len(httpConf.Http.TimeoutString) == 0 {
		httpConf = HttpConfig{
			Http: HttpConf{
				Timeout:             10 * time.Second,
				MaxIdleConns:        10,
				MaxConnsPerHost:     100,
				MaxIdleConnsPerHost: 10,
			},
		}
	} else {
		httpConf.Http.Timeout, err = time.ParseDuration(httpConf.Http.TimeoutString)
		if err != nil {
			panic(err)
		}
	}
	client := resty.New()
	client.SetDebug(httpConf.Http.Debug)
	client.SetTimeout(httpConf.Http.Timeout)
	// http连接池设置
	httpTransport := http.DefaultTransport.(*http.Transport).Clone()
	httpTransport.MaxIdleConns = int(httpConf.Http.MaxIdleConns)
	httpTransport.MaxConnsPerHost = int(httpConf.Http.MaxConnsPerHost)
	httpTransport.MaxIdleConnsPerHost = int(httpConf.Http.MaxIdleConnsPerHost)
	client.SetTransport(httpTransport)

	return client
}
