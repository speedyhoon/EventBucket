package main

import "fmt"
import "reflect"

func Call(name interface{}, params ... interface{}) map[reflect.Value][]reflect.Value{// (result []reflect.Value, err error)
		f := reflect.ValueOf(name)
	if len(params) != f.Type().NumIn() {
		fmt.Print("whoops there was an error!")
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	return map[reflect.Value][]reflect.Value{
		f: in,
	}
}

func main() {
	router := router()
	ioc(router["bar"])
}

func ioc(model map[reflect.Value][]reflect.Value){
	for function, parameters := range model{
		//TODO it would be nice to call the needed functions in order before calling the model I'm note sure if this is possible though
		function.Call(parameters)
	}
}

func router()(map[string]map [reflect.Value][]reflect.Value){
	elements := map[string]map[reflect.Value][]reflect.Value{
		"foo": Call(foo),
		"bar":Call(bar, 1, 4, 6),
	}
	return elements
}

func foo() {
	fmt.Println("we are running foo")
}

func bar(a, b, c int) {
	fmt.Println("we are running bar", a, b, c)
}
