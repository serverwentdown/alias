package main

import (
	"testing"

	_ "github.com/serverwentdown/alias"

	"github.com/coredns/coredns/plugin/test"
	ctest "github.com/coredns/coredns/test"
	"github.com/miekg/dns"
)

func TestAlias(t *testing.T) {
	name, rm, err := test.TempFile(".", `$ORIGIN example.org.
@	3600 IN	SOA   sns.dns.icann.org. noc.dns.icann.org. 2017042745 7200 3600 1209600 3600

    3600 IN NS    b.iana-servers.net.

    3600 IN CNAME www.foo
`)
	if err != nil {
		t.Fatalf("Failed to create zone: %s", err)
	}
	defer rm()

	name2, rm2, err2 := test.TempFile(".", `$ORIGIN foo.example.org.
@	3600 IN	SOA sns.dns.icann.org. noc.dns.icann.org. 2017042745 7200 3600 1209600 3600

    3600 IN NS  b.iana-servers.net.

www 3600 IN A   127.0.0.53
`)
	if err2 != nil {
		t.Fatalf("Failed to create zone: %s", err2)
	}
	defer rm2()

	corefile := `example.org:0 {
		file ` + name + ` example.org
		alias
	}
	foo.example.org:0 {
		file ` + name2 + ` foo.example.org
	}`

	i, udp, _, err := ctest.CoreDNSServerAndPorts(corefile)
	if err != nil {
		t.Fatalf("Could not get CoreDNS serving instance: %s", err)
	}
	defer i.Stop()

	m := new(dns.Msg)
	m.SetQuestion("example.org.", dns.TypeA)

	r, err := dns.Exchange(m, udp)
	if err != nil {
		t.Fatalf("Could not exchange msg: %s", err)
	}
	if r.Rcode == dns.RcodeServerFailure {
		t.Fatalf("Rcode should not be dns.RcodeServerFailure")
	}
	t.Log(r.Answer)
	if x := len(r.Answer); x != 1 {
		t.Errorf("Expected 1 RR in reply, got %d", x)
	}
	if x := r.Answer[0].(*dns.A).A.String(); x != "127.0.0.53" {
		t.Errorf("Failed to get address for CNAME, expected 127.0.0.53, got %s", x)
	}
}
