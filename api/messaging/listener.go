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
	"github.com/nsqio/go-nsq"
	"github.com/vil-coyote-acme/toto-build-common/message"
	"log"
)

type ListenerConfig struct {
	LookupAddr []string
	Topic      string
	Channel    string
	BuffSize   int
}

type Listener struct {
	conf     *ListenerConfig
	consumer *nsq.Consumer
}

// initialize new config for listener
func NewListenerConfig() *ListenerConfig {
	c := new(ListenerConfig)
	c.Topic = "jobs"
	c.Channel = "buld-agent"
	c.BuffSize = 10
	return c
}

// initialize new listener
func NewListener(conf *ListenerConfig) *Listener {
	l := new(Listener)
	l.conf = conf
	return l
}

// start listening for incoming ToWork
func (l *Listener) Start() chan message.ToWork {
	c := make(chan message.ToWork, l.conf.BuffSize)
	cons, err := nsq.NewConsumer(l.conf.Topic, l.conf.Channel, nsq.NewConfig())
	if err != nil {
		log.Panicf("error when trying to create a consumer for topic : %v and channel : %v", l.conf.Topic, l.conf.Channel)
	}
	l.consumer = cons
	// maybe possible to handle message in multiple goroutines
	log.Print("starting connecting broker listener")
	l.consumer.AddHandler(NewHandler(c))
	l.consumer.ConnectToNSQLookupds(l.conf.LookupAddr)
	log.Print("Listener started")
	return c
}

// todo tests
func (l *Listener) Stop() {
	if l.consumer != nil {
		l.consumer.Stop()
	}
}
