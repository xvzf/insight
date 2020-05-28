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

package protos

// ProtocolType defines the type of protocol
type ProtocolType uint8

// Based on https://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml
const (
	// TCP protocol
	TCP ProtocolType = 6
	// UDP protocol
	UDP ProtocolType = 17
	// ICMP4 protocol
	ICMP4 ProtocolType = 1
	// ICMP6 protocol
	ICMP6 ProtocolType = 58
)

// Stringr converts the protocol type to a string
func (pt ProtocolType) String() string {
	switch pt {
	case TCP:
		return "tcp"
	case UDP:
		return "udp"
	case ICMP4:
		return "icmp"
	case ICMP6:
		return "ipv6-icmp"
	default:
		return "UNDEFINED"
	}
}
