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

      Info("git add all new/mod/deleted files")
      Info("in dir: %s", directory)
      err = w.AddWithOptions(&git.AddOptions{ All: true})
      CheckIfError(err)

      return nil
    },

  }
}


