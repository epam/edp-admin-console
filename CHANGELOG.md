<a name="unreleased"></a>
## [Unreleased]


<a name="v2.11.2"></a>
## [v2.11.2] - 2022-02-22
### Features

- get codebase request - added commitMessagePattern-field  to the response. [EPMDEDP-8457](https://jiraeu.epam.com/browse/EPMDEDP-8457)

### Bug Fixes

- Fix changelog generation in GH Release Action [EPMDEDP-8386](https://jiraeu.epam.com/browse/EPMDEDP-8386)


<a name="v2.11.1"></a>
## [v2.11.1] - 2022-02-18
### Features

- get codebase request - added fields to the response. [EPMDEDP-8412](https://jiraeu.epam.com/browse/EPMDEDP-8412)
- get codebases request - added fields to the response. [EPMDEDP-8412](https://jiraeu.epam.com/browse/EPMDEDP-8412)

### Routine

- Update release flow for GH [EPMDEDP-8383](https://jiraeu.epam.com/browse/EPMDEDP-8383)


<a name="v2.11.0"></a>
## [v2.11.0] - 2022-02-14
### Features

- Update Makefile changelog target [EPMDEDP-8218](https://jiraeu.epam.com/browse/EPMDEDP-8218)
- k8s namespaced client draft [EPMDEDP-8229](https://jiraeu.epam.com/browse/EPMDEDP-8229)
- GetImageStreamFromStage implementation [EPMDEDP-8260](https://jiraeu.epam.com/browse/EPMDEDP-8260)
- implemented logic to get input image stream for the next (non-first) stage [EPMDEDP-8262](https://jiraeu.epam.com/browse/EPMDEDP-8262)
- reworked logic to get input image stream for the (non-first) stage [EPMDEDP-8262](https://jiraeu.epam.com/browse/EPMDEDP-8262)
- prepare chi-based v2 api route. [EPMDEDP-8264](https://jiraeu.epam.com/browse/EPMDEDP-8264)
- handler preparation [EPMDEDP-8265](https://jiraeu.epam.com/browse/EPMDEDP-8265)
- pipeline-stage handler [EPMDEDP-8265](https://jiraeu.epam.com/browse/EPMDEDP-8265)
- add SetupNamespacedClient func [EPMDEDP-8265](https://jiraeu.epam.com/browse/EPMDEDP-8265)
- get pipeline request using custom resource as a source. [EPMDEDP-8281](https://jiraeu.epam.com/browse/EPMDEDP-8281)
- get codebases request using custom resource as a source. [EPMDEDP-8282](https://jiraeu.epam.com/browse/EPMDEDP-8282)
- get sinlge-codebase request using custom resource as a source. [EPMDEDP-8324](https://jiraeu.epam.com/browse/EPMDEDP-8324)
- new fields to get codebase response [EPMDEDP-8325](https://jiraeu.epam.com/browse/EPMDEDP-8325)

### Code Refactoring

- Address golangci-lint issues [EPMDEDP-7945](https://jiraeu.epam.com/browse/EPMDEDP-7945)
- Remove unused helm chart [EPMDEDP-7997](https://jiraeu.epam.com/browse/EPMDEDP-7997)
- remove init() funcs. [EPMDEDP-8263](https://jiraeu.epam.com/browse/EPMDEDP-8263)

### Testing

- stub tests files. [EPMDEDP-8263](https://jiraeu.epam.com/browse/EPMDEDP-8263)
- use Makefile to make build and to make test. [EPMDEDP-8264](https://jiraeu.epam.com/browse/EPMDEDP-8264)
- add fake k8s config. [EPMDEDP-8264](https://jiraeu.epam.com/browse/EPMDEDP-8264)
- increase test coverage. [EPMDEDP-8264](https://jiraeu.epam.com/browse/EPMDEDP-8264)
- Add test for SetupNamespacedClient [EPMDEDP-8265](https://jiraeu.epam.com/browse/EPMDEDP-8265)

### Routine

- Update release CI pipelines [EPMDEDP-7847](https://jiraeu.epam.com/browse/EPMDEDP-7847)


<a name="v2.10.0"></a>
## [v2.10.0] - 2021-12-07
### Bug Fixes

- Empty 'Deployment Script' field in Admin Console by default [EPMDEDP-7280](https://jiraeu.epam.com/browse/EPMDEDP-7280)
- Use Default branch for branch creation [EPMDEDP-7552](https://jiraeu.epam.com/browse/EPMDEDP-7552)
- Fix CI Pipelines [EPMDEDP-7847](https://jiraeu.epam.com/browse/EPMDEDP-7847)

### Code Refactoring

- Set ormDebug to false by default [EPMDEDP-7847](https://jiraeu.epam.com/browse/EPMDEDP-7847)

### Formatting

- add explicit names [EPMDEDP-7847](https://jiraeu.epam.com/browse/EPMDEDP-7847)

### Routine

- Prepare for release [EPMDEDP-7847](https://jiraeu.epam.com/browse/EPMDEDP-7847)
- Update docker image [EPMDEDP-7895](https://jiraeu.epam.com/browse/EPMDEDP-7895)

### Documentation

- Update the link to documentation [EPMDEDP-7781](https://jiraeu.epam.com/browse/EPMDEDP-7781)


<a name="v2.9.0"></a>
## [v2.9.0] - 2021-12-07

<a name="v2.8.1"></a>
## [v2.8.1] - 2021-12-07

<a name="v2.8.0"></a>
## [v2.8.0] - 2021-12-07

<a name="v2.7.3"></a>
## [v2.7.3] - 2021-12-07

<a name="v2.7.2"></a>
## [v2.7.2] - 2021-12-07

<a name="v2.7.1"></a>
## [v2.7.1] - 2021-12-07

<a name="v2.7.0"></a>
## v2.7.0 - 2021-12-07
### Reverts

- [EPMDEDP-5591]: Add Import strategy in app.conf"
- [EPMDEDP-3929] Add new action type to display in Action Log table


[Unreleased]: https://github.com/epam/edp-admin-console/compare/v2.11.2...HEAD
[v2.11.2]: https://github.com/epam/edp-admin-console/compare/v2.11.1...v2.11.2
[v2.11.1]: https://github.com/epam/edp-admin-console/compare/v2.11.0...v2.11.1
[v2.11.0]: https://github.com/epam/edp-admin-console/compare/v2.10.0...v2.11.0
[v2.10.0]: https://github.com/epam/edp-admin-console/compare/v2.9.0...v2.10.0
[v2.9.0]: https://github.com/epam/edp-admin-console/compare/v2.8.1...v2.9.0
[v2.8.1]: https://github.com/epam/edp-admin-console/compare/v2.8.0...v2.8.1
[v2.8.0]: https://github.com/epam/edp-admin-console/compare/v2.7.3...v2.8.0
[v2.7.3]: https://github.com/epam/edp-admin-console/compare/v2.7.2...v2.7.3
[v2.7.2]: https://github.com/epam/edp-admin-console/compare/v2.7.1...v2.7.2
[v2.7.1]: https://github.com/epam/edp-admin-console/compare/v2.7.0...v2.7.1
