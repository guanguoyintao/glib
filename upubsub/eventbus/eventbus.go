package eventbus

import (
	"fmt"
	"git.umu.work/AI/uglib/upubsub"
	"sync"
)

// Observer 定义观察者必须实现的接口。
type Observer interface {
	Notify(event upubsub.Event)
}

// EventBus 表示处理事件分发的事件总线。
type EventBus struct {
	observers map[string][]Observer
	mu        sync.RWMutex
}

// NewEventBus 创建一个新的事件总线。
func NewEventBus() *EventBus {
	return &EventBus{
		observers: make(map[string][]Observer),
	}
}

// Subscribe 订阅一个观察者到特定的事件类型。
func (eb *EventBus) Subscribe(eventType string, observer Observer) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if _, ok := eb.observers[eventType]; !ok {
		eb.observers[eventType] = []Observer{observer}
	} else {
		eb.observers[eventType] = append(eb.observers[eventType], observer)
	}
}

// Unsubscribe 从特定的事件类型中取消订阅观察者。
func (eb *EventBus) Unsubscribe(eventType string, observer Observer) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if observers, ok := eb.observers[eventType]; ok {
		for i, obs := range observers {
			if obs == observer {
				// 从观察者列表中删除观察者
				eb.observers[eventType] = append(observers[:i], observers[i+1:]...)
				break
			}
		}
	}
}

// Publish 发布事件到特定事件类型的所有观察者。
func (eb *EventBus) Publish(eventType string, event upubsub.Event) {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	if observers, ok := eb.observers[eventType]; ok {
		// 遍历所有观察者并使用goroutine异步通知它们
		for _, observer := range observers {
			go func(observer Observer) {
				observer.Notify(event)
			}(observer)
		}
	}
}

// ConcreteObserver 是观察者的具体实现示例。
type ConcreteObserver struct {
	Name string
}

// Notify 实现Observer接口的Notify方法。
func (co *ConcreteObserver) Notify(event upubsub.Event) {
	fmt.Printf("%s 收到事件: %+v\n", co.Name, event)
}
