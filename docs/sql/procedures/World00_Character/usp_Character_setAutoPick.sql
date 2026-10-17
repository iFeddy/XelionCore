CREATE PROCEDURE [dbo].[usp_Character_setAutoPick]
	@nCharNo   int
,	@bAutoPick bit
,	@nRet      int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000 
	
	SET @nRet = 0

	IF NOT EXISTS(SELECT nCharNo FROM tCharacterSingleOption WHERE nCharNo = @nCharNo) BEGIN
		INSERT INTO tCharacterSingleOption (nCharNo, bAutoPick)
		VALUES (@nCharNo, @bAutoPick)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -3192
			RETURN
		END		
	END
	ELSE BEGIN
		UPDATE tCharacterSingleOption 
		   SET bAutoPick = @bAutoPick 
		 WHERE nCharNo = @nCharNo
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -3193
			RETURN
		END		
	END
END

0！뾦3
GO
