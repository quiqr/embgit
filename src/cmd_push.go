package main
import (
  "github.com/urfave/cli/v2"
  "github.com/go-git/go-git/v5"

)

func cmdPush() *cli.Command {

  return &cli.Command{
    Name:  "push",
    Usage: "push to remote",
    Flags: []cli.Flag{
      &cli.StringFlag{
        Name:    "ssh-key",
        Aliases: []string{"i"},
        Usage:   "alternative ssh-key from `FILE`",
      },
      &cli.BoolFlag{
        Name: "insecure",
        Aliases: []string{"s"},
        Usage:   "skip host key validation",
      },
    },
    Action: func(c *cli.Context) error {
      directory := c.Args().Get(0)
      auth := setAuth(c.String("ssh-key"), c.Bool("insecure"))

      r, err := git.PlainOpen(directory)
      CheckIfError(err)

      Info("git push")
      // push using default options
      err = r.Push(&git.PushOptions{
        Auth: auth,
      })
      CheckIfError(err)

      return nil
    },

  }
}


