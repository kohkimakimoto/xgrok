package server

import (
	"github.com/cjoudrey/gluahttp"
	"github.com/kohkimakimoto/gluaenv"
	"github.com/kohkimakimoto/gluafs"
	"github.com/kohkimakimoto/gluatemplate"
	"github.com/kohkimakimoto/gluayaml"
	gluajson "github.com/layeh/gopher-json"
	"github.com/otm/gluash"
	"github.com/yuin/gluare"
	"github.com/yuin/gopher-lua"
	"net/http"
)

func initLuaState(L *lua.LState) {
	registerTunnelClass(L)

	// modules
	L.PreloadModule("json", gluajson.Loader)
	L.PreloadModule("fs", gluafs.Loader)
	L.PreloadModule("yaml", gluayaml.Loader)
	L.PreloadModule("template", gluatemplate.Loader)
	L.PreloadModule("env", gluaenv.Loader)
	L.PreloadModule("http", gluahttp.NewHttpModule(&http.Client{}).Loader)
	L.PreloadModule("re", gluare.Loader)
	L.PreloadModule("sh", gluash.Loader)

	lserver := L.NewTable()
	L.SetGlobal("server", lserver)
	luserAuth := L.NewTable()
	L.SetGlobal("user_auth", luserAuth)
	lhooks := L.NewTable()
	L.SetGlobal("hooks", lhooks)
}


const LTunnelClass = "Tunnel*"

func registerTunnelClass(L *lua.LState) {
	mt := L.NewTypeMetatable(LTunnelClass)
	mt.RawSetString("__call", L.NewFunction(tunnelCall))
	mt.RawSetString("__index", L.NewFunction(tunnelIndex))
	mt.RawSetString("__newindex", L.NewFunction(tunnelNewindex))
}


func newLTunnel(L *lua.LState, tunnel *Tunnel) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = tunnel
	L.SetMetatable(ud, L.GetTypeMetatable(LTunnelClass))
	return ud
}

func checkTunnel(L *lua.LState) *Tunnel {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*Tunnel); ok {
		return v
	}
	L.ArgError(1, "Tunnel object expected")
	return nil
}

func tunnelCall(L *lua.LState) int {

	return 0
}

func tunnelIndex(L *lua.LState) int {
	tunnel := checkTunnel(L)
	index := L.CheckString(2)

	if index == "url" {
		L.Push(lua.LString(tunnel.url))
		return 1
	}

	L.Push(lua.LNil)
	return 1
}

func tunnelNewindex(L *lua.LState) int {

	return 0
}


func toLValue(L *lua.LState, value interface{}) lua.LValue {
	switch converted := value.(type) {
	case bool:
		return lua.LBool(converted)
	case float64:
		return lua.LNumber(converted)
	case string:
		return lua.LString(converted)
	case []interface{}:
		arr := L.CreateTable(len(converted), 0)
		for _, item := range converted {
			arr.Append(toLValue(L, item))
		}
		return arr
	case map[string]interface{}:
		tbl := L.CreateTable(0, len(converted))
		for key, item := range converted {
			tbl.RawSetH(lua.LString(key), toLValue(L, item))
		}
		return tbl
	}
	return lua.LNil
}
