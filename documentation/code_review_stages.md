# Using Lint Stages for Code Review

This section contains the description of new stages **dockerbuild-verify**, **dockerfile-lint** and h**elm-lint** which were added in the Code-Review pipeline.

These stages help obtain quick feedback on the validity of the code in the Code-Review pipeline in Kubernetes for all types of applications supported by EDP out of the box, thus providing insight on the correctness of the code taken from the repository.

  ![add-app3](../readme-resource/stages1.png "add-app3_2")
  
  The stages perform the following functions:
 
 - **dockerbuild-verify** stage collects artifacts and builds an image from the Dockerfile without push to registry. This stage is intended to check if the image is built.
  
 - **dockerfile-lint** launches the [_hadolint_](https://github.com/hadolint/hadolint) command in order to check the Dockerfile.
 
 - [**helm-lint**](https://github.com/helm/chart-testing#chart-testing) launches the _ct lint --charts-deploy-templates/_ command.


### Related Articles

* [Using Terraform Library in EDP](../documentation/terraform_stages.md)
