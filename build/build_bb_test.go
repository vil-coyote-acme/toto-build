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
	"github.com/vil-coyote-acme/toto-build-common/message"
	"github.com/vil-coyote-acme/toto-build-common/testtools"
	"testing"
	"toto-build-agent/build"
)

// test the build function
func Test_Should_Build_Test_Sources(t *testing.T) {
	// given
	reportChan := make(chan message.Report, 1)
	defer close(reportChan)
	mes := message.ToWork{int64(1), message.PACKAGE, "toto-build-agent/testapp", "go1.6", "https://github.com/vil-coyote-acme/toto-example.git"}
	// when
	build.BuildPackage(mes, reportChan)
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
	mes := message.ToWork{int64(1), message.PACKAGE, "toto-build-agent/testapp", "go1.6", "https://github.com/vil-coyote-acme/toto-example.git"}
	// when
	build.TestPackage(mes, reportChan)
	msg := <-reportChan
	end := <-reportChan
	// then
	t.Logf("Test the go test command with succes. Output : %s\n", msg.Logs)
	assert.Contains(t, testtools.FromSliceToString(msg.Logs), "toto-build-agent/testapp")
	assert.Equal(t, msg.Status, message.WORKING)
	assert.Equal(t, end.Status, message.SUCCESS)
}

func Test_TactPackage_should_failed_for_unknown_package(t *testing.T) {
	// given
	reportChan := make(chan message.Report, 1)
	defer close(reportChan)
	mes := message.ToWork{int64(1), message.PACKAGE, "plop/", "go1.6", "https://github.com/vil-coyote-acme/toto-example.git"}
	// when
	build.TestPackage(mes, reportChan)
	msg := <-reportChan
	end := <-reportChan
	// then
	t.Logf("test the exec command failure. Output : %s, %d", msg.Logs, len(msg.Logs))
	assert.Contains(t, testtools.FromSliceToString(msg.Logs), "can't load package: package plop")
	assert.Equal(t, msg.Status, message.WORKING)
	assert.Equal(t, end.Status, message.FAILED)
}

func Test_BuildPackage_should_failed_for_unknown_package(t *testing.T) {
	// given
	reportChan := make(chan message.Report, 1)
	defer close(reportChan)
	mes := message.ToWork{int64(1), message.PACKAGE, "plop/", "go1.6", "https://github.com/vil-coyote-acme/toto-example.git"}
	// when
	build.BuildPackage(mes, reportChan)
	msg := <-reportChan
	end := <-reportChan
	// then
	t.Logf("test the exec command failure. Output : %s, %d", msg.Logs, len(msg.Logs))
	assert.Contains(t, testtools.FromSliceToString(msg.Logs), "can't load package: package plop")
	assert.Equal(t, msg.Status, message.WORKING)
	assert.Equal(t, end.Status, message.FAILED)
}

// test the build function
func Test_ExecuteJob_Should_Build_Test_Sources(t *testing.T) {
	// given
	toWorkChan := make(chan message.ToWork)
	reportChan := make(chan message.Report, 1)
	defer close(reportChan)
	defer close(toWorkChan)
	// when
	build.ExecuteJob(toWorkChan, reportChan)
	toWorkChan <- message.ToWork{int64(1), message.PACKAGE, "toto-build-agent/testapp", "go1.6", "https://github.com/vil-coyote-acme/toto-example.git"}
	msg := <-reportChan
	end := <-reportChan
	t.Logf("Test the go build command with succes. Output : %s\n\r", msg.Logs)
	assert.Contains(t, testtools.FromSliceToString(msg.Logs), "toto-build-agent/testapp")
	assert.Equal(t, msg.Status, message.WORKING)
	assert.Equal(t, end.Status, message.SUCCESS)
}

func Test_ExecuteJob_Should_Test_Sources(t *testing.T) {
	// given
	toWorkChan := make(chan message.ToWork)
	reportChan := make(chan message.Report, 1)
	defer close(reportChan)
	defer close(toWorkChan)
	// when
	build.ExecuteJob(toWorkChan, reportChan)
	toWorkChan <- message.ToWork{int64(1), message.TEST, "toto-build-agent/testapp", "go1.6", "https://github.com/vil-coyote-acme/toto-example.git"}
	msg := <-reportChan
	end := <-reportChan
	t.Logf("Test the go build command with succes. Output : %s\n\r", msg.Logs)
	assert.Contains(t, testtools.FromSliceToString(msg.Logs), "toto-build-agent/testapp")
	assert.Equal(t, msg.Status, message.WORKING)
	assert.Equal(t, end.Status, message.SUCCESS)
}

func Test_ExecuteJob_Should_Reply_To_Hello(t *testing.T) {
	// given
	toWorkChan := make(chan message.ToWork)
	reportChan := make(chan message.Report, 1)
	defer close(reportChan)
	defer close(toWorkChan)
	// when
	build.ExecuteJob(toWorkChan, reportChan)
	toWorkChan <- message.ToWork{int64(1), message.HELLO, "", "go1.6", "https://github.com/vil-coyote-acme/toto-example.git"}
	msg := <-reportChan
	t.Logf("Test the go build command with succes. Output : %s\n\r", msg.Logs)
	assert.Contains(t, testtools.FromSliceToString(msg.Logs), "Hello")
	assert.Equal(t, msg.Status, message.SUCCESS)
}
