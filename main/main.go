package main

import (
	"fmt"
	i "ghthor/init"
	//"runtime"
	//"ghthor/node"
)

type Object interface {
	Init(interface {})(interface {})
	Dispose()
}

type test struct {
	i.InitVar
	name string
}

func (t *test) Init(name interface {}) (interface {}) {
	initVar := t.InitVar.Init(name)
	switch initVar.(type) {
		case i.Default:
			t.name = "tacos"
		case *test:
			t.name = name.(*test).name
		case string:
			t.name = initVar.(string)
		default:
			t.name = "Invalid Initialization Object"
	}
	t.InitArg = initVar
	return t
}

func (t *test) Dispose() {
	//close(t.tacos)
}


func main() {
	fmt.Printf("\n\n");

	temp := (new(test)).Init("tacos yay").(*test)
	tacos := new(test).Init(nil).(*test)
	bad := new(test).Init(10).(*test)
	clone := new(test).Init(temp).(*test)

	fmt.Println(temp.name)
	fmt.Println(tacos.name)
	fmt.Println(bad.name)
	fmt.Println(clone.name)

	fmt.Printf("\n\n");
}
