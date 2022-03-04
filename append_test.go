package main

import (
	"fmt"
	"testing"
)

func TestAppend(t *testing.T) {
	list:=[]element{
		{F:         10},
		{F:         3},
		{F:         14},
		{F:         12},
		{F:         5},
		{F:         5},
	}
	result:=make([]element,0)
	for _,item:=range list{
		result=Append(result,item)
	}
	_=result
}

func TestName(t *testing.T) {
	list:=[]int{0,0,1}
	fmt.Println(list[:3])
}