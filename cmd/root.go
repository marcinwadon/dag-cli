package cmd

import (
	node2 "dag-cli/domain/node"
	"dag-cli/infrastructure/config"
	"dag-cli/infrastructure/lb"
	"dag-cli/infrastructure/node"
	"dag-cli/pkg/color"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use:     "dag",
		Short:   "DAG Command Line Utility",
		Version: "v0.7.0",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.LoadConfig()
			if err != nil {
				return err
			}

			loc, _ := time.LoadLocation("UTC")
			now := time.Now().In(loc)

			fmt.Printf("External IP: %s\n", color.Ize(color.Yellow, cfg.ExternalIP))

			fmt.Printf("%s\n", color.Ize(color.Cyan, "--- Layer 0 ---"))

			nodeClient := node.GetClient("127.0.0.1", cfg.L0.Port.Public)
			nodeInfo, err := nodeClient.GetNodeInfo()
			if err != nil {
				fmt.Printf("%s\n", color.Ize(color.Red, "Node is down"))
			} else {
				l0Lb := lb.GetClient(cfg.L0.LoadBalancer)

				fmt.Printf("Node state: %s\n", getNodeState(nodeInfo.State))

				randomPeer, err := l0Lb.GetRandomReadyPeer()
				if err != nil {
					fmt.Printf("Version: %s\n", getVersion(nodeInfo.Version, "-"))
				} else {
					randomNodeClient := node.GetClient(randomPeer.Ip, randomPeer.PublicPort)
					peerInfo, err := randomNodeClient.GetNodeInfo()
					if err != nil {
						fmt.Printf("Version: %s\n", getVersion(nodeInfo.Version, "-"))
					} else {
						fmt.Printf("Version: %s\n", getVersion(nodeInfo.Version, peerInfo.Version))
					}
				}

				uptime := time.Unix(nodeInfo.Session / int64(time.Microsecond), 0)
				fmt.Printf("Uptime: %s (%s)\n", color.Ize(color.Yellow, fmt.Sprintf("%s", now.Sub(uptime))), color.Ize(color.Yellow, uptime.Format(time.RFC822)))

				info, err := l0Lb.FindPeerById(nodeInfo.Id)
				if err != nil {
					fmt.Printf("%s\n", color.Ize(color.Red, "Peer not found on load-balancer"))
				} else {
					fmt.Printf("Load-balancer: %s\n", getNodeState(info.State))
				}

			}

			fmt.Printf("%s\n", color.Ize(color.Cyan, "\n--- Layer 1 ---"))

			nodeClient = node.GetClient("127.0.0.1", cfg.L1.Port.Public)
			nodeInfo, err = nodeClient.GetNodeInfo()
			if err != nil {
				fmt.Printf("%s\n", color.Ize(color.Red, "Node is down"))
			} else {
				l1Lb := lb.GetClient(cfg.L1.LoadBalancer)

				fmt.Printf("Node state: %s\n", getNodeState(nodeInfo.State))

				randomPeer, err := l1Lb.GetRandomReadyPeer()
				if err != nil {
					fmt.Printf("Version: %s\n", getVersion(nodeInfo.Version, "-"))
				} else {
					randomNodeClient := node.GetClient(randomPeer.Ip, randomPeer.PublicPort)
					peerInfo, err := randomNodeClient.GetNodeInfo()
					if err != nil {
						fmt.Printf("Version: %s\n", getVersion(nodeInfo.Version, "-"))
					} else {
						fmt.Printf("Version: %s\n", getVersion(nodeInfo.Version, peerInfo.Version))
					}
				}

				uptime := time.Unix(nodeInfo.Session / int64(time.Microsecond), 0)
				fmt.Printf("Uptime: %s (%s)\n", color.Ize(color.Yellow, fmt.Sprintf("%s", now.Sub(uptime))), color.Ize(color.Yellow, uptime.Format(time.RFC822)))

				info, err := l1Lb.FindPeerById(nodeInfo.Id)
				if err != nil {
					fmt.Printf("%s\n", color.Ize(color.Red, "Peer not found on load-balancer"))
				} else {
					fmt.Printf("Load-balancer: %s\n", getNodeState(info.State))
				}

			}

			return nil
		},
	}
	Verbose bool
)

func getNodeState(state string) string {
	nodeState := node2.ParseString(state)
	if nodeState == node2.Offline || nodeState == node2.Leaving {
		return color.Ize(color.Red, state)
	}
	if nodeState == node2.Ready {
		return color.Ize(color.Green, state)
	}
	return color.Ize(color.Yellow, state)
}

func getVersion(nodeVersion string, clusterVersion string) string {
	if nodeVersion == clusterVersion {
		return color.Ize(color.Green, "v" + nodeVersion)
	}
	if clusterVersion == "-" {
		return fmt.Sprintf("node=%s, cluster=%s", color.Ize(color.Green, "v" + nodeVersion), color.Ize(color.Red, "unknown"))
	}
	return fmt.Sprintf("node=%s, cluster=%s", color.Ize(color.Red, "v" + nodeVersion), color.Ize(color.Green, "v" + clusterVersion))
}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return nil
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&Verbose, "verbose", false, "verbose output")
	_ = viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.SetConfigName("dag")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/dag")
	viper.AddConfigPath("$HOME/.dag")
	viper.AddConfigPath(".dag")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	cobra.CheckErr(err)
}
