CREATE PROCEDURE [dbo].[p_Char_CreateSetDefaultData]
/*
Character Create set default data
2005.6 By CJC

 input:
   All args
   
 output: 
   0 = error
   1 = ok
*/
-- Character columns
@nCharNo int,
@sLoginZone varchar(16),
@nLoginZoneX int,
@nLoginZoneY int,
@nHP smallint, 
@nSP smallint, 
@nHPS smallint,
@nSPS smallint,
@nMoney bigint,
-- Output var
@nRet int OUTPUT

AS

SET NOCOUNT ON

UPDATE tCharacter SET
	sLoginZone	= @sLoginZone,
	nLoginZoneX	= @nLoginZoneX,
	nLoginZoneY	= @nLoginZoneY,
	nHP			= @nHP,
	nSP			= @nSP,
	nHPS		= @nHPS,
	nSPS		= @nSPS,
	nMoney		= @nMoney
WHERE nCharNo = @nCharNo

IF @@ERROR <> 0 OR @@ROWCOUNT = 0
BEGIN
	SET @nRet = 0
	RETURN
END

SET @nRet = 1
-- end

0ᴹľ
GO
