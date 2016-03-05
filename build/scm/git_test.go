package scm_test

import (
	"testing"
	"os"
	"github.com/vil-coyote-acme/toto-build-common/message"
	"toto-build-agent/build/scm"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"log"
)

func Test_Fetch_successFull_When_No_Local_Repo(t *testing.T) {
	// given
	goPath := setGoPath("toto-tmp/")
	defer os.RemoveAll(goPath)
	mess := message.ToWork{int64(1), message.TEST, "myPkg", "go1.6", "https://github.com/vil-coyote-acme/toto-example.git"}
	// when
	err := scm.Fetch(mess)
	// then
	assert.Nil(t, err)
	f, errf := os.Stat(goPath + "src/example.go")
	assert.Nil(t, errf)
	assert.True(t, !f.IsDir())
}

func Test_Fetch_Failed_Bad_GOPATH(t *testing.T) {
	// given
	os.Setenv("GOPATH", "")
	mess := message.ToWork{int64(1), message.TEST, "myPkg", "go1.6", "https://github.com/vil-coyote-acme/toto-example.git"}
	// when
	err := scm.Fetch(mess)
	// then
	assert.NotNil(t, err)
	assert.True(t, os.IsNotExist(err))
}

func setGoPath(root string) string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Panic(err)
	}
	absDir := dir + "/" + root;
	os.Setenv("GOPATH", absDir)
	os.MkdirAll(absDir + "/src/", os.ModePerm)
	os.MkdirAll(absDir + "/pkg/", os.ModePerm)
	os.MkdirAll(absDir + "/bin/", os.ModePerm)
	_, errf := os.Stat(absDir + "/src/")
	if errf != nil {
		log.Panic(errf)
	}
	return absDir
}
