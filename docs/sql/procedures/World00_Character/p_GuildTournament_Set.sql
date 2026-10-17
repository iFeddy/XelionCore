CREATE PROCEDURE [dbo].[p_GuildTournament_Set]
/*
Set/Update Guild Tournament Data
2007.5 By CJC
*/
@nGTNo				int,
@nType				smallint,
@nData				int,
@nRet				int OUTPUT
AS
SET NOCOUNT ON
UPDATE tGuildTournament SET	nData = @nData WHERE nGTNo = @nGTNo AND nType = @nType
IF @@ERROR <> 0 OR @@ROWCOUNT = 0
BEGIN
	INSERT tGuildTournament (  nGTNo,  nType,  nData )
	VALUES					( @nGTNo, @nType, @nData )
END
SET @nRet = @@ERROR
-- end

0
GO
