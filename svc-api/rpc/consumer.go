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

	consumerproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/consumer"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
)

// RunEventConsumer defines the RPC call function for
// the GetUpdateService from update micro service
func RunEventConsumer(req consumerproto.ConsumerRequest) (*consumerproto.ConsumerResponse, error) {

	run := consumerproto.NewConsumerService(services.Consumer, services.Service.Client())

	resp, err := run.RunEventConsumer(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}

	return resp, err
}

// StopEventConsumer defines the RPC call function for
// the GetFirmwareInventory from update micro service
func StopEventConsumer(req consumerproto.ConsumerRequest) (*consumerproto.ConsumerResponse, error) {

	consume := consumerproto.NewConsumerService(services.Consumer, services.Service.Client())
	resp, err := consume.StopEventConsumer(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}

	return resp, err
}

