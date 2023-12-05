package uruleengine

import (
	"container/list"
	"context"
	"fmt"
	"git.umu.work/AI/uglib/ubiz"
	"git.umu.work/AI/uglib/uerrors"
	"git.umu.work/AI/uglib/ujson"
	"git.umu.work/be/goframework/logger"
	"strings"
	"sync"
)

var (
	rules = make(map[string]ubiz.RuleOperator)
	rw    = &sync.RWMutex{}
)

type Node struct {
	index  int          `json:"-"`
	Key    ubiz.RuleKey `json:"key"`
	Input  []string     `json:"input"`
	Output []string     `json:"output"`
}

type Edge struct {
	Source ubiz.RuleKey       `json:"source"`
	Target ubiz.RuleKey       `json:"target"`
	Weight map[string]float64 `json:"weight"`
}

type Graph struct {
	Nodes   []*Node                                 `json:"nodes"`
	nodeMap map[ubiz.RuleKey]*Node                  `json:"-"`
	Edges   []*Edge                                 `json:"edges"`
	edgeMap map[ubiz.RuleKey]map[ubiz.RuleKey]*Edge `json:"-"`
}

const (
	StartRuleKey ubiz.RuleKey = "start"
	EndRuleKey   ubiz.RuleKey = "end"
)

type RuleEngine struct {
	ctx         context.Context
	namespace   string
	indexInputs map[int][]*Node
	ruleInputs  map[ubiz.RuleKey]map[string]*ubiz.RuleParam
	ruleOutputs map[ubiz.RuleKey]map[string]interface{}
	states      map[string]interface{}
	// dBClient    *mysql.UcommonDBClient
}

func NewRuleEngine(ctx context.Context, namespace string) (ubiz.URuleEngine, error) {

	return &RuleEngine{
		ctx:         ctx,
		namespace:   namespace,
		indexInputs: make(map[int][]*Node),
		ruleInputs:  make(map[ubiz.RuleKey]map[string]*ubiz.RuleParam),
		ruleOutputs: make(map[ubiz.RuleKey]map[string]interface{}),
		states:      make(map[string]interface{}),
		// dBClient:    mysql.NewUcommonDB(ctx),
	}, nil
}

func (r *RuleEngine) RegisterRule(ctx context.Context, key ubiz.RuleKey, rule ubiz.RuleOperator) error {
	rw.Lock()
	rules[r.formatRuleKey(key)] = rule
	rw.Unlock()

	return nil
}

func (r *RuleEngine) formatRuleKey(key ubiz.RuleKey) string {
	return fmt.Sprintf("%+v.%+v", r.namespace, key)
}

func (r *RuleEngine) getAllRuleOperator() []ubiz.RuleOperator {
	rw.RLock()
	globalRuleMap := rules
	rw.RUnlock()
	operators := make([]ubiz.RuleOperator, 0, len(globalRuleMap))
	for key, value := range globalRuleMap {
		if strings.HasPrefix(key, r.namespace) {
			operators = append(operators, value)
		}
	}

	return operators
}

func (r *RuleEngine) registerRuleInput(ctx context.Context, weightMap map[string]float64, node *Node, output map[string]interface{}) {
	inputMap := make(map[string]struct{})
	for _, paramsKey := range node.Input {
		inputMap[paramsKey] = struct{}{}
	}
	for outKey, out := range output {
		_, ok := inputMap[outKey]
		if ok {
			o := out
			key := outKey
			_, ok = r.ruleInputs[node.Key]
			var weight *float64
			if ok {
				if weightMap != nil {
					w, ok := weightMap[key]
					if ok {
						weight = &w
					}
				}
				r.ruleInputs[node.Key][key] = &ubiz.RuleParam{
					Weight: weight,
					Value:  o,
				}
			} else {
				if weightMap != nil {
					w, ok := weightMap[key]
					if ok {
						weight = &w
					}
				}
				r.ruleInputs[node.Key] = map[string]*ubiz.RuleParam{key: {
					Weight: weight,
					Value:  o,
				}}
			}
		}
	}

	return
}

