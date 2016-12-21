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

func Save(v interface{}, fname string) error {
	if fname == "" {
		fname = getConfName()
	}

	txtBytes, _ := json.MarshalIndent(v, "", "\t")
	return ioutil.WriteFile(fname, txtBytes, 0644)
}
