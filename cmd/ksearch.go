package cmd

import (
	"bytes"

	// Load all known auth plugins
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"os"
	"strings"
	"sync"

	"github.com/arush-sal/ksearch/pkg/printers"
	"github.com/arush-sal/ksearch/pkg/util"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	"github.com/spf13/cobra"
)

var (
	resName, namespace, kinds string
)

var defaultResources = []string{
	"Pods",
	"ConfigMaps",
	"Endpoints",
	"Events",
	"LimitRanges",
	"Namespaces",
	"PersistentVolumes",
	"PersistentVolumeClaims",
	"PodTemplates",
	"ResourceQuotas",
	"Secrets",
	"Services",
	"ServiceAccounts",
	"DaemonSets",
	"Deployments",
	"ReplicaSets",
	"StatefulSets",
}

func effectiveResources(kinds string) []string {
	if kinds == "" {
		return append([]string(nil), defaultResources...)
	}

	return strings.Split(kinds, ",")
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "ksearch",
	Short:   "run ksearch --help to get the usage",
	Version: "v0.0.1",
	Long:    `ksearch is a command line tool to search for a given pattern in a Kubernetes cluster and will print all of the available resources in a cluster if none is provided`,
	Run: func(cmd *cobra.Command, args []string) {
		getter := make(chan interface{})
		effectiveResources := effectiveResources(kinds)

		cfg := config.GetConfigOrDie()
		clientset := kubernetes.NewForConfigOrDie(cfg)

		go util.Getter(namespace, clientset, kinds, getter)

		results := make([][]byte, len(effectiveResources))
		var wg sync.WaitGroup
		index := 0
		for resource := range getter {
			resultIndex := index
			index++

			wg.Add(1)
			go func(idx int, renderedResource interface{}) {
				defer wg.Done()

				var buffer bytes.Buffer
				printers.Printer(&buffer, renderedResource, resName)
				results[idx] = append([]byte(nil), buffer.Bytes()...)
			}(resultIndex, resource)
		}

		wg.Wait()

		for _, result := range results {
			if len(result) == 0 {
				continue
			}

			if _, err := os.Stdout.Write(result); err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&resName, "pattern", "p", "", "pattern you want to search for")
	rootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "namespace you want to search in")
	rootCmd.PersistentFlags().StringVarP(&kinds, "kinds", "k", "", "comma separated list of all the kinds that you want to include")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
}
