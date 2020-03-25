# AWS Cleaner 

## aws-cleanup

This tool will made cleaning on AWS object


### Build

Locally 
```sh
make build_local
```

Build container
```sh
make build
```

You can have a look on [Makefile](Makefile),in order to launch linter or swagger file generation

### Usage


```
Cleaning AWS Object like cloudwatch log stream

Usage:
  aws-cleanup [flags]
  aws-cleanup [command]

Available Commands:
  cloudwatch  cloudwatch cleaning
  help        Help about any command

Flags:
      --config string   config file (default is $HOME/.process-launcher.yaml)
      --debug           Set log level to Debug
  -h, --help            help for aws-cleanup
      --version         version

Use "aws-cleanup [command] --help" for more information about a command.
```





## Credit

Inspired by [aerostitch](https://github.com/aerostitch)

