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

package policyengine

import (
	"github.com/hyperledger/firefly-common/pkg/fftypes"
	"github.com/hyperledger/firefly-transaction-manager/internal/confirmations"
	"github.com/hyperledger/firefly-transaction-manager/pkg/apitypes"
	"github.com/hyperledger/firefly-transaction-manager/pkg/ffcapi"
)

type ManagedTXError struct {
	Time   *fftypes.FFTime    `json:"time"`
	Error  string             `json:"error,omitempty"`
	Mapped ffcapi.ErrorReason `json:"mapped,omitempty"`
}

// ManagedTXOutput is the structure stored into the operation in FireFly, that the policy
// engine can use to apply policy, and apply updates to
type ManagedTXOutput struct {
	FFTMName        string                             `json:"fftmName"`
	ID              string                             `json:"id"`
	Nonce           *fftypes.FFBigInt                  `json:"nonce"`
	Gas             *fftypes.FFBigInt                  `json:"gas"`
	TransactionHash string                             `json:"transactionHash,omitempty"`
	TransactionData string                             `json:"transactionData,omitempty"`
	GasPrice        *fftypes.JSONAny                   `json:"gasPrice"`
	PolicyInfo      *fftypes.JSONAny                   `json:"policyInfo"`
	FirstSubmit     *fftypes.FFTime                    `json:"firstSubmit,omitempty"`
	LastSubmit      *fftypes.FFTime                    `json:"lastSubmit,omitempty"`
	Request         *apitypes.TransactionRequest       `json:"request,omitempty"`
	Receipt         *ffcapi.TransactionReceiptResponse `json:"receipt,omitempty"`
	ErrorHistory    []*ManagedTXError                  `json:"errorHistory"`
	Confirmations   []confirmations.BlockInfo          `json:"confirmations,omitempty"`
}
