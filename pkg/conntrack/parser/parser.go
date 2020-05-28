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

package parser

import (
	"errors"
	"net"
	"regexp"
	"strconv"

	"github.com/xvzf/insight/pkg/conntrack"
	"github.com/xvzf/insight/pkg/flow"
	"github.com/xvzf/insight/pkg/protos"
)

const ipPortRegex = `src=(.*) dst=(.*) sport=(\d+) dport=(\d+)`

var extractor = regexp.MustCompile(`\s*\[(NEW|UPDATE|DESTROY)\]\s+(tcp|udp).*` + ipPortRegex + `.*` + ipPortRegex + `.*`)

// ParseConntrackEvent parses an event comming from conntrack -E
func ParseConntrackEvent(l string) (*conntrack.Event, error) {
	m := extractor.FindStringSubmatch(l)

	// Matched :-)
	if len(m) != 0 {

		var eType uint8
		var transport protos.ProtocolType

		switch m[1] {
		case "NEW":
			eType = conntrack.EventNew
		case "UPDATE":
			eType = conntrack.EventUpdate
		case "DESTROY":
			eType = conntrack.EventDestroy
		}

		switch m[2] {
		case "tcp":
			transport = protos.TCP
		case "udp":
			transport = protos.UDP
		}

		// We do not need error checking here; Regex makes sure we only get digits
		srcPort0, _ := strconv.Atoi(m[5])
		dstPort0, _ := strconv.Atoi(m[6])
		srcPort1, _ := strconv.Atoi(m[9])
		dstPort1, _ := strconv.Atoi(m[10])

		return &conntrack.Event{
			Type: eType,
			Entry: conntrack.Entry{
				FlowMeta0: flow.Meta{
					Transport: transport,
					Src:       net.ParseIP(m[3]),
					Dst:       net.ParseIP(m[4]),
					SrcPort:   uint16(srcPort0),
					DstPort:   uint16(dstPort0),
				},
				FlowMeta1: flow.Meta{
					Transport: transport,
					Src:       net.ParseIP(m[7]),
					Dst:       net.ParseIP(m[8]),
					SrcPort:   uint16(srcPort1),
					DstPort:   uint16(dstPort1),
				},
			},
		}, nil
	}

	return nil, errors.New("Failed to extract conntrack event")
}
