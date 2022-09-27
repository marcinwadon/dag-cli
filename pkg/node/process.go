package node

import (
	"fmt"
	"os"
)

func Start() {
	var attr = os.ProcAttr{
		Dir:   ".",
		Env:   os.Environ(),
		Files: []*os.File{
			os.Stdin,
			nil,
			nil,
		},
		Sys:   nil,
	}
	process, err := os.StartProcess("/bin/sleep", []string{"/bin/sleep", "100"}, &attr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("PID is %d\n", process.Pid)
	err = process.Release()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}