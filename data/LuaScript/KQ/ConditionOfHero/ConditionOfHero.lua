--------------------------------------------------------------------------------
--                   Condition of Hero (PvP Kingdom Quest)                    --
--------------------------------------------------------------------------------

require("common")

KQConditionOfHero = {}
KQConditionOfHero.State = {
    InitFlag = 1,
    Winners = {},
}

function KQConditionOfHero.Main(Field, InitFlag)
    local state = KQConditionOfHero.State
    state.InitFlag = InitFlag or 1

    -- Detener combate inicial y apagar FreePK
    StopBattle(Field, "PK")
    SetFreePK(Field, false)

    SetScriptFile(Field, "KQConditionOfHero")

    local players = WaitPlayerLogin(Field, 60)
    if not players or players == 0 then
        KQConditionOfHero.QuestFail(Field)
        return
    end

    Wait(7)

    -- Introducción del Alcalde Roumenus
    for i = 0, 4 do
        SendChatWindow(Field, "RouTownChiefRoumenus", "Intro" .. i)
        Wait(2)
    end

    -- Fase 1: Recolección y preparación de monstruos
    SetTimeLimit(Field, 3 * 60)
    for count = 0, 19 do
        SpawnRegenGroup(Field, "KDHero", "HERO" .. state.InitFlag)
    end
    WaitSignal("PrepDone")

    -- Fase 2: Inicio de Duelo PK
    SendChatWindow(Field, "RouTownChiefRoumenus", "DualStart")
    StartBattle(Field, "PK")
    SetFreePK(Field, true)

    SetTimeLimit(Field, 7 * 60)
    SetInterrupt(Field, "TimeOut", "KQConditionOfHero.OnPvPEnd")
    SetTimer(Field, "KQConditionOfHero.MineRegen", 30, true)

    WaitSignal("PvPFinished")
end

function KQConditionOfHero.MineRegen(Field)
    local state = KQConditionOfHero.State
    SpawnRegenGroup(Field, "KDHero", "HERO" .. state.InitFlag)
end

function KQConditionOfHero.OnPvPEnd(Field)
    local state = KQConditionOfHero.State
    SendChatWindow(Field, "RouTownChiefRoumenus", "DualStop")
    StopBattle(Field, "PK")
    SetFreePK(Field, false)

    -- Evaluar resultados y posiciones
    local resultList = GetKQResultList(Field, "ConditionOfHero")
    local count = 0
    for _, playerHandle in ipairs(resultList) do
        if count < 3 then
            SendQuestResult(Field, playerHandle, true)
            RewardIndividual(playerHandle, "HERO_" .. state.InitFlag .. "_" .. count)
            state.Winners[count + 1] = GetCharacterName(playerHandle)
        else
            SendQuestResult(Field, playerHandle, false)
            RewardIndividual(playerHandle, "HERO_" .. state.InitFlag .. "_3")
        end
        count = count + 1
    end

    SendChatWindow(Field, "RouTownChiefRoumenus", "DualResult", state.Winners[1] or "", state.Winners[2] or "", state.Winners[3] or "")

    -- Retorno a Elderine
    local returnTimes = {30, 20, 10, 5}
    for _, sec in ipairs(returnTimes) do
        BroadcastToField(Field, "KQReturn" .. sec)
        Wait(sec == 5 and 5 or (sec == 10 and 5 or 10))
    end

    TeleportAll(Field, "Eld", 17214, 13445)
    EndScenario(Field)
    SendSignal("PvPFinished")
end

function KQConditionOfHero.QuestFail(Field)
    StopBattle(Field, "PK")
    SetFreePK(Field, false)

    local returnTimes = {30, 20, 10, 5}
    for _, sec in ipairs(returnTimes) do
        BroadcastToField(Field, "KQReturn" .. sec)
        Wait(sec == 5 and 5 or (sec == 10 and 5 or 10))
    end

    TeleportAll(Field, "Eld", 17214, 13445)
    EndScenario(Field)
    SendSignal("PvPFinished")
end

return KQConditionOfHero
