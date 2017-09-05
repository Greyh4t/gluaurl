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

	L.PreloadModule("url", gluaurl.Loader)

	if err := L.DoString(`
		local url = require("url")
		local x = url.parse("http://www.example.com:8080/测试/?a=1&a=2&b=6&e=&c%5B%5D=3&c%5B%5D=4&c%5B%5D=5&d=%E6%B5%8B%E8%AF%95#xxx")
		print(x.scheme)
		print(x.host)
		print(x.port)
		print(x.rawpath)
		print(x.path)
		print(x.rawquery)
		for key, value in pairs(x.query) do
		    print(key, value)
		end

		print(x.fragment)
		options={
			scheme=x.scheme,
			host=x.host,
			path=x.path,
			rawquery=url.build_query_string(x.query),
			fragment=x.fragment
		}
		print(url.build(options))

		options={
			scheme="https",
			host="www.example.com:8888",
			path="index/",
			rawquery="a=1&b[]=2&c=3"
		}
		print(url.build(options))

		options={
			a=1,
			b=2,
			c={3,4},
			d={e=5,f="x"}
		}
		print(url.build_query_string(options))

		print(url.resolve("http://www.example.com","/index/test?a=1"))
	`); err != nil {
		panic(err)
	}
}
```
