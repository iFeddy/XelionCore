-- ========================================================================
-- Stored Procedures Oficiales de Fiesta Online - Base de Datos: World00_GameLog
-- Extraídos directamente de World00_GameLog.bak
-- Total procedimientos recuperados: 26 de 26
-- ========================================================================

-- ========================================================================
-- Procedure: p_GameLog_Add
-- ========================================================================
CREATE PROCEDURE [dbo].[p_GameLog_Add]
/*
 Game log add
 input:
  All Args = 로그에 기록될 정보
 output: 
  nLogNo = 추가된 로그의 번호(에러시 LogNo = 0 )
*/
@nLogNo bigint OUTPUT,
@nType int,
@nCharNo int,
@sMap varchar(16),
@nMapX int,
@nMapY int,
@nMapZ int,
@nTargetCharNo int,
@nTargetID int,
@nItemKey bigint,
@nInt1 int,
@nInt2 int,
@nInt3 int,
@nBigint1 bigint
AS
SET NOCOUNT ON
SET @nLogNo = 0
INSERT tGameLog (  nType,  nCharNo,  sMap,  nMapX,  nMapY,  nMapZ,  nTargetCharNo,  nTargetID,  nItemKey,  nInt1,  nInt2,  nInt3,  nBigint1 )
VALUES          ( @nType, @nCharNo, @sMap, @nMapX, @nMapY, @nMapZ, @nTargetCharNo, @nTargetID, @nItemKey, @nInt1, @nInt2, @nInt3, @nBigint1 )
SET @nLogNo = @@IDENTITY
-- end

0䌁넢Ķ
GO

-- ========================================================================
-- Procedure: p_GameLog_Add_Type1
-- ========================================================================
CREATE PROCEDURE [dbo].[p_GameLog_Add_Type1]
/*
 Game log add
 input:
  All Args = 로그에 기록될 정보
 output: 
  nLogNo = 추가된 로그의 번호(에러시 LogNo = 0 )
*/
@nLogNo bigint OUTPUT,
@nType int,
@nCharNo int,
@sMap varchar(16),
@nMapX int,
@nMapY int,
@nMapZ int
AS
SET NOCOUNT ON
SET @nLogNo = 0
INSERT tGameLog (  nType,  nCharNo,  sMap,  nMapX,  nMapY,  nMapZ )
VALUES          ( @nType, @nCharNo, @sMap, @nMapX, @nMapY, @nMapZ )
SET @nLogNo = @@IDENTITY
-- end

0簁ꕆķ
GO

-- ========================================================================
-- Procedure: p_GameLog_Add_Type2
-- ========================================================================
CREATE PROCEDURE [dbo].[p_GameLog_Add_Type2]
/*
 Game log add
 input:
  All Args = 로그에 기록될 정보
 output: 
  nLogNo = 추가된 로그의 번호(에러시 LogNo = 0 )
*/
@nLogNo bigint OUTPUT,
@nType int,
@nCharNo int,
@sMap varchar(16),
@nMapX int,
@nMapY int,
@nMapZ int,
@nInt1 int,
@nInt2 int
AS
SET NOCOUNT ON
SET @nLogNo = 0
INSERT tGameLog (  nType,  nCharNo,  sMap,  nMapX,  nMapY,  nMapZ,  nInt1,  nInt2 )
VALUES          ( @nType, @nCharNo, @sMap, @nMapX, @nMapY, @nMapZ, @nInt1, @nInt2 )
SET @nLogNo = @@IDENTITY
-- end

0딁饪ĸ
GO

