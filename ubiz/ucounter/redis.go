package ucounter

import (
	"context"
	"fmt"
	"git.umu.work/AI/uglib/ubiz"
	mysql "git.umu.work/AI/uglib/ubiz/mysql/ucommon"
	"git.umu.work/AI/uglib/ubiz/mysql/ucommon/model"
	"git.umu.work/AI/uglib/ubiz/mysql/ucommon/query"
	"git.umu.work/AI/uglib/uwrapper"
	"git.umu.work/be/goframework/accelerator/cache"
	"git.umu.work/be/goframework/logger"
	"github.com/go-redis/redis/v8"
	gerrors "github.com/pkg/errors"
	"gorm.io/gorm/clause"
	"strconv"
	"strings"
	"sync"
	"time"
)

// TODO:修改messageID

type RedisCounter struct {
	namespace   CounterNameSpaceType
	dBClient    *mysql.UcommonDBClient
	redisClient redis.Cmdable
}

const (
	cacheNameCounterCache    = "ai-ucounter-cache"
	cacheNameCounterKey      = "ai-ucounter:%s"
	cacheNameCounterWatchKey = "ai-ucounter:%s:%s:watch"
	cacheNameCounterLockKey  = "ai-ucounter:%s:%s:lock"
)

var (
	cronJobManager  = sync.Map{}
	cronJobInterval = 2 * time.Second
	lockTimeout     = 1 * time.Hour
)

func NewRedisCounter(ctx context.Context, namespace CounterNameSpaceType) (ubiz.UCounter, error) {
	redisClient, err := cache.GetRedis(ctx, cacheNameCounterCache)
	if err != nil {
		return nil, err
	}
	counter := &RedisCounter{
		namespace:   namespace,
		dBClient:    mysql.NewUcommonDB(ctx),
		redisClient: redisClient,
	}
	counter.cronJob(ctx)

	return counter, nil
}

func (c *RedisCounter) cronJob(ctx context.Context) {
	_, ok := cronJobManager.LoadOrStore(c.namespace, struct{}{})
	if !ok {
		go func() {
			defer uwrapper.GoroutineRecover(ctx)
			t := time.NewTimer(cronJobInterval)
			defer t.Stop()
			for {
				select {
				case <-t.C:
					err := c.counterCacheGC(ctx)
					if err != nil {
						logger.GetLogger(ctx).Error(err.Error())
					}
					t.Reset(cronJobInterval)
				}
			}
		}()
	}
}

func (c *RedisCounter) counterCacheGC(ctx context.Context) error {
	fields, err := c.redisClient.HKeys(ctx, fmt.Sprintf(cacheNameCounterKey, c.namespace)).Result()
	if err != nil && !gerrors.Is(err, redis.Nil) {
		return err
	}
	wg := sync.WaitGroup{}
	wg.Add(len(fields))
	for _, field := range fields {
		fieldTmp := field
		go func() {
			defer wg.Done()
			err := c.redisCounterPersistenceTx(ctx, fieldTmp)
			if err != nil {
				logger.GetLogger(ctx).Error(err.Error())
			}
		}()
	}
	wg.Wait()

	return nil
}

func (c *RedisCounter) redisCounterPersistenceTx(ctx context.Context, field string) error {
	watchKey := fmt.Sprintf(cacheNameCounterWatchKey, c.namespace, field)
	cacheKey := fmt.Sprintf(cacheNameCounterKey, c.namespace)
	callback := func(tx *redis.Tx) error {
		// 先查询下当前watch监听的key的值
		v, err := tx.HGet(ctx, cacheKey, field).Result()
		if err != nil && !gerrors.Is(err, redis.Nil) {
			return err
		}
		counterNumCache, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			return err
		}
		// 如果key的值没有改变的话，Pipelined函数才会调用成功
		_, err = tx.Pipelined(ctx, func(pipe redis.Pipeliner) error {
			m := query.Use(c.dBClient.DB).Counter
			msgID, err := strconv.ParseUint(field, 10, 64)
			if err != nil {
				return err
			}
			fmt.Println(m.CounterNum.ColumnName())
			err = m.WithContext(ctx).Clauses(clause.OnConflict{
				Columns: []clause.Column{
					{
						Table: m.TableName(),
						Name:  string(m.CounterKey.ColumnName()),
					},
					{
						Table: m.TableName(),
						Name:  string(m.MsgID.ColumnName()),
					},
				},
				DoUpdates: clause.AssignmentColumns([]string{
					string(m.CounterNum.ColumnName()),
				}),
			},
			).Create(&model.Counter{
				MsgID:      msgID,
				CounterKey: string(c.namespace),
				CounterNum: uint32(counterNumCache),
			})
			if err != nil {
				return err
			}
			pipe.HDel(ctx, cacheKey, field)
			// 在事务中修改监听的字符串，如果当前事务执行成功，其他修改该值的事务必然失败，以此达到同步的目的
			pipe.Incr(ctx, watchKey)
			pipe.Expire(ctx, watchKey, lockTimeout)

			return nil
		})
		return err
	}
	// 重试，乐观锁应用于发生数据碰撞概率较低的场景下，所以重试次数一般不会太多
	for i := 0; i < 20; i++ {
		// 使用Watch监听一些Key, 同时绑定一个回调函数callback, 监听Key后的逻辑写在callback这个回调函数里面
		client := c.redisClient.(*redis.Client)
		err := client.Watch(ctx, callback, watchKey)
		if err != nil && err != redis.TxFailedErr {
			return err
		} else if err == nil {
			break
		}
	}

	return nil
}

