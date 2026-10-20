package world_test

import (
	"net"
	"testing"

	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/server/world"
)

func TestPartyManager_FullFlow(t *testing.T) {
	pm := world.NewPartyManager()

	c1Server, c1Client := net.Pipe()
	defer c1Server.Close()
	defer c1Client.Close()

	c2Server, c2Client := net.Pipe()
	defer c2Server.Close()
	defer c2Client.Close()

	go func() {
		buf := make([]byte, 1024)
		for {
			if _, err := c1Client.Read(buf); err != nil {
				return
			}
		}
	}()
	go func() {
		buf := make([]byte, 1024)
		for {
			if _, err := c2Client.Read(buf); err != nil {
				return
			}
		}
	}()

	client1, err := network.NewClient(c1Server, nil)
	if err != nil {
		t.Fatalf("Failed to create client1: %v", err)
	}
	client1.AccountID = 1
	client1.CharacterID = 101
	client1.Username = "Player1"

	client2, err := network.NewClient(c2Server, nil)
	if err != nil {
		t.Fatalf("Failed to create client2: %v", err)
	}
	client2.AccountID = 2
	client2.CharacterID = 102
	client2.Username = "Player2"

	pm.RegisterOnline(client1, 101, "Arthur", 1, 10)
	pm.RegisterOnline(client2, 102, "Lancelot", 2, 12)

	// 2. Arthur invita a Lancelot
	pm.InviteMember(client1, "Lancelot")

	// 3. Lancelot acepta la invitación
	pm.AcceptInvite(client2)

	// 4. Cambiar regla de botín
	pm.SetLooting(client1, 1, 3)

	// 5. Publicar anuncio en buscador de grupos
	pm.AddFinderAd(client1, "Farming goblins in RouN")
	pm.ListFinderAds(client1)
	pm.DeleteFinderAd(client1)

	// 6. Transferir liderazgo
	pm.ChangeMaster(client1, "Lancelot")

	// 7. Salir del grupo
	pm.LeaveParty(client2)
}
