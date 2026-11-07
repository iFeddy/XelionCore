package world

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/iFeddy/XelionCore/pkg/data/classparam"
	"github.com/iFeddy/XelionCore/pkg/data/field"
	"github.com/iFeddy/XelionCore/pkg/data/shn"
)

// CharacterSHNFiles lista los 54 archivos .shn verificados por shndatafilechecksum.cpp en Character.exe.
var CharacterSHNFiles = []string{
	"AbStateView.shn", "Abstate.shn", "ActionViewInfo.shn", "ActiveSkill.shn", "ActiveSkillView.shn", "AttendReward.shn",
	"CharacterTitleData.shn", "CharacterTitleStateView.shn", "ChargedDeletableBuff.shn", "ChargedEffect.shn", "ChatColor.shn", "ClassName.shn",
	"CollectCardView.shn", "DiceDividind.shn", "EffectViewInfo.shn", "FaceInfo.shn", "GTIView.shn", "Gather.shn",
	"GradeItemOption.shn", "HairColorInfo.shn", "HairInfo.shn", "HolyPromiseReward.shn", "ItemDismantle.shn", "ItemInfo.shn",
	"ItemMoney.shn", "ItemShopView.shn", "ItemViewEquipTypeInfo.shn", "ItemViewInfo.shn", "MapInfo.shn", "MapLinkPoint.shn",
	"MapViewInfo.shn", "MapWayPoint.shn", "MarketSearchInfo.shn", "MiniHouse.shn", "MiniHouseEndure.shn", "MiniHouseFurniture.shn",
	"MiniHouseFurnitureObjEffect.shn", "MiniHouseObjAni.shn", "MinimonAutoUseItem.shn", "MinimonInfo.shn", "MobInfo.shn", "MobViewInfo.shn",
	"NPCViewInfo.shn", "PassiveSkill.shn", "PassiveSkillView.shn", "ProduceView.shn", "PupMain.shn", "Riding.shn",
	"SingleData.shn", "SubAbstate.shn", "TermExtendMatch.shn", "UpgradeInfo.shn", "WeaponAttrib.shn", "WeaponTitleData.shn",
}

// WorldManagerSHNFiles lista los 71 archivos .shn oficiales cargados por WorldManager.exe.
var WorldManagerSHNFiles = []string{
	"AbStateView.shn", "Abstate.shn", "ActionViewInfo.shn", "ActiveSkill.shn", "ActiveSkillView.shn", "AttendReward.shn",
	"AttendSchedule.shn", "BadNameFilter.shn", "CharacterTitleData.shn", "CharacterTitleStateView.shn", "ChargedDeletableBuff.shn", "ChargedEffect.shn",
	"ChatColor.shn", "ClassName.shn", "CollectCardView.shn", "DiceDividind.shn", "EffectViewInfo.shn", "GBHouse.shn",
	"GTIView.shn", "Gather.shn", "GradeItemOption.shn", "GuildAcademy.shn", "GuildAcademyLevelUp.shn", "GuildAcademyRank.shn",
	"GuildGradeData.shn", "GuildGradeScoreData.shn", "GuildLevelScoreData.shn", "GuildTournament.shn", "GuildTournamentRequire.shn", "GuildTournamentReward.shn",
	"ItemDismantle.shn", "ItemInfo.shn", "ItemMoney.shn", "ItemShopView.shn", "ItemViewEquipTypeInfo.shn", "ItemViewInfo.shn",
	"KQIsVote.shn", "KQTeam.shn", "KQVoteMajorityRate.shn", "KingdomQuest.shn", "KingdomQuestMap.shn", "MapInfo.shn",
	"MapLinkPoint.shn", "MapViewInfo.shn", "MapWayPoint.shn", "MarketSearchInfo.shn", "MiniHouse.shn", "MiniHouseEndure.shn",
	"MiniHouseFurniture.shn", "MiniHouseFurnitureObjEffect.shn", "MiniHouseObjAni.shn", "MinimonAutoUseItem.shn", "MinimonInfo.shn", "MobInfo.shn",
	"MobViewInfo.shn", "MsgWorldManager.shn", "NPCViewInfo.shn", "NpcSchedule.shn", "PassiveSkill.shn", "PassiveSkillView.shn",
	"ProduceView.shn", "PupMain.shn", "Riding.shn", "SingleData.shn", "SpamerReport.shn", "SubAbstate.shn",
	"TermExtendMatch.shn", "UpgradeInfo.shn", "UseClassTypeInfo.shn", "WeaponAttrib.shn", "WeaponTitleData.shn",
}

