# QiitaOrgWatcher

Bot tool for search posts in the specified Qiita Organization and notify Slack.


# Configurations

| NAME               | REQUIRED | TYPE               | NOTES                          |
|--------------------|----------|--------------------|--------------------------------|
| QIITA_TOKEN        | ○        | STRING             | Qiita access token             |
| QIITA_ORGANIZATION | ○        | STRING             | Qiita Organization ID          |
| CREATED            | --        | STRING(YYYY-MM-DD) | Created date for qiita article |
| SLACK_TOKEN        | ○        | STRING             | Slack access token             |
| SLACK_CHANNEL      | ○        | STRING             | Slack chanell                  |

# Quick Start

1. generate watcher access token
  * https://qiita.com/settings/applications
2. run application
```
$ <set envioromental variables>
$ export GO111MODULE=on
$ go mod init # only initialize
$ go mod tidy
$ go build
$ go run
```

# Deploy Google Cloud Function

Deploy to GCP Cloud Functions（HTTP）

```sh
$ gcloud functions --project ${project} deploy main \
--entry-point Receive \
--runtime go111 \
--set-env-vars QIITA_TOKEN=${qiita token},QIITA_ORGANIZATION=${qiita organization},SLACK_TOKEN=${slack token},SLACK_CHANNEL=${slack channel} \
--trigger-http
```

