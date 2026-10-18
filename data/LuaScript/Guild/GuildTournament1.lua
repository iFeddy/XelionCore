--------------------------------------------------------------------------------
--                      Guild Tournament 1 (Match Variant)                    --
--------------------------------------------------------------------------------

require("common")
require("Guild/GuildTournament")

GuildTournament1 = {}

function GuildTournament1.Main(Field, GuildA, GuildB, ReadyMin, PlayTimeMin)
    -- Ejecuta la lógica principal del torneo de gremios con reglas de variante 1
    return GuildTournament.Main(Field, GuildA, GuildB, ReadyMin or 1, PlayTimeMin or 20)
end

return GuildTournament1
