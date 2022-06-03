package main
import (
  "github.com/urfave/cli/v2"
  "gopkg.in/src-d/go-git.v4"

)

func cmdAddAll() *cli.Command {

  return &cli.Command{
    Name:  "all_add",
    Usage: "all_add",
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


