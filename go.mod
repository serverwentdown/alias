module github.com/serverwentdown/alias

go 1.13

require (
	github.com/caddyserver/caddy v1.0.3
	github.com/coredns/coredns v1.6.4
	github.com/miekg/dns v1.1.17
	golang.org/x/net v0.0.0-20190827160401-ba9fcec4b297
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.0.0+incompatible
	github.com/miekg/dns v1.1.3 => github.com/miekg/dns v1.1.17
	golang.org/x/net v0.0.0-20190813000000-74dc4d7220e7 => golang.org/x/net v0.0.0-20190827160401-ba9fcec4b297
)
