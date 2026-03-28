package cmd

import (
	"bytes"
	"encoding/base64"
	"fmt"

	// Load all known auth plugins
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"os"
	"strings"
	"sync"
	"time"

	"github.com/arush-sal/ksearch/pkg/cache"
	"github.com/arush-sal/ksearch/pkg/printers"
	"github.com/arush-sal/ksearch/pkg/util"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	"github.com/spf13/cobra"
)

var (
	resName, namespace, kinds string
	cacheTTL                  time.Duration
	noCache                   bool
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
		ttl, err := resolvedCacheTTL(cmd)
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}

		currentContext, err := currentContextName()
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}

		getter := make(chan interface{})
		resourceOrder := effectiveResources(kinds)
		key := cache.KeyFor(currentContext, namespace, kinds, resName)

		if !noCache {
			entry, err := cache.Read(key, ttl)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}

			if entry != nil {
				if err := writeCachedSections(entry.Sections); err != nil {
					cmd.PrintErrln(err)
					os.Exit(1)
				}
				return
			}
		}

		cfg := config.GetConfigOrDie()
		clientset := kubernetes.NewForConfigOrDie(cfg)

		go util.Getter(namespace, clientset, kinds, getter)

		results := make([]cache.SectionEntry, len(resourceOrder))
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
				results[idx] = cache.SectionEntry{
					Kind:   resourceOrder[idx],
					Output: buffer.String(),
				}
			}(resultIndex, resource)
		}

		wg.Wait()

		if err := cache.Write(key, cache.CacheMeta{
			Context:    currentContext,
			Namespace:  namespace,
			Kinds:      kinds,
			Pattern:    resName,
			TTLSeconds: int(ttl.Seconds()),
		}, results); err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}

		for _, result := range results {
			if len(result.Output) == 0 {
				continue
			}

			if _, err := os.Stdout.Write([]byte(result.Output)); err != nil {
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
	rootCmd.PersistentFlags().DurationVar(&cacheTTL, "cache-ttl", time.Minute, "duration before cache expires")
	rootCmd.PersistentFlags().BoolVar(&noCache, "no-cache", false, "skip cache for this invocation")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
}

func resolvedCacheTTL(cmd *cobra.Command) (time.Duration, error) {
	if cmd.Flags().Lookup("cache-ttl").Changed {
		return cacheTTL, nil
	}

	value := os.Getenv("KSEARCH_CACHE_TTL")
	if value == "" {
		return cacheTTL, nil
	}

	ttl, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("parse KSEARCH_CACHE_TTL: %w", err)
	}

	return ttl, nil
}

func currentContextName() (string, error) {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	cfg, err := rules.Load()
	if err != nil {
		return "", err
	}

	return cfg.CurrentContext, nil
}

func writeCachedSections(sections []cache.SectionEntry) error {
	for _, section := range sections {
		if section.Output == "" {
			continue
		}

		decoded, err := base64.StdEncoding.DecodeString(section.Output)
		if err != nil {
			return fmt.Errorf("decode cached section %s: %w", section.Kind, err)
		}

		if _, err := os.Stdout.Write(decoded); err != nil {
			return err
		}
	}

	return nil
}
