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
  - string
        (optional) webhook secret 
  -bot-server string
        The base URL the bot lives on
  -dry-run
        don't perform any actions, just print out the actions that would be taken if live
  -policies string
        The relative path to the policies file (default "./examples/.policies.yaml")
  -port int
        The port the bot listens on (default 7838)
  -token string
        The personal access token for the stated user
  -user string
        The Gitlab user this bot will act as
  -version
        display version of quetzal
  -webhook-endpoint string
        The webhook endpoint (default "/webhook-endpoint")

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

You can see an example `policies.yaml` in the `examples` directory.

### Policies

Policies are what drives `Quetzal`. There are 5 main properties to a policy

A single `Policy` is declared as part of an array of the `Policies` property.

#### Policy Properties
- [Name](#policy-name)
- [Resource](#policy-resource)
- [Conditions](#policy-conditions)
- [Limit](#policy-limit)
- [Actions](#policy-actions)

#### Policy Name
Is simply the name for this chosen policy 
```yaml
policies:
  - name: Awesome Policy
    #...
  - name: Round robin assignee
    #...
```

#### Policy Resource
The resource is the type of webhook this policy is for.
The available options are listed below and are the values of the `X-Gitlab-Event` header:

- Build Hook
- Deployment Hook
- Issue Hook
- Confidential Issue Hook
- Job Hook
- Merge Request Hook
- Note Hook
- Confidential Note Hook
- Pipeline Hook
- Push Hook
- Release Hook
- System Hook
- Tag Push Hook
- Wiki Page Hook

```yaml
policies:
  - name: Assign MR
    resource: Merge Request Hook
```

#### Policy Conditions
Policy Conditions allow a user to specify a series of conditions that confirm that this webhook event be processed.
All available options are type-safe and validated, once the policies file has been successfully parsed.


- [Date](#date-condition)
- [State](#state-condition)
- [Milestone](#milestone-condition)
- [Labels](#labels-condition)
- [Note](#note-condition)

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

#### State Condition
State must have the available options for hooks that have this property.
```yaml
policies:
  - name: Assign MR
    resource: Merge Request Hook
    conditions:
      state: open
```
`Merge`
- open
- close
- reopen
- update
- approved
- unapproved
- merge

`Issue`

- open
- close
- reopen
- update

`Release`

- create
- update

#### Milestone Condition
Milestone is the integer representation on the milestone
```yaml
policies:
  - name: Assign MR
    resource: Merge Request Hook
    conditions:
      milestone: 5
```

#### Labels Condition
The Labels condition accepts an array of labels by name to filter webhooks on.
The webhook must match all the provided labels on the policy to be valid.
```yaml
policies:
  - name: Assign MR
    resource: Merge Request Hook
    conditions:
      labels:
        - done
        - kittens       
```

#### Note Condition
The available options for `note` are:

| Property      | required | options                     |
| --------      | -------- | -------                     |
| noteType      | no       | `Commit`, `Issue`. Leaving blank will cause action on any `Note` webhook. **nb** `MergeRequest` and `Snippet` are not supported  |
| mentions      | no       | an array of mentioned users required to trigger action |
| command       | no       | any command you wish to use  |

The note condition allows your bot to respond to certain notes or even commands. As an example image the time a user mentions your bot with a specified command phrase.
`@yourbot show -help`
```yaml
    conditions:
      note:
        noteType: Issue
        mentions:
          - botUser
        command: show -help
```

#### Policy Limit

#### Policy Actions

Policy actions are what your bot performs when a webhook matching the policy pre-conditions is met.


### Contributions
All contributions are welcome, please open an issue/feature req at [GitLab](https://gitlab.com/jonny7/quetzal)

The Quetzal icon is attributed to [FreePick](https://www.freepik.com)
