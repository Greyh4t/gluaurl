# gluaurl

[![](https://travis-ci.org/cjoudrey/gluahttp.svg)](https://travis-ci.org/cjoudrey/gluahttp)

gluaurl is a package of url in go for [GopherLua](https://github.com/yuin/gopher-lua).

## Installation

```
go get github.com/Greyh4t/gluaurl
```

## Usage

```go
package main

import (
	"github.com/Greyh4t/gluaurl"
	"github.com/yuin/gopher-lua"
)

func main() {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("url", gluaurl.NewUrl().Loader)

	if err := L.DoString(`
		local url = require("url")
		local x = url.parse("http://www.jd.com:8080/测试/?a=1&a=2&b=3&c=%E6%B5%8B%E8%AF%95#xxx")
		print(x.schema)
		print(x.host)
		print(x.port)
		print(x.rawpath)
		print(x.path)
		print(x.rawquery)
		for key, value in pairs(x.query) do
		    print(key, table.concat(value,","))
		end
		print(x.fragment)
	`); err != nil {
		panic(err)
	}
}
```
