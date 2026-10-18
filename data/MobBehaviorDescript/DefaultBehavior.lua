-- DefaultBehavior.lua
-- Modernizado desde DefaultBehavior.ps para XelionCore ZoneService (Lua)
-- 100% fidelidad de logica, eventos KQ y comentarios oficiales preservados.

function main()
    local Enemy = nil
    local Perm = nil
    local EnemyName = nil
    while true do
        Enemy = cWhoIsTarget() -- 현재 타겟의 핸들을 얻어옴
        if Enemy ~= "" then -- 타겟이 있음
            EnemyName = cGetName(Enemy) -- 핸들(Enemy)의 이름을 EnemyName에 넣음

            Perm = cPermillage()

            if Perm < 300 then
                cChat(Enemy, EnemyName .. ", 게 섯거라!")
            elseif Perm < 700 then
                cChat(Enemy, "각오해라, " .. EnemyName .. "!")
            else
                cChat(Enemy, "적이다!!")
            end
        end

        Enemy = cWhoKillMe() -- 내가 죽었을 경우 나를 죽인 캐릭의 핸들을 얻어옴
        if Enemy ~= "" then -- 나를 죽인 캐릭이 있음 - 난 죽어있는 상태
            cChat(Enemy, "억울하다!!!")
            break
        end
    end
end

