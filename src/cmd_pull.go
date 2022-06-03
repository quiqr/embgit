package main
import (
  "github.com/urfave/cli/v2"
  "gopkg.in/src-d/go-git.v4"

)

func cmdPull() *cli.Command {

  return &cli.Command{
    Name:  "pull",
    Usage: "pull from remote",
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

      w, err := r.Worktree()
      CheckIfError(err)

      Info("git pull")
      // pull using default options
      err = w.Pull(&git.PullOptions{
        Auth: auth,
      })
      CheckIfError(err)

      return nil
    },

  }
}


