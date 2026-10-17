-- ========================================================================
-- Stored Procedures Oficiales de Fiesta Online - Base de Datos: StatisticsData
-- Extraídos directamente de StatisticsData.bak
-- Total procedimientos recuperados: 18 de 18
-- ========================================================================

-- ========================================================================
-- Procedure: pFindCCU
-- ========================================================================
create procedure [dbo].[pFindCCU]
/*written by amadeus
*/
@start datetime,
@end	datetime
as
select dDate,w0,w1,w0+w1 total from (
select convert(char(16),dDate,120) dDate ,max(case nWorld when 0 then nCount else 0 end) w0, max(case nWorld when 1 then nCount else 0 end) w1
from tSameTimeChar where dDate between @start and @end
group by convert(char(16),dDate,120)
) as t1
GO

-- ========================================================================
-- Procedure: p_DayClassLevelPlay_Add
-- ========================================================================
CREATE PROCEDURE [dbo].[p_DayClassLevelPlay_Add] 
/*
 DayClassLevelPlay add record
 input:
  See SQL
 output: 
  nIndex = 추가된 번호(에러시 = 0 )
*/
@nWorld	tinyint,
@nClass	tinyint,
@nLevel smallint,
@nPlayMin		int,
@nConnectCount	int,
@nIndex		int OUTPUT
AS
SET NOCOUNT ON
SET @nIndex = 0
INSERT tDayClassLevelPlay (  nWorld,  nClass,  nLevel,  nPlayMin,  nConnectCount )
VALUES                    ( @nWorld, @nClass, @nLevel, @nPlayMin, @nConnectCount)
SET @nIndex = @@IDENTITY
-- end

0␁쇡Ġ
GO

-- ========================================================================
-- Procedure: p_DayClassLevelPlay_Query
-- ========================================================================
CREATE PROCEDURE [dbo].[p_DayClassLevelPlay_Query] 
/*
  DayClassLevelPlay Query - 주어진 조건에 맞는 데이타를 모두검색해서 얻는다.
  input:
   - 주어진 조건절
  output:
   - Select된 데이타
*/
@WhereExpr nvarchar(256)
AS
SET NOCOUNT ON
DECLARE @SQLCMD nvarchar(512)
SET  @SQLCMD = N'SELECT nIndex, dDate, nWorld, nClass, nLevel, nPlayMin, nConnectCount FROM tDayClassLevelPlay WHERE ' + @WhereExpr
EXEC sp_executesql @SQLCMD
-- end

0㴁ㅾă
GO

-- ========================================================================
-- Procedure: p_DayClassLevel_Add
-- ========================================================================
CREATE PROCEDURE [dbo].[p_DayClassLevel_Add] 
/*
 DayClassLevel add record
 input:
  See SQL
 output: 
  nIndex = 추가된 번호(에러시 = 0 )
*/
@nWorld	tinyint,
@nClass	tinyint,
@nLevel smallint,
@nCount	int,
@nIndex		int OUTPUT
AS
SET NOCOUNT ON
SET @nIndex = 0
INSERT tDayClassLevel (  nWorld,  nClass,  nLevel,  nCount )
VALUES                ( @nWorld, @nClass, @nLevel, @nCount )
SET @nIndex = @@IDENTITY
-- end

0춼ğ
GO

-- ========================================================================
-- Procedure: p_DayClassLevel_Query
-- ========================================================================
CREATE PROCEDURE [dbo].[p_DayClassLevel_Query] 
/*
  DayClassLevel Query - 주어진 조건에 맞는 데이타를 모두검색해서 얻는다.
  input:
   - 주어진 조건절
  output:
   - Select된 데이타
*/
@WhereExpr nvarchar(256)
AS
SET NOCOUNT ON
DECLARE @SQLCMD nvarchar(512)
SET  @SQLCMD = N'SELECT nIndex, dDate, nWorld, nClass, nLevel, nCount FROM tDayClassLevel WHERE ' + @WhereExpr
EXEC sp_executesql @SQLCMD
-- end

0瘁▢Ą
GO

-- ========================================================================
-- Procedure: p_DayConnectionTime
-- ========================================================================
CREATE PROCEDURE [dbo].[p_DayConnectionTime] 
AS
SET NOCOUNT ON
declare @nMon int
IF DAY(GETDATE()) = 1
set @nMon =MONTH(GETDATE())-1 
ELSE 
set @nMon = MONTH(GETDATE()) 
select min(nPlayMin) nMin, max(nPlayMin) nMax, avg(nPlayMin) nAvg, nWorldNo, DAY(dDate) nDay 
from [AccountLog].[dbo].[tAccountLog]			
where MONTH(dDate) = @nMon and nType = 2 
group by nWorldNo, DAY(dDate) 
order by nWorldNo, DAY(dDate) 
return 

0
GO

-- ========================================================================
-- Procedure: p_DayUser
-- ========================================================================
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

