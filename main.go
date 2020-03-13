package main

import (
  "fmt"
  "log"
  "os"
  "strings"
  "sort"

  "github.com/urfave/cli/v2"
  "gopkg.in/src-d/go-git.v4"
)

//clone
//add
//commit
//push
//ssh -i

func main() {

  app := &cli.App{
    Flags: []cli.Flag{
      &cli.StringFlag{
        Name:  "lang, l",
        Value: "english",
        Usage: "Language for the greeting",
      },
      &cli.StringFlag{
        Name:  "config, c",
        Usage: "Load configuration from `FILE`",
      },
    },
    Commands: []*cli.Command{
      {
        Name:    "clone",
        Aliases: []string{"c"},
        Usage:   "complete a task on the list",
        Action:  func(c *cli.Context) error {

          Info("git clone https://github.com/src-d/go-git")

          _, err := git.PlainClone("/tmp/foo", false, &git.CloneOptions{
            URL:      "https://github.com/src-d/go-git",
            Progress: os.Stdout,
          })

          CheckIfError(err)

          return nil
        },
      },
      {
        Name:    "add",
        Aliases: []string{"a"},
        Usage:   "add a task to the list",
        Action:  func(c *cli.Context) error {
          return nil
        },
      },
    },
  }

  sort.Sort(cli.FlagsByName(app.Flags))
  sort.Sort(cli.CommandsByName(app.Commands))

  err := app.Run(os.Args)
  if err != nil {
    log.Fatal(err)
  }
}

// CheckArgs should be used to ensure the right command line arguments are
// passed before executing an example.
func CheckArgs(arg ...string) {
	if len(os.Args) < len(arg)+1 {
		Warning("Usage: %s %s", os.Args[0], strings.Join(arg, " "))
		os.Exit(1)
	}
}

// CheckIfError should be used to naively panics if an error is not nil.
func CheckIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}

// Info should be used to describe the example commands that are about to run.
func Info(format string, args ...interface{}) {
	fmt.Printf("\x1b[34;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

// Warning should be used to display a warning
func Warning(format string, args ...interface{}) {
	fmt.Printf("\x1b[36;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}
