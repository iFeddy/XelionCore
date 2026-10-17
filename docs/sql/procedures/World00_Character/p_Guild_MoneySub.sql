CREATE PROCEDURE [dbo].[p_Guild_MoneySub]

/*
Sub Guild Money
2007.11 By Lunar

Input:
	See SQL
	
Output:
	nRet =   0, OK

		1, Fail get Guild Money
		2, Guild Money Sub Failed(Not enough money)
		3, Guild Money Update Failed
		4, Member Money Update Failed
		5, Fail get Member Money
		6, Fail get Final Guild Money
*/

@nNo int,
@nMemberCharNo int,
@nMoneySub bigint,

@nMemberMoney bigint output,
@nTotalMoney bigint output,
@nRet tinyint output

AS

SET NOCOUNT ON

SET @nRet = 0
SET @nMemberMoney = 0
SET @nTotalMoney = 0

-- 길드의 금액을 얻는다
SELECT @nTotalMoney = nMoney FROM tGuild WHERE nNo = @nNo
IF @@ERROR <> 0
BEGIN
	SET @nRet = 1
	RETURN
END

-- 길드의 보유 금액이 요청한 금액보다 적으면 실패처리 한다.
IF @nMoneySub > @nTotalMoney
BEGIN
	SET @nRet = 2
	RETURN
END

-- 길드 금액처리
BEGIN TRAN

-- 1. 길드 금액 빼기
UPDATE tGuild SET nMoney = nMoney - @nMoneySub WHERE nNo = @nNo
SET @nRet = @@ROWCOUNT

IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
BEGIN
	ROLLBACK TRAN
	SET @nRet = 3
	RETURN
END

-- 2. 길원 금액 증가
EXEC p_Char_MoneyAdd @nMemberCharNo, @nMoneySub, @nRet OUTPUT
IF @nRet <> 1
BEGIN
	ROLLBACK TRAN
	SET @nRet = 4
	RETURN
END

-- 3. 길원 금액 얻기
EXEC p_Char_MoneyGet @nMemberCharNo, @nMemberMoney OUTPUT, @nRet OUTPUT
IF @nRet <> 1
BEGIN
	ROLLBACK TRAN
	SET @nRet = 5
	RETURN
END

-- 4. 길드 최종급액 얻기
SELECT @nTotalMoney = nMoney FROM tGuild WHERE nNo = @nNo
IF @@ROWCOUNT <> 1
BEGIN
	ROLLBACK TRAN
	SET @nRet = 6
	RETURN
END

-- 완료
SET @nRet = 0
COMMIT TRAN

-- end

0툁艘ı
GO
