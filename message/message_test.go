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
package message_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"toto-build/message"
)

func Test_ToWork_String(t *testing.T) {
	assert.Equal(t, "JobId: 1, Cmd: 0, Package: myPackage ",
		message.ToWork{int64(1), message.PACKAGE, "myPackage"}.String())
}

func Test_Error(t *testing.T) {
	assert.Equal(t, "myError", message.Error{"myError"}.Error())
}
