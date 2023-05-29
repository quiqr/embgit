package main
import (
  "fmt"
  "encoding/json"
  "time"
  "github.com/urfave/cli/v2"
  "github.com/go-git/go-git/v5"
  "github.com/go-git/go-git/v5/plumbing/object"
  "github.com/go-git/go-git/v5/storage/memory"
)

func cmdLogRemote() *cli.Command {

  return &cli.Command{
    Name:  "log_remote",
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
      //directory := c.Args().Get(1)
      auth := setAuth(c.String("ssh-key"), c.Bool("insecure"))
      //branch := c.String("branch")

      r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
        URL:      url,
        Auth:     auth,
      })

      CheckIfError(err)

      ref, err := r.Head()
      CheckIfError(err)

      since := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
      until := time.Date(2099, 1, 30, 0, 0, 0, 0, time.UTC)
      cIter, err := r.Log(&git.LogOptions{From: ref.Hash(), Since: &since, Until: &until})
      CheckIfError(err)

      var jsonCommits []*jsonCommitEntry

      err = cIter.ForEach(func(c *object.Commit) error {

        commitEntry := &jsonCommitEntry{
          Author:  c.Author.String(),
          Message: c.Message,
          Ref: c.Hash.String(),
          Date: c.Author.When.Format(DateFormat),
        }

        jsonCommits = append(jsonCommits, commitEntry)

        return nil
      })
      CheckIfError(err)

      data, _ := json.Marshal(jsonCommits)
      fmt.Println(string(data))

      return nil
    },

  }
}


