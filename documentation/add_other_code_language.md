## Add Other Code Language

There is an ability to extend the default code languages when creating a codebase with the clone strategy.

![other-language](../readme-resource/ac_other_language.png "other-language")

>_**NOTE**: The create strategy does not allow to customize the default code language set._

In order to customize the Build Tool list, perform the following:
* Navigate to OpenShift, and edit the edp-admin-console deployment by adding the necessary code language into
the BUILD TOOLS field.

![build-tools](../readme-resource/other_build_tool.png "build-tools")

>_**NOTE**: Use the comma sign to separate the code languages in order to make them available, e.g. maven, gradle._

* Add the Jenkins slave by following the
[Add Jenkins Slave](https://github.com/epam/edp-jenkins-operator/blob/master/documentation/add-jenkins-slave.md#add-jenkins-slave) instruction.

* As a result, the newly added Jenkins slave will be available in the **Select Jenkins Slave** dropdown list of the
Advanced Settings block during the codebase creation:

![jenkins-slave](../readme-resource/ac_jenkins_slave.png "jenkins-slave")

* Extend or modify the Jenkins provisioner by following the
[Add Job Provisioner](https://github.com/epam/edp-jenkins-operator/blob/master/documentation/add-job-provision.md) instruction.

If it is necessary to create Code Review and Build pipelines, add corresponding entries (e.g. stages[Build-application-docker], [Code-review-application-docker]). See the example below:

```java
...
stages['Code-review-application-docker'] = '[{"name": "gerrit-checkout"}' + "${commitValidateStage}" + ',{"name": "sonar"}]'
stages['Build-application-docker'] = '[{"name": "checkout"},{"name": "get-version"},{"name": "sonar"},' +
                                     '{"name": "build-image-kaniko"}' + "${createJFVStage}" + ',{"name": "git-tag"}]'
...
```
![jenkins-provisioner](../readme-resource/ac_jenkins_provisioner.png "jenkins-provisioner")

>_**NOTE**: Application is one of the available options. Another option might be to add a library. Please refer to the [Add Libraries](https://github.com/epam/edp-admin-console/blob/release/2.7/documentation/add_libraries.md) page for details._

#### Related Articles

* [Add Jenkins Slave](https://github.com/epam/edp-jenkins-operator/blob/master/documentation/add-jenkins-slave.md#add-jenkins-slave)
* [Add Job Provisioner](https://github.com/epam/edp-jenkins-operator/blob/master/documentation/add-job-provision.md)
* [Add Libraries](https://github.com/epam/edp-admin-console/blob/release/2.7/documentation/add_libraries.md)
