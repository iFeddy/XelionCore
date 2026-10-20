package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

// MySQLCharacterRepository implementa CharacterRepository interactuando directamente
// con las tablas oficiales tCharacter, tCharacterShape, tCharacterOptions e tItem
// replicando con 100% de fidelidad los 8 stored procedures oficiales del CharServer.
type MySQLCharacterRepository struct {
	db *sql.DB
}

// NewMySQLCharacterRepository crea una nueva instancia del repositorio MySQL para personajes.
func NewMySQLCharacterRepository(db *sql.DB) *MySQLCharacterRepository {
	return &MySQLCharacterRepository{db: db}
}

// GetCharactersByAccount implementa la lógica combinada de p_Char_GetListOfUserChar,
// p_Char_GetLoginData y p_Item_GetList para obtener todos los personajes activos de una cuenta.
func (r *MySQLCharacterRepository) GetCharactersByAccount(accountID int) ([]*Character, error) {
	query := `
		SELECT 
			c.nCharNo, c.nUserNo, c.sID, c.nSlotNo, c.nLevel, c.sLoginZone,
			c.nLoginZoneX, c.nLoginZoneY, c.nLoginZoneD,
			c.nHP, c.nSP, c.nHPS, c.nSPS, c.nMoney, c.nExp, c.nFame,
			s.nClass, s.nRace, s.nGender, s.nHairType, s.nHairColor, s.nFaceShape
		FROM tCharacter c
		INNER JOIN tCharacterShape s ON c.nCharNo = s.nCharNo
		WHERE c.nUserNo = ? AND c.bDeleted = 0
		ORDER BY c.nSlotNo ASC;
	`

	rows, err := r.db.Query(query, accountID)
	if err != nil {
		return nil, fmt.Errorf("error querying characters for user %d: %w", accountID, err)
	}
	defer rows.Close()

	var characters []*Character
	for rows.Next() {
		ch := &Character{
			Equips: make(map[byte]uint16),
		}

		var (
			charNo     int
			userNo     int
			sID        string
			slotNo     byte
			level      uint16
			loginZone  string
			zoneX      int
			zoneY      int
			zoneD      byte
			hp         int
			sp         int
			hps        int16
			sps        int16
			money      int64
			exp        int64
			fame       int
			class      byte
			race       byte
			gender     byte
			hairType   byte
			hairColor  byte
			faceShape  byte
		)

		err := rows.Scan(
			&charNo, &userNo, &sID, &slotNo, &level, &loginZone,
			&zoneX, &zoneY, &zoneD,
			&hp, &sp, &hps, &sps, &money, &exp, &fame,
			&class, &race, &gender, &hairType, &hairColor, &faceShape,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning character row: %w", err)
		}

		ch.ID = charNo
		ch.AccountID = userNo
		ch.Name = sID
		ch.Slot = slotNo
		ch.Job = class
		ch.Level = byte(level)
		ch.Male = (gender != 0)
		ch.Hair = hairType
		ch.HairColor = hairColor
		ch.Face = faceShape
		ch.TownName = loginZone
		ch.XPos = zoneX
		ch.YPos = zoneY
		ch.CurHP = hp
		ch.CurSP = sp

		// Asignar en tablas canónicas
		ch.TCharacter.CharNo = charNo
		ch.TCharacter.UserNo = userNo
		ch.TCharacter.CharName = sID
		ch.TCharacter.Slot = slotNo
		ch.TCharacter.Class = class
		ch.TCharacter.Level = byte(level)
		ch.TCharacter.IsDelete = 0

		ch.TCharacterShape.CharNo = charNo
		ch.TCharacterShape.Gender = gender
		ch.TCharacterShape.Hair = hairType
		ch.TCharacterShape.HairColor = hairColor
		ch.TCharacterShape.Face = faceShape

		ch.TCharacterState.CharNo = charNo
		ch.TCharacterState.XPos = zoneX
		ch.TCharacterState.YPos = zoneY
		ch.TCharacterState.HP = hp
		ch.TCharacterState.SP = sp
		ch.TCharacterState.HPStones = hps
		ch.TCharacterState.SPStones = sps
		ch.TCharacterState.Exp = exp
		ch.TCharacterState.Money = money
		ch.TCharacterState.Fame = fame

		// Cargar equipamiento desde tItem (nStorageType = 0 es Equip)
		equips, err := r.getCharacterEquips(charNo)
		if err != nil {
			return nil, err
		}
		ch.Equips = equips

		characters = append(characters, ch)
	}

	return characters, nil
}

