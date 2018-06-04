# Run-occasionally

Runs commands on a schedule.

## Single command

To run a single command use the commandline. The command is specified as the argument
and the interval is specified by flags.

There are two flags to specify the interval.

* `-interval [i]` - `i` is a string representing the interval. Examples: `2s`, `5m`, `4h`
* `-cron [c]` - `c` is a cron expression with added seconds. The field should consist of
  six parts. Documentation can be found [https://godoc.org/github.com/robfig/cron#hdr-CRON_Expression_Format](here).
  Examples: `*/5 * * * * *` - every five seconds. `54 37 13 * * mon` - every monday at 13:37:54.

## Multiple commands

To run multiple commands at different schedules, use a configuration file. yaml, json and toml is supported.

Name the file `run-occasionally.[yaml|json|toml]` and put it into the working directory of the application.

Example configuration file:

```yaml
jobs:
  - command: date
    interval: "5s"
  - command: echo hello world
    interval: "1s"
```
