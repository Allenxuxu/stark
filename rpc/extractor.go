package rpc

import (
	"fmt"
	"go/ast"
	"reflect"

	"github.com/Allenxuxu/stark/registry"
)

func extractEndpoints(handler interface{}) []*registry.Endpoint {
	typ := reflect.TypeOf(handler)
	hdlr := reflect.ValueOf(handler)
	name := reflect.Indirect(hdlr).Type().Name()

	var endpoints []*registry.Endpoint
	for m := 0; m < typ.NumMethod(); m++ {
		if e := extractEndpoint(typ.Method(m)); e != nil {
			e.Name = name + "." + e.Name
			endpoints = append(endpoints, e)
		}
	}

	return endpoints
}
func extractEndpoint(method reflect.Method) *registry.Endpoint {
	if method.PkgPath != "" {
		return nil
	}

	// todo do better
	in := &registry.Value{
		Name:   "in parameter",
		Type:   "",
		Values: nil,
	}
	out := &registry.Value{
		Name:   "out parameter",
		Type:   "",
		Values: nil,
	}

	mt := method.Type
	for i := 1; i < mt.NumIn(); i++ {
		in.Values = append(in.Values, extractValue(mt.In(i)))
	}

	for i := 0; i < mt.NumOut(); i++ {
		out.Values = append(out.Values, extractValue(mt.Out(i)))
	}

	ep := &registry.Endpoint{
		Name:     method.Name,
		Request:  in,
		Response: out,
		Metadata: make(map[string]string),
	}

	if method.Type.NumOut() == 1 {
		ep.Metadata = map[string]string{
			"stream": fmt.Sprintf("%v", true),
		}
	}

	return ep
}

func extractValue(v reflect.Type) *registry.Value {
	if v == nil {
		return nil
	}

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	arg := &registry.Value{
		Name: v.Name(),
		Type: v.Kind().String(),
	}

	return arg
}

func isExportedOrBuiltinType(t reflect.Type) bool {
	return ast.IsExported(t.Name()) || t.PkgPath() == ""
}
