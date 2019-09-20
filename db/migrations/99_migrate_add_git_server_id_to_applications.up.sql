update codebase
set git_server_id = gs.id
from git_server gs
where gs.name = 'gerrit'
  and codebase.type = 'application';