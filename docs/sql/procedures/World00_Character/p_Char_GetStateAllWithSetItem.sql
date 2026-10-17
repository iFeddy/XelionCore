CREATE PROCEDURE [dbo].[p_Char_GetStateAllWithSetItem]
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
		nSetItemIndexA,
		nSetItemValueA,
		nSetItemIndexB,
		nSetItemValueB,
		nSetItemIndexC,
		nSetItemValueC,
		nSetItemIndexD,
		nSetItemValueD,
		nCasterLevel,
		nCasterAttackPower		
	 FROM tCharacterState 
	WHERE nCharNo = @nCharNo
END

0舁䤣ş
GO
