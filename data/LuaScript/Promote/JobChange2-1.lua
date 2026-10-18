--------------------------------------------------------------------------------
--                     Job Promotion 2-1 (Class Change 2)                     --
--------------------------------------------------------------------------------

require("common")

JobChange2_1 = {}
JobChange2_1.State = {
    PlayerHandle = 0,
    Roumenos = 0,
    NPCs = {},
    Beggars = {},
    Zombies = {},
    ClickNum = 0,
    ClickHandle = 0,
}

function JobChange2_1.Main(Field)
    local state = JobChange2_1.State
    SetScriptFile(Field, "JobChange2-1")

    -- Spawn de NPCs de la aldea fantasma de Roumen
    local npcs = {
        {"RouSmithJames", 5937, 6644, 0},
        {"RouSoulMctJulia", 8032, 6469, 0},
        {"RouWeaponTitleMctZach", 4985, 5853, 171},
        {"RouSkillRubi", 7129, 6115, 173},
        {"RouItemMctPey", 7133, 5428, 0},
        {"RouStoreRaina", 6139, 5197, 342},
        {"RouRookieGuideRaemi", 7720, 4623, 80},
        {"RouGrandfatherRobin", 4478, 6735, 0},
        {"RouGuildBrom", 4948, 5015, 0},
    }

    for _, n in ipairs(npcs) do
        local h = SpawnNPC(Field, n[1], n[2], n[3], n[4], 1000, "Normal")
        table.insert(state.NPCs, h)
    end

    state.Roumenos = SpawnNPC(Field, "RouTownChiefRoumenus", 5988, 4677, 178, 1000, "Normal")
    table.insert(state.NPCs, state.Roumenos)

    -- Spawn de mendigos dispersos en círculos
    local centers = {
        {6242, 4646, 1000, 15},
        {5807, 6234, 500, 5},
    }
    for _, c in ipairs(centers) do
        for i = 1, c[4] do
            local x, y = GetCircularScatter(c[1], c[2], c[3])
            local d = math.random(0, 360)
            local h = SpawnMob(Field, "Slime", x, y, d, 1000, "Normal", "Beggar")
            ChangeToNPC(Field, h, "Slime")
            table.insert(state.Beggars, h)
        end
    end

    -- Esperar conexión del jugador
    state.PlayerHandle = WaitPlayerLogin(Field, 60)
    if not state.PlayerHandle or state.PlayerHandle == 0 then
        JobChange2_1.QuestFail(Field)
        return
    end

    -- Fase 1: El jugador interactúa con los habitantes extraños
    JobChange2_1.WaitArrival(Field)

    -- Fase 2: Niebla roja y transformación de los NPCs en Zombies
    JobChange2_1.Zombinaze(Field)

    -- Fase 3: Batalla contra la horda de zombies (20 min)
    JobChange2_1.ZombieBattle(Field)
end

function JobChange2_1.WaitArrival(Field)
    local state = JobChange2_1.State
    SetInterrupt(Field, "PlayerEliminate", "JobChange2_1.QuestFail")
    SetInterrupt(Field, "NPCClickAny", "JobChange2_1.OnNPCClick")
    WaitSignal("RoumenosReached")
end

function JobChange2_1.OnNPCClick(Field, ClickedNPC)
    local state = JobChange2_1.State
    SendNPCChat(ClickedNPC, "NoChat") -- Responden "..."

    if ClickedNPC == state.Roumenos then
        SendSignal("RoumenosReached")
    end
end

function JobChange2_1.Zombinaze(Field)
    local state = JobChange2_1.State
    ClearInterrupts(Field, "NPCClickAny")

    -- Efectos visuales de corrupción: luz rojiza y niebla espesa
    SetMapLight(Field, 150, 70, 70)
    SetMapFog(Field, 100, 1400)

    -- Todos los NPCs y mendigos gritan y se transforman en Job2_Zombie
    local allEntities = {}
    for _, h in ipairs(state.NPCs) do table.insert(allEntities, h) end
    for _, h in ipairs(state.Beggars) do table.insert(allEntities, h) end

    local screamNo = 0
    for _, npcHandle in ipairs(allEntities) do
        SendNPCChat(npcHandle, "Scream" .. screamNo)
        local zombieHandle = ExchangeToMob(Field, npcHandle, "Job2_Zombie")
        table.insert(state.Zombies, zombieHandle)
        SendNPCChat(zombieHandle, "Scream" .. screamNo)

        screamNo = (screamNo + 1) % 3
    end
end

function JobChange2_1.ZombieBattle(Field)
    local state = JobChange2_1.State
    SetTimeLimit(Field, 20 * 60)
    SetInterrupt(Field, "TimeOut", "JobChange2_1.QuestFail")
    SetInterrupt(Field, "MobEliminate", "JobChange2_1.QuestSuccess")
    SetInterrupt(Field, "DeadIndex", "Job2_Zombie", "JobChange2_1.OnZombieKilled")

    WaitSignal("BattleDone")
end

function JobChange2_1.OnZombieKilled(Field)
    local remaining = CountAliveMobs(Field, "Job2_Zombie")
    if remaining == 0 then
        JobChange2_1.QuestSuccess(Field)
        return
    end

    local displayCount = math.min(remaining, 10)
    BroadcastToField(Field, "ReftMob" .. displayCount)
end

function JobChange2_1.QuestSuccess(Field)
    local state = JobChange2_1.State
    RemoveItem(state.PlayerHandle, "Job2_STpiece1")
    CreateItem(state.PlayerHandle, "Job2_STpiece1-1", 1)
    Wait(1)

    SendQuestResult(Field, state.PlayerHandle, true)
    SetScenarioCode(state.PlayerHandle, 16)
    Wait(2)

    -- Cuenta regresiva de salida
    local returnTimes = {30, 20, 10, 5}
    for _, sec in ipairs(returnTimes) do
        BroadcastToField(Field, "KQReturn" .. sec)
        Wait(sec == 5 and 5 or (sec == 10 and 5 or 10))
    end

    TeleportAll(Field, "RouVal01", 12842, 6467)
    EndScenario(Field)
    SendSignal("BattleDone")
end

function JobChange2_1.QuestFail(Field)
    local state = JobChange2_1.State
    SendQuestResult(Field, state.PlayerHandle, false)

    local returnTimes = {30, 20, 10, 5}
    for _, sec in ipairs(returnTimes) do
        BroadcastToField(Field, "RouReturn" .. sec)
        Wait(sec == 5 and 5 or (sec == 10 and 5 or 10))
    end

    TeleportAll(Field, "RouVal01", 12842, 6467)
    EndScenario(Field)
    SendSignal("BattleDone")
end

return JobChange2_1
