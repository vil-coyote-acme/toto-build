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
package messaging

import (
	"bytes"
	"encoding/json"
	"github.com/nsqio/go-nsq"
	"github.com/vil-coyote-acme/toto-build-common/message"
	"log"
)

// messaging api handler
type Handler struct {
	buffer chan message.ToWork
}

// Handle the incoming message. If inccorect or if buffer is full, return one error
func (h *Handler) HandleMessage(mes *nsq.Message) error {
	toWork, e := unmarshallToWork(mes)
	if e == nil {
		select {
		case h.buffer <- toWork:
		default:
			e = message.Error{"the buffer is full"}
		}
	}
	return e
}

// create one new Handler using the given buffer
func NewHandler(buffer chan message.ToWork) *Handler {
	h := new(Handler)
	h.buffer = buffer
	return h
}

// unmarshall the ToWork message
func unmarshallToWork(mes *nsq.Message) (message.ToWork, error) {
	var toWork message.ToWork
	err := json.NewDecoder(bytes.NewBuffer(mes.Body)).Decode(&toWork)
	if err != nil {
		log.Printf("encountered an error during message unmarshall : %s", err.Error())
	}
	return toWork, err
}
