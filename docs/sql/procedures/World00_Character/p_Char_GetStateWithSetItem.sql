CREATE PROCEDURE [dbo].[p_Char_GetStateWithSetItem]
	@nCharNo int,
	@nStateNo smallint,
	@nStateLevel smallint OUTPUT,
	@nStateTime int OUTPUT,
	@nStateCount smallint OUTPUT,
	@nStatePowerDemage tinyint OUTPUT,
	@nStatePowerSP tinyint OUTPUT,
	@nStatePowerKeepTime tinyint OUTPUT,
	@nStatePowerCoolTime tinyint OUTPUT,
	@nSetItemIndexA smallint OUTPUT,
	@nSetItemValueA smallint OUTPUT,
	@nSetItemIndexB smallint OUTPUT,
	@nSetItemValueB smallint OUTPUT,
	@nSetItemIndexC smallint OUTPUT,
	@nSetItemValueC smallint OUTPUT,
	@nSetItemIndexD smallint OUTPUT,
	@nSetItemValueD smallint OUTPUT,
	@nCasterLevel smallint output,
	@nCasterAttackPower int output,			
	@nRet int OUTPUT
AS
BEGIN
	SET NOCOUNT ON
	SELECT 
		@nStateLevel			= nStateLevel, 
		@nStateTime			= nStateTime,
		@nStateCount			= nStateCount,	
		@nStatePowerDemage		= nStatePowerDemage,
		@nStatePowerSP		= nStatePowerSP,
		@nStatePowerKeepTime	= nStatePowerKeepTime,
		@nStatePowerCoolTime		= nStatePowerCoolTime,
		@nSetItemIndexA		= nSetItemIndexA,
		@nSetItemValueA		= nSetItemValueA,
		@nSetItemIndexB		= nSetItemIndexB,
		@nSetItemValueB		= nSetItemValueB,
		@nSetItemIndexC		= nSetItemIndexC,
		@nSetItemValueC		= nSetItemValueC,
		@nSetItemIndexD		= nSetItemIndexD,
		@nSetItemValueD		= nSetItemValueD,
		@nCasterLevel        = nCasterLevel,
		@nCasterAttackPower  = @nCasterAttackPower 		
	  FROM tCharacterState 
 	 WHERE nCharNo = @nCharNo AND nStateNo = @nStateNo

	SET @nRet = @@ROWCOUNT
END

0圁媚Ţ
GO