// StarterItem representa un objeto de inicio asignado a la nueva clase.
type StarterItem struct {
	ItemID uint16 `json:"item_id"`
	Amount int    `json:"amount"`
	Name   string `json:"name"`
}

// StarterShortcut representa un acceso directo inicial en la barra de habilidades.
type StarterShortcut struct {
	Slot   byte   `json:"slot"`
	CodeNo uint16 `json:"code_no"`
	Value  int    `json:"value"`
}

// MaxShape define los límites estéticos oficiales por clase.
type MaxShape struct {
	MaxFace      byte `json:"max_face"`
	MaxHair      byte `json:"max_hair"`
	MaxHairColor byte `json:"max_hair_color"`
}

// StarterQuest representa una misión inicial asignada por DefaultCharacterData.txt (#DEFINE QUEST).
type StarterQuest struct {
	ClassID byte `json:"class_id"`
	QuestID int  `json:"quest_id"`
}

// DefaultClassData contiene los parámetros base oficiales por clase leídos de default_character_data.json.
type DefaultClassData struct {
	ClassID       byte              `json:"class_id"`
	ClassName     string            `json:"class_name"`
	StartMap      string            `json:"start_map"`
	StartX        int               `json:"start_x"`
	StartY        int               `json:"start_y"`
	HP            int               `json:"hp"`
	SP            int               `json:"sp"`
	HPSoul        int16             `json:"hp_soul"`
	SPSoul        int16             `json:"sp_soul"`
	Money         int64             `json:"money"`
	InitLevel     byte              `json:"init_level"`
	InitExp       int64             `json:"init_exp"`
	MaxShape      MaxShape          `json:"max_shape"`
	StarterItems  []StarterItem     `json:"starter_items"`
	TutorialItems []StarterItem     `json:"tutorial_items"`
	Shortcuts     []StarterShortcut `json:"shortcuts"`
	StarterSkills []int             `json:"starter_skills"`
	StarterQuests []int             `json:"starter_quests"`
}

// DefaultCharacterDataStore almacena las configuraciones de todas las clases disponibles.
type DefaultCharacterDataStore struct {
	Classes       map[string]DefaultClassData `json:"classes"`
	Tutorials     map[string]DefaultClassData `json:"tutorials,omitempty"`
	ClassSkills   map[string][]int            `json:"class_skills,omitempty"`
	StarterQuests []StarterQuest              `json:"starter_quests"`
}

// LoadDefaultCharacterData lee el JSON de datos iniciales oficiales de personajes.
func LoadDefaultCharacterData(filePath string) (*DefaultCharacterDataStore, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", filePath, err)
	}

	var store DefaultCharacterDataStore
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, fmt.Errorf("failed to unmarshal %s: %w", filePath, err)
	}

	return &store, nil
}

// GetClassData busca la configuración base para una clase (1: Fighter, 6: Cleric, etc.).
func (s *DefaultCharacterDataStore) GetClassData(classID byte) (DefaultClassData, bool) {
	if s == nil || s.Classes == nil {
		return DefaultClassData{}, false
	}
	key := fmt.Sprintf("%d", classID)
	data, exists := s.Classes[key]
	if !exists && s.Tutorials != nil {
		data, exists = s.Tutorials[key]
	}
	return data, exists
}

// GetTutorialData busca la configuración de tutorial para una clase.
func (s *DefaultCharacterDataStore) GetTutorialData(classID byte) (DefaultClassData, bool) {
	if s == nil {
		return DefaultClassData{}, false
	}
	key := fmt.Sprintf("%d", classID)
	if s.Tutorials != nil {
		if data, exists := s.Tutorials[key]; exists {
			return data, true
		}
	}
	if s.Classes != nil {
		if data, exists := s.Classes[key]; exists {
			return data, true
		}
	}
	return DefaultClassData{}, false
}