-- ========================================================================
-- Procedure: p_GameLog_Add_Type3
-- ========================================================================
CREATE PROCEDURE [dbo].[p_GameLog_Add_Type3]
/*
 Game log add
 input:
  All Args = 로그에 기록될 정보
 output: 
  nLogNo = 추가된 로그의 번호(에러시 LogNo = 0 )
*/
@nLogNo bigint OUTPUT,
@nType int,
@nCharNo int,
@sMap varchar(16),
@nMapX int,
@nMapY int,
@nMapZ int,
@nTargetCharNo int,
@nTargetID int,
@nInt1 int
AS
SET NOCOUNT ON
SET @nLogNo = 0
INSERT tGameLog (  nType,  nCharNo,  sMap,  nMapX,  nMapY,  nMapZ,  nTargetCharNo,  nTargetID,  nInt1 )
VALUES          ( @nType, @nCharNo, @sMap, @nMapX, @nMapY, @nMapZ, @nTargetCharNo, @nTargetID, @nInt1 )
SET @nLogNo = @@IDENTITY
-- end

0趎Ĺ
GO

-- ========================================================================
-- Procedure: p_GameLog_Add_Type4
-- ========================================================================
CREATE PROCEDURE [dbo].[p_GameLog_Add_Type4]
/*
 Game log add
 input:
  All Args = 로그에 기록될 정보
 output: 
  nLogNo = 추가된 로그의 번호(에러시 LogNo = 0 )
*/
@nLogNo bigint OUTPUT,
@nType int,
@nCharNo int,
@sMap varchar(16),
@nMapX int,
@nMapY int,
@nMapZ int,
@nInt1 int,
@nInt2 int,
@nInt3 int
AS
SET NOCOUNT ON
SET @nLogNo = 0
INSERT tGameLog (  nType,  nCharNo,  sMap,  nMapX,  nMapY,  nMapZ,  nInt1,  nInt2,  nInt3 )
VALUES          ( @nType, @nCharNo, @sMap, @nMapX, @nMapY, @nMapZ, @nInt1, @nInt2, @nInt3 )
SET @nLogNo = @@IDENTITY
-- end
GO

-- ========================================================================
-- Procedure: p_GameLog_Add_Type5
-- ========================================================================
CREATE PROCEDURE [dbo].[p_GameLog_Add_Type5]
/*
 Game log add
 input:
  All Args = 로그에 기록될 정보
 output: 
  nLogNo = 추가된 로그의 번호(에러시 LogNo = 0 )
*/
@nLogNo bigint OUTPUT,
@nType int,
@nCharNo int,
@nInt1 int,
@nInt2 int
AS
SET NOCOUNT ON
SET @nLogNo = 0
INSERT tGameLog (  nType,  nCharNo,  nInt1,  nInt2 )
VALUES          ( @nType, @nCharNo, @nInt1, @nInt2 )
SET @nLogNo = @@IDENTITY
-- end

0态痗Ļ
GO

-- ========================================================================
-- Procedure: p_GameLog_Add_Type6
-- ========================================================================
CREATE PROCEDURE [dbo].[p_GameLog_Add_Type6]
/*
 Game log add
 input:
  All Args = 로그에 기록될 정보
 output: 
  nLogNo = 추가된 로그의 번호(에러시 LogNo = 0 )
*/
@nLogNo bigint OUTPUT,
@nType int,
@nCharNo int,
@sMap varchar(16),
@nMapX int,
@nMapY int,
@nMapZ int,
@nInt1 int,
@nBigint1 bigint
AS
SET NOCOUNT ON
SET @nLogNo = 0
INSERT tGameLog (  nType,  nCharNo,  sMap,  nMapX,  nMapY,  nMapZ,  nInt1,  nBigint1 )
VALUES          ( @nType, @nCharNo, @sMap, @nMapX, @nMapY, @nMapZ, @nInt1, @nBigint1 )
SET @nLogNo = @@IDENTITY
-- end

0餁槻ļ
GO

