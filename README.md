![Build Status](https://github.com/j0shgrant/inventoryd/actions/workflows/main.yml/badge.svg)

# inventoryd
`inventoryd` is a lightweight binary that uses [Ably's Presence API](https://ably.com/documentation/core-features/presence) to expose an inventory of running Docker containers across a pool of hosts.

## Getting started
To install `inventoryd`, simply run the Makefile like below:
```shell
$ make
```
In order to allow `inventoryd` to authenticate with the Ably Presence API, you most export your API key as an environment variable:
```shell
$ export INVENTORYD_ABLY_KEY=<YOUR_ABLY_API_KEY>
```
Then, export the below environment variables to ensure your `inventoryd` will run against the correct Ably channels:
```shell
$ export INVENTORYD_ENVIRONMENT="test" && \
  export INVENTORYD_REGION="eu-west-1" && \
  export INVENTORYD_CRON_SCHEDULE="* * * * *"
```

Once your environment variables have been configured, simply run the `inventoryd` binary:
```shell
$ inventoryd
```

## Contributors
[Josh Grant](https://github.com/j0shgrant)