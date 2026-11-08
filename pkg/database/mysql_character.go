package database

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// MySQLCharacterRepository implementa CharacterRepository interactuando 100% mediante
// Stored Procedures oficiales contra la base de datos xelion_online (World00_Character).
// Regla Mandatoria: CERO SQL inline (SELECT/INSERT/UPDATE/DELETE) en todo el repositorio.
type MySQLCharacterRepository struct {
	db *sql.DB
}

// NewMySQLCharacterRepository crea una nueva instancia del repositorio MySQL para personajes.
func NewMySQLCharacterRepository(db *sql.DB) *MySQLCharacterRepository {
	return &MySQLCharacterRepository{db: db}
}

// GetCharactersByAccount implementa la recuperación de personajes activos de una cuenta
// invocando p_Char_GetListOfUserChar, p_Char_GetAllData, p_Char_GetShape y p_Item_GetListType.
func (r *MySQLCharacterRepository) GetCharactersByAccount(accountID int) ([]*Character, error) {
	rows, err := r.db.Query("CALL p_Char_GetListOfUserChar(?)", accountID)
	if err != nil {
		return nil, fmt.Errorf("error calling p_Char_GetListOfUserChar(%d): %w", accountID, err)
	}
	defer rows.Close()

	var charIDs []int
	for rows.Next() {
		var charNo int
		if err := rows.Scan(&charNo); err != nil {
			return nil, fmt.Errorf("error scanning charNo from p_Char_GetListOfUserChar: %w", err)
		}
		charIDs = append(charIDs, charNo)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	_ = rows.Close()

	changeList, _ := r.GetCharacterChangeIDList(accountID)
	changeByChar := make(map[int]*CharIDChangeInfo, len(changeList))
	for _, entry := range changeList {
		changeByChar[entry.CharNo] = entry
	}

	var characters []*Character
	for _, id := range charIDs {
		ch, err := r.GetCharacterByID(id)
		if err != nil {
			continue
		}
		if entry, ok := changeByChar[ch.ID]; ok {
			ch.ChangeID = entry
		}
		characters = append(characters, ch)
	}
	return characters, nil
}

// GetCharacterByName localiza un personaje por su nombre oficial invocando p_Char_GetCharNo.
func (r *MySQLCharacterRepository) GetCharacterByName(name string) (*Character, error) {
	var charNo int
	row := r.db.QueryRow("CALL p_Char_GetCharNo(?, @out_nCharNo)", strings.TrimSpace(name))
	if err := row.Scan(&charNo); err != nil {
		return nil, fmt.Errorf("error calling p_Char_GetCharNo: %w", err)
	}
	if charNo <= 0 {
		return nil, ErrCharacterNotFound
	}
	return r.GetCharacterByID(charNo)
}

// GetCharacterByID recupera los datos canónicos de un personaje invocando p_Char_GetAllData,
// p_Char_GetShape y p_Item_GetListType (para equipamiento activo).
func (r *MySQLCharacterRepository) GetCharacterByID(id int) (*Character, error) {
	ch := &Character{
		Equips:        make(map[byte]uint16),
		EquipUpgrades: make(map[byte]byte),
	}

	// 1. Cargar tCharacter canónico
	rows, err := r.db.Query("CALL p_Char_GetAllData(?)", id)
	if err != nil {
		return nil, fmt.Errorf("error calling p_Char_GetAllData(%d): %w", id, err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, ErrCharacterNotFound
	}

	var (
		userNo, slotNo, level                                                                       int
		exp, money                                                                                  int64
		hp, sp, ap, hps, sps, ps, gs, fame, flags, adminLevel                                       int
		loginZone                                                                                   string
		loginZoneX, loginZoneY, loginZoneD, prisonMin, pkCount, pkYellowTime                        int
		strength, constitute, dexterity, intelligence, mentalPower, redistributePoint, playMin, lc int
		lastLoginDate, createDate                                                                   sql.NullTime
		createWorld, bDeleted                                                                       int
		deletedDate                                                                                 sql.NullTime
		sID                                                                                         string
	)

	err = rows.Scan(
		&userNo, &slotNo, &sID, &level, &exp, &hp, &sp, &ap, &hps, &sps, &ps, &gs,
		&money, &fame, &flags, &adminLevel, &loginZone, &loginZoneX, &loginZoneY, &loginZoneD,
		&prisonMin, &pkCount, &pkYellowTime, &strength, &constitute, &dexterity, &intelligence,
		&mentalPower, &redistributePoint, &playMin, &lc, &lastLoginDate, &createDate,
		&createWorld, &bDeleted, &deletedDate,
	)
	if err != nil {
		return nil, fmt.Errorf("error scanning p_Char_GetAllData: %w", err)
	}
	_ = rows.Close()

	if bDeleted != 0 {
		return nil, ErrCharacterNotFound
	}

	ch.ID = id
	ch.AccountID = userNo
	ch.Name = sID
	ch.Slot = byte(slotNo)
	ch.Level = byte(level)
	ch.AdminLevel = byte(adminLevel)
	ch.CreateWorld = byte(createWorld)
	ch.TownName = loginZone
	ch.XPos = loginZoneX
	ch.YPos = loginZoneY
	ch.CurHP = hp
	ch.CurSP = sp
	if deletedDate.Valid && bDeleted != 0 {
		t := deletedDate.Time
		ch.DeletedAt = &t
	}

	ch.TCharacter = TCharacter{
		CharNo:            id,
		UserNo:            userNo,
		CharName:          sID,
		Slot:              byte(slotNo),
		Level:             byte(level),
		IsDelete:          bDeleted,
		Strength:          strength,
		Constitute:        constitute,
		Dexterity:         dexterity,
		Intelligence:      intelligence,
		MentalPower:       mentalPower,
		RedistributePoint: redistributePoint,
	}

	ch.TCharacterState = TCharacterState{
		CharNo:   id,
		XPos:     loginZoneX,
		YPos:     loginZoneY,
		Rotation: byte(loginZoneD),
		HP:       hp,
		SP:       sp,
		HPStones: int16(hps),
		SPStones: int16(sps),
		Exp:      exp,
		Money:    money,
		Fame:     fame,
	}

	// 2. Cargar tCharacterShape canónico
	shapeRows, err := r.db.Query("CALL p_Char_GetShape(?)", id)
	if err == nil {
		defer shapeRows.Close()
		if shapeRows.Next() {
			var race, class, gender, hairType, hairColor, faceShape int
			if err := shapeRows.Scan(&race, &class, &gender, &hairType, &hairColor, &faceShape); err == nil {
				ch.Race = byte(race)
				ch.Job = byte(class)
				ch.Male = (gender != 0)
				ch.Hair = byte(hairType)
				ch.HairColor = byte(hairColor)
				ch.Face = byte(faceShape)
				ch.TCharacter.Class = byte(class)
				ch.TCharacterShape = TCharacterShape{
					CharNo:    id,
					Gender:    byte(gender),
					Hair:      byte(hairType),
					HairColor: byte(hairColor),
					Face:      byte(faceShape),
				}
			}
		}
		_ = shapeRows.Close()
	}

	// 3. Cargar equipamiento activo (StorageType = 8 oficial de Zone/Character; fallback a 0 si aplica)
	type equippedKeySlot struct {
		itemKey int64
		slot    byte
	}
	var equippedKeys []equippedKeySlot
	for _, stType := range []byte{8, 0} {
		eqRows, eqErr := r.db.Query("CALL p_Item_GetListType(?, ?)", id, stType)
		if eqErr != nil {
			continue
		}
		for eqRows.Next() {
			var itemKey int64
			var storage int16
			var itemID, flags int
			var dDate time.Time
			if scanErr := eqRows.Scan(&itemKey, &storage, &itemID, &flags, &dDate); scanErr == nil {
				slot := byte(storage)
				if _, exists := ch.Equips[slot]; !exists {
					ch.Equips[slot] = uint16(itemID)
					equippedKeys = append(equippedKeys, equippedKeySlot{itemKey: itemKey, slot: slot})
				}
			}
		}
		_ = eqRows.Close()
		if len(ch.Equips) > 0 {
			break
		}
	}

	// 4. Cargar nivel de mejora (IOT_NUMBER = 0) de las piezas equipadas vía p_Item_GetOptionAll
	for _, ek := range equippedKeys {
		opts, optErr := r.GetItemOptions(ek.itemKey)
		if optErr == nil {
			if upg, ok := opts[0]; ok && upg > 0 {
				ch.EquipUpgrades[ek.slot] = byte(upg & 0x0F)
			}
		}
	}

	// 5. Cargar estado de tutorial del personaje vía usp_Tutorial_getCharInfo
	if tutState, tutStep, tutErr := r.GetTutorialCharInfo(id); tutErr == nil {
		ch.TutorialState = int32(tutState)
		ch.TutorialStep = tutStep
	}

	ch.Sync()
	return ch, nil
}

// GetCharacterLoginData recupera datos de sesión de login para un personaje invocando p_Char_GetLoginData.
func (r *MySQLCharacterRepository) GetCharacterLoginData(charNo int) (*Character, error) {
	rows, err := r.db.Query("CALL p_Char_GetLoginData(?)", charNo)
	if err != nil {
		return nil, fmt.Errorf("error calling p_Char_GetLoginData(%d): %w", charNo, err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, ErrCharacterNotFound
	}
	_ = rows.Close()

	return r.GetCharacterByID(charNo)
}

// CreateCharacter crea un avatar oficial ejecutando p_Char_Create y p_Char_CreateSetDefaultData.
func (r *MySQLCharacterRepository) CreateCharacter(ch *Character) error {
	if ch.Slot > 5 {
		return ErrSlotOutOfRange
	}
	if strings.TrimSpace(ch.Name) == "" {
		return fmt.Errorf("character name cannot be empty")
	}

	genderByte := byte(0)
	if ch.Male {
		genderByte = 1
	}

	createWorld := ch.CreateWorld
	if createWorld == 0 {
		createWorld = 1
	}

	var newCharNo int
	row := r.db.QueryRow(
		"CALL p_Char_Create(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, @out_nCharNo)",
		ch.AccountID, createWorld, ch.AdminLevel, ch.Slot, ch.Name, ch.Race, ch.Job, genderByte, ch.Hair, ch.HairColor, ch.Face,
	)
	if err := row.Scan(&newCharNo); err != nil {
		return fmt.Errorf("error calling p_Char_Create: %w", err)
	}

	if newCharNo == -1 {
		return ErrSlotInUse
	}
	if newCharNo == -2 {
		return ErrCharacterExists
	}
	if newCharNo <= 0 {
		return fmt.Errorf("p_Char_Create failed with code %d", newCharNo)
	}

	if ch.Level == 0 {
		ch.Level = 1
	}

	ch.ID = newCharNo
	ch.TCharacter.CharNo = newCharNo
	ch.TCharacter.UserNo = ch.AccountID
	ch.TCharacter.CharName = ch.Name
	ch.TCharacter.Slot = ch.Slot
	ch.TCharacter.Class = ch.Job
	ch.TCharacter.Level = ch.Level
	ch.TCharacter.IsDelete = 0

	ch.TCharacterShape.CharNo = newCharNo
	ch.TCharacterShape.Gender = genderByte
	ch.TCharacterShape.Hair = ch.Hair
	ch.TCharacterShape.HairColor = ch.HairColor
	ch.TCharacterShape.Face = ch.Face

	hps := ch.TCharacterState.HPStones
	if hps == 0 {
		hps = 15
	}
	sps := ch.TCharacterState.SPStones
	if sps == 0 {
		sps = 7
	}
	money := ch.TCharacterState.Money

	ch.TCharacterState.CharNo = newCharNo
	ch.TCharacterState.Map = ch.Map
	ch.TCharacterState.XPos = ch.XPos
	ch.TCharacterState.YPos = ch.YPos
	ch.TCharacterState.HP = ch.CurHP
	ch.TCharacterState.SP = ch.CurSP
	ch.TCharacterState.HPStones = hps
	ch.TCharacterState.SPStones = sps

	town := ch.TownName
	if town == "" {
		return fmt.Errorf("character town/login zone cannot be empty")
	}
	hp := ch.CurHP
	if hp <= 0 {
		hp = 100
	}
	sp := ch.CurSP
	if sp <= 0 {
		sp = 50
	}

	var ret int
	err := r.db.QueryRow(
		"CALL p_Char_CreateSetDefaultData(?, ?, ?, ?, ?, ?, ?, ?, ?, @out_nRet)",
		ch.ID, town, ch.XPos, ch.YPos, hp, sp, hps, sps, money,
	).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Char_CreateSetDefaultData: %w", err)
	}

	// Persistir nivel y experiencia inicial cuando la clase arranca sobre nivel 1 (ej. Sentinel Lv 60)
	if ch.Level > 1 || ch.TCharacterState.Exp > 0 {
		_ = r.SaveCharLevelExp(ch.ID, int16(ch.Level), ch.TCharacterState.Exp)
	}

	// Persistir atributos base cuando la clase arranca sobre nivel 1
	if ch.TCharacter.Strength > 0 || ch.TCharacter.Dexterity > 0 || ch.TCharacter.Constitute > 0 {
		_ = r.SaveCharFreeStats(ch.ID, byte(ch.TCharacter.Strength), byte(ch.TCharacter.Constitute), byte(ch.TCharacter.Dexterity), byte(ch.TCharacter.Intelligence), byte(ch.TCharacter.MentalPower), byte(ch.TCharacter.RedistributePoint))
	}

	// Persistir equipamiento inicial si fue provisto (StorageType = 8 oficial)
	if len(ch.Equips) > 0 {
		for slot, itemID := range ch.Equips {
			itemKey := time.Now().UnixNano() + int64(slot)
			_ = r.CreateItem(&Item{
				ItemKey:     itemKey,
				Owner:       ch.ID,
				StorageType: 8,
				Storage:     int16(slot),
				ItemID:      itemID,
				Flags:       0,
				Date:        time.Now(),
			})
		}
	}

	// Persistir items de inventario iniciales si fueron provistos (StorageType = 9)
	if len(ch.Items) > 0 {
		for i, it := range ch.Items {
			it.Owner = ch.ID
			if it.StorageType == 0 {
				it.StorageType = 9
			}
			if it.ItemKey == 0 {
				it.ItemKey = time.Now().UnixNano() + int64(100+i)
			}
			_ = r.CreateItem(it)
		}
	}

	// Persistir shortcuts clásicos si fueron provistos
	if len(ch.Shortcuts) > 0 {
		_ = r.SaveCharacterShortcuts(ch.ID, ch.Shortcuts)
	}

	// Persistir shortcuts mejorados NA2016 si fueron provistos
	for _, sc := range ch.ImproveShortcuts {
		_ = r.SetImproveShortCut(ch.ID, sc.SlotNo, sc.CodeNo, sc.Value)
	}

	// Persistir skills iniciales si fueron provistos
	if len(ch.Skills) > 0 {
		_ = r.SaveCharacterSkills(ch.ID, ch.Skills)
	}

	// Persistir misiones iniciales si fueron provistas
	for _, q := range ch.Quests {
		_ = r.AddCharacterQuest(ch.ID, q.QuestNo, q.Status, q.Data)
	}

	// Inicializar mapeo de teclas y opciones de juego NA2016 desde tablas semilla
	_ = r.InitImproveKeyMap(ch.ID)
	_ = r.InitImproveGameOptions(ch.ID)

	// Persistir estado de tutorial inicial si aplica
	if ch.TutorialState > 0 || ch.TutorialStep > 0 {
		_ = r.SetTutorialCharInfo(ch.AccountID, ch.ID, byte(ch.TutorialState), ch.TutorialStep)
	}

	ch.Sync()
	return nil
}

// UpdateCharacter sincroniza el estado completo del personaje en tCharacter invocando p_Char_SetAllData.
func (r *MySQLCharacterRepository) UpdateCharacter(ch *Character) error {
	ch.Sync()
	var ret int
	err := r.db.QueryRow(
		"CALL p_Char_SetAllData(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, @out_nRet)",
		ch.ID, ch.AccountID, ch.Slot, ch.Name, ch.Level, ch.TCharacterState.Exp,
		ch.CurHP, ch.CurSP, 0, ch.TCharacterState.HPStones, ch.TCharacterState.SPStones,
		0, 0, ch.TCharacterState.Money, ch.TCharacterState.Fame, 0, 0,
		ch.TownName, ch.XPos, ch.YPos, ch.TCharacterState.Rotation,
		0, 0, 0,
		ch.TCharacter.Strength, ch.TCharacter.Constitute, ch.TCharacter.Dexterity,
		ch.TCharacter.Intelligence, ch.TCharacter.MentalPower, ch.TCharacter.RedistributePoint,
		0, 0, time.Now(), time.Now(), 1, ch.TCharacter.IsDelete, time.Now(),
	).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Char_SetAllData: %w", err)
	}
	return nil
}

// DeleteCharacter realiza la baja lógica de un personaje por cuenta y slot ejecutando p_Char_Delete,
// verificando previamente que no sea maestro activo de gremio (p_Guild_GetMyNo + p_Guild_Member_Get).
func (r *MySQLCharacterRepository) DeleteCharacter(accountID int, slot byte) error {
	chars, err := r.GetCharactersByAccount(accountID)
	if err != nil {
		return err
	}
	var targetID int
	for _, c := range chars {
		if c.Slot == slot {
			targetID = c.ID
			break
		}
	}
	if targetID == 0 {
		return ErrCharacterNotFound
	}

	isMaster, err := r.IsActiveGuildMaster(targetID)
	if err != nil {
		return err
	}
	if isMaster {
		return ErrCharacterIsGuildMaster
	}

	var ret int
	err = r.db.QueryRow("CALL p_Char_Delete(?, @out_nRet)", targetID).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Char_Delete(%d): %w", targetID, err)
	}
	return nil
}

