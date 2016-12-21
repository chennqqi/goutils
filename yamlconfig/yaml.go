package yamlconfig

import (
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

func getConfName() string {
	appName := os.Args[0]
	if strings.HasSuffix(appName, ".exe") {
		appName = appName[:len(appName)-len(".exe")]
	}
	return appName + ".yml"
}

func Load(pv interface{}, fname string) error {
	if fname == "" {
		fname = getConfName()
	}

	txtBytes, err := ioutil.ReadFile(fname)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(txtBytes, pv)
}

func Save(v interface{}, fname string) error {
	if fname == "" {
		fname = getConfName()
	}

	txtBytes, _ := yaml.Marshal(v)
	return ioutil.WriteFile(fname, txtBytes, 0644)
}
