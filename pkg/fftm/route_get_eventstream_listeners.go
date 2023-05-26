// Copyright © 2023 Kaleido, Inc.
//
// SPDX-License-Identifier: Apache-2.0
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

package fftm

import (
	"net/http"

	"github.com/hyperledger/firefly-common/pkg/ffapi"
	"github.com/hyperledger/firefly-transaction-manager/internal/persistence"
	"github.com/hyperledger/firefly-transaction-manager/internal/tmmsgs"
	"github.com/hyperledger/firefly-transaction-manager/pkg/apitypes"
)

var getEventStreamListeners = func(m *manager) *ffapi.Route {
	route := &ffapi.Route{
		Name:   "getEventStreamListeners",
		Path:   "/eventstreams/{streamId}/listeners",
		Method: http.MethodGet,
		PathParams: []*ffapi.PathParam{
			{Name: "streamId", Description: tmmsgs.APIParamStreamID},
		},
		Description:     tmmsgs.APIEndpointGetEventStreamListeners,
		JSONInputValue:  nil,
		JSONOutputValue: func() interface{} { return []*apitypes.Listener{} },
		JSONOutputCodes: []int{http.StatusOK},
	}
	if m.richQueryAPI {
		route.FilterFactory = persistence.ListenerFilters
		route.JSONHandler = func(r *ffapi.APIRequest) (output interface{}, err error) {
			return r.FilterResult(m.getStreamListenersRich(r.Req.Context(), r.PP["streamId"], r.Filter))
		}
	} else {
		// Very limited query support
		route.QueryParams = []*ffapi.QueryParam{
			{Name: "limit", Description: tmmsgs.APIParamLimit},
			{Name: "after", Description: tmmsgs.APIParamAfter},
		}
		route.JSONHandler = func(r *ffapi.APIRequest) (output interface{}, err error) {
			return m.getStreamListenersByCreateTime(r.Req.Context(), r.QP["after"], r.QP["limit"], r.PP["streamId"])
		}
	}
	return route
}
