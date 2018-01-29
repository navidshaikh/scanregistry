package main

import (
  "encoding/json"
  "fmt"
  "net/http"
  "io/ioutil"
  "time"
  "log"
  "flag"
  "os"
)

const (
  cccpregistry = "https://registry.centos.org"
  apiversion = "/v2"
  catalog = "/_catalog"
  tagslist = "/tags/list"
)

type repoList struct {
  Repositories []string `json:"repositories"`
}

func findReposInRegistry(regURL string) ([]string, error) {
  regClient := http.Client{
    Timeout: time.Second * 10,
  }

  req, err := http.NewRequest(http.MethodGet, regURL, nil)
  if err != nil {
    log.Fatal(err)
    return nil, err
  }

  req.Header.Set("User-Agent", "scanregistry")

  res, getErr := regClient.Do(req)

  if getErr != nil {
    log.Fatal(getErr)
    return nil, getErr
  }

  body, readErr := ioutil.ReadAll(res.Body)

  if readErr != nil {
    log.Fatal(readErr)
    return nil, readErr
  }

  repos := repoList{}

  jsonErr := json.Unmarshal(body, &repos)

  if jsonErr != nil {
    log.Fatal(jsonErr)
    return nil, jsonErr
  }
  return repos.Repositories, nil
}


func main(){
  var registryURL, repositoryURL string

  flag.StringVar(&registryURL, "registry-url", cccpregistry, "Registry URL to scan repositories from.")
  flag.StringVar(&repositoryURL, "repository-url", "", "Specific repository URL to scan.")

  flag.Parse()

  registryURL = fmt.Sprintf("%s%s%s", registryURL, apiversion, catalog)

  fmt.Printf("\nRegistry URL %v provided to scan\n", registryURL)
  repolist, err := findReposInRegistry(registryURL)

  if err != nil {
    log.Fatal("\nFailed to get repositories from registry %v", registryURL)
    os.Exit(1)
  }

  fmt.Printf("\nRepos: %s", repolist)
}
