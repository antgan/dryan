package common

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
)

//profile variables
type conf struct {
	PORT      string     `yaml:"port"`
	DATABASES []DataBase `yaml:"databases,flow"`
}

type DataBase struct {
	KEY      string                 `yaml:"key"`
	TYPE     string                 `yaml:"type"`
	HOST     string                 `yaml:"host"`
	NAME     string                 `yaml:"name"`
	USER     string                 `yaml:"user"`
	PASSWORD string                 `yaml:"password"`
	EXT      map[string]interface{} `yaml:",flow"`
}

var Config = &conf{}

func InitConf() {
	Config = Config.getConf()
}

func (c *conf) getConf() *conf {
	yamlFile, err := ioutil.ReadFile("conf.yml")
	if err != nil {
		fmt.Println(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		fmt.Println(err.Error())
	}
	return c
}

func init() {
	InitConf()
}

// Get value of given key of Database section.
func (d *DataBase) Ext(key string, defaultVal ...interface{}) interface{} {
	if v, exist := d.EXT[key]; exist {
		return v
	}

	if len(defaultVal) > 0 {
		return defaultVal[0]
	}

	return ""
}

func (d *DataBase) ExtString(key string, defaultVal ...interface{}) string {
	keyValue := d.Ext(key, defaultVal...)

	return fmt.Sprint(keyValue)
}

func (d *DataBase) ExtInt(key string, defaultVal ...interface{}) int {
	return d.Ext(key, defaultVal...).(int)
}

func (d *DataBase) ExtDuration(key string, defaultVal ...interface{}) time.Duration {
	keyValue := d.Ext(key, defaultVal...)

	str := fmt.Sprint(keyValue)
	t, _ := time.ParseDuration(str)
	return t
}
