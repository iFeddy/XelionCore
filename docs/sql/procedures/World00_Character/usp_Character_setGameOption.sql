CREATE PROCEDURE [dbo].[usp_Character_setGameOption]
	@nCharNo   int
,	@nOptionNo smallint	
,	@nValue    tinyint
,	@nRet      int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	SET @nRet = 0
	
	IF EXISTS(SELECT nCharNo FROM tCharacterGameOption WHERE nCharNo = @nCharNo AND nOptionNo = @nOptionNo) BEGIN
		UPDATE tCharacterGameOption 
		   SET nValue = @nValue
		 WHERE nCharNo = @nCharNo
		   AND nOptionNo = @nOptionNo 
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = -2792
			RETURN
		END		   
	END
	ELSE BEGIN
		INSERT INTO tCharacterGameOption (nCharNo, nOptionNo, nValue)
		VALUES (@nCharNo, @nOptionNo, @nValue)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -2792
			RETURN
		END		   		
	END		
END

0鴁⮾;
GO
