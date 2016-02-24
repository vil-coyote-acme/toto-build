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
	"flag"
	"toto-build-agent/api/messaging"
	"github.com/vil-coyote-acme/toto-build-common/broker"
	"github.com/vil-coyote-acme/toto-build-common/message"
	"toto-build-agent/build"
	"log"
	"github.com/nsqio/go-nsq"
	"encoding/json"
)

var (
	nsqLookUpHost string
	nsqLookUpPort string
	brokerAddr string
)

func main() {
	flag.StringVar(&brokerAddr, "broker-addr", "127.0.0.1", "address of the broker. Should be accessible from scheduler")
	flag.StringVar(&nsqLookUpHost, "lookup-addrr", "127.0.0.1", "address of the lookup service. Used by toto-build to get topic used for communications")
	flag.StringVar(&nsqLookUpPort, "lookup-port", "4161", "port of the lookup service.")
	flag.Parse()
}

func startListening() {
	// first start the broker
	log.Print("Start listening for jobs")
	embeddedBroker := broker.NewBroker()
	embeddedBroker.StartBroker() //todo configure ip and port

	// start the report producer
	reportChan := make(chan message.Report, 20)
	producerConf := messaging.NewProducerConfig()
	producerConf.NsqAddr = "127.0.0.1:4150"//todo configure ip and port
	producer := messaging.NewProducer(producerConf)
	producer.Start(reportChan)

	// start the work listener
	listenerConf := messaging.NewListenerConfig()
	listenerConf.LookupAddr = []string{nsqLookUpHost + ":" + nsqLookUpPort}
	listener := messaging.NewListener(listenerConf)
	// must create the topic BEFORE listing on it
	createTopic(listenerConf.Topic)
	toWorkChan := listener.Start()

	go func() {
		// todo : put this code inside listener.go
		for toWork := range toWorkChan {
			log.Printf("receive one job : %s", toWork)
			go executeJob(toWork, reportChan)
		}
	}()
}

func createTopic(topic string) {
	config := nsq.NewConfig()
	p, _ := nsq.NewProducer("127.0.0.1:4150", config)//todo configure ip and port
	mess := message.ToWork{int64(1), message.HELLO, "HELLO"}
	body, _ := json.Marshal(mess)// todo handle this error case
	p.Publish("jobs", body)
}

func executeJob(toWork message.ToWork, reportChan chan message.Report) {
	var logsChan chan string
	switch toWork.Cmd {
	case message.PACKAGE :
		logsChan = build.BuildPackage(toWork.Package)
	case message.TEST:
		logsChan = build.TestPackage(toWork.Package)
	case message.HELLO:
		reportChan <- message.Report{toWork.JobId, message.WORKING, []string{"Hello"}}
	default:
	// todo handle this case
	}
	if logsChan != nil {
		for log := range logsChan {
			// todo handle buffered logs
			// todo handle job status
			reportChan <- message.Report{toWork.JobId, message.WORKING, []string{log}}
		}
	}
}
