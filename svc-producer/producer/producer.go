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

//Package update ...
package producer

// ---------------------------------------------------------------------------------------
// IMPORT Section
//
import (
	//"github.com/ODIM-Project/ODIM/lib-utilities/common"
	//"github.com/ODIM-Project/ODIM/lib-utilities/config"
	producerproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/producer"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-producer/pcommon"
	"github.com/ODIM-Project/ODIM/svc-producer/pmodel"
	log "github.com/sirupsen/logrus"

	"net/http"
	"time"
	"github.com/go-redis/redis"
	"encoding/json"
)

func RunEventProducer(req *producerproto.ProducerRequest) response.RPC {
    var resp response.RPC
    client := pmodel.ConnectRedis()
    addProduceEventsKey(client, req)
    go produce(client)

    resp.Header = map[string]string{
		"Allow":             `"GET"`,
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	log.Info("Event producer started")
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success
	return resp
}

func StopEventProducer(req *producerproto.ProducerRequest) response.RPC {
    var resp response.RPC
    client := pmodel.ConnectRedis()
    defer client.Close()
    addProduceEventsKey(client, req)

    resp.Header = map[string]string{
		"Allow":             `"GET"`,
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	log.Info("Event producer stopped")
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success
	return resp
}

func addProduceEventsKey(client *redis.Client, req *producerproto.ProducerRequest) error{
    log.Info(req)
    err := client.Set("ProduceEvents","true",0).Err()
    if err!= nil{
        return err
    }
    return nil
}

func produce(client *redis.Client) {
    defer client.Close()
    dt := time.Now()
	event := pcommon.EventMessageData{
        OdataType : "#Event.v1_2_1.Event",
		Name      : "Event Array",
		Context   : "/redfish/v1/$metadata#Event.Event",
		Events    : []pcommon.Event{
					pcommon.Event{MemberID    : "5615",
					EventType         : "Alert",
					EventGroupID      : 1,
					EventID           : "8491",
					Severity          : "Informational",
					EventTimestamp    : dt.String(),
					Message           : "Successfully logged in using admin, from 10.18.1.216 and REDFISH.",
					MessageID         : "USR0030",
					OriginOfCondition : pcommon.OdataID{
											OdataID : "/redfish/v1/Managers/iDRAC.Embedded.1",
										},
				}},
    }

	body,_ := json.Marshal(event)
	i := 0;
	for {
	    produceEvents, _ := client.Get("ProduceEvents").Result()
        log.Info("ProduceEvents", produceEvents)
        if(produceEvents == "false"){
            break
        }
        strCMD := client.XAdd(&redis.XAddArgs{
            Stream: "test",
            Values: map[string]interface{}{
                "ID": i,
                "data": body,
                },
        })
        newID, err := strCMD.Result()
        if err != nil {
            log.Error("Event generation error:%v\n", err)
        } else {
            log.Info("Event generated:%v,%v\n", i,newID)
        }
        i++
    }
}
