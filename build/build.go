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
	"os/exec"
)

// get the version of go tools
func GoVersion() (string, error) {
	cmd := exec.Command("go", "version")
	return execCommand(cmd)
}

func execCommand(cmd *exec.Cmd) (string, error) {
	out, err := cmd.Output()
	if err == nil {
		return string(out[:len(out)]), nil
	} else {
		return err.Error(), err
	}
}
