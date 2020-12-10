package cmd

import (
	"crypto/tls"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/client"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/pborman/getopt/v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type Config struct {
	PublicRepos struct {
		Url []string `yaml:"urls,flow"`
	} `yaml:"public"`
	PrivateRepos struct {
		Url []string `yaml:"urls,flow"`
	} `yaml:"private"`
	Opsgenie struct {
		Api string `yaml:"api"`
	} `yaml:"opsgenie"`
}

var (
	config      Config
	gitClient   *http.Client
	timeout     = time.Second * 15
	verboseFlag bool
	configFile  = "./config.yaml"
	helpFlag    bool
	quietFlag   bool
	excludeFlag string
)

func init() {

	getopt.Flag(&verboseFlag, 'v', "be verbose").SetOptional()
	getopt.FlagLong(&configFile, "config", 'f', "path to config file").SetOptional()
	getopt.FlagLong(&timeout, "timeout", 't', "git client connect timeout ").SetOptional()
	getopt.FlagLong(&helpFlag, "help", 'h', "Displays help message").SetOptional()
	getopt.FlagLong(&excludeFlag, "exclude", 'e', "", "for exclude a repository from check, separated by coma").SetOptional()
	getopt.FlagLong(&quietFlag, "quiet", 'q', "see only warnings without additional information").SetOptional()

	configFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Printf("Error reading config file: %s\n", err)
		os.Exit(1)
	}

	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		fmt.Printf("Error parsing config file: %s\n", err)
		os.Exit(1)
	}

	gitClient = &http.Client{
		Transport: &http.Transport{ // accept any certificate (might be useful for testing)
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 15 * time.Second, // 15 second timeout
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // don't follow redirect
		},
	}

}

func Execute() {

	getopt.Parse()
	args := getopt.Args()
	if len(args) > 1 || helpFlag {
		UsageMessage()
	}

	ShowRepos(config)
	fmt.Println("tail:", args, excludeFlag)
	for i := range config.PrivateRepos.Url {
		CheckRepo(config.PrivateRepos.Url[i], true)
	}
	for i := range config.PrivateRepos.Url {
		CheckRepo(config.PublicRepos.Url[i], false)
	}
}

func ShowRepos(cfg Config) {
	fmt.Printf("Result: %v\n", cfg)
}

func UsageMessage() {
	getopt.PrintUsage(os.Stderr)
	os.Exit(0)
}

func CheckRepo(url string, condition bool) (int, error) {
	client.InstallProtocol("https", githttp.NewClient(gitClient))
	// Clone repository using the new client if the protocol is https://
	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{URL: url})
	if err != nil {
		if err.Error() == "authentication required" {
			fmt.Printf("Repo still private: %s\n", err)
			if condition {
				return 1, nil
			}
			return 0, nil
		} else {
			fmt.Printf("ERROR: %s\n", err)
			return -1, err
		}
	}
	head, err := r.Head()
	fmt.Println(head.Hash())
	if err != nil {
		return -1, err
	}
	if condition {
		return 1, nil
	}
	return 0, nil
}
