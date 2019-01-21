package main

import (
  "context"
  "flag"
  "fmt"
  "os"
  "time"
  "encoding/csv"
  "github.com/google/go-github/github"
  "golang.org/x/oauth2"
)

func main() {
  var (
    token   = flag.String("token", "", "Github auth token")
    org     = flag.String("org", "illuminateeducation", "Github organization")
    delay   = flag.Int("delay", 250, "Delay in ms between Github API requests to prevent rate-limiting")
    output  = flag.String("output", "results.csv", "Path to CSV to write results")
    per_page = flag.Int("per-page", 50, "Number of results per page")
  )

  flag.Parse()

  rate_limit_delay := time.Duration(*delay) * time.Millisecond

  if (*token == "") {
    fmt.Println("Must provide an oauth token with --token.")
    os.Exit(1)
  }

  ctx := context.Background()
  token_source := oauth2.StaticTokenSource(
    &oauth2.Token{AccessToken: *token},
  )
  token_client := oauth2.NewClient(ctx, token_source)
  client := github.NewClient(token_client)

  opt := &github.RepositoryListByOrgOptions{
    ListOptions: github.ListOptions{PerPage: *per_page},
  }
  var all_repos []*github.Repository

  fmt.Printf("Fetching %d results per page\n", *per_page)

  for {
    repos, resp, err := client.Repositories.ListByOrg(ctx, *org, opt)

    if err != nil {
      fmt.Println(err)
      os.Exit(1)
    }

    all_repos = append(all_repos, repos...)

    if resp.NextPage == 0 {
      fmt.Println("Reached last page")
      break
    }

    fmt.Printf("Fetched page %d/%d\n", opt.Page, resp.LastPage)

    opt.Page = resp.NextPage

    time.Sleep(rate_limit_delay)
  }

  fmt.Printf("Found %d repositories\n", len(all_repos))

  outfile, err := os.Create(*output)

  if err != nil {
    fmt.Printf("Error opening %s for writing: %s", *output, err)
  }

  csv_output := csv.NewWriter(outfile)
  records    := [][]string{}

  for i := 0; i < len(all_repos); i++ {
    var (
      repo_name         string
      repo_url          string
      repo_description  string
      repo_updated      string
      repo_commit_sha   string
      repo_committer    string
    )

    if all_repos[i].Name != nil {
      repo_name = *all_repos[i].Name
    }

    if all_repos[i].HTMLURL != nil {
      repo_url = *all_repos[i].HTMLURL
    }

    if all_repos[i].Description != nil {
      repo_description = *all_repos[i].Description
    }

    if all_repos[i].PushedAt != nil {
      repo_updated = all_repos[i].PushedAt.String()
    }

    opt := &github.CommitsListOptions{
      ListOptions: github.ListOptions{PerPage: 1},
    }

    commits, _, err := client.Repositories.ListCommits(ctx, *org, *all_repos[i].Name, opt)

    if err == nil {
      if commits[0].SHA != nil {
        repo_commit_sha = *commits[0].SHA
      }

      if commits[0].Committer != nil {
        repo_committer = *commits[0].Committer.HTMLURL
      }
    } else {
      fmt.Printf("Error getting commits: %s\n", err)
    }

    record := []string{repo_name, repo_url, repo_description, repo_updated, repo_commit_sha, repo_committer}
    records = append(records, record)

    time.Sleep(rate_limit_delay)
  }

  fmt.Println("Dumping data to CSV")

  if err = csv_output.WriteAll(records); err != nil {
    fmt.Printf("Error writing record to CSV: %s\n", err)
  }

  fmt.Println("Done")
}
