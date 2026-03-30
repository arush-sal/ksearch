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
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	"github.com/spf13/cobra"
)

var (
	resName, namespace, kinds string
	cacheTTL                  time.Duration
	noCache                   bool
	version                   = "dev"

	currentContextNameFn    = currentContextName
	readCacheFn             = cache.Read
	writeCachedSectionsFn   = writeCachedSections
	writeCacheFn            = cache.Write
	getConfigOrDieFn        = func() *rest.Config { return config.GetConfigOrDie() }
	newClientsetForConfigFn = func(cfg *rest.Config) kubernetes.Interface {
		return kubernetes.NewForConfigOrDie(cfg)
	}
	discoverResourcesFn = util.Discover
	getterFn            = util.Getter
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Short:   "run ksearch --help to get the usage",
	Version: version,
	Long:    `ksearch is a command line tool to search for a given pattern in a Kubernetes cluster and will print all of the available resources in a cluster if none is provided`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runRoot(cmd, args); err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}
	},
}

func runRoot(cmd *cobra.Command, args []string) error {
	ttl, err := resolvedCacheTTL(cmd)
	if err != nil {
		return err
	}

	currentContext, err := currentContextNameFn()
	if err != nil {
		return err
	}

	key := cache.KeyFor(currentContext, namespace, kinds, resName)

	if !noCache {
		entry, err := readCacheFn(key, ttl)
		if err != nil {
			return err
		}

		if entry != nil {
			return writeCachedSectionsFn(entry.Sections)
		}
	}

	cfg := getConfigOrDieFn()
	clientset := newClientsetForConfigFn(cfg)

	resources, err := discoverResourcesFn(clientset.Discovery(), kinds)
	if err != nil {
		return err
	}

	getter := make(chan util.FetchResult)
	go getterFn(namespace, clientset, cfg, resources, getter)

	results := make([]cache.SectionEntry, len(resources))
	var wg sync.WaitGroup
	index := 0
	for fetched := range getter {
		resultIndex := index
		index++

		wg.Add(1)
		go func(idx int, fetched util.FetchResult) {
			defer wg.Done()

			var buffer bytes.Buffer
			if fetched.Resource != nil {
				printers.Printer(&buffer, fetched.Resource, resName)
			}
			results[idx] = cache.SectionEntry{
				Kind:   fetched.Kind,
				Output: buffer.String(),
			}
		}(resultIndex, fetched)
	}

	wg.Wait()

	if err := writeCacheFn(key, cache.CacheMeta{
		Context:    currentContext,
		Namespace:  namespace,
		Kinds:      kinds,
		Pattern:    resName,
		TTLSeconds: int(ttl.Seconds()),
	}, results); err != nil {
		return err
	}

	for _, result := range results {
		if len(result.Output) == 0 {
			continue
		}

		if _, err := os.Stdout.Write([]byte(result.Output)); err != nil {
			return err
		}
	}

	return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func SetVersion(buildVersion string) {
	version = buildVersion
	rootCmd.Version = version
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Use = pluginName()
	rootCmd.PersistentFlags().StringVarP(&resName, "pattern", "p", "", "pattern you want to search for")
	rootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "namespace you want to search in")
	rootCmd.PersistentFlags().StringVarP(&kinds, "kinds", "k", "", "comma separated list of all the kinds that you want to include")
	rootCmd.PersistentFlags().DurationVar(&cacheTTL, "cache-ttl", time.Minute, "duration before cache expires")
	rootCmd.PersistentFlags().BoolVar(&noCache, "no-cache", false, "skip cache for this invocation")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
}

func pluginName() string {
	base := os.Args[0]
	if lastSeparator := strings.LastIndexAny(base, `/\`); lastSeparator >= 0 {
		base = base[lastSeparator+1:]
	}

	base = strings.TrimSuffix(base, ".exe")
	if strings.HasPrefix(base, "kubectl-") {
		return "kubectl " + strings.TrimPrefix(base, "kubectl-")
	}

	return base
}

func resolvedCacheTTL(cmd *cobra.Command) (time.Duration, error) {
	flag := cmd.Flags().Lookup("cache-ttl")
	if flag != nil && flag.Changed {
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
