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

//Package rpc ...
package rpc

import (
	"context"
	"fmt"

	producerproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/producer"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
)

// RunEventProducer defines the RPC call function for
// the GetUpdateService from update micro service
func RunEventProducer(req producerproto.ProducerRequest) (*producerproto.ProducerResponse, error) {

	update := producerproto.RunEventProducer(services.Producer, services.Service.Client())

	resp, err := update.RunEventProducer(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}

	return resp, err
}

// StopEventProducer defines the RPC call function for
// the GetFirmwareInventory from update micro service
func StopEventProducer(req producerproto.ProducerRequest) (*producerproto.ProducerResponse, error) {

	produce := producerproto.StopEventProducer(services.Producer, services.Service.Client())
	resp, err := produce.StopEventProducer(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}

	return resp, err
}
