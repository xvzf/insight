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

package communityid

import (
	"net"
	"testing"

	"github.com/xvzf/insight/pkg/flow"
	"github.com/xvzf/insight/pkg/flow/common"
	"github.com/xvzf/insight/pkg/protos"
)

type testPair struct {
	communityID string
	flow        flow.Meta
}

func TestFlowHashing(t *testing.T) {
	chash := NewHasher(0)
	for _, sample := range []testPair{
		// TCP Test Pairs
		testPair{
			"1:LQU9qZlK+B5F3KDmev6m5PMibrg=",
			flow.Meta{
				Transport: protos.TCP,
				Src:       net.ParseIP("128.232.110.120"),
				Dst:       net.ParseIP("66.35.250.204"),
				SrcPort:   34855,
				DstPort:   80,
			},
		},
		testPair{
			"1:LQU9qZlK+B5F3KDmev6m5PMibrg=",
			flow.Meta{
				Transport: protos.TCP,
				Src:       net.ParseIP("66.35.250.204"),
				Dst:       net.ParseIP("128.232.110.120"),
				SrcPort:   80,
				DstPort:   34855,
			},
		},
		// UDP TestPair
		testPair{
			"1:d/FP5EW3wiY1vCndhwleRRKHowQ=",
			flow.Meta{
				Transport: protos.UDP,
				Src:       net.ParseIP("192.168.1.52"),
				Dst:       net.ParseIP("8.8.8.8"),
				SrcPort:   54585,
				DstPort:   53,
			},
		},
		testPair{
			"1:d/FP5EW3wiY1vCndhwleRRKHowQ=",
			flow.Meta{
				Transport: protos.UDP,
				Src:       net.ParseIP("8.8.8.8"),
				Dst:       net.ParseIP("192.168.1.52"),
				SrcPort:   53,
				DstPort:   54585,
			},
		},
		// ICMPv4 TestPair
		testPair{
			"1:X0snYXpgwiv9TZtqg64sgzUn6Dk=",
			flow.Meta{
				Transport: protos.ICMP4,
				Src:       net.ParseIP("192.168.0.89"),
				Dst:       net.ParseIP("192.168.0.1"),
				IcmpType:  common.ICMPv4TypeEchoRequest,
				IcmpCode:  123,
			},
		},
		testPair{
			"1:X0snYXpgwiv9TZtqg64sgzUn6Dk=",
			flow.Meta{
				Transport: protos.ICMP4,
				Src:       net.ParseIP("192.168.0.1"),
				Dst:       net.ParseIP("192.168.0.89"),
				IcmpType:  common.ICMPv4TypeEchoReply,
				IcmpCode:  111,
			},
		},
		testPair{
			"1:X0snYXpgwiv9TZtqg64sgzUn6Dk=",
			flow.Meta{
				Transport: protos.ICMP4,
				Src:       net.ParseIP("192.168.0.1"),
				Dst:       net.ParseIP("192.168.0.89"),
				IcmpType:  0,
				IcmpCode:  8,
			},
		},
		testPair{
			"1:dGHyGvjMfljg6Bppwm3bg0LO8TY=",
			flow.Meta{
				Transport: protos.ICMP6,
				Src:       net.ParseIP("fe80::200:86ff:fe05:80da"),
				Dst:       net.ParseIP("fe80::260:97ff:fe07:69ea"),
				IcmpType:  135,
				IcmpCode:  0,
			},
		},
		testPair{
			"1:dGHyGvjMfljg6Bppwm3bg0LO8TY=",
			flow.Meta{
				Transport: protos.ICMP6,
				Src:       net.ParseIP("fe80::260:97ff:fe07:69ea"),
				Dst:       net.ParseIP("fe80::200:86ff:fe05:80da"),
				IcmpType:  136,
				IcmpCode:  0,
			},
		},
	} {
		if computed := chash.Hash(sample.flow); sample.communityID != computed {
			t.Errorf("Hasher generated different community-ID; %s != %s", sample.communityID, computed)
		}
	}
}
