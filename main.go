package main
import (
    "fmt"
    "os"
    "os/user"
    "myMonkey/repl"
)

func main() {
    user, err := user.Current()
    if err != nil {
        panic(err)
    }

    fmt.Printf("Hello %s! This is the Monkey programming language!\n",
        user.Username)
    fmt.Printf("Feel free to type in commands\n")
    // repl.Evaluate(os.Stdin, os.Stdout)
    repl.VM(os.Stdin, os.Stdout)
}

