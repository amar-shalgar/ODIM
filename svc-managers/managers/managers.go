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

//Package managers ...
package managers

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"

	dmtf "github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	managersproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/managers"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-managers/mgrcommon"
	"github.com/ODIM-Project/ODIM/svc-managers/mgrmodel"
	"github.com/ODIM-Project/ODIM/svc-managers/mgrresponse"
)

// GetManagersCollection will get the all the managers(odimra, Plugins, Servers)
func (e *ExternalInterface) GetManagersCollection(req *managersproto.ManagerRequest) (response.RPC, error) {
	var resp response.RPC
	resp.Header = map[string]string{
		"Allow":             `"GET"`,
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	managers := mgrresponse.ManagersCollection{
		OdataContext: "/redfish/v1/$metadata#ManagerCollection.ManagerCollection",
		OdataID:      "/redfish/v1/Managers",
		OdataType:    "#ManagerCollection.ManagerCollection",
		Description:  "Managers view",
		Name:         "Managers",
	}
	var members []dmtf.Link

	// Add servers as manager in manager collection
	managersCollectionKeysArray, err := e.DB.GetAllKeysFromTable("Managers")
	if err != nil || len(managersCollectionKeysArray) == 0 {
		log.Error("odimra Doesnt have Servers")
	}

	for _, key := range managersCollectionKeysArray {
		members = append(members, dmtf.Link{Oid: key})
	}
	managers.Members = members
	managers.MembersCount = len(members)
	resp.Body = managers
	resp.StatusCode = http.StatusOK
	return resp, nil
}

// GetManagers will fetch individual manager details with the given ID
func (e *ExternalInterface) GetManagers(req *managersproto.ManagerRequest) response.RPC {
	var resp response.RPC
	resp.Header = map[string]string{
		"Allow":             `"GET"`,
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}

	if req.ManagerID == config.Data.RootServiceUUID {
		manager, err := e.getManagerDetails(req.ManagerID)
		if err != nil {
			log.Error("error getting manager details : " + err.Error())
			errArgs := []interface{}{"Managers", req.ManagerID}
			errorMessage := err.Error()
			resp = common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errorMessage,
				errArgs, nil)
			return resp
		}
		resp.Body = manager
	} else {

		requestData := strings.Split(req.ManagerID, ":")
		if len(requestData) <= 1 {
			resp = e.getPluginManagerResoure(requestData[0], req.URL)
			return resp
		}
		uuid := requestData[0]
		data, err := e.DB.GetManagerByURL(req.URL)
		if err != nil {
			log.Error("error getting manager details : " + err.Error())
			var errArgs = []interface{}{}
			var statusCode int
			var StatusMessage string
			errorMessage := err.Error()
			if errors.DBKeyNotFound == err.ErrNo() {
				errArgs = []interface{}{"Managers", req.ManagerID}

				statusCode = http.StatusNotFound
				StatusMessage = response.ResourceNotFound
			} else {
				statusCode = http.StatusInternalServerError
				StatusMessage = response.InternalError
			}
			resp = common.GeneralError(int32(statusCode), StatusMessage, errorMessage,
				errArgs, nil)
			return resp
		}
		var managerData map[string]interface{}
		jerr := json.Unmarshal([]byte(data), &managerData)
		if jerr != nil {
			errorMessage := "error unmarshalling manager details: " + jerr.Error()
			log.Error(errorMessage)
			resp = common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage,
				nil, nil)
			return resp
		}
		// extracting the Manager Type from the  managerData
		var managerType string
		if val, ok := managerData["ManagerType"]; ok {
			managerType = val.(string)
		}

		if managerType != common.ManagerTypeService && managerType != "" {
			deviceData, err := e.getResourceInfoFromDevice(req.URL, uuid, requestData[1])
			if err != nil {
				log.Error("Device " + req.URL + " is unreachable: " + err.Error())
				// Updating the state
				managerData["Status"] = map[string]string{
					"State": "Absent",
				}
			} else {
				jerr := json.Unmarshal([]byte(deviceData), &managerData)
				if jerr != nil {
					errorMessage := "error unmarshaling manager details: " + jerr.Error()
					log.Error(errorMessage)
					resp = common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage,
						nil, nil)
					return resp
				}
			}
			err = e.DB.UpdateManagersData(req.URL, managerData)
			if err != nil {
				errorMessage := "error while saving manager details: " + err.Error()
				log.Error(errorMessage)
				resp = common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage,
					nil, nil)
				return resp
			}
			dataBytes, err := json.Marshal(managerData)
			if err != nil {
				errorMessage := "error while marshalling manager details: " + err.Error()
				log.Error(errorMessage)
				resp = common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage,
					nil, nil)
				return resp
			}
			data = string(dataBytes)
		}
		var resource map[string]interface{}
		json.Unmarshal([]byte(data), &resource)
		resp.Body = resource
	}
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success
	return resp
}

