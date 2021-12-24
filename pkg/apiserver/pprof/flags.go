package pprof

import (
	"flag"
	"strings"
)

type flags struct {
	args  []string
	s     flag.FlagSet
	usage []string
}

// Bool implements the plugin.FlagSet interface.
func (p *flags) Bool(o string, d bool, c string) *bool {
	return p.s.Bool(o, d, c)
}

// Int implements the plugin.FlagSet interface.
func (p *flags) Int(o string, d int, c string) *int {
	return p.s.Int(o, d, c)
}

// Float64 implements the plugin.FlagSet interface.
func (p *flags) Float64(o string, d float64, c string) *float64 {
	return p.s.Float64(o, d, c)
}

// String implements the plugin.FlagSet interface.
func (p *flags) String(o, d, c string) *string {
	return p.s.String(o, d, c)
}

// BoolVar implements the plugin.FlagSet interface.
func (p *flags) BoolVar(b *bool, o string, d bool, c string) {
	p.s.BoolVar(b, o, d, c)
}

// IntVar implements the plugin.FlagSet interface.
func (p *flags) IntVar(i *int, o string, d int, c string) {
	p.s.IntVar(i, o, d, c)
}

// Float64Var implements the plugin.FlagSet interface.
// the value of the flag.
func (p *flags) Float64Var(f *float64, o string, d float64, c string) {
	p.s.Float64Var(f, o, d, c)
}

// StringVar implements the plugin.FlagSet interface.
func (p *flags) StringVar(s *string, o, d, c string) {
	p.s.StringVar(s, o, d, c)
}

// StringList implements the plugin.FlagSet interface.
func (p *flags) StringList(o, d, c string) *[]*string {
	return &[]*string{p.s.String(o, d, c)}
}

// AddExtraUsage implements the plugin.FlagSet interface.
func (p *flags) AddExtraUsage(eu string) {
	p.usage = append(p.usage, eu)
}

// ExtraUsage implements the plugin.FlagSet interface.
func (p *flags) ExtraUsage() string {
	return strings.Join(p.usage, "\n")
}

// Parse implements the plugin.FlagSet interface.
func (p *flags) Parse(usage func()) []string {
	p.s.Usage = usage
	_ = p.s.Parse(p.args)
	args := p.s.Args()
	if len(args) == 0 {
		usage()
	}
	return args
}
