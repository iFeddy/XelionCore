CREATE   PROCEDURE [dbo].[usp_HolyPromise_setMoney]
/*
cen use reward
2009.4 By Luunar
 input:
	see SQL
 output:
	= 0 No Error
	< 0 Error
하위멤버가 돈을 사용하면 상위멤버에게 돈을 지급
*/

@nCharNo		int,
@nRewardMoney	bigint,
@nLevelGap		int,
@nRet			int OUTPUT

AS
BEGIN
	SET NOCOUNT ON

	DECLARE @nUpMemberCharNo	int
	DECLARE @nUpMemberLevel	int
	DECLARE @nMyLevel		int
	DECLARE @nType          tinyint

	SET @nUpMemberCharNo	= 0
	SET @nUpMemberLevel		= 0
	SET @nMyLevel		= 0
	SET @nType = 1 -- (1:입금/2:출금)

	IF ISNULL(@nCharNo, 0) <= 0 BEGIN
		SET @nRet = -2211
		RETURN
	END

	IF ISNULL(@nRewardMoney, 0) <= 0 BEGIN
		SET @nRet = -2212
		RETURN
	END


	IF EXISTS(SELECT nUpCharNo FROM tHolyPromise WHERE nCharNo = @nCharNo)
	BEGIN
		-- 상위 케릭터 번호 및 레벨을 가져와 본다.
		SELECT @nUpMemberCharNo = P.nUpCharNo, @nUpMemberLevel = C.nLevel FROM
			tHolyPromise AS P INNER JOIN tCharacter AS C
			ON C.nCharNo =  P.nUpCharNo
		WHERE P.nCharNo = @nCharNo AND C.bDeleted = 0

		-- 상위가 실존하면
		IF @nUpMemberCharNo > 0 and @nUpMemberLevel > 0
		BEGIN
			-- 나의 하위 맴버(나)의 Level 을 가져온다.
			SELECT @nMyLevel = nLevel FROM tCharacter WHERE nCharNo = @nCharNo
			-- 하위와 상위 맴버간 Level 차이가 일정량 이상을 벗어나면 보상해 주지 않는다.
			IF @nUpMemberLevel - @nMyLevel < @nLevelGap
			BEGIN
				SET @nRet = -2223
				RETURN
			END

			-- 상위 맴버에게 돈 보상 지급
			DECLARE @nSeqNo		int
			DECLARE @nRemainMoney	bigint

			SET @nSeqNo       = 0
			SET @nRemainMoney = 0

			BEGIN TRAN

			SELECT TOP 1 @nSeqNo = nSeqNo, @nRemainMoney = nRemainMoney 
			  FROM tHolyPromiseReward 
			 WHERE nCharNo = @nUpMemberCharNo 
			 ORDER BY nSeqNo DESC

			SET @nSeqNo       = @nSeqNo + 1
			SET @nRemainMoney = @nRemainMoney + @nRewardMoney

			INSERT INTO tHolyPromiseReward( nCharNo, nSeqNo, nType, nMoney, nRemainMoney, nReceiptCharNo )
			VALUES (@nUpMemberCharNo, @nSeqNo, @nType, @nRewardMoney, @nRemainMoney, @nCharNo)
			IF @@ERROR <> 0 BEGIN
				-- 파일 에러처리
				SET @nRet = -2292
				ROLLBACK TRAN
				RETURN
			END
			ELSE BEGIN
				SET @nRet = 0
				COMMIT TRAN
				RETURN			
			END
		END
		ELSE
		BEGIN
			-- 상위 맴버가 존재하지 않는다.
			SET @nRet = -2222
			RETURN
		END
	END
	ELSE
	BEGIN
		-- 홀리프라미스에 가입한 적이 없다. (큰 의미는 없음. 어차피 다음 단계에서 걸러짐)
		SET @nRet = -2221
		RETURN
	END
END

0숁癪
GO
