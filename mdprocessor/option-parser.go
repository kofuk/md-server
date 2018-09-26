package mdprocessor

import (
    "io"
)

type mdInput struct {
    Input *io.Reader
    Output *io.Writer
}
/*
func parseOptions() *[]mdInput {
    var inputs []mdInput
    input := mdInput{}
    var err error
    input.Input, err = os.Open(os.Args[1])
    if err != nil {
        log.Println("Error opening", os.Args[1], "; Ignoring ...")
        panic(err)
    }
    input.Output, err =
        os.OpenFile(os.Args[2], os.O_WRONLY | os.O_CREATE, 0666)
    if err != nil {
        log.Println("Error opening", os.Args[2], "for output; Ignoring input")
        input.Input.Close()
        panic(err)
    }
    inputs = append(inputs, input)
    return &inputs
}

func closeAllInput(inputs *[]mdInput) {
    for _, input := range *inputs {
        input.Input.Close()
        input.Output.Close()
    }
}

const USAGE =
`Usage:  mdc [FILE]...
With no FILE, or when FILE is only -, read standard input.

  --help     print this help and exit
  --version  print version information and exit
`

const VERSION =
`Markdown compiler 1.0
Copyright (C) 2018 KoFuk
`

func printHelpAndExit() {
    fmt.Println(USAGE)
    os.Exit(1)
}

func printVersionAndExit() {
    fmt.Println(VERSION)
    os.Exit(1)
}
*/