-- ========================================================================
-- Procedure: p_GameLog_Add_Type7
-- ========================================================================
CREATE PROCEDURE [dbo].[p_GameLog_Add_Type7]
/*
 Game log add
 input:
  All Args = 로그에 기록될 정보
 output: 
  nLogNo = 추가된 로그의 번호(에러시 LogNo = 0 )
*/
@nLogNo bigint OUTPUT,
@nType int,
@nCharNo int,
@sMap varchar(16),
@nMapX int,
@nMapY int,
@nMapZ int,
@nTargetCharNo int,
@nInt1 int,
@nInt2 int
AS
SET NOCOUNT ON
SET @nLogNo = 0
INSERT tGameLog (  nType,  nCharNo,  sMap,  nMapX,  nMapY,  nMapZ,  nTargetCharNo,  nInt1,  nInt2 )
VALUES          ( @nType, @nCharNo, @sMap, @nMapX, @nMapY, @nMapZ, @nTargetCharNo, @nInt1, @nInt2 )
SET @nLogNo = @@IDENTITY
-- end

0툁帟Ľ
GO

-- ========================================================================
-- Procedure: p_GameLog_Add_Type8
-- ========================================================================
CREATE PROCEDURE [dbo].[p_GameLog_Add_Type8]
/*
 Game log add
 input:
  All Args = 로그에 기록될 정보
 output: 
  nLogNo = 추가된 로그의 번호(에러시 LogNo = 0 )
*/
@nLogNo bigint OUTPUT,
@nType int,
@nCharNo int,
@sMap varchar(16),
@nMapX int,
@nMapY int,
@nMapZ int,
@nInt1 int
AS
SET NOCOUNT ON
SET @nLogNo = 0
INSERT tGameLog (  nType,  nCharNo,  sMap,  nMapX,  nMapY,  nMapZ,  nInt1 )
VALUES          ( @nType, @nCharNo, @sMap, @nMapX, @nMapY, @nMapZ, @nInt1 )
SET @nLogNo = @@IDENTITY
-- end

0ଁ剄ľ
GO

-- ========================================================================
-- Procedure: p_GameLog_Add_Type9
-- ========================================================================
CREATE PROCEDURE [dbo].[p_GameLog_Add_Type9]
/*
 Game log add
 input:
  All Args = 로그에 기록될 정보
 output: 
  nLogNo = 추가된 로그의 번호(에러시 LogNo = 0 )
*/
@nLogNo bigint OUTPUT,
@nType int,
@nCharNo int,
@sMap varchar(16),
@nMapX int,
@nMapY int,
@nMapZ int,
@nItemKey bigint,
@nInt1 int,
@nInt2 int
AS
SET NOCOUNT ON
SET @nLogNo = 0
INSERT tGameLog (  nType,  nCharNo,  sMap,  nMapX,  nMapY,  nMapZ,  nItemKey,  nInt1,  nInt2 )
VALUES          ( @nType, @nCharNo, @sMap, @nMapX, @nMapY, @nMapZ, @nItemKey, @nInt1, @nInt2 )
SET @nLogNo = @@IDENTITY
-- end
GO

-- ========================================================================
-- Procedure: p_GameLog_Add_TypeA
-- ========================================================================
CREATE PROCEDURE [dbo].[p_GameLog_Add_TypeA]
/*
 Game log add
 input:
  All Args = 로그에 기록될 정보
 output: 
  nLogNo = 추가된 로그의 번호(에러시 LogNo = 0 )
*/
@nLogNo bigint OUTPUT,
@nType int,
@nCharNo int,
@nTargetCharNo int,
@nTargetID int,
@nItemKey bigint,
@nInt1 int,
@nInt2 int,
@nInt3 int
AS
SET NOCOUNT ON
SET @nLogNo = 0
INSERT tGameLog (  nType,  nCharNo,  nTargetCharNo,  nTargetID,  nItemKey,  nInt1,  nInt2,  nInt3 )
VALUES          ( @nType, @nCharNo, @nTargetCharNo, @nTargetID, @nItemKey, @nInt1, @nInt2, @nInt3 )
SET @nLogNo = @@IDENTITY
-- end