// GetCharacterByName busca un personaje activo por su nombre (sID).
func (r *MySQLCharacterRepository) GetCharacterByName(name string) (*Character, error) {
	query := `
		SELECT 
			c.nCharNo, c.nUserNo, c.sID, c.nSlotNo, c.nLevel, c.sLoginZone,
			c.nLoginZoneX, c.nLoginZoneY, c.nLoginZoneD,
			c.nHP, c.nSP, c.nHPS, c.nSPS, c.nMoney, c.nExp, c.nFame,
			s.nClass, s.nRace, s.nGender, s.nHairType, s.nHairColor, s.nFaceShape
		FROM tCharacter c
		INNER JOIN tCharacterShape s ON c.nCharNo = s.nCharNo
		WHERE LOWER(c.sID) = LOWER(?) AND c.bDeleted = 0
		LIMIT 1;
	`

	ch := &Character{Equips: make(map[byte]uint16)}
	var (
		charNo     int
		userNo     int
		sID        string
		slotNo     byte
		level      uint16
		loginZone  string
		zoneX      int
		zoneY      int
		zoneD      byte
		hp         int
		sp         int
		hps        int16
		sps        int16
		money      int64
		exp        int64
		fame       int
		class      byte
		race       byte
		gender     byte
		hairType   byte
		hairColor  byte
		faceShape  byte
	)

	err := r.db.QueryRow(query, strings.TrimSpace(name)).Scan(
		&charNo, &userNo, &sID, &slotNo, &level, &loginZone,
		&zoneX, &zoneY, &zoneD,
		&hp, &sp, &hps, &sps, &money, &exp, &fame,
		&class, &race, &gender, &hairType, &hairColor, &faceShape,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCharacterNotFound
		}
		return nil, fmt.Errorf("error querying character by name %s: %w", name, err)
	}

	ch.ID = charNo
	ch.AccountID = userNo
	ch.Name = sID
	ch.Slot = slotNo
	ch.Job = class
	ch.Level = byte(level)
	ch.Male = (gender != 0)
	ch.Hair = hairType
	ch.HairColor = hairColor
	ch.Face = faceShape
	ch.TownName = loginZone
	ch.XPos = zoneX
	ch.YPos = zoneY
	ch.CurHP = hp
	ch.CurSP = sp

	ch.TCharacter.CharNo = charNo
	ch.TCharacter.UserNo = userNo
	ch.TCharacter.CharName = sID
	ch.TCharacter.Slot = slotNo
	ch.TCharacter.Class = class
	ch.TCharacter.Level = byte(level)

	ch.TCharacterShape.CharNo = charNo
	ch.TCharacterShape.Gender = gender
	ch.TCharacterShape.Hair = hairType
	ch.TCharacterShape.HairColor = hairColor
	ch.TCharacterShape.Face = faceShape

	ch.TCharacterState.CharNo = charNo
	ch.TCharacterState.XPos = zoneX
	ch.TCharacterState.YPos = zoneY
	ch.TCharacterState.HP = hp
	ch.TCharacterState.SP = sp
	ch.TCharacterState.HPStones = hps
	ch.TCharacterState.SPStones = sps
	ch.TCharacterState.Exp = exp
	ch.TCharacterState.Money = money
	ch.TCharacterState.Fame = fame

	equips, err := r.getCharacterEquips(charNo)
	if err != nil {
		return nil, err
	}
	ch.Equips = equips

	return ch, nil
}

