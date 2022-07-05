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

package ffcapi

import (
	"context"

	"github.com/hyperledger/firefly-common/pkg/fftypes"
)

type EventStreamStartRequest struct {
	ID               *fftypes.UUID              // UUID of the stream, which we be referenced in any future add/remove listener requests
	StreamContext    context.Context            // Context that will be cancelled when the event stream needs to stop - no further events will be consumed after this, so all pushes to the stream should select on the done channel too
	EventStream      chan<- *ListenerEvent      // The event stream to push events to as they are detected, and checkpoints regularly even if there are no events - remember to select on Done as well when pushing events
	BlockListener    chan<- *BlockHashEvent     // The connector should push new blocks to every stream, marking if it's possible blocks were missed (due to reconnect). The stream guarantees to always consume from this channel, until the stream context closes.
	InitialListeners []*EventListenerAddRequest // Initial list of event listeners to start with the stream - allows these to be started concurrently
}

type EventStreamStartResponse struct {
}
