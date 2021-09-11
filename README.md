[![Go Report Card](https://goreportcard.com/badge/gitlab.com/jonny7/quetzal)](https://goreportcard.com/report/gitlab.com/jonny7/quetzal) [![Maintainability](https://api.codeclimate.com/v1/badges/d87c674cf1e418ef430d/maintainability)](https://codeclimate.com/github/jonny7/quetzal/maintainability) [![codecov](https://codecov.io/gh/jonny7/quetzal/branch/main/graph/badge.svg?token=NYF3T02QGL)](https://codecov.io/gh/jonny7/quetzal)

> This main repository for Quetzal is [GitLab](https://gitlab.com/jonny7/quetzal). Please create issues and questions there. [Github](https://github.com/jonny7/quetzal) is a mirror of that repo
# Quetzal

Quetzal is a GitLab bot written in Go. It takes inspiration from the `GitLab Triage Bot`.

## Installation
The easiest way is using Docker.
```shell
docker run --name my-quetzal -d -p 7838:7838 jonny7/quetzal 
```

However, you can build the source yourself
```shell
go build -o quetzal ./cmd/quetzal
```

#### Usage
```shell
./quetzal -h
Usage of ./quetzal:
  -bot-server string
        The base URL the bot lives on (default "https://bot-bot.com")
  -dry-run
        don't perform any actions, just print out the actions that would be taken if live
  -policies string
        The relative path to the policies file (default "./.policies.yaml")
  -port int
        The port the bot listens on (default 7838)
  -token string
        The personal access token for the stated user (default "notareatoken")
  -user string
        The Gitlab user this bot will act as (default "username@gitlab.com")
  -version
        display version of quetzal
  -webhook-endpoint string
        The webhook endpoint (default "/webhook/path")
  -webhook-secret string
        The (optional) webhook secret  (default "1234abcd")
```

#### Versioning
Quetzal uses the SemVer specification. To query the binary, use the `-version` flag
```shell
./quetzal -version
# Quetzal version 1.1.1
```

### How Quetzal works
At its heart, Quetzal is a yaml driven policy based bot. It needs some config parameters (listed above) and a policy file. The policy file is `yaml` based and has a default location provided. Please note this is relative to the Quetzal binary.

| File type   | Default location       |
| ----------- | ---------------------- |
| .policies.yaml | ./.policies.yaml    |

You can see examples of both of these file in the `examples` directory.

### Policies

Policies are what drives `Quetzal`. There are 4 main properties to a policy, with them all be technically optional. 
Though most likely, you'll always want an `action` and `name`.

A single `Policy` is declared as part of an array of the `Policies` property.

#### Policy Properties
- [Name](#policy-name)
- [Conditions](#policy-conditions)
- Limit
- Actions

#### Policy Name
```yaml
policies:
  - name: Awesome Policy
    #...
  - name: Round robin assignee
    #...
```

#### Policy Conditions
Policy Conditions allow a user to specify a series of conditions that confirm that this webhook event be processed.
All available options are type-safe and validated, once the policies file has been successfully parsed.


- [Date](#date-condition)

#### Date Condition

The available options for `date` are as follows:

| Property      | required | options                     |
| --------      | -------- | -------                     |
| attribute     | yes     | `created_at` or `updated_at` |
| condition     | yes     | `older_than` or `newer_than` |
| intervalType  | yes     | `days`, `weeks`, `months`, `years` |
| interval      | yes     | any valid unsigned int       |

An example date condition could look like this:
```yaml
    conditions:
      date:
        attribute: created_at
        condition: older_than
        intervalType: days
        interval: 10
```



### Contributions
All contributions are welcome, please open an issue/feature req at [GitLab](https://gitlab.com/jonny7/quetzal)

The Quetzal icon is attributed to [FreePick](https://www.freepik.com)