// GetCharacterByID busca un personaje activo por su ID único (nCharNo).
func (r *MySQLCharacterRepository) GetCharacterByID(id int) (*Character, error) {
	query := `
		SELECT 
			c.nCharNo, c.nUserNo, c.sID, c.nSlotNo, c.nLevel, c.sLoginZone,
			c.nLoginZoneX, c.nLoginZoneY, c.nLoginZoneD,
			c.nHP, c.nSP, c.nHPS, c.nSPS, c.nMoney, c.nExp, c.nFame,
			s.nClass, s.nRace, s.nGender, s.nHairType, s.nHairColor, s.nFaceShape
		FROM tCharacter c
		INNER JOIN tCharacterShape s ON c.nCharNo = s.nCharNo
		WHERE c.nCharNo = ? AND c.bDeleted = 0
		LIMIT 1;
	`

	ch := &Character{Equips: make(map[byte]uint16)}
	var (
		charNo     int
		userNo     int
		sID        string
		slotNo     byte
		level      uint16
		loginZone  string
		zoneX      int
		zoneY      int
		zoneD      byte
		hp         int
		sp         int
		hps        int16
		sps        int16
		money      int64
		exp        int64
		fame       int
		class      byte
		race       byte
		gender     byte
		hairType   byte
		hairColor  byte
		faceShape  byte
	)

	err := r.db.QueryRow(query, id).Scan(
		&charNo, &userNo, &sID, &slotNo, &level, &loginZone,
		&zoneX, &zoneY, &zoneD,
		&hp, &sp, &hps, &sps, &money, &exp, &fame,
		&class, &race, &gender, &hairType, &hairColor, &faceShape,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCharacterNotFound
		}
		return nil, fmt.Errorf("error querying character by ID %d: %w", id, err)
	}

	ch.ID = charNo
	ch.AccountID = userNo
	ch.Name = sID
	ch.Slot = slotNo
	ch.Job = class
	ch.Level = byte(level)
	ch.Male = (gender != 0)
	ch.Hair = hairType
	ch.HairColor = hairColor
	ch.Face = faceShape
	ch.TownName = loginZone
	ch.XPos = zoneX
	ch.YPos = zoneY
	ch.CurHP = hp
	ch.CurSP = sp

	ch.TCharacter.CharNo = charNo
	ch.TCharacter.UserNo = userNo
	ch.TCharacter.CharName = sID
	ch.TCharacter.Slot = slotNo
	ch.TCharacter.Class = class
	ch.TCharacter.Level = byte(level)

	ch.TCharacterShape.CharNo = charNo
	ch.TCharacterShape.Gender = gender
	ch.TCharacterShape.Hair = hairType
	ch.TCharacterShape.HairColor = hairColor
	ch.TCharacterShape.Face = faceShape

	ch.TCharacterState.CharNo = charNo
	ch.TCharacterState.XPos = zoneX
	ch.TCharacterState.YPos = zoneY
	ch.TCharacterState.HP = hp
	ch.TCharacterState.SP = sp
	ch.TCharacterState.HPStones = hps
	ch.TCharacterState.SPStones = sps
	ch.TCharacterState.Exp = exp
	ch.TCharacterState.Money = money
	ch.TCharacterState.Fame = fame

	equips, err := r.getCharacterEquips(charNo)
	if err != nil {
		return nil, err
	}
	ch.Equips = equips

	return ch, nil
}

