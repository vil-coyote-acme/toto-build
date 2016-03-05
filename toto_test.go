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
	"github.com/nsqio/go-nsq"
	"github.com/stretchr/testify/assert"
	"github.com/vil-coyote-acme/toto-build-common/broker"
	"github.com/vil-coyote-acme/toto-build-common/message"
	"github.com/vil-coyote-acme/toto-build-common/testtools"
	"testing"
	"github.com/vil-coyote-acme/toto-build-common/logs"
)

func Test_Main_should_Parse_Arguments(t *testing.T) {
	// when
	main()
	// then
	assert.Equal(t, "127.0.0.1", brokerAddr)
	assert.Equal(t, "4150", brokerPort)
	assert.Equal(t, "127.0.0.1", nsqLookUpHost)
	assert.Equal(t, "4161", nsqLookUpPort)
}

func Test_Main_should_Start_An_Nsq_Service(t *testing.T) {
	//given
	initVar()
	b := startLookUp()
	defer b.Stop()
	// when
	startListening()
	defer graceFullShutDown()
	sendMsg()
	// then
	receip, consumer := testtools.SetupListener("report", b.LookUpHttpAddrr + ":" + b.LookUpHttpPort)
	defer close(receip)
	defer consumer.Stop()
	assert.NotNil(t, consumer)
	assert.NotNil(t, receip)
	// first get the hello from the agent, then the message with the build log and finally the last report with
	// the build status (should be SUCCESS here)
	hello := <-receip
	buildTrace := <-receip
	finalReport := <-receip

	assert.Contains(t, testtools.FromSliceToString(hello.Logs), "Hello")
	assert.Equal(t, hello.Status, message.SUCCESS)
	assert.Contains(t, testtools.FromSliceToString(buildTrace.Logs), "toto-build-agent/testapp")
	assert.Equal(t, buildTrace.Status, message.WORKING)
	assert.Equal(t, finalReport.Status, message.SUCCESS)
}

func initVar() {
	brokerAddr = "127.0.0.1"
	brokerPort = "4150"
	nsqLookUpHost = "127.0.0.1"
	nsqLookUpPort = "4161"
	logger = logs.NewLogger("[TOTO-BUILD main] : ", logs.NewConsoleAppender(logs.INFO))
}

func startLookUp() *broker.Broker {
	b := broker.NewBroker()
	b.StartLookUp()
	return b
}

func sendMsg() {
	// test message creation
	mess := message.ToWork{int64(1), message.TEST, "toto-build-agent/testapp", "go1.6", ""}
	body, _ := json.Marshal(mess)
	// message sending
	config := nsq.NewConfig()
	p, _ := nsq.NewProducer("127.0.0.1:4150", config)
	p.Publish("jobs", body)
}
