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
	"net"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/xvzf/insight/pkg/capture"
	"github.com/xvzf/insight/pkg/flow"
	"github.com/xvzf/insight/pkg/flow/common"
	"github.com/xvzf/insight/pkg/protos"
)

var (
	tcpTestSamples = []*capture.Sample{
		// TCP packets IPv4
		&capture.Sample{Transport: protos.TCP, IcmpType: 0, IcmpCode: 0, SrcPort: 123, DstPort: 345, Src: net.ParseIP("10.0.0.1"), Dst: net.ParseIP("10.0.0.2"), Bytes: 1024},
		&capture.Sample{Transport: protos.TCP, IcmpType: 0, IcmpCode: 0, SrcPort: 345, DstPort: 123, Src: net.ParseIP("10.0.0.2"), Dst: net.ParseIP("10.0.0.1"), Bytes: 2048},
		&capture.Sample{Transport: protos.TCP, IcmpType: 0, IcmpCode: 0, SrcPort: 123, DstPort: 345, Src: net.ParseIP("10.0.0.1"), Dst: net.ParseIP("10.0.0.2"), Bytes: 1024},
		// TCP packets IPv6
		&capture.Sample{Transport: protos.TCP, IcmpType: 0, IcmpCode: 0, SrcPort: 123, DstPort: 345, Src: net.ParseIP("2000:dead:beef::2345"), Dst: net.ParseIP("2000:dead:beef::1234"), Bytes: 1024},
		&capture.Sample{Transport: protos.TCP, IcmpType: 0, IcmpCode: 0, SrcPort: 345, DstPort: 123, Src: net.ParseIP("2000:dead:beef::1234"), Dst: net.ParseIP("2000:dead:beef::2345"), Bytes: 2048},
		&capture.Sample{Transport: protos.TCP, IcmpType: 0, IcmpCode: 0, SrcPort: 123, DstPort: 345, Src: net.ParseIP("2000:dead:beef::2345"), Dst: net.ParseIP("2000:dead:beef::1234"), Bytes: 1024},
	}

	tcpFlows = []*flow.Flow{
		// TCP IPv4
		&flow.Flow{
			Meta:     flow.Meta{Transport: protos.TCP, IcmpType: 0, IcmpCode: 0, SrcPort: 345, DstPort: 123, Src: net.ParseIP("10.0.0.2"), Dst: net.ParseIP("10.0.0.1")},
			Incoming: flow.Counters{Bytes: 2048, Packets: 1},
			Outgoing: flow.Counters{Bytes: 2048, Packets: 2},
		},
		// TCP IPv6
		&flow.Flow{
			Meta:     flow.Meta{Transport: protos.TCP, IcmpType: 0, IcmpCode: 0, SrcPort: 345, DstPort: 123, Src: net.ParseIP("2000:dead:beef::1234"), Dst: net.ParseIP("2000:dead:beef::2345")},
			Incoming: flow.Counters{Bytes: 2048, Packets: 1},
			Outgoing: flow.Counters{Bytes: 2048, Packets: 2},
		},
	}

	udpTestSamples = []*capture.Sample{
		// UDP packets IPv4
		&capture.Sample{Transport: protos.UDP, IcmpType: 0, IcmpCode: 0, SrcPort: 123, DstPort: 345, Src: net.ParseIP("10.0.0.1"), Dst: net.ParseIP("10.0.0.2"), Bytes: 1024},
		&capture.Sample{Transport: protos.UDP, IcmpType: 0, IcmpCode: 0, SrcPort: 345, DstPort: 123, Src: net.ParseIP("10.0.0.2"), Dst: net.ParseIP("10.0.0.1"), Bytes: 2048},
		&capture.Sample{Transport: protos.UDP, IcmpType: 0, IcmpCode: 0, SrcPort: 123, DstPort: 345, Src: net.ParseIP("10.0.0.1"), Dst: net.ParseIP("10.0.0.2"), Bytes: 1024},
		// UDP packets IPv6
		&capture.Sample{Transport: protos.UDP, IcmpType: 0, IcmpCode: 0, SrcPort: 123, DstPort: 345, Src: net.ParseIP("2000:dead:beef::2345"), Dst: net.ParseIP("2000:dead:beef::1234"), Bytes: 1024},
		&capture.Sample{Transport: protos.UDP, IcmpType: 0, IcmpCode: 0, SrcPort: 345, DstPort: 123, Src: net.ParseIP("2000:dead:beef::1234"), Dst: net.ParseIP("2000:dead:beef::2345"), Bytes: 2048},
		&capture.Sample{Transport: protos.UDP, IcmpType: 0, IcmpCode: 0, SrcPort: 123, DstPort: 345, Src: net.ParseIP("2000:dead:beef::2345"), Dst: net.ParseIP("2000:dead:beef::1234"), Bytes: 1024},
	}

	udpFlows = []*flow.Flow{
		// UDP IPv4
		&flow.Flow{
			Meta:     flow.Meta{Transport: protos.UDP, IcmpType: 0, IcmpCode: 0, SrcPort: 345, DstPort: 123, Src: net.ParseIP("10.0.0.2"), Dst: net.ParseIP("10.0.0.1")},
			Incoming: flow.Counters{Bytes: 2048, Packets: 1},
			Outgoing: flow.Counters{Bytes: 2048, Packets: 2},
		},
		// UDP IPv6
		&flow.Flow{
			Meta:     flow.Meta{Transport: protos.UDP, IcmpType: 0, IcmpCode: 0, SrcPort: 345, DstPort: 123, Src: net.ParseIP("2000:dead:beef::1234"), Dst: net.ParseIP("2000:dead:beef::2345")},
			Incoming: flow.Counters{Bytes: 2048, Packets: 1},
			Outgoing: flow.Counters{Bytes: 2048, Packets: 2},
		},
	}

	icmp4TestSamples = []*capture.Sample{
		// ICMPv4
		&capture.Sample{Transport: protos.ICMP4, IcmpType: common.ICMPv4TypeEchoRequest, IcmpCode: 0, SrcPort: 0, DstPort: 0, Src: net.ParseIP("10.0.0.1"), Dst: net.ParseIP("10.0.0.2"), Bytes: 512},
		&capture.Sample{Transport: protos.ICMP4, IcmpType: common.ICMPv4TypeEchoRequest, IcmpCode: 0, SrcPort: 0, DstPort: 0, Src: net.ParseIP("10.0.0.1"), Dst: net.ParseIP("10.0.0.2"), Bytes: 512},
		&capture.Sample{Transport: protos.ICMP4, IcmpType: common.ICMPv4TypeEchoReply, IcmpCode: 0, SrcPort: 0, DstPort: 0, Src: net.ParseIP("10.0.0.2"), Dst: net.ParseIP("10.0.0.1"), Bytes: 512},
	}

	icmp4Flows = []*flow.Flow{
		// ICMPv4
		&flow.Flow{
			Meta:     flow.Meta{Transport: protos.ICMP4, IcmpType: common.ICMPv4TypeEchoRequest, IcmpCode: 0, SrcPort: 0, DstPort: 0, Src: net.ParseIP("10.0.0.1"), Dst: net.ParseIP("10.0.0.2")},
			Incoming: flow.Counters{Bytes: 1024, Packets: 2},
			Outgoing: flow.Counters{Bytes: 512, Packets: 1},
		},
	}

	icmp6TestSamples = []*capture.Sample{
		// ICMPv6
		&capture.Sample{Transport: protos.ICMP6, IcmpType: common.ICMPv6TypeEchoReply, IcmpCode: 0, SrcPort: 0, DstPort: 0, Src: net.ParseIP("2000:dead:beef::2345"), Dst: net.ParseIP("2000:dead:beef::1234"), Bytes: 512},
		&capture.Sample{Transport: protos.ICMP6, IcmpType: common.ICMPv6TypeEchoRequest, IcmpCode: 0, SrcPort: 0, DstPort: 0, Src: net.ParseIP("2000:dead:beef::1234"), Dst: net.ParseIP("2000:dead:beef::2345"), Bytes: 512},
	}

	icmp6Flows = []*flow.Flow{
		&flow.Flow{
			Meta:     flow.Meta{Transport: protos.ICMP6, IcmpType: common.ICMPv6TypeEchoRequest, IcmpCode: 0, SrcPort: 0, DstPort: 0, Src: net.ParseIP("2000:dead:beef::1234"), Dst: net.ParseIP("2000:dead:beef::2345")},
			Incoming: flow.Counters{Bytes: 512, Packets: 1},
			Outgoing: flow.Counters{Bytes: 512, Packets: 1},
		},
	}

	testFlows   []*flow.Flow
	testSamples []*capture.Sample
)

