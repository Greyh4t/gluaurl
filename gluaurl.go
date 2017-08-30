package gluaurl

import (
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"strings"

	"github.com/yuin/gopher-lua"
)

var rBracket = regexp.MustCompile("\\[\\]$")

func Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"parse":              parse,
		"build":              build,
		"build_query_string": buildQueryString,
		"resolve":            resolve,
	})
	L.Push(mod)
	return 1
}

func parse(L *lua.LState) int {
	parsed := L.NewTable()
	obj, err := url.Parse(L.CheckString(1))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	parsed.RawSetString("scheme", lua.LString(obj.Scheme))
	parsed.RawSetString("rawquery", lua.LString(obj.RawQuery))
	parsed.RawSetString("path", lua.LString(obj.Path))
	parsed.RawSetString("rawpath", lua.LString(obj.RawPath))
	parsed.RawSetString("port", lua.LString(obj.Port()))
	parsed.RawSetString("fragment", lua.LString(obj.Fragment))
	parsed.RawSetString("scheme", lua.LString(obj.Scheme))

	if len(obj.Query()) > 0 {
		query := L.NewTable()
		for k, vlist := range obj.Query() {
			each := L.NewTable()
			for i, v := range vlist {
				each.Insert(i+1, lua.LString(v))
			}
			query.RawSetString(k, each)
		}
		parsed.RawSetString("query", query)
	}

	if obj.User != nil {
		parsed.RawSetString("username", lua.LString(obj.User.Username()))
		if password, hasPassword := obj.User.Password(); hasPassword {
			parsed.RawSetString("password", lua.LString(password))
		} else {
			parsed.RawSetString("password", lua.LNil)
		}
	} else {
		parsed.RawSetString("username", lua.LNil)
		parsed.RawSetString("password", lua.LNil)
	}

	L.Push(parsed)
	return 1
}

func build(L *lua.LState) int {
	options := L.CheckTable(1)

	buildUrl := url.URL{}

	if scheme := options.RawGetString("scheme"); scheme != lua.LNil {
		buildUrl.Scheme = scheme.String()
	}

	if username := options.RawGetString("username"); username != lua.LNil {
		if password := options.RawGetString("password"); password != lua.LNil {
			buildUrl.User = url.UserPassword(username.String(), password.String())
		} else {
			buildUrl.User = url.User(username.String())
		}
	}

	if host := options.RawGetString("host"); host != lua.LNil {
		buildUrl.Host = host.String()
	}

	if path := options.RawGetString("path"); path != lua.LNil {
		buildUrl.Path = path.String()
	}

	if query := options.RawGetString("rawquery"); query != lua.LNil {
		buildUrl.RawQuery = query.String()
	}

	if fragment := options.RawGetString("fragment"); fragment != lua.LNil {
		buildUrl.Fragment = fragment.String()
	}

	L.Push(lua.LString(buildUrl.String()))

	return 1
}

func buildQueryString(L *lua.LState) int {
	options := L.CheckTable(1)

	ret := make([]string, 0)

	options.ForEach(func(key, value lua.LValue) {
		toQueryString(key.String(), value, &ret)
	})

	sort.Strings(ret)

	L.Push(lua.LString(strings.Join(ret, "&")))

	return 1
}

func toQueryString(prefix string, lv lua.LValue, ret *[]string) {
	switch v := lv.(type) {
	case lua.LBool:
		*ret = append(*ret, url.QueryEscape(prefix)+"="+v.String())
		break

	case lua.LNumber:
		*ret = append(*ret, url.QueryEscape(prefix)+"="+v.String())
		break

	case lua.LString:
		*ret = append(*ret, url.QueryEscape(prefix)+"="+url.QueryEscape(v.String()))
		break

	case *lua.LTable:
		maxn := v.MaxN()
		if maxn == 0 {
			ret2 := make([]string, 0)
			v.ForEach(func(key lua.LValue, value lua.LValue) {
				toQueryString(prefix+"["+key.String()+"]", value, &ret2)
			})
			sort.Strings(ret2)
			*ret = append(*ret, strings.Join(ret2, "&"))
		} else {
			ret2 := make([]string, 0)
			for i := 1; i <= maxn; i++ {
				vi := v.RawGetInt(i)

				if rBracket.MatchString(prefix) {
					ret2 = append(ret2, url.QueryEscape(prefix)+"="+vi.String())
				} else {
					if vi.Type() == lua.LTTable {
						toQueryString(fmt.Sprintf("%s[%d]", prefix, i-1), vi, &ret2)
					} else {
						toQueryString(prefix+"[]", vi, &ret2)
					}
				}
			}
			*ret = append(*ret, strings.Join(ret2, "&"))
		}
		break
	}
}

func resolve(L *lua.LState) int {
	from := L.CheckString(1)
	to := L.CheckString(2)

	fromUrl, err := url.Parse(from)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	toUrl, err := url.Parse(to)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	resolvedUrl := fromUrl.ResolveReference(toUrl).String()
	L.Push(lua.LString(resolvedUrl))
	return 1
}
