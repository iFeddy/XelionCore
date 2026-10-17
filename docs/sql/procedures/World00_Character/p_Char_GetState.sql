CREATE PROCEDURE [dbo].[p_Char_GetState]
	@nCharNo int,
	@nStateNo smallint,
	@nStateLevel smallint OUTPUT,
	@nStateTime int OUTPUT,
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
		@nStateLevel         = nStateLevel, 
		@nStateTime          = nStateTime,
		@nStateCount         = nStateCount,	
		@nStatePowerDemage   = nStatePowerDemage,
		@nStatePowerSP       = nStatePowerSP,
		@nStatePowerKeepTime = nStatePowerKeepTime,
		@nStatePowerCoolTime = nStatePowerCoolTime,
		@nCasterLevel        = nCasterLevel,
		@nCasterAttackPower  = @nCasterAttackPower 
	  FROM tCharacterState 
	 WHERE nCharNo = @nCharNo AND nStateNo = @nStateNo
	SET @nRet = @@ROWCOUNT
END

0휁沶Ŝ
GO
