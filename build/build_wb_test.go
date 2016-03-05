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
	"testing"
)

// test the command execution
func Test_execCommand_should_failed_for_non_existing_command(t *testing.T) {
	// given
	reportChan := make(chan message.Report, 1)
	defer close(reportChan)
	mes := message.ToWork{int64(1), message.PACKAGE, "toto", "go1.6", "https://github.com/vil-coyote-acme/toto-example.git"}
	// when
	execCommand(exec.Command("toto", "isHappy"), mes, reportChan)
	out := <-reportChan
	// then
	t.Logf("test the exec command failure. Output : %s, %d", out.Logs, len(out.Logs))
	assert.Contains(t, testtools.FromSliceToString(out.Logs), "executable file not found in $PATH")
	assert.Equal(t, out.Status, message.FAILED)
}

func Test_hasError_Should_Return_False(t *testing.T) {
	// when
	res, mes := hasError(nil, nil, nil)
	// then
	assert.False(t, res)
	assert.Nil(t, mes)
}

func Test_hasError_Should_Return_True(t *testing.T) {
	// when
	res, mes := hasError(nil, nil, testtools.NewTestErr("my error"))
	// then
	assert.True(t, res)
	assert.Equal(t, "my error", mes[0])
}

func Test_consumeBuffer_with_empty_buf(t *testing.T) {
	// given
	buf := []string{}
	reportChan := make(chan message.Report)
	defer close(reportChan)
	// when
	cleanBuf := consumeBuffer(buf, int64(1), reportChan)
	// then
	assert.Equal(t, buf, cleanBuf)
	assert.True(t, len(cleanBuf) == 0)
}

func Test_consumeBuffer_with_non_empty_buf(t *testing.T) {
	// given
	buf := []string{"titi, toto"}
	reportChan := make(chan message.Report, 2)
	defer close(reportChan)
	// when
	cleanBuf := consumeBuffer(buf, int64(1), reportChan)
	rep := <-reportChan
	// then
	assert.True(t, len(cleanBuf) == 0)
	assert.Equal(t, buf, rep.Logs)
	assert.Equal(t, message.WORKING, rep.Status)
	assert.Equal(t, int64(1), rep.JobId)
}
