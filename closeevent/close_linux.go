package closeevent

import (
	"os"
	"os/signal"
	"syscall"
)

/*
syscall.SIGUSR1 linux only
*/

func CloseNotify(c chan os.Signal) {
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)
}
