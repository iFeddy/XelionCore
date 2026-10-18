--------------------------------------------------------------------------------
--                         OX Quiz Event Scenario Script                      --
--------------------------------------------------------------------------------

require("common")

OXField = {}
OXField.State = {
    Door = 0,
    Effect = 0,
    IsRunning = false,
}

function OXField.Main(Field)
    local state = OXField.State

    -- Construir y abrir puerta central divisoria entre las zonas O y X
    state.Door = BuildDoor(Field, "OX_gate", 6457, 3655, 0, 1000, "Normal")
    OpenDoor(Field, state.Door, "OXBlock")

    -- Registrar comandos de administración para el Game Master
    RegisterAdminCommand(Field, "ox_block", "OXField.StartBlocking")
    RegisterAdminCommand(Field, "ox_answer_o", "OXField.AnswerO")
    RegisterAdminCommand(Field, "ox_answer_x", "OXField.AnswerX")
    RegisterAdminCommand(Field, "ox_clear_block", "OXField.BlockClear")
    RegisterAdminCommand(Field, "ox_clear_arena", "OXField.ArenaClear")
end

-- Inicia la cuenta regresiva de 15 segundos para bloquear el centro
function OXField.StartBlocking(Field, AdminHandle)
    local state = OXField.State
    if state.IsRunning then
        SendChatMessage(AdminHandle, "El evento OX ya tiene una ronda en curso.")
        return
    end

    state.IsRunning = true
    BroadcastToField(Field, "¡Tienen 15 segundos para elegir bando (O o X) antes de que se cierre la puerta central!")
    SetTimeLimit(Field, 15)
    SetInterrupt(Field, "TimeOut", "OXField.CenterBlock")
end

function OXField.CenterBlock(Field)
    local state = OXField.State
    CloseDoor(Field, state.Door, "OXBlock")
    BroadcastToField(Field, "¡Puertas cerradas! ¡Las respuestas han quedado fijadas!")
    state.IsRunning = false
end

-- La respuesta correcta fue O (se elimina a los de la zona X)
function OXField.AnswerO(Field, AdminHandle)
    local state = OXField.State
    BroadcastBriefEffect(Field, 1) -- Efecto visual de O
    Wait(4)

    state.Effect = SpawnEffect(Field, "OX_act", 6942, 3639, 0, 5000, 1000)
    Wait(5)

    -- Expulsar a los jugadores de la zona perdedora (XZone01) a las gradas
    BanFromArea(Field, "XZone01", "OX_field", 6467, 4771)
    BroadcastToField(Field, "¡La respuesta correcta era O! Los participantes en X han sido descalificados.")
end

-- La respuesta correcta fue X (se elimina a los de la zona O)
function OXField.AnswerX(Field, AdminHandle)
    local state = OXField.State
    BroadcastBriefEffect(Field, 2) -- Efecto visual de X
    Wait(4)

    state.Effect = SpawnEffect(Field, "OX_act", 5940, 3646, 0, 5000, 1000)
    Wait(5)

    -- Expulsar a los jugadores de la zona perdedora (OZone01) a las gradas
    BanFromArea(Field, "OZone01", "OX_field", 6467, 4771)
    BroadcastToField(Field, "¡La respuesta correcta era X! Los participantes en O han sido descalificados.")
end

-- Reabrir puerta para la siguiente ronda de preguntas
function OXField.BlockClear(Field, AdminHandle)
    local state = OXField.State
    OpenDoor(Field, state.Door, "OXBlock")
    BroadcastToField(Field, "La puerta central está abierta. ¡Prepárense para la siguiente pregunta!")
end

-- Limpiar la arena completa enviando a todos a las gradas de espectadores
function OXField.ArenaClear(Field, AdminHandle)
    BanFromArea(Field, "OZone01", "OX_field", 6467, 4771)
    BanFromArea(Field, "XZone01", "OX_field", 6467, 4771)
    BroadcastToField(Field, "La arena de OX ha sido limpiada.")
end

return OXField
