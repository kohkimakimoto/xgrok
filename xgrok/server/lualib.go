package server

import (
	"fmt"
	"github.com/cjoudrey/gluahttp"
	"github.com/kohkimakimoto/gluaenv"
	"github.com/kohkimakimoto/gluafs"
	"github.com/kohkimakimoto/gluatemplate"
	"github.com/kohkimakimoto/gluayaml"
	"github.com/kohkimakimoto/xgrok/xgrok/msg"
	gluajson "github.com/layeh/gopher-json"
	"github.com/otm/gluash"
	"github.com/yuin/gluare"
	"github.com/yuin/gopher-lua"
	"net/http"
)

func initLuaState(L *lua.LState) {
	registerTunnelClass(L)
	registerMsgAuthClass(L)
	registerMsgAuthRespClass(L)
	registerMsgReqTunnelClass(L)
	registerMsgNewTunnelClass(L)

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

const LMsgAuthRespClass = "MsgAuthResp*"

func registerMsgAuthRespClass(L *lua.LState) {
	mt := L.NewTypeMetatable(LMsgAuthRespClass)
	mt.RawSetString("__call", L.NewFunction(msgAuthRespCall))
	mt.RawSetString("__index", L.NewFunction(msgAuthRespIndex))
	mt.RawSetString("__newindex", L.NewFunction(msgAuthRespNewindex))
}

func newLMsgAuthResp(L *lua.LState, authResp *msg.AuthResp) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = authResp
	L.SetMetatable(ud, L.GetTypeMetatable(LMsgAuthRespClass))
	return ud
}

func checkMsgAuthResp(L *lua.LState) *msg.AuthResp {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*msg.AuthResp); ok {
		return v
	}
	L.ArgError(1, "MsgAuthResp object expected")
	return nil
}

func msgAuthRespCall(L *lua.LState) int {

	return 0
}

func msgAuthRespIndex(L *lua.LState) int {
	msgAuthResp := checkMsgAuthResp(L)
	index := L.CheckString(2)

	switch index {
	case "version":
		L.Push(lua.LString(msgAuthResp.Version))
		return 1
	case "mm_version":
		L.Push(lua.LString(msgAuthResp.MmVersion))
		return 1
	case "client_id":
		L.Push(lua.LString(msgAuthResp.ClientId))
		return 1
	}

	L.Push(lua.LNil)
	return 1
}

func msgAuthRespNewindex(L *lua.LState) int {

	return 0
}

const LMsgNewTunnelClass = "MsgNewTunnel*"

func registerMsgNewTunnelClass(L *lua.LState) {
	mt := L.NewTypeMetatable(LMsgNewTunnelClass)
	mt.RawSetString("__call", L.NewFunction(msgNewTunnelCall))
	mt.RawSetString("__index", L.NewFunction(msgNewTunnelIndex))
	mt.RawSetString("__newindex", L.NewFunction(msgNewTunnelNewindex))
}

func newLMsgNewTunnel(L *lua.LState, newTunnel *msg.NewTunnel) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = newTunnel
	L.SetMetatable(ud, L.GetTypeMetatable(LMsgNewTunnelClass))
	return ud
}

func checkMsgNewTunnel(L *lua.LState) *msg.NewTunnel {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*msg.NewTunnel); ok {
		return v
	}
	L.ArgError(1, "MsgNewTunnel object expected")
	return nil
}

func msgNewTunnelCall(L *lua.LState) int {

	return 0
}

func msgNewTunnelIndex(L *lua.LState) int {
	newTunnel := checkMsgNewTunnel(L)
	index := L.CheckString(2)

	switch index {
	case "url":
		L.Push(lua.LString(newTunnel.Url))
		return 1
	case "public_url":
		L.Push(lua.LString(newTunnel.PublicUrl))
		return 1
	case "proto":
		L.Push(lua.LString(newTunnel.Protocol))
		return 1
	case "req_id":
		L.Push(lua.LString(newTunnel.ReqId))
		return 1
	case "custom_props":
		L.Push(newTunnel.LCustomProps)
		return 1
	}

	L.Push(lua.LNil)
	return 1
}

func msgNewTunnelNewindex(L *lua.LState) int {
	newTunnel := checkMsgNewTunnel(L)
	index := L.CheckString(2)

	switch index {
	case "url":
		newTunnel.Url = L.CheckString(3)
		return 0

	case "public_url":
		newTunnel.PublicUrl = L.CheckString(3)
		return 0

	case "proto":
		newTunnel.Protocol = L.CheckString(3)
		return 0
	case "req_id":
		newTunnel.ReqId = L.CheckString(3)
		return 1
	case "custom_props":
		newTunnel.LCustomProps = L.CheckTable(3)

		// parse to go value
		tb := newTunnel.LCustomProps
		tb.ForEach(func(k, v lua.LValue) {
			kvpair := v.(*lua.LTable)
			if kvpair != nil {
				kvpair.ForEach(func(key, value lua.LValue) {
					keystr := fmt.Sprint(toGoValue(key))
					valuestr := fmt.Sprint(toGoValue(value))

					newTunnel.CustomProps = append(newTunnel.CustomProps, msg.CustomProp{
						Key:   keystr,
						Value: valuestr,
					})
				})
			}

		})
		//
		return 1
	}

	return 0
}

const LMsgAuthClass = "MsgAuth*"

func registerMsgAuthClass(L *lua.LState) {
	mt := L.NewTypeMetatable(LMsgAuthClass)
	mt.RawSetString("__call", L.NewFunction(msgAuthCall))
	mt.RawSetString("__index", L.NewFunction(msgAuthIndex))
	mt.RawSetString("__newindex", L.NewFunction(msgAuthNewindex))
}

