package kafkapubsub

import (
	goframeworkMQ "git.umu.work/be/goframework/async/mq"
	"git.umu.work/be/goframework/config"
)

type KafkaName string

type KafkaServerConfig struct {
	Version string                    `json:"version" yaml:"version"`
	Config  goframeworkMQ.KafkaConfig `json:"KafkaConfig" yaml:"KafkaConfig"`
}

type KfkConfig struct {
	Name        KafkaName         `json:"Name" yaml:"Name"`
	KafkaConfig KafkaServerConfig `json:"KafkaConfig" yaml:"KafkaConfig"`
}

type KafkaConfig struct {
	Kafuka []KfkConfig `json:"kafuka" yaml:"kafuka"`
}

func NewKafkaConfig() map[KafkaName]KafkaServerConfig {
	var kfkConf KafkaConfig
	conf := config.GetConfig()
	kfkConfMap := make(map[KafkaName]KafkaServerConfig)
	if err := conf.Scan(&kfkConf); err != nil {
		panic(err)
	}
	for _, c := range kfkConf.Kafuka {
		kfkConfMap[c.Name] = c.KafkaConfig
	}

	return kfkConfMap
}
