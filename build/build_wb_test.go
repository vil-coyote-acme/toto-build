/*
Toto-build, the stupid Go continuous build server.

Toto-build is free software; you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation; either version 3 of the License, or
(at your option) any later version.

Toto-build is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program; if not, write to the Free Software Foundation,
Inc., 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301  USA
*/
package build

import (
	"github.com/stretchr/testify/assert"
	"github.com/vil-coyote-acme/toto-build-common/message"
	"github.com/vil-coyote-acme/toto-build-common/testtools"
	"os/exec"
	"strings"
	"testing"
)

// test the command execution
func Test_execCommand_should_failed_for_non_existing_command(t *testing.T) {
	c := execCommand(exec.Command("toto", "isHappy"))
	str := testtools.ConsumeStringChan(c)
	t.Logf("test the exec command failure. Output : %s", str)
	assert.True(t, strings.Contains(str, "executable file not found in $PATH"))
}

func Test_hasError_Should_Return_False(t *testing.T) {
	res, mes := hasError(nil, nil, nil)
	assert.False(t, res)
	assert.Nil(t, mes)
}

func Test_hasError_Should_Return_True(t *testing.T) {
	res, mes := hasError(nil, nil, testtools.NewTestErr("my error"))
	assert.True(t, res)
	assert.Equal(t, "my error", mes[0])
}

func Test_listenForLogs_without_agregate(t *testing.T) {
	// given
	toWork := message.ToWork{int64(1), message.TEST, "myPkg"}
	logsChan := make(chan string, 2)
	reportChan := make(chan message.Report)
	defer close(logsChan)
	defer close(reportChan)
	go listenForLogs(logsChan, reportChan, toWork)
	// when
	logsChan <- "toto"
	// then
	assert.Equal(t, message.Report{toWork.JobId, message.WORKING, []string{"toto"}}, <-reportChan)

}

func Test_listenForLogs_with_agregate(t *testing.T) {
	// given
	toWork := message.ToWork{int64(1), message.TEST, "myPkg"}
	logsChan := make(chan string, 2)
	reportChan := make(chan message.Report)
	defer close(logsChan)
	defer close(reportChan)
	go listenForLogs(logsChan, reportChan, toWork)
	// when
	for i := 0; i < 2; i++ { // test twice to check the reset of the internal buffer
		logsChan <- "toto"
		logsChan <- "titi"
		// then
		assert.Equal(t, message.Report{toWork.JobId, message.WORKING, []string{"toto", "titi"}}, <-reportChan)
	}
}
