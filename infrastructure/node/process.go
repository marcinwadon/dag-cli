package node

import (
	"dag-cli/domain/layer"
	"dag-cli/infrastructure/config"
	"dag-cli/infrastructure/lb"
	"dag-cli/pkg/pid"
	"fmt"
	"os"
	"syscall"
)

func GetL0Command(cfg config.Config) []string {
	return []string{
		"/usr/bin/java",
		fmt.Sprintf("-Xms%s", cfg.L0.Java.Xms),
		fmt.Sprintf("-Xmx%s", cfg.L0.Java.Xmx),
		fmt.Sprintf("-Xss%s", cfg.L0.Java.Xss),
		"-jar",
		cfg.L0.Path,
		"run-validator",
		"--seedlist",
		cfg.L0.SeedlistPath,
		"--env",
		cfg.L0.Env,
		"--ip",
		cfg.ExternalIP,
		"--public-port",
		fmt.Sprintf("%d", cfg.L0.Port.Public),
		"--p2p-port",
		fmt.Sprintf("%d", cfg.L0.Port.P2P),
		"--cli-port",
		fmt.Sprintf("%d", cfg.L0.Port.CLI),
	}
}

func GetL0Env(cfg config.Config) []string {
	return append([]string{
		fmt.Sprintf("CL_KEYSTORE=%s", cfg.Key.Path),
		fmt.Sprintf("CL_KEYALIAS=%s", cfg.Key.Alias),
		fmt.Sprintf("CL_PASSWORD=%s", cfg.Key.Password),
		fmt.Sprintf("CL_SNAPSHOT_STORED_PATH=%s", cfg.L0.SnapshotsPath),
	}, os.Environ()...)
}

func GetL1Command(cfg config.Config) []string {
	return []string{
		"/usr/bin/java",
		fmt.Sprintf("-Xms%s", cfg.L0.Java.Xms),
		fmt.Sprintf("-Xmx%s", cfg.L0.Java.Xmx),
		fmt.Sprintf("-Xss%s", cfg.L0.Java.Xss),
		"-jar",
		cfg.L1.Path,
		"run-validator",
		"--env",
		cfg.L1.Env,
		"--ip",
		cfg.ExternalIP,
		"--public-port",
		fmt.Sprintf("%d", cfg.L1.Port.Public),
		"--p2p-port",
		fmt.Sprintf("%d", cfg.L1.Port.P2P),
		"--cli-port",
		fmt.Sprintf("%d", cfg.L1.Port.CLI),
		"--l0-peer-id",
		cfg.L1.L0Peer.Id,
		"--l0-peer-host",
		cfg.L1.L0Peer.Host,
		"--l0-peer-port",
		fmt.Sprintf("%d", cfg.L1.L0Peer.Port),
	}
}

func GetL1Env(cfg config.Config) []string {
	return append([]string{
		fmt.Sprintf("CL_KEYSTORE=%s", cfg.Key.Path),
		fmt.Sprintf("CL_KEYALIAS=%s", cfg.Key.Alias),
		fmt.Sprintf("CL_PASSWORD=%s", cfg.Key.Password),
	}, os.Environ()...)
}

func Start(cfg config.Config, l layer.Layer) error {
	var p *pid.PID

	if l == layer.L0 {
		p = pid.New(cfg.L0.PidPath)
	} else {
		p = pid.New(cfg.L1.PidPath)
	}

	err := p.Load()
	if err == nil {
		return fmt.Errorf("PID file exists at %s; process is already running?", p.Path())
	}

	env := GetL0Env(cfg)
	if l == layer.L1 {
		env = GetL1Env(cfg)
	}

	dir := "l0"
	if l == layer.L1 {
		dir = "l1"
	}

	_ = os.MkdirAll(dir, os.ModePerm)

	var attr = os.ProcAttr{
		Dir: dir,
		Env: env,
		Files: []*os.File{
			os.Stdin,
			os.Stdout,
			os.Stderr,
		},
		Sys: nil,
	}

	var command []string
	if l == layer.L0 {
		command = GetL0Command(cfg)
	} else {
		command = GetL1Command(cfg)
	}

	process, err := os.StartProcess(command[0], command, &attr)
	if err != nil {
		return err
	}

	err = p.Save(process.Pid)
	if err != nil {
		return err
	}

	err = process.Release()
	if err != nil {
		return err
	}

	return nil
}

func Stop(cfg config.Config, l layer.Layer) error {
	var p *pid.PID

	if l == layer.L0 {
		p = pid.New(cfg.L0.PidPath)
	} else {
		p = pid.New(cfg.L1.PidPath)
	}

	if err := p.Load(); err != nil {
		return fmt.Errorf("no PID file exists at %s; process not running?", p.Path())
	}

	proc, err := os.FindProcess(p.PID)
	if err != nil {
		return err
	}

	err = proc.Signal(syscall.SIGTERM)
	if err != nil {
		return err
	}

	err = p.Free()
	if err != nil {
		return err
	}

	return nil
}

func Join(cfg config.Config, l layer.Layer) error {
	lbClient := lb.GetClient(cfg.L0.LoadBalancer)
	if l == layer.L1 {
		lbClient = lb.GetClient(cfg.L1.LoadBalancer)
	}

	randomPeer, err := lbClient.GetRandomReadyPeer()
	if err != nil {
		return err
	}

	host := "127.0.0.1"
	nodeClient := GetClient(host, cfg.L0.Port.CLI)
	if l == layer.L1 {
		nodeClient = GetClient(host, cfg.L1.Port.CLI)
	}

	err = nodeClient.Join(randomPeer.Id, randomPeer.Ip, randomPeer.P2PPort)
	if err != nil {
		return err
	}

	return nil
}
