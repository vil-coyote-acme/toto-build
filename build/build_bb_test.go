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
package build_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"toto-build-agent/build"
	"github.com/vil-coyote-acme/toto-build-common/message"
	"github.com/vil-coyote-acme/toto-build-common/testtools"
)

// test the printing of go tools versions
func Test_Should_Get_Go_Tools_Versions(t *testing.T) {
	// given
	reportChan := make(chan message.Report, 1)
	defer close(reportChan)
	build.GoVersion(int64(1), reportChan)
	// when
	msg := <-reportChan
	end := <-reportChan
	// then
	t.Logf("Test the go version command. Output : %s\n\r", msg.Logs)
	assert.Contains(t, testtools.FromSliceToString(msg.Logs), "go version")
	assert.Equal(t, msg.Status, message.WORKING)
	assert.Equal(t, end.Status, message.SUCCESS)
}

// test the build function
func Test_Should_Build_Test_Sources(t *testing.T) {
	// given
	reportChan := make(chan message.Report, 1)
	defer close(reportChan)
	build.BuildPackage("toto-build-agent/testapp", int64(1), reportChan)
	// when
	msg := <-reportChan
	end := <-reportChan
	t.Logf("Test the go build command with succes. Output : %s\n\r", msg.Logs)
	assert.Contains(t, testtools.FromSliceToString(msg.Logs), "toto-build-agent/testapp")
	assert.Equal(t, msg.Status, message.WORKING)
	assert.Equal(t, end.Status, message.SUCCESS)
}

func Test_Should_Test_Sources(t *testing.T) {
	// given
	reportChan := make(chan message.Report, 1)
	defer close(reportChan)
	build.TestPackage("toto-build-agent/testapp", int64(1), reportChan)
	// when
	msg := <-reportChan
	end := <-reportChan
	// then
	t.Logf("Test the go test command with succes. Output : %s\n", msg.Logs)
	assert.Contains(t, testtools.FromSliceToString(msg.Logs), "toto-build-agent/testapp")
	assert.Equal(t, msg.Status, message.WORKING)
	assert.Equal(t, end.Status, message.SUCCESS)
}

