//a simple yaml config load and save function package
package yamlconfig

import (
	"io/ioutil"

	"github.com/chennqqi/goutils/utils"
	"gopkg.in/yaml.v2"
)

func getConfName() string {
	return utils.ApplicationName() + ".yml"
}

/*
	load configure `fname` to an `pv` interface{}
	@return: error
	@detail: if `fname` is empty, it will auto get`${applicationName}.yml`
*/
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

/*
	save configure to `fname`
	@return: error
	@detail: if `fname` is empty, it will auto get`${applicationName}.yml`
*/
func Save(v interface{}, fname string) error {
	if fname == "" {
		fname = getConfName()
	}

	txtBytes, _ := yaml.Marshal(v)
	return ioutil.WriteFile(fname, txtBytes, 0644)
}
