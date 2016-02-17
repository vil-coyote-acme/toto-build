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
	"log"
)

type ProducerConfig struct {
	NsqAddr string
	Topic   string
}

type Producer struct {
	conf *ProducerConfig
	listeningChan chan string
}

func NewProducerConfig() *ProducerConfig {
	c := new (ProducerConfig)
	c.Topic = "report"
	return c
}

func NewProducer(conf *ProducerConfig) *Producer {
	p := new(Producer)
	p.conf = conf
	return p
}

// start waiting for report to send back to toto scheduler
func (prd *Producer) Start() chan string {
	config := nsq.NewConfig()
	p, errP := nsq.NewProducer(prd.conf.NsqAddr, config)
	if errP != nil {
		// todo implements test case to pass here
		log.Panicf("Error during nsq Producer init : %f", errP.Error())
	}
	prd.listeningChan = make(chan string)
	go func() {
		for msg := range prd.listeningChan {
			p.Publish(prd.conf.Topic, []byte(msg))
		}
	}();
	return prd.listeningChan
}

// stop the producer
func (prd *Producer) Stop() {
	// todo implements stop
}