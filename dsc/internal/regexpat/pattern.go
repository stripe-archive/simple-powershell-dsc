package regexpat

import (
	"net/http"
	"regexp"

	"goji.io/pattern"
)

type Pattern struct {
	re      *regexp.Regexp
	methods map[string]struct{}
}

func New(r string) *Pattern {
	return &Pattern{
		re: regexp.MustCompile(r),
	}
}

func NewWithMethods(r string, ms ...string) *Pattern {
	methods := make(map[string]struct{})
	for _, m := range ms {
		methods[m] = struct{}{}
	}
	return &Pattern{
		re:      regexp.MustCompile(r),
		methods: methods,
	}
}

func Get(r string) *Pattern    { return NewWithMethods(r, "GET") }
func Post(r string) *Pattern   { return NewWithMethods(r, "POST") }
func Put(r string) *Pattern    { return NewWithMethods(r, "PUT") }
func Delete(r string) *Pattern { return NewWithMethods(r, "DELETE") }

func (p *Pattern) Match(r *http.Request) *http.Request {
	if p.methods != nil {
		if _, ok := p.methods[r.Method]; !ok {
			return nil
		}
	}

	ctx := r.Context()
	path := pattern.Path(ctx)

	reMatches := p.re.FindStringSubmatch(path)
	if len(reMatches) == 0 {
		return nil
	}

	matches := make(map[pattern.Variable]string)
	for i, name := range p.re.SubexpNames() {
		if i == 0 {
			continue
		}
		matches[pattern.Variable(name)] = reMatches[i]
	}

	return r.WithContext(&regexMatch{ctx, matches})
}

func (p *Pattern) HTTPMethods() map[string]struct{} {
	return p.methods
}

func (p *Pattern) PathPrefix() string {
	prefix, _ := p.re.LiteralPrefix()
	return prefix
}