// CreateCharacter implementa fielmente los procedimientos p_Char_Create y p_Char_CreateSetDefaultData
// realizando la inserción atómica en tCharacter, tCharacterShape, tCharacterOptions e tItem.
func (r *MySQLCharacterRepository) CreateCharacter(ch *Character) error {
	if ch.Slot > 5 {
		return ErrSlotOutOfRange
	}

	trimmedName := strings.TrimSpace(ch.Name)
	if trimmedName == "" {
		return errors.New("character name cannot be empty")
	}

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// 1. Validar si el slot ya está ocupado en la cuenta (p_Char_Create línea 40)
	var existingSlotChar int
	err = tx.QueryRow(
		"SELECT nCharNo FROM tCharacter WHERE nUserNo = ? AND nSlotNo = ? AND bDeleted = 0 LIMIT 1;",
		ch.AccountID, ch.Slot,
	).Scan(&existingSlotChar)
	if err == nil {
		return ErrSlotInUse
	} else if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("failed checking slot availability: %w", err)
	}

	// 2. Validar si el nombre ya existe en el servidor (p_Char_Create línea 47)
	var existingNameChar int
	err = tx.QueryRow(
		"SELECT nCharNo FROM tCharacter WHERE LOWER(sID) = LOWER(?) AND bDeleted = 0 LIMIT 1;",
		trimmedName,
	).Scan(&existingNameChar)
	if err == nil {
		return ErrCharacterExists
	} else if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("failed checking name uniqueness: %w", err)
	}

	// Valores iniciales por defecto (p_Char_CreateSetDefaultData)
	loginZone := ch.TownName
	if loginZone == "" {
		loginZone = "RouN"
	}
	level := ch.Level
	if level == 0 {
		level = 1
	}
	hp := ch.CurHP
	if hp == 0 {
		hp = 100
	}
	sp := ch.CurSP
	if sp == 0 {
		sp = 100
	}

	// 3. Insertar en tCharacter
	res, err := tx.Exec(`
		INSERT INTO tCharacter (
			nUserNo, nSlotNo, sID, nLevel, nHP, nSP, nHPS, nSPS, nMoney,
			sLoginZone, nLoginZoneX, nLoginZoneY, nLoginZoneD, bDeleted
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0);
	`, ch.AccountID, ch.Slot, trimmedName, level, hp, sp, 15, 7, ch.TCharacterState.Money,
		loginZone, ch.XPos, ch.YPos, 0x55,
	)
	if err != nil {
		return fmt.Errorf("failed inserting into tCharacter: %w", err)
	}

	charNo64, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed getting inserted charNo: %w", err)
	}
	charNo := int(charNo64)

	// 4. Insertar en tCharacterShape (p_Char_Create línea 69)
	genderByte := byte(0)
	if ch.Male {
		genderByte = 1
	}
	_, err = tx.Exec(`
		INSERT INTO tCharacterShape (
			nCharNo, nClass, nRace, nGender, nHairType, nHairColor, nFaceShape
		) VALUES (?, ?, ?, ?, ?, ?, ?);
	`, charNo, ch.Job, 0, genderByte, ch.Hair, ch.HairColor, ch.Face)
	if err != nil {
		return fmt.Errorf("failed inserting into tCharacterShape: %w", err)
	}

	// 5. Insertar fila en blanco en tCharacterOptions (p_Char_Create línea 78)
	_, err = tx.Exec(`INSERT INTO tCharacterOptions (nCharNo) VALUES (?);`, charNo)
	if err != nil {
		return fmt.Errorf("failed inserting into tCharacterOptions: %w", err)
	}

	// 6. Insertar equipamiento inicial si fue provisto
	if ch.Equips != nil {
		for slot, itemID := range ch.Equips {
			if itemID == 0 || itemID == 0xFFFF {
				continue
			}
			_, err = tx.Exec(`
				INSERT INTO tItem (nStorageType, nOwner, nStorage, nItemID, nFlags, dDate)
				VALUES (0, ?, ?, ?, 0, NOW());
			`, charNo, slot, itemID)
			if err != nil {
				return fmt.Errorf("failed inserting starter equip item: %w", err)
			}
		}
	}

	// 7. Insertar items de inventario iniciales si fueron provistos
	if ch.Items != nil {
		for _, it := range ch.Items {
			if it == nil || it.ItemID == 0 {
				continue
			}
			stType := it.StorageType
			if stType == 0 {
				stType = 9 // Default a inventario (Box 9)
			}
			_, err = tx.Exec(`
				INSERT INTO tItem (nStorageType, nOwner, nStorage, nItemID, nFlags, dDate)
				VALUES (?, ?, ?, ?, ?, NOW());
			`, stType, charNo, it.Storage, it.ItemID, it.Flags)
			if err != nil {
				return fmt.Errorf("failed inserting starter inventory item: %w", err)
			}
		}
	}

	// 8. Insertar shortcuts iniciales si fueron provistos
	if ch.Shortcuts != nil {
		for _, sc := range ch.Shortcuts {
			if sc == nil {
				continue
			}
			codeType := sc.CodeType
			if codeType == 0 {
				codeType = 1
			}
			_, err = tx.Exec(`
				INSERT INTO tCharacterShortCut (nCharNo, nSlotNo, nCodeType, nCodeNo, nValue)
				VALUES (?, ?, ?, ?, ?);
			`, charNo, sc.SlotNo, codeType, sc.CodeNo, sc.Value)
			if err != nil {
				return fmt.Errorf("failed inserting starter shortcut: %w", err)
			}
		}
	}

	// 9. Insertar habilidades iniciales si fueron provistas
	if ch.Skills != nil {
		for _, sk := range ch.Skills {
			if sk == nil || sk.SkillNo == 0 {
				continue
			}
			lvl := sk.SkillLevel
			if lvl == 0 {
				lvl = 1
			}
			_, err = tx.Exec(`
				INSERT INTO tSkill (nCharNo, nSkillNo, nSkillLevel)
				VALUES (?, ?, ?);
			`, charNo, sk.SkillNo, lvl)
			if err != nil {
				return fmt.Errorf("failed inserting starter skill: %w", err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit character creation: %w", err)
	}

	// Sincronizar identificador asignado
	ch.ID = charNo
	ch.TCharacter.CharNo = charNo
	ch.TCharacterShape.CharNo = charNo
	ch.TCharacterState.CharNo = charNo
	ch.TCharacterKeyMap.CharNo = charNo
	ch.TCharacterShortCut.CharNo = charNo
	ch.TCharacterGameOpt.CharNo = charNo
	ch.TCharacterSingleOpt.CharNo = charNo

	return nil
}

// DeleteCharacter implementa fielmente p_Char_Delete con soft-delete oficial:
// marca bDeleted = 1, dDeletedDate = NOW() y renombra el sID con formato 'nombre~YYYYMMDD#0'
// para liberar el nickname inmediatamente sin perder historial ni registros de auditoría.
func (r *MySQLCharacterRepository) DeleteCharacter(accountID int, slot byte) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var (
		charNo int
		sID    string
	)
	err = tx.QueryRow(
		"SELECT nCharNo, sID FROM tCharacter WHERE nUserNo = ? AND nSlotNo = ? AND bDeleted = 0 LIMIT 1;",
		accountID, slot,
	).Scan(&charNo, &sID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrCharacterNotFound
		}
		return fmt.Errorf("error finding character to delete: %w", err)
	}

	// Generar nombre de borrado oficial
	cleanName := sID
	if idx := strings.Index(cleanName, "~"); idx != -1 {
		cleanName = cleanName[:idx]
	}
	dateSuffix := time.Now().Format("20060102")
	changedID := fmt.Sprintf("%s~%s#0", cleanName, dateSuffix)

	_, err = tx.Exec(`
		UPDATE tCharacter
		SET bDeleted = 1, dDeletedDate = NOW(), sID = ?
		WHERE nCharNo = ?;
	`, changedID, charNo)
	if err != nil {
		return fmt.Errorf("failed executing soft delete on tCharacter: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit character deletion: %w", err)
	}

	return nil
}

