--------------------------------------------------------------------------------
--                     Job Promotion 1 (First Class Change)                   --
--------------------------------------------------------------------------------

require("common")

JobChange1 = {}
JobChange1.State = {
    PlayerHandle = 0,
    ShadowHandle = 0,
    ShadowMobIndex = "",
    Doors = {},
    ChiefA = 0,
    ChiefD = 0,
    ChiefAA = 0,
    ChiefDD = 0,
    Statues = {},
    StatusBreakInterval = 10,
}

function JobChange1.Main(Field)
    local state = JobChange1.State

    -- Construir puertas del laberinto
    state.Doors[1] = BuildDoor(Field, "KQ_Gate4", 1018, 1560, 0, 1000, "Normal")
    CloseDoor(Field, state.Doors[1], "Door01")

    state.Doors[2] = BuildDoor(Field, "KQ_Gate4", 2098, 3291, 270, 1000, "Normal")
    CloseDoor(Field, state.Doors[2], "Door02")

    state.Doors[3] = BuildDoor(Field, "KQ_Gate4", 3845, 3266, 90, 1000, "Normal")
    CloseDoor(Field, state.Doors[3], "Door03")

    state.Doors[4] = BuildDoor(Field, "KQ_Gate4", 5162, 4629, 0, 1000, "Normal")
    OpenDoor(Field, state.Doors[4], "Door04")

    SetScriptFile(Field, "JobChange1")

    -- Esperar conexión del jugador
    state.PlayerHandle = WaitPlayerLogin(Field, 60)
    if not state.PlayerHandle or state.PlayerHandle == 0 then
        JobChange1.QuestFail(Field)
        return
    end

    -- Identificar clase para generar la sombra contraria
    local charClass = GetPlayerClass(state.PlayerHandle)
    state.ShadowMobIndex = "Job1_" .. charClass .. "-m"

    -- Fases del escenario
    JobChange1.FirstContact(Field)
    JobChange1.FightSkeletons(Field)
    JobChange1.RoadToLastRoom(Field)
    JobChange1.LastBattle(Field)
end

function JobChange1.FirstContact(Field)
    local state = JobChange1.State
    SetInterrupt(Field, "TimeOut", "JobChange1.QuestFail")
    SetInterrupt(Field, "PlayerEliminate", "JobChange1.QuestFail")

    -- Primer contacto: La sombra aparece y habla
    state.ShadowHandle = SpawnMob(Field, state.ShadowMobIndex, 1000, 1200, 0, 1000, "Normal")
    SendNPCChat(state.ShadowHandle, "FirstContact0")
    Wait(2)
    SendNPCChat(state.ShadowHandle, "FirstContact1")
    Wait(2)

    -- Abrir puerta 1
    OpenDoor(Field, state.Doors[1], "Door01")
    SetInterrupt(Field, "AreaEntry", "Zone_Mob01", "JobChange1.OnEnterSkelRoom")
end

function JobChange1.OnEnterSkelRoom(Field)
    local state = JobChange1.State
    CloseDoor(Field, state.Doors[1], "Door01")
    CloseDoor(Field, state.Doors[2], "Door02")
    CloseDoor(Field, state.Doors[3], "Door03")

    -- Invocación de esqueletos
    SendNPCChat(state.ShadowHandle, "SkellBreed0")
    Wait(2)
    SendNPCChat(state.ShadowHandle, "SkellBreed2")

    local coords = {
        {"Job1_Skeleton", 2773, 3306},
        {"Job1_Skeleton", 3619, 3786},
        {"Job1_Skeleton", 2340, 2868},
        {"Job1_SkelArcher", 2446, 3302},
        {"Job1_SkelArcher", 3582, 2803},
        {"Job1_SkelArcher", 2476, 3774},
    }
    for _, mob in ipairs(coords) do
        SpawnMob(Field, mob[1], mob[2], mob[3], 0, 1000, "Normal")
    end

    Vanish(Field, state.ShadowHandle)
    SetInterrupt(Field, "MobEliminate", "JobChange1.OnSkeletonsDefeated")
end

function JobChange1.FightSkeletons(Field)
    -- Se espera a que todos los esqueletos sean eliminados
    WaitSignal("SkelBreak")
end

function JobChange1.OnSkeletonsDefeated(Field)
    SendSignal("SkelBreak")
end

function JobChange1.RoadToLastRoom(Field)
    local state = JobChange1.State
    -- Invocar Kebings atados
    for i = 1, 3 do
        local h = SpawnMob(Field, "Job1_Kebing", 5004, 4056, 0, 1000, "Normal")
        ApplyAbnormalState(h, "StaQuestEntangle", 1, 4200000000)
    end

    OpenDoor(Field, state.Doors[3], "Door03")

    SetInterrupt(Field, "AreaEntry", "Zone_Mob03", "JobChange1.OnFreeKebings")
    SetInterrupt(Field, "AreaEntry", "Zone_Mob04", "JobChange1.OnLightOff")
    WaitSignal("RoadComplete")
