package shn

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ChecksumFileList contiene la lista oficial exacta y ordenada de los 49 archivos .shn
// extraída de la tabla consecutiva de WorldManager.exe (índices 0 a 48).
var ChecksumFileList = [49]string{
	"ChargedDeletableBuff.shn",
	"MinimonAutoUseItem.shn",
	"MinimonInfo.shn",
	"TermExtendMatch.shn",
	"ChatColor.shn",
	"PupMain.shn",
	"ItemMoney.shn",
	"MarketSearchInfo.shn",
	"SingleData.shn",
	"ItemViewEquipTypeInfo.shn",
	"GTIView.shn",
	"CollectCardView.shn",
	"ProduceView.shn",
	"PassiveSkillView.shn",
	"NPCViewInfo.shn",
	"MobViewInfo.shn",
	"MapViewInfo.shn",
	"ItemViewInfo.shn",
	"ItemShopView.shn",
	"EffectViewInfo.shn",
	"CharacterTitleStateView.shn",
	"ActiveSkillView.shn",
	"AbStateView.shn",
	"MapWayPoint.shn",
	"MapLinkPoint.shn",
	"ActionViewInfo.shn",
	"DiceDividind.shn",
	"MiniHouseEndure.shn",
	"MiniHouseFurnitureObjEffect.shn",
	"WeaponTitleData.shn",
	"WeaponAttrib.shn",
	"UpgradeInfo.shn",
	"SubAbstate.shn",
	"Riding.shn",
	"PassiveSkill.shn",
	"MobInfo.shn",
	"MiniHouseObjAni.shn",
	"MiniHouseFurniture.shn",
	"MiniHouse.shn",
	"MapInfo.shn",
	"ItemInfo.shn",
	"ItemDismantle.shn",
	"GradeItemOption.shn",
	"Gather.shn",
	"ClassName.shn",
	"ChargedEffect.shn",
	"CharacterTitleData.shn",
	"ActiveSkill.shn", // Índice 47
	"Abstate.shn",     // Índice 48
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

// NewChecksumEngine calcula los checksums MD5 oficiales de los 49 archivos .shn
// buscando en el directorio provisto (ej. "data/shine").
func NewChecksumEngine(shineDir string) (*ChecksumEngine, error) {
	engine := &ChecksumEngine{}

	for idx, filename := range ChecksumFileList {
		path := filepath.Join(shineDir, filename)
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read anti-cheat SHN %s (index %d): %w", filename, idx, err)
		}

		// Calcular MD5 canónico del archivo en disco
		hash := md5.Sum(data)
		hexStr := strings.ToLower(hex.EncodeToString(hash[:]))

		var name8 Name8
		copy(name8[:], hexStr)

		engine.checksums[idx] = name8
		engine.rawHashes[idx] = hexStr
	}

	return engine, nil
}

// VerifyDataFiles implementa CShnDataFileCheckSum::DataFileCheck
// comparando el array Name8[49] provisto por el cliente contra los hashes oficiales del servidor.
// Si son idénticos, retorna (true, 0).
// Si alguno difiere, retorna (false, nWrongDataFileIndex) donde nWrongDataFileIndex es el índice del archivo adulterado.
func (e *ChecksumEngine) VerifyDataFiles(clientChecksums [49]Name8) (bool, int) {
	for i := 0; i < 49; i++ {
		if clientChecksums[i] != e.checksums[i] {
			return false, i
		}
	}
	return true, 0
}

// GetChecksums retorna el array oficial Name8[49] para tests o inspección.
func (e *ChecksumEngine) GetChecksums() [49]Name8 {
	return e.checksums
}

// GetFileName retorna el nombre de archivo correspondiente a un índice (0..48).
func (e *ChecksumEngine) GetFileName(index int) string {
	if index < 0 || index >= 49 {
		return ""
	}
	return ChecksumFileList[index]
}
