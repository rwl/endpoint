// Copyright 2013 Richard lincoln All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"net/http"
	"strings"
	"fmt"
)

type route struct {
	pathPrefix string
	handler    func(http.ResponseWriter, *http.Request)
}

type router struct {
	routes []*route
}

func newRouter() *router {
	return &router{}
}

func (h *router) HandleFunc(pathPrefix string, handler func(http.ResponseWriter, *http.Request)) {
	r := &route{pathPrefix, handler}
	h.routes = append(h.routes, r)
}

// Routes the request if it matches one of our reserved URLs. This also
// handles OPTIONS CORS requests.
func (h *router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range h.routes {
		if strings.HasPrefix(r.URL.Path, route.pathPrefix) {
			route.handler(w, r)
			return
		}
	}

	if r.Method == "OPTIONS" {
		corsHandler := newCheckCorsHeaders(r)
		if corsHandler.allowCorsRequest {
			// The server returns 200 rather than 204, for some reason.
			corsHandler.updateHeaders(w.Header())
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "")
		}
	}
}
