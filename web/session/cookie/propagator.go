package cookie

import "net/http"

type PropagatorOption func(p *Propagator)

func WithCookieOption(opt func(c *http.Cookie)) PropagatorOption {
	return func(p *Propagator) {
		p.cookieOpt = opt
	}
}

type Propagator struct {
	cookieName string
	cookieOpt  func(c *http.Cookie)
}

func NewPropagator(cookieName string, opts ...PropagatorOption) *Propagator {
	return &Propagator{
		cookieName: cookieName,
		cookieOpt:  func(c *http.Cookie) {},
	}
}

func (p *Propagator) Inject(id string, writer http.ResponseWriter) error {
	cookie := &http.Cookie{
		Name:  p.cookieName,
		Value: id,
	}
	p.cookieOpt(cookie)
	http.SetCookie(writer, cookie)
	return nil
}

func (p *Propagator) Extract(req *http.Request) (string, error) {
	cookie, err := req.Cookie(p.cookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func (p *Propagator) Remove(writer http.ResponseWriter) error {
	cookie := &http.Cookie{
		Name:   p.cookieName,
		MaxAge: -1,
	}
	p.cookieOpt(cookie)
	http.SetCookie(writer, cookie)
	return nil
}
