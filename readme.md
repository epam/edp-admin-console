# EDP Admin Console

## Overview
The Admin Console management tool provides user interface to give clients an opportunity to manage business entities:
* Create Codebases as Applications, Libraries and Autotests;
* Create/Update CD Pipelines;

_*NOTE*: To interact with Admin Console via REST API, explore the [Create Codebase Entity](documentation/rest-api.md) page._

## Add Other Code Language

There is an ability to extend the default code languages when creating a codebase with the clone strategy.  
![other-language](readme-resource/ac_other_language.png "other-language")

_**NOTE**: The create strategy does not allow to customize the default code language set._
 
In order to customize the Build Tool list, perform the following:
1. Navigate to OpenShift, and edit the edp-admin-console deployment config map by adding the necessary code language into the BUILD TOOLS field. 
![build-tools](readme-resource/other_build_tool.png "build-tools")

_**NOTE**: Use the comma sign to separate the code languages in order to make them available, e.g. maven, gradle._

##  How to check the availability of the job-provision in Admin Console

To check the availability of the job-provision in Admin Console expand the Advanced Settings block during the codebase creation and see all job provisions you were created: 

 ![provisioner-ac](readme-resource/as_job_provision.png "provisioner-ac")