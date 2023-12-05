// 缓存一致性用cache aside策略
package udconfig

import (
	"context"
	"fmt"
	"git.umu.work/AI/uglib/ubiz"
	mysql "git.umu.work/AI/uglib/ubiz/mysql/ucommon"
	"git.umu.work/AI/uglib/ubiz/mysql/ucommon/query"
	"git.umu.work/AI/uglib/ucontext"
	"git.umu.work/AI/uglib/uwrapper"
	"git.umu.work/be/goframework/accelerator/cache"
	"git.umu.work/be/goframework/logger"
	"github.com/micro/go-micro/v2/errors"
	"net/http"
	"sync"
	"time"
)

type DConfigNameSpaceType string

type DBDConfig struct {
	regedit     map[string]func(ctx context.Context, config interface{}) (value interface{}, err error)
	namespace   DConfigNameSpaceType
	dBClient    *mysql.UcommonDBClient
	cacheClient *cache.Cache
	mu          *sync.Mutex
}

const (
	cacheNameDConfigCache         = "ai-udconfig"
	cacheNameDConfigKey           = "ai-udconfig:%s:%s"
	cacheNameEnterpriseDConfigKey = "ai-enterprise-udconfig:%d:%s:%s"
	cacheConfigExpiryTime         = time.Hour
)

func NewDBDConfig(ctx context.Context, namespace DConfigNameSpaceType) (ubiz.UDConfig, error) {
	cacheClient, err := cache.GetCache(ctx, cacheNameDConfigCache)
	if err != nil {
		return nil, err
	}
	dConfig := &DBDConfig{
		namespace:   namespace,
		dBClient:    mysql.NewUcommonDB(ctx),
		cacheClient: cacheClient,
		regedit:     make(map[string]func(ctx context.Context, config interface{}) (value interface{}, err error)),
		mu:          &sync.Mutex{},
	}

	return dConfig, nil
}

func (c *DBDConfig) RegisterDecoder(ctx context.Context, key string, decodeFunc func(ctx context.Context, config interface{}) (value interface{}, err error)) error {
	c.mu.Lock()
	c.regedit[key] = decodeFunc
	c.mu.Unlock()

	return nil
}

func (c *DBDConfig) GetConfig(ctx context.Context, key string) (value interface{}, err error) {
	decodeFunc, ok := c.regedit[key]
	if !ok {
		logger.GetLogger(ctx).Error(fmt.Sprintf("dynamic config %s, is not registered", key))
		return nil, errors.New("DCONFIG_NOT_REGISTERED_ERROR", fmt.Sprintf("动态配置%s未注册", key), http.StatusInternalServerError)
	}
	// 获取配置
	var conf string
	cacheKey := fmt.Sprintf(cacheNameDConfigKey, c.namespace, key)
	ok = c.cacheClient.Exists(ctx, cacheKey)
	if ok {
		err = c.cacheClient.Get(ctx, cacheKey, &conf)
		if err != nil {
			logger.GetLogger(ctx).Warn(err.Error())
			return nil, err
		}
	} else {
		// cache没有去db取
		m := query.Use(c.dBClient.DB).Dconfig
		dos, err := m.WithContext(ctx).Select(m.Content).Where(m.Namespace.Eq(string(c.namespace)), m.Key.Eq(key)).Order(m.Version.Desc()).Limit(1).Find()
		if err != nil {
			logger.GetLogger(ctx).Warn(err.Error())
			return nil, err
		}
		if len(dos) == 0 {
			return nil, errors.New("DCONFIG_CONFIG_QUANTITY_ERROR", fmt.Sprintf("动态配置，配置数量不正确,为%d个", len(dos)), http.StatusInternalServerError)
		}
		conf = dos[0].Content
		// 写缓存
		uwrapper.GoExecHandlerWithoutTimeout(ucontext.NewUValueContext(ctx), func(ctx context.Context) error {
			err = c.cacheConfig(ctx, key, dos[0].Content)
			if err != nil {
				logger.GetLogger(ctx).Error(err.Error())
			}
			return nil
		})
	}

	// 调用注册的解析函数进行config decode
	result, err := decodeFunc(ctx, conf)
	if err != nil {
		logger.GetLogger(ctx).Warn(err.Error())
		return nil, err
	}

	return result, nil
}

// GetEnterpriseConfig TODO: template func
func (c *DBDConfig) GetEnterpriseConfig(ctx context.Context, enterpriseID uint64, key string) (value interface{}, err error) {
	decodeFunc, ok := c.regedit[key]
	if !ok {
		logger.GetLogger(ctx).Error(fmt.Sprintf("%s dynamic config %s, is not registered", enterpriseID, key))
		return nil, errors.New("DCONFIG_NOT_REGISTERED_ERROR", fmt.Sprintf("动态配置%s未注册", key), http.StatusInternalServerError)
	}
	// 获取配置
	var conf string
	cacheKey := fmt.Sprintf(cacheNameEnterpriseDConfigKey, enterpriseID, c.namespace, key)
	ok = c.cacheClient.Exists(ctx, cacheKey)
	if ok {
		err = c.cacheClient.Get(ctx, cacheKey, &conf)
		if err != nil {
			logger.GetLogger(ctx).Warn(err.Error())
			return nil, err
		}
	} else {
		// cache没有去db取
		m := query.Use(c.dBClient.DB).DconfigEnterprise
		dos, err := m.WithContext(ctx).Select(m.Content).Where(
			m.EnterpriseID.Eq(int64(enterpriseID)),
			m.Namespace.Eq(string(c.namespace)),
			m.Key.Eq(key),
		).Order(m.Version.Desc()).Limit(1).Find()
		if err != nil {
			logger.GetLogger(ctx).Warn(err.Error())
			return nil, err
		}
		if len(dos) == 0 {
			return nil, errors.New("DCONFIG_CONFIG_QUANTITY_ERROR", fmt.Sprintf("动态配置，配置数量不正确,为%d个", len(dos)), http.StatusInternalServerError)
		}
		conf = dos[0].Content
		// 写缓存
		uwrapper.GoExecHandlerWithoutTimeout(ucontext.NewUValueContext(ctx), func(ctx context.Context) error {
			err = c.cacheEnterpriseConfig(ctx, enterpriseID, key, dos[0].Content)
			if err != nil {
				logger.GetLogger(ctx).Error(err.Error())
			}
			return nil
		})
	}

	// 调用注册的解析函数进行config decode
	result, err := decodeFunc(ctx, conf)
	if err != nil {
		logger.GetLogger(ctx).Warn(err.Error())
		return nil, err
	}

	return result, nil
}

func (c *DBDConfig) cacheEnterpriseConfig(ctx context.Context, enterpriseID uint64, key, conf string) (err error) {
	cacheKey := fmt.Sprintf(cacheNameEnterpriseDConfigKey, enterpriseID, c.namespace, key)
	err = c.cacheClient.Set(&cache.Item{
		Ctx:   ctx,
		Key:   cacheKey,
		Value: conf,
		TTL:   cacheConfigExpiryTime,
	})
	if err != nil {
		logger.GetLogger(ctx).Warn(err.Error())
		return err
	}

	return nil
}

func (c *DBDConfig) cacheConfig(ctx context.Context, key, conf string) (err error) {
	cacheKey := fmt.Sprintf(cacheNameDConfigKey, c.namespace, key)
	err = c.cacheClient.Set(&cache.Item{
		Ctx:   ctx,
		Key:   cacheKey,
		Value: conf,
		TTL:   cacheConfigExpiryTime,
	})
	if err != nil {
		logger.GetLogger(ctx).Warn(err.Error())
		return err
	}

	return nil
}
