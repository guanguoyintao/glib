package uruleengine

import (
	"context"
	"git.umu.work/AI/uglib/ubiz"
	"reflect"
	"testing"
)

func TestRuleEngine_CalcRuleFlow(t *testing.T) {
	type args struct {
		ctx       context.Context
		namespace string
		dag       string
		params    map[string]interface{}
		operators map[ubiz.RuleKey]ubiz.RuleOperator
	}
	tests := []struct {
		name    string
		args    args
		want    map[ubiz.RuleKey]map[string]interface{}
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				ctx: nil,
				dag: `
					{
						"nodes": [
							{
								"key": "add",
								"input": [
								   "a",
								   "b"
								],
								"output": [
									"a",
								   "c"
								]
							},
							{
								"key": "divide",
								"input": [
									"c"
								],
								"output": [
									"d"
								]
							},
							{
								"key": "times",
								"input": [
									"d",
									"a"
								],
								"output": [
									"e"
								]
							}
						],
						"edges": [
							{
								"source": "start",
								"target": "add"  
							},
							{
								"source": "add",
								"target": "divide"
							},
							{
								"source": "add",
								"target": "times"
							},
							{
								"source": "divide",
								"target": "times"
							}
						]
					}
					`,
				params: map[string]interface{}{
					"a": 30,
					"b": 20,
				},
				operators: map[ubiz.RuleKey]ubiz.RuleOperator{
					"add": func(ctx context.Context, params map[string]*ubiz.RuleParam) (map[string]interface{}, error) {
						a := params["a"].Value
						b := params["b"].Value
						out := make(map[string]interface{})
						out["c"] = a.(int) + b.(int)
						out["a"] = a

						return out, nil
					},
					"times": func(ctx context.Context, params map[string]*ubiz.RuleParam) (map[string]interface{}, error) {
						a := params["a"].Value
						d := params["d"].Value
						out := make(map[string]interface{})

						out["e"] = a.(int) * 100 * d.(int)

						return out, nil
					},
					"divide": func(ctx context.Context, params map[string]*ubiz.RuleParam) (map[string]interface{}, error) {
						c := params["c"].Value
						out := make(map[string]interface{})

						out["d"] = c.(int) / 2

						return out, nil
					},
				},
			},
			want:    map[ubiz.RuleKey]map[string]interface{}{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := NewRuleEngine(tt.args.ctx, tt.args.namespace)
			if err != nil {
				t.Error(err)
			}
			for ruleKey, operator := range tt.args.operators {
				err = r.RegisterRule(tt.args.ctx, ruleKey, operator)
				if err != nil {
					t.Error(err)
				}
			}
			got, err := r.CalcRuleFlow(tt.args.ctx, tt.args.dag, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalcRuleFlow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CalcRuleFlow() got = %v, want %v", got, tt.want)
			}
		})
	}
}
