# lib

**Common lib**
- Managed by version: [What is the version ?
](https://semver.org/)
- Please `note here` what changes in each version

## Version

### v1.0.0
- Initialize logger util
### v1.0.1
- Fix production mod
- Because some package have init function so maybe when you use `logger` with
 `production mod`, log in `init func` will log with `develop mode` (default is `develop mode`)
- Use:  
    Call in first line of main func:
    ```go
  package main
  import (
      ...
      "gitlab.ghn.vn/logistics/ts/lib/utils/logger"
  )
  
  func main(){
      logger.Initialize(config.Config.Production)
      ...
  }
    ```