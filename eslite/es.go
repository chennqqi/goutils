//eslite: Package for a simple elasticsearch writer, support es1~es6
package eslite

import (
	"errors"
)

//elasticsearch add pipeline feature after version 5.0
var (
	ErrNotSupportPipeline = errors.New("Only elasticv5,v6 support pipeline")
)

//elasticsearch writer interface
type ESLite interface {
	Open(host string, port int, userName, pass string) error
	Close()
	Begin() error

	SetPipeline(pipeline string) error

	Write(index string, id string,
		typ string, v interface{}) error

	WriteDirect(index string, id string,
		typ string, v interface{}) error

	Commit() error
}
