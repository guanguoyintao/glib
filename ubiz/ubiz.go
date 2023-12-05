package ubiz

import (
	"context"
)

type RuleKey string
type RuleParam struct {
	Weight *float64
	Value  interface{}
}

type RuleOperator interface {
	Operate(ctx context.Context, state map[string]interface{}, input map[string]*RuleParam) (map[string]interface{}, error)
}

// UDConfig Dynamic Config
type UDConfig interface {
	RegisterDecoder(ctx context.Context, key string, decodeFunc func(ctx context.Context, config interface{}) (value interface{}, err error)) error
	GetConfig(ctx context.Context, key string) (value interface{}, err error)
	GetEnterpriseConfig(ctx context.Context, enterpriseID uint64, key string) (value interface{}, err error)
}

// UCounter 业务计数器
type UCounter interface {
	Get(ctx context.Context, key uint64) (uint32, error)
	Incr(ctx context.Context, key uint64, value uint32) (uint32, error)
	Decr(ctx context.Context, key uint64, value uint32) (uint32, error)
}

// URuleEngine 规则引擎
type URuleEngine interface {
	RegisterRule(ctx context.Context, key RuleKey, rule RuleOperator) error
	CalcRuleFlow(ctx context.Context, dag string, params map[string]interface{}) (map[RuleKey]map[string]interface{}, error)
	GetResult(ctx context.Context, key RuleKey) (map[string]interface{}, error)
}
