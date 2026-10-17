CREATE PROCEDURE [dbo].[p_Char_GetStateWithRestTime]
	@nCharNo int,
	@nStateNo smallint,
	@nStateLevel smallint OUTPUT,
	@nStateTime int OUTPUT,
	@nRestTime int OUTPUT,
	@nStateCount smallint OUTPUT,
	@nStatePowerDemage tinyint OUTPUT,
	@nStatePowerSP tinyint OUTPUT,
	@nStatePowerKeepTime tinyint OUTPUT,
	@nStatePowerCoolTime tinyint OUTPUT,
	@nCasterLevel smallint output,
	@nCasterAttackPower int output,			
	@nRet int OUTPUT
AS
BEGIN
	SET NOCOUNT ON
	SELECT 
		@nStateLevel			= nStateLevel, 
		@nStateTime			= nStateTime,
		@nRestTime			= nRestTime,
		@nStateCount			= nStateCount,	
		@nStatePowerDemage		= nStatePowerDemage,
		@nStatePowerSP			= nStatePowerSP,
		@nStatePowerKeepTime	= nStatePowerKeepTime,
		@nStatePowerCoolTime	= nStatePowerCoolTime,
		@nCasterLevel        = nCasterLevel,
		@nCasterAttackPower  = @nCasterAttackPower 		
	  FROM tCharacterState 
	 WHERE nCharNo = @nCharNo AND nStateNo = @nStateNo

	SET @nRet = @@ROWCOUNT
END

0ḁ晶š
GO
