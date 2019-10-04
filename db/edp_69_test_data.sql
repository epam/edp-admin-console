INSERT INTO third_party_service (id, name, description, version) VALUES (1, 'activemq', 'Openshift template for ActiveMQ service', '5.15.6');
INSERT INTO third_party_service (id, name, description, version) VALUES (2, 'activemq-ephemeral', 'Openshift template for ActiveMQ service', '5.15.6');
INSERT INTO third_party_service (id, name, description, version) VALUES (3, 'hystrix', 'Openshift template for Hystrix dashboard service', '1.2.0');
INSERT INTO third_party_service (id, name, description, version) VALUES (4, 'keycloak', 'Openshift template for KeyCloak service', '3.4.3.Final');
INSERT INTO third_party_service (id, name, description, version) VALUES (5, 'keycloak-ephemeral', 'Openshift template for KeyCloak service', '3.4.3.Final');
INSERT INTO third_party_service (id, name, description, version) VALUES (6, 'rabbitmq', 'Openshift template for RabbitMQ service', '3.7.15-management');
INSERT INTO third_party_service (id, name, description, version) VALUES (7, 'rabbitmq-ephemeral', 'Openshift template for RabbitMQ service', '3.7.15-management');
INSERT INTO third_party_service (id, name, description, version) VALUES (8, 'turbine', 'Openshift template for Turbine service', '1.2.0');
INSERT INTO third_party_service (id, name, description, version) VALUES (9, 'zipkin-ephemeral', 'Openshift template for Zipkin service', '2.6.0');

insert into git_server(id, name, available, hostname) values (1, 'gerrit', true, 'gerrit-mr-3013-2-edp-cicd.delivery.aws.main.edp.projects.epam.com');

INSERT INTO codebase (id, type, name, language, framework, build_tool, strategy, repository_url, route_site, route_path, database_kind, database_version, database_capacity, database_storage, status, test_report_framework, description, git_server_id) VALUES (2, 'application', 'bar-service', 'java', 'springboot', 'maven', 'clone', 'https://git.epam.com/epmd-edp/examples/basic/bar-service.git', '', '', '', '', '', '', 'active', '', '', 1);
INSERT INTO codebase (id, type, name, language, framework, build_tool, strategy, repository_url, route_site, route_path, database_kind, database_version, database_capacity, database_storage, status, test_report_framework, description, git_server_id) VALUES (1, 'application', 'foo-service', 'java', 'springboot', 'maven', 'clone', 'https://git.epam.com/epmd-edp/examples/basic/foo-service.git', '', '', '', '', '', '', 'active', '', '', 1);
INSERT INTO codebase (id, type, name, language, framework, build_tool, strategy, repository_url, route_site, route_path, database_kind, database_version, database_capacity, database_storage, status, test_report_framework, description, git_server_id) VALUES (4, 'autotests', 'foobar-tests', 'java', null, 'maven', 'clone', 'https://git.epam.com/epmd-edp/examples/basic/foobar-tests.git', '', '', '', '', '', '', 'active', 'allure', 'foobar-tests', 1);
INSERT INTO codebase (id, type, name, language, framework, build_tool, strategy, repository_url, route_site, route_path, database_kind, database_version, database_capacity, database_storage, status, test_report_framework, description, git_server_id) VALUES (3, 'application', 'zuul', 'java', 'springboot', 'maven', 'clone', 'https://git.epam.com/epmd-edp/examples/basic/zuul.git', 'zuul', '/', '', '', '', '', 'active', '', '', 1);

INSERT INTO codebase_docker_stream (id, codebase_id, oc_image_stream_name) VALUES (1, 1, 'foo-service-master');
INSERT INTO codebase_docker_stream (id, codebase_id, oc_image_stream_name) VALUES (2, 2, 'bar-service-master');
INSERT INTO codebase_docker_stream (id, codebase_id, oc_image_stream_name) VALUES (3, 3, 'zuul-master');
INSERT INTO codebase_docker_stream (id, codebase_id, oc_image_stream_name) VALUES (4, 2, 'dev-foobar-sit-bar-service-verified');
INSERT INTO codebase_docker_stream (id, codebase_id, oc_image_stream_name) VALUES (5, 1, 'dev-foobar-sit-foo-service-verified');
INSERT INTO codebase_docker_stream (id, codebase_id, oc_image_stream_name) VALUES (6, 3, 'dev-foobar-sit-zuul-verified');
INSERT INTO codebase_docker_stream (id, codebase_id, oc_image_stream_name) VALUES (7, 2, 'dev-foobar-qa-bar-service-verified');
INSERT INTO codebase_docker_stream (id, codebase_id, oc_image_stream_name) VALUES (8, 1, 'dev-foobar-qa-foo-service-verified');
INSERT INTO codebase_docker_stream (id, codebase_id, oc_image_stream_name) VALUES (9, 3, 'dev-foobar-qa-zuul-verified');
INSERT INTO codebase_docker_stream (id, codebase_id, oc_image_stream_name) VALUES (10, 3, 'zuul-fix1');

INSERT INTO codebase_branch (id, name, codebase_id, from_commit, output_codebase_docker_stream_id, status) VALUES (4, 'master', 4, '', null, 'active');
INSERT INTO codebase_branch (id, name, codebase_id, from_commit, output_codebase_docker_stream_id, status) VALUES (3, 'master', 3, '', 3, 'active');
INSERT INTO codebase_branch (id, name, codebase_id, from_commit, output_codebase_docker_stream_id, status) VALUES (2, 'master', 2, '', 2, 'active');
INSERT INTO codebase_branch (id, name, codebase_id, from_commit, output_codebase_docker_stream_id, status) VALUES (1, 'master', 1, '', 1, 'active');
INSERT INTO codebase_branch (id, name, codebase_id, from_commit, output_codebase_docker_stream_id, status) VALUES (5, 'fix1', 3, '', 10, 'active');

INSERT INTO cd_pipeline (id, name, status) VALUES (1, 'dev', 'active');

INSERT INTO cd_stage (id, name, cd_pipeline_id, description, trigger_type, quality_gate, jenkins_step_name, "order", status) VALUES (2, 'foobar-qa', 1, 'QA environment', 'manual', 'manual', 'qa', 1, 'active');
INSERT INTO cd_stage (id, name, cd_pipeline_id, description, trigger_type, quality_gate, jenkins_step_name, "order", status) VALUES (1, 'foobar-sit', 1, 'SIT environment', 'manual', 'autotests', 'sit', 0, 'active');

INSERT INTO cd_pipeline_codebase_branch (cd_pipeline_id, codebase_branch_id) VALUES (1, 2);
INSERT INTO cd_pipeline_codebase_branch (cd_pipeline_id, codebase_branch_id) VALUES (1, 1);
INSERT INTO cd_pipeline_codebase_branch (cd_pipeline_id, codebase_branch_id) VALUES (1, 3);

INSERT INTO cd_pipeline_third_party_service (cd_pipeline_id, third_party_service_id) VALUES (1, 6);

INSERT INTO cd_stage_codebase_branch (cd_stage_id, codebase_branch_id) VALUES (1, 4);