package main

import (
	"fmt"

	sw "github.com/dt665m/steamworksgo"
)

func main() {
	const (
		STEAM_KEY               = ""
		STEAM_TICKET_B64        = ""
		CUB_TICKET       uint32 = 0
		APP_ID           uint32 = 0
		ACCOUNT_ID       uint32 = 0
	)

	steam, err := sw.NewSteamWorks(STEAM_KEY, APP_ID, 3600)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = steam.VerifyAppTicket(STEAM_TICKET_B64, CUB_TICKET, ACCOUNT_ID)
	fmt.Println(err)
}
