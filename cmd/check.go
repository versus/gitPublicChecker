package cmd

const (
	ERROR = iota
	PRIVATE
	PUBLIC
)

type GitCheckResult struct {
	Result   int
	Url      string
	Expected int
	Error    error
	Check    bool
}

func NewGitChecker(url string, expected int, res int, err error, check bool) GitCheckResult {
	return GitCheckResult{
		Result:   res,
		Error:    err,
		Url:      url,
		Expected: expected,
		Check:    check,
	}
}
