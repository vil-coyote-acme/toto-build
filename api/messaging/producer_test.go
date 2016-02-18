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
	"encoding/json"
	"github.com/nsqio/go-nsq"
	"github.com/stretchr/testify/assert"
	"github.com/vil-coyote-acme/toto-build-common/message"
	"testing"
	"time"
	"toto-build-common/testtools"
)

type handlerTest struct {
	receip chan message.Report
}

func (h *handlerTest) HandleMessage(mes *nsq.Message) (e error) {
	var report message.Report
	json.Unmarshal(mes.Body, &report)
	h.receip <- report
	return e
}

func Test_NewProducerConfig(t *testing.T) {
	// when
	c := NewProducerConfig()
	// then
	assert.Equal(t, "report", c.Topic)
	assert.Equal(t, "", c.NsqAddr)
}

func Test_Producer_Start(t *testing.T) {
	// given
	c := NewProducerConfig()
	c.NsqAddr = "127.0.0.1:4150"
	p := NewProducer(c)
	// and broker initialization
	b := testtools.NewBroker()
	b.Start()
	defer b.Stop()
	// when
	ch := make(chan message.Report, 20)
	p.Start(ch)
	mes := message.Report{int64(1), message.PENDING, []string{"test"}}
	ch <- mes
	// and test listener
	receip, consumer := setupListener(c.Topic)
	// then
	assert.Equal(t, mes, <-receip)
	consumer.Stop()
}

func setupListener(topic string) (chan message.Report, *nsq.Consumer) {
	duration, _ := time.ParseDuration("300ms")
	time.Sleep(duration)
	handler := new(handlerTest)
	receip := make(chan message.Report, 2)
	handler.receip = receip
	consumer, _ := nsq.NewConsumer(topic, "scheduler", nsq.NewConfig())
	consumer.AddHandler(handler)
	consumer.ConnectToNSQLookupds([]string{"127.0.0.1:4161"})
	return receip, consumer
}