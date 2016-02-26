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
	"io"
	"os/exec"
	"strings"
	"log"
	"github.com/vil-coyote-acme/toto-build-common/message"
)

// will lauch job on incoming toWork
func ExecuteJob(toWorkChan chan message.ToWork, reportChan chan message.Report) {
	go func() {
		for toWork := range toWorkChan {
			log.Printf("receive one job : %s", toWork)
			var logsChan chan string
			switch toWork.Cmd {
			case message.PACKAGE :
				logsChan = BuildPackage(toWork.Package)
			case message.TEST:
				logsChan = TestPackage(toWork.Package)
			case message.HELLO:
				reportChan <- message.Report{toWork.JobId, message.WORKING, []string{"Hello"}}
			default:
			// todo handle this case
			}
			if logsChan != nil {
				go listenForLogs(logsChan, reportChan, toWork)
			}
		}
	}()
}

// listen for log from a job and push it to the report chan
func listenForLogs(logsChan chan string, reportChan chan message.Report, toWork message.ToWork) {
	for log := range logsChan {
		// todo handle buffered logs
		// todo handle job status
		reportChan <- message.Report{toWork.JobId, message.WORKING, []string{log}}
	}

}

// get the version of go tools
func GoVersion() chan string {
	cmd := exec.Command("go", "version")
	return execCommand(cmd)
}

// call the build command
func BuildPackage(pkg string) chan string {
	// todo next here: support many options !
	cmd := exec.Command("go", "build", "-v", pkg)
	return execCommand(cmd)
}

func TestPackage(pkg string) chan string {
	cmd := exec.Command("go", "test", "-cover", pkg)
	return execCommand(cmd)
}

// execute one command
func execCommand(cmd *exec.Cmd) chan string {
	log.Printf("Start executing one command : %s", cmd)
	c := make(chan string, 50)
	go func() {
		defer close(c)
		stdout, errPipe1 := cmd.StdoutPipe()
		stderr, errPipe2 := cmd.StderrPipe()
		errCmd := cmd.Start()
		hasErr, errMes := hasError(errPipe1, errPipe2, errCmd)
		if hasErr {
			c <- strings.Join(errMes, "\n\r")
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
	}()
	return c
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