0紁㪌ŀ
GO

-- ========================================================================
-- Procedure: p_GameLog_Add_TypeB
-- ========================================================================
CREATE PROCEDURE [dbo].[p_GameLog_Add_TypeB]
/*
 Game log add
 input:
  All Args = 로그에 기록될 정보
 output: 
  nLogNo = 추가된 로그의 번호(에러시 LogNo = 0 )
*/
@nLogNo bigint OUTPUT,
@nType int,
@nCharNo int,
@sMap varchar(16),
@nMapX int,
@nMapY int,
@nMapZ int,
@nTargetCharNo int,
@nTargetID int,
@nItemKey bigint,
@nInt1 int,
@nInt2 int,
@nInt3 int
AS
SET NOCOUNT ON
SET @nLogNo = 0
INSERT tGameLog (  nType,  nCharNo,  sMap,  nMapX,  nMapY,  nMapZ,  nTargetCharNo,  nTargetID,  nItemKey,  nInt1,  nInt2,  nInt3 )
VALUES          ( @nType, @nCharNo, @sMap, @nMapX, @nMapY, @nMapZ, @nTargetCharNo, @nTargetID, @nItemKey, @nInt1, @nInt2, @nInt3 )
SET @nLogNo = @@IDENTITY
-- end

0똁⺰Ł
GO

-- ========================================================================
-- Procedure: p_GameLog_Add_TypeC
-- ========================================================================
CREATE PROCEDURE [dbo].[p_GameLog_Add_TypeC]
/*
 Game log add
 input:
  All Args = 로그에 기록될 정보
 output: 
  nLogNo = 추가된 로그의 번호(에러시 LogNo = 0 )
*/
@nLogNo bigint OUTPUT,
@nType int,
@nCharNo int,
@sMap varchar(16),
@nMapX int,
@nMapY int,
@nMapZ int,
@nTargetID int,
@nItemKey bigint,
@nInt1 int,
@nInt2 int,
@nInt3 int
AS
SET NOCOUNT ON
SET @nLogNo = 0
INSERT tGameLog (  nType,  nCharNo,  sMap,  nMapX,  nMapY,  nMapZ,  nTargetID,  nItemKey,  nInt1,  nInt2,  nInt3 )
VALUES          ( @nType, @nCharNo, @sMap, @nMapX, @nMapY, @nMapZ, @nTargetID, @nItemKey, @nInt1, @nInt2, @nInt3 )
SET @nLogNo = @@IDENTITY
-- end

0⋔ł
GO

-- ========================================================================
-- Procedure: p_GameLog_Add_TypeD
-- ========================================================================
CREATE PROCEDURE [dbo].[p_GameLog_Add_TypeD]
/*
 Game log add
 input:
  All Args = 로그에 기록될 정보
 output: 
  nLogNo = 추가된 로그의 번호(에러시 LogNo = 0 )
*/
@nLogNo bigint OUTPUT,
@nType int,
@nCharNo int,
@nItemKey bigint,
@nInt1 int,
@nInt2 int
AS
SET NOCOUNT ON
SET @nLogNo = 0
INSERT tGameLog (  nType,  nCharNo,  nItemKey,  nInt1,  nInt2 )
VALUES          ( @nType, @nCharNo, @nItemKey, @nInt1, @nInt2 )
SET @nLogNo = @@IDENTITY
-- end

0⠁᛹Ń
GO

