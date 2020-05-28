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

package capture

import (
	"errors"
	"fmt"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/xvzf/insight/pkg/flow"
	"github.com/xvzf/insight/pkg/protos"
)

// Sample contains flow data (IP based) annotated with TCP and UDP metadata
type Sample struct {
	Transport protos.ProtocolType // Protocol Type
	IcmpType  uint16              // ICMP Type (in case of protocol = ICMPv4 or ICMPv6)
	IcmpCode  uint16              // ICMP Code (in case of protocol = ICMPv4 or ICMPv6)
	SrcPort   uint16              // Source port (in case of protocol = UDP or TCP)
	DstPort   uint16              // Destination port (in case of protocol = UDP or TCP)
	Src       net.IP              // Source IP
	Dst       net.IP              // Destination IP
	Bytes     uint16              // Packet size
}

// FlowMeta extracts flow metadata from a packet
func (s Sample) FlowMeta() flow.Meta {
	return flow.Meta{
		Transport: s.Transport,
		Src:       s.Src,
		Dst:       s.Dst,
		DstPort:   s.DstPort,
		SrcPort:   s.SrcPort,
		IcmpType:  s.IcmpType,
		IcmpCode:  s.IcmpCode,
	}
}

func (s Sample) String() string {
	return fmt.Sprintf(
		"PROTOCOL: %s [%s]:%d -> [%s]:%d",
		s.Transport,
		s.Src,
		s.SrcPort,
		s.Dst,
		s.DstPort,
	)
}

// getTransportMeta tries to extract TCP/UDP Metadata (srcport, dstport) from a given packet
func getMeta(p gopacket.Packet) (uint16, uint16, uint16, uint16, protos.ProtocolType) {

	if tcpLayer := p.Layer(layers.LayerTypeTCP); tcpLayer != nil {
		tcp, _ := tcpLayer.(*layers.TCP)
		return uint16(tcp.SrcPort), uint16(tcp.DstPort), 0, 0, protos.TCP
	}

	if udpLayer := p.Layer(layers.LayerTypeUDP); udpLayer != nil {
		udp, _ := udpLayer.(*layers.UDP)
		return uint16(udp.SrcPort), uint16(udp.DstPort), 0, 0, protos.UDP
	}

	if icmp4Layer := p.Layer(layers.LayerTypeICMPv4); icmp4Layer != nil {
		icmp, _ := icmp4Layer.(*layers.ICMPv4)
		return 0, 0, uint16(icmp.TypeCode.Type()), uint16(icmp.TypeCode.Code()), protos.ICMP4
	}

	if icmp6Layer := p.Layer(layers.LayerTypeICMPv6); icmp6Layer != nil {
		icmp, _ := icmp6Layer.(*layers.ICMPv6)
		return 0, 0, uint16(icmp.TypeCode.Type()), uint16(icmp.TypeCode.Code()), protos.ICMP6
	}

	return 0, 0, 0, 0, 0
}

// NewSample parses an IP packet and generates a Sample
func NewSample(p gopacket.Packet) (*Sample, error) {

	// Parse IPv4 packet
	if v4layer := p.Layer(layers.LayerTypeIPv4); v4layer != nil {
		v4, _ := v4layer.(*layers.IPv4)
		srcPort, dstPort, icmpType, icmpCode, proto := getMeta(p)
		return &Sample{
			Src:       v4.SrcIP,
			Dst:       v4.DstIP,
			IcmpType:  icmpType,
			IcmpCode:  icmpCode,
			Transport: proto,
			Bytes:     v4.Length,
			SrcPort:   srcPort,
			DstPort:   dstPort,
		}, nil
	}

	// Parse IPv6 packet
	if v6layer := p.Layer(layers.LayerTypeIPv6); v6layer != nil {
		v6, _ := v6layer.(*layers.IPv6)
		srcPort, dstPort, icmpType, icmpCode, proto := getMeta(p)
		return &Sample{
			Src:       v6.SrcIP,
			Dst:       v6.DstIP,
			IcmpType:  icmpType,
			IcmpCode:  icmpCode,
			Transport: proto,
			Bytes:     v6.Length,
			SrcPort:   srcPort,
			DstPort:   dstPort,
		}, nil
	}

	// No IP packet/invalid header
	return nil, errors.New("not an IP packet")
}
