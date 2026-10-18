-- GordonMaster.lua
-- Modernizado desde GordonMaster.ps para XelionCore ZoneService (Lua)
-- 100% fidelidad de logica, eventos KQ y comentarios oficiales preservados.

function RegenGorDonMaster()
    print("Regen RegenGorDonMaster")
    cMobRegen("Aggressive", "KDEchoCave", "KDEchoCaveArea01", "Breeder", "EventManagerA")
end


function RegenBat()
    print("Regen RegenBat")
    cMobRegen("Aggressive", "KDEchoCave", "KDEchoCaveArea07", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDEchoCave", "KDEchoCaveArea08", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDEchoCave", "KDEchoCaveArea11", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDEchoCave", "KDEchoCaveArea24", "Breeder", "EventManagerA")
end


function RegenFireViVi()
    print("Regen RegenFireViVi")
    cMobRegen("Aggressive", "KDEchoCave", "KDEchoCaveArea10", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDEchoCave", "KDEchoCaveArea12", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDEchoCave", "KDEchoCaveArea39", "Breeder", "EventManagerA")
end


function RegenSkelArcher()
    print("Regen RegenSkelArcher")
    cMobRegen("Aggressive", "KDEchoCave", "KDEchoCaveArea04", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDEchoCave", "KDEchoCaveArea09", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDEchoCave", "KDEchoCaveArea18", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDEchoCave", "KDEchoCaveArea19", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDEchoCave", "KDEchoCaveArea21", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDEchoCave", "KDEchoCaveArea28", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDEchoCave", "KDEchoCaveArea30", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDEchoCave", "KDEchoCaveArea33", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDEchoCave", "KDEchoCaveArea37", "Breeder", "EventManagerA")
end


function RegenSkeleton()
    print("Regen RegenSkeleton")
    cMobRegen("Aggressive", "KDEchoCave", "KDEchoCaveArea02", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDEchoCave", "KDEchoCaveArea03", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDEchoCave", "KDEchoCaveArea31", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDEchoCave", "KDEchoCaveArea35", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDEchoCave", "KDEchoCaveArea36", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDEchoCave", "KDEchoCaveArea42", "Breeder", "EventManagerA")
end


function RegenSkelKnight()
    print("Regen RegenSkelKnight")
    cMobRegen("Aggressive", "KDEchoCave", "KDEchoCaveArea14", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDEchoCave", "KDEchoCaveArea15", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDEchoCave", "KDEchoCaveArea20", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDEchoCave", "KDEchoCaveArea25", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDEchoCave", "KDEchoCaveArea26", "Breeder", "EventManagerA")
end


function RegenSkelWarrior()
    print("Regen RegenSkelWarrior")
    cMobRegen("Aggressive", "KDEchoCave", "KDEchoCaveArea05", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDEchoCave", "KDEchoCaveArea23", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDEchoCave", "KDEchoCaveArea27", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDEchoCave", "KDEchoCaveArea32", "Breeder", "EventManagerA")
end


function RegenSpider()
    print("Regen RegenSpider")
    cMobRegen("Aggressive", "KDEchoCave", "KDEchoCaveArea13", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDEchoCave", "KDEchoCaveArea17", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDEchoCave", "KDEchoCaveArea38", "Breeder", "EventManagerA")
end


function RegenZombie()
    print("Regen RegenZombie")
    cMobRegen("Aggressive", "KDEchoCave", "KDEchoCaveArea06", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDEchoCave", "KDEchoCaveArea29", "Breeder", "EventManagerA")
end

------------------------------------------------------------------------
function main()
    local MyName = nil
    local MyHandle = nil
    MyName = cMyName()
    MyHandle = cMyHandle()

    local KQResult = nil

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

    local NumGoden = nil

    NumGoden = cHowManyMob("KQ_GorDonMaster")
    if NumGoden == 0 then
        print("IsTimeUp : Exit 2")
        KQResult = "Success"
        return false
    end
    return true
end

------------------------------------------------------------------------
function EventManagerAMobAppend()
    local Mobs = nil
    Mobs = cHowManyMob("KQ_" .. AppendMobName)

    print("KQ_" .. AppendMobName .. "[" .. Mobs .. "]")

    if Mobs < (_G[ "Num" .. AppendMobName] or 0) then
        if _G["Regen" .. AppendMobName] then _G["Regen" .. AppendMobName]() end
    end
end


function EventManagerA()

    print("EventManagerA : Entry")

    local CurMap = nil
    CurMap = cCurrentMap() -- 현재 주인몹이 있는 맵 리턴

    -- 플레이어가 들어오기를 기다림
    WaitPlayerEntry()

    -- 모든 몹 리젠
    RegenGorDonMaster()
    RegenBat()
    RegenFireViVi()
    RegenSkelArcher()
    RegenSkeleton()
    RegenSkelKnight()
    RegenSkelWarrior()
    RegenSpider()
    RegenZombie()


    -- 기본 몹수 얻기 위해
    local NumBat = nil
    local NumFireViVi = nil
    local NumGhost = nil
    local NumOldFox = nil
    local NumSkelArcher = nil
    local NumSkeleton = nil
    local NumSkelKnight = nil
    local NumSkelWarrior = nil
    local NumSpider = nil
    local NumZombie = nil

    NumBat = cHowManyMob("KQ_Bat")
    NumFireViVi = cHowManyMob("KQ_FireViVi")
    NumSkelArcher = cHowManyMob("KQ_SkelArcher")
    NumSkeleton = cHowManyMob("KQ_Skeleton")
    NumSkelKnight = cHowManyMob("KQ_SkelKnight")
    NumSkelWarrior = cHowManyMob("KQ_SkelWarrior")
    NumSpider = cHowManyMob("KQ_Spider")
    NumZombie = cHowManyMob("KQ_Zombie")

    NumBat = NumBat / 2
    NumFireViVi = NumFireViVi / 2
    NumGhost = NumGhost / 2
    NumOldFox = NumOldFox / 2
    NumSkelArcher = NumSkelArcher / 2
    NumSkeleton = NumSkeleton / 2
    NumSkelKnight = NumSkelKnight / 2
    NumSkelWarrior = NumSkelWarrior / 2
    NumSpider = NumSpider / 2
    NumZombie = NumZombie / 2

    local AppendMobName = nil
    -- Bloque [MainLoop]
    repeat
        while true do
            -- 모자란 몹 추가
            AppendMobName = "Bat"
            EventManagerAMobAppend()
            AppendMobName = "FireViVi"
            EventManagerAMobAppend()
            AppendMobName = "SkelArcher"
            EventManagerAMobAppend()
            AppendMobName = "Skeleton"
            EventManagerAMobAppend()
            AppendMobName = "SkelKnight"
            EventManagerAMobAppend()
            AppendMobName = "SkelWarrior"
            EventManagerAMobAppend()
            AppendMobName = "Spider"
            EventManagerAMobAppend()
            AppendMobName = "Zombie"
            EventManagerAMobAppend()

            if not IsTimeUp() then break end
        end
    until true

    cWaitSec(3)
    if KQResult == "Success" then -- 성공
        cReward("KingdomQuest") -- 리워드를 줌
        cBroadcast("all", "고든마스터의 습격을 성공적으로 막아냈습니다.")
        cEndOfEvent("KingdomQuest", 1) -- 킹덤퀘스트 성공으로 마무리
    else
        cBroadcast("all", "고든마스터의 습격을 막는데 실패하였습니다.")
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

------------------------------------------------------------------------
function KQ_GorDonMaster()
    local MaxHP = nil
    MaxHP = cCurrentHP() -- 맥스 HP 측정

    local NextSummonHPRate = nil
    local NextTeleportHPRate = nil

    NextSummonHPRate = 910
    NextTeleportHPRate = 800

    local CurHP = nil
    local CurHPRate = nil
    -- Bloque [MainLoop]
    repeat
        while true do
            CurHP = cCurrentHP()
            CurHPRate = (CurHP * 1000) / MaxHP -- 주의 : 이런 수식은 뒤에서부터 실행됨!!!!

            -- 소환체크
            if CurHPRate < NextSummonHPRate then
                cShout("all", "", "깨어나거라, 어둠의 존재들이여..")
                cMobSummon("Aggressive", "KQ_GorDon", 5, "Breeder", "KQ_GorDonMaster")
                NextSummonHPRate = NextSummonHPRate - 91
            end

            -- 텔레포트 체크
            if CurHPRate < NextTeleportHPRate then
                cShout("all", "", "어리석은 것들 같으니, 능력이 되면 날 찾아봐라.")
                cSkillBlast("me", "GoMasterTeleportStart")
                NextTeleportHPRate = NextTeleportHPRate - 200
            end
            if not IsTimeUp() then break end
        end
    until true

    if KQResult == "Success" then -- 성공
        cShout("All", "", "안돼!!!!")
    else
        cShout("All", "", "후후후, 오랜만에 질좋은 영혼들을 잔뜩 모았군...")
    end
end

