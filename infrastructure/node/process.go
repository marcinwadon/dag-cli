package node

import (
	"dag-cli/domain/layer"
	"dag-cli/domain/node"
	"dag-cli/infrastructure/config"
	"dag-cli/infrastructure/lb"
	"dag-cli/pkg/pid"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

func GetL0Command(cfg config.Config) []string {
	return []string{
		"/usr/bin/java",
		fmt.Sprintf("-Xms%s", cfg.L0.Java.Xms),
		fmt.Sprintf("-Xmx%s", cfg.L0.Java.Xmx),
		fmt.Sprintf("-Xss%s", cfg.L0.Java.Xss),
		"-jar",
		cfg.GetL0JarFilename(),
		"run-validator",
		"--seedlist",
		cfg.GetL0SeedlistFilename(),
		"--env",
		cfg.Env,
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

func GetL1Command(cfg config.Config, l0Peer node.L0Peer) []string {
	return []string{
		"/usr/bin/java",
		fmt.Sprintf("-Xms%s", cfg.L1.Java.Xms),
		fmt.Sprintf("-Xmx%s", cfg.L1.Java.Xmx),
		fmt.Sprintf("-Xss%s", cfg.L1.Java.Xss),
		"-jar",
		cfg.GetL1JarFilename(),
		"run-validator",
		"--env",
		cfg.Env,
		"--ip",
		cfg.ExternalIP,
		"--public-port",
		fmt.Sprintf("%d", cfg.L1.Port.Public),
		"--p2p-port",
		fmt.Sprintf("%d", cfg.L1.Port.P2P),
		"--cli-port",
		fmt.Sprintf("%d", cfg.L1.Port.CLI),
		"--l0-peer-id",
		l0Peer.Id,
		"--l0-peer-host",
		l0Peer.Host,
		"--l0-peer-port",
		fmt.Sprintf("%d", l0Peer.Port),
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
		p = pid.New(cfg.GetL0PidFilename())
	} else {
		p = pid.New(cfg.GetL1PidFilename())
	}

	err := p.Load()
	if err == nil {
		return fmt.Errorf("PID file exists at %s; process is already running?", p.Path())
	}

	env := GetL0Env(cfg)
	if l == layer.L1 {
		env = GetL1Env(cfg)
	}

	name := "l0"
	if l == layer.L1 {
		name = "l1"
	}

	_ = os.MkdirAll(name, os.ModePerm)

	var sysproc = &syscall.SysProcAttr{
		Credential: nil,
		Noctty:     true,
	}

	filename := filepath.Join(cfg.WorkingDir, fmt.Sprintf("%s.log", name))
	temp, err := os.Create(filename)
	if err != nil {
		return err
	}

	var attr = os.ProcAttr{
		Dir: name,
		Env: env,
		Files: []*os.File{
			os.Stdin,
			temp,
			temp,
		},
		Sys: sysproc,
	}

	var command []string
	if l == layer.L0 {
		command = GetL0Command(cfg)
	} else {
		var l0Peer node.L0Peer
		if cfg.L1.L0Peer.Discovery {
			lbClient := lb.GetClient(cfg.L0.LoadBalancer)

			peer, err := lbClient.GetRandomReadyPeer()
			if err != nil {
				return err
			}

			l0Peer = node.L0PeerFromPeerInfo(*peer)
		} else {
			l0Peer = node.L0Peer{
				Id:   cfg.L1.L0Peer.Id,
				Host: cfg.L1.L0Peer.Host,
				Port: cfg.L1.L0Peer.Port,
			}
		}

		command = GetL1Command(cfg, l0Peer)
	}

	if cfg.Verbose {
		fmt.Printf("Start command:\n%s\n", strings.Join(command, " "))
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
		p = pid.New(cfg.GetL0PidFilename())
	} else {
		p = pid.New(cfg.GetL1PidFilename())
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

	if cfg.Verbose {
		fmt.Printf("Joining to: %s:%d (id: %s)\n", randomPeer.Ip, randomPeer.P2PPort, randomPeer.Id)
	}

	err = nodeClient.Join(randomPeer.Id, randomPeer.Ip, randomPeer.P2PPort)
	if err != nil {
		return err
	}

	return nil
}
