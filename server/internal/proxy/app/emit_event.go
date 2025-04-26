package app

import (
	"fmt"
	"proxy-srv/internal/proxy/domain"

	"github.com/ethereum/go-ethereum/common"
)

func (a *app) EmitNewParticipantJoined(roomId string, newParticipant string) error {
	resp, err := a.smc.GetParticipantsAndTracksOfRoom(roomId)
	if err != nil {
		return err
	}
	participants := resp.Arg0
	pName := ""
	for _, p := range participants {
		if p.WalletAddress.String() == newParticipant {
			pName = p.Name
		}
	}
	defer func() {
		if r := recover(); r != nil {
			a.errChan <- fmt.Errorf("emit new participant joined %v", r)
		}
	}()
	for _, p := range participants {
		go func(addr common.Address) {
			err := a.smc.EmitEventToFrontend(roomId, addr.String(),
				domain.EventNewParticipantJoined{
					EventName:       domain.EventNewParticipantJoinedName,
					Participant:     newParticipant,
					ParticipantName: pName,
				})
			if err != nil {
				a.errChan <- fmt.Errorf("emit new participant joined %v", err)
			}
		}(p.WalletAddress)
	}
	return nil
}
