# Use Lint Stages for Code Review

This section contains the description of new stages **dockerbuild-verify**, **dockerfile-lint** and **helm-lint** which were added in the Code-Review pipeline.

These stages help obtain quick call on the validity of the code in the Code-Review pipeline in Kubernetes for all types of applications supported by EDP out of the box. 

  ![add_custom_lib2](../customization_resources/stages1.png)
  
Inspect the functions performed by the following stages:
 
1. **dockerbuild-verify** stage collects artifacts and builds an image from the Dockerfile without push to registry. This stage is intended to check if the image is built.
  
2. **dockerfile-lint** stage launches the [_hadolint_](https://github.com/hadolint/hadolint) command in order to check the Dockerfile.
 
3. [**helm-lint**](https://github.com/helm/chart-testing#chart-testing) stage launches the _ct lint --charts-deploy-templates/_ command in order to validate the chart.




### Related Articles

* [Use Terraform Library in EDP](../documentation/cicd_customization/terraform_stages.md)
* [EDP Pipeline Framework](../documentation/cicd_customization/edp_pipeline_framework.md)
