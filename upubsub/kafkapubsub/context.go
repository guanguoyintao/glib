package kafkapubsub

import (
	"context"
	"fmt"
	"git.umu.work/AI/uglib/uerrors"
	goframeworkMQ "git.umu.work/be/goframework/async/mq"
	"sync"
)

var kfkContext *kafkaContext

type kfkContextKey string

var keyKafukaInContext = kfkContextKey("kfk")

type kafkaContext struct {
	mutex  sync.RWMutex
	kfkMap map[string]*goframeworkMQ.KafkaConfig
}

func (c *kafkaContext) GetKfk(kfkName string) (*goframeworkMQ.KafkaConfig, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if db, ok := c.kfkMap[kfkName]; ok {
		return db, nil
	}

	return nil, fmt.Errorf("erro %w, default kafuka =%s", uerrors.UErrorKafukaNotExitError, kfkName)
}

func (c *kafkaContext) Set(kfkName string, kfk *goframeworkMQ.KafkaConfig) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.kfkMap == nil {
		c.kfkMap = make(map[string]*goframeworkMQ.KafkaConfig)
	}
	c.kfkMap[kfkName] = kfk
}

func (c *kafkaContext) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, keyKafukaInContext, c)
}
