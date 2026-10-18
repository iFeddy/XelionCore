--------------------------------------------------------------------------------
--                   Halloween Bat Hunt Kingdom Quest Scenario                --
--------------------------------------------------------------------------------

require("common")

KQHBat = {}
KQHBat.State = {
    Maul = "KQ_InvincibleHammer",
    Winners = {},
}

function KQHBat.Main(Field)
    local state = KQHBat.State
    StopBattle(Field, "PK")

    local players = WaitPlayerLogin(Field, 60)
    if not players or players == 0 then
        KQHBat.QuestEnd(Field)
        return
    end

    -- Entregar martillos invencibles a los participantes
    local playerList = GetFieldPlayers(Field)
    for _, player in ipairs(playerList) do
        CreateItem(player, state.Maul, 1)
    end

    BroadcastToField(Field, "¡La Cacería de Murciélagos de Halloween ha comenzado!")
    SetTimeLimit(Field, 10 * 60)
    SetInterrupt(Field, "TimeOut", "KQHBat.OnHuntEnd")

    -- Oleadas de murciélagos
    for round = 1, 5 do
        SpawnRegenGroup(Field, "KQ_HBat_0" .. round, "HBat")
        Wait(60)
    end
end

function KQHBat.OnHuntEnd(Field)
    local state = KQHBat.State
    local resultList = GetKQResultList(Field, "HBat")
    for i, p in ipairs(resultList) do
        if i <= 3 then
            SendQuestResult(Field, p, true)
            RewardIndividual(p, "HBAT_WINNER_" .. i)
            state.Winners[i] = GetCharacterName(p)
        else
            SendQuestResult(Field, p, false)
            RewardIndividual(p, "HBAT_PARTICIPANT")
        end
    end

    BroadcastToField(Field, "¡Cacería finalizada! Ganadores: " .. (state.Winners[1] or "Nadie"))
    Wait(15)
    TeleportAll(Field, "Rou", 6208, 6405)
    EndScenario(Field)
end

function KQHBat.QuestEnd(Field)
    TeleportAll(Field, "Rou", 6208, 6405)
    EndScenario(Field)
end

return KQHBat
