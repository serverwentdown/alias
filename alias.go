package alias

import (
	"github.com/coredns/coredns/plugin"

	"github.com/miekg/dns"

	"golang.org/x/net/context"
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

// WriteMsg records the status code and calls the
// underlying ResponseWriter's WriteMsg method.
func (r *ResponseModifier) WriteMsg(res *dns.Msg) error {
	// Guess zone based on authority section.
	var zone string
	for i := 0; i < len(res.Ns); i++ {
		rr := res.Ns[i]
		if rr.Header().Rrtype == dns.TypeNS {
			zone = rr.Header().Name
		}
	}

	// Find and delete CNAME record on that zone, storing the canonical name.
	var cname string
	for i := 0; i < len(res.Answer); i++ {
		rr := res.Answer[i]
		if rr.Header().Rrtype == dns.TypeCNAME && rr.Header().Name == zone {
			cname = rr.(*dns.CNAME).Target
			// Remove the CNAME record
			res.Answer = append(res.Answer[:i], res.Answer[i+1:]...)
			break
		}
	}

	// Rename all the records with the above canonical name to the zone name
	for i := 0; i < len(res.Answer); i++ {
		rr := res.Answer[i]
		if rr.Header().Name == cname {
			rr.Header().Name = zone
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
