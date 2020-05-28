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

import "testing"

import "github.com/xvzf/insight/pkg/protos"

import "net"

import "github.com/google/go-cmp/cmp"

import "github.com/xvzf/insight/pkg/flow/common"

type testPairFlowMeta struct {
	raw    Meta // raw Flow metadata
	golden Meta // expected after correction
}

func TestWithCorrectedSourceICMPv6(t *testing.T) {
	for _, toTest := range []testPairFlowMeta{
		// Correct source and destination
		testPairFlowMeta{
			raw:    Meta{Src: net.ParseIP("10.0.0.1"), Dst: net.ParseIP("10.0.0.2"), IcmpType: common.ICMPv6TypeEchoRequest},
			golden: Meta{Src: net.ParseIP("10.0.0.1"), Dst: net.ParseIP("10.0.0.2"), IcmpType: common.ICMPv6TypeEchoRequest},
		},
		testPairFlowMeta{
			raw:    Meta{Src: net.ParseIP("10.0.0.1"), Dst: net.ParseIP("10.0.0.2"), IcmpType: common.ICMPv6TypeRouterAdvertisement},
			golden: Meta{Src: net.ParseIP("10.0.0.1"), Dst: net.ParseIP("10.0.0.2"), IcmpType: common.ICMPv6TypeRouterAdvertisement},
		},
		testPairFlowMeta{
			raw:    Meta{Src: net.ParseIP("10.0.0.1"), Dst: net.ParseIP("10.0.0.2"), IcmpType: common.ICMPv6TypeNeighborAdvertisement},
			golden: Meta{Src: net.ParseIP("10.0.0.1"), Dst: net.ParseIP("10.0.0.2"), IcmpType: common.ICMPv6TypeNeighborAdvertisement},
		},
		testPairFlowMeta{
			raw:    Meta{Src: net.ParseIP("10.0.0.1"), Dst: net.ParseIP("10.0.0.2"), IcmpType: common.ICMPv6TypeMLDv1MulticastListenerQueryMessage},
			golden: Meta{Src: net.ParseIP("10.0.0.1"), Dst: net.ParseIP("10.0.0.2"), IcmpType: common.ICMPv6TypeMLDv1MulticastListenerQueryMessage},
		},
		// Flipped source and destination
		testPairFlowMeta{
			raw:    Meta{Src: net.ParseIP("10.0.0.1"), Dst: net.ParseIP("10.0.0.2"), IcmpType: common.ICMPv6TypeEchoReply},
			golden: Meta{Src: net.ParseIP("10.0.0.2"), Dst: net.ParseIP("10.0.0.1"), IcmpType: common.ICMPv6TypeEchoRequest},
		},
		testPairFlowMeta{
			raw:    Meta{Src: net.ParseIP("10.0.0.1"), Dst: net.ParseIP("10.0.0.2"), IcmpType: common.ICMPv6TypeRouterSolicitation},
			golden: Meta{Src: net.ParseIP("10.0.0.2"), Dst: net.ParseIP("10.0.0.1"), IcmpType: common.ICMPv6TypeRouterAdvertisement},
		},
		testPairFlowMeta{
			raw:    Meta{Src: net.ParseIP("10.0.0.1"), Dst: net.ParseIP("10.0.0.2"), IcmpType: common.ICMPv6TypeNeighborSolicitation},
			golden: Meta{Src: net.ParseIP("10.0.0.2"), Dst: net.ParseIP("10.0.0.1"), IcmpType: common.ICMPv6TypeNeighborAdvertisement},
		},
		testPairFlowMeta{
			raw:    Meta{Src: net.ParseIP("10.0.0.1"), Dst: net.ParseIP("10.0.0.2"), IcmpType: common.ICMPv6TypeMLDv1MulticastListenerReportMessage},
			golden: Meta{Src: net.ParseIP("10.0.0.2"), Dst: net.ParseIP("10.0.0.1"), IcmpType: common.ICMPv6TypeMLDv1MulticastListenerQueryMessage},
		},
		// Should not be touched as it is not a "connection tracking" type
		testPairFlowMeta{
			raw:    Meta{Src: net.ParseIP("10.0.0.1"), Dst: net.ParseIP("10.0.0.2"), IcmpType: 1337},
			golden: Meta{Src: net.ParseIP("10.0.0.1"), Dst: net.ParseIP("10.0.0.2"), IcmpType: 1337},
		},
	} {
		// Set ICMPv4 Protocol
		toTest.raw.Transport, toTest.golden.Transport = protos.ICMP6, protos.ICMP6

		// Test if the source detection was successful
		if !cmp.Equal(toTest.raw.WithCorrectedSource(), toTest.golden) {
			t.Errorf(
				"[%s] expected  [%s]:%d -> [%s]:%d ,got: [%s]:%d -> [%s]:%d",
				toTest.golden.Transport,
				toTest.golden.Src, toTest.golden.IcmpType, toTest.golden.Dst, toTest.golden.IcmpCode,
				toTest.raw.Src, toTest.raw.IcmpType, toTest.raw.Dst, toTest.raw.IcmpCode,
			)
		}
	}
}
func TestWithCorrectedSourceICMPv4(t *testing.T) {
	for _, toTest := range []testPairFlowMeta{
		// Correct source and destination
		testPairFlowMeta{
			raw:    Meta{Src: net.ParseIP("10.0.0.1"), Dst: net.ParseIP("10.0.0.2"), IcmpType: common.ICMPv4TypeAddressMaskRequest},
			golden: Meta{Src: net.ParseIP("10.0.0.1"), Dst: net.ParseIP("10.0.0.2"), IcmpType: common.ICMPv4TypeAddressMaskRequest},
		},
		testPairFlowMeta{
			raw:    Meta{Src: net.ParseIP("10.0.0.1"), Dst: net.ParseIP("10.0.0.2"), IcmpType: common.ICMPv4TypeInfoRequest},
			golden: Meta{Src: net.ParseIP("10.0.0.1"), Dst: net.ParseIP("10.0.0.2"), IcmpType: common.ICMPv4TypeInfoRequest},
		},
		testPairFlowMeta{
			raw:    Meta{Src: net.ParseIP("10.0.0.1"), Dst: net.ParseIP("10.0.0.2"), IcmpType: common.ICMPv4TypeTimestampRequest},
			golden: Meta{Src: net.ParseIP("10.0.0.1"), Dst: net.ParseIP("10.0.0.2"), IcmpType: common.ICMPv4TypeTimestampRequest},
		},
		testPairFlowMeta{
			raw:    Meta{Src: net.ParseIP("10.0.0.1"), Dst: net.ParseIP("10.0.0.2"), IcmpType: common.ICMPv4TypeEchoRequest},
			golden: Meta{Src: net.ParseIP("10.0.0.1"), Dst: net.ParseIP("10.0.0.2"), IcmpType: common.ICMPv4TypeEchoRequest},
		},
		testPairFlowMeta{
			raw:    Meta{Src: net.ParseIP("10.0.0.1"), Dst: net.ParseIP("10.0.0.2"), IcmpType: common.ICMPv4TypeRouterAdvertisement},
			golden: Meta{Src: net.ParseIP("10.0.0.1"), Dst: net.ParseIP("10.0.0.2"), IcmpType: common.ICMPv4TypeRouterAdvertisement},
		},
		// Flipped source and destination
		testPairFlowMeta{
			raw:    Meta{Src: net.ParseIP("10.0.0.1"), Dst: net.ParseIP("10.0.0.2"), IcmpType: common.ICMPv4TypeAddressMaskReply},
			golden: Meta{Src: net.ParseIP("10.0.0.2"), Dst: net.ParseIP("10.0.0.1"), IcmpType: common.ICMPv4TypeAddressMaskRequest},
		},
		testPairFlowMeta{
			raw:    Meta{Src: net.ParseIP("10.0.0.1"), Dst: net.ParseIP("10.0.0.2"), IcmpType: common.ICMPv4TypeInfoReply},
			golden: Meta{Src: net.ParseIP("10.0.0.2"), Dst: net.ParseIP("10.0.0.1"), IcmpType: common.ICMPv4TypeInfoRequest},
		},
		testPairFlowMeta{
			raw:    Meta{Src: net.ParseIP("10.0.0.1"), Dst: net.ParseIP("10.0.0.2"), IcmpType: common.ICMPv4TypeTimestampReply},
			golden: Meta{Src: net.ParseIP("10.0.0.2"), Dst: net.ParseIP("10.0.0.1"), IcmpType: common.ICMPv4TypeTimestampRequest},
		},
		testPairFlowMeta{
			raw:    Meta{Src: net.ParseIP("10.0.0.1"), Dst: net.ParseIP("10.0.0.2"), IcmpType: common.ICMPv4TypeEchoReply},
			golden: Meta{Src: net.ParseIP("10.0.0.2"), Dst: net.ParseIP("10.0.0.1"), IcmpType: common.ICMPv4TypeEchoRequest},
		},
		testPairFlowMeta{
			raw:    Meta{Src: net.ParseIP("10.0.0.1"), Dst: net.ParseIP("10.0.0.2"), IcmpType: common.ICMPv4TypeRouterSolicitation},
			golden: Meta{Src: net.ParseIP("10.0.0.2"), Dst: net.ParseIP("10.0.0.1"), IcmpType: common.ICMPv4TypeRouterAdvertisement},
		},
		// Should not be touched as it is not a "connection tracking" type
		testPairFlowMeta{
			raw:    Meta{Src: net.ParseIP("10.0.0.1"), Dst: net.ParseIP("10.0.0.2"), IcmpType: 1337},
			golden: Meta{Src: net.ParseIP("10.0.0.1"), Dst: net.ParseIP("10.0.0.2"), IcmpType: 1337},
		},
	} {
		// Set ICMPv4 Protocol
		toTest.raw.Transport, toTest.golden.Transport = protos.ICMP4, protos.ICMP4

		// Test if the source detection was successful
		if !cmp.Equal(toTest.raw.WithCorrectedSource(), toTest.golden) {
			t.Errorf(
				"[%s] expected  [%s]:%d -> [%s]:%d ,got: [%s]:%d -> [%s]:%d",
				toTest.golden.Transport,
				toTest.golden.Src, toTest.golden.IcmpType, toTest.golden.Dst, toTest.golden.IcmpCode,
				toTest.raw.Src, toTest.raw.IcmpType, toTest.raw.Dst, toTest.raw.IcmpCode,
			)
		}
	}
}

