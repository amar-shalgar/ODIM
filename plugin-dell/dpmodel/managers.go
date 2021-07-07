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

//Package dpmodel
package dpmodel

// VirtualMediaInsert struct is used to hold the insert virtual media request payload
type VirtualMediaInsert struct {
	Image                string `json:"Image"`
	Inserted             bool   `json:"Inserted"`
	Password             string `json:"Password,omitempty"`
	TransferMethod       string `json:"TransferMethod,omitempty"`
	TransferProtocolType string `json:"TransferProtocolType,omitempty"`
	UserName             string `json:"UserName,omitempty"`
	WriteProtected       bool   `json:"WriteProtected"`
}

// VirtualMediaInsertSouthBound struct is used for device request payload
type VirtualMediaInsertSouthBound struct {
	Image          string `json:"Image"`
	Inserted       bool   `json:"Inserted"`
	WriteProtected bool   `json:"WriteProtected"`
}