func (r *RuleEngine) calcRuleBFS(ctx context.Context, graph *Graph, startNodes []*Node) error {
	if graph.Nodes == nil {
		return nil
	}
	// 创建一个队列并将起始节点放入队列中
	r.indexInputs[0] = make([]*Node, 0)
	queue := list.New()
	path := ""
	visited := make(map[ubiz.RuleKey]struct{})
	for _, node := range startNodes {
		queue.PushBack(node)
		r.indexInputs[0] = append(r.indexInputs[0], node)
	}

	for queue.Len() > 0 {
		element := queue.Front()
		queue.Remove(element) // 出队

		currentNode := element.Value.(*Node) // 取出队列头部的节点

		// 获取当前节点对应的算子
		rw.RLock()
		operator, ok := rules[r.formatRuleKey(currentNode.Key)]
		rw.RUnlock()
		if !ok {
			return uerrors.UErrorRuleNotRegister
		}

		// 从当前节点中获取输入参数
		input, ok := r.ruleInputs[currentNode.Key]
		if !ok {
			continue
		}

		// 获取已经计算过节点的输出
		output, ok := r.ruleOutputs[currentNode.Key]
		if !ok {
			output = make(map[string]interface{}, 0)
		}
		// 执行当前节点的操作，并获取输出
		operatorReady := true
		for _, inputParam := range currentNode.Input {
			_, ok = input[inputParam]
			if !ok {
				operatorReady = false
				break
			}
		}
		_, ok = visited[currentNode.Key]
		if ok {
			operatorReady = false
		}
		if operatorReady {
			o, err := operator.Operate(r.ctx, r.states, input)
			if err != nil {
				logger.GetLogger(ctx).Warn(err.Error())
				return err
			}
			// 遍历过节点注册 visited 里面
			path += fmt.Sprintf("-->[%+v]%+v{%+v}", input, currentNode.Key, o)
			logger.GetLogger(ctx).Info(fmt.Sprintf("path: %+v", path))
			visited[currentNode.Key] = struct{}{}
			for k, v := range o {
				output[k] = v
			}
		}
		// 将当前节点的输出记录下来
		r.ruleOutputs[currentNode.Key] = output

		// 将当前节点的输出注册到同一层其他节点输入
		inputNodes, ok := r.indexInputs[currentNode.index]
		if ok {
			for _, inputNode := range inputNodes {
				r.registerRuleInput(ctx, nil, inputNode, output)
			}
		}

		edges, ok := graph.edgeMap[currentNode.Key]
		if !ok {
			continue
		}
		// 如果没有计算完的节点重新加入到队列里面
		_, isVisited := visited[currentNode.Key]
		isRepeatable := false
		if queue.Len() > 0 {
			tailNode := queue.Back().Value.(*Node) // 取出队列尾部的节点
			isRepeatable = tailNode.Key == currentNode.Key
		}
		if !isVisited && !isRepeatable {
			queue.PushBack(currentNode)
			continue
		}

		// 处理当前节点的子节点
		for ruleKey, edge := range edges {
			node, ok := graph.nodeMap[ruleKey]
			if !ok {
				return uerrors.UErrorRuleNotRegister
			}
			// 将当前节点的输出注册到子节点的输入中
			r.registerRuleInput(ctx, edge.Weight, node, output)
			// 将子节点加入队列，以便下一轮处理
			if node.index <= currentNode.index {
				node.index = currentNode.index + 1
			}
			queue.PushBack(node)
			_, ok = r.indexInputs[node.index]
			if !ok {
				r.indexInputs[node.index] = make([]*Node, 0)
			}
			r.indexInputs[node.index] = append(r.indexInputs[node.index], node)
		}
	}

	return nil
}

// CalcRuleFlow 计算规则引擎的规则图。
// 根据提供的有向无环图 (graph) 计算规则流程，并以广度优先的方式处理规则，最终返回规则的输出结果。
func (r *RuleEngine) CalcRuleFlow(ctx context.Context, g string, params map[string]interface{}) (map[ubiz.RuleKey]map[string]interface{}, error) {
	var graph Graph
	err := ujson.Unmarshal([]byte(g), &graph)
	if err != nil {
		logger.GetLogger(ctx).Warn(err.Error())
		return nil, err
	}

	// 重构node 数组变成hash map
	graph.nodeMap = make(map[ubiz.RuleKey]*Node, len(graph.Nodes))
	for _, node := range graph.Nodes {
		graph.nodeMap[node.Key] = node
	}

	// 重构edge数组变成二维hash map
	graph.edgeMap = make(map[ubiz.RuleKey]map[ubiz.RuleKey]*Edge, 0)
	for _, edge := range graph.Edges {
		_, ok := graph.edgeMap[edge.Source]
		if !ok {
			graph.edgeMap[edge.Source] = make(map[ubiz.RuleKey]*Edge, 0)
		}
		graph.edgeMap[edge.Source][edge.Target] = edge
	}

	// 为graph中的所有节点注册输入
	startNodes := make([]*Node, 0)
	startRuleEdge, ok := graph.edgeMap[StartRuleKey]
	if !ok {
		return nil, uerrors.UErrorRuleInvalidGraph
	}
	for ruleKey, edge := range startRuleEdge {
		node, ok := graph.nodeMap[ruleKey]
		if !ok {
			return nil, uerrors.UErrorRuleNotRegister
		}
		startNodes = append(startNodes, node)
		r.registerRuleInput(ctx, edge.Weight, node, params)
	}

	// 使用广度优先遍历计算规则
	err = r.calcRuleBFS(ctx, &graph, startNodes)
	if err != nil {
		logger.GetLogger(ctx).Warn(err.Error())
		return nil, err
	}

	// 返回计算后的规则输出
	return r.ruleOutputs, nil
}

func (r *RuleEngine) GetResult(ctx context.Context, key ubiz.RuleKey) (map[string]interface{}, error) {
	res, ok := r.ruleOutputs[key]
	if !ok {
		return nil, uerrors.UErrorRuleHasNoResult
	}

	return res, nil
}
