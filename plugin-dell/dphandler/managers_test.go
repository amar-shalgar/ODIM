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

//Package dphandler
package dphandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	iris "github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
	"github.com/ODIM-Project/ODIM/plugin-dell/config"
	"github.com/ODIM-Project/ODIM/plugin-dell/dpmodel"
	"github.com/ODIM-Project/ODIM/plugin-dell/dpresponse"
)

func mockManagers(username, password, url string, w http.ResponseWriter) {
	body := `{"data": "success"}`
	if url == "/ODIM/v1/Managers/1/VirtualMedia/1/Actions/VirtualMedia.InsertMedia" && username == "admin" {
		e, _ := json.Marshal(body)
		w.WriteHeader(http.StatusOK)
		w.Write(e)
		return
	}
	if url == "/ODIM/v1/Managers/1/VirtualMedia/1/Actions/VirtualMedia.InsertMedia" && username != "admin" {
		e, _ := json.Marshal(body)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(e)
		return
	}
	if url == "/ODIM/v1/Managers/1/VirtualMedia/1/Actions/VirtualMedia.EjectMedia" && username == "admin" {
		e, _ := json.Marshal(body)
		w.WriteHeader(http.StatusOK)
		w.Write(e)
		return
	}
	if url == "/ODIM/v1/Managers/1/VirtualMedia/1/Actions/VirtualMedia.EjectMedia" && username != "admin" {
		e, _ := json.Marshal(body)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(e)
		return
	}
	return
}

func TestGetManagerCollection(t *testing.T) {
	config.SetUpMockConfig(t)
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/ODIM/v1")

	redfishRoutes.Get("/Managers", GetManagersCollection)

	dpresponse.PluginToken = "token"
	e := httptest.New(t, mockApp)

	var deviceDetails = dpmodel.Device{
		Host: "",
	}
	//Unit Test for success scenario
	e.GET("/ODIM/v1/Managers").WithJSON(deviceDetails).Expect().Status(http.StatusOK)

	//Case for invalid token
	e.GET("/ODIM/v1/Managers").WithHeader("X-Auth-Token", "Invalidtoken").WithJSON(deviceDetails).Expect().Status(http.StatusUnauthorized)

}

func TestGetManager(t *testing.T) {
	config.SetUpMockConfig(t)
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/ODIM/v1")

	redfishRoutes.Get("/Managers", GetManagersInfo)

	dpresponse.PluginToken = "token"
	e := httptest.New(t, mockApp)

	var deviceDetails = dpmodel.Device{
		Host: "",
	}
	//Unit Test for success scenario
	e.GET("/ODIM/v1/Managers").WithJSON(deviceDetails).Expect().Status(http.StatusOK)

	//Case for invalid token
	e.GET("/ODIM/v1/Managers").WithHeader("X-Auth-Token", "Invalidtoken").WithJSON(deviceDetails).Expect().Status(http.StatusUnauthorized)

}

func TestVirtualMediaActions(t *testing.T) {
	config.SetUpMockConfig(t)
	deviceHost := "localhost"
	devicePort := "1234"
	ts := startTestServer(mockManagers)
	// Start the server.
	ts.StartTLS()
	defer ts.Close()
    time.Sleep(1 * time.Second)
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1")
	redfishRoutes.Post("/Managers/1/VirtualMedia/1/Actions/VirtualMedia.InsertMedia", VirtualMediaActions)
	redfishRoutes.Post("/Managers/1/VirtualMedia/1/Actions/VirtualMedia.EjectMedia", VirtualMediaActions)
	dpresponse.PluginToken = "token"

	test := httptest.New(t, mockApp)
	attributes := map[string]interface{}{"Image": "abc"}
	attributeByte, _ := json.Marshal(attributes)
	requestBody := map[string]interface{}{
		"ManagerAddress": fmt.Sprintf("%s:%s", deviceHost, devicePort),
		"UserName":       "admin",
		"Password":       []byte("P@$$w0rd"),
		"PostBody":       attributeByte,
	}
	test.POST("/redfish/v1/Managers/1/VirtualMedia/1/Actions/VirtualMedia.InsertMedia").WithJSON(requestBody).Expect().Status(http.StatusOK)
	requestBody["UserName"] = "invalid"
	test.POST("/redfish/v1/Managers/1/VirtualMedia/1/Actions/VirtualMedia.InsertMedia").WithJSON(requestBody).Expect().Status(http.StatusBadRequest)

	requestBody = map[string]interface{}{
		"ManagerAddress": fmt.Sprintf("%s:%s", deviceHost, devicePort),
		"UserName":       "admin",
		"Password":       []byte("P@$$w0rd"),
	}
	test.POST("/redfish/v1/Managers/1/VirtualMedia/1/Actions/VirtualMedia.EjectMedia").WithJSON(requestBody).Expect().Status(http.StatusOK)
	requestBody["UserName"] = "invalid"
	test.POST("/redfish/v1/Managers/1/VirtualMedia/1/Actions/VirtualMedia.EjectMedia").WithJSON(requestBody).Expect().Status(http.StatusBadRequest)
}
