CREATE PROCEDURE [dbo].[p_Char_GetStateAll]
	@nCharNo int
AS
BEGIN
	SET NOCOUNT ON
	SELECT 
		nStateNo, 
		nStateLevel, 
		nStateTime,
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

0㨁闥ŝ
GO
