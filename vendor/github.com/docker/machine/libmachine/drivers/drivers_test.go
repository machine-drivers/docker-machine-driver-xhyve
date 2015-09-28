package drivers

import (
	"testing"

	"github.com/docker/machine/libmachine/mcnflag"
)

func TestGetCreateFlags(t *testing.T) {
	Register("foo", &RegisteredDriver{
		GetCreateFlags: func() []mcnflag.Flag {
			return []mcnflag.Flag{
				mcnflag.Flag{
					Name:   "a",
					Value:  "",
					Usage:  "",
					EnvVar: "",
				},
				mcnflag.Flag{
					Name:   "b",
					Value:  "",
					Usage:  "",
					EnvVar: "",
				},
				mcnflag.Flag{
					Name:   "c",
					Value:  "",
					Usage:  "",
					EnvVar: "",
				},
			}
		},
	})
	Register("bar", &RegisteredDriver{
		GetCreateFlags: func() []mcnflag.Flag {
			return []mcnflag.Flag{
				mcnflag.Flag{
					Name:   "d",
					Value:  "",
					Usage:  "",
					EnvVar: "",
				},
				mcnflag.Flag{
					Name:   "e",
					Value:  "",
					Usage:  "",
					EnvVar: "",
				},
				mcnflag.Flag{
					Name:   "f",
					Value:  "",
					Usage:  "",
					EnvVar: "",
				},
			}
		},
	})

	expected := []string{"-a \t", "-b \t", "-c \t", "-d \t", "-e \t", "-f \t"}

	// test a few times to catch offset issue
	// if it crops up again
	for i := 0; i < 5; i++ {
		flags := GetCreateFlags()
		for j, e := range expected {
			if flags[j].String() != e {
				t.Fatal("Flags are out of order")
			}
		}
	}
}