// ChangeCharacterName renombra un personaje verificando colisión con p_Char_GetCharNo e invocando p_Char_ChangeID.
func (r *MySQLCharacterRepository) ChangeCharacterName(charNo int, newName string) error {
	trimmed := strings.TrimSpace(newName)
	if existing, err := r.GetCharacterByName(trimmed); err == nil && existing != nil && existing.ID != charNo {
		return ErrCharacterExists
	}

	var ret int
	err := r.db.QueryRow("CALL p_Char_ChangeID(?, ?, @out_nRet)", charNo, trimmed).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Char_ChangeID: %w", err)
	}
	if ret == 0 {
		return ErrCharacterNotFound
	}
	return nil
}

// ChangeCharacterClass actualiza la clase de un personaje invocando p_Char_ChangeClass.
func (r *MySQLCharacterRepository) ChangeCharacterClass(charNo int, newClass byte) error {
	_, err := r.db.Exec("CALL p_Char_ChangeClass(?, ?)", charNo, newClass)
	if err != nil {
		return fmt.Errorf("error calling p_Char_ChangeClass: %w", err)
	}
	return nil
}

// SetCharacterShape actualiza la apariencia física de un personaje invocando p_Char_SetShape.
func (r *MySQLCharacterRepository) SetCharacterShape(charNo int, race, class, gender, hair, hairColor, face byte) error {
	var ret int
	err := r.db.QueryRow("CALL p_Char_SetShape(?, ?, ?, ?, ?, ?, ?, @out_nRet)", charNo, race, class, gender, hair, hairColor, face).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Char_SetShape: %w", err)
	}
	return nil
}

