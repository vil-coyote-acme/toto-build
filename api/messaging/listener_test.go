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
	"testing"
	"github.com/stretchr/testify/assert"
	"toto-build-agent/api/messaging"
	"github.com/vil-coyote-acme/toto-build-common/testtools"
	"github.com/vil-coyote-acme/toto-build-common/message"
	"encoding/json"
	"github.com/nsqio/go-nsq"
)

func Test_NewConfig(t *testing.T) {
	// when
	conf := messaging.NewConfig()
	// then
	assert.Equal(t, 1, len(conf.LookupAddr))
	assert.Equal(t, "127.0.0.1:4161", conf.LookupAddr[0])
	assert.Equal(t, "buld-agent", conf.Channel)
	assert.Equal(t, "jobs", conf.Topic)
	assert.Equal(t, 10, conf.BuffSize)
}

func Test_Reception_Of_One_ToWork_Message(t *testing.T) {
	// given the listener
	c := messaging.NewConfig()
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
	incomingChan := l.StartListening()
	// then
	assert.Equal(t, message.ToWork{int64(1), message.TEST, "myPkg"}, <- incomingChan)
}

