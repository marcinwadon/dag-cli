package pid

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
)

type PID struct {
	PID  int `json:"pid"`
	path string
}

func New(path string) *PID {
	return &PID{path: path}
}

func (pid *PID) Save(processPID int) error {
	var err error

	pid.PID = processPID

	data, err := json.Marshal(pid)
	if err != nil {
		return err
	}

	path := pid.Path()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}

		return ioutil.WriteFile(path, data, 0644)
	}

	return fmt.Errorf("PID file exists already at '%s'", path)
}

func (pid *PID) Load() error {
	data, err := ioutil.ReadFile(pid.Path())
	if err != nil {
		return fmt.Errorf("no PID file exists at %s; process not running?", pid.Path())
	}

	tmpPid := New(pid.Path())
	err = json.Unmarshal(data, &tmpPid)
	if err != nil {
		return err
	}

	proc, err := tmpPid.Process()
	if err != nil {
		return err
	}
	err = proc.Signal(syscall.Signal(0))
	if err != nil {
		_ = tmpPid.Free()
		return fmt.Errorf("found pid file at %s, but process is not running - removing pid file", pid.Path())
	}

	return json.Unmarshal(data, &pid)
}

func (pid *PID) Free() error {
	if _, err := os.Stat(pid.Path()); os.IsNotExist(err) {
		return nil
	}

	return os.Remove(pid.Path())
}

func (pid *PID) Path() string {
	return pid.path
}

func (pid *PID) Process() (*os.Process, error) {
	if pid.PID == 0 {
		return nil, errors.New("PID has not yet been saved or loaded")
	}

	return os.FindProcess(pid.PID)
}

func (pid *PID) Kill() error {
	proc, err := pid.Process()
	if err != nil {
		return err
	}

	return proc.Kill()
}

func (pid *PID) Signal(sig os.Signal) error {
	proc, err := pid.Process()
	if err != nil {
		return err
	}

	return proc.Signal(sig)
}