// CharacterLogin registra el ingreso de sesión de un personaje invocando p_Char_Login.
func (r *MySQLCharacterRepository) CharacterLogin(charNo int) error {
	var ret int
	err := r.db.QueryRow("CALL p_Char_Login(?, @out_nRet)", charNo).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Char_Login: %w", err)
	}
	return nil
}

// CharacterLogout registra la salida de sesión y tiempo jugado invocando p_Char_Logout.
func (r *MySQLCharacterRepository) CharacterLogout(charNo int, playMin int) error {
	var ret int
	err := r.db.QueryRow("CALL p_Char_Logout(?, ?, @out_nRet)", charNo, playMin).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Char_Logout: %w", err)
	}
	return nil
}

// SaveCharLocation persiste la posición y mapa de un personaje invocando p_Char_SaveLocation.
func (r *MySQLCharacterRepository) SaveCharLocation(charNo int, mapName string, x, y int, rotation byte) error {
	_, err := r.db.Exec("CALL p_Char_SaveLocation(?, ?, ?, ?, 0, '', 0, 0)", charNo, mapName, x, y)
	if err != nil {
		return fmt.Errorf("error calling p_Char_SaveLocation: %w", err)
	}
	return nil
}

// SaveCharStats persiste los vitales de combate (HP, SP, piedras) invocando p_Char_SaveStat.
func (r *MySQLCharacterRepository) SaveCharStats(charNo int, hp, sp, hps, sps int16) error {
	_, err := r.db.Exec("CALL p_Char_SaveStat(?, 0, 0, ?, ?, ?, ?)", charNo, hps, sps, hp, sp)
	if err != nil {
		return fmt.Errorf("error calling p_Char_SaveStat: %w", err)
	}
	return nil
}

