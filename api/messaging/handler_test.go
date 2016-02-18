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
package messaging_test

import (
	"encoding/json"
	"github.com/nsqio/go-nsq"
	"github.com/stretchr/testify/assert"
	"github.com/vil-coyote-acme/toto-build-common/message"
	"testing"
	"toto-build-agent/api/messaging"
)

func Test_HandleMessage_Should_Return_Error_when_No_ToWork(t *testing.T) {
	// given
	c := make(chan message.ToWork, 5)
	defer close(c)
	h := messaging.NewHandler(c)
	var messId [16]byte
	// when
	err := h.HandleMessage(nsq.NewMessage(messId, make([]byte, 256)))
	// then
	assert.NotNil(t, err)
}

func Test_HandleMessage_Should_Emit_ToWork(t *testing.T) {
	// given
	c := make(chan message.ToWork, 5)
	defer close(c)
	h := messaging.NewHandler(c)
	var messId [16]byte
	mess := message.ToWork{int64(1), message.TEST, "myPkg"}
	body, _ := json.Marshal(mess)
	// when
	err := h.HandleMessage(nsq.NewMessage(messId, body))
	// then
	assert.Nil(t, err)
	assert.Equal(t, mess, <-c)
}

func Test_HandleMessage_Should_Return_Error_when_Buffer_full(t *testing.T) {
	// given
	c := make(chan message.ToWork, 2)
	defer close(c)
	h := messaging.NewHandler(c)
	var messId [16]byte
	mess := message.ToWork{int64(3), message.TEST, "myPkg3"}
	body, _ := json.Marshal(mess)
	// and
	c <- message.ToWork{int64(1), message.TEST, "myPkg1"}
	c <- message.ToWork{int64(2), message.TEST, "myPkg2"}
	// when
	err := h.HandleMessage(nsq.NewMessage(messId, body))
	// then
	assert.NotNil(t, err)
	assert.Equal(t, "the buffer is full", err.Error())
}
