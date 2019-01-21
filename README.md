# github-repo-list

## Build

    cd src
    go get
    go build

## Usage

    github-repo-list --token GITHUB_TOKEN

### Required args

| flag | description |
| --- | --- |
| --token | A Github Personal Access Token |

### Optional args

| flag | description | default |
| --- | --- | --- |
| --org| Name of Github organization to query on | illuminateed |
| --delay| Time in ms to wait between Github API requests | 250 |
| --output | Path to file where CSV data will be written | results.csv |
| --per-page| Results per page to retrieve from Github API | 50 |
