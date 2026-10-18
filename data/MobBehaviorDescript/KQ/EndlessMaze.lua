-- EndlessMaze.lua
-- Modernizado desde EndlessMaze.ps para XelionCore ZoneService (Lua)
-- 100% fidelidad de logica, eventos KQ y comentarios oficiales preservados.

function RegenZaco()
    cMobRegen("Aggressive", "KDEnMaze", "KDEndless_MiroArea01", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDEnMaze", "KDEndless_MiroArea02", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDEnMaze", "KDEndless_MiroArea03", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDEnMaze", "KDEndless_MiroArea04", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDEnMaze", "KDEndless_MiroArea05", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDEnMaze", "KDEndless_MiroArea07", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDEnMaze", "KDEndless_MiroArea08", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDEnMaze", "KDEndless_MiroArea10", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDEnMaze", "KDEndless_MiroArea11", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDEnMaze", "KDEndless_MiroArea12", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDEnMaze", "KDEndless_MiroArea13", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDEnMaze", "KDEndless_MiroArea14", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDEnMaze", "KDEndless_MiroArea15", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDEnMaze", "KDEndless_MiroArea16", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDEnMaze", "KDEndless_MiroArea17", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDEnMaze", "KDEndless_MiroArea18", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDEnMaze", "KDEndless_MiroArea19", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDEnMaze", "KDEndless_MiroArea20", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDEnMaze", "KDEndless_MiroArea21", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDEnMaze", "KDEndless_MiroArea22", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDEnMaze", "KDEndless_MiroArea24", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDEnMaze", "KDEndless_MiroArea25", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDEnMaze", "KDEndless_MiroArea26", "Breeder", "EventManagerA")
    cMobRegen("Aggressive", "KDEnMaze", "KDEndless_MiroArea32", "Breeder", "EventManagerA")
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
function EventManagerA()

    print("EventManagerA : Entry")

    local CurMap = nil
    CurMap = cCurrentMap() -- 현재 주인몹이 있는 맵 리턴

    -- 플레이어가 들어오기를 기다림
    WaitPlayerEntry()



    -- 모든 몹 리젠
    RegenZaco()

    -- Bloque [MainLoop]
    repeat
        while true do
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

