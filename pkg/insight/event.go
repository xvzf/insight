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

package insight

import (
	"net"
	"os"
	"time"

	"github.com/xvzf/insight/pkg/flow"
)

func init() {
	var err error
	hostname, err = os.Hostname()
	if err != nil {
		hostname = "undefined"
	}
}

const ECSversion = "1.4"

var hostname string

// Agent in ECS
type Agent struct {
	HostName string `json:"hostname"`
	Type     string `json:"type"`
}

// ECS version tag
type ECS struct {
	Version string `json:"version"`
}

// EventDescription in ECS
type EventDescription struct {
	Duration time.Duration `json:"duration"`
	Kind     string        `json:"kind"`
	Action   string        `json:"action"`
	Category string        `json:"category"`
	Dataset  string        `json:"dataset"`
	Start    time.Time     `json:"start"`
	End      time.Time     `json:"end"`
}

// EndpointDescription in ECS
type EndpointDescription struct {
	Address string `json:"address"`
	IP      net.IP `json:"ip"`
	Port    uint16 `json:"port"`
	Bytes   uint64 `json:"bytes"`
	Packets uint64 `json:"packets"`
}

// NetworkDescription in ECS
type NetworkDescription struct {
	Type        string `json:"type"`
	Bytes       uint64 `json:"bytes"`
	Packets     uint64 `json:"packets"`
	Transport   string `json:"transport"`
	CommunityID string `json:"community_id"`
}

// Event contains the event metadata passed to logstash
type Event struct {
	Type        string               `json:"type"`
	ECS         *ECS                 `json:"ecs"`
	Agent       *Agent               `json:"agent"`
	Event       *EventDescription    `json:"event"`
	Source      *EndpointDescription `json:"source"`
	Destination *EndpointDescription `json:"destination"`
	Network     *NetworkDescription  `json:"network"`
}

// NewFromFlows generates an event for every flow
func NewFromFlows(flows []*flow.Flow) []*Event {
	var buf []*Event

	for _, f := range flows {
		e := NewFromFlow(f)
		buf = append(buf, e)
	}

	return buf
}

// NewFromFlow generates a new event based on a flow
func NewFromFlow(f *flow.Flow) *Event {
	ipVersion := "ipv6"
	if f.Meta.Src.To4() != nil {
		ipVersion = "ipv4"
	}
	return &Event{
		Agent: &Agent{
			HostName: hostname,
			Type:     "insight",
		},
		ECS: &ECS{
			Version: ECSversion,
		},
		Event: &EventDescription{
			Duration: f.End.Sub(f.Start),
			Kind:     "event",
			Action:   "network_flow",
			Category: "network_traffic",
			Dataset:  "flow",
			Start:    f.Start,
			End:      f.End,
		},
		Source: &EndpointDescription{
			Address: f.Meta.Src.String(),
			IP:      f.Meta.Src,
			Port:    f.Meta.SrcPort,
			Bytes:   f.Incoming.Bytes,
			Packets: f.Incoming.Packets,
		},
		Destination: &EndpointDescription{
			Address: f.Meta.Dst.String(),
			IP:      f.Meta.Dst,
			Port:    f.Meta.DstPort,
			Bytes:   f.Outgoing.Bytes,
			Packets: f.Outgoing.Packets,
		},
		Network: &NetworkDescription{
			Type:        ipVersion,
			Bytes:       f.Incoming.Bytes + f.Outgoing.Bytes,
			Packets:     f.Incoming.Packets + f.Outgoing.Packets,
			Transport:   f.Meta.Transport.String(),
			CommunityID: f.CommunityID,
		},
	}
}
