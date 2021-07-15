# Use Lint Stages for Code Review

This section contains the description of stages **dockerbuild-verify**, **dockerfile-lint** and **helm-lint**
which were added in the Code-Review pipeline.

These stages help obtain quick response on the validity of the code in the Code-Review pipeline in Kubernetes for all types
of applications supported by EDP out of the box.

  ![add_custom_lib2](../customization_resources/stages1.png)
  
Inspect the functions performed by the following stages:
 
1. **dockerbuild-verify** stage collects artifacts and builds an image from the Dockerfile without push to registry.
This stage is intended to check if the image is built.

>_**NOTE**: **dockerbuild-verify** stage is not a default one. To add this stage, navigate to the Jenkins job, 
select **Configure** tab and add the following parameters {"name": "build"},{"name": "dockerbuild-verify"}._
  
2. [**dockerfile-lint**](https://github.com/hadolint/hadolint) stage launches the _hadolint Dockerfile_ command
in order to check the Dockerfile.
 
3. [**helm-lint**](https://github.com/helm/chart-testing#chart-testing) stage launches
the _ct lint --charts-deploy-templates/_ command in order to validate the chart.




### Related Articles

* [Use Terraform Library in EDP](../cicd_customization/terraform_stages.md)
* [EDP Pipeline Framework](../cicd_customization/edp_pipeline_framework.md)
* [Promote Docker Images from ECR to Docker Hub](../cicd_customization/ecr_to_docker_stage.md)

