package world

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/iFeddy/XelionCore/pkg/data/shn"
)

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
	Shortcuts     []StarterShortcut `json:"shortcuts"`
	StarterSkills []int             `json:"starter_skills"`
}

// DefaultCharacterDataStore almacena las configuraciones de todas las clases disponibles.
type DefaultCharacterDataStore struct {
	Classes map[string]DefaultClassData `json:"classes"`
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
	key := fmt.Sprintf("%d", classID)
	data, exists := s.Classes[key]
	return data, exists
}

// BadNameFilter gestiona las palabras y nombres restringidos por el motor oficial.
type BadNameFilter struct {
	exactBadNames map[string]bool
	subStrings    []string
}

// LoadBadNameFilter lee BadNameFilter.shn oficial con el parser binario de XelionCore.
func LoadBadNameFilter(filePath string) (*BadNameFilter, error) {
	file, err := shn.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load %s: %w", filePath, err)
	}

	filter := &BadNameFilter{
		exactBadNames: make(map[string]bool, file.RecordCount),
		subStrings:    []string{"admin", "gm", "master", "support", "fiesta", "gamigo", "onson"},
	}

	for _, row := range file.Rows {
		name := strings.TrimSpace(row.GetString("BadName"))
		if name != "" {
			filter.exactBadNames[strings.ToLower(name)] = true
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

	// Coincidencia exacta en lista negra oficial
	if f.exactBadNames[lower] {
		return true
	}

	// Nombres reservados de sistema / staff
	for _, sub := range f.subStrings {
		if strings.Contains(lower, sub) {
			return true
		}
	}

	return false
}

// WorldDataStore administra todos los datos estáticos, tablas SHN y validadores oficiales de WorldService.
type WorldDataStore struct {
	DataDir      string
	BadNames     *BadNameFilter
	DefaultData  *DefaultCharacterDataStore
	TutorialData *DefaultCharacterDataStore
	Checksums    *shn.ChecksumEngine

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

	// LoadedFiles contiene todos los SHNs cargados para acceso por clave
	LoadedFiles map[string]*shn.File

	validFaces  map[byte]bool
	validHairs  map[byte]bool
	validColors map[byte]bool
}

// LoadWorldDataStore carga e inicializa la totalidad de los archivos SHN y datos estáticos
// oficiales requeridos por WorldManager y Character (22 de coordinación + 4 de Character + Checksums).
func LoadWorldDataStore(dataDir string) (*WorldDataStore, error) {
	store := &WorldDataStore{
		DataDir:     dataDir,
		LoadedFiles: make(map[string]*shn.File),
		validFaces:  make(map[byte]bool),
		validHairs:  make(map[byte]bool),
		validColors: make(map[byte]bool),
	}

	// 1. Cargar BadNameFilter
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

	// 4. Inicializar motor de Checksums de 49 archivos
	if chkEngine, err := shn.NewChecksumEngine(dataDir); err == nil {
		store.Checksums = chkEngine
	} else {
		log.Printf("[WorldDataStore] Warning: Failed to init ChecksumEngine: %v", err)
	}

	// 5. Cargar los 4 SHN de Character.exe
	loadSHN := func(name string, target **shn.File) {
		p := filepath.Join(dataDir, name)
		f, err := shn.ReadFile(p)
		if err != nil {
			log.Printf("[WorldDataStore] Warning: Could not load %s: %v", name, err)
			return
		}
		*target = f
		store.LoadedFiles[name] = f
	}

	loadSHN("FaceInfo.shn", &store.FaceInfo)
	loadSHN("HairInfo.shn", &store.HairInfo)
	loadSHN("HairColorInfo.shn", &store.HairColorInfo)
	loadSHN("HolyPromiseReward.shn", &store.HolyPromiseReward)

	// 6. Cargar los 22 SHN de coordinación de WorldManager.exe
	loadSHN("GuildAcademy.shn", &store.GuildAcademy)
	loadSHN("GuildAcademyLevelUp.shn", &store.GuildAcademyLevelUp)
	loadSHN("GuildAcademyRank.shn", &store.GuildAcademyRank)
	loadSHN("GuildGradeData.shn", &store.GuildGradeData)
	loadSHN("GuildGradeScoreData.shn", &store.GuildGradeScoreData)
	loadSHN("GuildLevelScoreData.shn", &store.GuildLevelScoreData)
	loadSHN("GuildTournament.shn", &store.GuildTournament)
	loadSHN("GuildTournamentRequire.shn", &store.GuildTournamentRequire)
	loadSHN("GuildTournamentReward.shn", &store.GuildTournamentReward)
	loadSHN("GBHouse.shn", &store.GBHouse)
	loadSHN("KingdomQuest.shn", &store.KingdomQuest)
	loadSHN("KingdomQuestMap.shn", &store.KingdomQuestMap)
	loadSHN("KQTeam.shn", &store.KQTeam)
	loadSHN("KQIsVote.shn", &store.KQIsVote)
	loadSHN("KQVoteMajorityRate.shn", &store.KQVoteMajorityRate)
	loadSHN("SpamerReport.shn", &store.SpamerReport)
	loadSHN("AttendReward.shn", &store.AttendReward)
	loadSHN("AttendSchedule.shn", &store.AttendSchedule)
	loadSHN("NpcSchedule.shn", &store.NpcSchedule)
	loadSHN("SingleData.shn", &store.SingleData)
	loadSHN("UseClassTypeInfo.shn", &store.UseClassTypeInfo)
	loadSHN("MsgWorldManager.shn", &store.MsgWorldManager)

	// BadNameFilter.shn también se registra en LoadedFiles
	if store.LoadedFiles["BadNameFilter.shn"] == nil {
		if f, err := shn.ReadFile(filepath.Join(dataDir, "BadNameFilter.shn")); err == nil {
			store.LoadedFiles["BadNameFilter.shn"] = f
		}
	}

	// Popular mapas de validación rápida a partir de las tablas SHN oficiales
	// Exclusivamente Grade == 0 (opciones base iniciales gratuitas).
	// Grade 1 (Beauty Shop), Grade 4 (Premium Coupon) y Grade 5 (Deluxe Coupon)
	// quedan estrictamente bloqueados en la creación de personajes para evitar exploits por edición de SHN del cliente.
	if store.FaceInfo != nil {
		for _, row := range store.FaceInfo.Rows {
			if row.GetInt("Grade") == 0 {
				store.validFaces[byte(row.GetInt("ID"))] = true
			}
		}
	}
	if store.HairInfo != nil {
		for _, row := range store.HairInfo.Rows {
			if row.GetInt("Grade") == 0 {
				store.validHairs[byte(row.GetInt("ID"))] = true
			}
		}
	}
	if store.HairColorInfo != nil {
		for _, row := range store.HairColorInfo.Rows {
			if row.GetInt("Grade") == 0 {
				store.validColors[byte(row.GetInt("ID"))] = true
			}
		}
	}

	return store, nil
}

// ValidateAppearance valida formalmente que la cara, peinado y color solicitados
// para la creación del personaje existan en las tablas oficiales FaceInfo, HairInfo y HairColorInfo.
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
