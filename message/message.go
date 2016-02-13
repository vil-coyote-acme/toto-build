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
package message

import "fmt"

// type of available commands
type Command int
const (
	PACKAGE Command = iota
	TEST
)

// job execution state
type State int
const (
	PENDING State = iota
	WORKING
	SUCCESS
	FAILED
)

// message to run new job
type ToWork struct {
	JobId int64
	Cmd Command
	Package string
	// todo string function
	// todo add repository and credentials
}

func (t ToWork) String() string {
	return fmt.Sprintf("JobId: %d, Cmd: %d, Package: %s ", t.JobId, t.Cmd, t.Package)
}

// message to report job execution
type Report struct {
	JobId int64
	Status State
	Logs []string
}

type Error struct {
	Mes string
}

func (e Error) Error() string {
	return e.Mes
}
