package main

import (
  "fmt"
  //"github.com/davecgh/go-spew/spew"
  "path/filepath"
  "github.com/spf13/viper"
  "encoding/json"
  "golang.org/x/exp/slices"
  "strings"
  "bytes"
  "net/http"

  "github.com/urfave/cli/v2"
  "github.com/go-git/go-git/v5"
  "github.com/go-git/go-git/v5/plumbing/object"
  "github.com/go-git/go-git/v5/storage/memory"

)

type responseDictType struct {
  HugoVersion string
  HugoTheme string
  QuiqrFormsEndPoints int
  QuiqrModel string
  Screenshot string
}


func cmdLsRemote() *cli.Command {

  return &cli.Command{
    Name: "quiqr_repo_show",
    Usage: "quiqr_repo_show",
    Action: func(c *cli.Context) error {

      url := c.Args().Get(0)
      showCase(url)
      return nil
    },
  }
}


func showCase(url string){

  r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
    URL: url,
  })

  CheckIfError(err)

  ref, err := r.Head()
  CheckIfError(err)

  commit, err := r.CommitObject(ref.Hash())
  CheckIfError(err)

  tree, err := commit.Tree()
  CheckIfError(err)

  var hugover string;
  var hugotheme string;
  formEndPoints := 0;
  var screenshot1 string;
  var quiqrModel string;

  tree.Files().ForEach(func(f *object.File) error {

    //GET HUGO VERSIONS
    if(strings.HasPrefix(f.Name, "quiqr/model/base.")){

      contents, _ := f.Contents()
      var extension = strings.TrimLeft(filepath.Ext(f.Name), ".")
      viper.SetConfigType(extension)
      viper.ReadConfig(bytes.NewBuffer([]byte(contents)))
      hugover, _ = viper.Get("hugover").(string)
      quiqrModel = "base"
    }

    if(strings.HasPrefix(f.Name, "quiqr/model/includes/singles.")){
      quiqrModel = quiqrModel + " singles"
    }

    if(strings.HasPrefix(f.Name, "quiqr/model/includes/collections.")){
      quiqrModel = quiqrModel + " collections"
    }

    if(strings.HasPrefix(f.Name, "quiqr/model/includes/collections.")){
      quiqrModel = quiqrModel + " menu"
    }

    if(strings.HasPrefix(f.Name, "quiqr/etalage/screenshots/")){

      imgExts := []string{"jpg", "png", "git", "jpeg"}
      extension := strings.ToLower(strings.TrimLeft(filepath.Ext(f.Name), "."))
      if(slices.Contains(imgExts, extension)){
        //spew.Dump(f.Name)
        contents, _ := f.Contents()
        contentsBytes := []byte(contents)


        var base64Encoding string

        // Determine the content type of the image file
        mimeType := http.DetectContentType(contentsBytes)

        // Prepend the appropriate URI scheme header depending
        // on the MIME type
        switch mimeType {
        case "image/jpeg":
          base64Encoding += "data:image/jpeg;base64,"
        case "image/png":
          base64Encoding += "data:image/png;base64,"
        }

        // Append the base64 encoded output
        base64Encoding += toBase64(contentsBytes)

        // Print the full base64 representation of the image
        screenshot1 = base64Encoding
      }
    }

    if(strings.HasPrefix(f.Name, "quiqr/forms/")){
      formEndPoints++
      //spew.Dump(f.Name)
    }

    if(strings.HasPrefix(f.Name, "config.")){
      contents, _ := f.Contents()
      var extension = strings.TrimLeft(filepath.Ext(f.Name), ".")
      viper.SetConfigType(extension)
      viper.ReadConfig(bytes.NewBuffer([]byte(contents)))
      hugotheme, _ = viper.Get("theme").(string)
    }

    return nil
  })

  responseDict := &responseDictType{
    HugoVersion: hugover,
    HugoTheme: hugotheme,
    QuiqrFormsEndPoints: formEndPoints,
    QuiqrModel: quiqrModel,
    Screenshot: screenshot1,
  }
  responseJson, _ := json.Marshal(responseDict)
  fmt.Println(string(responseJson))
}



