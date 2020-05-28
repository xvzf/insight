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
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"

	"github.com/sirupsen/logrus"
)

func init() {
	log = logrus.WithFields(logrus.Fields{
		"context": "capture",
		"mode":    "pcap",
	})
}

// Logger
var log *logrus.Entry

// Capturer is the main interface for capturing packets
type Capturer interface {
	Packets() chan gopacket.Packet
	Filter(string) error
	Close()
}

// pcapHandle contains the handle in order to allow filter changes and graceful shutdown
type pcapHandle struct {
	handle *pcap.Handle
}

// Open creates a packet capture based on pcapHandle with the jumbo frames enabled
// Filtering is done by pcapHandle.
func Open(device string) (Capturer, error) {
	// Snapshot length is 9038, @TODO
	handle, err := pcap.OpenLive(device, 9038, true, pcap.BlockForever)

	if err != nil {
		log.Error("Could not open device ", device)
		return nil, err
	}
	log.Info("Opened device ", device)

	return &pcapHandle{
		handle: handle,
	}, nil
}

// Closes the capture
func (p *pcapHandle) Close() {
	defer p.handle.Close()
}

// Set a pcap filter
func (p *pcapHandle) Filter(filter string) error {
	res := p.handle.SetBPFFilter(filter)

	if res != nil {
		log.Error("Could not set BPF filter ", filter)
	} else {
		log.Info("BPF filter set ", filter)
	}

	return res
}

// Packets returns a channel producing every packet
func (p *pcapHandle) Packets() chan gopacket.Packet {
	log.Info("Starting packet stream")
	source := gopacket.NewPacketSource(p.handle, p.handle.LinkType())
	// Improve performance
	source.Lazy = true
	source.NoCopy = true
	return source.Packets()
}
