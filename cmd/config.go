package cmd

type Config struct {
	PublicRepos struct {
		Url []string `yaml:"urls,flow"`
	} `yaml:"public"`
	PrivateRepos struct {
		Url []string `yaml:"urls,flow"`
	} `yaml:"private"`
	Opsgenie struct {
		Api string `yaml:"api"`
	} `yaml:"opsgenie" binding:"required"`
}
