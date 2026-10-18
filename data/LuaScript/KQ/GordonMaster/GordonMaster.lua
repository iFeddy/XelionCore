--------------------------------------------------------------------------------
--                      Gordon Master Kingdom Quest Scenario                  --
--------------------------------------------------------------------------------

require("common")

KQGordonMaster = {}
KQGordonMaster.State = {
    DoorHero = 0,
    DoorMidBoss = 0,
    GlassDoor2Boss = 0,
    GordonMasterHandle = 0,
}

function KQGordonMaster.Main(Field)
    local state = KQGordonMaster.State

    -- Construir portones de la mazmorra de Gordon
    state.DoorHero = BuildDoor(Field, "KQ_Gate4", 989, 4298, 0, 1000, "Normal")
    state.DoorMidBoss = BuildDoor(Field, "KQ_Gate4", 2815, 4327, 0, 1000, "Normal")
    state.GlassDoor2Boss = BuildDoor(Field, "KQ_Gate4", 2834, 10023, 0, 1000, "Normal")

    CloseDoor(Field, state.DoorHero, "GorDonDoor_01")
    CloseDoor(Field, state.DoorMidBoss, "GorDonDoor_02")
    CloseDoor(Field, state.GlassDoor2Boss, "GorDonDoor_03")

    SetScriptFile(Field, "KQGordonMaster")

    -- Esperar jugadores
    local players = WaitPlayerLogin(Field, 60)
    if not players or players == 0 then
        KQGordonMaster.QuestFail(Field)
        return
    end

    SetTimeLimit(Field, 45 * 60)
    SetInterrupt(Field, "TimeOut", "KQGordonMaster.QuestFail")

    -- Fases
    KQGordonMaster.Phase1_HeroGate(Field)
    KQGordonMaster.Phase2_MidBoss(Field)
    KQGordonMaster.Phase3_GordonBoss(Field)
end

function KQGordonMaster.Phase1_HeroGate(Field)
    local state = KQGordonMaster.State
    SpawnRegenGroup(Field, "KQ_Gordon_01", "GordonMaster")
    SetInterrupt(Field, "MobEliminate", "KQGordonMaster.OnHeroGateClear")
    WaitSignal("HeroGateOpen")
end

function KQGordonMaster.OnHeroGateClear(Field)
    local state = KQGordonMaster.State
    OpenDoor(Field, state.DoorHero, "GorDonDoor_01")
    BroadcastToField(Field, "¡El portón hacia la cámara intermedia se ha abierto!")
    SendSignal("HeroGateOpen")
end

function KQGordonMaster.Phase2_MidBoss(Field)
    local state = KQGordonMaster.State
    SpawnRegenGroup(Field, "KQ_Gordon_02", "GordonMaster")
    SetInterrupt(Field, "MobEliminate", "KQGordonMaster.OnMidBossClear")
    WaitSignal("MidBossOpen")
end

function KQGordonMaster.OnMidBossClear(Field)
    local state = KQGordonMaster.State
    OpenDoor(Field, state.DoorMidBoss, "GorDonDoor_02")
    OpenDoor(Field, state.GlassDoor2Boss, "GorDonDoor_03")
    BroadcastToField(Field, "¡Se han abierto las puertas a la sala del Maestro Gordon!")
    SendSignal("MidBossOpen")
end

function KQGordonMaster.Phase3_GordonBoss(Field)
    local state = KQGordonMaster.State
    state.GordonMasterHandle = SpawnMob(Field, "GordonMaster", 2800, 11000, 0, 1000, "Normal")
    SetInterrupt(Field, "DeadHandle", state.GordonMasterHandle, "KQGordonMaster.QuestSuccess")
    WaitSignal("KQDone")
end

function KQGordonMaster.QuestSuccess(Field)
    BroadcastToField(Field, "¡El Maestro Gordon ha caído! ¡Kingdom Quest completada con éxito!")
    RewardKQ(Field, "KQGordonMaster", true)
    Wait(15)
    TeleportAll(Field, "Eld", 17214, 13445)
    EndScenario(Field)
    SendSignal("KQDone")
end

function KQGordonMaster.QuestFail(Field)
    BroadcastToField(Field, "¡Misión fallida en la mazmorra de Gordon!")
    RewardKQ(Field, "KQGordonMaster", false)
    Wait(15)
    TeleportAll(Field, "Eld", 17214, 13445)
    EndScenario(Field)
    SendSignal("KQDone")
end

return KQGordonMaster
