//Package for a simple elasticsearch writer, support es1~es5
package eslite

type ESLite interface {
	Open(host string, port int, userName, pass string) error
	Close()
	Begin() error
	Write(index string, id string,
		typ string, v interface{}) error

	WriteDirect(index string, id string,
		typ string, v interface{}) error

	Commit() error
}
