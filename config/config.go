package config

import (
	logger "tls_epoll_server/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Server   Server   `yaml:"server"`
	Ring     Ring     `yaml:"ring"`
	WorkPool WorkPool `yaml:"workpool"`
	TLS      TLS      `yaml:"tls"`
	Log      Log      `yaml:"log"`
}
type TLS struct {
	Ca string	`default:"" yaml:"ca"`
	Crt string	`default:"" yaml:"crt"`
	Key string	`default:"" yaml:"key"`
	Phrase string `default:"" yaml:"phrase"`
}
type Server struct {
	Address string `default:"" yaml:"address"`
	Port int `default:8080 yaml:"port"`
}
type Ring struct {
	Capacity int    `default:10000 yaml:"capacity"`
	Metric   Metric `yaml:"metric"`
}
type Metric struct {
	Tick int `default:60 yaml:"tick"`
}
type WorkPool struct {
	Capacity int `default:1000 yaml:"capacity"`
	Queue int `default:100000 yaml:"queue"`
	PreMalloc bool `default:true yaml:"preMalloc"`
}
type Log struct {
	Path string `default:"/var/log/server" yaml:"path"`
}
func (this *Config) Init(file string){
	f,err:=ioutil.ReadFile(file)
	if err!=nil {
		panic(err)
	}
	err=yaml.Unmarshal(f,this)
	if err!=nil {
		panic(err)
	}
	//log.Printf("Parse config file success:%+v",this)
	logger.InitLogger(this.Log.Path)
	logger.Logger.Infof("Parse config file success ...")
}