func (e *ExternalInterface) getManagerDetails(id string) (mgrmodel.Manager, error) {
	var mgr mgrmodel.Manager
	var mgrData mgrmodel.RAManager

	data, err := e.DB.GetManagerByURL("/redfish/v1/Managers/" + id)
	if err != nil {
		return mgr, fmt.Errorf("unable to retrieve manager information: %v", err)
	}
	if err := json.Unmarshal([]byte(data), &mgrData); err != nil {
		return mgr, fmt.Errorf("unable to marshal manager information: %v", err)
	}
	return mgrmodel.Manager{
		OdataContext:    "/redfish/v1/$metadata#Manager.Manager",
		OdataID:         "/redfish/v1/Managers/" + id,
		OdataType:       "#Manager.v1_3_3.Manager",
		Name:            mgrData.Name,
		ManagerType:     mgrData.ManagerType,
		ID:              mgrData.ID,
		UUID:            mgrData.UUID,
		FirmwareVersion: mgrData.FirmwareVersion,
		Status: &mgrmodel.Status{
			State: mgrData.State,
		},
	}, nil
}

// GetManagersResource is used to fetch resource data. The function is supposed to be used as part of RPC
// For getting system resource information,  parameters need to be passed GetSystemsRequest .
// GetManagersResource holds the  Uuid,Url and Resourceid ,
// Url will be parsed from that search key will created
// There will be two return values for the fuction. One is the RPC response, which contains the
// status code, status message, headers and body and the second value is error.
func (e *ExternalInterface) GetManagersResource(req *managersproto.ManagerRequest) response.RPC {
	var resp response.RPC
	resp.Header = map[string]string{
		"Allow":             `"GET"`,
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}

	requestData := strings.Split(req.ManagerID, ":")
	if len(requestData) <= 1 {
		resp = e.getPluginManagerResoure(requestData[0], req.URL)
		return resp
	}
	uuid := requestData[0]
	urlData := strings.Split(req.URL, "/")
	var tableName string
	if req.ResourceID == "" {
		resourceName := urlData[len(urlData)-1]
		tableName = common.ManagersResource[resourceName]
	} else {
		tableName = urlData[len(urlData)-2]
	}

	data, err := e.DB.GetResource(tableName, req.URL)
	if err != nil {
		if errors.DBKeyNotFound == err.ErrNo() {
			var err error
			if data, err = e.getResourceInfoFromDevice(req.URL, uuid, requestData[1]); err != nil {
				errorMessage := "unable to get resource details from device: " + err.Error()
				log.Error(errorMessage)
				errArgs := []interface{}{tableName, req.ManagerID}
				return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errorMessage, errArgs, nil)
			}
		} else {
			errorMessage := "unable to get managers details: " + err.Error()
			log.Error(errorMessage)
			return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, []interface{}{}, nil)
		}
	}

	var resource map[string]interface{}
	json.Unmarshal([]byte(data), &resource)
	resp.Body = resource
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success

	return resp
}

