package main
import (
  "github.com/urfave/cli/v2"
  "github.com/go-git/go-git/v5"
  "github.com/go-git/go-git/v5/config"
  "github.com/go-git/go-git/v5/storage/memory"
  "log"
)

func cmdLsRemote() *cli.Command {

  return &cli.Command{
    Name: "ls_remote",
    /*
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
    */
    Usage: "ls_remote",
    Action: func(c *cli.Context) error {

      url := c.Args().Get(0)


      // Create the remote with repository URL
      rem := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
        Name: "origin",
        URLs: []string{url},
      })

      log.Print("Fetching tags...")

      // We can then use every Remote functions to retrieve wanted information
      refs, err := rem.List(&git.ListOptions{})
      if err != nil {
        log.Fatal(err)
      }

      // Filters the references list and only keeps tags
      var tags []string
      for _, ref := range refs {
        if ref.Name().IsTag() {
          tags = append(tags, ref.Name().Short())
        }
      }

      log.Printf("Refs found: %v", refs)

      if len(tags) == 0 {
        log.Println("No tags!")
        return nil
      }



      return nil
    },

  }
}