end

function JobChange1.OnFreeKebings(Field)
    -- Liberar a los Kebings del entangle
    RemoveAbnormalStateFromGroup(Field, "Job1_Kebing", "StaQuestEntangle")
end

function JobChange1.OnLightOff(Field)
    local state = JobChange1.State
    CloseDoor(Field, state.Doors[4], "Door04")
    SetMapLight(Field, 0, 0, 0) -- Luces apagadas

    -- Antorchas encendidas en la oscuridad
    local torches = {
        {109, 4832}, {106, 4161}, {578, 3690}, {1245, 3689},
        {1719, 4157}, {1720, 4826}, {1250, 5300}
    }
    for i, t in ipairs(torches) do
        state.Statues[i] = SpawnEffect(Field, "Job1_Dn_Gob" .. string.format("%02d", i), t[1], t[2], 180, 3600000, 1000)
    end

    -- Jefes Marlone
    state.ChiefA = SpawnMob(Field, "Job1_MarloneFighter", 5867, 6072, 0, 1000, "Normal")
    state.ChiefD = SpawnMob(Field, "Job1_MarloneArcher", 4700, 6049, 0, 1000, "Normal")

    SendSignal("RoadComplete")
end

function JobChange1.LastBattle(Field)
    local state = JobChange1.State
    SetInterrupt(Field, "AreaEntry", "Zone_Mob05", "JobChange1.OnLightOn")
    WaitSignal("BattleFinished")
end

function JobChange1.OnLightOn(Field)
    local state = JobChange1.State

    -- Secuencia de encendido gradual de la sala final
    local lightValues = {36, 72, 108, 144, 180, 216, 255}
    for _, lv in ipairs(lightValues) do
        SetMapLight(Field, lv, lv, lv)
        WaitMs(500)
    end

    -- Inmovilizar al jugador brevemente para la escena
    ApplyAbnormalState(state.PlayerHandle, "StaQuestEntangle", 1, 4200000000)

    -- Aparición de la Sombra Final
    state.ShadowHandle = SpawnMob(Field, state.ShadowMobIndex, 1380, 4500, 90, 1000, "Normal", state.PlayerHandle)
    SendNPCChat(state.ShadowHandle, "Last0")
    Wait(2)
    SendNPCChat(state.ShadowHandle, "Last1")
    Wait(2)
    SendNPCChat(state.ShadowHandle, "Last2")

    RemoveAbnormalState(state.PlayerHandle, "StaQuestEntangle")

    -- Jefes adicionales si no fueron eliminados
    if state.ChiefA and state.ChiefA ~= 0 then
        Vanish(Field, state.ChiefA)
        state.ChiefAA = SpawnMob(Field, "Job1_MarloneFighter", 1380, 4500, 0, 1000, "Normal")
        SetAggro(state.ChiefAA, state.PlayerHandle, 1)
    end
    if state.ChiefD and state.ChiefD ~= 0 then
        Vanish(Field, state.ChiefD)
        state.ChiefDD = SpawnMob(Field, "Job1_MarloneArcher", 1380, 4500, 0, 1000, "Normal")
        SetAggro(state.ChiefDD, state.PlayerHandle, 1)
    end

    -- Regla de victoria al derrotar a la sombra
    SetInterrupt(Field, "DeadHandle", state.ShadowHandle, "JobChange1.QuestSuccess")

    -- Rompimiento periódico de estatuas que liberan goblins
    SetTimer(Field, "JobChange1.StatusBreak", 10, true)
end

function JobChange1.StatusBreak(Field)
    local state = JobChange1.State
    for i, s in ipairs(state.Statues) do
        if s ~= 0 then
            state.Statues[i] = 0
            Vanish(Field, s)
            SpawnMob(Field, "Job1_Goblin", 1200, 4500, 0, 1000, "Normal")
            break
        end
    end
end

function JobChange1.QuestSuccess(Field)
    local state = JobChange1.State
    SpawnEffectOnObject(state.PlayerHandle, "Job1_Dn_ShadowEnd", 4000, 1000)
    Wait(3)
    CompleteQuestScenario(Field, state.PlayerHandle, true)
    SendSignal("BattleFinished")
end

function JobChange1.QuestFail(Field)
    local state = JobChange1.State
    if state.PlayerHandle and state.PlayerHandle ~= 0 then
        CompleteQuestScenario(Field, state.PlayerHandle, false)
        Teleport(state.PlayerHandle, "Rou", 6208, 6405)
    end
    EndScenario(Field)
end

return JobChange1
