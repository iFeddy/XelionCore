CREATE PROCEDURE [dbo].[usp_Character_getGameOption]
	@nCharNo int
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	SET @nRet = 0
	
	IF EXISTS(SELECT nCharNo FROM tCharacterGameOption WHERE nCharNo = @nCharNo) BEGIN
		SELECT T1.nOptionNo, ISNULL(T2.nValue, T1.nDefaultValue) AS nValue
		  FROM tGameOptionInit AS T1 LEFT OUTER JOIN (
				SELECT * FROM tCharacterGameOption WHERE nCharNo = @nCharNo) AS T2
					ON T1.nOptionNo = T2.nOptionNo
		 WHERE T1.bIsActive = 1
	END
	ELSE BEGIN
		SELECT nOptionNo, nDefaultValue AS nValue
		  FROM tGameOptionInit 
		 WHERE bIsActive = 1
	END
END

0렁엧.
GO
