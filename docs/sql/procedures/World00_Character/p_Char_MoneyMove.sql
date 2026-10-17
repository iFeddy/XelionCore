CREATE PROCEDURE [dbo].[p_Char_MoneyMove]
/*
Move money
2008.2 By CJC
 input:
	see SQL
 output:
	@nRet	1 = OK
	        0 = Error
*/
	@nCharNoFrom int,
	@nCharNoTo   int,
	@nMoveMoney  bigint,
	@nRet        tinyint OUTPUT
AS
BEGIN
	SET NOCOUNT ON

	IF @nMoveMoney < 0
		BEGIN
			SET @nRet = 0
			RETURN
		END

	IF @nMoveMoney = 0
		BEGIN
			SET @nRet = 1
			RETURN
		END
		
	DECLARE @money as bigint

	BEGIN TRAN

	IF NOT EXISTS(SELECT nCharNo FROM tCharacter WHERE nCharNo = @nCharNoFrom AND bDeleted = 0)
	BEGIN
		SET @nRet = 0
		ROLLBACK TRAN
		RETURN
	END

	IF NOT EXISTS(SELECT nCharNo FROM tCharacter WHERE nCharNo = @nCharNoTo AND bDeleted = 0)
	BEGIN
		SET @nRet = 0
		ROLLBACK TRAN
		RETURN
	END

	SELECT @money = nMoney FROM tCharacter WHERE nCharNo = @nCharNoFrom AND bDeleted = 0
	IF @money >= @nMoveMoney
		BEGIN
			-- SUB MONEY
			UPDATE tCharacter SET nMoney = @money - @nMoveMoney WHERE nCharNo = @nCharNoFrom AND bDeleted = 0
			IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
			BEGIN
				SET @nRet =0
				ROLLBACK TRAN
				RETURN
			END
			
			-- ADD MONEY
			UPDATE tCharacter SET nMoney = nMoney + @nMoveMoney WHERE nCharNo = @nCharNoTo AND bDeleted = 0
			IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
			BEGIN
				SET @nRet = 0
				ROLLBACK TRAN
				RETURN
			END
	
			SET @nRet = 1
			COMMIT TRAN
		END
	ELSE
		BEGIN
			SET @nRet = 0
			ROLLBACK TRAN
			RETURN
		END

	SET NOCOUNT OFF
END

0
GO
