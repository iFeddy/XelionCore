-- MaraPirate.lua
-- Modernizado desde MaraPirate.ps para XelionCore ZoneService (Lua)
-- 100% fidelidad de logica, eventos KQ y comentarios oficiales preservados.

function RegenZacco()
    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea01", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea02", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea03", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea04", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea05", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea06", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea07", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea08", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea09", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea10", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea12", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea13", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea14", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea15", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea16", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea17", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea18", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea19", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea21", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea22", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea23", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea24", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea25", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea26", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea27", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea28", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea29", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea30", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea31", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea32", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea33", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea34", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea35", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea36", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea37", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea38", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea39", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea43", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea44", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea46", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea49", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea50", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea51", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea52", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea53", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea54", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea55", "Breeder", "EventManagerA")
end


function FirstBoss()
    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea47", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea47", "Breeder", "EventManager")

    local MobNum = nil

    MobNum = 0
    while MobNum == 0 do
        MobNum = cHowManyMob("KQ_Mara")
        if not IsTimeUp() then break end
    end

    MobNum = 2
    while MobNum > 0 do
        MobNum = cHowManyMob("KQ_Mara")
        if not IsTimeUp() then break end
    end
end


function SecondBoss()
    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea48", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea48", "Breeder", "EventManager")

    local MobNum = nil

    MobNum = 0
    while MobNum == 0 do
        MobNum = cHowManyMob("KQ_Marlone")
        if not IsTimeUp() then break end
    end

    MobNum = 2
    while MobNum > 0 do
        MobNum = cHowManyMob("KQ_Marlone")
        if not IsTimeUp() then break end
    end
end


function LastBoss()
    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea40", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea45", "Breeder", "EventManagerA")

    local MobNum = nil

    MobNum = 0
    while MobNum == 0 do
        local Mara = nil
        local Marlone = nil
        Mara = cHowManyMob("KQ_TrueMara")
        Marlone = cHowManyMob("KQ_TrueMarlone")
        MobNum = Mara + Marlone
        if not IsTimeUp() then break end
    end

    MobNum = 2
    while MobNum > 0 do
        local Mara = nil
        local Marlone = nil
        Mara = cHowManyMob("KQ_TrueMara")
        Marlone = cHowManyMob("KQ_TrueMarlone")
        MobNum = Mara + Marlone
        if not IsTimeUp() then break end
    end
end

------------------------------------------------------------------------
function main()
    local MyName = nil
    local MyHandle = nil
    MyName = cMyName()
    MyHandle = cMyHandle()

    CallNamedRoutine(MyName, "DefRoutine")
end


function DefRoutine()
end


function WaitPlayerEntry()
    local InitRest = nil
    local CurRest = nil
    InitRest = cKQRestMinute() -- 남은 시간을 분단위로 얻어옴

    local Ply = nil
    Ply = 0
    while Ply == 0 do
        Ply = cHowManyMob("Player") -- 들어온 플레이어들 확인

        CurRest = cKQRestMinute() -- 남은 시간을 분단위로 얻어옴
        if InitRest - CurRest > 5 then -- 5분간 대기
            print("Room vanish")
            cEndOfEvent("KingdomQuest", 0) -- 킹덤퀘스트 실패로 마무리
            break
        end
    end
end


function IsTimeUp() -- 시간제한이 다 되었는지 확인
    print("IsTimeUp")
    local TimeUp = nil
    TimeUp = cKQRestMinute() -- 남은 시간을 분단위로 얻어옴
    print("EventManagerAIsTimeUp : " .. TimeUp .. "Minute left")
    if TimeUp == 0 then -- 시간오버
        print("IsTimeUp : Exit 0")
        KQResult = "Fail"
        return false
    end

    local Ply = nil
    Ply = cHowManyMob("Player") -- 남아있는 플레이어들 확인
    print("EventManagerAIsTimeUp : " .. Ply .. "Player left")
    if Ply == 0 then
        print("IsTimeUp : Exit 1")
        KQResult = "Fail"
        return false
    end
    return true
