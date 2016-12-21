package yamlconfig

import (
	"fmt"
	"testing"
)

type SubConfig struct {
	V string `yaml:"value"`
	T int    `yaml:"type"`
	S string `yame:"chn"`
}

type Config struct {
	Name   string    `json:"name" yaml:"name"`
	Age    int       `json:"age" yaml:"age"`
	Scores []int     ` yaml:",flow""`
	Sub    SubConfig `yaml:"subconfig"`
}

func Test_config(t *testing.T) {
	cfg1 := Config{
		Name:   "Jim",
		Age:    18,
		Scores: []int{100, 90, 88},
		Sub: SubConfig{
			S: "汉语测试",
		},
	}

	fmt.Println(cfg1)
	Save(&cfg1, "t.yml")

	var cfg2 Config
	Load(&cfg2, "t.yml")
	fmt.Println(cfg1)

	if cfg2.Name == cfg1.Name && cfg2.Age == cfg1.Age {
		t.Log("OK")
	}
}
