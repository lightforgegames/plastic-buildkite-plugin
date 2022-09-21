# PlasticSCM Support for buildkite

A [Buildkite Plugin](https://buildkite.com/docs/agent/v3/plugins) that allows you to use [PlasticSCM](https://www.plasticscm.com/) instead of git in buildkite

## Example

Add the following to your `pipeline.yml`:

```yml
steps:
  - command: ls
    plugins:
      - lightforgegames/plastic: ~
```

