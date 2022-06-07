package main
import (
  "github.com/urfave/cli/v2"
  "github.com/go-git/go-git/v5"
)

func cmdAddAll() *cli.Command {

  return &cli.Command{
    Name:  "add_all",
    Usage: "add_all",
    Action: func(c *cli.Context) error {

      directory := c.Args().Get(0)

      r, err := git.PlainOpen(directory)
      CheckIfError(err)

      w, err := r.Worktree()
      CheckIfError(err)

      Info("git add all new files")
      _, err = w.Add(".")
      CheckIfError(err)

      return nil
    },

  }
}


