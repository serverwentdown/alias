package main

import (
	_ "github.com/coredns/coredns/core/plugin"
	_ "github.com/serverwentdown/alias"

	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/coremain"
)

func init() {
	// Insert alias directive before file directive
	var i int
	for i = 0; i < len(dnsserver.Directives); i++ {
		if dnsserver.Directives[i] == "file" {
			break
		}
	}
	dnsserver.Directives = append(dnsserver.Directives, "")
	copy(dnsserver.Directives[i+1:], dnsserver.Directives[i:])
	dnsserver.Directives[i] = "alias"
}

func main() {
	coremain.Run()
}
