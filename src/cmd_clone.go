package main

import (
  "github.com/urfave/cli/v2"
  "github.com/go-git/go-git/v5"

  "os"
)

func cmdClone() *cli.Command {

  return &cli.Command{
    Name:  "clone",
    Usage: "complete a task on the list",
    Flags: []cli.Flag{
      &cli.StringFlag{
        Name:    "ssh-key",
        Aliases: []string{"i"},
        Usage:   "alternative ssh-key from `FILE`",
      },
      &cli.StringFlag{
        Name:    "branch",
        Aliases: []string{"b"},
        Usage:   "other then default branch",
      },
      &cli.BoolFlag{
        Name: "insecure",
        Aliases: []string{"s"},
        Usage:   "skip host key validation",
      },
    },
    Action: func(c *cli.Context) error {

      url := c.Args().Get(0)
      directory := c.Args().Get(1)
      auth := setAuth(c.String("ssh-key"), c.Bool("insecure"))
      //branch := c.String("branch")

      Info("git clone %s %s", url, directory)

      _, err2 := git.PlainClone(directory, false, &git.CloneOptions{
        URL:      url,
        Progress: os.Stdout,
        Auth:     auth,
      })

      CheckIfError(err2)

      return nil
    },
  }
}

