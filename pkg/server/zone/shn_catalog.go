package zone

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/iFeddy/XelionCore/pkg/data/shn"
)

// SHNTableRef define un archivo .shn oficial de Zone.exe y sus ubicaciones en data/.
type SHNTableRef struct {
	Index    int
	FileName string
	Paths    []string
}

// RequiredSHNTables contiene las 180 tablas .shn oficiales en disco requeridas por Zone.exe (Sección 3.1).
var RequiredSHNTables = []SHNTableRef{
	{Index: 1, FileName: "AbState.shn", Paths: []string{"data/Abstate.shn"}},
	{Index: 2, FileName: "AbStateSaveTypeInfo.shn", Paths: []string{"data/AbStateSaveTypeInfo.shn"}},
	{Index: 3, FileName: "AbStateView.shn", Paths: []string{"data/AbStateView.shn", "data/View/AbStateView.shn"}},
	{Index: 4, FileName: "AccUpGrade.shn", Paths: []string{"data/AccUpgrade.shn"}},
	{Index: 5, FileName: "AccUpGradeInfo.shn", Paths: []string{"data/AccUpGradeInfo.shn"}},
	{Index: 6, FileName: "ActionEffectAbState.shn", Paths: []string{"data/ActionEffectAbState.shn"}},
	{Index: 7, FileName: "ActionEffectItem.shn", Paths: []string{"data/ActionEffectItem.shn"}},
	{Index: 8, FileName: "ActionRangeFactor.shn", Paths: []string{"data/ActionRangeFactor.shn"}},
	{Index: 9, FileName: "ActionViewInfo.shn", Paths: []string{"data/ActionViewInfo.shn", "data/View/ActionViewInfo.shn"}},
	{Index: 10, FileName: "ActiveSkill.shn", Paths: []string{"data/ActiveSkill.shn"}},
	{Index: 11, FileName: "ActiveSkillGroup.shn", Paths: []string{"data/ActiveSkillGroup.shn"}},
	{Index: 12, FileName: "ActiveSkillInfoServer.shn", Paths: []string{"data/ActiveSkillInfoServer.shn"}},
	{Index: 13, FileName: "ActiveSkillView.shn", Paths: []string{"data/ActiveSkillView.shn", "data/View/ActiveSkillView.shn"}},
	{Index: 14, FileName: "AdminLvSet.shn", Paths: []string{"data/AdminLvSet.shn"}},
	{Index: 15, FileName: "AnnounceData.shn", Paths: []string{"data/AnnounceData.shn"}},
	{Index: 16, FileName: "AreaSkill.shn", Paths: []string{"data/AreaSkill.shn"}},
	{Index: 17, FileName: "BadNameFilter.shn", Paths: []string{"data/BadNameFilter.shn"}},
	{Index: 18, FileName: "BelongDice.shn", Paths: []string{"data/BelongDice.shn"}},
	{Index: 19, FileName: "BelongTypeInfo.shn", Paths: []string{"data/BelongTypeInfo.shn"}},
	{Index: 20, FileName: "BMP.shn", Paths: []string{"data/BMP.shn"}},
	{Index: 21, FileName: "BRAccUpGrade.shn", Paths: []string{"data/BRAccUpgrade.shn"}},
	{Index: 22, FileName: "BRAccUpgradeInfo.shn", Paths: []string{"data/BRAccUpgradeInfo.shn"}},
	{Index: 23, FileName: "CharacterTitleData.shn", Paths: []string{"data/CharacterTitleData.shn"}},
	{Index: 24, FileName: "CharacterTitleStateServer.shn", Paths: []string{"data/CharacterTitleStateServer.shn"}},
	{Index: 25, FileName: "CharacterTitleStateView.shn", Paths: []string{"data/CharacterTitleStateView.shn", "data/View/CharacterTitleStateView.shn"}},
	{Index: 26, FileName: "ChargedDeletableBuff.shn", Paths: []string{"data/ChargedDeletableBuff.shn"}},
	{Index: 27, FileName: "ChargedEffect.shn", Paths: []string{"data/ChargedEffect.shn"}},
	{Index: 28, FileName: "ChatColor.shn", Paths: []string{"data/ChatColor.shn"}},
	{Index: 29, FileName: "ClassName.shn", Paths: []string{"data/ClassName.shn"}},
	{Index: 30, FileName: "CollectCard.shn", Paths: []string{"data/CollectCard.shn"}},
	{Index: 31, FileName: "CollectCardDropRate.shn", Paths: []string{"data/CollectCardDropRate.shn"}},
	{Index: 32, FileName: "CollectCardMobGroup.shn", Paths: []string{"data/CollectCardMobGroup.shn"}},
	{Index: 33, FileName: "CollectCardReward.shn", Paths: []string{"data/CollectCardReward.shn"}},
	{Index: 34, FileName: "CollectCardStarRate.shn", Paths: []string{"data/CollectCardStarRate.shn"}},
	{Index: 35, FileName: "CollectCardTitle.shn", Paths: []string{"data/CollectCardTitle.shn"}},
	{Index: 36, FileName: "CollectCardView.shn", Paths: []string{"data/CollectCardView.shn", "data/View/CollectCardView.shn"}},
	{Index: 37, FileName: "DamageLvGapEVP.shn", Paths: []string{"data/DamageLvGapEVP.shn"}},
	{Index: 38, FileName: "DamageLvGapPVE.shn", Paths: []string{"data/DamageLvGapPVE.shn"}},
	{Index: 39, FileName: "DamageLvGapPVP.shn", Paths: []string{"data/DamageLvGapPVP.shn"}},
	{Index: 40, FileName: "DiceDividind.shn", Paths: []string{"data/DiceDividind.shn"}},
	{Index: 41, FileName: "DiceGame.shn", Paths: []string{"data/DiceGame.shn"}},
	{Index: 42, FileName: "DiceRate.shn", Paths: []string{"data/DiceRate.shn"}},
	{Index: 43, FileName: "EffectViewInfo.shn", Paths: []string{"data/EffectViewInfo.shn", "data/View/EffectViewInfo.shn"}},
	{Index: 44, FileName: "EnchantSocketRate.shn", Paths: []string{"data/EnchantSocketRate.shn"}},
	{Index: 45, FileName: "FieldLvCondition.shn", Paths: []string{"data/FieldLvCondition.shn"}},
	{Index: 46, FileName: "FriendPointReward.shn", Paths: []string{"data/FriendPointReward.shn"}},
	{Index: 47, FileName: "Gather.shn", Paths: []string{"data/Gather.shn"}},
	{Index: 48, FileName: "GBBanTime.shn", Paths: []string{"data/GBBanTime.shn"}},
	{Index: 49, FileName: "GBDiceDividind.shn", Paths: []string{"data/GBDiceDividind.shn"}},
	{Index: 50, FileName: "GBDiceGame.shn", Paths: []string{"data/GBDiceGame.shn"}},
	{Index: 51, FileName: "GBDiceRate.shn", Paths: []string{"data/GBDiceRate.shn"}},
	{Index: 52, FileName: "GBEventCode.shn", Paths: []string{"data/GBEventCode.shn"}},
	{Index: 53, FileName: "GBExchangeMaxCoin.shn", Paths: []string{"data/GBExchangeMaxCoin.shn"}},
	{Index: 54, FileName: "GBHouse.shn", Paths: []string{"data/GBHouse.shn"}},
	{Index: 55, FileName: "GBJoinGameMember.shn", Paths: []string{"data/GBJoinGameMember.shn"}},
	{Index: 56, FileName: "GBReward.shn", Paths: []string{"data/GBReward.shn"}},
	{Index: 57, FileName: "GBSMAll.shn", Paths: []string{"data/GBSMAll.shn"}},
	{Index: 58, FileName: "GBSMBetCoin.shn", Paths: []string{"data/GBSMBetCoin.shn"}},
	{Index: 59, FileName: "GBSMCardRate.shn", Paths: []string{"data/GBSMCardRate.shn"}},
	{Index: 60, FileName: "GBSMCenter.shn", Paths: []string{"data/GBSMCenter.shn"}},
	{Index: 61, FileName: "GBSMGroup.shn", Paths: []string{"data/GBSMGroup.shn"}},
	{Index: 62, FileName: "GBSMJPRate.shn", Paths: []string{"data/GBSMJPRate.shn"}},
	{Index: 63, FileName: "GBSMLine.shn", Paths: []string{"data/GBSMLine.shn"}},
	{Index: 64, FileName: "GBSMNPC.shn", Paths: []string{"data/GBSMNPC.shn"}},
	{Index: 65, FileName: "GBTaxRate.shn", Paths: []string{"data/GBTaxRate.shn"}},
	{Index: 66, FileName: "GradeItemOption.shn", Paths: []string{"data/GradeItemOption.shn"}},
	{Index: 67, FileName: "GroupAbState.shn", Paths: []string{"data/GroupAbState.shn"}},
	{Index: 68, FileName: "GTIBreedSubject.shn", Paths: []string{"data/GTIBreedSubject.shn"}},
	{Index: 69, FileName: "GTIGetRate.shn", Paths: []string{"data/GTIGetRate.shn"}},
	{Index: 70, FileName: "GTIGetRateGap.shn", Paths: []string{"data/GTIGetRateGap.shn"}},
	{Index: 71, FileName: "GTIServer.shn", Paths: []string{"data/GTIServer.shn"}},
	{Index: 72, FileName: "GTIView.shn", Paths: []string{"data/GTIView.shn", "data/View/GTIView.shn"}},
	{Index: 73, FileName: "GTWinScore.shn", Paths: []string{"data/GTWinScore.shn"}},
	{Index: 74, FileName: "GuildTournamentLvGap.shn", Paths: []string{"data/GuildTournamentLvGap.shn"}},
	{Index: 75, FileName: "GuildTournamentMasterBuff.shn", Paths: []string{"data/GuildTournamentMasterBuff.shn"}},
	{Index: 76, FileName: "GuildTournamentOccupy.shn", Paths: []string{"data/GuildTournamentOccupy.shn"}},
	{Index: 77, FileName: "GuildTournamentScore.shn", Paths: []string{"data/GuildTournamentScore.shn"}},
	{Index: 78, FileName: "GuildTournamentSkill.shn", Paths: []string{"data/GuildTournamentSkill.shn"}},
	{Index: 79, FileName: "ItemAction.shn", Paths: []string{"data/ItemAction.shn"}},
	{Index: 80, FileName: "ItemActionCondition.shn", Paths: []string{"data/ItemActionCondition.shn"}},
	{Index: 81, FileName: "ItemActionEffect.shn", Paths: []string{"data/ItemActionEffect.shn"}},
	{Index: 82, FileName: "ItemDismantle.shn", Paths: []string{"data/ItemDismantle.shn"}},
	{Index: 83, FileName: "ItemDropLog.shn", Paths: []string{"data/ItemDropLog.shn"}},
	{Index: 84, FileName: "ItemInfo.shn", Paths: []string{"data/ItemInfo.shn"}},
	{Index: 85, FileName: "ItemInfoServer.shn", Paths: []string{"data/ItemInfoServer.shn"}},
	{Index: 86, FileName: "ItemInvenDel.shn", Paths: []string{"data/ItemInvenDel.shn"}},
	{Index: 87, FileName: "ItemMerchantInfo.shn", Paths: []string{"data/ItemMerchantInfo.shn"}},
	{Index: 88, FileName: "ItemMix.shn", Paths: []string{"data/ItemMix.shn"}},
	{Index: 89, FileName: "ItemMoney.shn", Paths: []string{"data/ItemMoney.shn"}},
	{Index: 90, FileName: "ItemOptions.shn", Paths: []string{"data/ItemOptions.shn"}},
	{Index: 91, FileName: "ItemServerEquipTypeInfo.shn", Paths: []string{"data/ItemServerEquipTypeInfo.shn"}},
	{Index: 92, FileName: "ItemShop.shn", Paths: []string{"data/ItemShop.shn"}},
	{Index: 93, FileName: "ItemShopView.shn", Paths: []string{"data/ItemShopView.shn", "data/View/ItemShopView.shn"}},
	{Index: 94, FileName: "ItemSort.shn", Paths: []string{"data/ItemSort.shn"}},
	{Index: 95, FileName: "ItemUpgrade.shn", Paths: []string{"data/ItemUpgrade.shn"}},
	{Index: 96, FileName: "ItemUseEffect.shn", Paths: []string{"data/ItemUseEffect.shn"}},
	{Index: 97, FileName: "ItemViewEquipTypeInfo.shn", Paths: []string{"data/ItemViewEquipTypeInfo.shn", "data/View/ItemViewEquipTypeInfo.shn"}},
	{Index: 98, FileName: "ItemViewInfo.shn", Paths: []string{"data/ItemViewInfo.shn", "data/View/ItemViewInfo.shn"}},
	{Index: 99, FileName: "JobEquipInfo.shn", Paths: []string{"data/JobEquipInfo.shn"}},
	{Index: 100, FileName: "KingdomQuestRew.shn", Paths: []string{"data/KingdomQuestRew.shn"}},
	{Index: 101, FileName: "KQItem.shn", Paths: []string{"data/KQItem.shn"}},
	{Index: 102, FileName: "LCGroupRate.shn", Paths: []string{"data/LCGroupRate.shn"}},
	{Index: 103, FileName: "LCReward.shn", Paths: []string{"data/LCReward.shn"}},
	{Index: 104, FileName: "MapBuff.shn", Paths: []string{"data/MapBuff.shn"}},
	{Index: 105, FileName: "MapInfo.shn", Paths: []string{"data/MapInfo.shn"}},
	{Index: 106, FileName: "MapLinkPoint.shn", Paths: []string{"data/MapLinkPoint.shn"}},
	{Index: 107, FileName: "MapViewInfo.shn", Paths: []string{"data/MapViewInfo.shn", "data/View/MapViewInfo.shn"}},
	{Index: 108, FileName: "MapWayPoint.shn", Paths: []string{"data/MapWayPoint.shn"}},
	{Index: 109, FileName: "MarketSearchInfo.shn", Paths: []string{"data/MarketSearchInfo.shn"}},
	{Index: 110, FileName: "MiniHouse.shn", Paths: []string{"data/MiniHouse.shn"}},
	{Index: 111, FileName: "MiniHouseEndure.shn", Paths: []string{"data/MiniHouseEndure.shn"}},
	{Index: 112, FileName: "MiniHouseFurniture.shn", Paths: []string{"data/MiniHouseFurniture.shn"}},
	{Index: 113, FileName: "MiniHouseFurnitureObjEffect.shn", Paths: []string{"data/MiniHouseFurnitureObjEffect.shn"}},
	{Index: 114, FileName: "MiniHouseObjAni.shn", Paths: []string{"data/MiniHouseObjAni.shn"}},
	{Index: 115, FileName: "MinimonAutoUseItem.shn", Paths: []string{"data/MinimonAutoUseItem.shn"}},
	{Index: 116, FileName: "MinimonInfo.shn", Paths: []string{"data/MinimonInfo.shn"}},
	{Index: 117, FileName: "MobAbStateDropSetting.shn", Paths: []string{"data/MobAbStateDropSetting.shn"}},
	{Index: 118, FileName: "MobAutoAction.shn", Paths: []string{"data/MobAutoAction.shn"}},
	{Index: 119, FileName: "MobConditionServer.shn", Paths: []string{"data/MobConditionServer.shn"}},
	{Index: 120, FileName: "MobInfo.shn", Paths: []string{"data/MobInfo.shn"}},
	{Index: 121, FileName: "MobInfoServer.shn", Paths: []string{"data/MobInfoServer.shn"}},
	{Index: 122, FileName: "MobKillAble.shn", Paths: []string{"data/MobKillAble.shn"}},
	{Index: 123, FileName: "MobKillAnnounce.shn", Paths: []string{"data/MobKillAnnounce.shn"}},
	{Index: 124, FileName: "MobKillLog.shn", Paths: []string{"data/MobKillLog.shn"}},
	{Index: 125, FileName: "MobLifeTime.shn", Paths: []string{"data/MobLifeTime.shn"}},
	{Index: 126, FileName: "MobRegenAni.shn", Paths: []string{"data/MobRegenAni.shn"}},
	{Index: 127, FileName: "MobResist.shn", Paths: []string{"data/MobResist.shn"}},
	{Index: 128, FileName: "MobSpecies.shn", Paths: []string{"data/MobSpecies.shn"}},
	{Index: 129, FileName: "MobViewInfo.shn", Paths: []string{"data/MobViewInfo.shn", "data/View/MobViewInfo.shn"}},
	{Index: 130, FileName: "MobWeapon.shn", Paths: []string{"data/MobWeapon.shn"}},
	{Index: 131, FileName: "MoverAbility.shn", Paths: []string{"data/MoverAbility.shn"}},
	{Index: 132, FileName: "MoverHG.shn", Paths: []string{"data/MoverHG.shn"}},
	{Index: 133, FileName: "MoverItem.shn", Paths: []string{"data/MoverItem.shn"}},
	{Index: 134, FileName: "MoverMain.shn", Paths: []string{"data/MoverMain.shn"}},
	{Index: 135, FileName: "MoverUpgradeEffect.shn", Paths: []string{"data/MoverUpgradeEffect.shn"}},
	{Index: 136, FileName: "MoverUseSkill.shn", Paths: []string{"data/MoverUseSkill.shn"}},
	{Index: 137, FileName: "MultiHitType.shn", Paths: []string{"data/MultiHitType.shn"}},
	{Index: 138, FileName: "MysteryVaultServer.shn", Paths: []string{"data/MysteryVaultServer.shn"}},
	{Index: 139, FileName: "NpcSchedule.shn", Paths: []string{"data/NpcSchedule.shn"}},
	{Index: 140, FileName: "NPCViewInfo.shn", Paths: []string{"data/NPCViewInfo.shn", "data/View/NPCViewInfo.shn"}},
	{Index: 141, FileName: "PartyBonusByLvDiff.shn", Paths: []string{"data/PartyBonusByLvDiff.shn"}},
	{Index: 142, FileName: "PartyBonusByMember.shn", Paths: []string{"data/PartyBonusByMember.shn"}},
	{Index: 143, FileName: "PartyBonusLimit.shn", Paths: []string{"data/PartyBonusLimit.shn"}},
	{Index: 144, FileName: "PassiveSkill.shn", Paths: []string{"data/PassiveSkill.shn"}},
	{Index: 145, FileName: "PassiveSkillView.shn", Paths: []string{"data/PassiveSkillView.shn", "data/View/PassiveSkillView.shn"}},
	{Index: 146, FileName: "Produce.shn", Paths: []string{"data/Produce.shn"}},
	{Index: 147, FileName: "ProduceView.shn", Paths: []string{"data/ProduceView.shn", "data/View/ProduceView.shn"}},
	{Index: 148, FileName: "PSkillSetAbstate.shn", Paths: []string{"data/PSkillSetAbstate.shn"}},
	{Index: 149, FileName: "PupCase.shn", Paths: []string{"data/PupCase.shn"}},
	{Index: 150, FileName: "PupCaseDesc.shn", Paths: []string{"data/PupCaseDesc.shn"}},
	{Index: 151, FileName: "PupFactorCondition.shn", Paths: []string{"data/PupFactorCondition.shn"}},
	{Index: 152, FileName: "PupMain.shn", Paths: []string{"data/PupMain.shn"}},
	{Index: 153, FileName: "PupMind.shn", Paths: []string{"data/PupMind.shn"}},
	{Index: 154, FileName: "PupPriority.shn", Paths: []string{"data/PupPriority.shn"}},
	{Index: 155, FileName: "PupServer.shn", Paths: []string{"data/PupServer.shn"}},
	{Index: 156, FileName: "QuestData.shn", Paths: []string{"data/QuestData.shn"}},
	{Index: 157, FileName: "QuestSpecies.shn", Paths: []string{"data/QuestSpecies.shn"}},
	{Index: 158, FileName: "RandomOption.shn", Paths: []string{"data/RandomOption.shn"}},
	{Index: 159, FileName: "RandomOptionCount.shn", Paths: []string{"data/RandomOptionCount.shn"}},
	{Index: 160, FileName: "RareMoverEachRate.shn", Paths: []string{"data/RareMoverEachRate.shn"}},
	{Index: 161, FileName: "RareMoverRate.shn", Paths: []string{"data/RareMoverRate.shn"}},
	{Index: 162, FileName: "RareMoverSubRate.shn", Paths: []string{"data/RareMoverSubRate.shn"}},
	{Index: 163, FileName: "ReactionType.shn", Paths: []string{"data/ReactionType.shn"}},
	{Index: 164, FileName: "Riding.shn", Paths: []string{"data/Riding.shn"}},
	{Index: 165, FileName: "SetEffect.shn", Paths: []string{"data/SetEffect.shn"}},
	{Index: 166, FileName: "ShineReward.shn", Paths: []string{"data/ShineReward.shn"}},
	{Index: 167, FileName: "SingleData.shn", Paths: []string{"data/SingleData.shn"}},
	{Index: 168, FileName: "SpamerPenalty.shn", Paths: []string{"data/SpamerPenalty.shn"}},
	{Index: 169, FileName: "SpamerPenaltyRule.shn", Paths: []string{"data/SpamerPenaltyRule.shn"}},
	{Index: 170, FileName: "StateField.shn", Paths: []string{"data/StateField.shn"}},
	{Index: 171, FileName: "StateItem.shn", Paths: []string{"data/StateItem.shn"}},
	{Index: 172, FileName: "StateMob.shn", Paths: []string{"data/StateMob.shn"}},
	{Index: 173, FileName: "SubAbState.shn", Paths: []string{"data/SubAbstate.shn"}},
	{Index: 174, FileName: "TermExtendMatch.shn", Paths: []string{"data/TermExtendMatch.shn"}},
	{Index: 175, FileName: "ToggleSkill.shn", Paths: []string{"data/ToggleSkill.shn"}},
	{Index: 176, FileName: "TownPortal.shn", Paths: []string{"data/TownPortal.shn"}},
	{Index: 177, FileName: "UpGradeInfo.shn", Paths: []string{"data/UpgradeInfo.shn"}},
	{Index: 178, FileName: "UseClassTypeInfo.shn", Paths: []string{"data/UseClassTypeInfo.shn"}},
	{Index: 179, FileName: "WeaponAttrib.shn", Paths: []string{"data/WeaponAttrib.shn"}},
	{Index: 180, FileName: "WeaponTitleData.shn", Paths: []string{"data/WeaponTitleData.shn"}},
}

