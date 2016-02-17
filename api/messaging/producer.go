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

type ProducerConfig struct {
	NsqAddr string
	Topic   string
}

type Producer struct {
	conf *ProducerConfig
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
func (p *Producer) Start() chan string {
	return
}

// stop the producer
func (p *Producer) Stop() {

}