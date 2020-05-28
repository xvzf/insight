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

package common

// From github.com/google/gopacket/layers/
const (
	ICMPv4TypeEchoReply           = 0
	ICMPv4TypeEchoRequest         = 8
	ICMPv4TypeRouterAdvertisement = 9
	ICMPv4TypeRouterSolicitation  = 10
	ICMPv4TypeTimestampRequest    = 13
	ICMPv4TypeTimestampReply      = 14
	ICMPv4TypeInfoRequest         = 15
	ICMPv4TypeInfoReply           = 16
	ICMPv4TypeAddressMaskRequest  = 17
	ICMPv4TypeAddressMaskReply    = 18

	// The following are from RFC 4443
	ICMPv6TypeDestinationUnreachable = 1
	ICMPv6TypePacketTooBig           = 2
	ICMPv6TypeTimeExceeded           = 3
	ICMPv6TypeParameterProblem       = 4
	ICMPv6TypeEchoRequest            = 128
	ICMPv6TypeEchoReply              = 129

	// The following are from RFC 4861
	ICMPv6TypeRouterSolicitation    = 133
	ICMPv6TypeRouterAdvertisement   = 134
	ICMPv6TypeNeighborSolicitation  = 135
	ICMPv6TypeNeighborAdvertisement = 136
	ICMPv6TypeRedirect              = 137

	// The following are from RFC 2710
	ICMPv6TypeMLDv1MulticastListenerQueryMessage  = 130
	ICMPv6TypeMLDv1MulticastListenerReportMessage = 131
	ICMPv6TypeMLDv1MulticastListenerDoneMessage   = 132

	// The following are from RFC 3810
	ICMPv6TypeMLDv2MulticastListenerReportMessageV2 = 143
)

var (
	// Used for mapping ICMP (IPv4) requests and replies in order to consolidate them in one flow
	icmpV4EquivReply = map[uint16]uint16{
		ICMPv4TypeEchoReply:          ICMPv4TypeEchoRequest,
		ICMPv4TypeTimestampReply:     ICMPv4TypeTimestampRequest,
		ICMPv4TypeInfoReply:          ICMPv4TypeInfoRequest,
		ICMPv4TypeRouterSolicitation: ICMPv4TypeRouterAdvertisement,
		ICMPv4TypeAddressMaskReply:   ICMPv4TypeAddressMaskRequest,
	}
	icmpV4EquivRequest = map[uint16]uint16{
		ICMPv4TypeEchoRequest:         ICMPv4TypeEchoReply,
		ICMPv4TypeTimestampRequest:    ICMPv4TypeTimestampReply,
		ICMPv4TypeInfoRequest:         ICMPv4TypeInfoReply,
		ICMPv4TypeRouterAdvertisement: ICMPv4TypeRouterSolicitation,
		ICMPv4TypeAddressMaskRequest:  ICMPv4TypeAddressMaskReply,
	}
	icmpV4Equiv map[uint16]uint16 // filled at module initialization

	// Same for ICMPv6 (IPv6)

	icmpV6EquivReply = map[uint16]uint16{
		ICMPv6TypeEchoReply:                           ICMPv6TypeEchoRequest,
		ICMPv6TypeRouterSolicitation:                  ICMPv6TypeRouterAdvertisement,
		ICMPv6TypeNeighborSolicitation:                ICMPv6TypeNeighborAdvertisement,
		ICMPv6TypeMLDv1MulticastListenerReportMessage: ICMPv6TypeMLDv1MulticastListenerQueryMessage,
	}
	icmpV6EquivRequest = map[uint16]uint16{
		ICMPv6TypeEchoRequest:                        ICMPv6TypeEchoReply,
		ICMPv6TypeRouterAdvertisement:                ICMPv6TypeRouterSolicitation,
		ICMPv6TypeNeighborAdvertisement:              ICMPv6TypeNeighborSolicitation,
		ICMPv6TypeMLDv1MulticastListenerQueryMessage: ICMPv6TypeMLDv1MulticastListenerReportMessage,
	}
	icmpV6Equiv map[uint16]uint16 // filled at module initialization
)

func init() {
	icmpV4Equiv = make(map[uint16]uint16)
	icmpV6Equiv = make(map[uint16]uint16)

	// Initialize icmp equivalent maps
	for k, v := range icmpV4EquivReply {
		icmpV4Equiv[k] = v
	}
	for k, v := range icmpV4EquivRequest {
		icmpV4Equiv[k] = v
	}
	for k, v := range icmpV6EquivReply {
		icmpV6Equiv[k] = v
	}
	for k, v := range icmpV6EquivRequest {
		icmpV6Equiv[k] = v
	}
}

// GetICMPv4PortEquivalents Returns ICMPv4 port equivalents for hashing as well as if the
// ICMP flow is just going one-direction
func GetICMPv4PortEquivalents(t, c uint16) (uint16, uint16, bool) {
	if v, ok := icmpV4Equiv[t]; ok {
		return t, v, false // Is not one-way
	}
	return t, c, true // Is one-way
}

// GetICMPv6PortEquivalents Returns ICMPv6 port equivalents for hashing as well as if the
// ICMP flow is just going one-direction
func GetICMPv6PortEquivalents(t, c uint16) (uint16, uint16, bool) {
	if v, ok := icmpV6Equiv[t]; ok {
		return t, v, false // Is not one-way
	}
	return t, c, true // Is one-way
}

// GetICMPv4RequestType provides an ICMPv4 request type for a given ICMPv4 response type
func GetICMPv4RequestType(t uint16) (uint16, bool) {
	v, ok := icmpV4EquivReply[t]
	return v, ok
}

// GetICMPv6RequestType provides an ICMPv4 request type for a given ICMPv4 response type
func GetICMPv6RequestType(t uint16) (uint16, bool) {
	v, ok := icmpV6EquivReply[t]
	return v, ok
}
