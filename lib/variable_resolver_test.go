package lib

import (
	"flag"
	"testing"
)

func TestVariableResolver(t *testing.T)  {
	//os.Setenv("MY_VAR", "my value")
	//os.Setenv("MY_VAR_1", "my value 1")

	flag.Set("myVar", "my value")

	//var myVar string

	r := NewVariableResolver()

	r.
		//Define(&myVar, "myVar", "MY_VAR", "my variable").
		Resolve()

	t.Log(*r.Defs[0].Val)

}
