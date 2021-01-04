package cmd

import (
	"crypto/tls"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/pborman/getopt/v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"
)

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

func ShowRepos(cfg Config) {
	fmt.Printf("%v\n", cfg)
}

func UsageMessage() {
	getopt.PrintUsage(os.Stderr)
	os.Exit(0)
}

func checkRepo(url string, exp int, wg *sync.WaitGroup, c chan GitCheckResult) {
	// Clone repository using the new client if the protocol is https://
	var err error
	var ret int
	var check bool = false

	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{URL: url})
	if err != nil {
		if err.Error() == "authentication required" {
			ret = PRIVATE
		} else {
			ret = ERROR
		}
	} else {
		_, err = r.Head()
		//fmt.Println(head.Hash())
		if err != nil {
			ret = ERROR
		} else {
			ret = PUBLIC
		}
	}
	if ret == exp {
		check = true
	}
	c <- NewGitChecker(url, exp, ret, err, check)
	wg.Done()
}
