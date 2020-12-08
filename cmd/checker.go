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

const ConfigFile = "config.yaml"

type Config struct {
	Repos    []string `yaml:"urls,flow"`
	Opsgenie struct {
		Api string `yaml:"api"`
	} `yaml:"opsgenie"`
}

var (
	config    Config
	gitClient *http.Client
)

func init() {
	configFile, err := ioutil.ReadFile(ConfigFile)
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
	helpFlag := getopt.BoolLong("help", 'h', "Displays help message")
	excludeFlag := getopt.StringLong("exclude", 'e', "", "for exclude a file from check")
	quietFlag := getopt.BoolLong("quiet", 'q', "see only warnings without additional information")

	getopt.Parse()
	args := getopt.Args()
	if len(args) > 1 || *helpFlag {
		UsageMessage()
	}

	ShowRepos(config)
	fmt.Println(*excludeFlag, "tail:", args, *quietFlag)
	CheckRepo(config.Repos[0])
}

func ShowRepos(cfg Config) {
	fmt.Printf("Result: %v\n", cfg)
}

func UsageMessage() {
	getopt.PrintUsage(os.Stderr)
	os.Exit(0)
}

func CheckRepo(url string) {
	client.InstallProtocol("https", githttp.NewClient(gitClient))
	// Clone repository using the new client if the protocol is https://
	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{URL: url})
	if err != nil {
		if err.Error() == "authentication required" {
			fmt.Printf("Repo still private: %s\n", err)
		} else {
			fmt.Printf("ERROR: %s\n", err)
		}
		return
	}
	head, err := r.Head()
	fmt.Println(head.Hash())
}