// SaveCharLevelExp persiste el nivel y la experiencia acumulada invocando p_Char_SetLevelExp.
func (r *MySQLCharacterRepository) SaveCharLevelExp(charNo int, level int16, exp int64) error {
	var ret int
	err := r.db.QueryRow("CALL p_Char_SetLevelExp(?, ?, ?, @out_nRet)", charNo, level, exp).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Char_SetLevelExp: %w", err)
	}
	return nil
}

// SaveCharMoney actualiza el dinero del personaje invocando p_Char_MoneySet.
func (r *MySQLCharacterRepository) SaveCharMoney(charNo int, money int64) error {
	var ret int
	err := r.db.QueryRow("CALL p_Char_MoneySet(?, ?, @out_nRet)", charNo, money).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Char_MoneySet: %w", err)
	}
	return nil
}

// AddCharMoney añade dinero a la cuenta del personaje invocando p_Char_MoneyAdd.
func (r *MySQLCharacterRepository) AddCharMoney(charNo int, amount int64) error {
	var ret int
	err := r.db.QueryRow("CALL p_Char_MoneyAdd(?, ?, @out_nRet)", charNo, amount).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Char_MoneyAdd: %w", err)
	}
	return nil
}

// SubCharMoney descuenta dinero del personaje invocando p_Char_MoneySub.
func (r *MySQLCharacterRepository) SubCharMoney(charNo int, amount int64) (bool, error) {
	var ret int
	err := r.db.QueryRow("CALL p_Char_MoneySub(?, ?, @out_nRet)", charNo, amount).Scan(&ret)
	if err != nil {
		return false, fmt.Errorf("error calling p_Char_MoneySub: %w", err)
	}
	return ret == 1, nil
}

