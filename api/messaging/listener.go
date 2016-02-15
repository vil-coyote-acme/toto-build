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
	"github.com/vil-coyote-acme/toto-build-common/message"
	"github.com/nsqio/go-nsq"
	"log"
)

type Config struct {
	LookupAddr []string
	Topic      string
	Channel    string
	BuffSize   int
}

type Listener struct {
	conf *Config
}

func NewConfig() *Config {
	c := new(Config)
	c.LookupAddr = []string{"127.0.0.1:4161"}
	c.Topic = "jobs"
	c.Channel = "buld-agent"
	c.BuffSize = 10
	return c
}

func NewListener(conf *Config) *Listener {
	l := new(Listener)
	l.conf = conf
	return l
}

func (l *Listener) StartListening() chan message.ToWork {
	c := make(chan message.ToWork, l.conf.BuffSize)
	cons, err := nsq.NewConsumer(l.conf.Topic, l.conf.Channel, nsq.NewConfig())
	if err != nil {
		log.Panicf("error when trying to create a consumer for topic : %v and channel : %v", l.conf.Topic, l.conf.Channel)
	}
	// maybe possible to handle message in multiple goroutines
	cons.AddHandler(NewHandler(c))
	cons.ConnectToNSQLookupds(l.conf.LookupAddr)
	return c
}