// UpdateCharacter persiste el estado actual del personaje (HP, SP, posición, dinero, experiencia).
func (r *MySQLCharacterRepository) UpdateCharacter(ch *Character) error {
	_, err := r.db.Exec(`
		UPDATE tCharacter
		SET nHP = ?, nSP = ?, sLoginZone = ?, nLoginZoneX = ?, nLoginZoneY = ?, nMoney = ?, nExp = ?, nLevel = ?
		WHERE nCharNo = ? AND bDeleted = 0;
	`, ch.CurHP, ch.CurSP, ch.TownName, ch.XPos, ch.YPos, ch.TCharacterState.Money, ch.TCharacterState.Exp, ch.Level, ch.ID)
	if err != nil {
		return fmt.Errorf("error updating character %d: %w", ch.ID, err)
	}
	return nil
}

// getCharacterEquips obtiene el mapa de slots y IDs equipados (nStorageType = 0).
func (r *MySQLCharacterRepository) getCharacterEquips(charNo int) (map[byte]uint16, error) {
	rows, err := r.db.Query(`
		SELECT nStorage, nItemID FROM tItem WHERE nOwner = ? AND nStorageType = 0;
	`, charNo)
	if err != nil {
		return nil, fmt.Errorf("error fetching equips for char %d: %w", charNo, err)
	}
	defer rows.Close()

	equips := make(map[byte]uint16)
	for rows.Next() {
		var (
			slot   byte
			itemID uint16
		)
		if err := rows.Scan(&slot, &itemID); err != nil {
			return nil, fmt.Errorf("error scanning equip item: %w", err)
		}
		equips[slot] = itemID
	}
	return equips, nil
}

