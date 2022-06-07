package main
import (
  "fmt"
  "time"
  "github.com/urfave/cli/v2"
  "github.com/go-git/go-git/v5"
  "github.com/go-git/go-git/v5/plumbing/object"
)

func cmdCommit() *cli.Command {

  return &cli.Command{
    Name: "commit",
    Flags: []cli.Flag{
      &cli.StringFlag{
        Name:     "message",
        Aliases:  []string{"m"},
        Usage:    "commit message `MESSAGE`",
        Required: true,
      },
      &cli.StringFlag{
        Name:     "author-email",
        Aliases:  []string{"e"},
        Usage:    "commit author `EMAIL`",
        Required: true,
      },
      &cli.StringFlag{
        Name:     "author-name",
        Aliases:  []string{"n"},
        Usage:    "commit author `NAME`",
        Required: true,
      },
      &cli.BoolFlag{Name: "all", Aliases: []string{"a"}},
    },
    Usage: "commit",
    Action: func(c *cli.Context) error {

      directory := c.Args().Get(0)

      Info("commit with message: %s in %s", c.String("message"), directory)
      r, err := git.PlainOpen(directory)
      CheckIfError(err)

      w, err := r.Worktree()
      CheckIfError(err)

      commit, err := w.Commit(c.String("message"), &git.CommitOptions{
        All: c.Bool("all"),
        Author: &object.Signature{
          Name:  c.String("author-name"),
          Email: c.String("author-email"),
          When:  time.Now(),
        },
      })

      CheckIfError(err)

      Info("git show -s")
      obj, err := r.CommitObject(commit)
      CheckIfError(err)

      fmt.Println(obj)

      return nil
    },

  }
}

