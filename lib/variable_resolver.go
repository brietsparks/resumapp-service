package lib

import (
	"flag"
	"os"
)

type VariableDef struct {
	Val *string
	Flag string
	Env string
	Desc string
	IsRequired bool
	ErrMissing string
}

type VariableResolver struct {
	Defs []VariableDef
}

func NewVariableResolver() *VariableResolver {
    return &VariableResolver{}
}

func (r *VariableResolver) Define(val *string, flag string, env string, desc string, isRequired bool, errMissing string) *VariableResolver {
	r.Defs = append(r.Defs, VariableDef{
		Val: val,
		Flag: flag,
		Env: env,
		Desc: desc,
		IsRequired: isRequired,
		ErrMissing: errMissing,
	})
	return r
}

func (r *VariableResolver) Resolve() {
	for _, def := range r.Defs {
		flag.StringVar(def.Val, def.Flag, "", def.Desc)
	}

	flag.Parse()

	for _, def := range r.Defs {
		if *def.Val == "" {
			v := os.Getenv(def.Env)
			def.Val = &v
		}
	}
}
