package tasker

import (
	"io/ioutil"
	"path/filepath"
	"runtime"
)

func loadScript(conf string) string {
	_, b, _, _ := runtime.Caller(0)
	script, _ := ioutil.ReadFile(filepath.Dir(b) + conf)

	return string(script)
}
