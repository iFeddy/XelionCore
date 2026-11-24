package shn

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ChecksumFileList contiene la lista oficial exacta y ordenada de los 49 archivos .shn
// extraída directamente del enum oficial SHN_DATA_FILE_INDEX (Zone.pdb / WorldManager.pdb, índices 0 a 48).
var ChecksumFileList = [49]string{
	"Abstate.shn",                     // 0: SHN_Abstate
	"ActiveSkill.shn",                 // 1: SHN_ActiveSkill
	"CharacterTitleData.shn",          // 2: SHN_CharacterTitleData
	"ChargedEffect.shn",               // 3: SHN_ChargedEffect
	"ClassName.shn",                   // 4: SHN_ClassName
	"Gather.shn",                      // 5: SHN_Gather
	"GradeItemOption.shn",             // 6: SHN_GradeItemOption
	"ItemDismantle.shn",               // 7: SHN_ItemDismantle
	"ItemInfo.shn",                    // 8: SHN_ItemInfo
	"MapInfo.shn",                     // 9: SHN_MapInfo
	"MiniHouse.shn",                   // 10: SHN_MiniHouse
	"MiniHouseFurniture.shn",          // 11: SHN_MiniHouseFurniture
	"MiniHouseObjAni.shn",             // 12: SHN_MiniHouseObjAni
	"MobInfo.shn",                     // 13: SHN_MobInfo
	"PassiveSkill.shn",                // 14: SHN_PassiveSkill
	"Riding.shn",                      // 15: SHN_Riding
	"SubAbstate.shn",                  // 16: SHN_SubAbstate
	"UpgradeInfo.shn",                 // 17: SHN_UpgradeInfo
	"WeaponAttrib.shn",                // 18: SHN_WeaponAttrib
	"WeaponTitleData.shn",             // 19: SHN_WeaponTitleData
	"MiniHouseFurnitureObjEffect.shn", // 20: SHN_MiniHouseFurnitureObjEffect
	"MiniHouseEndure.shn",             // 21: SHN_MiniHouseFurnitureEndure
	"DiceDividind.shn",                // 22: SHN_DiceDividind
	"ActionViewInfo.shn",              // 23: SHN_ActionViewInfo
	"MapLinkPoint.shn",                // 24: SHN_MapLinkPoint
	"MapWayPoint.shn",                 // 25: SHN_MapWayPoint
	"AbStateView.shn",                 // 26: SHN_AbStateView
	"ActiveSkillView.shn",             // 27: SHN_ActiveSkillView
	"CharacterTitleStateView.shn",     // 28: SHN_CharacterTitleStateView
	"EffectViewInfo.shn",              // 29: SHN_EffectViewInfo
	"ItemShopView.shn",                // 30: SHN_ItemShopView
	"ItemViewInfo.shn",                // 31: SHN_ItemViewInfo
	"MapViewInfo.shn",                 // 32: SHN_MapViewInfo
	"MobViewInfo.shn",                 // 33: SHN_MobViewInfo
	"NPCViewInfo.shn",                 // 34: SHN_NPCViewInfo
	"PassiveSkillView.shn",            // 35: SHN_PassiveSkillView
	"ProduceView.shn",                 // 36: SHN_ProduceView
	"CollectCardView.shn",             // 37: SHN_CollectCardView
	"GTIView.shn",                     // 38: SHN_GTIView
	"ItemViewEquipTypeInfo.shn",       // 39: SHN_ItemViewEquipTypeInfo
	"SingleData.shn",                  // 40: SHN_SingleData
	"MarketSearchInfo.shn",            // 41: SHN_MarketSearchInfo
	"ItemMoney.shn",                   // 42: SHN_ItemMoney
	"PupMain.shn",                     // 43: SHN_PupMain
	"ChatColor.shn",                   // 44: SHN_ChatColor
	"TermExtendMatch.shn",             // 45: SHN_TermExtendMatch
	"MinimonInfo.shn",                 // 46: SHN_MinimonInfo
	"MinimonAutoUseItem.shn",          // 47: SHN_MinimonAutoUseItem
	"ChargedDeletableBuff.shn",        // 48: SHN_ChargedDeletableBuff
}

// Name8 representa la estructura canónica oficial de 32 bytes (8 DWORDs)
// correspondiente al hash MD5 hexadecimal que el cliente envía en PROTO_NC_MAP_LOGIN_REQ.
type Name8 [32]byte

// String convierte el hash Name8 en string legible.
func (n Name8) String() string {
	return string(n[:])
}

// ChecksumEngine implementa CShnDataFileCheckSum para verificación anti-cheat.
type ChecksumEngine struct {
	checksums [49]Name8
	rawHashes [49]string
}

// findFileCaseInsensitive busca un archivo en los directorios indicados ignorando mayúsculas/minúsculas.
func findFileCaseInsensitive(filename string, dirs []string) (string, error) {
	lowerName := strings.ToLower(filename)
	for _, dir := range dirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if strings.ToLower(e.Name()) == lowerName {
				return filepath.Join(dir, e.Name()), nil
			}
		}
	}
	return "", os.ErrNotExist
}