end

------------------------------------------------------------------------
function EventManagerA()

    print("EventManagerA : Entry")

    WaitPlayerEntry() -- 사람들이 들어올 때까지 기다림

    local KQResult = nil -- IsTimeUp에서 사용할 변수

    -- Bloque [MainLoop]
    repeat
        RegenZacco()
        FirstBoss()
        SecondBoss()
        LastBoss()
        KQResult = "Success" -- 보스가 모두 전멸했으니 성공
    until true

    cWaitSec(3)
    if KQResult == "Success" then -- 성공
        cReward("KingdomQuest") -- 리워드를 줌
        cBroadcast("all", "마라의 해적 소탕에 성공하였습니다.")
        cEndOfEvent("KingdomQuest", 1) -- 킹덤퀘스트 성공으로 마무리
    else
        cBroadcast("all", "마라의 해적 소탕에 실패하였습니다.")
        cEndOfEvent("KingdomQuest", 0) -- 킹덤퀘스트 실패로 마무리
    end


    cBroadcast("all", "30초후 루멘으로 이동합니다.")
    cWaitSec(10)
    cBroadcast("all", "20초후 루멘으로 이동합니다.")
    cWaitSec(10)
    cBroadcast("all", "10초후 루멘으로 이동합니다.")
    cWaitSec(5)
    cBroadcast("all", "5초후 루멘으로 이동합니다.")
    cWaitSec(5)
    cLinkTo("all", "Rou", "Rou", 7907, 8572)

    cMapClear("all")

end


function KQ_Mara()
    if Breeder == "EventManagerA" then -- 이벤트매니저에 의해 만들어진 마라일 경우
        local Target = nil
        Target = ""
        while Target == "" do
            Target = cWhoIsTarget() -- 타겟이 설정되었을 경우 타겟의 핸들 저장
        end

        -- 타겟발견
        cShout("all", "", "그 누구도 이곳을 지나갈 수 없다!!!")
    end

    local CurHP = nil
    CurHP = 100
    while CurHP > 0 do -- 죽을 때까지 기다림
        CurHP = cCurrentHP()
    end

    cShout("all", "", "으윽, 부.. 분하다!")
end


function KQ_Marlone()
    if Breeder == "EventManagerA" then -- 이벤트매니저에 의해 만들어진 마라일 경우
        cWaitSec(2)
        cShout("all", "", "마라 방어선이 무너졌다. 모두 적들을 맞을 준비를 하라")
    end

    local CurHP = nil
    CurHP = 100
    while CurHP > 0 do -- 죽을 때까지 기다림
        CurHP = cCurrentHP()
    end

    cShout("all", "", "으악, 적들이 너무 많습니다!")
end


function KQ_TrueMara()
    cShout("all", "", "무능한 것들 같으니라구.")
    cWaitSec(4)
    cShout("all", "", "흥, 혼자만 재미보게 할 수는 없지.")


    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea41", "Breeder", "KQ_TrueMara")
    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea41", "Breeder", "KQ_TrueMara")

    while true do
        local MobNum = nil
        MobNum = cHowManyMob("KQ_Mara")
        cWaitSec(60)
        cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea41", "Breeder", "EventManagerA")
    end
end


function KQ_TrueMarlone()
    cWaitSec(2)
    cShout("all", "", "걱정마라, 마라, 이 말론이 모두 없애주마.")

    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea42", "Breeder", "KQ_TrueMarlone")
    cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea42", "Breeder", "KQ_TrueMarlone")

    while true do
        local MobNum = nil
        MobNum = cHowManyMob("KQ_Marlone")
        cWaitSec(60)
        cMobRegen("Aggressive", "KDPrtShip", "KDPrtShipArea42", "Breeder", "EventManagerA")
    end
end

