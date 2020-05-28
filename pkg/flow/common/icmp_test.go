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

import "testing"

func TestGetICMPv4PortEquivalents(t *testing.T) {
	for mapFrom, mapTo := range icmpV4Equiv {
		p0, p1, oneWay := GetICMPv4PortEquivalents(mapFrom, 0)

		if !oneWay && p0 != mapFrom && p1 != mapTo {
			t.Errorf(
				"Expected (%d, %d, %t) got (%d, %d, %t)",
				mapFrom, mapTo, false,
				p0, p1, oneWay,
			)
		}
	}
}

func TestGetICMPv4PortNoEquivalents(t *testing.T) {
	for k, v := range map[uint16]uint16{
		1337: 1337,
		1338: 1338,
		1234: 1234,
	} {
		p0, p1, oneWay := GetICMPv4PortEquivalents(k, v)

		if oneWay && p0 != k && p1 != v {
			t.Errorf(
				"Expected (%d, %d, %t) got (%d, %d, %t)",
				k, v, false,
				p0, p1, oneWay,
			)
		}
	}
}

func TestGetICMPv6PortEquivalents(t *testing.T) {
	for mapFrom, mapTo := range icmpV6Equiv {
		p0, p1, oneWay := GetICMPv6PortEquivalents(mapFrom, 0)

		if !oneWay && p0 != mapFrom && p1 != mapTo {
			t.Errorf(
				"Expected (%d, %d, %t) got (%d, %d, %t)",
				mapFrom, mapTo, false,
				p0, p1, oneWay,
			)
		}
	}
}

func TestGetICMPv6PortNoEquivalents(t *testing.T) {
	for k, v := range map[uint16]uint16{
		1337: 1337,
		1338: 1338,
		1234: 1234,
	} {
		p0, p1, oneWay := GetICMPv6PortEquivalents(k, v)

		if oneWay && p0 != k && p1 != v {
			t.Errorf(
				"Expected (%d, %d, %t) got (%d, %d, %t)",
				k, v, false,
				p0, p1, oneWay,
			)
		}
	}
}

func TestGetICMPv4RequestType(t *testing.T) {
	for k, v := range map[uint16]uint16{
		ICMPv4TypeEchoReply:        ICMPv4TypeEchoRequest,
		ICMPv4TypeAddressMaskReply: ICMPv4TypeAddressMaskRequest,
	} {
		res, ok := GetICMPv4RequestType(k)
		if !ok {
			t.Errorf("%d should have been mapped to a request type", k)
		} else if res != v {
			t.Errorf("Expected %d, got %d", v, res)
		}
	}
}

func TestGetICMPv6RequestType(t *testing.T) {
	for k, v := range map[uint16]uint16{
		ICMPv6TypeEchoReply: ICMPv6TypeEchoRequest,
	} {
		res, ok := GetICMPv6RequestType(k)
		if !ok {
			t.Errorf("%d should have been mapped to a request type", k)
		} else if res != v {
			t.Errorf("Expected %d, got %d", v, res)
		}
	}
}
