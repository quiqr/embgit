package main

import (
  "fmt"
  "io/ioutil"
  "log"
  "os"
  "sort"
  "strings"

  "github.com/urfave/cli/v2"
  "github.com/go-git/go-git/v5/plumbing/transport"

  "golang.org/x/crypto/ssh"
  ssh2 "github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

const version = "v0.5.1"

func setAuth(keyfilepath string, ignoreHostkey bool) transport.AuthMethod {

  if keyfilepath != "" {

    pem, _ := ioutil.ReadFile(keyfilepath)
    signer, err := ssh.ParsePrivateKey(pem)
    if err != nil {
      fmt.Println("Could not read keyfile")
      os.Exit(1)
    }
    auth := &ssh2.PublicKeys{User: "git", Signer: signer}
    if ignoreHostkey {
      auth.HostKeyCallback = ssh.InsecureIgnoreHostKey()
    }
    return auth

  } else {
    return nil
  }

}

func main() {

  app := &cli.App{
    Commands: []*cli.Command{
      {
        Name:  "version",
        Usage: "display version",
        Action: func(c *cli.Context) error {
          fmt.Printf("embgit %s\n", version)
          fmt.Printf("Copyright Quiqr. 2020-2022\n")
          return nil
        },
      },
    },
  }

  app.Commands = append(app.Commands,
    cmdAddAll(),
    cmdClone(),
    cmdCommit(),
    cmdPull(),
    cmdPush(),
    cmdResetHard(),
    cmdFingerprint(),
    cmdKeyGen(),
    cmdKeyGenEcdsa(),
    cmdRepoShowHugotheme(),
    cmdRepoShowQuiqrsite())

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


