// Copyright 2020 Oz Tiram <oz.tiram@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package zerotier

import (
	"context"
	"io"
	"net"
	"os"
	"strings"
	"time"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/metrics"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

// Define log to be a logger with the plugin name in it. This way we can just use log.Info and
// friends to log.
var log = clog.NewWithPlugin("zerotier")

// Zerotier struct
type Zerotier struct {
	Url           string
	Token         string
	CacheDuration time.Duration // CacheDuration -
	Next          plugin.Handler
}

func (n Zerotier) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {

	answers := []dns.RR{}
	state := request.Request{W: w, Req: r}

	ip_addresses := query(n.Url, n.Token, strings.TrimRight(state.QName(), "."), n.CacheDuration)
	// no IP is found in zerotier pass processing to the next plugin
	if len(ip_addresses) == 0 {
		return plugin.NextOrFailure(n.Name(), n.Next, ctx, w, r)
	}

	// Export metric with the server label set to the current
	// server handling the request.
	requestCount.WithLabelValues(metrics.WithServer(ctx)).Inc()

	for i := 0; i < len(ip_addresses); i++ {
		rec := new(dns.A)
		rec.Hdr = dns.RR_Header{Name: state.QName(), Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: uint32(n.CacheDuration.Seconds())}
		rec.A = net.ParseIP(ip_addresses[i])
		answers = append(answers, rec)
	}

	m := new(dns.Msg)
	m.Answer = answers
	m.SetReply(r)
	w.WriteMsg(m)

	return dns.RcodeSuccess, nil
}

// Name implements the Handler interface.
func (n Zerotier) Name() string { return "zerotier" }

// Make out a reference to os.Stdout so we can easily overwrite it for testing.
var out io.Writer = os.Stdout
