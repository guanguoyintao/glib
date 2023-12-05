package umiddleware

import (
	"context"
	"fmt"
	"git.umu.work/AI/uglib/ucontext"
	"git.umu.work/AI/uglib/uerrors"
	"git.umu.work/AI/uglib/uhash"
	"git.umu.work/AI/uglib/ujson"
	"git.umu.work/AI/uglib/umetadata"
	"git.umu.work/AI/uglib/umonitor"
	"git.umu.work/AI/uglib/uwrapper"
	"git.umu.work/be/goframework/accelerator/cache"
	"git.umu.work/be/goframework/common"
	"git.umu.work/be/goframework/logger"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/server"
	"time"
)

func NewTimeoutCallWrapper(requestTimeout, dialTimeout time.Duration) client.CallWrapper {
	return func(cf client.CallFunc) client.CallFunc {
		return func(ctx context.Context, node *registry.Node, req client.Request, rsp interface{}, opts client.CallOptions) error {
			opts.RequestTimeout = requestTimeout
			opts.DialTimeout = dialTimeout

			return cf(ctx, node, req, rsp, opts)
		}
	}
}

func NewWithoutTimeoutHandlerWrapper() server.HandlerWrapper {
	return func(h server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			ctx = ucontext.NewUValueContext(ctx)
			return h(ctx, req, rsp)
		}
	}
}

// NewTimeoutTrackerHandlerWrapper 服务超时监控
func NewTimeoutTrackerHandlerWrapper() server.HandlerWrapper {
	return func(h server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) (err error) {
			// 获取 Context 中的截止时间
			deadline, ok := ctx.Deadline()
			if ok {
				// 计算剩余时间
				remainingTime := time.Until(deadline)
				logger.GetLogger(ctx).Info(fmt.Sprintf("Remaining Time: %v\n", remainingTime))
			} else {
				logger.GetLogger(ctx).Info("Deadline not set")
			}
			done := make(chan error)
			uwrapper.GoExecHandlerWithoutTimeout(ctx, func(ctx context.Context) error {
				err := h(ctx, req, rsp)
				if err != nil {
					logger.GetLogger(ctx).Warn(err.Error())
					done <- err
					return err
				}
				done <- nil
				return nil
			})
			for {
				select {
				case <-ctx.Done():
					// 检查上下文错误
					switch ctx.Err() {
					case context.DeadlineExceeded:
						// 上下文超时
						logger.GetLogger(ctx).Info("handler execution timed out")
						return uerrors.UErrorTimeout
					case context.Canceled:
						// 上下文主动取消
						logger.GetLogger(ctx).Info("handler execution canceled")
						return uerrors.UErrorCanceled
					}
				case err = <-done:
					if err != nil {
						logger.GetLogger(ctx).Warn(err.Error())
						return err
					}
					return nil
				}
			}
		}
	}
}

// NewMultiLevelCacheHandlerWrapper 多级缓存
func NewMultiLevelCacheHandlerWrapper(cacheName, serverName string, methodWhitelist []string) server.HandlerWrapper {
	return func(h server.HandlerFunc) server.HandlerFunc {
		// cache aside
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			logger.GetLogger(ctx).Info(fmt.Sprintf("%s multi level cache handler", serverName))
			c, err := cache.GetCache(ctx, cacheName)
			if err != nil {
				logger.GetLogger(ctx).Warn(err.Error())
				return err
			}
			for _, method := range methodWhitelist {
				if req.Method() == method {
					// 生成缓存key
					reqJson, err := ujson.Marshal(req.Body())
					if err != nil {
						logger.GetLogger(ctx).Warn(err.Error())
						return err
					}
					reqHash, err := uhash.HashMurmurHash340(string(reqJson))
					if err != nil {
						logger.GetLogger(ctx).Warn(err.Error())
						return err
					}
					cacheKey := fmt.Sprintf("%s:req:%s:%s", serverName, req.Method(), reqHash)
					logger.GetLogger(ctx).Debug(fmt.Sprintf("cache key is %+v", cacheKey))
					// 判断缓存hit
					ok := c.Exists(ctx, cacheKey)
					if ok {
						// 从缓存中获取结果
						var rspString string
						err = c.Get(ctx, cacheKey, &rspString)
						if err != nil {
							return err
						}
						err = ujson.Unmarshal([]byte(rspString), &rsp)
						if err != nil {
							logger.GetLogger(ctx).Warn(err.Error())
							return err
						}
						return nil
					}
					// 执行handler
					// 如果从缓存中获取数据失败，则从数据库中获取数据
					err = h(ctx, req, rsp)
					if err != nil {
						logger.GetLogger(ctx).Warn(err.Error())
						return err
					}
					uwrapper.GoExecHandlerWithoutTimeout(ctx, func(innerCtx context.Context) error {
						// 设置缓存
						rspJson, err := ujson.Marshal(rsp)
						if err != nil {
							logger.GetLogger(innerCtx).Warn(err.Error())
						}
						data := string(rspJson)
						err = c.Set(&cache.Item{
							Ctx:   innerCtx,
							Key:   cacheKey,
							Value: data,
						})
						if err != nil {
							logger.GetLogger(innerCtx).Warn(err.Error())
						}
						return nil
					})
					return nil
				}
			}

			return h(ctx, req, rsp)
		}
	}
}

// NewPlaybackHandlerWrapper 流量重放
func NewPlaybackHandlerWrapper() server.HandlerWrapper {
	return func(h server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			header := umetadata.ToMap(ctx)
			hr := req.Header()
			if hr != nil {
				for k, v := range hr {
					header[k] = v
				}
			}
			startTime := time.Now()
			logger.GetLogger(ctx).Info(fmt.Sprintf("api request: Service:%+v, Method:%+v, Endpoint:%+v, Header:%+v, Body:%+v",
				req.Service(), req.Method(), req.Endpoint(), header, req.Body()))
			err := h(ctx, req, rsp)
			if err != nil {
				logger.GetLogger(ctx).Error(fmt.Sprintf("api err: %+v", err))
				return err
			}
			logger.GetLogger(ctx).Info(fmt.Sprintf("api response: %+v", rsp), "api_cost", time.Since(startTime).Milliseconds())

			return nil
		}
	}
}

// NewMonitorHandlerWrapper 监控中间件
func NewMonitorHandlerWrapper(serviceName string) server.HandlerWrapper {
	return func(h server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			startTime := time.Now()
			var errCode int64
			var errMsg string
			err := h(ctx, req, rsp)
			if err != nil {
				logger.GetLogger(ctx).Warn(err.Error())
				umuErr := common.AsErrorUmu(err)
				errCode = int64(umuErr.GetCode())
				errMsg = umuErr.Error()
			}
			umonitor.IncRequest(ctx, &umonitor.RequestLabels{
				ServiceName:   serviceName,
				Client:        "rpc",
				InterfaceName: req.Endpoint(),
				ErrCode:       errCode,
				ErrMsg:        errMsg,
				Begin:         startTime,
			})

			return err
		}
	}
}

// NewTracerHandlerWrapper 链路信息中间件
func NewTracerHandlerWrapper() server.HandlerWrapper {
	return func(h server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			header := umetadata.ToMap(ctx)
			logger.GetLogger(ctx).Debug(fmt.Sprintf("metadata is %+v", header))
			ctx = umetadata.Merge(ctx, header)
			err := h(ctx, req, rsp)
			if err != nil {
				return err
			}

			return err
		}
	}
}
