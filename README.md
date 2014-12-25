xsrf [![Build Status](https://drone.io/github.com/tango-contrib/xsrf/status.png)](https://drone.io/github.com/tango-contrib/xsrf/latest) [![](http://gocover.io/_badge/github.com/tango-contrib/xsrf)](http://gocover.io/github.com/tango-contrib/xsrf)
======

Middleware xsrf is a xsrf checker for [Tango](https://github.com/lunny/tango). 

## Installation

    go get github.com/tango-contrib/xsrf

## Simple Example

```Go
type XsrfAction struct {
    render.Render
    xsrf.Checker
}

func (x *XsrfAction) Get() {
    x.RenderFile("test.html", render.T{
        "XsrfFormHtml": xsrf.XsrfFormHtml,
    })
}

func (x *XsrfAction) Post() {
    // before this call, xsrf will be checked
}
```

If you don't want some action do not check, then
```Go
type NoCheckAction struct {
    xsrf.NoCheck
}
```
will be ok.