package logger

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/colinyl/lib4go/concurrent"
)

const (
	ILevel_OFF = iota
	ILevel_Info
	ILevel_Error
	ILevel_Fatal
	ILevel_Debug
	ILevel_ALL
)
const (
	SLevel_OFF   = "Off"
	SLevel_Info  = "Info"
	SLevel_Error = "Error"
	SLevel_Fatal = "Fatal"
	SLevel_Debug = "Debug"
	SLevel_ALL   = "All"
)

type LoggerAppender struct {
	Type  string
	Level string
	Path  string
}
type LoggerLayout struct {
	Level   int
	Content string
}
type LoggerConfig struct {
	Name     string
	Appender *LoggerAppender
}

//LoggerEvent 日志
type LoggerEvent struct {
	Level   string
	Now     time.Time
	Name    string
	Content string
	Path    string
}

//ILogger 日志接口
type ILogger interface {
	Info(content ...interface{})
	Infof(format string, content ...interface{})
	Debug(content ...interface{})
	Debugf(format string, a ...interface{})
	Error(content ...interface{})
	Errorf(format string, a ...interface{})
	Fatal(content ...interface{})
	Fatalf(format string, a ...interface{})
}

//Logger 日志组件
type Logger struct {
	Name       string
	Level      string
	Config     LoggerConfig
	DataChan   chan *LoggerEvent
	OpenSysLog bool
}

var sysDefaultConfig concurrent.ConcurrentMap //map[string]*LoggerConfig
var sysLoggers concurrent.ConcurrentMap       //map[string]*Logger
var levelMap map[string]int

func init() {
	levelMap = map[string]int{
		SLevel_OFF:   ILevel_OFF,
		SLevel_Info:  ILevel_Info,
		SLevel_Error: ILevel_Error,
		SLevel_Fatal: ILevel_Fatal,
		SLevel_Debug: ILevel_Debug,
		SLevel_ALL:   ILevel_ALL,
	}
	sysDefaultConfig = concurrent.NewConcurrentMap()
	sysLoggers = concurrent.NewConcurrentMap()
	readLoggerConfig()

}

//Get 根据sourceName创建新的日志组件
func Get(name string, sourceName string, openSysLog bool) (*Logger, error) {
	return getLogger(name, sourceName, openSysLog)
}

//New 获取指定日志组件
func New(name string, openSysLog bool) (*Logger, error) {
	return getLogger(name, name, openSysLog)
}

//--------------------以下是私有函数--------------------------------------------
func getLogger(name string, sourceName string, openSysLog bool) (logger *Logger, err error) {
	l := sysLoggers.Get(name)
	if l != nil {
		logger = l.(*Logger)
		return
	}
	logger, err = createLogger(name, sourceName, openSysLog)
	if err != nil {
		return nil, err
	}
	sysLoggers.Set(name, logger)
	return
}
func createLogger(name string, sourceName string, openSysLog bool) (log *Logger, err error) {
	objConfig := sysDefaultConfig.Get(sourceName)
	if objConfig == nil {
		objConfig = sysDefaultConfig.Get("*")
	}
	if objConfig == nil {
		return nil, fmt.Errorf("logger %s is invalid", name)
	}
	config := objConfig.(LoggerConfig)
	var dataChan chan *LoggerEvent
	dataChan = make(chan *LoggerEvent, 1000000)
	log = &Logger{Name: name, Level: config.Appender.Level, Config: config,
		DataChan: dataChan, OpenSysLog: openSysLog}
	go FileAppenderWrite(dataChan)
	return
}
func (l *Logger) SetLevel(level string) {
	if _, b := levelMap[level]; !b {
		return
	}
	l.Config.Appender.Level = level
}
func (l *Logger) SetPath(path string) {
	l.Config.Appender.Path = path
}

func (l *Logger) Info(content ...interface{}) {
	l.print(SLevel_Info, fmt.Sprint(content...))
}
func (l *Logger) Infof(format string, content ...interface{}) {
	l.Info(fmt.Sprintf(format, content...))
}

func (l *Logger) Debug(content ...interface{}) {
	l.print(SLevel_Debug, fmt.Sprint(content...))
}
func (l *Logger) Debugf(format string, a ...interface{}) {
	l.Debug(fmt.Sprintf(format, a...))
}
func (l *Logger) Error(content ...interface{}) {
	l.print(SLevel_Error, fmt.Sprint(content...))
}
func (l *Logger) Errorf(format string, a ...interface{}) {
	l.Error(fmt.Sprintf(format, a...))
}
func (l *Logger) Fatal(content ...interface{}) {
	l.print(SLevel_Fatal, fmt.Sprint(content...))
}
func (l *Logger) Fatalf(format string, a ...interface{}) {
	l.Fatal(fmt.Sprintf(format, a...))
}
func (l *Logger) Print(content ...interface{}) {
	l.print(SLevel_Info, fmt.Sprint(content...))
}
func (l *Logger) Printf(format string, a ...interface{}) {
	l.Fatal(fmt.Sprintf(format, a...))
}
func (l *Logger) print(level string, content string) {
	if strings.EqualFold(content, "") {
		return
	}
	events := l.getEvents(level, content)
	for _, event := range events {
		l.DataChan <- event
	}

	if l.OpenSysLog {
		log.SetFlags(log.Ldate | log.Lmicroseconds)
		log.Println(content)
	}

}
func (l *Logger) getEvents(level string, content string) (events map[string]*LoggerEvent) {
	events = make(map[string]*LoggerEvent)
	currentLevel := levelMap[level]
	if currentLevel <= levelMap[l.Level] && currentLevel > ILevel_OFF && currentLevel < ILevel_ALL {
		event := &LoggerEvent{Level: level, Name: l.Name, Now: time.Now(), Content: content,
			Path: l.Config.Appender.Path}
		events[level] = event
	}

	/*for k, v := range levelMap {
		if v <= level && v <= logLevel && v > ILevel_OFF {
			event := &LoggerEvent{Level: k, Name: l.Name, Now: time.Now(), Content: content,
				Path: l.Config.Appender.Path}
			events[k] = event
		}
	}*/
	return
}

//------------------------------------------

func createConfig(config []LoggerConfig) {
	defer func() {
		fmt.Printf("创建日志文件错误:%s\n", recover())
	}()
	data, _ := json.Marshal(config)
	ioutil.WriteFile("lib4go.logger.json", data, os.ModeAppend)
}
func exists(p string) bool {
	_, err := os.Stat(p)
	return err == nil || os.IsExist(err)
}
func readFromFile() ([]LoggerConfig, error) {
	if !exists("./lib4go.logger.json") {
		return nil, errors.New("lib4go.logger.json not exists,using default logger config and create it")
	}
	configs := []LoggerConfig{}
	bytes, err := ioutil.ReadFile("./lib4go.logger.json")
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
	configs[0].Appender = &LoggerAppender{Level: "All", Type: "FileAppender", Path: "./logs/%level/%name/%pid_%date.log"}
	return configs[:]
}
