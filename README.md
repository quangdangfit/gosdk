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
- Logger:  
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
  
 - Mgo wrapper:
    - Example  
    ```go
    package main
    import (
       ...
       transport/lib/utils/dbs"
    )
    
    func main(){
        dbConfig := dbs.DBConfig{
            MongoDBHosts: "localhost:27017",
            Database: "testdb",
        }
        
        db := dbs.Connect(dbConfig)
       
        #Define model           
        type Brand struct {
            Code string `json:"code" bson:"code"`
            Name string `json:"name" bson:"name"`
        }
   
        var results = []Brand{}
        collectionName := "brand"
        filter := bson.M{"code": "code"}
        
        err = db.FindMany(collectionName, filter, "_id", &results)
        if err != nil {
           ...
        }    
...
    }
    ```