// GetCharacterItems recupera todos los items asociados a un personaje desde tItem.
func (r *MySQLCharacterRepository) GetCharacterItems(charNo int) ([]*Item, error) {
	rows, err := r.db.Query(`
		SELECT nItemKey, nStorageType, nOwner, nStorage, nItemID, nFlags, dDate
		FROM tItem
		WHERE nOwner = ?
		ORDER BY nStorageType, nStorage;
	`, charNo)
	if err != nil {
		return nil, fmt.Errorf("error querying items for char %d: %w", charNo, err)
	}
	defer rows.Close()

	var items []*Item
	for rows.Next() {
		var it Item
		if err := rows.Scan(&it.ItemKey, &it.StorageType, &it.Owner, &it.Storage, &it.ItemID, &it.Flags, &it.Date); err != nil {
			return nil, fmt.Errorf("error scanning item for char %d: %w", charNo, err)
		}
		items = append(items, &it)
	}
	return items, nil
}

// GetCharacterShortcuts recupera la barra de atajos oficial desde tCharacterShortCut.
func (r *MySQLCharacterRepository) GetCharacterShortcuts(charNo int) ([]*Shortcut, error) {
	rows, err := r.db.Query(`
		SELECT nCharNo, nSlotNo, nCodeType, nCodeNo, nValue
		FROM tCharacterShortCut
		WHERE nCharNo = ?
		ORDER BY nSlotNo;
	`, charNo)
	if err != nil {
		return nil, fmt.Errorf("error querying shortcuts for char %d: %w", charNo, err)
	}
	defer rows.Close()

	var shortcuts []*Shortcut
	for rows.Next() {
		var sc Shortcut
		if err := rows.Scan(&sc.CharNo, &sc.SlotNo, &sc.CodeType, &sc.CodeNo, &sc.Value); err != nil {
			return nil, fmt.Errorf("error scanning shortcut for char %d: %w", charNo, err)
		}
		shortcuts = append(shortcuts, &sc)
	}
	return shortcuts, nil
}

// GetCharacterSkills recupera las habilidades aprendidas desde tSkill.
func (r *MySQLCharacterRepository) GetCharacterSkills(charNo int) ([]*CharacterSkill, error) {
	rows, err := r.db.Query(`
		SELECT nCharNo, nSkillNo, nSkillLevel
		FROM tSkill
		WHERE nCharNo = ?
		ORDER BY nSkillNo;
	`, charNo)
	if err != nil {
		return nil, fmt.Errorf("error querying skills for char %d: %w", charNo, err)
	}
	defer rows.Close()

	var skills []*CharacterSkill
	for rows.Next() {
		var sk CharacterSkill
		if err := rows.Scan(&sk.CharNo, &sk.SkillNo, &sk.SkillLevel); err != nil {
			return nil, fmt.Errorf("error scanning skill for char %d: %w", charNo, err)
		}
		skills = append(skills, &sk)
	}
	return skills, nil
}

// SaveCharacterItem persiste un item individual en tItem.
func (r *MySQLCharacterRepository) SaveCharacterItem(item *Item) error {
	res, err := r.db.Exec(`
		INSERT INTO tItem (nStorageType, nOwner, nStorage, nItemID, nFlags, dDate)
		VALUES (?, ?, ?, ?, ?, NOW());
	`, item.StorageType, item.Owner, item.Storage, item.ItemID, item.Flags)
	if err != nil {
		return fmt.Errorf("error inserting item for char %d: %w", item.Owner, err)
	}
	id, err := res.LastInsertId()
	if err == nil {
		item.ItemKey = id
	}
	return nil
}

