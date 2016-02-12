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
	"io"
	"bufio"
)

// get the version of go tools
func GoVersion() (chan string) {
	cmd := exec.Command("go", "version")
	return execCommand(cmd)
}

// call the build command
func BuildPackage(pkg string) (chan string) {
	// todo next here: support many options !
	cmd := exec.Command("go", "build", "-v", pkg)
	return execCommand(cmd)
}

func TestPackage(pkg string) (chan string) {
	cmd := exec.Command("go", "test", pkg)
	return execCommand(cmd)
}

// execute one command
func execCommand(cmd *exec.Cmd) (chan string) {
	c := make(chan string, 10)
	go func() {
		defer close(c)
		stdout, errPipe1 := cmd.StdoutPipe()
		stderr, errPipe2 := cmd.StderrPipe()
		errCmd := cmd.Start()
		if errCmd != nil {
			c <- errCmd.Error()
		} else {
			// todo check this errors in unit test. With mock ?
			if errPipe1 != nil {
				c <- errPipe1.Error()
			} else if errPipe2 != nil {
				c <- errPipe2.Error()
			} else {
				multi := io.MultiReader(stdout, stderr)
				in := bufio.NewScanner(multi)
				for in.Scan() {
					c <- in.Text()
				}
				if in.Err() != nil {
					c <- in.Err().Error()
				}
			}
		}
	}()
	return c
}
