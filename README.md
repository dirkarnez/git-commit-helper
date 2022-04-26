git-commit-helper
=================
### TODOs
- auto detect repo origin provider if it is a repo
- Colored flags
- adhoc reject
- profiles
  - base
    - reject any unignored `node_modules` folder
      - **FORCE DELETE**
    - reject any archive file format other than `.zip`
      - by config
    - **FORCE DELETE MODE** macOS trash: `__MACOSX`, `.DS_Store`
    - safe full path length
    - auto write .gitignore?
  - github (with base)
    - less than 100MB per file
    - less than 2GB per commit / push
      - auto commit (split big changes)
    - auto create
      - repo
        - "main"
      - workflows
      - mirroring?
