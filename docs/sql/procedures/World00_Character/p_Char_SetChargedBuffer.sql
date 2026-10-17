CREATE PROCEDURE [dbo].[p_Char_SetChargedBuffer]
/*
Char Set charged buffer
2006.6 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK, 
	        ? = Error code
*/

@nCharNo int,
@nKey int,
@nID smallint,
@dUseDate datetime,
@dEndDate datetime,

@nRet int OUTPUT

AS

SET NOCOUNT ON

IF EXISTS(SELECT nCharNo FROM tCharacterChargedBuffer WHERE nCharNo = @nCharNo AND nKey = @nKey )
BEGIN
	UPDATE tCharacterChargedBuffer SET 
		nID = @nID,
		dUseDate = @dUseDate,
		dEndDate = @dEndDate
	WHERE nCharNo = @nCharNo AND nKey = @nKey
END
ELSE
BEGIN
	INSERT tCharacterChargedBuffer (  nCharNo,  nKey,  nID,  dUseDate,  dEndDate )
	VALUES                         ( @nCharNo, @nKey, @nID, @dUseDate, @dEndDate )
END

SET @nRet = @@ERROR
RETURN

-- end

0ⲈŶ
GO
