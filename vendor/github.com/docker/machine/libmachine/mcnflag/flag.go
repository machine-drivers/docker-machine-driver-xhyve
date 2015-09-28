package mcnflag

type Flag struct {
	EnvVar string
	Name   string
	Usage  string
	Value  interface{}
}
