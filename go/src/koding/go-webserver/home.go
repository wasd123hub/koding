package main

import "net/http"

// type LoggedInUser struct {
//   Account       *models.Account
//   Machines      []*modelhelper.MachineContainer
//   Workspaces    []*models.Workspace
//   Group         *models.Group
//   Username      string
//   SessionId     string
//   Impersonating bool
// }

type LoggedInUser map[string]interface{}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	userInfo, err := fetchUserInfo(w, r)
	if err != nil {
		writeLoggedOutHomeToResp(w)
		return

	}
	onItem := make(chan Item, 0)
	onDone := make(chan LoggedInUser, 1)
	onError := make(chan error, 1)

	outputter := &Outputter{OnItem: onItem, OnError: onError}

	go collectItems(onItem, onDone, 4)

	go sendAccount(userInfo.Account, outputter)
	go fetchMachines(userInfo.UserId, outputter)
	go fetchWorkspaces(userInfo.AccountId, outputter)
	go fetchSocial(userInfo.AccountId, outputter)

	select {
	case <-onError:
		writeLoggedOutHomeToResp(w)
	case resp := <-onDone:
		writeLoggedInHomeToResp(w, resp)
	}
}

func collectItems(onItem <-chan Item, onDone chan<- LoggedInUser, max int) {
	resp := LoggedInUser{}

	for i := 0; i < max; i++ {
		item := <-onItem
		resp[item.Name] = item.Data
	}
}
