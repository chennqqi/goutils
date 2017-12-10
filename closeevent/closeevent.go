//signal notify close event(win:INT,TREM)
package closeevent

import (
	"os"
	"os/signal"

	"github.com/Sirupsen/logrus"
)

func Wait(stopcall func(os.Signal), signals ...os.Signal) {
	quitChan := make(chan os.Signal, 1)
	defer close(quitChan)
	if len(signals) > 0 {
		signal.Notify(quitChan, signals...)
	} else {
		CloseNotify(quitChan)
	}

	sig := <-quitChan
	logrus.Info("[main:main] quit, recv signal ", sig)
	if stopcall != nil {
		stopcall(sig)
	}
}
