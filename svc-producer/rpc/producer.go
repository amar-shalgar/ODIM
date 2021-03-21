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
	log "github.com/sirupsen/logrus"
	"net/http"

	producerproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/producer"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-producer/producer"
)

// RunEventProducer is an rpc handler
func RunEventProducer(ctx context.Context, req *producerproto.UpdateRequest, resp *producerproto.UpdateResponse) error {
    fillProtoResponse(resp, producer.RunEventProducer(req))
	return nil
}

// StopEventProducer is an rpc handler
func StopEventProducer(ctx context.Context, req *producerproto.UpdateRequest, resp *producerproto.UpdateResponse) error {
	fillProtoResponse(resp, producer.RunEventProducer(req))
	return nil
}
