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
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hyperledger/firefly-common/pkg/fftypes"
	"github.com/hyperledger/firefly-transaction-manager/mocks/ffcapimocks"
	"github.com/hyperledger/firefly-transaction-manager/pkg/apitypes"
	"github.com/hyperledger/firefly-transaction-manager/pkg/ffcapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNonceStaleStateContention(t *testing.T) {

	_, m, cancel := newTestManager(t)
	defer cancel()

	// Write a stale record to persistence
	oldTime := fftypes.FFTime(time.Now().Add(-100 * time.Hour))
	err := m.persistence.WriteTransaction(m.ctx, &apitypes.ManagedTX{
		ID:         "stale1",
		Status:     apitypes.TxStatusSucceeded,
		SequenceID: apitypes.UUIDVersion1(),
		Nonce:      fftypes.NewFFBigInt(1000), // old nonce
		Created:    &oldTime,                  // old record
		Request: &apitypes.TransactionRequest{
			TransactionInput: ffcapi.TransactionInput{
				TransactionHeaders: ffcapi.TransactionHeaders{
					From: "0x12345",
				},
			},
		},
	}, true)
	assert.NoError(t, err)

	mFFC := m.connector.(*ffcapimocks.API)

	mFFC.On("NextNonceForSigner", mock.Anything, mock.MatchedBy(func(nonceReq *ffcapi.NextNonceForSignerRequest) bool {
		return "0x12345" == nonceReq.Signer
	})).Return(&ffcapi.NextNonceForSignerResponse{
		Nonce: fftypes.NewFFBigInt(1111),
	}, ffcapi.ErrorReason(""), nil)

	locked1 := make(chan struct{})
	done1 := make(chan struct{})
	done2 := make(chan struct{})

	go func() {
		defer close(done1)

		ln, err := m.assignAndLockNonce(context.Background(), "ns1:"+fftypes.NewUUID().String(), "0x12345")
		assert.NoError(t, err)
		assert.Equal(t, uint64(1111), ln.nonce)
		close(locked1)

		time.Sleep(1 * time.Millisecond)
		ln.spent = &apitypes.ManagedTX{
			ID:         "ns1:" + fftypes.NewUUID().String(),
			Nonce:      fftypes.NewFFBigInt(int64(ln.nonce)),
			Status:     apitypes.TxStatusPending,
			SequenceID: apitypes.UUIDVersion1(),
			Created:    &oldTime, // old record
			Request: &apitypes.TransactionRequest{
				TransactionInput: ffcapi.TransactionInput{
					TransactionHeaders: ffcapi.TransactionHeaders{
						From: "0x12345",
					},
				},
			},
		}
		err = m.persistence.WriteTransaction(m.ctx, ln.spent, true)
		assert.NoError(t, err)
		ln.complete(context.Background())
	}()

	go func() {
		defer close(done2)

		<-locked1
		ln, err := m.assignAndLockNonce(context.Background(), "ns2:"+fftypes.NewUUID().String(), "0x12345")
		assert.NoError(t, err)

		assert.Equal(t, uint64(1112), ln.nonce)

		ln.complete(context.Background())

	}()

	<-done1
	<-done2

}

func TestNonceError(t *testing.T) {

	_, m, cancel := newTestManager(t)
	defer cancel()

	mFFC := m.connector.(*ffcapimocks.API)

	mFFC.On("NextNonceForSigner", mock.Anything, mock.Anything).Return(nil, ffcapi.ErrorReason(""), fmt.Errorf("pop"))

	_, err := m.sendManagedTransaction(context.Background(), &apitypes.TransactionRequest{
		TransactionInput: ffcapi.TransactionInput{
			TransactionHeaders: ffcapi.TransactionHeaders{
				From: "0x12345",
			},
		},
	})
	assert.Regexp(t, "pop", err)

	m.mux.Lock()
	locked, isLocked := m.lockedNonces["0x12345"]
	assert.Nil(t, locked)
	assert.False(t, isLocked)
	m.mux.Unlock()

}
