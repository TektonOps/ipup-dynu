# Contributing to ipup-dynu

We want to make contributing to this project as easy and transparent as
possible.

## Project structure

- `main.go` - main function with entrypoint to app.
- `ineternal` - Contains internal packages
- `config` - Configuration params
- `dynu` - Updates the dns records on dynu
- `ip` - Getting public ip address 
- `utils` - Helper function


## Commits

Commit messages should be well formatted, and to make that "standardized", we
are using Conventional Commits.

```shell

  <type>[<scope>]: <short summary>
     │      │             │
     │      │             └─> Summary in present tense. Not capitalized. No
     |      |                 period at the end. 
     │      │
     │      └─> Scope (optional): eg. common, compiler, authentication, core
     │
     └─> Type: chore, docs, feat, fix, refactor, style, or test.
     
```

You can follow the documentation on
[their website](https://www.conventionalcommits.org).

## Pull Requests

We actively welcome your pull requests.

1. Fork the repo and create your branch from `main`.
2. If you've added code that should be tested, add tests.

## Issues

We use GitHub issues to track public bugs. Please ensure your description is
clear and has sufficient instructions to be able to reproduce the issue.

## License

By contributing to plogs, you agree that your contributions will be licensed
under the LICENSE file in the root directory of this source tree.
