package script

import (
	"os"
	"time"

	"github.com/arsgo/lib4go/concurrent"
)

type LuaScriptWatch struct {
	callback func()
	files    *concurrent.ConcurrentMap
	lastTime time.Time
}

//NewLuaScriptWatch 构建脚本监控文件
func NewLuaScriptWatch(callback func()) *LuaScriptWatch {
	w := &LuaScriptWatch{callback: callback, lastTime: time.Now()}
	w.files = concurrent.NewConcurrentMap()
	go w.watch()
	return w
}

//AppendFile 添加监控文件
func (w *LuaScriptWatch) AppendFile(path string) {
	w.files.Set(path, path)
}

func (w *LuaScriptWatch) watch() {
	tk := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-tk.C:
			if !w.checkFile() {
				w.callback()
			}
		}
	}
}

//checkFile 检查文件最后修改时间
func (w *LuaScriptWatch) checkFile() bool {
	files := w.files.GetAll()
	for path := range files {
		fileinfo, err := os.Stat(path)
		if err != nil {
			continue
		}
		if fileinfo.ModTime().Sub(w.lastTime) > 0 {
			w.lastTime = fileinfo.ModTime()
			return false
		}
	}
	return true
}
