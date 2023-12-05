package snowflake

import (
	"fmt"
	"sync"
	"time"
)

const (
	// 时间戳位数
	timestampBits = 42
	// 数据中心ID位数
	datacenterIDBits = 5
	// 机器ID位数
	machineIDBits = 5
	// 序列号位数
	sequenceBits = 12
)

// SnowflakeUID 结构体
type SnowflakeUID struct {
	mutex         sync.Mutex
	timestamp     int64
	datacenterID  int64
	machineID     int64
	sequence      int64
	lastTimestamp int64
	timeProvider  func() int64
}

// NewSnowflake 创建一个 SnowflakeUID 实例
func NewSnowflake(datacenterID, machineID int64) *SnowflakeUID {
	return &SnowflakeUID{
		mutex:        sync.Mutex{},
		datacenterID: datacenterID,
		machineID:    machineID,
		timeProvider: defaultTimeProvider,
	}
}

func (s *SnowflakeUID) NewString() string {
	id := s.NewInt()

	return fmt.Sprintf("%d", id)
}

// NewInt 生成唯一ID
func (s *SnowflakeUID) NewInt() int64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	timestamp := s.currentTimestamp()

	if timestamp < s.lastTimestamp {
		timestamp = s.lastTimestamp
	}

	if timestamp == s.lastTimestamp {
		s.sequence = (s.sequence + 1) & ((1 << sequenceBits) - 1)
		if s.sequence == 0 {
			// 自旋等待获取下一个时间戳
			timestamp = s.waitNextMillis(timestamp)
		}
	} else {
		s.sequence = 0
	}

	s.lastTimestamp = timestamp

	id := (timestamp << (datacenterIDBits + machineIDBits + sequenceBits)) |
		(s.datacenterID << (machineIDBits + sequenceBits)) |
		(s.machineID << sequenceBits) |
		s.sequence

	return id
}

// currentTimestamp 获取当前时间戳
func (s *SnowflakeUID) currentTimestamp() int64 {
	if s.timeProvider != nil {
		return s.timeProvider()
	}
	return time.Now().UnixNano() / 1e6
}

// waitNextMillis 等待下一毫秒，直到获取到一个大于lastTimestamp的新时间戳
func (s *SnowflakeUID) waitNextMillis(lastTimestamp int64) int64 {
	timestamp := s.currentTimestamp()
	for timestamp <= lastTimestamp {
		time.Sleep(time.Millisecond)
		timestamp = s.currentTimestamp()
	}
	return timestamp
}

// 默认时间提供者，使用系统当前时间戳
func defaultTimeProvider() int64 {
	return time.Now().UnixNano() / 1e6
}
