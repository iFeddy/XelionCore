package world_test

import (
	"net"
	"testing"

	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/server/world"
)

func TestFriendManager_FullFlow(t *testing.T) {
	memRepo := database.NewMemoryCharacterRepositoryWithFile("")

	// Crear dos personajes en memoria
	char1 := &database.Character{
		ID:        201,
		AccountID: 1,
		Name:      "Gandalf",
		Job:       4, // Mage
		Level:     50,
		TownName:  "RouN",
	}
	char2 := &database.Character{
		ID:        202,
		AccountID: 2,
		Name:      "Frodo",
		Job:       5, // Trickster
		Level:     20,
		TownName:  "RouN",
	}

	_ = memRepo.CreateCharacter(char1)
	_ = memRepo.CreateCharacter(char2)

	pm := world.NewPartyManager()
	fm := world.NewFriendManager(memRepo, pm)

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
	client1.CharacterID = char1.ID
	client1.Username = "Gandalf"

	client2, err := network.NewClient(c2Server, nil)
	if err != nil {
		t.Fatalf("Failed to create client2: %v", err)
	}
	client2.AccountID = 2
	client2.CharacterID = char2.ID
	client2.Username = "Frodo"

	pm.RegisterOnline(client1, char1.ID, "Gandalf", 4, 50)
	pm.RegisterOnline(client2, char2.ID, "Frodo", 5, 20)

	// 1. Gandalf solicita amistad a Frodo
	fm.RequestFriend(client1, "Frodo")

	// 2. Frodo acepta la solicitud
	fm.AcceptFriend(client2)

	// 3. Verificar que ambos personajes son amigos en el repositorio
	fList1, err := memRepo.GetFriends(char1.ID)
	if err != nil || len(fList1) != 1 || fList1[0].CharName != "Frodo" {
		t.Fatalf("Expected Gandalf to have friend Frodo, got %+v (err: %v)", fList1, err)
	}
	fList2, err := memRepo.GetFriends(char2.ID)
	if err != nil || len(fList2) != 1 || fList2[0].CharName != "Gandalf" {
		t.Fatalf("Expected Frodo to have friend Gandalf, got %+v (err: %v)", fList2, err)
	}

	// 4. Despachar lista de amigos
	fm.SendFriendList(client1)

	// 5. Consultar puntos de amistad
	fm.SendFriendPoints(client1)

	// 6. Eliminar amistad
	fm.DeleteFriend(client1, "Frodo")

	// 7. Verificar que ya no son amigos
	fListAfter, _ := memRepo.GetFriends(char1.ID)
	if len(fListAfter) != 0 {
		t.Fatalf("Expected 0 friends after deletion, got %d", len(fListAfter))
	}
}