// BadNameFilter gestiona las palabras y nombres restringidos por el motor oficial exclusivamente vía BadNameFilter.shn.
type BadNameFilter struct {
	exactBadNames     map[string]bool
	substringBadNames []string
}

// LoadBadNameFilter lee BadNameFilter.shn oficial con el parser binario de XelionCore sin cadenas hardcodeadas.
func LoadBadNameFilter(filePath string) (*BadNameFilter, error) {
	file, err := shn.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load %s: %w", filePath, err)
	}

	filter := &BadNameFilter{
		exactBadNames:     make(map[string]bool, file.RecordCount),
		substringBadNames: make([]string, 0, file.RecordCount),
	}

	for _, row := range file.Rows {
		name := strings.TrimSpace(row.GetString("BadName"))
		if name == "" {
			continue
		}
		lower := strings.ToLower(name)
		nType := row.GetInt("Type")
		if nType == 0 {
			filter.exactBadNames[lower] = true
		} else {
			filter.substringBadNames = append(filter.substringBadNames, lower)
		}
	}

	return filter, nil
}

// IsBadName valida si un nombre propuesto para personaje infringe las normas oficiales.
func (f *BadNameFilter) IsBadName(name string) bool {
	if f == nil {
		return false
	}

	lower := strings.ToLower(strings.TrimSpace(name))

	// 1. Coincidencia exacta en lista negra oficial (Type == 0)
	if f.exactBadNames[lower] {
		return true
	}

	// 2. Subcadena prohibida según BadNameFilter.shn (Type != 0)
	for _, sub := range f.substringBadNames {
		if strings.Contains(lower, sub) {
			return true
		}
	}

	return false
}

// AttendRewardDef representa la configuración de recompensa diaria por AttendReward.shn.
type AttendRewardDef struct {
	Index   byte
	Type    int
	Count   byte
	ItemID  uint16
	ItemInx string
}

// WorldDataStore administra todos los datos estáticos, tablas SHN y validadores oficiales de WorldService.
type WorldDataStore struct {
	DataDir      string
	BadNames     *BadNameFilter
	DefaultData  *DefaultCharacterDataStore
	TutorialData *DefaultCharacterDataStore
	Checksums    *shn.ChecksumEngine

	// Stores de mapas y parámetros de clase
	Fields      *field.Store
	ClassParams *classparam.Store

	// 4 archivos específicos de Character.exe
	FaceInfo          *shn.File
	HairInfo          *shn.File
	HairColorInfo     *shn.File
	HolyPromiseReward *shn.File

	// 22 archivos de coordinación de WorldManager.exe
	GuildAcademy           *shn.File
	GuildAcademyLevelUp     *shn.File
	GuildAcademyRank       *shn.File
	GuildGradeData         *shn.File
	GuildGradeScoreData    *shn.File
	GuildLevelScoreData    *shn.File
	GuildTournament        *shn.File
	GuildTournamentRequire *shn.File
	GuildTournamentReward  *shn.File
	GBHouse                *shn.File
	KingdomQuest           *shn.File
	KingdomQuestMap        *shn.File
	KQTeam                 *shn.File
	KQIsVote               *shn.File
	KQVoteMajorityRate     *shn.File
	SpamerReport           *shn.File
	AttendReward           *shn.File
	AttendSchedule         *shn.File
	NpcSchedule            *shn.File
	SingleData             *shn.File
	UseClassTypeInfo       *shn.File
	MsgWorldManager        *shn.File

	// Tablas globales críticas
	ItemInfo  *shn.File
	ClassName *shn.File
	MapInfo   *shn.File

	// 74+ archivos SHN únicos cargados y cacheados
	LoadedFiles map[string]*shn.File

	// Mappings rápidos de validación y resolución desde SHN
	validFaces    map[byte]bool
	validHairs    map[byte]bool
	validColors   map[byte]bool
	faceGrades    map[byte]byte
	hairGrades    map[byte]byte
	colorGrades   map[byte]byte
	itemEquips    map[uint16]byte
	inxToItemID   map[string]uint16
	styleCoupons  [6]uint16
	attendRewards map[byte]AttendRewardDef
	classNames    map[byte]string
}

