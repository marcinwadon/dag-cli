package config

import (
	"dag-cli/pkg/ip"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

func (ac *Autoconfig) AskKey(current KeyConfig) KeyConfig {
	path := ac.Ask("Path to keystore .p12 file", current.Path)
	alias := ac.Ask("Key alias", current.Alias)

	_default := "*****"
	if current.Password == "" {
		_default = ""
	}
	password := ac.Ask("Key password", _default)

	if password == "*****" || password == "" {
		password = current.Password
	}

	abs, _ := filepath.Abs(path)

	return KeyConfig{
		Path:     abs,
		Alias:    alias,
		Password: password,
	}
}

func (ac *Autoconfig) AskWorkingDir(current string) string {
	wd := ac.Ask("Working directory", current)
	path, _ := filepath.Abs(wd)

	return path
}

func (ac *Autoconfig) AskExternalIP(current string) string {
	externalIP := current
	if current == "" {
		externalIP = ip.GetExternalIP()
	}

	ask := ac.Ask("External IP", externalIP)

	return ask
}

func (ac *Autoconfig) AskBlockExplorer(current string) string {
	ask := ac.Ask("Block-explorer URL", current)

	return ask
}

func (ac *Autoconfig) AskGithub(current GithubConfig) GithubConfig {
	org := ac.Ask("Organization", current.Organization)
	repo := ac.Ask("Repository", current.Repository)
	l0Filename := ac.Ask("L0 filename", current.L0Filename)
	l1Filename := ac.Ask("L1 filename", current.L1Filename)
	seedlistFilename := ac.Ask("Seedlist filename", current.SeedlistFilename)

	return GithubConfig{
		Organization:     org,
		Repository:       repo,
		L0Filename:       l0Filename,
		L1Filename:       l1Filename,
		SeedlistFilename: seedlistFilename,
	}
}

func (ac *Autoconfig) AskJava(current JavaConfig) JavaConfig {
	xmx := ac.Ask("Xmx (max. memory allocation)", current.Xmx)
	xms := ac.Ask("Xms (startup memory allocation)", current.Xms)
	xss := ac.Ask("Xss (thread stack size)", current.Xss)

	return JavaConfig{
		Xmx: xmx,
		Xms: xms,
		Xss: xss,
	}
}

func (ac *Autoconfig) AskPorts(current PortConfig) PortConfig {
	publicS := ac.Ask("Public port", fmt.Sprintf("%d", current.Public))
	p2pS := ac.Ask("P2P port", fmt.Sprintf("%d", current.P2P))
	cliS := ac.Ask("CLI port", fmt.Sprintf("%d", current.CLI))

	public, _ := strconv.Atoi(publicS)
	p2p, _ := strconv.Atoi(p2pS)
	cli, _ := strconv.Atoi(cliS)

	return PortConfig{
		Public: public,
		P2P:    p2p,
		CLI:    cli,
	}
}

func (ac *Autoconfig) AskL0(current L0Config) L0Config {
	lb := ac.Ask("Load balancer URL", current.LoadBalancer)
	snapshotsPath := ac.Ask("Path to store snapshots", current.SnapshotsPath)
	snapshotsPathAbs, _ := filepath.Abs(snapshotsPath)

	fmt.Printf("--- L0 Java process --- \n")
	java := ac.AskJava(current.Java)
	ports := ac.AskPorts(current.Port)

	return L0Config{
		SnapshotsPath: snapshotsPathAbs,
		Java:          java,
		LoadBalancer:  lb,
		Port:          ports,
	}
}

func (ac *Autoconfig) AskL0Peer(current L0Peer) L0Peer {
	_default := "yes"
	if !current.Discovery {
		_default = "no"
	}

	discovery := ac.Ask("Auto-discovery L0 peer to fetch data from (yes|no)", _default)
	if strings.ToLower(discovery) == "y" || strings.ToLower(discovery) == "yes" {
		return L0Peer{
			Discovery: true,
		}
	}

	id := ac.Ask("Peer to fetch - ID", current.Id)
	host := ac.Ask("Peer to fetch - Host", current.Host)
	portS := ac.Ask("Peer to fetch - Public port", fmt.Sprintf("%d", current.Port))
	port, _ := strconv.Atoi(portS)

	return L0Peer{
		Discovery: false,
		Id:        id,
		Host:      host,
		Port:      port,
	}
}

func (ac *Autoconfig) AskL1(current L1Config) L1Config {
	lb := ac.Ask("Load balancer URL", current.LoadBalancer)

	fmt.Printf("--- L1 Java process --- \n")
	java := ac.AskJava(current.Java)
	ports := ac.AskPorts(current.Port)

	l0peer := ac.AskL0Peer(current.L0Peer)

	return L1Config{
		L0Peer:       l0peer,
		Java:         java,
		LoadBalancer: lb,
		Port:         ports,
	}
}
