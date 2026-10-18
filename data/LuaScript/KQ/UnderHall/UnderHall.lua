--------------------------------------------------------------------------------
--                     UnderHall Kingdom Quest Scenario                       --
--------------------------------------------------------------------------------

require("common")

KQUnderHall = {}
KQUnderHall.State = {
    BossRoboHandle = 0,
}

function KQUnderHall.Main(Field)
    local state = KQUnderHall.State
    SetScriptFile(Field, "KQUnderHall")

    local players = WaitPlayerLogin(Field, 60)
    if not players or players == 0 then
        KQUnderHall.QuestFail(Field)
        return
    end

    Wait(20)
    SetTimeLimit(Field, 30 * 60)
    SetInterrupt(Field, "TimeOut", "KQUnderHall.QuestFail")

    -- 5 fases del UnderHall
    for floor = 1, 4 do
        SpawnRegenGroup(Field, "KQ_UnderHall_0" .. floor, "UnderHall")
        SetInterrupt(Field, "MobEliminate", "KQUnderHall.OnFloorClear" .. floor)
        WaitSignal("Floor" .. floor .. "Done")
    end

    -- Fase 5: Jefe final Boss Robo
    state.BossRoboHandle = SpawnMob(Field, "KQ_BossRobo", 5000, 5000, 0, 1000, "Normal")
    SetInterrupt(Field, "DeadHandle", state.BossRoboHandle, "KQUnderHall.QuestSuccess")
    WaitSignal("KQDone")
end

for floor = 1, 4 do
    KQUnderHall["OnFloorClear" .. floor] = function(Field)
        BroadcastToField(Field, "¡Piso " .. floor .. " completado! Avanzando al siguiente nivel...")
        SendSignal("Floor" .. floor .. "Done")
    end
end

function KQUnderHall.QuestSuccess(Field)
    BroadcastToField(Field, "¡El Jefe Mecánico Boss Robo ha sido destruido! ¡KQ Completada!")
    RewardKQ(Field, "KQUnderHall", true)
    Wait(15)
    TeleportAll(Field, "Urg", 7757, 6776)
    EndScenario(Field)
    SendSignal("KQDone")
end

function KQUnderHall.QuestFail(Field)
    BroadcastToField(Field, "¡Misión fallida en el UnderHall!")
    RewardKQ(Field, "KQUnderHall", false)
    Wait(15)
    TeleportAll(Field, "Urg", 7757, 6776)
    EndScenario(Field)
    SendSignal("KQDone")
end

return KQUnderHall
