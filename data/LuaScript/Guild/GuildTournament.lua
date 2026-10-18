--------------------------------------------------------------------------------
--                       Guild Tournament Battleground                        --
--------------------------------------------------------------------------------

require("common")

GuildTournament = {}
GuildTournament.State = {
    GuildA = "",
    GuildB = "",
    DoorRed = 0,
    DoorBlue = 0,
    Flags = {},
    Portals = {},
    GatesA = {},
    GatesB = {},
}

function GuildTournament.Main(Field, GuildA, GuildB, ReadyMin, PlayTimeMin)
    local state = GuildTournament.State
    state.GuildA = GuildA or "GuildA"
    state.GuildB = GuildB or "GuildB"
    ReadyMin = ReadyMin or 1
    PlayTimeMin = PlayTimeMin or 30

    -- Construir y cerrar portones iniciales de ambos bandos
    state.DoorRed = BuildDoor(Field, "GT_BigGate", 1964, 2124, 35, 4000, "Normal")
    state.DoorBlue = BuildDoor(Field, "GT_BigGate", 10483, 10518, 210, 4000, "Normal")
    CloseDoor(Field, state.DoorRed, "RedDoor")
    CloseDoor(Field, state.DoorBlue, "BlueDoor")

    -- Banderas de captura (1 a 5)
    local flagCoords = {
        {"Normal_Flag00", 2968, 2868},
        {"Normal_Flag01", 6547, 6268},
        {"Normal_Flag02", 9711, 9559},
        {"Normal_Flag03", 2475, 10710},
        {"Normal_Flag04", 10600, 2247},
    }
    for i, fc in ipairs(flagCoords) do
        local h = SpawnNPC(Field, fc[1], fc[2], fc[3], 0, 2000, "Normal")
        state.Flags[i] = h
        RegisterCaptureFlag(Field, h, "Red_Flag0" .. (i - 1), "Blue_Flag0" .. (i - 1))
    end

    -- Portales de captura (1 a 5)
    local portalCoords = {
        {5084, 3326}, {5361, 7909}, {7949, 9399}, {4118, 10360}, {10313, 4050}
    }
    for i, pc in ipairs(portalCoords) do
        local h = SpawnNPC(Field, "GT_TargetGate_Normal", pc[1], pc[2], 0, 2000, "Normal")
        state.Portals[i] = h
        RegisterGTPortal(Field, h, state.Flags[i], "GT_TargetGate_Red", "GT_TargetGate_Blue")
    end

    -- Puertas de base del Gremio A (Rojo)
    local gatesACoords = {
        {1784, 1045}, {1607, 1032}, {1429, 1041}, {1954, 1050}, {2114, 1055}
    }
    for i, gc in ipairs(gatesACoords) do
        local h = SpawnNPC(Field, "GuildGate00", gc[1], gc[2], 240, 1000, "Normal")
        state.GatesA[i] = h
        RegisterGuildGate(Field, h, state.GuildA, state.Flags[i], "GuildGate0" .. i)
    end

    -- Puertas de base del Gremio B (Azul)
    local gatesBCoords = {
        {10878, 11400}, {10700, 11400}, {10520, 11400}, {11050, 11400}, {11230, 11400}
    }
    for i, gc in ipairs(gatesBCoords) do
        local h = SpawnNPC(Field, "GuildGate00", gc[1], gc[2], 60, 1000, "Normal")
        state.GatesB[i] = h
        RegisterGuildGate(Field, h, state.GuildB, state.Flags[i], "GuildGate0" .. i)
    end

    -- Tiempo de preparación (cuenta regresiva)
    BroadcastToField(Field, "¡El Torneo de Gremios comenzará en " .. ReadyMin .. " minuto(s)!")
    Wait(ReadyMin * 60)

    -- Abrir portones principales
    OpenDoor(Field, state.DoorRed, "RedDoor")
    OpenDoor(Field, state.DoorBlue, "BlueDoor")
    BroadcastToField(Field, "¡Las puertas se han abierto! ¡Que comience la batalla!")

    -- Límite de tiempo del combate
    SetTimeLimit(Field, PlayTimeMin * 60)
    SetInterrupt(Field, "TimeOut", "GuildTournament.OnTimeOut")
    SetInterrupt(Field, "GuildCoreDestroyed", "GuildTournament.OnCoreDestroyed")
end

function GuildTournament.OnCoreDestroyed(Field, WinningGuild, LosingGuild)
    BroadcastToField(Field, "¡El núcleo del gremio " .. LosingGuild .. " ha sido destruido!")
    BroadcastToField(Field, "¡Victoria para " .. WinningGuild .. "!")
    GuildTournament.EndMatch(Field, WinningGuild)
end

function GuildTournament.OnTimeOut(Field)
    -- Si expira el tiempo, gana el gremio con más banderas y puntos de captura
    local winGuild = EvaluateMatchScore(Field)
    BroadcastToField(Field, "¡Tiempo agotado! El gremio vencedor es: " .. winGuild)
    GuildTournament.EndMatch(Field, winGuild)
end

function GuildTournament.EndMatch(Field, WinningGuild)
    RewardGuildTournament(Field, WinningGuild)
    Wait(15)
    TeleportAll(Field, "Eld", 17214, 13445)
    EndScenario(Field)
end

return GuildTournament