// VirtualMediaActions
func (e *ExternalInterface) VirtualMediaActions(req *managersproto.ManagerRequest) response.RPC {
    var resp response.RPC
    /*// Checking payload for EjectMedia
    if strings.Contains(req.URL,"VirtualMedia.EjectMedia"){
        if len(req.RequestBody) != 0{
            errorMessage := "Request body must be null from odimra"
		    log.Error(errorMessage)
		    resp = common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errorMessage, []interface{}{}, nil)
		    return resp
        }
    }*/
    if strings.Contains(req.URL,"VirtualMedia.InsertMedia"){
        var vmInsert mgrmodel.VirtualMediaInsert
        // unmarshalling the volume
        err := json.Unmarshal(req.RequestBody, &vmInsert)
        if err != nil {
            errorMessage := "while unmarshaling the virtual media insert request: " + err.Error()
            log.Error(errorMessage)
            resp = common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errorMessage, []interface{}{}, nil)
            return resp
        }
        log.Info(vmInsert)
	    // Validating the request JSON properties for case sensitive
        invalidProperties, err := common.RequestParamsCaseValidator(req.RequestBody, vmInsert)
        if err != nil {
            errMsg := "error while validating request parameters for virtual media insert: " + err.Error()
            log.Error(errMsg)
            return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
        } else if invalidProperties != "" {
            errorMessage := "error: one or more properties given in the request body are not valid, ensure properties are listed in uppercamelcase "
            log.Error(errorMessage)
            response := common.GeneralError(http.StatusBadRequest, response.PropertyUnknown, errorMessage, []interface{}{invalidProperties}, nil)
            return response
        }
    }
	requestData := strings.Split(req.ManagerID, ":")
	uuid := requestData[0]

    target, gerr := mgrmodel.GetTarget(uuid)
	if gerr != nil {
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, gerr.Error(), nil, nil)
	}
	// Get the Plugin info
	plugin, gerr := mgrmodel.GetPluginData(target.PluginID)
	if gerr != nil {
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, gerr.Error(), nil, nil)
	}
	var contactRequest mgrcommon.PluginContactRequest
	contactRequest.ContactClient = e.Device.ContactClient
	contactRequest.Plugin = plugin

	if strings.EqualFold(plugin.PreferredAuthType, "XAuthToken") {
		token := mgrcommon.GetPluginToken(contactRequest)
		if token == "" {
			var errorMessage = "error: Unable to create session with plugin " + plugin.ID
			return common.GeneralError(http.StatusInternalServerError, response.InternalError, fmt.Sprintf(errorMessage), nil, nil)
		}

		contactRequest.Token = token
	} else {
		contactRequest.BasicAuth = map[string]string{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		}

	}
	decryptedPasswordByte, err := e.Device.DecryptDevicePassword(target.Password)
	if err != nil {
		errorMessage := "error while trying to decrypt device password: " + err.Error()
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, fmt.Sprintf(errorMessage), nil, nil)
	}
	contactRequest.DeviceInfo = map[string]interface{}{
		"ManagerAddress": target.ManagerAddress,
		"UserName":       target.UserName,
		"Password":       decryptedPasswordByte,
		"PostBody":  req.RequestBody,
	}
	//replace the uuid:system id with the system to the @odata.id from request url
	contactRequest.OID = strings.Replace(req.URL, uuid+":"+req.ManagerID, req.ManagerID, -1)
	contactRequest.HTTPMethodType = http.MethodPost
	log.Info("kkkkkkkkkkkkkkk1",string(req.RequestBody))
	target.PostBody = req.RequestBody
	body, _, getResp, err := mgrcommon.ContactPlugin(contactRequest, "error while getting the details "+contactRequest.OID+": ")
	if err != nil {
		/*if getResp.StatusCode == http.StatusUnauthorized && strings.EqualFold(contactRequest.Plugin.PreferredAuthType, "XAuthToken") {
			if body, _, _, err = RetryManagersOperation(contactRequest, "error while getting the details "+contactRequest.OID+": "); err != nil {
				return "", fmt.Errorf("error while trying to get data from plugin: %v", err)
			}
		} else {*/
			return common.GeneralError(http.StatusInternalServerError, response.InternalError, err.Error(), nil, nil)
		//}
	}
	resp.StatusCode = getResp.StatusCode
	//resp.StatusMessage = response.Success
	err = json.Unmarshal(body, &resp.Body)
	if err != nil {
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, err.Error(), nil, nil)
	}
	return resp
}

