package main
import (
  "fmt"
  "os"
  "time"
  "github.com/urfave/cli/v2"
  "github.com/go-git/go-git/v5"

  "github.com/go-git/go-git/v5/plumbing/object"
  "github.com/go-git/go-git/v5/storage/memory"
)

func cmdLog() *cli.Command {

  return &cli.Command{
    Name:  "log",
    Usage: "show logs from remote",
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

      url := c.Args().Get(0)
      directory := c.Args().Get(1)
      auth := setAuth(c.String("ssh-key"), c.Bool("insecure"))
      //branch := c.String("branch")

      Info("git log %s %s", url, directory)

      r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
        URL:      url,
        Progress: os.Stdout,
        Auth:     auth,
      })

      CheckIfError(err)

      // Gets the HEAD history from HEAD, just like this command:
      Info("git log")

      // ... retrieves the branch pointed by HEAD
      ref, err := r.Head()
      CheckIfError(err)

      // ... retrieves the commit history
      since := time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)
      until := time.Date(2019, 7, 30, 0, 0, 0, 0, time.UTC)
      cIter, err := r.Log(&git.LogOptions{From: ref.Hash(), Since: &since, Until: &until})
      CheckIfError(err)

      // ... just iterates over the commits, printing it
      err = cIter.ForEach(func(c *object.Commit) error {
        fmt.Println(c)

        return nil
      })
      CheckIfError(err)

      return nil
    },

  }
}


