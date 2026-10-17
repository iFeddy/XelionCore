CREATE PROCEDURE [dbo].[usp_Character_getSingleOption]
	@nCharNo      int
,	@nRet         int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000 
	
	SET @nRet = 0

	IF NOT EXISTS(SELECT nCharNo FROM tCharacterSingleOption WHERE nCharNo = @nCharNo) BEGIN
		INSERT INTO tCharacterSingleOption (nCharNo)
		VALUES (@nCharNo)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -3192
			RETURN
		END
	END

	SELECT bMysteryBolt
		,  bAutoPick 
		,  dRestExpLastExecuteTime
		,  nChatBan
	  FROM tCharacterSingleOption 
	 WHERE nCharNo = @nCharNo 		
END

0䅲Ĳ
GO