// SaveCharFreeStats persiste los puntos de atributos redistribuidos invocando usp_Character_setStatPoint.
func (r *MySQLCharacterRepository) SaveCharFreeStats(charNo int, str, con, dex, intel, men, freePoints byte) error {
	var ret int
	err := r.db.QueryRow(
		"CALL usp_Character_setStatPoint(?, 500, ?, ?, ?, ?, ?, 0, 0, 0, @out_nRet)",
		charNo, str, con, dex, intel, men,
	).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling usp_Character_setStatPoint: %w", err)
	}
	return nil
}

// GetCharacterStates obtiene la lista de estados anormales persistidos invocando p_Char_GetStateAll.
func (r *MySQLCharacterRepository) GetCharacterStates(charNo int) ([]*CharacterStateData, error) {
	rows, err := r.db.Query("CALL p_Char_GetStateAll(?)", charNo)
	if err != nil {
		return nil, fmt.Errorf("error calling p_Char_GetStateAll(%d): %w", charNo, err)
	}
	defer rows.Close()

	var states []*CharacterStateData
	for rows.Next() {
		s := &CharacterStateData{CharNo: charNo}
		if err := rows.Scan(
			&s.StateNo, &s.StateLevel, &s.StateTime, &s.StateCount,
			&s.PowerDamage, &s.PowerSP, &s.PowerKeepTime, &s.PowerCoolTime,
		); err == nil {
			states = append(states, s)
		}
	}
	return states, nil
}