-- ========================================================================
-- Procedure: p_GameLog_Add_TypeE
-- ========================================================================
CREATE PROCEDURE [dbo].[p_GameLog_Add_TypeE]
/*
 Game log add
 input:
  All Args = 로그에 기록될 정보
 output: 
  nLogNo = 추가된 로그의 번호(에러시 LogNo = 0 )
*/
@nLogNo bigint OUTPUT,
@nType int,
@sMap varchar(16),
@nMapX int,
@nMapY int,
@nMapZ int,
@nTargetID int,
@nItemKey bigint,
@nInt1 int,
@nInt2 int
AS
SET NOCOUNT ON
SET @nLogNo = 0
INSERT tGameLog (  nType,  sMap,  nMapX,  nMapY,  nMapZ,  nTargetID,  nItemKey,  nInt1,  nInt2 )
VALUES          ( @nType, @sMap, @nMapX, @nMapY, @nMapZ, @nTargetID, @nItemKey, @nInt1, @nInt2 )
SET @nLogNo = @@IDENTITY
-- end
GO

-- ========================================================================
-- Procedure: p_GameLog_Add_TypeF
-- ========================================================================
CREATE PROCEDURE [dbo].[p_GameLog_Add_TypeF]
/*
 Game log add
 input:
  All Args = 로그에 기록될 정보
 output: 
  nLogNo = 추가된 로그의 번호(에러시 LogNo = 0 )
*/
@nLogNo bigint OUTPUT,
@nType int,
@nCharNo int,
@sMap varchar(16),
@nMapX int,
@nMapY int,
@nMapZ int,
@nTargetID int,
@nItemKey bigint,
@nInt1 int,
@nInt2 int,
@nInt3 int,
@nBigint1 bigint
AS
SET NOCOUNT ON
SET @nLogNo = 0
INSERT tGameLog (  nType,  nCharNo,  sMap,  nMapX,  nMapY,  nMapZ,  nTargetID,  nItemKey,  nInt1,  nInt2,  nInt3,  nBigint1 )
VALUES          ( @nType, @nCharNo, @sMap, @nMapX, @nMapY, @nMapZ, @nTargetID, @nItemKey, @nInt1, @nInt2, @nInt3, @nBigint1 )
SET @nLogNo = @@IDENTITY
-- end

0騁ａń
GO

-- ========================================================================
-- Procedure: p_GameLog_Add_TypeG
-- ========================================================================
CREATE PROCEDURE [dbo].[p_GameLog_Add_TypeG]
/*
 Game log add
 input:
  All Args = 로그에 기록될 정보
 output: 
  nLogNo = 추가된 로그의 번호(에러시 LogNo = 0 )
*/
@nLogNo bigint OUTPUT,
@nType int,
@nCharNo int,
@sMap varchar(16),
@nMapX int,
@nMapY int,
@nMapZ int,
@nTargetID int,
@nItemKey bigint,
@nInt1 int,
@nInt2 int,
@nBigint1 bigint
AS
SET NOCOUNT ON
SET @nLogNo = 0
INSERT tGameLog (  nType,  nCharNo,  sMap,  nMapX,  nMapY,  nMapZ,  nTargetID,  nItemKey,  nInt1,  nInt2,  nBigint1 )
VALUES          ( @nType, @nCharNo, @sMap, @nMapX, @nMapY, @nMapZ, @nTargetID, @nItemKey, @nInt1, @nInt2, @nBigint1 )
SET @nLogNo = @@IDENTITY
-- end

0팁Ņ
GO

-- ========================================================================
-- Procedure: p_GameLog_Add_TypeH
-- ========================================================================
CREATE PROCEDURE [dbo].[p_GameLog_Add_TypeH]
/*
 Game log add
 input:
  All Args = 로그에 기록될 정보
 output: 
  nLogNo = 추가된 로그의 번호(에러시 LogNo = 0 )
*/
@nLogNo bigint OUTPUT,
@nType int,
@nCharNo int,
@sMap varchar(16),
@nMapX int,
@nMapY int,
@nMapZ int,
@nTargetCharNo int
AS
SET NOCOUNT ON
SET @nLogNo = 0
INSERT tGameLog (  nType,  nCharNo,  sMap,  nMapX,  nMapY,  nMapZ,  nTargetCharNo )
VALUES          ( @nType, @nCharNo, @sMap, @nMapX, @nMapY, @nMapZ, @nTargetCharNo )
SET @nLogNo = @@IDENTITY
-- end

