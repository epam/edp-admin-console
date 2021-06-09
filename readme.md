# EDP Admin Console
## Overview
Admin Console is a central management tool in the EDP ecosystem that provides the ability to deploy infrastructure, project resources and new technologies in a simple way. 
Using Admin Console enables to manage business entities:
* Create Codebases as Applications, Libraries and Autotests;
* Create/Update CD Pipelines;

_**NOTE**: To interact with Admin Console via REST API, explore the [Create Codebase Entity](documentation/rest-api.md) page._

![overview-page](readme-resource/ac_overview_page.png "overview-page") 

- <strong>Navigation bar </strong>– consists of seven sections: Overview, Continuous Delivery, Applications, Autotests, Libraries, and Delivery Dashboard Diagram. 
Click the necessary section to add an entity, open a home page or check the diagram.
- <strong>User name</strong> – displays the registered user name. 
- <strong>Main links</strong> – displays the corresponding links to the major adjusted toolset, to the management tool and to the OpenShift cluster.

Admin Console is a complete tool allowing to manage and control the added to the environment codebases (applications, autotests, libraries) as well as to create a CD pipeline and check the visualization diagram. 
Inspect the main features available in Admin Console by following the corresponding link:

1. [Add Applications](documentation/add_applications.md)
2. [Add Autotests](documentation/add_autotests.md)
3. [Add Libraries](documentation/add_libraries.md)
4. [Add CD Pipelines](documentation/add_CD_pipelines.md)
5. [Delivery Dashboard Diagram](documentation/d_d_diagram.md)

_**NOTE**: The Admin Console link is available on the OpenShift overview page for your CI/CD project._

### Related Articles

* [GitHub Integration](documentation/github-integration.md)
* [GitLab Integration](documentation/gitlab-integration.md)
* [Local Development](documentation/local_development.md)
---
* [Add Jenkins Slave](https://github.com/epam/edp-jenkins-operator/blob/master/documentation/add-jenkins-slave.md#add-jenkins-slave)
* [Add Job Provision](https://github.com/epam/edp-jenkins-operator/blob/master/documentation/add-job-provision.md#add-job-provision)
* [Add Other Code Language](documentation/add_other_code_language.md)
* [Adjust GitLab CI Tool](documentation/ci-tool.md)
* [Adjust Import Strategy](documentation/import-strategy.md)
* [Adjust Integration With Jira Server](documentation/jira-server.md)
* [Adjust VCS Integration With Jira Server](documentation/jira_vcs_integration.md)
----
* [Add a New Custom Global Pipeline Library](documentation/cicd_customization/add_new_custom_global_pipeline_lib.md)
* [Associate IAM Roles With Service Accounts](documentation/enable_irsa.md)
* [Clone Project via Git Bash Terminal](documentation/cicd_customization/clone_project_using_gitbash.md)
* [Customize CD Pipeline](documentation/cicd_customization/customize-deploy-pipeline.md)
* [Customize CI Pipeline](documentation/cicd_customization/customize_ci_pipeline.md)
* [EDP Pipeline Framework](documentation/cicd_customization/edp_pipeline_framework.md)
* [EDP Stages](documentation/edp-stages.md)
* [Promote Docker Images from ECR to Docker Hub](documentation/cicd_customization/ecr_to_docker_stage.md)
* [Run Functional Autotest](documentation/cicd_customization/run_functional_autotest.md)
* [Use Lint Stages for Code Review](documentation/cicd_customization/code_review_stages.md)
* [Use Open Policy Agent Library in EDP](documentation/cicd_customization/opa_stages.md)
* [Use Terraform Library in EDP](documentation/cicd_customization/terraform_stages.md)