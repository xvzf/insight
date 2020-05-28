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
	"net"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/xvzf/insight/pkg/conntrack"
	"github.com/xvzf/insight/pkg/flow"
	"github.com/xvzf/insight/pkg/protos"
)

func TestParserValid(t *testing.T) {
	for line, expected := range map[string]*conntrack.Event{
		"[DESTROY] udp      17 src=10.42.2.26 dst=10.43.0.10 sport=49129 dport=53 src=10.42.1.4 dst=10.42.2.26 sport=53 dport=49129 [ASSURED]": &conntrack.Event{
			Type: conntrack.EventDestroy,
			Entry: conntrack.Entry{
				FlowMeta0: flow.Meta{
					Transport: protos.UDP, Src: net.ParseIP("10.42.2.26"), Dst: net.ParseIP("10.43.0.10"),
					SrcPort: 49129, DstPort: 53,
				},
				FlowMeta1: flow.Meta{
					Transport: protos.UDP, Src: net.ParseIP("10.42.1.4"), Dst: net.ParseIP("10.42.2.26"),
					SrcPort: 53, DstPort: 49129,
				},
			},
		},
		"[UPDATE] udp      17 src=10.42.2.26 dst=10.43.0.10 sport=49129 dport=53 src=10.42.1.4 dst=10.42.2.26 sport=53 dport=49129 [ASSURED]": &conntrack.Event{
			Type: conntrack.EventUpdate,
			Entry: conntrack.Entry{
				FlowMeta0: flow.Meta{
					Transport: protos.UDP, Src: net.ParseIP("10.42.2.26"), Dst: net.ParseIP("10.43.0.10"),
					SrcPort: 49129, DstPort: 53,
				},
				FlowMeta1: flow.Meta{
					Transport: protos.UDP, Src: net.ParseIP("10.42.1.4"), Dst: net.ParseIP("10.42.2.26"),
					SrcPort: 53, DstPort: 49129,
				},
			},
		},
		"[NEW] udp      17 src=10.42.2.26 dst=10.43.0.10 sport=49129 dport=53 src=10.42.1.4 dst=10.42.2.26 sport=53 dport=49129 [ASSURED]": &conntrack.Event{
			Type: conntrack.EventNew,
			Entry: conntrack.Entry{
				FlowMeta0: flow.Meta{
					Transport: protos.UDP, Src: net.ParseIP("10.42.2.26"), Dst: net.ParseIP("10.43.0.10"),
					SrcPort: 49129, DstPort: 53,
				},
				FlowMeta1: flow.Meta{
					Transport: protos.UDP, Src: net.ParseIP("10.42.1.4"), Dst: net.ParseIP("10.42.2.26"),
					SrcPort: 53, DstPort: 49129,
				},
			},
		},
		" [UPDATE] tcp      6 86400 ESTABLISHED src=10.42.2.28 dst=10.43.183.150 sport=53296 dport=8080 src=10.42.1.16 dst=10.42.2.28 sport=8080 dport=53296 [ASSURED]": &conntrack.Event{
			Type: conntrack.EventUpdate,
			Entry: conntrack.Entry{
				FlowMeta0: flow.Meta{
					Transport: protos.TCP, Src: net.ParseIP("10.42.2.28"), Dst: net.ParseIP("10.43.183.150"),
					SrcPort: 53296, DstPort: 8080,
				},
				FlowMeta1: flow.Meta{
					Transport: protos.TCP, Src: net.ParseIP("10.42.1.16"), Dst: net.ParseIP("10.42.2.28"),
					SrcPort: 8080, DstPort: 53296,
				},
			},
		},
	} {
		computed, err := ParseConntrackEvent(line)
		if err != nil {
			t.Error(err)
		}
		if !cmp.Equal(computed, expected) {
			t.Errorf(cmp.Diff(expected, computed))
		}
	}
}

func TestParserInvalidPorts(t *testing.T) {
	for _, line := range []string{
		"[DESTROY] udp      17 src=10.42.2.26 dst=10.43.0.10 sport=-49129 dport=53 src=10.42.1.4 dst=10.42.2.26 sport=53 dport=49129 [ASSURED]",
		"[DESTROY] udp      17 src=10.42.2.26 dst=10.43.0.10 sport=49129 dport=-53 src=10.42.1.4 dst=10.42.2.26 sport=53 dport=49129 [ASSURED]",
		"[DESTROY] udp      17 src=10.42.2.26 dst=10.43.0.10 sport=49129 dport=53 src=10.42.1.4 dst=10.42.2.26 sport=-53 dport=49129 [ASSURED]",
		"[DESTROY] udp      17 src=10.42.2.26 dst=10.43.0.10 sport=49129 dport=53 src=10.42.1.4 dst=10.42.2.26 sport=53 dport=-49129 [ASSURED]",
	} {
		_, err := ParseConntrackEvent(line)
		if err == nil {
			t.Error("Invaid port, should've failed")
		}
	}
}
