--------------------------------------------------------------------------------
--                            Wedding Ceremony Script                         --
--------------------------------------------------------------------------------

require("common")

Wedding = {}
Wedding.State = {
    GroomName = "",
    BrideName = "",
    GroomHandle = 0,
    BrideHandle = 0,
    Celebrant = 0,   -- WeddingDreian
    Book = 0,        -- PledgeBook
    Door = 0,        -- OX_gate
    Gate = 0,        -- MapLinkGate
    ClickGroom = 0,
    ClickBride = 0,
}

function Wedding.Main(Field, GroomName, BrideName)
    local state = Wedding.State
    state.GroomName = GroomName or ""
    state.BrideName = BrideName or ""

    -- Límite de tiempo: 15 minutos
    SetTimeLimit(Field, 15 * 60)
    SetInterrupt(Field, "TimeOut", "Wedding.OnWeddingEnd")

    -- Spawn de NPCs del altar y puertas
    state.Book = SpawnNPC(Field, "PledgeBook", 557, 1137, 0, 1000, "Normal")
    state.Celebrant = SpawnNPC(Field, "WeddingDreian", 501, 1106, 0, 1000, "Normal")

    state.Door = BuildDoor(Field, "OX_gate", 559, 560, 0, 1000, "Normal")
    OpenDoor(Field, state.Door, "door1")

    SetScriptFile(Field, "Wedding")

    -- Tiempo de espera previo (recordatorios cada 1 minuto)
    for i = 1, 4 do
        Wait(60)
        SendChatWindow(Field, "WeddingDreian", "Wedding1")
    end
    Wait(60)

    -- Inicio de la ceremonia nupcial
    CloseDoor(Field, state.Door, "door1")

    state.GroomHandle = FindPlayer(Field, state.GroomName, "Wedding")
    state.BrideHandle = FindPlayer(Field, state.BrideName, "Wedding")

    -- Teletransporte de invitados a las bancas y novios a la entrada
    TeleportAll(Field, "Split", 304, 833, 808, 833)
    Teleport(state.GroomHandle, 536, 580)
    Teleport(state.BrideHandle, 582, 580)

    -- Caminata de novios al primer punto
    PlayerWalk(state.GroomHandle, 536, 620)
    PlayerWalk(state.BrideHandle, 582, 620)

    SendChatWindow(Field, "WeddingDreian", "Wedding2")
    Wait(10)

    SendChatWindow(Field, "WeddingDreian", "Wedding3")
    Wait(10)

    -- Caminata de los novios hacia el altar
    PlayerWalk(state.GroomHandle, 536, 1078)
    PlayerWalk(state.BrideHandle, 582, 1078)

    -- Palabras del celebrante Dreian (votos)
    for step = 4, 12 do
        SendChatWindow(Field, "WeddingDreian", "Wedding" .. step)
        Wait(10)
    end

    -- Lectura del libro de votos (PledgeBook)
    for step = 1, 9 do
        SendChatWindow(Field, "PledgeBook", "Book" .. step)
        Wait(10)
    end

    SendChatWindow(Field, "WeddingDreian", "Wedding13")
    Wait(10)

    state.GroomHandle = FindPlayer(Field, state.GroomName, "Wedding")
    state.BrideHandle = FindPlayer(Field, state.BrideName, "Wedding")

    -- Completar matrimonio oficialmente
    CompleteWedding(state.GroomHandle, state.BrideHandle)

    -- Aplicar estado anormal/buff de casados por 3 días (259.200.000 ms)
    ApplyAbnormalState(state.GroomHandle, "StaWeddingSub02", 1, 259200000)
    ApplyAbnormalState(state.BrideHandle, "StaWeddingSub02", 1, 259200000)

    SendChatWindow(Field, "WeddingDreian", "Wedding14")
    SendChatWindow(Field, "WeddingDreian", "Wedding15")

    -- Abrir puerta y permitir salida
    OpenDoor(Field, state.Door, "door1")

    state.GroomHandle = FindPlayer(Field, state.GroomName, "Wedding")
    state.BrideHandle = FindPlayer(Field, state.BrideName, "Wedding")

    PlayerWalk(state.GroomHandle, 536, 580)
    PlayerWalk(state.BrideHandle, 582, 580)

    -- Portal de retorno a Elderine
    state.Gate = SpawnNPC(Field, "MapLinkGate", 549, 129, 0, 1000, "Normal")
    SetNPCInteractHandler(state.Gate, "Wedding.OnGateClick")

    -- 5 minutos para fotos y celebración antes de cerrar la sala
    Wait(5 * 60)
    Wedding.OnWeddingEnd(Field)
end

function Wedding.OnGateClick(Field, PlayerHandle)
    Teleport(PlayerHandle, "Eld", 17214, 13445)
end

function Wedding.OnWeddingEnd(Field)
    local state = Wedding.State
    TeleportAll(Field, "Eld", 17214, 13445)
    Vanish(Field, state.Book)
    Vanish(Field, state.Celebrant)
    Vanish(Field, state.Gate)
    EndScenario(Field)
end

return Wedding
