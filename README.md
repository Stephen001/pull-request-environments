# pull-request-environments
A simple service that can be used to listen for GitHub pull request events, and clean up any associated ephemeral environments.

## Intended Usecase

The idea with this service is that you can create, update or remove pull request infrastructure when your developers create, update or close pull requests. For example, you could use this tool to clean up kubernetes namespaces automatically created by your CI system, which you could use as preview environments for changes.
