create type "action" as enum
  ('codebase_registration', 'gerrit_repository_provisioning', 'jenkins_configuration',
    'perf_registration', 'setup_deployment_templates', 'codebase_branch_registration');