-- ========================================================================
-- Procedure: p_DayUser_Add
-- ========================================================================
CREATE PROCEDURE [dbo].[p_DayUser_Add] 
/*
 DayUser add
 input:
  See SQL
 output: 
  nIndex = 추가된 번호(에러시 = 0 )
*/
@nUserCount		int,
@nConnectCount	int,
@nIndex		int OUTPUT
AS
SET NOCOUNT ON
SET @nIndex = 0
INSERT tDayUser (  nUserCount,  nConnectCount )
VALUES          ( @nUserCount, @nConnectCount )
SET @nIndex = @@IDENTITY
-- end

0崁똅ġ
GO

-- ========================================================================
-- Procedure: p_DayUser_Query
-- ========================================================================
CREATE PROCEDURE [dbo].[p_DayUser_Query] 
/*
  DayUser Query - 주어진 조건에 맞는 데이타를 모두검색해서 얻는다.
  input:
   - 주어진 조건절
  output:
   - Select된 데이타
*/
@WhereExpr nvarchar(256)
AS
SET NOCOUNT ON
DECLARE @SQLCMD nvarchar(512)
SET  @SQLCMD = N'SELECT nIndex, dDate, nUserCount, nConnectCount FROM tDayUser WHERE ' + @WhereExpr
EXEC sp_executesql @SQLCMD
-- end

0෪Ć
GO

-- ========================================================================
-- Procedure: p_SameTimeChar_Add
-- ========================================================================
CREATE PROCEDURE [dbo].[p_SameTimeChar_Add] 
/*
 SameTimeChar add
 input:
  See SQL
 output: 
  nIndex = 추가된 번호(에러시 = 0 )
*/
@nWorld	tinyint,
@nCount		int,
@nIndex		int OUTPUT
AS
SET NOCOUNT ON
---
SET @nIndex = 0
 
INSERT tSameTimeChar (  nWorld,  nCount )
VALUES               ( @nWorld, @nCount )
SET @nIndex = @@IDENTITY
-- end

0܁ﴬě
GO

-- ========================================================================
-- Procedure: p_SameTimeChar_Query
-- ========================================================================
CREATE PROCEDURE [dbo].[p_SameTimeChar_Query] 
/*
  SameTimeChar Query - 주어진 조건에 맞는 데이타를 모두검색해서 얻는다.
  input:
   - 주어진 조건절
  output:
   - Select된 데이타
*/
@WhereExpr nvarchar(256)
AS
SET NOCOUNT ON
DECLARE @SQLCMD nvarchar(512)
SET  @SQLCMD = N'SELECT nIndex, dDate, nWorld, nCount FROM tSameTimeChar WHERE ' + @WhereExpr
EXEC sp_executesql @SQLCMD
-- end

0℁ȏć
GO

-- ========================================================================
-- Procedure: p_SameTimeWorld
-- ========================================================================
CREATE  PROCEDURE [dbo].[p_SameTimeWorld]
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000
declare 
@CurrentTime datetime,
@TimeText nvarchar(24)
SELECT TOP 1 @TimeText= convert(nvarchar(24),dDate,0), @CurrentTime=dDate  FROM tSameTimeChar WITH(READUNCOMMITTED)
WHERE dDate >= DATEADD(minute, -5, GETDATE() )
order by dDate desc
IF @@ROWCOUNT = 0 
	BEGIN
		SELECT GETDATE() AS CurTime, 0 AS CurCount
	END
ELSE
	BEGIN
		SELECT @CurrentTime AS CurTime , SUM(nCount) AS CurCount  FROM tSameTimeChar WITH(READUNCOMMITTED)
		WHERE convert(nvarchar(24),dDate,0) = @TimeText
	END
END

0䀁
GO

-- ========================================================================
-- Procedure: p_UserAgeReport
-- ========================================================================
CREATE PROCEDURE [dbo].[p_UserAgeReport]
	@S_DATE NVARCHAR(10) = '',
	@E_DATE NVARCHAR(10) = ''
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000
	DECLARE @S_YEAR AS INT
	DECLARE @S_MONTH AS INT
	DECLARE @S_DAY AS INT
	DECLARE @E_YEAR AS INT
	DECLARE @E_MONTH AS INT
	DECLARE @E_DAY AS INT
	IF @S_DATE = ''
		SET @S_DATE = DATEADD(dd, -1, GETDATE())
	IF @E_DATE = ''
		SET @E_DATE = DATEADD(dd, -1, GETDATE())
	SET @S_YEAR  = DATENAME(yyyy, @S_DATE)
	SET @S_MONTH = DATENAME(mm, @S_DATE)
	SET @S_DAY   = DATENAME(dd, @S_DATE)
	SET @E_YEAR  = DATENAME(yyyy, @E_DATE)
	SET @E_MONTH = DATENAME(mm, @E_DATE)
	SET @E_DAY   = DATENAME(dd, @E_DATE)
	PRINT '시작일 : ' + @S_DATE
	PRINT '종료일 : ' + @E_DATE
	PRINT @S_MONTH
	PRINT @E_MONTH
	SELECT CASE(nGender)
				WHEN 0 THEN '여자'
				WHEN 1 THEN '남자'
		   END AS [성별], nCount AS [인원]
	  FROM tGender(NOLOCK)
	 WHERE nYear >= @S_YEAR AND nMonth >= @S_MONTH AND nDay >= @S_DAY
	   AND nYear <= @E_YEAR AND nMonth <= @E_MONTH AND nDay <= @E_DAY
	SET LOCK_TIMEOUT -1
	SET NOCOUNT OFF
