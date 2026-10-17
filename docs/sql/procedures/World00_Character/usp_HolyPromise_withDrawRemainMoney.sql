CREATE   PROCEDURE [dbo].[usp_HolyPromise_withDrawRemainMoney]
/*
withdraw cen reward
2009.4 By Luunar
 input:
	see SQL
 output:
	= 0 No Error
	< 0 Error
*/

@nCharNo		int,
@nRet			int OUTPUT

AS
BEGIN
	SET NOCOUNT ON

	DECLARE @nType tinyint
	SET @nType = 2 -- (1:입금/2:출금)

	IF ISNULL(@nCharNo, 0) <= 0 BEGIN
		SET @nRet = -2211
		RETURN
	END

	IF EXISTS(SELECT nMoney FROM tHolyPromiseReward WHERE nCharNo = @nCharNo)
	BEGIN
		-- 트랜잭션 시작
		BEGIN TRAN
			DECLARE @nRewardMoney	bigint
			DECLARE @nSeqNo		int

			-- 보상 금액 가져오기
			SELECT TOP 1 @nSeqNo = nSeqNo, @nRewardMoney = nRemainMoney FROM tHolyPromiseReward WHERE nCharNo = @nCharNo ORDER BY nSeqNo DESC

			-- 지급 금액이 0원보다 작거나 같은 경우엔 지급해 주지 않는다.
			IF @nRewardMoney <= 0
			BEGIN
				SET @nRet = -2226
				ROLLBACK TRAN
				RETURN
			END

			-- 보상 금액 지급
			EXEC p_Char_MoneyAdd @nCharNo, @nRewardMoney, @nRet OUTPUT
			IF @nRet <> 1
			BEGIN
				SET @nRet = -2227
				ROLLBACK TRAN
				RETURN
			END

			-- 보상 기록 기재
			SET @nSeqNo = @nSeqNo + 1

			INSERT INTO tHolyPromiseReward( nCharNo, nSeqNo, nType, nMoney, nRemainMoney, nReceiptCharNo )
			VALUES (@nCharNo, @nSeqNo, @nType, @nRewardMoney, 0, 0)
			IF @@ERROR <> 0
			BEGIN
				SET @nRet = -2292
				ROLLBACK TRAN
				RETURN
			END

			-- 결과
			SET @nRet = 0
			SELECT @nRewardMoney
		COMMIT TRAN
	END
	ELSE
	BEGIN
		-- 통장이 개설된 적이 없다.
		SET @nRet = -2225
		RETURN
	END
END

0ﬁ檎
GO