// LoadWorldDataStore inicializa y valida el ecosistema estático completo de WorldService.
func LoadWorldDataStore(dataDir string) (*WorldDataStore, error) {
	store := &WorldDataStore{
		DataDir:       dataDir,
		LoadedFiles:   make(map[string]*shn.File),
		validFaces:    make(map[byte]bool),
		validHairs:    make(map[byte]bool),
		validColors:   make(map[byte]bool),
		faceGrades:    make(map[byte]byte),
		hairGrades:    make(map[byte]byte),
		colorGrades:   make(map[byte]byte),
		itemEquips:    make(map[uint16]byte),
		inxToItemID:   make(map[string]uint16),
		attendRewards: make(map[byte]AttendRewardDef),
		classNames:    make(map[byte]string),
	}

	// 1. Cargar BadNameFilter exclusivamente desde BadNameFilter.shn
	badNamePath := filepath.Join(dataDir, "BadNameFilter.shn")
	if filter, err := LoadBadNameFilter(badNamePath); err == nil {
		store.BadNames = filter
	} else {
		log.Printf("[WorldDataStore] Warning: Failed to load %s: %v", badNamePath, err)
	}

	// 2. Cargar default_character_data.json
	defaultPath := filepath.Join(dataDir, "default_character_data.json")
	if defData, err := LoadDefaultCharacterData(defaultPath); err == nil {
		store.DefaultData = defData
	} else {
		log.Printf("[WorldDataStore] Warning: Failed to load %s: %v", defaultPath, err)
	}

	// 3. Cargar tutorial_character_data.json
	tutorialPath := filepath.Join(dataDir, "tutorial_character_data.json")
	if tutData, err := LoadDefaultCharacterData(tutorialPath); err == nil {
		store.TutorialData = tutData
	} else {
		log.Printf("[WorldDataStore] Warning: Failed to load %s: %v", tutorialPath, err)
	}

	// 4. Cargar Fields (World/Map)
	mapDir := filepath.Join(dataDir, "World", "Map")
	if fStore, err := field.Load(mapDir); err == nil {
		store.Fields = fStore
	} else {
		log.Printf("[WorldDataStore] Warning: Failed to load maps from %s: %v", mapDir, err)
	}

	// 5. Cargar ClassParams (World/ClassParam, 27 clases)
	classParamDir := filepath.Join(dataDir, "World", "ClassParam")
	if cpStore, err := classparam.Load(classParamDir); err == nil {
		store.ClassParams = cpStore
	} else {
		log.Printf("[WorldDataStore] Warning: Failed to load classparams from %s: %v", classParamDir, err)
	}

	// 6. Inicializar motor de Checksums de 49 archivos
	if chkEngine, err := shn.NewChecksumEngine(dataDir, filepath.Join(dataDir, "View")); err == nil {
		store.Checksums = chkEngine
	} else {
		log.Printf("[WorldDataStore] Warning: Failed to init ChecksumEngine: %v", err)
	}

	// 7. Cargar el universo completo de los 74+ archivos SHN únicos (54 de Character U 71 de WorldManager)
	uniqueSHNs := make(map[string]bool)
	for _, name := range CharacterSHNFiles {
		uniqueSHNs[name] = true
	}
	for _, name := range WorldManagerSHNFiles {
		uniqueSHNs[name] = true
	}

	for name := range uniqueSHNs {
		p := filepath.Join(dataDir, name)
		if _, err := os.Stat(p); os.IsNotExist(err) {
			p = filepath.Join(dataDir, "View", name)
		}
		f, err := shn.ReadFile(p)
		if err != nil {
			log.Printf("[WorldDataStore] Warning: Could not load %s: %v", name, err)
			continue
		}
		store.LoadedFiles[name] = f
	}

	// Asignar referencias tipadas de compatibilidad
	store.FaceInfo = store.LoadedFiles["FaceInfo.shn"]
	store.HairInfo = store.LoadedFiles["HairInfo.shn"]
	store.HairColorInfo = store.LoadedFiles["HairColorInfo.shn"]
	store.HolyPromiseReward = store.LoadedFiles["HolyPromiseReward.shn"]
	store.GuildAcademy = store.LoadedFiles["GuildAcademy.shn"]
	store.GuildAcademyLevelUp = store.LoadedFiles["GuildAcademyLevelUp.shn"]
	store.GuildAcademyRank = store.LoadedFiles["GuildAcademyRank.shn"]
	store.GuildGradeData = store.LoadedFiles["GuildGradeData.shn"]
	store.GuildGradeScoreData = store.LoadedFiles["GuildGradeScoreData.shn"]
	store.GuildLevelScoreData = store.LoadedFiles["GuildLevelScoreData.shn"]
	store.GuildTournament = store.LoadedFiles["GuildTournament.shn"]
	store.GuildTournamentRequire = store.LoadedFiles["GuildTournamentRequire.shn"]
	store.GuildTournamentReward = store.LoadedFiles["GuildTournamentReward.shn"]
	store.GBHouse = store.LoadedFiles["GBHouse.shn"]
	store.KingdomQuest = store.LoadedFiles["KingdomQuest.shn"]
	store.KingdomQuestMap = store.LoadedFiles["KingdomQuestMap.shn"]
	store.KQTeam = store.LoadedFiles["KQTeam.shn"]
	store.KQIsVote = store.LoadedFiles["KQIsVote.shn"]
	store.KQVoteMajorityRate = store.LoadedFiles["KQVoteMajorityRate.shn"]
	store.SpamerReport = store.LoadedFiles["SpamerReport.shn"]
	store.AttendReward = store.LoadedFiles["AttendReward.shn"]
	store.AttendSchedule = store.LoadedFiles["AttendSchedule.shn"]
	store.NpcSchedule = store.LoadedFiles["NpcSchedule.shn"]
	store.SingleData = store.LoadedFiles["SingleData.shn"]
	store.UseClassTypeInfo = store.LoadedFiles["UseClassTypeInfo.shn"]
	store.MsgWorldManager = store.LoadedFiles["MsgWorldManager.shn"]
	store.ItemInfo = store.LoadedFiles["ItemInfo.shn"]
	store.ClassName = store.LoadedFiles["ClassName.shn"]
	store.MapInfo = store.LoadedFiles["MapInfo.shn"]

	// 8. Construir tablas de búsqueda rápida de apariencia (FaceInfo, HairInfo, HairColorInfo)
	if store.FaceInfo != nil {
		for _, row := range store.FaceInfo.Rows {
			id := byte(row.GetInt("ID"))
			grade := byte(row.GetInt("Grade"))
			if grade != 255 {
				store.faceGrades[id] = grade
				if grade == 0 {
					store.validFaces[id] = true
				}
			}
		}
	}
	if store.HairInfo != nil {
		for _, row := range store.HairInfo.Rows {
			id := byte(row.GetInt("ID"))
			grade := byte(row.GetInt("Grade"))
			if grade != 255 {
				store.hairGrades[id] = grade
				if grade == 0 {
					store.validHairs[id] = true
				}
			}
		}
	}
	if store.HairColorInfo != nil {
		for _, row := range store.HairColorInfo.Rows {
			id := byte(row.GetInt("ID"))
			grade := byte(row.GetInt("Grade"))
			if grade != 255 {
				store.colorGrades[id] = grade
				if grade == 0 {
					store.validColors[id] = true
				}
			}
		}
	}

	// 9. Construir tablas de búsqueda rápida de ítems y cupones de belleza (ItemInfo.shn)
	if store.ItemInfo != nil {
		for _, row := range store.ItemInfo.Rows {
			id := row.GetUint16("ID")
			equip := byte(row.GetInt("Equip"))
			inxName := strings.TrimSpace(row.GetString("InxName"))
			store.itemEquips[id] = equip
			if inxName != "" {
				store.inxToItemID[inxName] = id
			}
			switch inxName {
			case "HairShop01":
				store.styleCoupons[1] = id
			case "HairShop02":
				store.styleCoupons[2] = id
			case "HairShop03":
				store.styleCoupons[3] = id
			case "HairShop04":
				store.styleCoupons[4] = id
			case "HairShop05":
				store.styleCoupons[5] = id
			}
		}
	}

	// 10. Construir tabla de nombres de clase oficiales (ClassName.shn)
	if store.ClassName != nil {
		for _, row := range store.ClassName.Rows {
			classID := byte(row.GetInt("ClassID"))
			engName := strings.TrimSpace(row.GetString("acEngName"))
			if classID > 0 && engName != "" {
				store.classNames[classID] = engName
			}
		}
	}

	// 11. Construir tabla de recompensas de asistencia diaria (AttendReward.shn)
	if store.AttendReward != nil {
		for _, row := range store.AttendReward.Rows {
			arID := byte(row.GetInt("AR_ID"))
			arType := row.GetInt("AR_Type")
			arCount := byte(row.GetInt("AR_Count"))
			itemInx := strings.TrimSpace(row.GetString("AR_ItemInx"))
			var itemID uint16
			if id, ok := store.inxToItemID[itemInx]; ok {
				itemID = id
			}
			store.attendRewards[arID] = AttendRewardDef{
				Index:   arID,
				Type:    arType,
				Count:   arCount,
				ItemID:  itemID,
				ItemInx: itemInx,
			}
		}
	}

	return store, nil
}

