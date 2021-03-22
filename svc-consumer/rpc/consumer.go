//(C) Copyright [2020] Hewlett Packard Enterprise Development LP
//
//Licensed under the Apache License, Version 2.0 (the "License"); you may
//not use this file except in compliance with the License. You may obtain
//a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//License for the specific language governing permissions and limitations
// under the License.

package rpc

import (
	"context"
	//log "github.com/sirupsen/logrus"
	//"net/http"

	consumerproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/consumer"
	//"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-consumer/consumer"
)

type Consumer struct{}

// RunEventConsumer is an rpc handler
func (p *Consumer) RunEventConsumer(ctx context.Context, req *consumerproto.ConsumerRequest, resp *consumerproto.ConsumerResponse) error {
    fillProtoResponse(resp, consumer.RunEventConsumer(req))
	return nil
}

// StopEventConsumer is an rpc handler
func  (p *Consumer) StopEventConsumer(ctx context.Context, req *consumerproto.ConsumerRequest, resp *consumerproto.ConsumerResponse) error {
	fillProtoResponse(resp, consumer.StopEventConsumer(req))
	return nil
}
