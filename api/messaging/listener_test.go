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
	"github.com/vil-coyote-acme/toto-build-common/testtools"
	"testing"
	"toto-build-agent/api/messaging"
)

func Test_NewConfig(t *testing.T) {
	// when
	conf := messaging.NewListenerConfig()
	// then
	assert.Equal(t, "buld-agent", conf.Channel)
	assert.Equal(t, "jobs", conf.Topic)
	assert.Equal(t, 10, conf.BuffSize)
}

func Test_Reception_Of_One_ToWork_Message(t *testing.T) {
	// given the listener
	c := messaging.NewListenerConfig()
	c.LookupAddr = []string{"127.0.0.1:4161"}
	l := messaging.NewListener(c)
	// broker initialization
	b := testtools.NewBroker()
	b.Start()
	defer b.Stop()
	// test message creation
	mess := message.ToWork{int64(1), message.TEST, "myPkg"}
	body, _ := json.Marshal(mess)
	// message sending
	config := nsq.NewConfig()
	p, _ := nsq.NewProducer("127.0.0.1:4150", config)
	p.Publish(c.Topic, body)
	// when
	incomingChan := l.Start()
	// then
	assert.Equal(t, message.ToWork{int64(1), message.TEST, "myPkg"}, <-incomingChan)
}
