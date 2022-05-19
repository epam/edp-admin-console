<a name="unreleased"></a>
## [Unreleased]

### Routine

- Update current development version [EPMDEDP-8832](https://jiraeu.epam.com/browse/EPMDEDP-8832)


<a name="v2.12.0"></a>
## [v2.12.0] - 2022-05-17
### Features

- Add Container as a library [EPMDEDP-8341](https://jiraeu.epam.com/browse/EPMDEDP-8341)
- Implement v2 get "/overview" request [EPMDEDP-8382](https://jiraeu.epam.com/browse/EPMDEDP-8382)
- prepare chi-based v2 html rendering [EPMDEDP-8384](https://jiraeu.epam.com/browse/EPMDEDP-8384)
- add cd pipeline overview page [EPMDEDP-8399](https://jiraeu.epam.com/browse/EPMDEDP-8399)
- get /v2/admin/edp/cd-pipeline/create request implemented. [EPMDEDP-8400](https://jiraeu.epam.com/browse/EPMDEDP-8400)
- implemented GET /admin/edp/cd-pipeline/{pipelineName}/update request [EPMDEDP-8401](https://jiraeu.epam.com/browse/EPMDEDP-8401)
- post /v2/admin/edp/cd-pipeline//{pipelineName}/update request implemented. [EPMDEDP-8402](https://jiraeu.epam.com/browse/EPMDEDP-8402)
- post /v2/admin/edp/cd-pipeline request implemented. [EPMDEDP-8403](https://jiraeu.epam.com/browse/EPMDEDP-8403)
- implemented v2/admin/edp/application/overview route [EPMDEDP-8404](https://jiraeu.epam.com/browse/EPMDEDP-8404)
- implemented GET method /v2/admin/edp/application/create route [EPMDEDP-8405](https://jiraeu.epam.com/browse/EPMDEDP-8405)
- add post application handler route [EPMDEDP-8406](https://jiraeu.epam.com/browse/EPMDEDP-8406)
- post /v2/admin/edp/cd-pipeline/delete request implemented. [EPMDEDP-8407](https://jiraeu.epam.com/browse/EPMDEDP-8407)
- get codebase overview request [EPMDEDP-8408](https://jiraeu.epam.com/browse/EPMDEDP-8408)
- get codebase update page [EPMDEDP-8409](https://jiraeu.epam.com/browse/EPMDEDP-8409)
- get codebases request - added fields to the response. [EPMDEDP-8412](https://jiraeu.epam.com/browse/EPMDEDP-8412)
- get codebase request - added fields to the response. [EPMDEDP-8412](https://jiraeu.epam.com/browse/EPMDEDP-8412)
- add auth for v2/admin [EPMDEDP-8416](https://jiraeu.epam.com/browse/EPMDEDP-8416)
- get codebase request - added commitMessagePattern-field  to the response. [EPMDEDP-8457](https://jiraeu.epam.com/browse/EPMDEDP-8457)
- add new fields to responses. [EPMDEDP-8479](https://jiraeu.epam.com/browse/EPMDEDP-8479)
- role access middleware [EPMDEDP-8483](https://jiraeu.epam.com/browse/EPMDEDP-8483)
- add post branch handler route [EPMDEDP-8484](https://jiraeu.epam.com/browse/EPMDEDP-8484)
- post /v2/admin/edp/codebase/:name/update request implemented. [EPMDEDP-8495](https://jiraeu.epam.com/browse/EPMDEDP-8495)
- POST /v2/admin/edp/codebase request implemented. [EPMDEDP-8508](https://jiraeu.epam.com/browse/EPMDEDP-8508)
- Implement v2 get "/cd-pipeline/{pipelinename}/overview" request [EPMDEDP-8802](https://jiraeu.epam.com/browse/EPMDEDP-8802)
- Implemented POST /v2/admin/edp/codebase/branch/delete request in chi-router. [EPMDEDP-8814](https://jiraeu.epam.com/browse/EPMDEDP-8814)
- Add additional logging [EPMDEDP-8821](https://jiraeu.epam.com/browse/EPMDEDP-8821)
- Ability to set suffix for branch version [EPMDEDP-8945](https://jiraeu.epam.com/browse/EPMDEDP-8945)

### Bug Fixes

- fix page rendering when using docker image [EPMDEDP-8384](https://jiraeu.epam.com/browse/EPMDEDP-8384)
- Fix changelog generation in GH Release Action [EPMDEDP-8386](https://jiraeu.epam.com/browse/EPMDEDP-8386)
- fix SetupAuthController if auth disabled [EPMDEDP-8482](https://jiraeu.epam.com/browse/EPMDEDP-8482)
- change logic of applications payload [EPMDEDP-8538](https://jiraeu.epam.com/browse/EPMDEDP-8538)
- Revert action-target for the delete_confirmation_template.html template. [EPMDEDP-8785](https://jiraeu.epam.com/browse/EPMDEDP-8785)
- Fix existed v2 pages [EPMDEDP-8795](https://jiraeu.epam.com/browse/EPMDEDP-8795)
- CSRF verification fixed for the operations: Create codebase branch, delete application. [EPMDEDP-8807](https://jiraeu.epam.com/browse/EPMDEDP-8807)
- CSRF verification fixed for the operation: delete codebase branch. [EPMDEDP-8810](https://jiraeu.epam.com/browse/EPMDEDP-8810)
- Jenkins branch url fix [EPMDEDP-8811](https://jiraeu.epam.com/browse/EPMDEDP-8811)
- Hide xsrfdata [EPMDEDP-8813](https://jiraeu.epam.com/browse/EPMDEDP-8813)
- Add annotation for non-first stage [EPMDEDP-8815](https://jiraeu.epam.com/browse/EPMDEDP-8815)
- Failed to fet the CodebaseImageStream CR when branch name contains '/' [EPMDEDP-8883](https://jiraeu.epam.com/browse/EPMDEDP-8883)
- Make sure we store applications for CDPipelines [EPMDEDP-8929](https://jiraeu.epam.com/browse/EPMDEDP-8929)

### Routine

- Update release flow for GH [EPMDEDP-8383](https://jiraeu.epam.com/browse/EPMDEDP-8383)
- Update current development version [EPMDEDP-8383](https://jiraeu.epam.com/browse/EPMDEDP-8383)
- extended logging for the get stage info request. [EPMDEDP-8497](https://jiraeu.epam.com/browse/EPMDEDP-8497)
- Update base docker image to alpine 3.15.4 [EPMDEDP-8853](https://jiraeu.epam.com/browse/EPMDEDP-8853)
- Update changelog [EPMDEDP-9185](https://jiraeu.epam.com/browse/EPMDEDP-9185)

### BREAKING CHANGE:


use gorilla's csrf implementation instead of beego's xsrf.


<a name="v2.11.5"></a>
## [v2.11.5] - 2022-05-17

<a name="v2.11.4"></a>
## [v2.11.4] - 2022-03-15
### Routine

- extended logging for the get stage info request. [EPMDEDP-8497](https://jiraeu.epam.com/browse/EPMDEDP-8497)


<a name="v2.11.3"></a>
## [v2.11.3] - 2022-02-28
### Features

- add new fields to responses. [EPMDEDP-8479](https://jiraeu.epam.com/browse/EPMDEDP-8479)


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


[Unreleased]: https://github.com/epam/edp-admin-console/compare/v2.12.0...HEAD
[v2.12.0]: https://github.com/epam/edp-admin-console/compare/v2.11.5...v2.12.0
[v2.11.5]: https://github.com/epam/edp-admin-console/compare/v2.11.4...v2.11.5
[v2.11.4]: https://github.com/epam/edp-admin-console/compare/v2.11.3...v2.11.4
[v2.11.3]: https://github.com/epam/edp-admin-console/compare/v2.11.2...v2.11.3
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
