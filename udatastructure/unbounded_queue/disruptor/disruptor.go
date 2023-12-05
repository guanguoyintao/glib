package disruptorqueue

import (
	"context"
	"runtime"
	"sync/atomic"
)

// 缓存结构体，用于存放元素及相应的序号
type cache struct {
	putNo uint32      // 存放元素的序号
	getNo uint32      // 取出元素的序号
	value interface{} // 存放的值
}

// Queue 无锁队列
type Queue struct {
	capacity uint32  // 队列容量
	capMod   uint32  // 容量模数，用于快速计算位置
	putPos   uint32  // 存放位置
	getPos   uint32  // 取出位置
	cache    []cache // 缓存数组
}

// NewQueue 创建一个新的队列实例
func NewQueue(capacity uint32) *Queue {
	q := new(Queue)
	q.capacity = minQuantity(capacity) // 将容量向上取最近的2的幂
	q.capMod = q.capacity - 1
	q.putPos = 0
	q.getPos = 0
	q.cache = make([]cache, q.capacity)
	for i := range q.cache {
		ca := &q.cache[i]
		ca.getNo = uint32(i)
		ca.putNo = uint32(i)
	}
	ca := &q.cache[0]
	ca.getNo = q.capacity
	ca.putNo = q.capacity
	return q
}

// Enqueue 存放元素到队列
func (q *Queue) Enqueue(ctx context.Context, item interface{}) (ok bool, quantity uint32) {
	var putPos, putPosNew, getPos, posCnt uint32
	var ca *cache
	capMod := q.capMod

	getPos = atomic.LoadUint32(&q.getPos)
	putPos = atomic.LoadUint32(&q.putPos)

	// 计算队列中元素数量
	if putPos >= getPos {
		posCnt = putPos - getPos
	} else {
		posCnt = capMod + (putPos - getPos)
	}

	// 如果队列已满，让出CPU时间片并返回失败
	if posCnt >= capMod-1 {
		runtime.Gosched()
		return false, posCnt
	}

	// 计算新的存放位置
	putPosNew = putPos + 1
	// 使用原子操作尝试更新存放位置，如果失败则让出CPU时间片并返回失败
	if !atomic.CompareAndSwapUint32(&q.putPos, putPos, putPosNew) {
		runtime.Gosched()
		return false, posCnt
	}

	ca = &q.cache[putPosNew&capMod]

	// 自旋等待，直到可以存放元素
	for {
		getNo := atomic.LoadUint32(&ca.getNo)
		putNo := atomic.LoadUint32(&ca.putNo)
		if putPosNew == putNo && getNo == putNo {
			ca.value = item
			atomic.AddUint32(&ca.putNo, q.capacity)
			return true, posCnt + 1
		} else {
			runtime.Gosched()
		}
	}
}

// Dequeue 从队列中取出元素
func (q *Queue) Dequeue(ctx context.Context) (item interface{}, ok bool, quantity uint32) {
	var putPos, getPos, getPosNew, posCnt uint32
	var ca *cache
	capMod := q.capMod

	putPos = atomic.LoadUint32(&q.putPos)
	getPos = atomic.LoadUint32(&q.getPos)

	// 计算队列中元素数量
	if putPos >= getPos {
		posCnt = putPos - getPos
	} else {
		posCnt = capMod + (putPos - getPos)
	}

	// 如果队列为空，让出CPU时间片并返回失败
	if posCnt < 1 {
		runtime.Gosched()
		return nil, false, posCnt
	}

	// 计算新的取出位置
	getPosNew = getPos + 1
	// 使用原子操作尝试更新取出位置，如果失败则让出CPU时间片并返回失败
	if !atomic.CompareAndSwapUint32(&q.getPos, getPos, getPosNew) {
		runtime.Gosched()
		return nil, false, posCnt
	}

	ca = &q.cache[getPosNew&capMod]

	// 自旋等待，直到可以取出元素
	for {
		getNo := atomic.LoadUint32(&ca.getNo)
		putNo := atomic.LoadUint32(&ca.putNo)
		if getPosNew == getNo && getNo == putNo-q.capacity {
			item = ca.value
			ca.value = nil
			atomic.AddUint32(&ca.getNo, q.capacity)
			return item, true, posCnt - 1
		} else {
			runtime.Gosched()
		}
	}
}

// 将容量向上取最近的2的幂
func minQuantity(v uint32) uint32 {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v++
	return v
}
