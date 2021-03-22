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

//Package consumer ...
package consumer

// ---------------------------------------------------------------------------------------
// IMPORT Section
//
import (
	//"github.com/ODIM-Project/ODIM/lib-utilities/common"
	//"github.com/ODIM-Project/ODIM/lib-utilities/config"
	consumerproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/consumer"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	//"github.com/ODIM-Project/ODIM/svc-consumer/ccommon"
	"github.com/ODIM-Project/ODIM/svc-consumer/cmodel"
	log "github.com/sirupsen/logrus"

	"net/http"
	"time"
	"github.com/go-redis/redis"
	//"encoding/json"
)

func RunEventConsumer(req *consumerproto.ConsumerRequest) response.RPC {
    var resp response.RPC
    client := cmodel.ConnectRedis()
    addConsumeEventsKey(client, req, "true")
    go consume(client)

    resp.Header = map[string]string{
		"Allow":             `"GET"`,
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	log.Info("Event consumer started")
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success
	return resp
}

func StopEventConsumer(req *consumerproto.ConsumerRequest) response.RPC {
    var resp response.RPC
    client := cmodel.ConnectRedis()
    //defer client.Close()
    addConsumeEventsKey(client, req, "false")

    resp.Header = map[string]string{
		"Allow":             `"GET"`,
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	log.Info("Event consumer stopped")
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success
	return resp
}

func addConsumeEventsKey(client *redis.Client, req *consumerproto.ConsumerRequest, val string) error{
    log.Info(req)
    err := client.Set("ConsumeEvents", val,0).Err()
    if err!= nil{
        return err
    }
    return nil
}

func consume(client *redis.Client) {
	//defer client.Close()
	for {
		consumeEvents, _ := client.Get("ConsumeEvents").Result()
        //log.Info("ConsumeEvents", consumeEvents)
        if(consumeEvents == "false"){
            break
        }
		res := client.XRead(&redis.XReadArgs{
				Streams: []string{"test", "$"},
				Count:   2,
				Block:   100 * time.Millisecond,
		})
		respBody, _ := res.Result()
		if(len(respBody)>0){
			log.Info("Event received: ",respBody)
		}
	}
}
