/**
 * Copyright (C) 2015 Red Hat, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package restPack

import (
	restful "github.com/emicklei/go-restful"
	"gopkg.in/vmihailenco/msgpack.v2"
)

const MIME_MSGPACK = "application/x-msgpack" // Accept or Content-Type used in Consumes() and/or Produces()

// NewEntityAccessorMPack returns a new EntityReaderWriter for accessing MessagePack content.
// This package is not initialized with such an accessor using the MIME_MSGPACK contentType.
func NewEntityAccessorMsgPack() restful.EntityReaderWriter {
	return entityMsgPackAccess{}
}

// entityOctetAccess is a EntityReaderWriter for Octet encoding
type entityMsgPackAccess struct {
}

// Read unmarshalls the value from byte slice and using msgpack to unmarshal
func (e entityMsgPackAccess) Read(req *restful.Request, v interface{}) error {
	return msgpack.NewDecoder(req.Request.Body).Decode(v)
}

// Write marshals the value to byte slice and set the Content-Type Header.
func (e entityMsgPackAccess) Write(resp *restful.Response, status int, v interface{}) error {
	if v == nil {
		resp.WriteHeader(status)
		// do not write a nil representation
		return nil
	}
	resp.WriteHeader(status)
	return msgpack.NewEncoder(resp).Encode(v)
}