// SaveCharacterShortcuts persiste una lista de atajos en tCharacterShortCut.
func (r *MySQLCharacterRepository) SaveCharacterShortcuts(charNo int, shortcuts []*Shortcut) error {
	for _, sc := range shortcuts {
		_, err := r.db.Exec(`
			INSERT INTO tCharacterShortCut (nCharNo, nSlotNo, nCodeType, nCodeNo, nValue)
			VALUES (?, ?, ?, ?, ?)
			ON DUPLICATE KEY UPDATE nCodeType = VALUES(nCodeType), nCodeNo = VALUES(nCodeNo), nValue = VALUES(nValue), sys_update_date = NOW();
		`, charNo, sc.SlotNo, sc.CodeType, sc.CodeNo, sc.Value)
		if err != nil {
			return fmt.Errorf("error saving shortcut slot %d for char %d: %w", sc.SlotNo, charNo, err)
		}
	}
	return nil
}

// SaveCharacterSkills persiste habilidades aprendidas en tSkill.
func (r *MySQLCharacterRepository) SaveCharacterSkills(charNo int, skills []*CharacterSkill) error {
	for _, sk := range skills {
		_, err := r.db.Exec(`
			INSERT INTO tSkill (nCharNo, nSkillNo, nSkillLevel)
			VALUES (?, ?, ?)
			ON DUPLICATE KEY UPDATE nSkillLevel = VALUES(nSkillLevel);
		`, charNo, sk.SkillNo, sk.SkillLevel)
		if err != nil {
			return fmt.Errorf("error saving skill %d for char %d: %w", sk.SkillNo, charNo, err)
		}
	}
	return nil
}

// GetCharacterOptions recupera las opciones binarias persistidas de un personaje en tCharacterOptions
// (implementa p_Char_GetOptShortCutData, p_Char_GetOptShortCutSize, p_Char_GetOptGame, etc.).
func (r *MySQLCharacterRepository) GetCharacterOptions(charNo int) (*CharacterOptions, error) {
	query := `
		SELECT nCharNo, sShortCutData, sShortCutSize, sVideo, sSound, sGame, sWindowsPos, sKeyMapping
		FROM tCharacterOptions
		WHERE nCharNo = ?
		LIMIT 1;
	`
	opts := &CharacterOptions{CharNo: charNo}
	var (
		cNo        int
		scData     []byte
		scSize     []byte
		video      []byte
		sound      []byte
		game       []byte
		windowsPos []byte
		keyMapping []byte
	)

	err := r.db.QueryRow(query, charNo).Scan(
		&cNo, &scData, &scSize, &video, &sound, &game, &windowsPos, &keyMapping,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return opts, nil
		}
		return nil, fmt.Errorf("error querying character options for %d: %w", charNo, err)
	}

	opts.ShortCutData = scData
	opts.ShortCutSize = scSize
	opts.Video = video
	opts.Sound = sound
	opts.Game = game
	opts.WindowsPos = windowsPos
	opts.KeyMapping = keyMapping

	return opts, nil
}

// SaveCharacterOption guarda una opción binaria en tCharacterOptions según el procedimiento oficial:
// p_Char_SetOptShortCutData, p_Char_SetOptShortCutSize, p_Char_SetOptGame, etc.
func (r *MySQLCharacterRepository) SaveCharacterOption(charNo int, optType string, data []byte) error {
	var colName string
	switch optType {
	case "shortcut_data":
		colName = "sShortCutData"
	case "shortcut_size":
		colName = "sShortCutSize"
	case "game":
		colName = "sGame"
	case "sound":
		colName = "sSound"
	case "window_pos":
		colName = "sWindowsPos"
	case "key_mapping":
		colName = "sKeyMapping"
	case "video":
		colName = "sVideo"
	default:
		return fmt.Errorf("unknown character option type: %s", optType)
	}

	query := fmt.Sprintf(`
		INSERT INTO tCharacterOptions (nCharNo, %s)
		VALUES (?, ?)
		ON DUPLICATE KEY UPDATE %s = VALUES(%s);
	`, colName, colName, colName)

	_, err := r.db.Exec(query, charNo, data)
	if err != nil {
		return fmt.Errorf("error saving %s for char %d: %w", colName, charNo, err)
	}
	return nil
}