func newLMsgAuth(L *lua.LState, auth *msg.Auth) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = auth
	L.SetMetatable(ud, L.GetTypeMetatable(LMsgAuthClass))
	return ud
}

func checkMsgAuth(L *lua.LState) *msg.Auth {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*msg.Auth); ok {
		return v
	}
	L.ArgError(1, "MsgAuth object expected")
	return nil
}

func msgAuthCall(L *lua.LState) int {

	return 0
}

func msgAuthIndex(L *lua.LState) int {
	msgAuth := checkMsgAuth(L)
	index := L.CheckString(2)

	switch index {
	case "version":
		L.Push(lua.LString(msgAuth.Version))
		return 1
	case "mm_version":
		L.Push(lua.LString(msgAuth.MmVersion))
		return 1
	case "user":
		L.Push(lua.LString(msgAuth.User))
		return 1
	case "password":
		L.Push(lua.LString(msgAuth.Password))
		return 1
	case "os":
		L.Push(lua.LString(msgAuth.OS))
		return 1
	case "arch":
		L.Push(lua.LString(msgAuth.Arch))
		return 1
	case "client_id":
		L.Push(lua.LString(msgAuth.ClientId))
		return 1
	}

	L.Push(lua.LNil)
	return 1
}

func msgAuthNewindex(L *lua.LState) int {
	auth := checkMsgAuth(L)
	index := L.CheckString(2)

	switch index {
	case "version":
		auth.Version = L.CheckString(3)
		return 0
	case "mm_version":
		auth.Version = L.CheckString(3)
		return 0
	case "user":
		auth.User = L.CheckString(3)
		return 0
	case "password":
		auth.Password = L.CheckString(3)
		return 0
	case "os":
		auth.OS = L.CheckString(3)
		return 0
	case "arch":
		auth.Arch = L.CheckString(3)
		return 0
	case "client_id":
		auth.ClientId = L.CheckString(3)
		return 0
	}

	return 0
}

const LMsgReqTunnelClass = "MsgReqTunnel*"

func registerMsgReqTunnelClass(L *lua.LState) {
	mt := L.NewTypeMetatable(LMsgReqTunnelClass)
	mt.RawSetString("__call", L.NewFunction(msgReqTunnelCall))
	mt.RawSetString("__index", L.NewFunction(msgReqTunnelIndex))
	mt.RawSetString("__newindex", L.NewFunction(msgReqTunnelNewindex))
}

func newLMsgReqTunnel(L *lua.LState, auth *msg.ReqTunnel) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = auth
	L.SetMetatable(ud, L.GetTypeMetatable(LMsgReqTunnelClass))
	return ud
}

func checkMsgReqTunnel(L *lua.LState) *msg.ReqTunnel {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*msg.ReqTunnel); ok {
		return v
	}
	L.ArgError(1, "MsgReqTunnel object expected")
	return nil
}

func msgReqTunnelCall(L *lua.LState) int {

	return 0
}

func msgReqTunnelIndex(L *lua.LState) int {
	msgReqTunnel := checkMsgReqTunnel(L)
	index := L.CheckString(2)

	switch index {
	case "req_id":
		L.Push(lua.LString(msgReqTunnel.ReqId))
		return 1
	case "protocol":
		L.Push(lua.LString(msgReqTunnel.Protocol))
		return 1
	case "hostname":
		L.Push(lua.LString(msgReqTunnel.Hostname))
		return 1
	case "subdomain":
		L.Push(lua.LString(msgReqTunnel.Subdomain))
		return 1
	case "http_auth":
		L.Push(lua.LString(msgReqTunnel.HttpAuth))
		return 1
	case "remote_port":
		L.Push(lua.LNumber(msgReqTunnel.RemotePort))
		return 1
	}

	L.Push(lua.LNil)
	return 1
}

func msgReqTunnelNewindex(L *lua.LState) int {
	reqTunnel := checkMsgReqTunnel(L)
	index := L.CheckString(2)

	switch index {
	case "req_id":
		reqTunnel.ReqId = L.CheckString(3)
		return 0
	case "protocol":
		reqTunnel.Protocol = L.CheckString(3)
		return 0
	case "hostname":
		reqTunnel.Hostname = L.CheckString(3)
		return 0
	case "subdomain":
		reqTunnel.Subdomain = L.CheckString(3)
		return 0
	case "http_auth":
		reqTunnel.HttpAuth = L.CheckString(3)
		return 0
	case "remote_port":
		reqTunnel.RemotePort = uint16(L.CheckNumber(3))
		return 0

	}

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

// This code inspired by https://github.com/yuin/gluamapper/blob/master/gluamapper.go
func toGoValue(lv lua.LValue) interface{} {
	switch v := lv.(type) {
	case *lua.LNilType:
		return nil
	case lua.LBool:
		return bool(v)
	case lua.LString:
		return string(v)
	case lua.LNumber:
		return float64(v)
	case *lua.LTable:
		maxn := v.MaxN()
		if maxn == 0 { // table
			ret := make(map[string]interface{})
			v.ForEach(func(key, value lua.LValue) {
				keystr := fmt.Sprint(toGoValue(key))
				ret[keystr] = toGoValue(value)
			})
			return ret
		} else { // array
			ret := make([]interface{}, 0, maxn)
			for i := 1; i <= maxn; i++ {
				ret = append(ret, toGoValue(v.RawGetInt(i)))
			}
			return ret
		}
	default:
		return v
	}
}