func (c *RedisCounter) Get(ctx context.Context, key uint64) (uint32, error) {
	ret, err := c.redisClient.HGet(ctx, fmt.Sprintf(cacheNameCounterKey, c.namespace), strconv.FormatUint(key, 10)).Result()
	if err != nil {
		if gerrors.Is(err, redis.Nil) {
			// cache没有去db取
			m := query.Use(c.dBClient.DB).Counter
			do, err := m.WithContext(ctx).Select(m.CounterNum).Where(m.CounterKey.Eq(string(c.namespace)), m.MsgID.Eq(key)).Find()
			if err != nil {
				return 0, err
			}
			var counterNum int32
			if len(do) == 0 {
				counterNum = 0
			} else {
				counterNum = int32(do[0].CounterNum)
			}
			go func() {
				err := c.redisCounterSetTx(ctx, key, counterNum)
				if err != nil {
					logger.GetLogger(ctx).Error(err.Error())
				}
			}()
			return uint32(counterNum), nil
		}
		return 0, err
	}
	counter, err := strconv.ParseUint(ret, 10, 32)
	if err != nil {
		return 0, err
	}

	return uint32(counter), nil
}

func (c *RedisCounter) redisCounterSetTx(ctx context.Context, key uint64, counterNum int32) error {
	field := strconv.FormatUint(key, 10)
	watchKey := fmt.Sprintf(cacheNameCounterWatchKey, c.namespace, field)
	lockKey := fmt.Sprintf(cacheNameCounterLockKey, c.namespace, field)
	cacheKey := fmt.Sprintf(cacheNameCounterKey, c.namespace)
	redisClient := c.redisClient.(*redis.Client)
	ok, err := redisClient.SetNX(ctx, lockKey, time.Now().Unix(), 1*time.Minute).Result()
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}
	callback := func(tx *redis.Tx) error {
		// 如果key的值没有改变的话，Pipelined函数才会调用成功
		res, err := tx.Pipelined(ctx, func(pipe redis.Pipeliner) error {
			// 在事务中修改监听的字符串，如果当前事务执行成功，其他修改该值的事务必然失败，以此达到同步的目的
			pipe.Incr(ctx, watchKey)
			pipe.Expire(ctx, watchKey, lockTimeout)
			pipe.HSet(ctx, cacheKey, field, int64(counterNum))

			return nil
		})
		if err != nil {
			return err
		}
		fmt.Println(res)

		return nil
	}
	// 重试，乐观锁应用于发生数据碰撞概率较低的场景下，所以重试次数一般不会太多
	for i := 0; i < 20; i++ {
		// 使用Watch监听一些Key, 同时绑定一个回调函数callback, 监听Key后的逻辑写在callback这个回调函数里面
		err := redisClient.Watch(ctx, callback, watchKey)
		if err != nil && err != redis.TxFailedErr {
			return err
		} else if err == nil {
			break
		}
	}

	return nil
}

func (c *RedisCounter) Incr(ctx context.Context, key uint64, value uint32) (uint32, error) {
	res, err := c.redisCounterIncrTx(ctx, key, int32(value))
	if err != nil {
		return 0, err
	}

	return res, nil
}

func (c *RedisCounter) redisCounterIncrTx(ctx context.Context, key uint64, increment int32) (uint32, error) {
	redisClient := c.redisClient.(*redis.Client)
	var result int32
	field := strconv.FormatUint(key, 10)
	watchKey := fmt.Sprintf(cacheNameCounterWatchKey, c.namespace, field)
	cacheKey := fmt.Sprintf(cacheNameCounterKey, c.namespace)
	callback := func(tx *redis.Tx) error {
		// 如果key的值没有改变的话，Pipelined函数才会调用成功
		res, err := tx.Pipelined(ctx, func(pipe redis.Pipeliner) error {
			// 在事务中修改监听的字符串，如果当前事务执行成功，其他修改该值的事务必然失败，以此达到同步的目的
			pipe.Incr(ctx, watchKey)
			pipe.Expire(ctx, watchKey, lockTimeout)
			pipe.HIncrBy(ctx, cacheKey, field, int64(increment))
			return nil
		})
		if err != nil {
			return err
		}
		counterNumber, err := strconv.ParseInt(strings.Split(res[2].String(), " ")[4], 10, 64)
		if err != nil {
			return err
		}
		if int32(counterNumber) == increment {
			// 首次增加需要和数据库进行merge
			m := query.Use(c.dBClient.DB).Counter
			do, err := m.WithContext(ctx).Select(m.CounterNum).Where(m.CounterKey.Eq(string(c.namespace)), m.MsgID.Eq(key)).Find()
			if err != nil {
				return err
			}
			if len(do) == 1 {
				result = int32(do[0].CounterNum) + increment
				c.redisCounterSetTx(ctx, key, result)
			} else {
				result = increment
			}
		} else {
			result = int32(counterNumber)
		}

		return nil
	}
	// 重试，乐观锁应用于发生数据碰撞概率较低的场景下，所以重试次数一般不会太多
	for i := 0; i < 20; i++ {
		// 使用Watch监听一些Key, 同时绑定一个回调函数callback, 监听Key后的逻辑写在callback这个回调函数里面
		err := redisClient.Watch(ctx, callback, watchKey)
		if err != nil && err != redis.TxFailedErr {
			return 0, err
		} else if err == nil {
			break
		}
	}

	return uint32(result), nil
}

func (c *RedisCounter) Decr(ctx context.Context, key uint64, value uint32) (uint32, error) {
	var negativeFlag int32 = -1
	res, err := c.redisCounterIncrTx(ctx, key, negativeFlag*int32(value))
	if err != nil {
		return 0, err
	}

	return res, nil
}