// OptionalSHNTables documenta las 5 referencias .shn presentes en binarios pero ausentes en disco/AstraServer.
var OptionalSHNTables = map[string]string{
	"BMPIndex.shn":     "Opcional: las máscaras .bmp se resuelven directamente por nombre desde data/AreaBMP/",
	"DiecGame.shn":     "Errata tipográfica legacy en Zone.exe reemplazada por DiceGame.shn",
	"MoveMain.shn":     "Referencia legacy reemplazada por MoverMain.shn",
	"MoverUpgrade.shn": "Referencia legacy absorbida por MoverSub.shn y UpgradeInfo.shn",
	"Quest.shn":        "Referencia legacy reemplazada por QuestData.shn",
}

// LoadAllSHNTables carga e indexa en memoria RAM las 180 tablas .shn requeridas en disco.
func LoadAllSHNTables(dataDir string) (map[string]*shn.File, error) {
	tables := make(map[string]*shn.File, len(RequiredSHNTables)*2)

	for _, ref := range RequiredSHNTables {
		var loaded *shn.File
		var err error

		// Intentar cargar por cada ruta declarada
		for _, relPath := range ref.Paths {
			// Normalizar ruta: si relPath empieza con "data/", quitar el prefijo
			p := relPath
			if strings.HasPrefix(p, "data/") {
				p = strings.TrimPrefix(p, "data/")
			}
			fullPath := filepath.Join(dataDir, p)

			if _, statErr := os.Stat(fullPath); statErr == nil {
				loaded, err = shn.ReadFile(fullPath)
				if err == nil {
					break
				}
			}
		}

		// Si falló por las rutas declaradas, buscar directamente en dataDir o dataDir/View
		if loaded == nil {
			directPath := filepath.Join(dataDir, ref.FileName)
			if _, statErr := os.Stat(directPath); statErr == nil {
				loaded, err = shn.ReadFile(directPath)
			}
		}
		if loaded == nil {
			viewPath := filepath.Join(dataDir, "View", ref.FileName)
			if _, statErr := os.Stat(viewPath); statErr == nil {
				loaded, err = shn.ReadFile(viewPath)
			}
		}

		if loaded == nil {
			return nil, fmt.Errorf("error crítico cargando .shn requerido #%d (%s): %v", ref.Index, ref.FileName, err)
		}

		base := strings.TrimSuffix(ref.FileName, filepath.Ext(ref.FileName))
		tables[base] = loaded
		tables[strings.ToLower(base)] = loaded
	}

	// Comprobar las 5 referencias opcionales (documentar en log sin fallar si faltan)
	for optName, desc := range OptionalSHNTables {
		optPath := filepath.Join(dataDir, optName)
		if _, err := os.Stat(optPath); err == nil {
			f, err := shn.ReadFile(optPath)
			if err == nil {
				base := strings.TrimSuffix(optName, filepath.Ext(optName))
				tables[base] = f
				tables[strings.ToLower(base)] = f
			}
		} else {
			log.Printf("[Zone][Data] Nota: SHN opcional '%s' ausente en disco (%s)", optName, desc)
		}
	}

	return tables, nil
}
