CREATE PROCEDURE [dbo].[p_Char_SetStateWithSetItem]
	@nCharNo int,
	@nStateNo smallint,
	@nStateLevel smallint,
	@nStateTime int,
	@nStateCount smallint,
	@nStatePowerDemage tinyint,
	@nStatePowerSP tinyint,
	@nStatePowerKeepTime tinyint,
	@nStatePowerCoolTime tinyint,
	@nSetItemIndexA smallint,
	@nSetItemValueA smallint,
	@nSetItemIndexB smallint,
	@nSetItemValueB smallint,
	@nSetItemIndexC smallint,
	@nSetItemValueC smallint,
	@nSetItemIndexd smallint,
	@nSetItemValueD smallint,
	@nCasterLevel smallint,
	@nCasterAttackPower int,
	@nRet int = 0 OUTPUT
AS
BEGIN
	SET NOCOUNT ON
	
	IF EXISTS(SELECT nCharNo FROM tCharacterState WHERE nCharNo = @nCharNo AND nStateNo = @nStateNo) BEGIN
		UPDATE tCharacterState 
		   SET nStateLevel			= @nStateLevel
		     , nStateTime			= @nStateTime
		     , nStateCount			= @nStateCount
		     , nStatePowerDemage	= @nStatePowerDemage
		     , nStatePowerSP		= @nStatePowerSP
		     , nStatePowerKeepTime	= @nStatePowerKeepTime
		     , nStatePowerCoolTime	= @nStatePowerCoolTime
		     , nSetItemIndexA		= @nSetItemIndexA
		     , nSetItemValueA		= @nSetItemValueA
		     , nSetItemIndexB		= @nSetItemIndexB
		     , nSetItemValueB		= @nSetItemValueB
		     , nSetItemIndexC		= @nSetItemIndexC
		     , nSetItemValueC		= @nSetItemValueC
		     , nSetItemIndexD		= @nSetItemIndexD
		     , nSetItemValueD		= @nSetItemValueD
		     , nCasterLevel         = @nCasterLevel 
		     , nCasterAttackPower   = @nCasterAttackPower 
		 WHERE nCharNo = @nCharNo AND nStateNo = @nStateNo

		SET @nRet = @@ERROR
	END
	ELSE BEGIN
		INSERT tCharacterState (nCharNo, nStateNo, nStateLevel, nStateTime, nStateCount, nStatePowerDemage, nStatePowerSP, nStatePowerKeepTime, nStatePowerCoolTime, nSetItemIndexA, nSetItemValueA, nSetItemIndexB, nSetItemValueB, nSetItemIndexC, nSetItemValueC, nSetItemIndexD, nSetItemValueD, nCasterLevel, nCasterAttackPower)
		VALUES (@nCharNo, @nStateNo, @nStateLevel, @nStateTime, @nStateCount, @nStatePowerDemage, @nStatePowerSP, @nStatePowerKeepTime, @nStatePowerCoolTime, @nSetItemIndexA, @nSetItemValueA, @nSetItemIndexB, @nSetItemValueB, @nSetItemIndexC, @nSetItemValueC, @nSetItemIndexD, @nSetItemValueD, @nCasterLevel, @nCasterAttackPower)

		SET @nRet = @@ERROR
	END
END

0뀁퍅ā
GO