// ValidateAppearance valida formalmente que la cara, peinado y color solicitados
// para la creación del personaje existan en las tablas oficiales y sean de Grade == 0 (gratuitas).
func (s *WorldDataStore) ValidateAppearance(classID byte, hair, hairColor, face byte) bool {
	if s == nil {
		return true
	}

	// 1. Validar Face en FaceInfo.shn oficial
	if len(s.validFaces) > 0 && !s.validFaces[face] {
		return false
	}

	// 2. Validar Hair en HairInfo.shn oficial
	if len(s.validHairs) > 0 && !s.validHairs[hair] {
		return false
	}

	// 3. Validar HairColor en HairColorInfo.shn oficial
	if len(s.validColors) > 0 && !s.validColors[hairColor] {
		return false
	}

	return true
}

// ValidateStyleChange valida el cambio de apariencia en la tienda de belleza y determina el grado de cupón requerido.
func (s *WorldDataStore) ValidateStyleChange(hair, hairColor, face byte) (requiredGrade byte, ok bool) {
	if s == nil {
		return 0, true
	}

	fg, fOk := s.faceGrades[face]
	hg, hOk := s.hairGrades[hair]
	cg, cOk := s.colorGrades[hairColor]
	if !fOk || !hOk || !cOk {
		return 0, false
	}

	maxGrade := fg
	if hg > maxGrade {
		maxGrade = hg
	}
	if cg > maxGrade {
		maxGrade = cg
	}
	return maxGrade, true
}

// GetItemEquipSlot retorna el slot oficial de equipamiento de un ítem según ItemInfo.shn (0 si no es equipable).
func (s *WorldDataStore) GetItemEquipSlot(itemID uint16) byte {
	if s == nil {
		return 0
	}
	return s.itemEquips[itemID]
}

// GetStyleCouponID retorna el ItemID oficial del cupón de belleza correspondiente al grado (1..5).
func (s *WorldDataStore) GetStyleCouponID(grade byte) uint16 {
	if s == nil || grade < 1 || int(grade) >= len(s.styleCoupons) {
		return 0
	}
	return s.styleCoupons[grade]
}

// GetAttendReward retorna la definición oficial de recompensa de asistencia según AttendReward.shn.
func (s *WorldDataStore) GetAttendReward(index byte) (AttendRewardDef, bool) {
	if s == nil {
		return AttendRewardDef{}, false
	}
	rew, ok := s.attendRewards[index]
	return rew, ok
}

// GetClassName retorna el nombre oficial en inglés de la clase según ClassName.shn.
func (s *WorldDataStore) GetClassName(classID byte) string {
	if s == nil {
		return ""
	}
	return s.classNames[classID]
}
