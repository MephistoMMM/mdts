package conf

import "testing"

type TestConf struct {
	K1 int    `json:"Key1"`
	K2 string `json:"Key2"`
}

var defaultConf = TestConf{
	K1: 1,
	K2: "2",
}

func TestLoadConf(t *testing.T) {
	var config TestConf

	if err := LoadConfig("./test/config.json", &defaultConf, &config); err != nil {
		if err == ErrConfNotExist {
			t.Log("Config File Not Exist, Created It.")
		} else {
			t.Error(err)
		}
	}

	t.Log(config)

}