0ఁņ
GO

-- ========================================================================
-- Procedure: p_GameLog_Add_TypeI
-- ========================================================================
CREATE PROCEDURE [dbo].[p_GameLog_Add_TypeI]
/*
 Game log add
 input:
  All Args = 로그에 기록될 정보
 output: 
  nLogNo = 추가된 로그의 번호(에러시 LogNo = 0 )
*/
@nLogNo bigint OUTPUT,
@nType int,
@nCharNo int,
@sMap varchar(16),
@nMapX int,
@nMapY int,
@nMapZ int,
@nTargetCharNo int,
@nTargetID int
AS
SET NOCOUNT ON
SET @nLogNo = 0
INSERT tGameLog (  nType,  nCharNo,  sMap,  nMapX,  nMapY,  nMapZ,  nTargetCharNo,  nTargetID )
VALUES          ( @nType, @nCharNo, @sMap, @nMapX, @nMapY, @nMapZ, @nTargetCharNo, @nTargetID )
SET @nLogNo = @@IDENTITY
-- end

0䔁Ň
GO

-- ========================================================================
-- Procedure: p_GameLog_Add_TypeJ
-- ========================================================================
CREATE PROCEDURE [dbo].[p_GameLog_Add_TypeJ]
/*
 Game log add
 input:
  All Args = 로그에 기록될 정보
 output: 
  nLogNo = 추가된 로그의 번호(에러시 LogNo = 0 )
*/
@nLogNo bigint OUTPUT,
@nType int,
@nCharNo int,
@sMap varchar(16),
@nMapX int,
@nMapY int,
@nMapZ int,
@nTargetID int,
@nItemKey bigint,
@nInt1 int,
@nInt2 int
AS
SET NOCOUNT ON
SET @nLogNo = 0
INSERT tGameLog (  nType,  nCharNo,  sMap,  nMapX,  nMapY,  nMapZ,  nTargetID,  nItemKey,  nInt1,  nInt2 )
VALUES          ( @nType, @nCharNo, @sMap, @nMapX, @nMapY, @nMapZ, @nTargetID, @nItemKey, @nInt1, @nInt2 )
SET @nLogNo = @@IDENTITY
-- end
GO

-- ========================================================================
-- Procedure: p_GameLog_Add_TypeK
-- ========================================================================
CREATE PROCEDURE [dbo].[p_GameLog_Add_TypeK]
/*
 Game log add
 input:
  All Args = 로그에 기록될 정보
 output: 
  nLogNo = 추가된 로그의 번호(에러시 LogNo = 0 )
*/
@nLogNo bigint OUTPUT,
@nType int,
@nCharNo int,
@sMap varchar(16),
@nMapX int,
@nMapY int,
@nMapZ int,
@nInt1 int,
@nInt2 int,
@nBigint1 bigint
AS
SET NOCOUNT ON
SET @nLogNo = 0
INSERT tGameLog (  nType,  nCharNo,  sMap,  nMapX,  nMapY,  nMapZ,  nInt1,  nInt2,  nBigint1 )
VALUES          ( @nType, @nCharNo, @sMap, @nMapX, @nMapY, @nMapZ, @nInt1, @nInt2, @nBigint1 )
SET @nLogNo = @@IDENTITY
-- end

0뜁쏶ŉ
GO

-- ========================================================================
-- Procedure: p_GameLog_Add_TypeL
-- ========================================================================
CREATE PROCEDURE [dbo].[p_GameLog_Add_TypeL]
@nLogNo bigint OUTPUT,
@nType int,
@nCharNo int,
@nTargetCharNo int,
@nTargetID int,
@nItemKey bigint,
@nInt1 int,
@nInt2 int,
@nInt3 int,
@nBigint1 bigint
AS
BEGIN
SET NOCOUNT ON
SET @nLogNo = 0
INSERT tGameLog (  nType,  nCharNo,  nTargetCharNo,  nTargetID,  nItemKey,  nInt1,  nInt2,  nInt3, nBigint1 )
VALUES          ( @nType, @nCharNo, @nTargetCharNo, @nTargetID, @nItemKey, @nInt1, @nInt2, @nInt3, @nBigint1 )
SET @nLogNo = @@IDENTITY

