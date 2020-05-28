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
	"bytes"
	"crypto"
	_ "crypto/sha1" // Include SHA1 algorithm
	"encoding/base64"
	"encoding/binary"
	"net"

	"github.com/xvzf/insight/pkg/flow"
	"github.com/xvzf/insight/pkg/flow/common"
	"github.com/xvzf/insight/pkg/protos"
)

// Hasher computes the CommunityID of a flow based on this specification https://github.com/corelight/community-id-spec
type Hasher interface {
	Hash(f flow.Meta) string
}

type hasher struct {
	seed [2]byte
}

// NewHasher a new Hasher
func NewHasher(seed uint16) Hasher {
	h := &hasher{}
	binary.BigEndian.PutUint16(h.seed[:], seed)
	return h
}

func rawIP(ip net.IP) []byte {
	if v4 := ip.To4(); v4 != nil {
		return v4
	}
	return ip
}

func extractTuple(f flow.Meta) ([]byte, []byte, uint16, uint16) {
	// Incase we have an ICMP flow, set src and dst port according to the specification
	switch f.Transport {
	case protos.ICMP4:
		f.SrcPort, f.DstPort, _ = common.GetICMPv4PortEquivalents(f.IcmpType, f.IcmpCode)
	case protos.ICMP6:
		f.SrcPort, f.DstPort, _ = common.GetICMPv6PortEquivalents(f.IcmpType, f.IcmpCode)
	}

	cmp := bytes.Compare(f.Src, f.Dst)
	if cmp < 0 || (cmp == 0 && f.SrcPort < f.DstPort) {
		return rawIP(f.Src), rawIP(f.Dst), f.SrcPort, f.DstPort
	}
	return rawIP(f.Dst), rawIP(f.Src), f.DstPort, f.SrcPort
}

// Hash creates the community ID for a flow
func (ch *hasher) Hash(f flow.Meta) string {

	// SHA1 hasher
	h := crypto.SHA1.New()
	h.Write(ch.seed[:])

	ip0, ip1, p0, p1 := extractTuple(f)

	binary.Write(h, binary.BigEndian, ip0)
	binary.Write(h, binary.BigEndian, ip1)
	h.Write([]byte{byte(f.Transport), 0})
	binary.Write(h, binary.BigEndian, p0)
	binary.Write(h, binary.BigEndian, p1)

	// defined output as in the specification
	return "1:" + base64.StdEncoding.EncodeToString(h.Sum(nil))
}
