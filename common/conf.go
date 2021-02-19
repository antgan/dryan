package common

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

//profile variables
type Conf struct {
	PORT string `yaml:"port"`
}

var Config *Conf

func InitConf() {
	Config = Config.getConf()
}

func (c *Conf) getConf() *Conf {
	yamlFile, err := ioutil.ReadFile("conf.yaml")
	if err != nil {
		fmt.Println(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		fmt.Println(err.Error())
	}
	return c
}
