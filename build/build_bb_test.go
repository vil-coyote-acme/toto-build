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
	"github.com/vil-coyote-acme/toto-build-common/testtools"
	"strings"
	"testing"
	"toto-build-agent/build"
)

// test the printing of go tools versions
func Test_Should_Get_Go_Tools_Versions(t *testing.T) {
	c := build.GoVersion()
	out := testtools.ConsumeStringChan(c)
	t.Logf("Test the go version command. Output : %s\n\r", out)
	assert.True(t, strings.Contains(out, "go version"))
}

// test the build function
func Test_Should_Build_Test_Sources(t *testing.T) {
	c := build.BuildPackage("toto-build-agent/testapp")
	out := testtools.ConsumeStringChan(c)
	t.Logf("Test the go build command with succes. Output : %s\n\r", out)
	assert.Contains(t, strings.TrimSpace(out), "toto-build-agent/testapp")
}

func Test_Should_Test_Sources(t *testing.T) {
	c := build.TestPackage("toto-build-agent/testapp")
	out := testtools.ConsumeStringChan(c)
	t.Logf("Test the go test command with succes. Output : %s\n", out)
}
