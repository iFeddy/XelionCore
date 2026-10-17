CREATE PROCEDURE [dbo].[p_Guild_MoneyAdd]

/*
Add Guild Money
2006.2 By CJC

Input:
	See SQL
	
Output:
	nRet = 0,  Error(Not found guild/member, not enough member's money)
	nRet = 1,  OK
*/

@nNo int,
@nMemberCharNo int,
@nMoneyAdd bigint,

@nMemberMoney bigint output,
@nTotalMoney bigint output,
@nRet tinyint output

AS

SET NOCOUNT ON

SET @nRet = 0

-- 길원의 금액을 얻는다
EXEC p_Char_MoneyGet @nMemberCharNo, @nMemberMoney OUTPUT, @nRet OUTPUT
IF @nRet <> 1 
	RETURN

-- 길원의 보유 금액이 요청한 금액보다 적으면 실패처리 한다.
IF @nMoneyAdd > @nMemberMoney
	RETURN


-- 길드 금액처리
BEGIN TRAN

-- 1. 길원 금액 빼기
EXEC p_Char_MoneySub @nMemberCharNo, @nMoneyAdd, @nRet OUTPUT
IF @nRet <> 1
BEGIN
	ROLLBACK TRAN
	SET @nRet = 0
	RETURN
END

-- 2. 길원의 남은 금액 얻기
EXEC p_Char_MoneyGet @nMemberCharNo, @nMemberMoney OUTPUT, @nRet OUTPUT
IF @nRet <> 1
BEGIN
	ROLLBACK TRAN
	SET @nRet = 0
	RETURN
END

-- 3. 길드 금액 증가
UPDATE tGuild SET nMoney = nMoney + @nMoneyAdd WHERE nNo = @nNo
SET @nRet = @@ROWCOUNT

IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
BEGIN
	ROLLBACK TRAN
	SET @nRet = 0
	RETURN
END

-- 4. 길드 최종급액 얻기
SELECT @nTotalMoney = nMoney FROM tGuild WHERE nNo = @nNo
IF @@ROWCOUNT <> 1
BEGIN
	ROLLBACK TRAN
	SET @nRet = 0
	RETURN
END

-- 완료
SET @nRet = 1
COMMIT TRAN

-- end

0餁踴İ
GO
