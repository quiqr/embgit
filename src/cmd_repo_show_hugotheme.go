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

type responseHugothemeDictType struct {
  Name string
  License string
  LicenseLink string
  Description string
  MinHugoVersion string
  Author string
  Screenshot string
  ExampleSite bool
  Features []string
  Tags []string
  Homepage string
  Demosite string
  AuthorHomepage string
}


func cmdRepoShowHugotheme() *cli.Command {

  return &cli.Command{
    Name: "repo_show_hugotheme",
    Usage: "repo_show_hugotheme",
    Action: func(c *cli.Context) error {

      url := c.Args().Get(0)
      showCaseHugotheme(url)
      return nil
    },
  }
}


func showCaseHugotheme(url string){

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

  var minHugoVersion string
  var themeName string
  var themeLicense string
  var themeLicenseLink string
  var themeDescription string
  var themeAuthor string
  var screenshot1 string
  var themeExampleSite bool
  var themeFeatures []string
  var themeTags []string
  var themeHomepage string
  var themeDemosite string
  var themeAuthorHomepage string

  tree.Files().ForEach(func(f *object.File) error {

    if(strings.HasPrefix(f.Name, "exampleSite/")){
      themeExampleSite = true;
    }

    if(strings.HasPrefix(f.Name, "images/screenshot")){

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

    if(strings.HasPrefix(f.Name, "theme.")){
      contents, _ := f.Contents()
      var extension = strings.TrimLeft(filepath.Ext(f.Name), ".")
      viper.SetConfigType(extension)
      viper.ReadConfig(bytes.NewBuffer([]byte(contents)))
      minHugoVersion, _ = viper.Get("min_version").(string)
      themeName, _ = viper.Get("name").(string)
      themeLicense, _ = viper.Get("license").(string)
      themeLicenseLink, _ = viper.Get("licenselink").(string)
      themeDescription, _ = viper.Get("description").(string)
      themeHomepage, _ = viper.Get("homepage").(string)
      themeDemosite, _ = viper.Get("demosite").(string)
      themeAuthor, _ = viper.Get("author.name").(string)
      themeAuthorHomepage, _ = viper.Get("author.homepage").(string)
      themeFeatures = viper.GetStringSlice("features")
      themeTags = viper.GetStringSlice("tags")
    }

    return nil
  })

  responseDict := &responseHugothemeDictType{
    Screenshot: screenshot1,
    Name: themeName,
    License: themeLicense,
    LicenseLink: themeLicenseLink,
    Description: themeDescription,
    MinHugoVersion: minHugoVersion,
    Author: themeAuthor,
    ExampleSite: themeExampleSite,
    Features: themeFeatures,
    Tags: themeTags,
    Homepage: themeHomepage,
    Demosite: themeDemosite,
    AuthorHomepage: themeAuthorHomepage,
  }
  responseJson, _ := json.Marshal(responseDict)
  fmt.Println(string(responseJson))
}



