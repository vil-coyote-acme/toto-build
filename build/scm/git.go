package scm

import (
	"github.com/vil-coyote-acme/toto-build-common/message"
	"os"
	"github.com/libgit2/git2go"
	"github.com/vil-coyote-acme/toto-build-common/logs"
)

var (
	logger *logs.Logger = logs.NewLogger("[GIT-FETCHER] : ", logs.NewConsoleAppender(logs.INFO))
)

func Fetch(mes message.ToWork) (err error) {
	goPath := os.Getenv("GOPATH")
	_, errGoPath := os.Open(goPath + "/src")
	if errGoPath != nil {
		logger.Errorf("error while checking %s", goPath + "/src")
		err =errGoPath
		return
	}
	_, errClone := git.Clone(mes.RepoUrl, goPath + "/src", new(git.CloneOptions))
	if errClone != nil {
		logger.Errorf("error while trying to checkout %s on ", mes.RepoUrl, goPath + "/src")
		return errClone
	}
	logger.Infof("successfully clone %s into %s", mes.RepoUrl, goPath + "/src")
	return
}


