package controller

import (
	// "archive/zip"
	"encoding/json"
	// "fmt"
	// "github.com/eaciit/colony-core/v0"
	"github.com/eaciit/acl"
	"github.com/eaciit/colony-manager/helper"
	"github.com/eaciit/dbox"
	// _ "github.com/eaciit/dbox/dbc/jsons"
	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/toolkit"
	// "io"
	// "io/ioutil"
	// "os"
	// "path/filepath"
	// "strings"
	// "time"
)

type UserController struct {
	App
}

func CreateUserController(s *knot.Server) *UserController {
	var controller = new(UserController)
	controller.Server = s
	return controller
}

func (a *UserController) GetUser(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson
	a.InitialSetDatabase()
	tUser := new(acl.User)

	arrm := make([]toolkit.M, 0, 0)
	c, err := acl.Find(tUser, nil, toolkit.M{}.Set("take", 0))
	if err == nil {
		err = c.Fetch(&arrm, 0, false)
	}

	if err != nil {
		return helper.CreateResult(true, nil, err.Error())
	} else {
		return helper.CreateResult(true, arrm, "")
	}

}
func (a *UserController) FindUser(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson
	a.InitialSetDatabase()
	payload := map[string]interface{}{}
	err := r.GetPayload(&payload)
	tUser := new(acl.User)
	err = acl.FindByID(tUser, payload["_id"].(string))
	if err != nil {
		return helper.CreateResult(true, nil, err.Error())
	} else {
		return helper.CreateResult(true, tUser, "")
	}

}
func (a *UserController) DeleteUser(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson
	a.InitialSetDatabase()
	tUser := new(acl.User)
	payload := map[string]interface{}{}
	err := r.GetPayload(&payload)
	err = acl.FindByID(tUser, payload["_id"].(string))

	if err != nil {
		return helper.CreateResult(true, nil, err.Error())
	}

	err = acl.Delete(tUser)
	if err != nil {
		return helper.CreateResult(true, nil, err.Error())
	}

	return helper.CreateResult(true, nil, "success")
}

func (a *UserController) SaveUser(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson
	a.InitialSetDatabase()
	// payload := struct {
	// 	LoginID  string
	// 	FullName string
	// 	Email    string
	// 	Password string
	// 	Enable   bool
	// 	Groups   []string
	// }{}
	// err := r.GetPayload(&payload)
	// if err != nil {
	// 	return helper.CreateResult(false, nil, err.Error())
	// }
	// fmt.Println(payload.Groups)

	payload := map[string]interface{}{}
	err := r.GetPayload(&payload)
	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}
	user := payload["user"].(map[string]interface{})
	groups := user["Groups"]
	var group []string
	for _, v := range groups.([]interface{}) {
		group = append(group, v.(string))
	}
	initUser := new(acl.User)
	id := toolkit.RandomString(32)
	initUser.ID = id
	initUser.LoginID = user["LoginID"].(string)
	initUser.FullName = user["FullName"].(string)
	initUser.Email = user["Email"].(string)
	initUser.Password = user["Password"].(string)
	initUser.Enable = user["Enable"].(bool)
	initUser.Groups = group

	err = acl.Save(initUser)

	if err != nil {
		return helper.CreateResult(true, nil, err.Error())
	}

	err = acl.ChangePassword(initUser.ID, user["Password"].(string))
	if err != nil {
		return helper.CreateResult(true, nil, err.Error())
	}

	var grant map[string]interface{}
	for _, p := range payload["grants"].([]interface{}) {
		dat := []byte(p.(string))
		if err = json.Unmarshal(dat, &grant); err != nil {
			return helper.CreateResult(true, nil, err.Error())
		}
		AccessID := grant["AccessID"].(string)
		Accessvalue := grant["AccessValue"]
		for _, v := range Accessvalue.([]interface{}) {
			switch v {
			case "AccessCreate":
				initUser.Grant(AccessID, acl.AccessCreate)
			case "AccessRead":
				initUser.Grant(AccessID, acl.AccessRead)
			case "AccessUpdate":
				initUser.Grant(AccessID, acl.AccessUpdate)
			case "AccessDelete":
				initUser.Grant(AccessID, acl.AccessDelete)
			case "AccessSpecial1":
				initUser.Grant(AccessID, acl.AccessSpecial1)
			case "AccessSpecial2":
				initUser.Grant(AccessID, acl.AccessSpecial2)
			case "AccessSpecial3":
				initUser.Grant(AccessID, acl.AccessSpecial3)
			case "AccessSpecial4":
				initUser.Grant(AccessID, acl.AccessSpecial4)
			}
		}
	}
	err = acl.Save(initUser)

	return helper.CreateResult(true, nil, "sukses")
}

func (a *UserController) prepareconnection() (conn dbox.IConnection, err error) {
	conn, err = dbox.NewConnection("mongo",
		&dbox.ConnectionInfo{"localhost:27017", "valegrab", "", "", toolkit.M{}.Set("timeout", 3)})
	if err != nil {
		return
	}
	err = conn.Connect()
	return
}

func (a *UserController) InitialSetDatabase() error {
	conn, err := a.prepareconnection()

	if err != nil {
		return err
	}

	err = acl.SetDb(conn)
	if err != nil {
		return err
	}
	return nil
}
