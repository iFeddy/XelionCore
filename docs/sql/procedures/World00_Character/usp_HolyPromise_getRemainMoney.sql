CREATE  PROCEDURE [dbo].[usp_HolyPromise_getRemainMoney]
/*
get cen reward
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

	DECLARE @nRemainMoney bigint
	SET @nRemainMoney = 0

	IF ISNULL(@nCharNo, 0) <= 0 BEGIN
		SET @nRet = -2211
		RETURN
	END

	-- 잔액이 얼마나 있는지 가져온다.
	SELECT TOP 1 @nRemainMoney = nRemainMoney 
	  FROM tHolyPromiseReward 
	 WHERE nCharNo = @nCharNo 
	 ORDER BY nSeqNo DESC

	-- 만약 잔액이 0원 보다 많다면
	IF @nRemainMoney > 0
	BEGIN
		-- 잔액을 바로 출력해주고 끝
		SELECT @nRemainMoney
		SET @nRet = 0
		RETURN
	END
	ELSE
	-- 잔액이 0원일 경우엔 상위로써 등록되어 있는지 여부를 체크해야 함
	BEGIN
		-- 현재 상위로써 등록되어 있다면
		IF EXISTS(SELECT TOP 1 nUpCharNo FROM tHolyPromise WHERE nUpCharNo = @nCharNo)
		BEGIN
			-- 잔액을 출력해주고 끝(다만 잔액은 언제나 0원일 것이다.)
			SELECT 0
			SET @nRet = 0
			RETURN
		END
		ELSE
		-- 상위로써 등록되어 있지 않는 경우
		BEGIN
			-- 홀리 프라미스 설명창 띄우라고 명령
			SET @nRet = -2224
			RETURN
		END	
	END
END

0倁踢
GO
