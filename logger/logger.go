package logger

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/arsgo/lib4go/concurrent"
)

//Logger 日志组件
type Logger struct {
	Name     string
	Level    string
	Config   LoggerConfig
	DataChan chan LoggerEvent
	session  string
}

var sysDefaultConfig *concurrent.ConcurrentMap //map[string]*LoggerConfig
var sysLoggers *concurrent.ConcurrentMap       //map[string]*Logger
var levelMap map[string]int
var sysLogger ILogger
var currentSession int32
var logCreateLock sync.Mutex
var dataChan chan LoggerEvent
var isDebug bool

func init() {
	levelMap = map[string]int{
		SLevel_OFF:   ILevel_OFF,
		SLevel_Info:  ILevel_Info,
		SLevel_Error: ILevel_Error,
		SLevel_Fatal: ILevel_Fatal,
		SLevel_Debug: ILevel_Debug,
		SLevel_ALL:   ILevel_ALL,
	}
	currentSession = 100
	isDebug = true
	sysDefaultConfig = concurrent.NewConcurrentMap()
	sysLoggers = concurrent.NewConcurrentMap()
	dataChan = make(chan LoggerEvent, 1)
	go FileAppenderWrite(dataChan)
	readLoggerConfig()
	sysLogger = &NilLogger{}
	log.SetFlags(log.Ldate | log.Lmicroseconds)
}
func SetDebug(s bool) {
	isDebug = s
}

//Get Logger
func Get(name string) (ILogger, error) {
	return getLogger(name, name, "", true)
}

//New 根据日志组件名称创建新的日志组件
func New(name string) (ILogger, error) {
	return getLogger(name, name, "", false)
}

//NewSession 根据session创建新的日志
func NewSession(name string, session string) (ILogger, error) {	
	return getLogger(name, name, session, false)
}
func (s *Logger) GetName() string {
	return s.Name
}

//--------------------以下是私有函数--------------------------------------------
func getLogger(name string, sourceName string, session string, getFromCache bool) (logger ILogger, err error) {
	if getFromCache {
		logCreateLock.Lock()
		defer logCreateLock.Unlock()
		l := sysLoggers.Get(name)
		if l != nil {
			logger = l.(*Logger)
			return
		}
	}
	logger, err = createLogger(name, sourceName, session)
	if err != nil {
		return sysLogger, err
	}
	if getFromCache {
		sysLoggers.Set(name, logger)
	}
	return
}
func createLogger(name string, sourceName string, session string) (log *Logger, err error) {
	objConfig := sysDefaultConfig.Get(sourceName)
	if objConfig == nil {
		objConfig = sysDefaultConfig.Get("*")
	}
	if objConfig == nil {
		return nil, fmt.Errorf("logger %s is invalid", name)
	}
	config := objConfig.(LoggerConfig)
	log = &Logger{Name: name, Level: config.Appender.Level, Config: config,
		DataChan: dataChan, session: session}
	if strings.EqualFold(session, "") {
		log.session = createSession()
	}

	return
}
func createSession() string {
	return fmt.Sprintf("%s%d", time.Now().Format("150405"), atomic.AddInt32(&currentSession, 1))
}
func (l *Logger) recover() {
	if r := recover(); r != nil {
		sysLogger.Fatal(r)
	}
}
