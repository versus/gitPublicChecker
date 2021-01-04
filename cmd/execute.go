package cmd

import (
	"fmt"
	"github.com/go-git/go-git/v5/plumbing/transport/client"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/pborman/getopt/v2"
	"sync"
)

func Execute() {

	getopt.Parse()
	args := getopt.Args()
	if len(args) > 1 || helpFlag {
		UsageMessage()
	}

	ShowRepos(config)
	fmt.Println("tail:", args, excludeFlag)
	client.InstallProtocol("https", githttp.NewClient(gitClient))

	var wg sync.WaitGroup
	channel := make(chan GitCheckResult)
	for i := range config.PublicRepos.Url {
		wg.Add(1)
		go checkRepo(config.PublicRepos.Url[i], PUBLIC, &wg, channel)
	}
	for i := range config.PrivateRepos.Url {
		wg.Add(1)
		go checkRepo(config.PrivateRepos.Url[i], PRIVATE, &wg, channel)
	}

	repos := len(config.PublicRepos.Url) + len(config.PrivateRepos.Url)

	for i := 0; i < repos; i++ {
		message := <-channel
		fmt.Println(message.Url, "is expected ", message.Check)
	}

	wg.Wait()
	close(channel)
}
