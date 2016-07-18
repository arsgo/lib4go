package logger

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

var loggerPath string

func init() {
	loggerPath, _ = filepath.Abs("./ars.logger.json")
}

//------------------------------------------

func createConfig(config []LoggerConfig) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("创建日志文件错误:%v\n", r)
		}

	}()
	data, _ := json.Marshal(config)
	ioutil.WriteFile(loggerPath, data, os.ModeAppend)
}
func exists(p string) bool {
	_, err := os.Stat(p)
	return err == nil || os.IsExist(err)
}
func readFromFile() ([]LoggerConfig, error) {
	if !exists(loggerPath) {
		return nil, errors.New("lib4go.logger.json not exists,using default logger config and create it")
	}
	configs := []LoggerConfig{}

	bytes, err := ioutil.ReadFile(loggerPath)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(bytes, &configs); err != nil {
		fmt.Println("can't Unmarshal lib4go.logger.json: ", err.Error())
		return nil, err
	}
	return configs, nil
}

func readLoggerConfig() {
	configs, err := readFromFile()
	if err != nil {
		configs = getConfig()
		createConfig(configs)
	}
	for i := 0; i < len(configs); i++ {
		sysDefaultConfig.Set(configs[i].Name, configs[i])
	}
}
func getConfig() []LoggerConfig {
	configs := &[1]LoggerConfig{}
	configs[0] = LoggerConfig{}
	configs[0].Name = "*"
	configs[0].Appender = &LoggerAppender{Level: "All", Type: "FileAppender", Path: "./logs/%level/%name/%date.log"}
	return configs[:]
}
func getCaller(index int) string {
	_, file, line, _ := runtime.Caller(index)
	return fmt.Sprintf("%s%d", filepath.Base(file), line)
}