// SetCharacterState almacena un estado anormal con tiempo restante invocando p_Char_SetStateWithRestTime.
func (r *MySQLCharacterRepository) SetCharacterState(charNo int, stateNo int16, level int16, stateTime int, restTime int, count int16, pd, ps, pk, pc byte) error {
	var ret int
	err := r.db.QueryRow(
		"CALL p_Char_SetStateWithRestTime(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, @out_nRet)",
		charNo, stateNo, level, stateTime, restTime, count, pd, ps, pk, pc,
	).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Char_SetStateWithRestTime: %w", err)
	}
	return nil
}

// GetCharacterTitles recupera los títulos activos invocando p_Char_GetTitleAll.
func (r *MySQLCharacterRepository) GetCharacterTitles(charNo int) ([]*CharacterTitleData, error) {
	rows, err := r.db.Query("CALL p_Char_GetTitleAll(?)", charNo)
	if err != nil {
		return nil, fmt.Errorf("error calling p_Char_GetTitleAll(%d): %w", charNo, err)
	}
	defer rows.Close()

	var titles []*CharacterTitleData
	for rows.Next() {
		var titleType int16
		var title string
		if err := rows.Scan(&titleType, &title); err == nil {
			titles = append(titles, &CharacterTitleData{CharNo: charNo, TitleType: titleType, Title: title})
		}
	}
	return titles, nil
}

// SetCharacterTitle asigna un título de personaje invocando p_Char_SetTitle.
func (r *MySQLCharacterRepository) SetCharacterTitle(charNo int, titleType int16, title string) error {
	var ret int
	err := r.db.QueryRow("CALL p_Char_SetTitle(?, ?, 0, 0, @out_nRet)", charNo, titleType).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Char_SetTitle: %w", err)
	}
	return nil
}

