package jsonconfig

import (
	"fmt"
	"testing"
)

type Config struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func Test_config(t *testing.T) {
	cfg1 := Config{
		Name: "Jim",
		Age:  18,
	}

	fmt.Println(cfg1)
	Save(&cfg1, "")

	var cfg2 Config
	Load(&cfg2, "")
	fmt.Println(cfg1)

	if cfg2 == cfg1 {
		t.Log("OK")
	}
}