func expandChecksumDirs(dirs []string) []string {
	if len(dirs) == 0 {
		dirs = []string{"data", "data/View"}
	}
	allDirs := make([]string, 0, len(dirs)*2)
	for _, d := range dirs {
		allDirs = append(allDirs, d)
		viewSub := filepath.Join(d, "View")
		if fi, err := os.Stat(viewSub); err == nil && fi.IsDir() {
			allDirs = append(allDirs, viewSub)
		}
	}
	return allDirs
}

func computeSHNFileMD5FromDisk(filename string, idx int, allDirs []string) (Name8, string, error) {
	path, err := findFileCaseInsensitive(filename, allDirs)
	if err != nil {
		return Name8{}, "", fmt.Errorf("failed to locate anti-cheat SHN %s (index %d) in dirs %v: %w", filename, idx, allDirs, err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return Name8{}, "", fmt.Errorf("failed to read anti-cheat SHN %s (index %d): %w", filename, idx, err)
	}

	// En el cliente oficial (CDataReader::Read en Zone.exe y Fiesta.bin):
	// Si dataMode == 1 (offset 28:32), el cuerpo a partir de offset 36 se desencripta in-place con CDataReader::Encription.
	// El MD5 se calcula sobre el buffer con la cabecera (36 bytes) + cuerpo desencriptado.
	if len(data) >= 36 {
		dataMode := binary.LittleEndian.Uint32(data[28:32])
		if dataMode == 1 {
			DecryptInPlace(data[36:])
		}
	}

	hash := md5.Sum(data)
	hexStr := hex.EncodeToString(hash[:])

	var name8 Name8
	copy(name8[:], hexStr)
	return name8, hexStr, nil
}

// NewChecksumEngineFromFiles inicializa CShnDataFileCheckSum reutilizando los hashes MD5
// ya calculados sobre los buffers desencriptados de files (0 re-lecturas de disco para tablas ya cargadas).
// Para cualquier archivo ausente en files, realiza fallback a fallbackDirs.
func NewChecksumEngineFromFiles(files map[string]*File, fallbackDirs ...string) (*ChecksumEngine, error) {
	byLower := make(map[string]*File, len(files)*2)
	for k, f := range files {
		if f == nil {
			continue
		}
		lk := strings.ToLower(filepath.Base(k))
		byLower[lk] = f
		base := strings.TrimSuffix(lk, filepath.Ext(lk))
		byLower[base] = f
	}

	var allDirs []string
	engine := &ChecksumEngine{}

	for idx, filename := range ChecksumFileList {
		lowerFull := strings.ToLower(filename)
		lowerBase := strings.TrimSuffix(lowerFull, ".shn")

		f := byLower[lowerFull]
		if f == nil {
			f = byLower[lowerBase]
		}
		if f != nil {
			if name8, ok := f.MD5Checksum(); ok {
				engine.checksums[idx] = name8
				engine.rawHashes[idx] = name8.String()
				continue
			}
		}

		if allDirs == nil {
			allDirs = expandChecksumDirs(fallbackDirs)
		}
		name8, hexStr, err := computeSHNFileMD5FromDisk(filename, idx, allDirs)
		if err != nil {
			return nil, err
		}
		engine.checksums[idx] = name8
		engine.rawHashes[idx] = hexStr
	}

	return engine, nil
}

// NewChecksumEngine calcula los checksums MD5 oficiales de los 49 archivos .shn
// buscando en los directorios provistos (ej. "data", "data/View").
func NewChecksumEngine(dirs ...string) (*ChecksumEngine, error) {
	allDirs := expandChecksumDirs(dirs)
	engine := &ChecksumEngine{}

	for idx, filename := range ChecksumFileList {
		name8, hexStr, err := computeSHNFileMD5FromDisk(filename, idx, allDirs)
		if err != nil {
			return nil, err
		}
		engine.checksums[idx] = name8
		engine.rawHashes[idx] = hexStr
	}

	return engine, nil
}

// VerifyDataFiles implementa CShnDataFileCheckSum::DataFileCheck
// comparando el array Name8[49] provisto por el cliente contra los hashes oficiales del servidor.
// Si son idénticos, retorna (true, 0).
// Si alguno difiere, retorna (false, nWrongDataFileIndex) donde nWrongDataFileIndex es el índice (0..48) del archivo adulterado.
func (e *ChecksumEngine) VerifyDataFiles(clientChecksums [49]Name8) (bool, int) {
	for i := 0; i < 49; i++ {
		// Comparación insensible a mayúsculas/minúsculas de los 32 caracteres ASCII
		if !strings.EqualFold(string(clientChecksums[i][:]), string(e.checksums[i][:])) {
			return false, i
		}
	}
	return true, 0
}

// GetChecksums retorna el array oficial Name8[49] para tests o inspección.
func (e *ChecksumEngine) GetChecksums() [49]Name8 {
	return e.checksums
}

// GetRawHash retorna el string hexadecimal de un índice.
func (e *ChecksumEngine) GetRawHash(index int) string {
	if index < 0 || index >= 49 {
		return ""
	}
	return e.rawHashes[index]
}

// GetFileName retorna el nombre de archivo correspondiente a un índice (0..48).
func (e *ChecksumEngine) GetFileName(index int) string {
	if index < 0 || index >= 49 {
		return ""
	}
	return ChecksumFileList[index]
}
