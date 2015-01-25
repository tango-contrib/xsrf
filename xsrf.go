package xsrf

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/go-xweb/uuid"
	"github.com/lunny/tango"
)

const (
	XSRF_TAG string = "_xsrf"
)

type Xsrfer interface {
	CheckXsrf() bool
}

type NoCheck struct {
}

func (NoCheck) CheckXsrf() bool {
	return false
}

type XsrfChecker interface {
	SetXsrfValue(string)
}

type Checker struct {
	XsrfValue string
}

func (Checker) CheckXsrf() bool {
	return true
}

func (c *Checker) SetXsrfValue(v string) {
	c.XsrfValue = v
}

func (c *Checker) XsrfFormHtml() template.HTML {
	return template.HTML(fmt.Sprintf(`<input type="hidden" name="%v" value="%v" />`,
		XSRF_TAG, c.XsrfValue))
}

type Xsrf struct {
	timeout time.Duration
}

func New(timeout time.Duration) *Xsrf {
	return &Xsrf{
		timeout: timeout,
	}
}

// NewCookie is a helper method that returns a new http.Cookie object.
// Duration is specified in seconds. If the duration is zero, the cookie is permanent.
// This can be used in conjunction with ctx.SetCookie.
func newCookie(name string, value string, age int64) *http.Cookie {
	var utctime time.Time
	if age == 0 {
		// 2^31 - 1 seconds (roughly 2038)
		utctime = time.Unix(2147483647, 0)
	} else {
		utctime = time.Unix(time.Now().Unix()+age, 0)
	}
	return &http.Cookie{Name: name, Value: value, Expires: utctime}
}

func (xsrf *Xsrf) Handle(ctx *tango.Context) {
	var action interface{}
	if action = ctx.Action(); action == nil {
		ctx.Next()
		return
	}

	// if action implements check xsrf option and ask not check then return
	if checker, ok := action.(Xsrfer); ok && !checker.CheckXsrf() {
		ctx.Next()
		return
	}

	var val string = ""
	cookie, err := ctx.Req().Cookie(XSRF_TAG)
	if err != nil {
		val = uuid.NewRandom().String()
		cookie = newCookie(XSRF_TAG, val, int64(xsrf.timeout))
		ctx.Header().Set("Set-Cookie", cookie.String())
	} else {
		val = cookie.Value
	}

	if c, ok := action.(XsrfChecker); ok {
		c.SetXsrfValue(val)
	}

	if ctx.Req().Method == "POST" {
		res, err := ctx.Req().Cookie(XSRF_TAG)
		formVal := ctx.Req().FormValue(XSRF_TAG)

		if err != nil || res.Value == "" || res.Value != formVal {
			ctx.Abort(http.StatusInternalServerError, "xsrf token error.")
			ctx.Error("xsrf token error.")
			return
		}
	}

	ctx.Next()
}
