package world_test

import (
	"net"
	"testing"
	"time"

	"github.com/iFeddy/XelionCore/pkg/database"
	"github.com/iFeddy/XelionCore/pkg/network"
	"github.com/iFeddy/XelionCore/pkg/protocol"
	"github.com/iFeddy/XelionCore/pkg/server/world"
)

func TestGuildManager_FullFlow(t *testing.T) {
	memRepo := database.NewMemoryCharacterRepositoryWithFile("")

	// Crear personajes para las pruebas
	master := &database.Character{
		ID:        1001,
		AccountID: 1,
		Name:      "Aragorn",
		Job:       1, // Fighter
		Level:     60,
		TownName:  "RouN",
	}
	member := &database.Character{
		ID:        1002,
		AccountID: 2,
		Name:      "Legolas",
		Job:       3, // Archer
		Level:     55,
		TownName:  "RouN",
	}
	student := &database.Character{
		ID:        1003,
		AccountID: 3,
		Name:      "Pippin",
		Job:       2, // Cleric
		Level:     15,
		TownName:  "RouN",
	}

	_ = memRepo.CreateCharacter(master)
	_ = memRepo.CreateCharacter(member)
	_ = memRepo.CreateCharacter(student)

	gm := world.NewGuildManager(memRepo)

	// Conexiones simuladas mediante net.Pipe
	c1Server, c1Client := net.Pipe()
	defer c1Server.Close()
	defer c1Client.Close()

	c2Server, c2Client := net.Pipe()
	defer c2Server.Close()
	defer c2Client.Close()

	c3Server, c3Client := net.Pipe()
	defer c3Server.Close()
	defer c3Client.Close()

	drain := func(conn net.Conn) {
		buf := make([]byte, 2048)
		for {
			if _, err := conn.Read(buf); err != nil {
				return
			}
		}
	}
	go drain(c1Client)
	go drain(c2Client)
	go drain(c3Client)

	client1, err := network.NewClient(c1Server, nil)
	if err != nil {
		t.Fatalf("failed to create client1: %v", err)
	}
	client1.CharacterID = master.ID
	client1.Username = master.Name

	client2, err := network.NewClient(c2Server, nil)
	if err != nil {
		t.Fatalf("failed to create client2: %v", err)
	}
	client2.CharacterID = member.ID
	client2.Username = member.Name

	client3, err := network.NewClient(c3Server, nil)
	if err != nil {
		t.Fatalf("failed to create client3: %v", err)
	}
	client3.CharacterID = student.ID
	client3.Username = student.Name

	// Registrar jugadores online
	gm.RegisterOnline(client1, master.ID, master.Name)
	gm.RegisterOnline(client2, member.ID, member.Name)
	gm.RegisterOnline(client3, student.ID, student.Name)

	// 1. Crear Gremio
	gm.HandleGuildMake(client1, "Fellowship", "pass123", 0)

	// Verificar en repositorio
	guild, myMember, err := memRepo.GetGuildByCharID(master.ID)
	if err != nil || guild == nil || myMember == nil {
		t.Fatalf("expected master to belong to newly created guild, got err: %v", err)
	}
	if guild.Name != "Fellowship" {
		t.Errorf("expected guild name 'Fellowship', got '%s'", guild.Name)
	}
	if myMember.Grade != 0 {
		t.Errorf("expected master grade 0, got %d", myMember.Grade)
	}

	// 2. Listado de Gremios
	gm.HandleGuildList(client1, 0)
	guilds, err := memRepo.GetAllGuilds()
	if err != nil || len(guilds) != 1 {
		t.Fatalf("expected 1 guild in directory, got %d", len(guilds))
	}
	if guilds[0].GuildName != "Fellowship" {
		t.Errorf("expected directory guild 'Fellowship', got '%s'", guilds[0].GuildName)
	}

	// 3. Invitar a un Miembro
	gm.HandleInvite(client1, "Legolas")

	// 4. Aceptar Invitación
	gm.HandleJoinAck(client2, "Fellowship", 1)

	// Verificar que ahora Legolas es miembro
	_, legolasMember, err := memRepo.GetGuildByCharID(member.ID)
	if err != nil || legolasMember == nil {
		t.Fatalf("expected Legolas to be in guild, got err: %v", err)
	}
	if legolasMember.Grade != 4 {
		t.Errorf("expected new member to have grade 4 (Newbie), got %d", legolasMember.Grade)
	}

	// 5. Cambiar Rango a SubMaestro (Grade 1)
	gm.HandleGrade(client1, "Legolas", 1)
	_, updatedLegolas, _ := memRepo.GetGuildByCharID(member.ID)
	if updatedLegolas.Grade != 1 {
		t.Errorf("expected Legolas promoted to grade 1, got %d", updatedLegolas.Grade)
	}

	// 6. Actualizar Anuncio del Gremio
	gm.HandleNotify(client1, "Defend Middle-Earth!")
	guildAfterNotify, _ := memRepo.GetGuild(guild.ID)
	if guildAfterNotify.Notify != "Defend Middle-Earth!" {
		t.Errorf("expected notice 'Defend Middle-Earth!', got '%s'", guildAfterNotify.Notify)
	}

	// 7. Chat de Gremio
	gm.HandleChat(client1, "For Frodo!")

	// 8. Expulsar a un Miembro
	gm.HandleVanish(client1, "Legolas")
	_, kickedLegolas, _ := memRepo.GetGuildByCharID(member.ID)
	if kickedLegolas != nil {
		t.Errorf("expected Legolas to be kicked from guild, but still found")
	}

	// 9. Academia de Gremio (Creación automática, listado e ingreso)
	_ = memRepo.CreateGuildAcademy(guild.ID, master.ID)
	gm.HandleAcademyList(client3, 0)

	// Pippin ingresa a la academia
	gm.HandleAcademyJoin(client3, "Fellowship")
	students, err := memRepo.GetAcademyMembers(guild.ID)
	if err != nil || len(students) != 1 {
		t.Fatalf("expected 1 student in academy, got %d", len(students))
	}
	if students[0].CharName != "Pippin" {
		t.Errorf("expected student 'Pippin', got '%s'", students[0].CharName)
	}

	// Chat de Academia
	gm.HandleAcademyChat(client3, "What about second breakfast?")

	// Pippin abandona la academia
	gm.HandleAcademyLeave(client3)
	studentsAfterLeave, _ := memRepo.GetAcademyMembers(guild.ID)
	if len(studentsAfterLeave) != 0 {
		t.Errorf("expected 0 students after leave, got %d", len(studentsAfterLeave))
	}

	// 10. Disolver Gremio
	gm.HandleDismiss(client1)
	guildAfterDismiss, _ := memRepo.GetGuild(guild.ID)
	if guildAfterDismiss.DismissStatus != 1 {
		t.Errorf("expected guild dismiss status 1, got %d", guildAfterDismiss.DismissStatus)
	}

	// 11. Desconectar master
	gm.UnregisterOnline(master.ID)
}

