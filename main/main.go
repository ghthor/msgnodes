package main

import (
	"fmt"
	i "ghthor/init"
	//"runtime"
	//"ghthor/node"
)

type test struct {
	i.InitVar
	name string
}

func (t *test) Init(name interface {}) (interface {}) {
	initVar := t.InitVar.Init(name)
	switch initVar.(type) {
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

func main() {
	fmt.Printf("\n\n");

	temp := new(test).Init("I've Been Awake to Long I guess").(*test)
	taco := new(test).Init("Taco").(*test)
	clone := new(test).Init(temp).(*test)

	fmt.Println(temp.name, "\n",  taco.name, "\n",clone.name)
	fmt.Printf("\n\n");
}
