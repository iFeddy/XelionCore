--------------------------------------------------------------------------------
--                         Event Monkey Boss Scenario                         --
--------------------------------------------------------------------------------

require("common")

MonkeyEvent = {}
MonkeyEvent.State = {
    BossHandle = 0,
    BossX = 0,
    BossY = 0,
    BossIndex = "Event_Monkey_Boss",
    SummonIndex = "Event_Monkey_Sub",
}

function MonkeyEvent.Main(Field, X, Y)
    local state = MonkeyEvent.State
    state.BossX = X or 5000
    state.BossY = Y or 5000

    -- Spawn del Rey Mono del evento con defensa potenciada
    state.BossHandle = SpawnMob(Field, state.BossIndex, state.BossX, state.BossY, 90, 1000, "Normal")
    ApplyAbnormalState(state.BossHandle, "StaMonkeyDefense", 1, 4000000000)

    -- Interrupciones: invoca monitos al bajar del 90% de HP
    SetInterrupt(Field, "HPLow", state.BossHandle, 900, "MonkeyEvent.OnSummonMonkeys")
    SetInterrupt(Field, "DeadHandle", state.BossHandle, "MonkeyEvent.OnBossDead")
end

function MonkeyEvent.OnSummonMonkeys(Field)
    local state = MonkeyEvent.State
    -- Cura al jefe para reiniciar la fase de invocación
    SetHP(state.BossHandle, GetMaxHP(state.BossHandle))

    -- Invoca horda de monos pequeños en círculos
    for i = 1, 6 do
        local x, y = GetCircularScatter(state.BossX, state.BossY, 400)
        SpawnMob(Field, state.SummonIndex, x, y, 0, 1000, "Normal")
    end
end

function MonkeyEvent.OnBossDead(Field)
    local state = MonkeyEvent.State
    BroadcastToField(Field, "¡El Rey Mono ha sido derrotado! ¡Recojan las recompensas del evento!")
    -- Recompensas del evento
    Vanish(Field, state.BossHandle)
end

return MonkeyEvent