END

0렚Ŋ
GO

-- ========================================================================
-- Procedure: p_GameLog_Get
-- ========================================================================
CREATE PROCEDURE [dbo].[p_GameLog_Get]
/*
  GameLog Get - 주어진 조건에 맞는 로그를 모두검색해서 얻는다.
  input:
   - 주어진 조건절
  output:
   - Select된 데이타
*/
@WhereExpr nvarchar(256)
AS
SET NOCOUNT ON
DECLARE @SQLCMD nvarchar(512)
SET @SQLCMD = N'SELECT nLogNo, dDate, nCharNo, nType, sMap, nMapX, nMapY, nMapZ, nTargetCharNo, nTargetID, nItemKey, nInt1, nInt2, nInt3, nBigint1 FROM tGameLog WHERE ' + @WhereExpr
EXEC sp_executesql @SQLCMD
-- end

0⤁갿ŋ
GO

-- ========================================================================
-- Procedure: p_GameLog_GetCount
-- ========================================================================
CREATE PROCEDURE [dbo].[p_GameLog_GetCount]
/*
  GameLog Get Count - 주어진 조건에 맞는 로그를 모두검색해서 얻는다.
  input:
   - 주어진 조건절
  output:
   - @Count = 해당 조건에 만족하는 레코드수
*/
@WhereExpr nvarchar(256),
@Count int OUTPUT
AS
SET NOCOUNT ON
DECLARE @SQLCMD nvarchar(512)
SET @SQLCMD = N'SELECT @RowCount = COUNT(*) FROM tGameLog WHERE ' + @WhereExpr
EXEC sp_executesql @SQLCMD, N'@RowCount AS int OUTPUT', @RowCount = @Count OUTPUT
-- end

0戁ꁣŌ
GO

-- ========================================================================
-- Procedure: p_GameLog_GetLog
-- ========================================================================
CREATE PROCEDURE [dbo].[p_GameLog_GetLog]
	@QueryExpr nvarchar(256)
AS
BEGIN
	SET NOCOUNT ON
	DECLARE @SQLCMD nvarchar(512)
	SET @SQLCMD = N'SELECT a.nLogNo, a.dDate, a.nCharNo, a.nType, a.sMap, a.nMapX, a.nMapY, a.nTargetCharNo, a.nTargetID, a.nItemKey, a.nInt1, a.nInt2, a.nInt3, a.nBigint1, b.sID
			  FROM tGameLog a INNER JOIN [W00CHAR].World00_Character.dbo.tCharacter b ON a.nCharNo = b.nCharNo ' + @QueryExpr
	EXEC sp_executesql @SQLCMD
	SET NOCOUNT OFF
END
GO

-- ========================================================================
-- Procedure: p_GameLog_GetLogCount
-- ========================================================================
CREATE PROCEDURE [dbo].[p_GameLog_GetLogCount]
	@QueryExpr nvarchar(256),
	@Count int OUTPUT
	
AS
BEGIN
	SET NOCOUNT ON
	DECLARE @SQLCMD nvarchar(512)
	SET @SQLCMD = N'SELECT @RowCount = COUNT(*)
			  FROM tGameLog a INNER JOIN [W00CHAR].World00_Character.dbo.tCharacter b ON a.nCharNo = b.nCharNo ' + @QueryExpr
	EXEC sp_executesql @SQLCMD, N'@RowCount AS int OUTPUT', @RowCount = @Count OUTPUT
	SET NOCOUNT OFF
END

0缁攘Q
GO

