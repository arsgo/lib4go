package forever

import (
	"fmt"
	"os"
	"syscall"
)

func (f *forever) dup(name string) {
	path := fmt.Sprintf("./%s.dup", f.name)
	fmt.Println("启动调试模式, 所有日志将记录到文件:", path)
	logFile, _ := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_SYNC, 0755)
	syscall.Dup2(int(logFile.Fd()), 1)
	syscall.Dup2(int(logFile.Fd()), 2)
}