func TestWithCorrectedSourceTCPUDP(t *testing.T) {
	for _, toTest := range []testPairFlowMeta{
		// Correct source & destination
		testPairFlowMeta{
			raw:    Meta{Src: net.ParseIP("10.0.0.1"), Dst: net.ParseIP("10.0.0.2"), SrcPort: 50124, DstPort: 443},
			golden: Meta{Src: net.ParseIP("10.0.0.1"), Dst: net.ParseIP("10.0.0.2"), SrcPort: 50124, DstPort: 443},
		},
		// Dst and Source swapped
		testPairFlowMeta{
			raw:    Meta{Src: net.ParseIP("10.0.0.1"), Dst: net.ParseIP("10.0.0.2"), SrcPort: 443, DstPort: 50124},
			golden: Meta{Src: net.ParseIP("10.0.0.2"), Dst: net.ParseIP("10.0.0.1"), SrcPort: 50124, DstPort: 443},
		},
		// Dst and Source (Ports) are equal
		testPairFlowMeta{
			raw:    Meta{Src: net.ParseIP("10.0.0.1"), Dst: net.ParseIP("10.0.0.2"), SrcPort: 1337, DstPort: 1337},
			golden: Meta{Src: net.ParseIP("10.0.0.1"), Dst: net.ParseIP("10.0.0.2"), SrcPort: 1337, DstPort: 1337},
		},
	} {
		// TCP and UDP is handled equally
		for _, p := range []protos.ProtocolType{protos.UDP, protos.TCP} {
			// Set protocol type
			toTest.raw.Transport, toTest.golden.Transport = p, p

			// Test if the source detection was successful
			if !cmp.Equal(toTest.raw.WithCorrectedSource(), toTest.golden) {
				t.Errorf(
					"[%s] expected  [%s]:%d -> [%s]:%d ,got: [%s]:%d -> [%s]:%d",
					toTest.golden.Transport,
					toTest.golden.Src, toTest.golden.SrcPort, toTest.golden.Dst, toTest.golden.DstPort,
					toTest.raw.Src, toTest.raw.SrcPort, toTest.raw.Dst, toTest.raw.DstPort,
				)
			}
		}
	}
}

func TestNewFlow(t *testing.T) {
	fm := Meta{Src: net.ParseIP("10.0.0.1"), Dst: net.ParseIP("10.0.0.2"), IcmpType: common.ICMPv4TypeRouterSolicitation}

	f := New(fm)

	if !cmp.Equal(f.Meta, fm) {
		t.Error(cmp.Diff(fm, f.Meta))
	}
}
