module github.com/serverwentdown/alias/coredns

go 1.13

replace github.com/serverwentdown/alias => ../

require (
	github.com/coredns/coredns v1.6.9
	github.com/miekg/dns v1.1.29
	github.com/serverwentdown/alias v0.0.0-00010101000000-000000000000
)
