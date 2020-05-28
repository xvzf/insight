/* 
 *  MIT License
 *  
 *  Copyright (c) 2020 Matthias Riegler <me@xvzf.tech>
 *  
 *  Permission is hereby granted, free of charge, to any person obtaining a copy
 *  of this software and associated documentation files (the "Software"), to deal
 *  in the Software without restriction, including without limitation the rights
 *  to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *  copies of the Software, and to permit persons to whom the Software is
 *  furnished to do so, subject to the following conditions:
 *  
 *  The above copyright notice and this permission notice shall be included in all
 *  copies or substantial portions of the Software.
 *  
 *  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *  AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *  LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *  SOFTWARE.
 *  
 */

package container

import (
	"errors"
	"sync"
	"time"

	"github.com/xvzf/insight/pkg/capture"
	"github.com/xvzf/insight/pkg/flow"
	"github.com/xvzf/insight/pkg/flow/communityid"
)

// Container takes in samples and aggregates them into flows
type Container interface {
	Add(s *capture.Sample) error
	Dump() []*flow.Flow
}

type data struct {
	sync.Mutex
	flows map[string]*flow.Flow
}

type container struct {
	data   data
	hasher communityid.Hasher
	start  time.Time
}

// New creates a new flow container
func New() Container {
	return &container{
		data: data{
			flows: make(map[string]*flow.Flow),
		},
		hasher: communityid.NewHasher(0),
		start:  time.Now(),
	}
}

// Adds a sample to the flowtable
func (c *container) Add(s *capture.Sample) error {

	// Check for nullpointers
	if s == nil {
		return errors.New("sample cannot be nil")
	}

	// Generate CommunityID for the packet
	fm := s.FlowMeta()
	cID := c.hasher.Hash(fm)

	// Lock data container
	c.data.Lock()
	defer c.data.Unlock()

	// Create a new flow if it is not already in the flowtable
	if _, ok := c.data.flows[cID]; !ok {
		c.data.flows[cID] = flow.New(fm.WithCorrectedSource())
	}

	f, ok := c.data.flows[cID]
	if !ok {
		return errors.New("internal error on mapping operation")
	}

	// Update counters
	if f.Meta.Src.Equal(s.Src) {
		// Incoming (Src -> Dst)
		f.Incoming.Packets++
		f.Incoming.Bytes += uint64(s.Bytes)
	} else {
		// Outgoing (Dst -> Src)
		f.Outgoing.Packets++
		f.Outgoing.Bytes += uint64(s.Bytes)
	}

	return nil
}

// Dump dumps all flows in the container
func (c *container) Dump() []*flow.Flow {

	// Lock data container
	c.data.Lock()
	defer c.data.Unlock()

	var buf []*flow.Flow

	// end timestamp
	end := time.Now()

	// Iterate over the hashmap and set the CommunityID attribute
	for cID, flow := range c.data.flows {
		// Update flow parameters
		flow.CommunityID = cID
		flow.Start = c.start
		flow.End = end

		// Add to result buffer
		buf = append(buf, flow)
	}

	return buf
}
