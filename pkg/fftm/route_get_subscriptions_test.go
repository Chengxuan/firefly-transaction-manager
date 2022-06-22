// Copyright © 2022 Kaleido, Inc.
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
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/hyperledger/firefly-transaction-manager/mocks/ffcapimocks"
	"github.com/hyperledger/firefly-transaction-manager/pkg/apitypes"
	"github.com/hyperledger/firefly-transaction-manager/pkg/ffcapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetSubscriptions(t *testing.T) {

	url, m, done := newTestManager(t, func(w http.ResponseWriter, r *http.Request) {})
	defer done()

	err := m.Start()
	assert.NoError(t, err)

	mfc := m.connector.(*ffcapimocks.API)
	mfc.On("EventListenerVerifyOptions", mock.Anything, mock.Anything).Return(&ffcapi.EventListenerVerifyOptionsResponse{}, ffcapi.ErrorReason(""), nil)
	mfc.On("EventListenerAdd", mock.Anything, mock.Anything).Return(&ffcapi.EventListenerAddResponse{}, ffcapi.ErrorReason(""), nil)
	mfc.On("EventListenerRemove", mock.Anything, mock.Anything).Return(&ffcapi.EventListenerRemoveResponse{}, ffcapi.ErrorReason(""), nil).Maybe()

	// Create a stream
	var es1 apitypes.EventStream
	res, err := resty.New().R().SetBody(&apitypes.EventStream{Name: strPtr("stream1")}).SetResult(&es1).Post(url + "/eventstreams")
	assert.NoError(t, err)

	// Create some listeners
	var l1, l2 apitypes.Listener
	res, err = resty.New().R().SetBody(&apitypes.Listener{Name: strPtr("listener1"), StreamID: es1.ID}).SetResult(&l1).Post(url + "/subscriptions")
	assert.NoError(t, err)
	res, err = resty.New().R().SetBody(&apitypes.Listener{Name: strPtr("listener2"), StreamID: es1.ID}).SetResult(&l2).Post(url + "/subscriptions")
	assert.NoError(t, err)

	// Then get it
	var listeners []*apitypes.Listener
	res, err = resty.New().R().
		SetResult(&listeners).
		Get(url + "/subscriptions?limit=1&after=" + l1.ID.String())
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode())

	assert.Len(t, listeners, 1)
	assert.Equal(t, l2.ID, listeners[0].ID)
	assert.Equal(t, es1.ID, listeners[0].StreamID)

	mfc.AssertExpectations(t)

}
