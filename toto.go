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
package main

import (
	"encoding/json"
	"flag"
	"github.com/nsqio/go-nsq"
	"github.com/vil-coyote-acme/toto-build-common/broker"
	"github.com/vil-coyote-acme/toto-build-common/message"
	"log"
	"toto-build-agent/api/messaging"
	"toto-build-agent/build"
)

var (
	nsqLookUpHost  string
	nsqLookUpPort  string
	brokerAddr     string
	brokerPort     string
	embeddedBroker *broker.Broker
	listener       *messaging.Listener
	toWorkChan     chan message.ToWork
	reportChan     chan message.Report
)

func main() {
	flag.StringVar(&brokerAddr, "broker-addr", "127.0.0.1", "address of the broker. Should be accessible from scheduler")
	flag.StringVar(&brokerPort, "broker-port", "4150", "port of the broker. Should be accessible from scheduler")
	flag.StringVar(&nsqLookUpHost, "lookup-addrr", "127.0.0.1", "address of the lookup service. Used by toto-build to get topic used for communications")
	flag.StringVar(&nsqLookUpPort, "lookup-port", "4161", "port of the lookup service.")
	flag.Parse()
}

// start listening for toWork
// will start an embeded broker, a report producer and a toWork listener
func startListening() {
	// first start the broker
	log.Print("Start listening for jobs")
	embeddedBroker = broker.NewBroker()
	embeddedBroker.BrokerAddr = brokerAddr
	embeddedBroker.BrokerPort = brokerPort
	embeddedBroker.StartBroker()

	// start the report producer
	reportChan = make(chan message.Report, 20)
	producerConf := messaging.NewProducerConfig()
	producerConf.NsqAddr = brokerAddr + ":" + brokerPort
	producer := messaging.NewProducer(producerConf)
	producer.Start(reportChan)

	// start the work listener
	listenerConf := messaging.NewListenerConfig()
	listenerConf.LookupAddr = []string{nsqLookUpHost + ":" + nsqLookUpPort}
	listener = messaging.NewListener(listenerConf)
	// must create the topic BEFORE listing on it
	sayHello(listenerConf.Topic)
	toWorkChan = listener.Start()

	// and then start executing incoming job
	build.ExecuteJob(toWorkChan, reportChan)
}

func graceFullShutDown() {
	embeddedBroker.Stop()
	listener.Stop()
	close(toWorkChan)
	close(reportChan)

}

// will publish the first report : a hello one
func sayHello(topic string) {
	config := nsq.NewConfig()
	p, _ := nsq.NewProducer(brokerAddr+":"+brokerPort, config)
	body, _ := json.Marshal(message.ToWork{int64(1), message.HELLO, "HELLO"}) // todo handle this error case
	p.Publish("jobs", body)
}
