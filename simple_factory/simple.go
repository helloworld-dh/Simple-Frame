package simple_factory

import "fmt"

type Api interface {
	say(name string) string
}

type IPhone struct {
	name string
}

func (I *IPhone) say(name string) string {
	fmt.Printf("Phone is %s", name)
	return name
}

type HuaWei struct {
	name string
}

func (H *HuaWei) say(name string) string {
	fmt.Printf("HuaWei is %s", name)
	return name
}

func NewApi(name string) Api {
	if name == "iphone" {
		return &IPhone{
			name: name,
		}
	} else {
		return &HuaWei{name: name}
	}
}
