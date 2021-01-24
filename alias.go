package alias

import (
	"context"
	"math"

	"github.com/coredns/coredns/plugin"

	"github.com/miekg/dns"
)

// Rewrite is plugin to rewrite requests internally before being handled.
type Alias struct {
	Next plugin.Handler
}

// ServeDNS implements the plugin.Handler interface.
func (al Alias) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	mw := NewResponseModifier(w)
	return plugin.NextOrFailure(al.Name(), al.Next, ctx, mw, r)
}

// Name implements the Handler interface.
func (al Alias) Name() string { return "alias" }

type ResponseModifier struct {
	dns.ResponseWriter
}

// Returns a dns.Msg modifier that replaces CNAME on root zones with other records.
func NewResponseModifier(w dns.ResponseWriter) *ResponseModifier {
	return &ResponseModifier{
		ResponseWriter: w,
	}
}

func min(a, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}

// WriteMsg records the status code and calls the
// underlying ResponseWriter's WriteMsg method.
func (r *ResponseModifier) WriteMsg(res *dns.Msg) error {
	// Guess zone based on authority section.
	var zone string
	for _, rr := range res.Ns {
		if rr.Header().Rrtype == dns.TypeNS {
			zone = rr.Header().Name
		}
	}

	// Find and delete CNAME record on that zone, storing the canonical name.
	var (
		cname string = zone
		ttl   uint32 = math.MaxUint32
	)
	for i := 0; i < len(res.Answer); {
		rr := res.Answer[i]
		if rr.Header().Rrtype == dns.TypeCNAME && rr.Header().Name == cname {
			cname = rr.(*dns.CNAME).Target
			ttl = min(ttl, rr.(*dns.CNAME).Header().Ttl)
			// Remove the CNAME record
			res.Answer = append(res.Answer[:i], res.Answer[i+1:]...)
			continue
		}
		i++
	}

	// Rename all the records with the above canonical name to the zone name
	for _, rr := range res.Answer {
		if rr.Header().Name == cname {
			rr.Header().Name = zone
			rr.Header().Ttl = min(ttl, rr.Header().Ttl)
		}
	}

	return r.ResponseWriter.WriteMsg(res)
}

// Write is a wrapper that records the size of the message that gets written.
func (r *ResponseModifier) Write(buf []byte) (int, error) {
	n, err := r.ResponseWriter.Write(buf)
	return n, err
}

// Hijack implements dns.Hijacker. It simply wraps the underlying
// ResponseWriter's Hijack method if there is one, or returns an error.
func (r *ResponseModifier) Hijack() {
	r.ResponseWriter.Hijack()
	return
}
