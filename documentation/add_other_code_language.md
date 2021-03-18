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