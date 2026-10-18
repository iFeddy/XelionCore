--------------------------------------------------------------------------------
--                         Honeying Kingdom Quest Scenario                    --
--------------------------------------------------------------------------------

require("common")

KQHoneying = {}
KQHoneying.State = {
    Doors = {},
    EffDoors = {},
    BossHandle = 0,
}

function KQHoneying.Main(Field)
    local state = KQHoneying.State

    -- Construir puertas entre pisos
    state.Doors[1] = BuildDoor(Field, "KQ_SlimeGate", 9860, 6094, 272, 1000, "Normal")
    state.Doors[2] = BuildDoor(Field, "KQ_SlimeGate", 6692, 3944, 6, 1000, "Normal")
    state.Doors[3] = BuildDoor(Field, "KQ_SlimeGate", 5894, 6098, 88, 1000, "Normal")

    for i = 1, 3 do
        CloseDoor(Field, state.Doors[i], "CloseGate0" .. i)
        state.EffDoors[i] = SpawnEffectOnObject(state.Doors[i], "KQ_SlimeGate", 3600000, 1000)
    end

    SetScriptFile(Field, "KQHoneying")

    -- Esperar a los jugadores
    local players = WaitPlayerLogin(Field, 60)
    if not players or players == 0 then
        KQHoneying.QuestFail(Field)
        return
    end

    Wait(5)
    SetTimeLimit(Field, 35 * 60)
    SetInterrupt(Field, "TimeOut", "KQHoneying.QuestFail")

    -- Fases por piso
    KQHoneying.FirstFloor(Field)
    KQHoneying.SecondFloor(Field)
    KQHoneying.ThirdFloor(Field)
    KQHoneying.TopFloor(Field)
end

function KQHoneying.FirstFloor(Field)
    local state = KQHoneying.State
    SpawnRegenGroup(Field, "KQ_Honey_01", "Honeying")
    SetInterrupt(Field, "MobEliminate", "KQHoneying.OnFloor1Clear")
    WaitSignal("Floor1Cleared")
end

function KQHoneying.OnFloor1Clear(Field)
    local state = KQHoneying.State
    OpenDoor(Field, state.Doors[1], "CloseGate01")
    Vanish(Field, state.EffDoors[1])
    BroadcastToField(Field, "¡La puerta del primer piso se ha abierto!")
    SendSignal("Floor1Cleared")
end

function KQHoneying.SecondFloor(Field)
    local state = KQHoneying.State
    SpawnRegenGroup(Field, "KQ_Honey_02", "Honeying")
    SetInterrupt(Field, "MobEliminate", "KQHoneying.OnFloor2Clear")
    WaitSignal("Floor2Cleared")
end

function KQHoneying.OnFloor2Clear(Field)
    local state = KQHoneying.State
    OpenDoor(Field, state.Doors[2], "CloseGate02")
    Vanish(Field, state.EffDoors[2])
    BroadcastToField(Field, "¡La puerta del segundo piso se ha abierto!")
    SendSignal("Floor2Cleared")
end

function KQHoneying.ThirdFloor(Field)
    local state = KQHoneying.State
    SpawnRegenGroup(Field, "KQ_Honey_03", "Honeying")
    SetInterrupt(Field, "MobEliminate", "KQHoneying.OnFloor3Clear")
    WaitSignal("Floor3Cleared")
end

function KQHoneying.OnFloor3Clear(Field)
    local state = KQHoneying.State
    OpenDoor(Field, state.Doors[3], "CloseGate03")
    Vanish(Field, state.EffDoors[3])
    BroadcastToField(Field, "¡La puerta final hacia la cima se ha abierto!")
    SendSignal("Floor3Cleared")
end

function KQHoneying.TopFloor(Field)
    local state = KQHoneying.State
    SpawnRegenGroup(Field, "KQ_Honey_04", "Honeying")
    state.BossHandle = SpawnMob(Field, "KQ_KingHoneying", 5000, 5000, 0, 1000, "Normal")
    SetInterrupt(Field, "DeadHandle", state.BossHandle, "KQHoneying.QuestSuccess")
    WaitSignal("KQDone")
end

function KQHoneying.QuestSuccess(Field)
    BroadcastToField(Field, "¡Han derrotado al Rey Honeying! ¡Kingdom Quest Completada!")
    RewardKQ(Field, "KQHoneying", true)
    Wait(15)
    TeleportAll(Field, "Rou", 6208, 6405)
    EndScenario(Field)
    SendSignal("KQDone")
end

function KQHoneying.QuestFail(Field)
    BroadcastToField(Field, "¡Tiempo agotado o misión fallida!")
    RewardKQ(Field, "KQHoneying", false)
    Wait(15)
    TeleportAll(Field, "Rou", 6208, 6405)
    EndScenario(Field)
    SendSignal("KQDone")
end

return KQHoneying
