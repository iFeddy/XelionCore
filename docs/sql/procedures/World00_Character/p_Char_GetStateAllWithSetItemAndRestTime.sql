CREATE PROCEDURE [dbo].[p_Char_GetStateAllWithSetItemAndRestTime]
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

0묁㵇Š
GO
