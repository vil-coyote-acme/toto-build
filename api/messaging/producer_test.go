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
	"github.com/stretchr/testify/assert"
	"github.com/vil-coyote-acme/toto-build-common/message"
	"testing"
	"github.com/vil-coyote-acme/toto-build-common/broker"
	"github.com/vil-coyote-acme/toto-build-common/testtools"
)

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
	b := broker.NewBroker()
	b.Start()
	defer b.Stop()
	// when
	ch := make(chan message.Report, 20)
	p.Start(ch)
	mes := message.Report{int64(1), message.PENDING, []string{"test"}}
	ch <- mes
	// and test listener
	receip, consumer := testtools.SetupListener(c.Topic)
	// then
	assert.Equal(t, mes, <-receip)
	consumer.Stop()
}