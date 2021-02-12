package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"

	"golang.org/x/crypto/ssh"
	ssh2 "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"

	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)
const version = "v0.2.7"

func setAuth(keyfilepath string, ignoreHostkey bool) transport.AuthMethod {
	//var auth transport.AuthMethod

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
				Name:  "clone",
				Usage: "complete a task on the list",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "ssh-key",
						Aliases: []string{"i"},
						Usage:   "alternative ssh-key from `FILE`",
					},
					&cli.BoolFlag{Name: "insecure", Aliases: []string{"s"}},
				},
				Action: func(c *cli.Context) error {

					url := c.Args().Get(0)
					directory := c.Args().Get(1)
					auth := setAuth(c.String("ssh-key"), c.Bool("insecure"))

					Info("git clone %s %s", url, directory)

					_, err2 := git.PlainClone(directory, false, &git.CloneOptions{
						URL:      url,
						Progress: os.Stdout,
						Auth:     auth,
					})

					CheckIfError(err2)

					return nil
				},
			},
			{
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
			},
			{
				Name:  "alladd",
				Usage: "alladd",
				Action: func(c *cli.Context) error {

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
				Name:  "push",
				Usage: "push to remote",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "ssh-key",
						Aliases: []string{"i"},
						Usage:   "alternative ssh-key from `FILE`",
					},
					&cli.BoolFlag{Name: "insecure", Aliases: []string{"s"}},
				},
				Action: func(c *cli.Context) error {
					directory := c.Args().Get(0)
					auth := setAuth(c.String("ssh-key"), c.Bool("insecure"))

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
			{
				Name:  "keygen",
				Usage: "create passwordless ssl key pair",
				//        Flags: []cli.Flag{
				//          &cli.StringFlag{
				//            Name:  "filename",
				//            Aliases: []string{"f"},
				//            Usage: "alternative ssh-key from `FILE`",
				//          }
				//        },
				Action: func(c *cli.Context) error {
					//directory := c.Args().Get(0)
					//filename := c.String("filename"))

					savePrivateFileTo := "./id_rsa_pogo"
					savePublicFileTo := "./id_rsa_pogo.pub"
					bitSize := 4096

					privateKey, err := generatePrivateKey(bitSize)
					publicKeyBytes, err := generatePublicKey(&privateKey.PublicKey)
					privateKeyBytes := encodePrivateKeyToPEM(privateKey)

					err = writeKeyToFile(privateKeyBytes, savePrivateFileTo)
					err = writeKeyToFile([]byte(publicKeyBytes), savePublicFileTo)
					CheckIfError(err)

					return nil
				},
			},
			{
				Name:  "fingerprint",
				Usage: "get fingerprint from ssl key pair",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "ssh-key",
						Aliases: []string{"i"},
						Usage:   "alternative ssh-key from `FILE`",
					},
				},
				Action: func(c *cli.Context) error {

					keyfilepath := c.String("ssh-key")
					priv, err := ioutil.ReadFile(keyfilepath)

					block, _ := pem.Decode([]byte(priv))
					if block == nil || block.Type != "RSA PRIVATE KEY" {
						log.Fatal("failed to decode PEM block containing public key")
					}
					key, err := x509.ParsePKCS1PrivateKey(block.Bytes)

					publicKeyBytes, err := generatePublicKey(&key.PublicKey)

					ppubk,_,_,_,err := ssh.ParseAuthorizedKey(publicKeyBytes);
					f := ssh.FingerprintLegacyMD5(ppubk)
					fmt.Printf("%s\n", f)

					if err != nil {
						log.Fatal(err)
					}

					return nil
				},
			},
			{
				Name:  "version",
				Usage: "display version",
				Action: func(c *cli.Context) error {
					Info("embgit %s, Copyright PoppyGo B.V. 2020-2021, www.poppygo.io", version)
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

// generatePrivateKey creates a RSA Private Key of specified byte size
func generatePrivateKey(bitSize int) (*rsa.PrivateKey, error) {
	// Private Key generation
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, err
	}

	// Validate Private Key
	err = privateKey.Validate()
	if err != nil {
		return nil, err
	}

	log.Println("Private Key generated")
	return privateKey, nil
}

// encodePrivateKeyToPEM encodes Private Key from RSA to PEM format
func encodePrivateKeyToPEM(privateKey *rsa.PrivateKey) []byte {
	// Get ASN.1 DER format
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)

	// pem.Block
	privBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privDER,
	}

	// Private key in PEM format
	privatePEM := pem.EncodeToMemory(&privBlock)

	return privatePEM
}

// generatePublicKey take a rsa.PublicKey and return bytes suitable for writing to .pub file
// returns in the format "ssh-rsa ..."
func generatePublicKey(privatekey *rsa.PublicKey) ([]byte, error) {
	publicRsaKey, err := ssh.NewPublicKey(privatekey)
	if err != nil {
		return nil, err
	}

	pubKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)

	log.Println("Public key generated")
	return pubKeyBytes, nil
}

// writePemToFile writes keys to a file
func writeKeyToFile(keyBytes []byte, saveFileTo string) error {
	err := ioutil.WriteFile(saveFileTo, keyBytes, 0600)
	if err != nil {
		return err
	}

	log.Printf("Key saved to: %s", saveFileTo)
	return nil
}
