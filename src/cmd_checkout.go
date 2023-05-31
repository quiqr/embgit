package main

import (
  "github.com/urfave/cli/v2"
  "github.com/go-git/go-git/v5"
  "github.com/go-git/go-git/v5/plumbing"
)

func cmdCheckout() *cli.Command {

  return &cli.Command{
    Name:  "checkout",
    Usage: "checkout a commit",
    Flags: []cli.Flag{
      &cli.StringFlag{
        Name:    "ref",
        Aliases: []string{"r"},
        Usage:   "commit reference to checkout",
      },
    },
    Action: func(c *cli.Context) error {

      directory := c.Args().Get(1)
      ref := c.String("ref")

      r, err := git.PlainOpen(directory)
      CheckIfError(err)

      w, err := r.Worktree()
      CheckIfError(err)

      Info("git checkout %s", ref)
      err = w.Checkout(&git.CheckoutOptions{
        Hash: plumbing.NewHash(ref),
      })
      CheckIfError(err)

      return nil
    },
  }
}

