# Go common

**Common lib in golang (datbase wrapper, logger, ...)**
- Managed by version: [What is the version ?
](https://semver.org/)
- Please `note here` what changes in each version

###Logger:  
  Call in first line of main func:
  ```go
  package main
  import (
      ...
      "gitlab.com/quangdangfit/gocommon/utils/logger"
  )
  
  func main(){
      logger.Initialize(config.Config.Production)
      ...
  }
  ```
  
###Mgo wrapper:
   ```go
    package main
    import (
       ...
       "gopkg.in/mgo.v2/bson"
       db "gitlab.com/quangdangfit/gocommon/database"
    )
    
    func main(){
        dbConfig := db.DBConfig{
        		Hosts:        "localhost:27017",
        		AuthDatabase: "admin",
        		AuthUserName: "",
        		AuthPassword: "",
        		Database:     "testdb",
        	}
        
        db := db.Connect(dbConfig)
       
        //Define model           
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
    }
   ```
