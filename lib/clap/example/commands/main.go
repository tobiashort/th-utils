package main

import (
	"fmt"

	"github.com/tobiashort/th-utils/lib/clap"
)

type Args struct {
	Command any `clap:"cmd,mandatory,desc='The cmd to run'"`

	List any `clap:"cmdopt,desc='List all members'"`

	Add struct {
		Name string `clap:"positional,mandatory"`
	} `clap:"cmdopt,desc='Adds a member'"`

	Remove struct {
		Name string `clap:"positional,mandatory"`
	} `clap:"cmdopt,desc='Removes a member'"`

	Foo any `clap:"cmdopt,desc='The foo cmd'"`
	Bar any `clap:"cmdopt,desc='The bar cmd'"`

	Baz struct {
		Command any `clap:"cmd,mandatory,desc='The cmd to run'"`
		Boom    any `clap:"cmdopt,desc='The boom cmd'"`
	} `clap:"cmdopt,desc='The baz cmd'"`
}

func main() {
	args := Args{}
	clap.Parse(&args)

	switch args.Command {
	case &args.List:
		fmt.Println("1: Alice")
		fmt.Println("2: Bob")
		fmt.Println("3: Chris")
	case &args.Add:
		fmt.Println("Added " + args.Add.Name)
	case &args.Remove:
		fmt.Println("Removed " + args.Remove.Name)
	case &args.Foo:
		fmt.Println("foo")
	case &args.Bar:
		fmt.Println("bar")
	case &args.Baz:
		switch args.Baz.Command {
		case &args.Baz.Boom:
			fmt.Println("boom!")
		default:
			fmt.Printf("%v\n", &args.Baz.Command)
			fmt.Printf("%v\n", &args.Baz.Boom)
			panic("baz: unreachable")
		}
	default:
		panic("unreachable")
	}
}
