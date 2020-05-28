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

package flow

import (
	"net"
	"time"

	"github.com/xvzf/insight/pkg/flow/common"
	"github.com/xvzf/insight/pkg/protos"
)

// Counters contains flow counters
type Counters struct {
	Bytes   uint64
	Packets uint64
}

// Meta contains flow metadata
type Meta struct {
	Transport protos.ProtocolType
	Src       net.IP
	Dst       net.IP
	IcmpType  uint16
	IcmpCode  uint16
	DstPort   uint16
	SrcPort   uint16
}

// Flow contains flow data
type Flow struct {
	Meta        Meta      // Flow Src/Dst & Protocol Information
	Incoming    Counters  // Incoming counters
	Outgoing    Counters  // Outgoing counters
	CommunityID string    // CommunityID
	Start       time.Time // Start time
	End         time.Time // Stop time
}

// New creates a new Flow
func New(fm Meta) *Flow {
	return &Flow{
		Meta:     fm,
		Incoming: Counters{},
		Outgoing: Counters{},
	}
}

// WithCorrectedSource is a super simple helper function for determin which one is the source IP
// and updates the values inside accordingly
func (m Meta) WithCorrectedSource() Meta {
	switch m.Transport {
	case protos.TCP, protos.UDP:
		// Assume that min(SrcPort, DstPort) is the server -> destination
		if m.DstPort > m.SrcPort {
			// Swap around, otherwise it's fine
			m.Src, m.SrcPort, m.Dst, m.DstPort = m.Dst, m.DstPort, m.Src, m.SrcPort
		}
	case protos.ICMP4:
		if v, ok := common.GetICMPv4RequestType(m.IcmpType); ok {
			m.Src, m.Dst = m.Dst, m.Src
			m.IcmpType = v
		}
	case protos.ICMP6:
		if v, ok := common.GetICMPv6RequestType(m.IcmpType); ok {
			m.Src, m.Dst = m.Dst, m.Src
			m.IcmpType = v
		}
	}

	return m
}
