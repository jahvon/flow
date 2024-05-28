# Contributing

## How We Develop

This projects uses GitHub to host code, to track issues and feature requests, and to open and review pull requests.
We follow [GitHub Flow](https://docs.github.com/en/get-started/using-github/github-flow) to manage our development process.
In short, this means that all code changes happen through pull requests with descriptive commit messages, 
and that all changes are reviewed before they are merged.

By participating in this project, you agree to abide our [code of conduct](CODE_OF_CONDUCT.md).

### Submitting a Bug Report or Feature Request

We use GitHub issues to track public bugs and feature requests. Please fill out as much as you can in the issue template, 
and provide as much additional context as possible. 

When creating an issue, add any [labels](https://github.com/jahvon/flow/labels) that may be relevant.
The `triage` label will be added automatically and will be removed once the issue has been reviewed by a maintainer.

### Submitting Pull Requests

We actively welcome pull requests. Here are the steps to contribute:

1. Fork the repo and create your branch from `main`.
2. If you've added code that should be tested, add tests.
3. If you've changed Configs/APIs, update the documentation.
4. Ensure the pre-commit commands pass. This includes tests, linting, and code generation. See the [Development](DEVELOPMENT.md) guide for more information.
5. Open that pull request!

See the [Development guide](../DEVELOPMENT.md) for more information building and testing this project locally.

#### Best Practices
 - Try to keep your PR is up-to-date with the latest changes from the main branch.
 - Include a clear description of the problem and solution; especially noting side effects and testing details for the change.
 - Include a clear PR title (not a generic one like "Fixes issue").
 - Include a reference to the issue you are fixing (if applicable).
 - Include screenshots and animated GIFs in your pull request whenever possible.
 - Use descriptive commit messages. We loosely follow the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) standard.

### Review Process

PRs and issues will be reviewed by the maintainers as soon as possible. Please be patient when waiting for a response.
We welcome community members to review PRs as well, but only maintainers can merge them.

## License

Any contributions you make will be under the Apache 2.0 Software License. 
In short, when you submit code changes, your submissions are understood to be under the same [Apache License](https://choosealicense.com/licenses/apache-2.0/) that covers the project. 
