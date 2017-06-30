//a simple json config file load and save function.
package jsonconfig

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
)

func getConfName() string {
	appName := os.Args[0]
	if strings.HasSuffix(appName, ".exe") {
		appName = appName[:len(appName)-len(".exe")]
	}
	return appName + ".json"
}

//load json config file to struct `pv`, if not given fname, use $APPNAME.json
func Load(pv interface{}, fname string) error {
	if fname == "" {
		fname = getConfName()
	}

	txtBytes, err := ioutil.ReadFile(fname)
	if err != nil {
		return err
	}
	return json.Unmarshal(txtBytes, pv)
}

//save struct to fname file, if not given fname, use $APPNAME.json
func Save(v interface{}, fname string) error {
	if fname == "" {
		fname = getConfName()
	}

	txtBytes, _ := json.MarshalIndent(v, "", "\t")
	return ioutil.WriteFile(fname, txtBytes, 0644)
}