func init() {
	// Initialize combined datastructures for testing
	testSamples = make([]*capture.Sample, 0)
	testSamples = append(testSamples, tcpTestSamples...)
	testSamples = append(testSamples, udpTestSamples...)
	testSamples = append(testSamples, icmp4TestSamples...)
	testSamples = append(testSamples, icmp6TestSamples...)

	testFlows = make([]*flow.Flow, 0)
	testFlows = append(testFlows, tcpFlows...)
	testFlows = append(testFlows, udpFlows...)
	testFlows = append(testFlows, icmp4Flows...)
	testFlows = append(testFlows, icmp6Flows...)

}

func TestNewContainer(t *testing.T) {
	c := New()
	cRaw, ok := c.(*container)
	if !ok {
		t.Error("Container not initalized")
	}

	if cRaw.start.Sub(time.Now()) > 0 {
		t.Error("Start time not set correctly")
	}

	if cRaw.data.flows == nil {
		t.Error("datastructures not initalized")
	}

}

func TestContainerAddInvalidSample(t *testing.T) {
	c := New()

	if err := c.Add(nil); err == nil {
		t.Error("Invalid error, expected error got nil")
	}
}

func TestContainerDump(t *testing.T) {
	c := New()

	for _, s := range testSamples {
		if c.Add(s) != nil {
			t.Errorf("Failed to add sample %s", s)
			return
		}
	}

	for _, f := range c.Dump() {
		// Reset timers and CommunityID so we can compare the rest
		f.Start = time.Time{}
		f.End = time.Time{}
		f.CommunityID = ""

		found := false
		for _, testFlow := range testFlows {
			if f.Start.Sub(f.End) > 0 {
				t.Error("Time ranges are invalid")
			}

			found = found || cmp.Equal(f, testFlow)
		}
		if !found {
			t.Error(f)
			t.Error("failed to find flow in expected flow results")
		}
	}
}