END
	

0츁इě
GO

-- ========================================================================
-- Procedure: p_UserAge_Query
-- ========================================================================
CREATE PROCEDURE [dbo].[p_UserAge_Query] 
/*
  UserAge Query - 주어진 조건에 맞는 데이타를 모두검색해서 얻는다.
  input:
   - 주어진 조건절
  output:
   - Select된 데이타
*/
@WhereExpr nvarchar(256)
AS
SET NOCOUNT ON
DECLARE @SQLCMD nvarchar(512)
SET  @SQLCMD = N'SELECT nLogNo, nYear, nMonth, nDay, nWeek, nWeekDay, nAge, nCount FROM tAge WHERE ' + @WhereExpr
EXEC sp_executesql @SQLCMD
-- end

0威ć
GO

-- ========================================================================
-- Procedure: p_UserGenderTest
-- ========================================================================
CREATE  PROCEDURE [dbo].[p_UserGenderTest]
	@DATE	NVARCHAR(10) = ''
AS
BEGIN 
	DECLARE @YEAR AS INT
	DECLARE @MONTH AS INT
	DECLARE @DAY AS INT
	DECLARE @LOG_DATE AS NVARCHAR(10)
	DECLARE @BASE_DATE AS NVARCHAR(10)
	IF @DATE = ''
		SET @LOG_DATE = CONVERT(NVARCHAR(10), DATEADD(dd, -1, GETDATE()), 120)
	ELSE
		SET @LOG_DATE = @DATE
	PRINT @LOG_DATE
	PRINT '--------------'
	SET @YEAR  = DATENAME(yyyy, @LOG_DATE)
	SET @MONTH = DATENAME(mm, @LOG_DATE)
	SET @DAY   = DATENAME(dd, @LOG_DATE)
	SET @BASE_DATE = CONVERT(NVARCHAR(10), DATEADD(dd, 1, @LOG_DATE), 120)
	PRINT @BASE_DATE
	PRINT '--------------'
	PRINT @YEAR
	PRINT '--------------'
END
EXEC p_UserGender
GO

-- ========================================================================
-- Procedure: p_UserGender_Query
-- ========================================================================
CREATE PROCEDURE [dbo].[p_UserGender_Query] 
/*
  UserGender Query - 주어진 조건에 맞는 데이타를 모두검색해서 얻는다.
  input:
   - 주어진 조건절
  output:
   - Select된 데이타
*/
@WhereExpr nvarchar(256)
AS
SET NOCOUNT ON
DECLARE @SQLCMD nvarchar(512)
SET  @SQLCMD = N'SELECT nLogNo, nYear, nMonth, nDay, nWeek, nWeekDay, nGender, nCount FROM tGender WHERE ' + @WhereExpr
EXEC sp_executesql @SQLCMD
-- end

0찁ĉ
GO

-- ========================================================================
-- Procedure: p_UserLocation_Query
-- ========================================================================
CREATE PROCEDURE [dbo].[p_UserLocation_Query] 
/*
  UserLocation Query - 주어진 조건에 맞는 데이타를 모두검색해서 얻는다.
  input:
   - 주어진 조건절
  output:
   - Select된 데이타
*/
@WhereExpr nvarchar(256)
AS
SET NOCOUNT ON
DECLARE @SQLCMD nvarchar(512)
SET  @SQLCMD = N'SELECT nLogNo, nYear, nMonth, nDay, nWeek, nWeekDay, sLocation, nCount FROM tLocation WHERE ' + @WhereExpr
EXEC sp_executesql @SQLCMD
-- end

0嬁譹
GO

-- ========================================================================
-- Procedure: p_getQuestTimes
-- ========================================================================
CREATE procedure [dbo].[p_getQuestTimes]
@nQuestNo int,
@top smallint
as
	set nocount on;
	select * from (
		select top (@top) 0 ServerNo, c.sID,nTimes, convert(char(17),dLastComplete,113) dLastComplete from w00_Character.dbo.tQuestTimes q, w00_Character.dbo.tCharacter c
		where nQuestNo=@nQuestNo and q.nCharNo=c.nCharNo and c.bDeleted=0 
		order by nTimes desc, q.dLastComplete  
	) as t0
	union all
	select * from (	
		select top (@top) 1 ServerNo, c.sID,nTimes, convert(char(17),dLastComplete,113) dLastComplete from w01_Character.dbo.tQuestTimes q, w01_Character.dbo.tCharacter c
		where nQuestNo=@nQuestNo and q.nCharNo=c.nCharNo and c.bDeleted=0 
		order by nTimes desc, q.dLastComplete  
	) as t1
0
GO

