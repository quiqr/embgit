package main
import (
  "github.com/urfave/cli/v2"
  "gopkg.in/src-d/go-git.v4"

)

func cmdResetHard() *cli.Command {

  return &cli.Command{
    Name: "reset_hard",
    Usage: "commit",
    Action: func(c *cli.Context) error {

      directory := c.Args().Get(0)

      Info("Hard reset in %s", directory)
      r, err := git.PlainOpen(directory)
      CheckIfError(err)

      w, err := r.Worktree()
      CheckIfError(err)

      head, err := r.Head()
      if err != nil {
        return err
      }

      if err := w.Reset(&git.ResetOptions{
        Mode:   git.HardReset,
        Commit: head.Hash(),
      }); err != nil {
        return err
      }

      //          commit, err := w.Commit(c.String("message"), &git.CommitOptions{
      //All: c.Bool("all"),
      //Author: &object.Signature{
      //Name:  c.String("author-name"),
      //Email: c.String("author-email"),
      //When:  time.Now(),
      //},
      //})

      CheckIfError(err)

      return nil
    },
  }
}



