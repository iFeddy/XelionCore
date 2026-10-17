CREATE PROCEDURE [dbo].[p_DayUser]
/*
 Day User - 일일 개별 접속자, 접속횟수
 Desc:
	Update tDayUser Statstics table
	AccountLog.dbo.tAccountLog에서 마지막으로 로그인 날짜가 어제인 
	사람수, 접속횟수를 구해서 어제의 날짜로 테이블에 기록시킨다.
*/
AS
SET NOCOUNT ON
-- GameLog's nType for LOGOUT 
DECLARE @LoginID int	
SET		@LoginID = 1
-- 검사 시간을 얻는다.(어제의 날짜 범위)
DECLARE @Date		datetime
DECLARE @StartDate	datetime 
DECLARE @EndDate	datetime 
SET @Date = GETDATE()
SET @EndDate = STR( YEAR(@Date) ) + STR( MONTH(@Date) ) + STR( DAY(@Date) )
SET @StartDate = @EndDate - 1
-- 중복기록 방지를 위해 레코드 있는 경우 삭제(어제날짜)
DELETE tDayUser WHERE dDate = @StartDate
-- 접속자수, 접속횟수 얻기
DECLARE @nUserCount    int
DECLARE @nConnectCount int
SELECT @nUserCount = COUNT(DISTINCT(nUserNo)), @nConnectCount = COUNT(*) 
--FROM [ACCOUNT-LOG].AccountLog.dbo.tAccountLog 
FROM AccountLog.dbo.tAccountLog 
WHERE (nType = @LoginID) AND (dDate >= @StartDate) AND (dDate < @EndDate)
-- Table에 기록하기
DECLARE @Index int
exec p_DayUser_Add @nUserCount, @nConnectCount, @Index OUTPUT
IF @Index <> 0 UPDATE tDayUser SET dDate = @StartDate WHERE nIndex = @Index -- 날짜를 변경
-- end

<阁ꨩĢ
GO