// GetChargedBuffs recupera los buffs pagos de tCharacterChargedBuffer invocando p_Char_GetChargedBufferAll.
func (r *MySQLCharacterRepository) GetChargedBuffs(charNo int) ([]*ChargedBuffData, error) {
	rows, err := r.db.Query("CALL p_Char_GetChargedBufferAll(?)", charNo)
	if err != nil {
		return nil, fmt.Errorf("error calling p_Char_GetChargedBufferAll(%d): %w", charNo, err)
	}
	defer rows.Close()

	var buffs []*ChargedBuffData
	for rows.Next() {
		b := &ChargedBuffData{}
		if err := rows.Scan(&b.Key, &b.ID, &b.UseDate, &b.EndDate); err == nil {
			buffs = append(buffs, b)
		}
	}
	return buffs, nil
}

// SetChargedBuff registra o actualiza un buff pago invocando p_Char_SetChargedBuffer.
func (r *MySQLCharacterRepository) SetChargedBuff(charNo int, key int, id int16, useDate, endDate time.Time) error {
	var ret int
	err := r.db.QueryRow("CALL p_Char_SetChargedBuffer(?, ?, ?, ?, ?, @out_nRet)", charNo, key, id, useDate, endDate).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Char_SetChargedBuffer: %w", err)
	}
	return nil
}

// DelChargedBuff elimina un buff pago invocando p_Char_DelChargedBuff.
func (r *MySQLCharacterRepository) DelChargedBuff(charNo int, key int) error {
	var ret int
	err := r.db.QueryRow("CALL p_Char_DelChargedBuff(?, ?, @out_nRet)", charNo, key).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Char_DelChargedBuff: %w", err)
	}
	return nil
}

// GetRestExpExecTime obtiene la última fecha de ejecución de RestExp invocando usp_Char_GetRestExpExecTime.
func (r *MySQLCharacterRepository) GetRestExpExecTime(charNo int) (time.Time, error) {
	var t time.Time
	row := r.db.QueryRow("CALL usp_Char_GetRestExpExecTime(?, @out_nRet)", charNo)
	if err := row.Scan(&t); err != nil {
		return time.Now(), nil
	}
	return t, nil
}

// SetRestExpExecTime actualiza la fecha de RestExp invocando usp_Char_SetRestExpExecTime.
func (r *MySQLCharacterRepository) SetRestExpExecTime(charNo int, t time.Time) error {
	var ret int
	err := r.db.QueryRow("CALL usp_Char_SetRestExpExecTime(?, @out_nRet)", charNo).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling usp_Char_SetRestExpExecTime: %w", err)
	}
	return nil
}

// GetAutoPick consulta el estado de recolección automática invocando usp_Char_GetAutoPick.
func (r *MySQLCharacterRepository) GetAutoPick(charNo int) (bool, error) {
	var autoPick int
	row := r.db.QueryRow("CALL usp_Char_GetAutoPick(?, @out_nRet)", charNo)
	if err := row.Scan(&autoPick); err != nil {
		return false, nil
	}
	return autoPick != 0, nil
}

// SetAutoPick actualiza el estado de recolección automática invocando usp_Char_SetAutoPick.
func (r *MySQLCharacterRepository) SetAutoPick(charNo int, autoPick bool) error {
	val := byte(0)
	if autoPick {
		val = 1
	}
	var ret int
	err := r.db.QueryRow("CALL usp_Char_SetAutoPick(?, ?, @out_nRet)", charNo, val).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling usp_Char_SetAutoPick: %w", err)
	}
	return nil
}

// SaveCharLink actualiza el mapa de retorno de LinkScroll invocando usp_LinkScroll_SetDestMap.
func (r *MySQLCharacterRepository) SaveCharLink(charNo int, linkZone string, linkX, linkY int) error {
	var ret int
	err := r.db.QueryRow("CALL usp_LinkScroll_SetDestMap(?, ?, @out_nRet)", charNo, linkZone).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling usp_LinkScroll_SetDestMap: %w", err)
	}
	return nil
}

// AddCharExp añade experiencia al personaje invocando p_Char_AddExp.
func (r *MySQLCharacterRepository) AddCharExp(charNo int, exp int64) error {
	_, err := r.db.Exec("CALL p_Char_AddExp(?, ?)", charNo, exp)
	if err != nil {
		return fmt.Errorf("error calling p_Char_AddExp: %w", err)
	}
	return nil
}

// AddCharFame incrementa la fama del personaje invocando p_Char_FameAdd.
func (r *MySQLCharacterRepository) AddCharFame(charNo int, fame int) error {
	var ret int
	err := r.db.QueryRow("CALL p_Char_FameAdd(?, ?, @out_nRet)", charNo, fame).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Char_FameAdd: %w", err)
	}
	return nil
}

