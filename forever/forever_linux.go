package forever

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func (f *forever) dup(name string) {
	logFile, _ := os.OpenFile(fmt.Sprintf("./%s.dup", f.name), os.O_WRONLY|os.O_CREATE|os.O_SYNC, 0755)
	syscall.Dup2(int(logFile.Fd()), 1)
	syscall.Dup2(int(logFile.Fd()), 2)
}
