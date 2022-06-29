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

type responseQuiqrsiteDictType struct {
  HugoVersion                string
  HugoTheme                  string
  QuiqrFormsEndPoints        int
  QuiqrModel                 string
  QuiqrEtalageName           string
  QuiqrEtalageDescription    string
  QuiqrEtalageHomepage       string
  QuiqrEtalageDemoUrl        string
  QuiqrEtalageLicense        string
  QuiqrEtalageLicenseURL     string
  QuiqrEtalageAuthor         string
  QuiqrEtalageAuthorHomepage string
  QuiqrEtalageScreenshots    []string
  Screenshot                 string
  ScreenshotImageType        string
}


func cmdRepoShowQuiqrsite() *cli.Command {

  return &cli.Command{
    Name: "repo_show_quiqrsite",
    Usage: "repo_show_quiqrsite",
    Flags: []cli.Flag{
      &cli.BoolFlag{
        Name:    "skip-base64-screenshot",
        Aliases: []string{"S"},
        Usage:   "do not output 1st found screenshot image",
      },
    },
    Action: func(c *cli.Context) error {
      skipBase64Screenshot := c.Bool("skip-base64-screenshot")
      url := c.Args().Get(0)
      showCaseQuiqrsite(url, skipBase64Screenshot)
      return nil
    },
  }
}


func showCaseQuiqrsite(url string, skipBase64Screenshot bool){

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

  var hugover string
  var hugotheme string
  formEndPoints := 0
  screenshotBase64 := ""
  var screenshotBase64ImageType string
  var quiqrModel string
  var quiqrEtalageName string
  var quiqrEtalageDescription string
  var quiqrEtalageHomepage string
  var quiqrEtalageDemoURL string
  var quiqrEtalageLicense string
  var quiqrEtalageLicenseURL string
  var quiqrEtalageAuthor string
  var quiqrEtalageAuthorHomepage string
  quiqrEtalageScreenshots := []string{}

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

    if(strings.HasPrefix(f.Name, "quiqr/etalage/template.json")){
      contents, _ := f.Contents()
      viper.SetConfigType("json")
      viper.ReadConfig(bytes.NewBuffer([]byte(contents)))

      quiqrEtalageName, _           = viper.Get("name").(string)
      quiqrEtalageDescription, _    = viper.Get("description").(string)
      quiqrEtalageHomepage, _       = viper.Get("homepage").(string)
      quiqrEtalageDemoURL, _        = viper.Get("demoURL").(string)
      quiqrEtalageLicense, _        = viper.Get("license").(string)
      quiqrEtalageLicenseURL, _     = viper.Get("licenseURL").(string)
      quiqrEtalageAuthor, _         = viper.Get("author").(string)
      quiqrEtalageAuthorHomepage, _ = viper.Get("authorHomepage").(string)
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
      quiqrEtalageScreenshots = append(quiqrEtalageScreenshots, f.Name)

      if(!skipBase64Screenshot && screenshotBase64 == ""){

        imgExts := []string{"jpg", "png", "gif", "jpeg"}
        extension := strings.ToLower(strings.TrimLeft(filepath.Ext(f.Name), "."))
        if(slices.Contains(imgExts, extension)){
          //spew.Dump(f.Name)
          contents, _ := f.Contents()
          contentsBytes := []byte(contents)

          var base64Encoding string

          mimeType := http.DetectContentType(contentsBytes)

          switch mimeType {
          case "image/jpeg":
            base64Encoding += "data:image/jpeg;base64,"
            screenshotBase64ImageType = "jpg"
          case "image/gif":
            base64Encoding += "data:image/gif;base64,"
            screenshotBase64ImageType = "gif"
          case "image/png":
            base64Encoding += "data:image/png;base64,"
            screenshotBase64ImageType = "png"
          }

          base64Encoding += toBase64(contentsBytes)
          screenshotBase64 = base64Encoding
        }
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

  responseDict := &responseQuiqrsiteDictType{

    HugoVersion:                hugover,
    HugoTheme:                  hugotheme,
    QuiqrFormsEndPoints:        formEndPoints,
    QuiqrModel:                 quiqrModel,
    QuiqrEtalageName:           quiqrEtalageName,
    QuiqrEtalageDescription:    quiqrEtalageDescription,
    QuiqrEtalageHomepage:       quiqrEtalageHomepage,
    QuiqrEtalageDemoUrl:        quiqrEtalageDemoURL,
    QuiqrEtalageLicense:        quiqrEtalageLicense,
    QuiqrEtalageLicenseURL:     quiqrEtalageLicenseURL,
    QuiqrEtalageAuthor:         quiqrEtalageAuthor,
    QuiqrEtalageAuthorHomepage: quiqrEtalageAuthorHomepage,
    QuiqrEtalageScreenshots:    quiqrEtalageScreenshots,
    Screenshot:                 screenshotBase64,
    ScreenshotImageType:        screenshotBase64ImageType,

  }
  responseJson, _ := json.Marshal(responseDict)
  fmt.Println(string(responseJson))
}



