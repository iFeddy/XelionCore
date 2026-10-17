CREATE PROCEDURE [dbo].[p_Char_GetStateAllWithRestTime]
	@nCharNo int
AS
BEGIN
	SET NOCOUNT ON
	SELECT 
		nStateNo, 
		nStateLevel, 
		nStateTime,
		nRestTime,
		nStateCount,
		nStatePowerDemage,
		nStatePowerSP,
		nStatePowerKeepTime,
		nStatePowerCoolTime,
		nCasterLevel,
		nCasterAttackPower 		
	 FROM tCharacterState 
	WHERE nCharNo = @nCharNo
END

0䤁哿Ş
GO
