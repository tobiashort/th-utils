package main

import (
	"fmt"
	"os"

	"github.com/tobiashort/th-utils/lib/ansi"
	"github.com/tobiashort/th-utils/lib/cfmt"
)

func main() {
	cfmt.Print("This ", "is ", "#r{red}\n")
	cfmt.Print("This ", "is ", "#g{green}\n")
	cfmt.Print("This ", "is ", "#y{yellow}\n")
	cfmt.Print("This ", "is ", "#b{blue}\n")
	cfmt.Print("This ", "is ", "#p{purple}\n")
	cfmt.Print("This ", "is ", "#c{cyan}\n")

	fmt.Printf("\n")

	cfmt.Printf("This is #r{%s}\n", "red")
	cfmt.Printf("This is #g{%s}\n", "green")
	cfmt.Printf("This is #y{%s}\n", "yellow")
	cfmt.Printf("This is #b{%s}\n", "blue")
	cfmt.Printf("This is #p{%s}\n", "purle")
	cfmt.Printf("This is #c{%s}\n", "cyan")

	fmt.Printf("\n")

	cfmt.Println("This", "is", "#r{red}")
	cfmt.Println("This", "is", "#g{green}")
	cfmt.Println("This", "is", "#y{yellow}")
	cfmt.Println("This", "is", "#b{blue}")
	cfmt.Println("This", "is", "#p{purple}")
	cfmt.Println("This", "is", "#c{cyan}")

	fmt.Printf("\n")

	cfmt.Cprint("c", "This ", "is ", "#r{red}\n")
	cfmt.Cprint("p", "This ", "is ", "#g{green}\n")
	cfmt.Cprint("b", "This ", "is ", "#y{yellow}\n")
	cfmt.Cprint("y", "This ", "is ", "#b{blue}\n")
	cfmt.Cprint("g", "This ", "is ", "#p{purple}\n")
	cfmt.Cprint("r", "This ", "is ", "#c{cyan}\n")

	fmt.Printf("\n")

	cfmt.Cprintf("c", "This is #r{red}\n")
	cfmt.Cprintf("p", "This is #g{green}\n")
	cfmt.Cprintf("b", "This is #y{yellow}\n")
	cfmt.Cprintf("y", "This is #b{blue}\n")
	cfmt.Cprintf("g", "This is #p{purple}\n")
	cfmt.Cprintf("r", "This is #c{cyan}\n")

	fmt.Printf("\n")

	cfmt.Cfprintf(os.Stdout, "c", "This is #r{red}\n")
	cfmt.Cfprintf(os.Stdout, "p", "This is #g{green}\n")
	cfmt.Cfprintf(os.Stdout, "b", "This is #y{yellow}\n")
	cfmt.Cfprintf(os.Stdout, "y", "This is #b{blue}\n")
	cfmt.Cfprintf(os.Stdout, "g", "This is #p{purple}\n")
	cfmt.Cfprintf(os.Stdout, "r", "This is #c{cyan}\n")

	fmt.Printf("\n")

	cfmt.Cprintln("c", "This is #r{red}")
	cfmt.Cprintln("p", "This is #g{green}")
	cfmt.Cprintln("b", "This is #y{yellow}")
	cfmt.Cprintln("y", "This is #b{blue}")
	cfmt.Cprintln("g", "This is #p{purple}")
	cfmt.Cprintln("r", "This is #c{cyan}")

	fmt.Printf("\n")

	cfmt.Fprint(os.Stdout, "This ", "is ", "#r{red}\n")
	cfmt.Fprint(os.Stdout, "This ", "is ", "#g{green}\n")
	cfmt.Fprint(os.Stdout, "This ", "is ", "#y{yellow}\n")
	cfmt.Fprint(os.Stdout, "This ", "is ", "#b{blue}\n")
	cfmt.Fprint(os.Stdout, "This ", "is ", "#p{purple}\n")
	cfmt.Fprint(os.Stdout, "This ", "is ", "#c{cyan}\n")

	fmt.Printf("\n")

	cfmt.Fprintf(os.Stdout, "This is #r{%s}\n", "red")
	cfmt.Fprintf(os.Stdout, "This is #g{%s}\n", "green")
	cfmt.Fprintf(os.Stdout, "This is #y{%s}\n", "yellow")
	cfmt.Fprintf(os.Stdout, "This is #b{%s}\n", "blue")
	cfmt.Fprintf(os.Stdout, "This is #p{%s}\n", "purle")
	cfmt.Fprintf(os.Stdout, "This is #c{%s}\n", "cyan")

	fmt.Printf("\n")

	cfmt.Fprintln(os.Stdout, "This", "is", "#r{red}")
	cfmt.Fprintln(os.Stdout, "This", "is", "#g{green}")
	cfmt.Fprintln(os.Stdout, "This", "is", "#y{yellow}")
	cfmt.Fprintln(os.Stdout, "This", "is", "#b{blue}")
	cfmt.Fprintln(os.Stdout, "This", "is", "#p{purple}")
	cfmt.Fprintln(os.Stdout, "This", "is", "#c{cyan}")

	fmt.Printf("\n")

	fmt.Print(cfmt.Sprint("This ", "is ", "#r{red}\n"))
	fmt.Print(cfmt.Sprint("This ", "is ", "#g{green}\n"))
	fmt.Print(cfmt.Sprint("This ", "is ", "#y{yellow}\n"))
	fmt.Print(cfmt.Sprint("This ", "is ", "#b{blue}\n"))
	fmt.Print(cfmt.Sprint("This ", "is ", "#p{purple}\n"))
	fmt.Print(cfmt.Sprint("This ", "is ", "#c{cyan}\n"))

	fmt.Printf("\n")

	fmt.Print(cfmt.Sprintf("This is #r{%s}\n", "red"))
	fmt.Print(cfmt.Sprintf("This is #g{%s}\n", "green"))
	fmt.Print(cfmt.Sprintf("This is #y{%s}\n", "yellow"))
	fmt.Print(cfmt.Sprintf("This is #b{%s}\n", "blue"))
	fmt.Print(cfmt.Sprintf("This is #p{%s}\n", "purle"))
	fmt.Print(cfmt.Sprintf("This is #c{%s}\n", "cyan"))

	fmt.Printf("\n")

	fmt.Print(cfmt.Sprintln("This", "is", "#r{red}"))
	fmt.Print(cfmt.Sprintln("This", "is", "#g{green}"))
	fmt.Print(cfmt.Sprintln("This", "is", "#y{yellow}"))
	fmt.Print(cfmt.Sprintln("This", "is", "#b{blue}"))
	fmt.Print(cfmt.Sprintln("This", "is", "#p{purple}"))
	fmt.Print(cfmt.Sprintln("This", "is", "#c{cyan}"))

	fmt.Printf("\n")

	cfmt.Println("This", "is", "#rB{red and bold}")
	cfmt.Println("This", "is", "#gU{green and underlined}")
	cfmt.Println("This", "is", "#yR{yellow reversed}")

	fmt.Printf("\n")

	cfmt.Begin(ansi.DecorPurple)
	fmt.Println("what follows now  is...")
	fmt.Println("...in purple.")
	cfmt.End()
	fmt.Println("and now back to normal")

	fmt.Printf("\n")
	cfmt.Println(`#y{\{\{\escaped\\\}\}}`)
}
