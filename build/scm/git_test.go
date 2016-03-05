package scm_test

import (
	"testing"
	"os"
)

func Test_Fetch_successFull_When_No_Local_Repo(t *testing.T) {
	// given
	defer os.RemoveAll(os.Getenv("GOPATH") + "/src/")
	//mes :=
}