// SubCharFame decrementa la fama del personaje invocando p_Char_FameSub.
func (r *MySQLCharacterRepository) SubCharFame(charNo int, fame int) error {
	var ret int
	err := r.db.QueryRow("CALL p_Char_FameSub(?, ?, @out_nRet)", charNo, fame).Scan(&ret)
	if err != nil {
		return fmt.Errorf("error calling p_Char_FameSub: %w", err)
	}
	return nil
}

// GetCharFame recupera la fama del personaje invocando p_Char_FameGet.
func (r *MySQLCharacterRepository) GetCharFame(charNo int) (int, error) {
	var fame, ret int
	err := r.db.QueryRow("CALL p_Char_FameGet(?, @out_nFame, @out_nRet)", charNo).Scan(&fame, &ret)
	if err != nil {
		return 0, fmt.Errorf("error calling p_Char_FameGet: %w", err)
	}
	return fame, nil
}

// GetCharacterChangeIDList consulta los personajes con cambio obligatorio/habilitado de nombre invocando usp_Character_getChangeIDList.
func (r *MySQLCharacterRepository) GetCharacterChangeIDList(userNo int) ([]*CharIDChangeInfo, error) {
	rows, err := r.db.Query("CALL usp_Character_getChangeIDList(?, @out_nRet)", userNo)
	if err != nil {
		return nil, fmt.Errorf("error calling usp_Character_getChangeIDList(%d): %w", userNo, err)
	}
	defer rows.Close()

	var list []*CharIDChangeInfo
	for rows.Next() {
		var rowNo uint32
		var uNo, cNo int
		var bInit byte
		if scanErr := rows.Scan(&rowNo, &uNo, &cNo, &bInit); scanErr == nil {
			list = append(list, &CharIDChangeInfo{
				RowNo:        rowNo,
				UserNo:       uNo,
				CharNo:       cNo,
				NeedChangeID: 1,
				Init:         bInit,
			})
		}
	}
	return list, nil
}

// ExecuteCharacterChangeID ejecuta el cambio transaccional de nombre en tCharacterChangeID invocando usp_Character_setChangeID.
func (r *MySQLCharacterRepository) ExecuteCharacterChangeID(rowNo int, charNo int, newName string, isInit byte, userNo int) (int, error) {
	var ret int
	err := r.db.QueryRow(
		"CALL usp_Character_setChangeID(?, ?, ?, ?, ?, @out_nRet)",
		rowNo, charNo, strings.TrimSpace(newName), isInit, userNo,
	).Scan(&ret)
	if err != nil {
		return -3026, fmt.Errorf("error calling usp_Character_setChangeID: %w", err)
	}
	return ret, nil
}

// IsActiveGuildMaster verifica si un personaje es líder activo (nGrade = 0) de un gremio invocando p_Guild_GetMyNo y p_Guild_Member_Get.
func (r *MySQLCharacterRepository) IsActiveGuildMaster(charNo int) (bool, error) {
	var guildNo int
	err := r.db.QueryRow("CALL p_Guild_GetMyNo(?, @out_nNo)", charNo).Scan(&guildNo)
	if err != nil || guildNo <= 0 {
		return false, nil
	}

	rows, err := r.db.Query("CALL p_Guild_Member_Get(?, ?)", guildNo, charNo)
	if err != nil {
		return false, fmt.Errorf("error calling p_Guild_Member_Get(%d, %d): %w", guildNo, charNo, err)
	}
	defer rows.Close()

	if rows.Next() {
		var (
			joinDate, lastLogin sql.NullTime
			grade, flags        int
			expRatio            int16
			intro, sID, mapName string
			class, level        int
		)
		if scanErr := rows.Scan(&joinDate, &grade, &flags, &expRatio, &intro, &sID, &class, &level, &mapName, &lastLogin); scanErr == nil {
			return grade == 0, nil
		}
	}
	return false, nil
}

// GetPrisonInfo consulta si el personaje posee una sanción en tPrison invocando p_Prison_Get.
func (r *MySQLCharacterRepository) GetPrisonInfo(charNo int) (*PrisonInfo, error) {
	rows, err := r.db.Query("CALL p_Prison_Get(?)", charNo)
	if err != nil {
		return nil, fmt.Errorf("error calling p_Prison_Get(%d): %w", charNo, err)
	}
	defer rows.Close()

	if rows.Next() {
		var reason, remark string
		if scanErr := rows.Scan(&reason, &remark); scanErr == nil {
			return &PrisonInfo{Reason: reason, Remark: remark}, nil
		}
	}
	return nil, nil
}
