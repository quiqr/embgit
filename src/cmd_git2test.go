package main
import (
  "github.com/urfave/cli/v2"

//  "reflect"
  "fmt"
  git2go "github.com/libgit2/git2go/v34"
)

var (
  Git2GoCommitRevWalk *git2go.Commit
)


func cmdGit2test() *cli.Command {

  return &cli.Command{
    Name:  "git2test",
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


      repo, err := git2go.OpenRepository(".")
      if err != nil {
        return err
      }
      defer repo.Free()

      headRef, err := repo.Head()
      if err != nil {
        return err
      }
      defer headRef.Free()

      revWalk, err := repo.Walk()
      if err != nil {
        return err
      }
      defer revWalk.Free()

      if err := revWalk.Push(headRef.Target()); err != nil {
        return err
      }

      revWalk.Sorting(git2go.SortTime)

      count := 0
      if err := revWalk.Iterate(func(commit *git2go.Commit) bool {
        defer commit.Free()
        count++
				Git2GoCommitRevWalk = commit
//				fmt.Println("type:", reflect.TypeOf(commit))
        fmt.Println(commit.Author().Name)
        fmt.Println(commit.Id().String())
        fmt.Println(commit.Message())
        //fmt.Println(commit.ExtractSignature())
        fmt.Println(commit.Committer().When)
        return true
      }); err != nil {
        return err
      }

      return nil

    },

  }
}


