package main

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/quangdangfit/gocommon/utils/config"
)

func main() {
	config.LoadConfig("config")

	/* config file contain:
	database:
	  host: localhost
	  name: database
	  user: quang

	*/

	fmt.Println(viper.GetString("database.host"))
}