func TestGuildProtocol_PacketStructures(t *testing.T) {
	// Verificar tamaños y offsets de los paquetes según PDB oficial

	// PROTO_NC_GUILD_MAKE_ACK: 27 bytes payload
	ack1 := protocol.New_PROTO_NC_GUILD_MAKE_ACK(0, 42, "TestGuild", 1, time.Now())
	if ack1.OpCode() != protocol.PROTO_NC_GUILD_MAKE_ACK {
		t.Errorf("expected opcode 0x%04X, got 0x%04X", protocol.PROTO_NC_GUILD_MAKE_ACK, ack1.OpCode())
	}
	// Total bytes con cabecera (1 byte size + 2 bytes netcmd + 27 bytes payload = 30 bytes)
	bytes1 := ack1.BuildWirePacket()
	if len(bytes1) != 30 {
		t.Errorf("expected 30 total bytes for PROTO_NC_GUILD_MAKE_ACK, got %d", len(bytes1))
	}

	// PROTO_NC_GUILD_MEMBER_INVITE_ACK: 22 bytes payload + 3 bytes header = 25 bytes
	ack2 := protocol.New_PROTO_NC_GUILD_MEMBER_INVITE_ACK("TestTarget", 0)
	bytes2 := ack2.BuildWirePacket()
	if len(bytes2) != 25 {
		t.Errorf("expected 25 total bytes for PROTO_NC_GUILD_MEMBER_INVITE_ACK, got %d", len(bytes2))
	}

	// PROTO_NC_GUILD_MEMBER_JOIN_REQ: 36 bytes payload + 3 bytes header = 39 bytes
	reqJoin := protocol.New_PROTO_NC_GUILD_MEMBER_JOIN_REQ("MyGuild", "Inviter")
	bytesJoin := reqJoin.BuildWirePacket()
	if len(bytesJoin) != 39 {
		t.Errorf("expected 39 total bytes for PROTO_NC_GUILD_MEMBER_JOIN_REQ, got %d", len(bytesJoin))
	}

	// PROTO_NC_GUILD_MEMBER_GRADE_ACK: 23 bytes payload + 3 bytes header = 26 bytes
	ackGrade := protocol.New_PROTO_NC_GUILD_MEMBER_GRADE_ACK("CharTarget", 2, 0)
	bytesGrade := ackGrade.BuildWirePacket()
	if len(bytesGrade) != 26 {
		t.Errorf("expected 26 total bytes for PROTO_NC_GUILD_MEMBER_GRADE_ACK, got %d", len(bytesGrade))
	}
}
