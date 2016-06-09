package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/colinyl/lib4go/concurrent"
)

var fileAppenders concurrent.ConcurrentMap

func init() {
	fileAppenders = concurrent.NewConcurrentMap()
}

//FileAppenderWriterEntity fileappender
type FileAppenderWriterEntity struct {
	LastUse    int64
	Path       string
	FileEntity *os.File
	Log        *log.Logger
	Data       chan *LoggerEvent
	Close      chan int
}

func getFileAppender(data *LoggerEvent) (f *FileAppenderWriterEntity, err error) {
	path := getAppendPath(data)
	entity := fileAppenders.Get(path)
	if entity != nil {
		f = entity.(*FileAppenderWriterEntity)
		return
	}
	entity, err = createFileHandler(path)
	if err != nil {
		return
	}
	fileAppenders.Set(path, entity)
	f = entity.(*FileAppenderWriterEntity)
	go f.writeLoop()
	go f.checkAppender()

	return
}

//FileAppenderWrite 1. 循环等待写入数据超时时长为1分钟，有新数据过来时先翻译文件输出路径，并查询缓存的实体对象，
//如果存在则调用该对象并输出，不存在则创建, 并输出
//超时后检查所有缓存对象，超过1分钟未使用的请除出缓存，并继续循环
func FileAppenderWrite(dataChan chan *LoggerEvent) {
	for {
		select {
		case data, b := <-dataChan:
			{
				if b {
					f, er := getFileAppender(data)
					if er == nil {
						f.Data <- data
					}
				}
			}
		}
	}
}
func getAppendPath(event *LoggerEvent) string {
	var resultString string
	resultString = event.Path
	formater := make(map[string]string)
	formater["date"] = event.Now.Format("20060102")
	formater["year"] = event.Now.Format("2006")
	formater["mm"] = event.Now.Format("01")
	formater["mi"] = event.Now.Format("04")
	formater["dd"] = event.Now.Format("02")
	formater["hh"] = event.Now.Format("15")
	formater["ss"] = event.Now.Format("05")
	formater["level"] = strings.ToLower(event.Level)
	formater["name"] = event.Name
	formater["pid"] = fmt.Sprintf("%d", os.Getpid())
	for i, v := range formater {
		match, _ := regexp.Compile("%" + i)
		resultString = match.ReplaceAllString(resultString, v)
	}
	path, _ := filepath.Abs(resultString)
	return path
}
func (entity *FileAppenderWriterEntity) checkAppender() {
	ticker := time.NewTicker(time.Minute)
LOOP:
	for {
		select {
		case <-ticker.C:
			{
				currentTime := time.Now().Unix()
				if (currentTime - entity.LastUse) >= 60 {
					fileAppenders.Delete(entity.Path)
					entity.FileEntity.Close()
					break LOOP
				}
			}
		}
	}
}
func (entity *FileAppenderWriterEntity) writeLoop() {
LOOP:
	for {
		select {
		case e := <-entity.Data:
			{
				entity.writelog2file(e)
			}
		case <-entity.Close:
			break LOOP
		}

	}
}

func (entity *FileAppenderWriterEntity) writelog2file(logEvent *LoggerEvent) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("write log exception ", r)
		}
	}()
	if levelMap[logEvent.Level] == ILevel_Info {
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
		Path: path, Log: logger, FileEntity: logFile, Data: make(chan *LoggerEvent, 1000000),
		Close: make(chan int, 1)}, nil
}
