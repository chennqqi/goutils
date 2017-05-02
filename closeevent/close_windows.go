package closeevent

import (
	"os"
	"os/signal"
	"syscall"
)

func CloseNotify(c chan os.Signal) {
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
}
