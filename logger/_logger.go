/*
日志处理类，用于记录系统日志,包括，写入文件，数据库，控制台等
*/
package logger
/*
import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"
)

const (
	ILevel_OFF = iota
	ILevel_Debug
	ILevel_Warn
	ILevel_Info
	ILevel_Error
	ILevel_Fatal
	ILevel_ALL
)
const (
	SLevel_OFF   = "Off"
	SLevel_Debug = "Debug"
	SLevel_Warn  = "warn"
	SLevel_Info  = "Info"
	SLevel_Error = "Error"
	SLevel_Fatal = "Fatal"
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

type LoggerEvent struct {
	Level   string
	Now     time.Time
	Name    string
	Content string
	Path    string
}

var sysDefaultConfig map[string]*LoggerConfig
var sysLoggers map[string]*Logger
var logLocker sync.Mutex
var configLocker sync.Mutex
var levelMap []string
var levelIndexs map[string]int

type ILogger interface {
	Info(content ...interface{})
	Infof(format string, content ...interface{})
	Debug(content ...interface{})
	Debugf(format string, a ...interface{})
	Warn(content ...interface{})
	Warnf(format string, a ...interface{})
	Error(content ...interface{})
	Errorf(format string, a ...interface{})
	Fatal(content ...interface{})
	Fatalf(format string, a ...interface{})
}

type Logger struct {
	Name       string
	Level      string
	Config     *LoggerConfig
	DataChan   chan *LoggerEvent
	OpenSysLog bool
}

func init() {
	levelMap = []string{"Debug", "Warn", "Info", "Error", "Fatal"}
	levelIndexs = make(map[string]int, 7)
	levelIndexs["Off"] = ILevel_OFF
	levelIndexs["Debug"] = ILevel_Debug
	levelIndexs["Warn"] = ILevel_Warn
	levelIndexs["Info"] = ILevel_Info
	levelIndexs["Error"] = ILevel_Error
	levelIndexs["Fatal"] = ILevel_Fatal
	levelIndexs["All"] = ILevel_ALL
	sysLoggers = make(map[string]*Logger)
	_, err := readLoggerConfig()
	if err != nil {
		fmt.Println(err)
	}
}
func Get(nName string, sourceName string, openSysLog bool) (*Logger, error) {
	return newLogger(nName, sourceName, openSysLog)
}
func New(name string, openSysLog bool) (*Logger, error) {
	return newLogger(name, name, openSysLog)
}

func (l *Logger) Info(content ...interface{}) {
	l.doWrite(SLevel_Info, fmt.Sprint(content...))
}
func (l *Logger) Infof(format string, content ...interface{}) {
	l.Info(fmt.Sprintf(format, content...))
}

func (l *Logger) Debug(content ...interface{}) {
	l.doWrite(SLevel_Debug, fmt.Sprint(content...))
}
func (l *Logger) Debugf(format string, a ...interface{}) {
	l.Debug(fmt.Sprintf(format, a...))
}
func (l *Logger) Warn(content ...interface{}) {
	l.doWrite(SLevel_Warn, fmt.Sprint(content...))
}
func (l *Logger) Warnf(format string, a ...interface{}) {
	l.Warn(fmt.Sprintf(format, a...))
}
func (l *Logger) Error(content ...interface{}) {
	l.doWrite(SLevel_Error, fmt.Sprint(content...))
}
func (l *Logger) Errorf(format string, a ...interface{}) {
	l.Error(fmt.Sprintf(format, a...))
}
func (l *Logger) Fatal(content ...interface{}) {
	l.doWrite(SLevel_Fatal, fmt.Sprint(content...))
}
func (l *Logger) Fatalf(format string, a ...interface{}) {
	l.Fatal(fmt.Sprintf(format, a...))
}
func (l *Logger) Printf(format string, a ...interface{}) {
	l.Info(fmt.Sprintf(format, a...))
}

func (l *Logger) Close() {
	close(l.DataChan)
}

//--------------------以下是私有函数--------------------------------------------
func newLogger(name string, sourceName string, openSysLog bool) (*Logger, error) {
	logger, b := sysLoggers[name]
	if b {
		return logger, nil
	}
	logger, err := createLogger(name, sourceName)
	if err != nil {
		return nil, err
	}
	logger.OpenSysLog = openSysLog
	sysLoggers[name] = logger
	return logger, nil
}

func createLogger(name string, configName string) (*Logger, error) {
	config, b := sysDefaultConfig[configName]
	if b == false {
		config, b = sysDefaultConfig["*"]
	}
	if b == false {
		return nil, fmt.Errorf(fmt.Sprintf("logger %s is invalid", name))
	}
	var dataChan chan *LoggerEvent
	dataChan = make(chan *LoggerEvent, 100000)
	log := &Logger{Name: name, Level: config.Appender.Level, Config: config,
		DataChan: dataChan}
	go FileAppenderWrite(log.DataChan)
	return log, nil
}

func (l *Logger) doWrite(level string, content string) {
		if strings.EqualFold(content, "") {
			return
		}
		if levelIndexs[l.Level] < levelIndexs[level] {
			return
		}
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("write log exception ", r)
			}
		}()
		event := &LoggerEvent{Level: level, Name: l.Name, Now: time.Now(), Content: content,
			Path: l.Config.Appender.Path}
		l.DataChan <- event
	if l.OpenSysLog {
		log.SetFlags(log.Ldate | log.Lmicroseconds)
		log.Println(content)
	}
}

func getDefaultConfigLogger() []*LoggerConfig {
	configs := &[1]*LoggerConfig{}
	configs[0] = &LoggerConfig{}
	configs[0].Name = "*"
	configs[0].Appender = &LoggerAppender{Level: "All", Type: "FileAppender", Path: "./logs/%level/%name/%pid_%date.log"}
	return configs[:]
}
func createDefautConfig(config []*LoggerConfig) {
	data, _ := json.Marshal(config)
	ioutil.WriteFile("lib4go.logger.json", data, os.ModeAppend)
}
func exists(p string) bool {
	_, err := os.Stat(p)
	return err == nil || os.IsExist(err)
}
func readLoggerConfigFromFile() ([]*LoggerConfig, error) {
	if !exists("./lib4go.logger.json") {
		return nil, errors.New("lib4go.logger.json not exists,using default logger config and create it")
	}

	configs := []*LoggerConfig{}
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

func readLoggerConfig() (map[string]*LoggerConfig, error) {
	if sysDefaultConfig == nil {
		configLocker.Lock()
		defer configLocker.Unlock()
		if sysDefaultConfig == nil {
			sysDefaultConfig = make(map[string]*LoggerConfig)
			configs, err := readLoggerConfigFromFile()
			if err != nil {
				configs = getDefaultConfigLogger()
				createDefautConfig(configs)
			}
			for i := 0; i < len(configs); i++ {
				sysDefaultConfig[configs[i].Name] = configs[i]
			}
		}
	}
	return sysDefaultConfig, nil
}

//----------------------------FileAppenderWriter-----------------------------------------------------
type FileAppenderWriterEntity struct {
	LastUse    int64
	Path       string
	FileEntity *os.File
	Log        *log.Logger
}

//FileAppenderWrite 1. 循环等待写入数据超时时长为1分钟，有新数据过来时先翻译文件输出路径，并查询缓存的实体对象，
//如果存在则调用该对象并输出，不存在则创建, 并输出
//超时后检查所有缓存对象，超过1分钟未使用的请除出缓存，并继续循环
func FileAppenderWrite(dataChan chan *LoggerEvent) {
	appenders := make(map[string]*FileAppenderWriterEntity)
LOOP:
	for {
	FIRSTLOOP:
		for {
			select {
			case data, e := <-dataChan:
				{
					if e {
						wirtelog2file(appenders, data)
					} else {
						break LOOP
					}

				}
			case <-time.After(time.Second * 60):
				{
					break FIRSTLOOP
				}
			}
		}
		//检查超时请求
		currentTime := time.Now().Unix()
		for k, v := range appenders {
			if (currentTime - v.LastUse) >= 60 {
				v.FileEntity.Close()
				delete(appenders, k)
			}
		}
	}

}
func transferPath(event *LoggerEvent) string {

	var resultString string
	resultString = event.Path
	formater := make(map[string]string)
	formater["date"] = time.Now().Format("20060102")
	formater["year"] = time.Now().Format("2006")
	formater["mm"] = time.Now().Format("01")
	formater["mi"] = time.Now().Format("04")
	formater["dd"] = time.Now().Format("02")
	formater["hh"] = time.Now().Format("15")
	formater["ss"] = time.Now().Format("05")
	formater["level"] = event.Level
	formater["name"] = event.Name
	formater["pid"] = fmt.Sprintf("%d", os.Getpid())
	for i, v := range formater {
		match, _ := regexp.Compile("%" + i)
		resultString = match.ReplaceAllString(resultString, v)
	}
	path, _ := filepath.Abs(resultString)
	return path
}
func wirtelog2file(appenders map[string]*FileAppenderWriterEntity, logEvent *LoggerEvent) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("write log exception ", r)
		}
	}()
	path := transferPath(logEvent)
	if appenders[path] == nil {
		ent, err := createFileHandler(path)
		if err != nil {
			fmt.Println(err)
		}
		appenders[path] = ent
	}
	entity, b := appenders[path]
	if b == false {
		return
	}
	if levelIndexs[logEvent.Level] == ILevel_Info {
		entity.Log.SetFlags(log.Ldate | log.Lmicroseconds)
	} else {
		entity.Log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
	}
	entity.Log.Printf("%s\r\n", logEvent.Content)
	entity.LastUse = time.Now().Unix()

}
func createFileHandler(path string) (*FileAppenderWriterEntity, error) {
	dir := filepath.Dir(path)
	er := os.MkdirAll(dir, 0777)
	if er != nil {
		return nil, fmt.Errorf(fmt.Sprintf("can't create dir %s", dir))
	}
	logFile, logErr := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if logErr != nil {
		return nil, fmt.Errorf(fmt.Sprintf("Fail to find file %s", path))
	}
	logger := log.New(logFile, "", log.Ldate|log.Lmicroseconds)
	return &FileAppenderWriterEntity{LastUse: time.Now().Unix(),
		Path: path, Log: logger, FileEntity: logFile}, nil
}
*/