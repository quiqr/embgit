package main

import (
  "log"

  "github.com/urfave/cli/v2"
  "golang.org/x/crypto/ssh"
  //ssh2 "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"


  "fmt"
  "io/ioutil"
  "crypto/rand"
  "crypto/rsa"
  "crypto/x509"
  "crypto/ecdsa"
  "crypto/elliptic"

  "encoding/pem"
)


func cmdFingerprint() *cli.Command {

  return &cli.Command{
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
      if block != nil {
        if block.Type == "RSA PRIVATE KEY" {
          key, err := x509.ParsePKCS1PrivateKey(block.Bytes)

          publicKeyBytes, err := generatePublicKey(&key.PublicKey)
          if err != nil {
            log.Fatal("Failed to parse private key")
          }

          ppubk,_,_,_,err := ssh.ParseAuthorizedKey(publicKeyBytes);
          f := ssh.FingerprintLegacyMD5(ppubk)
          fmt.Printf("%s\n", f)

        } else if block.Type == "OPENSSH PRIVATE KEY" {

          key, err := ssh.ParsePrivateKey(priv)
          if err != nil {
            log.Fatal("Failed to parse private key")
          }
          publicKeyBytes := ssh.MarshalAuthorizedKey(key.PublicKey())

          ppubk,_,_,_,err := ssh.ParseAuthorizedKey(publicKeyBytes);
          f := ssh.FingerprintLegacyMD5(ppubk)
          fmt.Printf("%s\n", f)
        } else{
          log.Fatal("failed to decode PEM block containing public key")
        }
      } else{
        log.Fatal("nothing to decode")
      }
      if err != nil {
        log.Fatal(err)
      }

      return nil
    },

  }

}

func cmdKeyGen() *cli.Command {

  return &cli.Command{
    Name:  "keygen",
    Usage: "create passwordless ssl key pair",
    Action: func(c *cli.Context) error {
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
  }
}

func cmdKeyGenEcdsa() *cli.Command {
  return &cli.Command{
    Name:  "keygen_ecdsa",
    Usage: "create passwordless ecdsa ssl key pair",
    Action: func(c *cli.Context) error {
      savePrivateFileTo := "./id_ecdsa_quiqr"
      savePublicFileTo := "./id_ecdsa_quiqr.pub"
      //bitSize := 4096

      //privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
      //pubKeyBytes := ssh.MarshalAuthorizedKey(privateKey.PublicKey)

      privateKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
      publicKey := &privateKey.PublicKey
      privateKeyBytes, _ := encodeECDSAToPEM(privateKey, publicKey)
      publicKeyBytes, err := generatePublicKeyECDSA(&privateKey.PublicKey)

      err = writeKeyToFile(privateKeyBytes, savePrivateFileTo)
      err = writeKeyToFile([]byte(publicKeyBytes), savePublicFileTo)

      CheckIfError(err)

      return nil
    },
  }
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


func encodeECDSAToPEM(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) ([]byte, []byte) {
  x509Encoded, _ := x509.MarshalECPrivateKey(privateKey)
  pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: x509Encoded})

  x509EncodedPub, _ := x509.MarshalPKIXPublicKey(publicKey)
  pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "EC PUBLIC KEY", Bytes: x509EncodedPub})

  return pemEncoded, pemEncodedPub
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

func generatePublicKeyECDSA(privatekey *ecdsa.PublicKey) ([]byte, error) {
  publicKey, err := ssh.NewPublicKey(privatekey)
  if err != nil {
    return nil, err
  }

  pubKeyBytes := ssh.MarshalAuthorizedKey(publicKey)

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
