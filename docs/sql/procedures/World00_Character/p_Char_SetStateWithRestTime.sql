CREATE PROCEDURE [dbo].[p_Char_SetStateWithRestTime]
	@nCharNo int,
	@nStateNo smallint,
	@nStateLevel smallint,
	@nStateTime int,
	@nRestTime int,
	@nStateCount smallint,
	@nStatePowerDemage tinyint,
	@nStatePowerSP tinyint,
	@nStatePowerKeepTime tinyint,
	@nStatePowerCoolTime tinyint,
	@nCasterLevel smallint,
	@nCasterAttackPower int,	
	@nRet int OUTPUT
AS
BEGIN
	SET NOCOUNT ON

	IF EXISTS(SELECT nCharNo FROM tCharacterState WHERE nCharNo = @nCharNo AND nStateNo = @nStateNo ) BEGIN
		UPDATE tCharacterState SET 
			nStateLevel			= @nStateLevel,
			nStateTime			= @nStateTime,
			nRestTime			= @nRestTime,
			nStateCount			= @nStateCount,	
			nStatePowerDemage	= @nStatePowerDemage,
			nStatePowerSP		= @nStatePowerSP,
			nStatePowerKeepTime	= @nStatePowerKeepTime,
			nStatePowerCoolTime	= @nStatePowerCoolTime,
		    nCasterLevel         = @nCasterLevel,
		    nCasterAttackPower   = @nCasterAttackPower 					
		WHERE nCharNo = @nCharNo AND nStateNo = @nStateNo
		
		SET @nRet = @@ERROR
	END
	ELSE BEGIN
		INSERT tCharacterState (  nCharNo,  nStateNo,  nStateLevel,  nStateTime, nRestTime,  nStateCount,  nStatePowerDemage,  nStatePowerSP,  nStatePowerKeepTime,  nStatePowerCoolTime, nCasterLevel, nCasterAttackPower )
		VALUES                 ( @nCharNo, @nStateNo, @nStateLevel, @nStateTime, @nRestTime, @nStateCount, @nStatePowerDemage, @nStatePowerSP, @nStatePowerKeepTime, @nStatePowerCoolTime, @nCasterLevel, @nCasterAttackPower )

		SET @nRet = @@ERROR
	END
END

0䴁ꨗĀ
GO
