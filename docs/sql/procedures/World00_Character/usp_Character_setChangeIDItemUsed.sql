CREATE PROCEDURE [dbo].[usp_Character_setChangeIDItemUsed]
	@nCharNo int
,	@bInit   bit
,	@nUserNo int
,	@nItemKey bigint
,	@nItemOption smallint
,	@nItemLot smallint
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	DECLARE @nGrade AS tinyint
	SET @nRet = 0

	-- 아이템 소유자 확인
	EXEC p_Item_IsValidOwner @nItemKey, @nCharNo, @nRet OUTPUT
	IF @nRet <> 0 BEGIN
		SET @nRet = -3036
		RETURN
	END

	-- impossible guildmaster
	SELECT @nGrade = nGrade
	  FROM tCharacter T1 LEFT OUTER JOIN tGuildMember T2 ON T1.nCharNo = T2.nCharNo
	 WHERE T1.nCharNo = @nCharNo AND T1.bDeleted = 0
	   AND bDeleted = 0
	IF @nGrade IS NOT NULL AND @nGrade = 0 AND @bInit = 1 BEGIN
		SET @nRet = -3025
		RETURN
	END	   
	   
	IF EXISTS(SELECT nRowNo FROM tCharacterChangeID WHERE nCharNo = @nCharNo AND bFlag = 0) BEGIN
		SET @nRet = -3022
		RETURN
	END
	ELSE IF EXISTS(SELECT nCharNo FROM tCharacter WHERE nUserNo = @nUserNo AND bDeleted = 0) BEGIN
		BEGIN TRAN
		
		-- 아이템 사용 기록
		INSERT INTO tCharacterChangeID (nUserNo, nCharNo, bInit)
		VALUES (@nUserNo, @nCharNo, @bInit)
		IF @@ERROR <> 0 OR @@IDENTITY = 0 BEGIN
			SET @nRet = -3092
			ROLLBACK TRAN
			RETURN
		END	
		
		-- 아이템 삭제 (사용후 Lot = 0 인 경우)
		IF @nItemLot = 0 BEGIN
			EXEC p_Item_Delete @nItemKey, @nRet OUTPUT
			IF @nRet <> 0 BEGIN
				ROLLBACK TRAN
				SET @nRet = -3037
				RETURN
			END
		END
		-- 아이템 차감하기 (사용후 Lot > 0 인 경우)
		ELSE BEGIN
			EXEC p_Item_SetOption @nItemKey, @nItemOption, @nItemLot, @nRet OUTPUT
			IF @nRet <> 0 BEGIN
				ROLLBACK TRAN
				SET @nRet = -3038
				RETURN
			END
		END
		
		COMMIT TRAN		
	END
	ELSE BEGIN
		SET @nRet = -3023
		RETURN	
	END
	
	SET NOCOUNT OFF
END

0뤁嬭7
GO
