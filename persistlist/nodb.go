// persist list
package persistlist

import (
	"encoding/json"
	"os"

	"fmt"

	"github.com/lunny/nodb"
	"github.com/lunny/nodb/config"
)

const (
	DEFAULT_KEYNAME = "_perist_chan"
)

type PersistList interface {
	Pop(v interface{}) error
	Push(v interface{}) (int64, error)
	Len() (int64, error)
	Close()
}

type nodbList struct {
	key string

	db   *nodb.DB
	inst *nodb.Nodb
}

func NewNodbList(indexDir, keyname string) (PersistList, error) {
	var list nodbList
	cfg := new(config.Config)
	cfg.DataDir = indexDir
	if keyname == "" {
		keyname = DEFAULT_KEYNAME
	}

	err := os.MkdirAll(cfg.DataDir, 0755)
	if !os.IsExist(err) && err != nil {
		fmt.Println("mkdir leveldb dir failed, error: \n", err)
		return nil, err
	}

	dbs, err := nodb.Open(cfg)
	if err != nil {
		fmt.Printf("nodb: error opening db: %v", err)
		return nil, err
	}

	db, _ := dbs.Select(0)

	list.db = db
	list.inst = dbs
	list.key = keyname
	return &list, nil
}

func (c *nodbList) Close() {
	if c.db != nil {
		c.inst.Close()
		c.inst = nil
	}
}

func (c *nodbList) Len() (int64, error) {
	return c.db.LLen([]byte(c.key))
}

func (c *nodbList) Push(v interface{}) (int64, error) {
	//write ToDisk
	txt, _ := json.Marshal(v)
	return c.db.LPush([]byte(c.key), txt)
}

func (c *nodbList) Pop(v interface{}) error {
	//read fromDisk
	txt, err := c.db.RPop([]byte(c.key))
	if err != nil {
		return err
	}
	return json.Unmarshal(txt, v)
}
