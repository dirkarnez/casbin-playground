package main

import (
	"flag"
	"fmt"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter"
	"github.com/jinzhu/gorm"
	// gormadapter "github.com/casbin/gorm-adapter/v3"
	// _ "github.com/go-sql-driver/mysql"
)

var (
	mode string
)

func main() {
	flag.StringVar(&mode, "mode", "file", "Mode: file, gorm. Default file mode")

	fmt.Println(fmt.Sprintf("mode = %s", mode))
	var enforcer *casbin.Enforcer = nil

	switch mode {
	case "file":
		// set policy
		enforcer, _ = casbin.NewEnforcer("rbac_model.conf", "rbac_policy.csv")
	case "gorm":
		// Initialize a Gorm adapter and use it in a Casbin enforcer:
		// The adapter will use the MySQL database named "casbin".
		// If it doesn't exist, the adapter will create it automatically.
		// You can also use an already existing gorm instance with gormadapter.NewAdapterByDB(gormInstance)
		var gormInstance *gorm.DB = nil
		adapter := gormadapter.NewAdapterByDB(gormInstance)
		enforcer, _ = casbin.NewEnforcer("rbac_model.conf", adapter)
		// Or you can use an existing DB "abc" like this:
		// The adapter will use the table named "casbin_rule".
		// If it doesn't exist, the adapter will create it automatically.
		// a := gormadapter.NewAdapter("mysql", "mysql_username:mysql_password@tcp(127.0.0.1:3306)/abc", true)

		// Load the policy from DB.
		enforcer.LoadPolicy()

		// Modify the policy.
		// e.AddPolicy(...)
		// e.RemovePolicy(...)

		// Save the policy back to DB.
		//e.SavePolicy()
	}

	user := "alice"
	resource_list := []string{"data1", "data3"}
	operation := "read"

	for _, resource := range resource_list {

		// apply policy
		if res, _ := enforcer.Enforce(user, resource, operation); res {
			// permit alice to read data1
			fmt.Println(fmt.Sprintf("%s to %s resource %s? permitted!", user, operation, resource))
		} else {
			// deny the request, show an error
			fmt.Println(fmt.Sprintf("%s to %s resource %s? denied!", user, operation, resource))
		}
	}
}
