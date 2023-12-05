package uobserver

type ServiceConf struct {
	RunMode string `yaml:"run_mode" json:"run_mode"`
}

// ServiceConfig 服务相关的配置
type ServiceConfig struct {
	Service ServiceConf `yaml:"service" json:"service"`
}
