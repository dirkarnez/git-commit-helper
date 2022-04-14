git-commit-helper
=================
### TODOs
- auto detect repo origin provider if it is a repo
- profiles
  - base
    - reject any `node_modules` folder
    - reject any archive file format other than `.zip`
    - safe full path length
    - auto write .gitignore?
  - github (with base)
    - less than 100MB per file
    - less than 2GB per commit / push
    - auto create
      - repo
        - "main"
      - workflows
