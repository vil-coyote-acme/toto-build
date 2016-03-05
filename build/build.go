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
	"bufio"
	"fmt"
	"github.com/vil-coyote-acme/toto-build-common/message"
	"io"
	"os/exec"
	"syscall"
	"github.com/vil-coyote-acme/toto-build-common/logs"
)

// todo limit the number of goroutines !

var (
	logger = logs.NewLogger("[BUILD-ENGINE] : ", logs.NewConsoleAppender(logs.INFO))
)

// will lauch job on incoming toWork
func ExecuteJob(toWorkChan chan message.ToWork, reportChan chan message.Report) {
	// one goroutine for launching jobs. May (must ?) be merge with routine in handler ?
	go func() {
		for toWork := range toWorkChan {
			logger.Infof("receive one job : %s", toWork)
			switch toWork.Cmd {
			case message.PACKAGE:
				BuildPackage(toWork, reportChan)
			case message.TEST:
				TestPackage(toWork, reportChan)
			case message.HELLO:
				reportChan <- message.Report{toWork.JobId, message.SUCCESS, []string{"Hello"}}
			default:
			// todo handle this case
			}
		}
	}()
}

// call the build command
func BuildPackage(mes message.ToWork, reportChan chan message.Report) {
	// todo next here: support many options !
	cmd := exec.Command("go", "build", "-v", "-a", mes.Package)
	execCommand(cmd, mes, reportChan)
}

func TestPackage(mes message.ToWork, reportChan chan message.Report) {
	cmd := exec.Command("go", "test", "-cover", mes.Package)
	execCommand(cmd, mes, reportChan)
}

// execute one command
func execCommand(cmd *exec.Cmd, mes message.ToWork, reportChan chan message.Report) {
	logger.Infof("Start executing one command : %v", cmd)
	go func() {
		stdout, errPipe1 := cmd.StdoutPipe()
		stderr, errPipe2 := cmd.StderrPipe()
		errCmd := cmd.Start()
		hasErr, errMes := hasError(errPipe1, errPipe2, errCmd)
		if hasErr {
			reportChan <- message.Report{mes.JobId, message.FAILED, errMes}
		} else {
			multi := io.MultiReader(stdout, stderr)
			in := bufio.NewScanner(multi)
			// todo see if a more efficient way is possible
			buf := []string{}
			for in.Scan() {
				s := in.Text()
				buf = append(buf, s)
				if len(buf) > 4 {
					buf = consumeBuffer(buf, mes.JobId, reportChan)
				}
			}
			// finish the rest of the buffer
			consumeBuffer(buf, mes.JobId, reportChan)
			if inErr := in.Err(); inErr != nil {
				reportChan <- message.Report{mes.JobId, message.FAILED, []string{inErr.Error()}}
			} else {
				pushEndReport(cmd, mes.JobId, reportChan)
			}
		}
	}()
}

// Detect error and return mes error
func hasError(errors ...error) (res bool, errMess []string) {
	for _, err := range errors {
		if err != nil {
			errMess = append(errMess, err.Error())
			res = true
		}
	}
	return
}

// if not empty, consume all the buffer and send a report
func consumeBuffer(buf []string, jobId int64, reportChan chan message.Report) []string {
	if len(buf) > 0 {
		reportChan <- message.Report{jobId, message.WORKING, buf}
		return []string{}
	}
	return buf
}

// must be called at the end of the command, when there is nothing more to
// read from stdout and stderr
// will send a success or a failure report message through the given chanel
func pushEndReport(cmd *exec.Cmd, jobId int64, reportChan chan message.Report) {
	if err := cmd.Wait(); err != nil {
		// we have a failure for the command.
		// we will now try to get the exit code
		if exiterr, ok := err.(*exec.ExitError); ok {
			// the error from cmd.Wait is an *exec.ExitError. It means that the
			// exit status is != 0

			/*
			 * here is a interesting comments getting from stackoverflow. I let it
			 *
			 */
			// This works on both Unix and Windows. Although package
			// syscall is generally platform dependent, WaitStatus is
			// defined for both Unix and Windows and in both cases has
			// an ExitStatus() method with the same signature.
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				errMesg := fmt.Sprintf("Exit Status: %d", status.ExitStatus())
				reportChan <- message.Report{jobId, message.FAILED, []string{errMesg}}
			}
		} else {
			errMesg := fmt.Sprintf("Unknow error : %v", err)
			reportChan <- message.Report{jobId, message.FAILED, []string{errMesg}}
		}
	} else {
		// the command is in success state.
		reportChan <- message.Report{jobId, message.SUCCESS, []string{}}
	}
}