func (e *ExternalInterface) getPluginManagerResoure(managerID, reqURI string) response.RPC {
	var resp response.RPC
	data, dberr := e.DB.GetManagerByURL("/redfish/v1/Managers/" + managerID)
	if dberr != nil {
		log.Error("unable to get manager details : " + dberr.Error())
		var errArgs = []interface{}{"Managers", managerID}
		errorMessage := dberr.Error()
		resp = common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errorMessage,
			errArgs, nil)
		return resp
	}
	var managerData map[string]interface{}
	jerr := json.Unmarshal([]byte(data), &managerData)
	if jerr != nil {
		errorMessage := "unable to unmarshal manager details: " + jerr.Error()
		log.Error(errorMessage)
		resp = common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage,
			nil, nil)
		return resp
	}
	var pluginID = managerData["Name"].(string)
	// Get the Plugin info
	plugin, gerr := e.DB.GetPluginData(pluginID)
	if gerr != nil {
		log.Error("unable to get manager details : " + gerr.Error())
		var errArgs = []interface{}{"Plugin", pluginID}
		errorMessage := gerr.Error()
		resp = common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errorMessage,
			errArgs, nil)
		return resp
	}
	var req mgrcommon.PluginContactRequest

	req.ContactClient = e.Device.ContactClient
	req.Plugin = plugin

	if strings.EqualFold(plugin.PreferredAuthType, "XAuthToken") {
		token := mgrcommon.GetPluginToken(req)
		if token == "" {
			var errorMessage = "unable to create session with plugin " + plugin.ID
			return common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage,
				[]interface{}{}, nil)
		}
		req.Token = token
	} else {
		req.BasicAuth = map[string]string{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		}

	}
	req.OID = reqURI
	var errorMessage = "unable to get the details " + reqURI + ": "
	var header = map[string]string{"Content-type": "application/json; charset=utf-8"}
	body, _, getResponse, err := mgrcommon.ContactPlugin(req, errorMessage)
	if err != nil {
		if getResponse.StatusCode == http.StatusUnauthorized && strings.EqualFold(req.Plugin.PreferredAuthType, "XAuthToken") {
			if body, _, getResponse, err = mgrcommon.RetryManagersOperation(req, errorMessage); err != nil {
				resp.StatusCode = getResponse.StatusCode
				json.Unmarshal(body, &resp.Body)
				resp.Header = header
				return resp
			}
		} else {
			resp.StatusCode = getResponse.StatusCode
			json.Unmarshal(body, &resp.Body)
			resp.Header = header
			return resp
		}
	}
	return fillResponse(body)

}

func fillResponse(body []byte) response.RPC {
	var resp response.RPC
	data := string(body)
	//replacing the response with north bound translation URL
	for key, value := range config.Data.URLTranslation.NorthBoundURL {
		data = strings.Replace(data, key, value, -1)
	}
	var respData map[string]interface{}
	err := json.Unmarshal([]byte(data), &respData)
	if err != nil {
		log.Error(err.Error())
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, err.Error(),
			[]interface{}{}, nil)
	}
	resp.Body = respData
	resp.Header = map[string]string{
		"Allow":             `"GET"`,
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success
	return resp

}

func (e *ExternalInterface) getResourceInfoFromDevice(reqURL, uuid, systemID string) (string, error) {
	var getDeviceInfoRequest = mgrcommon.ResourceInfoRequest{
		URL:                   reqURL,
		UUID:                  uuid,
		SystemID:              systemID,
		ContactClient:         e.Device.ContactClient,
		DecryptDevicePassword: e.Device.DecryptDevicePassword,
	}
	return e.Device.GetDeviceInfo(getDeviceInfoRequest)

}
