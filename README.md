[![Go Report Card](https://goreportcard.com/badge/gitlab.com/jonny7/quetzal)](https://goreportcard.com/report/gitlab.com/jonny7/quetzal) [![Maintainability](https://api.codeclimate.com/v1/badges/d87c674cf1e418ef430d/maintainability)](https://codeclimate.com/github/jonny7/quetzal/maintainability) [![codecov](https://codecov.io/gh/jonny7/quetzal/branch/main/graph/badge.svg?token=NYF3T02QGL)](https://codecov.io/gh/jonny7/quetzal)

> This main repository for Quetzal is [GitLab](https://gitlab.com/jonny7/quetzal). Please create issues and questions there. [Github](https://github.com/jonny7/quetzal) is a mirror of that repo
# Quetzal

Quetzal is a GitLab bot written in Go. It takes inspiration from the `GitLab Triage Bot`.

## Installation
The easiest way is using Docker.
```shell
docker run --name my-quetzal -d -p 7838:7838 jonny7/quetzal 
```

However, you can build the source and run it yourself
```shell
go build -o quetzal ./cmd
# if your config and policy files are in the default locations then no other command is needed
./quetzal 
# otherwise, paths are relative to the quetzal binary
./quetzal -config="./path/to/file" -policies="../path/to/policies"
```

#### Versioning
Quetzal uses the SemVer specification. To query the binary, use the `-version` flag
```shell
./quetzal -version
# Quetzal version 1.1.1
```

### How Quetzal works
At its heart, Quetzal is a yaml driven policy based bot. It needs 2 things, a config file and a policy file. Both are `yaml` based and have default locations provided. Please note these are relative to the Quetzal binary.

| File type   | Default location       |
| ----------- | ---------------------- |
| config.yaml | ./config.yaml          |
| .policies.yaml | ./.policies.yaml    |

You can see examples of both of these file in the `examples` directory.

#### Config Options
| Property  | Usage |
| --------- | ----- |
| user      | The Gitlab user this bot will act as |
| token     | The personal access token for the stated `user` |
| repoHost  | The base URL for the GitLab instance |
| botServer | The base URL for where Quetzal is running |
| endpoint  | The endpoint where GitLab will be sending webhooks |
| secret    | (optional) the secret used by Quetzal to confirm the legitimacy of the webhook|
| port      | The port Quetzal will run on |
| policyPath| The relative path from ./quetzal to the policy file | 

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
