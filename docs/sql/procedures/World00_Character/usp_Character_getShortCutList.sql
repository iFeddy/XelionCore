CREATE PROCEDURE [dbo].[usp_Character_getShortCutList]
	@nCharNo int,
	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	SET @nRet = 0
	
	IF EXISTS(SELECT nCharNo FROM tCharacter WHERE nCharNo = @nCharNo AND bDeleted = 0) BEGIN
		SELECT nSlotNo, nCodeNo, nValue
		  FROM tCharacterShortCut
		 WHERE nCharNo = @nCharNo AND nCodeNo > 0
	END
	ELSE BEGIN
		SET @nRet = -2721
	END
	
	SET NOCOUNT OFF
END

0ꠁ䵎ı
GO
