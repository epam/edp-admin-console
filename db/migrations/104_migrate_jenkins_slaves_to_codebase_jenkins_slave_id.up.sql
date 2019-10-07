update codebase c
set jenkins_slave_id = js.id
from jenkins_slave js
where js.name = c.build_tool;