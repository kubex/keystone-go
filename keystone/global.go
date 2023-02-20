package keystone

import "github.com/kubex/definitions-go/app"

var defaultSetGlobalAppID *app.GlobalAppID

func SetMutationGlobalAppID(id app.GlobalAppID) {
	defaultSetGlobalAppID = &id
}

func ClearMutationGlobalAppID() {
	defaultSetGlobalAppID = nil
}
