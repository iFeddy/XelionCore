CREATE PROCEDURE [dbo].[usp_Guild_TokenSub]
/*
Sub money
2008.12 By Lunar
 input:
	see SQL
 output:
	@nRet	1 = OK
	        0 = Error, or Minus value
*/
@nGuildNo int,
@nSubToken bigint,
@nRet int OUTPUT
AS
SET NOCOUNT ON
DECLARE @nGuildToken bigint
EXEC usp_Guild_GetToken @nGuildNo, @nGuildToken OUTPUT, @nRet OUTPUT
IF @nRet <> 0
	RETURN
-- Sub money
SET @nGuildToken = @nGuildToken - @nSubToken
-- Check minus
IF @nGuildToken < 0
BEGIN
	SET @nRet = -1
	RETURN	
END
EXEC usp_Guild_SetToken @nGuildNo, @nGuildToken, @nRet OUTPUT
-- End

0
GO
