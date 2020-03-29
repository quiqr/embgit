package main

import (
  "fmt"
  "log"
  "os"
  "strings"
  "sort"
  "io/ioutil"
  "time"

  "github.com/urfave/cli/v2"
  "gopkg.in/src-d/go-git.v4"
  "gopkg.in/src-d/go-git.v4/plumbing/transport"
  "gopkg.in/src-d/go-git.v4/plumbing/object"

  "golang.org/x/crypto/ssh"
  ssh2 "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)
func setAuth(keyfilepath string) transport.AuthMethod {
  //var auth transport.AuthMethod

  if keyfilepath != "" {

    pem, _ := ioutil.ReadFile(keyfilepath)
    signer, err := ssh.ParsePrivateKey(pem)
    if err != nil {
      fmt.Println("Could not read keyfile")
      os.Exit(1);
    }
    return &ssh2.PublicKeys{User: "git", Signer: signer}

  } else {
    return nil
  }

}

func main() {

  app := &cli.App{
    Commands: []*cli.Command{
      {
        Name:    "clone",
        Usage:   "complete a task on the list",
        Flags: []cli.Flag{
          &cli.StringFlag{
            Name:  "ssh-key",
            Aliases: []string{"i"},
            Usage: "alternative ssh-key from `FILE`",
          },
        },
        Action:  func(c *cli.Context) error {

          url := c.Args().Get(0)
          directory := c.Args().Get(1)
          auth := setAuth(c.String("ssh-key"))

          Info("git clone %s %s", url, directory)

          _, err2 := git.PlainClone(directory, false, &git.CloneOptions{
            URL:      url,
            Progress: os.Stdout,
            Auth: auth,
          })

          CheckIfError(err2)

          return nil
        },
      },
      {
        Name:    "commit",
        Flags: []cli.Flag{
          &cli.StringFlag{
            Name:  "message",
            Aliases: []string{"m"},
            Usage: "commit message `MESSAGE`",
            Required: true,
          },
          &cli.StringFlag{
            Name:  "author-email",
            Aliases: []string{"e"},
            Usage: "commit author `EMAIL`",
            Required: true,
          },
          &cli.StringFlag{
            Name:  "author-name",
            Aliases: []string{"n"},
            Usage: "commit author `NAME`",
            Required: true,
          },
        },
        Usage:   "commit",
        Action:  func(c *cli.Context) error {

          directory := c.Args().Get(0)

          Info("commit with message: %s in %s", c.String("message"), directory)
          r, err := git.PlainOpen(directory)
          CheckIfError(err)

          w, err := r.Worktree()
          CheckIfError(err)

          commit, err := w.Commit(c.String("message"), &git.CommitOptions{
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
      },
      {
        Name:    "alladd",
        Usage:   "alladd",
        Action:  func(c *cli.Context) error {

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
      },
      {
        Name:    "push",
        Usage:   "push to remote",
        Flags: []cli.Flag{
          &cli.StringFlag{
            Name:  "ssh-key",
            Aliases: []string{"i"},
            Usage: "alternative ssh-key from `FILE`",
          },
        },
        Action:  func(c *cli.Context) error {
          directory := c.Args().Get(0)
          auth := setAuth(c.String("ssh-key"))

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
      },
    },
  }

  sort.Sort(cli.FlagsByName(app.Flags))
  sort.Sort(cli.CommandsByName(app.Commands))

  err := app.Run(os.Args)
  if err != nil {
    log.Fatal(err)
  }
}

// CheckArgs should be used to ensure the right command line arguments are
// passed before executing an example.
func CheckArgs(arg ...string) {
  if len(os.Args) < len(arg)+1 {
    Warning("Usage: %s %s", os.Args[0], strings.Join(arg, " "))
    os.Exit(1)
  }
}

// CheckIfError should be used to naively panics if an error is not nil.
func CheckIfError(err error) {
  if err == nil {
    return
  }

  fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
  os.Exit(1)
}

// Info should be used to describe the example commands that are about to run.
func Info(format string, args ...interface{}) {
  fmt.Printf("\x1b[34;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

// Warning should be used to display a warning
func Warning(format string, args ...interface{}) {
  fmt.Printf("\x1b[36;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}
