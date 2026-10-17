-- ========================================================================
-- Stored Procedures Oficiales de Fiesta Online - Base de Datos: OperatorTool
-- Extraídos directamente de OperatorTool.bak
-- Total procedimientos recuperados: 17 de 17
-- ========================================================================

-- ========================================================================
-- Procedure: p_OPLog_Add
-- ========================================================================
CREATE PROCEDURE [dbo].[p_OPLog_Add]
/*
 Operator log add
 input:
  All Args = 로그에 기록될 정보
 output: 
  nLogNo = 추가된 로그의 번호(에러시 LogNo = 0 )
*/
@nLogType	tinyint,
@nOperNo	smallint,
@sIP		nvarchar(20),
@nType		smallint,
@nUserNo	int,
@sServer	nvarchar(50),
@nCharNo	int,
@sDesc		nvarchar(200),
@nLogNo		int OUTPUT
AS
BEGIN
SET NOCOUNT ON
SET @nLogNo = 0
DECLARE @STMT	AS NVARCHAR(200)
SET @STMT = 'INSERT tOpLog ' + 
	CASE @nLogType
		WHEN 0 THEN		'( nOperNo,		sIP,	nType,	nUserNo,	sServer,	nCharNo,	sDesc )
				VALUES	 ( @nOperNo,	@sIP,	@nType, @nUserNo,	@sServer,	@nCharNo,	@sDesc )'
		WHEN 1 THEN		'( nOperNo,		sIP,	nType,	sDesc )
				VALUES	 ( @nOperNo,	@sIP,	@nType, @sDesc )'
		WHEN 2 THEN		'( nOperNo,		sIP,	nType,	nUserNo,	sServer,	sDesc )
				VALUES	 ( @nOperNo,	@sIP,	@nType, @nUserNo,	@sServer,	@sDesc )'
		WHEN 3 THEN		'( nOperNo,		sIP,	nType,	sServer,	sDesc )
				VALUES	 ( @nOperNo,	@sIP,	@nType, @sServer,	@sDesc )'
	END
PRINT @STMT
EXEC sp_executesql @STMT
SET @nLogNo = @@IDENTITY
END
-- end

0긁蒀ż
GO

-- ========================================================================
-- Procedure: p_OPLog_Add_All
-- ========================================================================
CREATE PROCEDURE [dbo].[p_OPLog_Add_All]
/*
 Operator log add
 input:
  All Args = 로그에 기록될 정보
 output: 
  nLogNo = 추가된 로그의 번호(에러시 LogNo = 0 )
*/
@nOperNo	smallint,
@sIP		nvarchar(20),
@nType		smallint,
@nUserNo	int,
@sServer	nvarchar(50),
@nCharNo	int,
@sDesc		nvarchar(200),
@nLogNo		int OUTPUT
AS
SET NOCOUNT ON
SET @nLogNo = 0
INSERT tOpLog ( nOperNo,	sIP,	nType,	nUserNo,	sServer,	nCharNo,	sDesc )
VALUES        ( @nOperNo,	@sIP,	@nType, @nUserNo,	@sServer,	@nCharNo,	@sDesc )
SET @nLogNo = @@IDENTITY
-- end

0㸁웄ċ
GO

-- ========================================================================
-- Procedure: p_OPLog_Add_ToIP
-- ========================================================================
CREATE PROCEDURE [dbo].[p_OPLog_Add_ToIP]
/*
 Operator log add
 input:
  All Args = 로그에 기록될 정보
 output: 
  nLogNo = 추가된 로그의 번호(에러시 LogNo = 0 )
*/
@nOperNo	smallint,
@sIP		nvarchar(20),
@nType		smallint,
@sServer	nvarchar(50),
@sDesc		nvarchar(200),
@nLogNo		int OUTPUT
AS
SET NOCOUNT ON
SET @nLogNo = 0
INSERT tOpLog ( nOperNo,	sIP,	nType,	sServer,	sDesc )
VALUES        ( @nOperNo,	@sIP,	@nType, @sServer,	@sDesc )
SET @nLogNo = @@IDENTITY
-- end

0省뫨Č
GO

-- ========================================================================
-- Procedure: p_OPLog_Add_ToUserNo
-- ========================================================================
CREATE PROCEDURE [dbo].[p_OPLog_Add_ToUserNo]
/*
 Operator log add
 input:
  All Args = 로그에 기록될 정보
 output: 
  nLogNo = 추가된 로그의 번호(에러시 LogNo = 0 )
*/
@nOperNo	smallint,
@sIP		nvarchar(20),
@nType		smallint,
@nUserNo	int,
@sServer	nvarchar(50),
@sDesc		nvarchar(200),
@nLogNo		int OUTPUT
AS
SET NOCOUNT ON
SET @nLogNo = 0
INSERT tOpLog ( nOperNo,	sIP,	nType,	nUserNo,	sServer,	sDesc )
VALUES        ( @nOperNo,	@sIP,	@nType, @nUserNo,	@sServer,	@sDesc )
SET @nLogNo = @@IDENTITY
-- end
GO

-- ========================================================================
-- Procedure: p_OPLog_Add_ToWorldNo
-- ========================================================================
CREATE PROCEDURE [dbo].[p_OPLog_Add_ToWorldNo]
/*
 Operator log add
 input:
  All Args = 로그에 기록될 정보
 output: 
  nLogNo = 추가된 로그의 번호(에러시 LogNo = 0 )
*/
@nOperNo	smallint,
@sIP		nvarchar(20),
@nType		smallint,
@sServer	nvarchar(50),
@sDesc		nvarchar(200),
@nLogNo		int OUTPUT
AS
SET NOCOUNT ON
SET @nLogNo = 0
INSERT tOpLog ( nOperNo,	sIP,	nType,	sServer,	sDesc )
VALUES        ( @nOperNo,	@sIP,	@nType, @sServer,	@sDesc )
SET @nLogNo = @@IDENTITY
-- end

0ꌰĎ
GO

-- ========================================================================
-- Procedure: p_OPLog_Get
-- ========================================================================
CREATE PROCEDURE [dbo].[p_OPLog_Get]
/*
  OPLog Get - 주어진 조건에 맞는 로그를 모두검색해서 얻는다.
  input:
   - 주어진 조건절
  output:
   - Select된 데이타
*/
@WhereExpr nvarchar(256)
AS
SET NOCOUNT ON
DECLARE @SQLCMD nvarchar(512)
SET @SQLCMD = N'SELECT nNo, dDate, nOperNo, sIP, nType, nUserNo, sServer, nCharNo, sDesc FROM tOpLog ' + @WhereExpr
EXEC sp_executesql @SQLCMD
-- end

0碤Ž
GO

-- ========================================================================
-- Procedure: p_OPLog_GetCount
-- ========================================================================
CREATE PROCEDURE [dbo].[p_OPLog_GetCount]
/*
  OPLog Get Count - 주어진 조건에 맞는 로그를 모두검색해서 얻는다.
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
SET @SQLCMD = N'SELECT @RowCount = COUNT(*) FROM tOpLog ' + @WhereExpr
EXEC sp_executesql @SQLCMD, N'@RowCount AS int OUTPUT', @RowCount = @Count OUTPUT
-- end

0Ċ
GO

-- ========================================================================
-- Procedure: p_OPTMenu_GetAuth
-- ========================================================================
CREATE PROCEDURE [dbo].[p_OPTMenu_GetAuth]
@sMenuName	NVARCHAR(20)
AS
SET NOCOUNT ON
SELECT auth FROM tMenuAuth WHERE menu = @sMenuName

0威ć
GO

-- ========================================================================
-- Procedure: p_OPTMenu_GetAuthAll
-- ========================================================================
CREATE PROCEDURE [dbo].[p_OPTMenu_GetAuthAll]
AS
SET NOCOUNT ON
SELECT no, menu, auth FROM tMenuAuth

0錁Ĉ
GO

-- ========================================================================
-- Procedure: p_OPTMenu_GetCount
-- ========================================================================
CREATE PROCEDURE [dbo].[p_OPTMenu_GetCount]
@nCount tinyint OUTPUT
AS
SET NOCOUNT ON
SELECT @nCount = COUNT(*) FROM tMenuAuth

0찁ĉ
GO

-- ========================================================================
-- Procedure: p_OPTMenu_SetAuth
-- ========================================================================
CREATE PROCEDURE [dbo].[p_OPTMenu_SetAuth]
@sMenuName	NVARCHAR(20),
@nAuth		tinyint,
@nRet		int OUTPUT
AS
SET NOCOUNT ON
UPDATE tMenuAuth SET auth = @nAuth WHERE menu = @sMenuName
IF (@@ERROR = 0) AND (@@ROWCOUNT = 1)
	SET @nRet = 1
ELSE
	SET @nRet = 0

0ԁ튠Ċ
GO

-- ========================================================================
-- Procedure: p_Operator_Add
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Operator_Add]
/*
Operator add
input:
  See SQL
output: 
  nOperNo = 추가된 Operator의 번호(에러시 = 0 )
*/
@nOperNo	int OUTPUT,
@sOperID	nvarchar(20),
@sOperPW	nvarchar(20),
@nOperLevel	tinyint
AS
SET NOCOUNT ON
SET @nOperNo = 0
INSERT tOperator (  sOperID,  sOperPW,  nOperLevel )
VALUES           ( @sOperID, @sOperPW, @nOperLevel )
IF @@ERROR = 0
	SET @nOperNo = @@IDENTITY
	
-- end

0∁靕ď
GO

-- ========================================================================
-- Procedure: p_Operator_Del
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Operator_Del]
/*
Operator delete
input:
  See SQL
output: 
  @nOperLevel = 0, Delete failed
                1, Delete successed
*/
@nOperNo	int,
@nOperLevel tinyint OUTPUT
AS
SET NOCOUNT ON
SET @nOperLevel = 0
DELETE FROM tOperator WHERE nOperNo = @nOperNo
IF @@ROWCOUNT > 0
	SET @nOperLevel = 1
	
-- end

0
GO

-- ========================================================================
-- Procedure: p_Operator_GetAll
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Operator_GetAll]
/*
  Operator Get All
  output:
   - Select된 데이타
*/
AS
SET NOCOUNT ON
SELECT nOperNo, sOperID, nOperLevel, dRegDate FROM tOperator
-- end

0鐁羝đ
GO

-- ========================================================================
-- Procedure: p_Operator_Logon
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Operator_Logon]
/*
Operator logon
2004.12 By CJC
- Check ID.Password
Output:
	@nOperNo = 0 Failed logon (ID,Password missmatch) 
			 = ? Successed logon ( @nOperLevel = Operator level )
*/
@sOperID	nvarchar(20),
@sOperPW	nvarchar(20),
@nOperNo	int OUTPUT,
@nOperLevel	int OUTPUT
AS
SET NOCOUNT ON
SET @nOperNo = 0
IF EXISTS(SELECT * FROM tOperator WHERE (sOperID = @sOperID) AND (sOperPW = @sOperPW) )
BEGIN
	SELECT @nOperNo = nOperNo, @nOperLevel = nOperLevel 
	FROM tOperator WHERE (sOperID = @sOperID) AND (sOperPW = @sOperPW)
END
-- end

0
GO

-- ========================================================================
-- Procedure: p_Operator_SetLevel
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Operator_SetLevel]
/*
Operator set level
input:
  See SQL
output: 
  nRet = 1 TRUE
	   = 0 FALSE
*/
@nOperNo	int,
@nOperLevel	tinyint,
@nRet		int OUTPUT
AS
SET NOCOUNT ON
UPDATE tOperator SET nOperLevel = @nOperLevel WHERE nOperNo = @nOperNo
IF (@@ERROR = 0) AND (@@ROWCOUNT = 1)
	SET @nRet = 1
ELSE
	SET @nRet = 0
-- end

0
GO

-- ========================================================================
-- Procedure: p_Operator_SetPW
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Operator_SetPW]
/*
Operator set password
input:
  See SQL
output: 
  nRet = 1 TRUE
	   = 0 FALSE
*/
@nOperNo	int,
@sOperPW	nvarchar(20),
@nRet		int OUTPUT
AS
SET NOCOUNT ON
UPDATE tOperator SET sOperPW = @sOperPW WHERE nOperNo = @nOperNo
IF (@@ERROR = 0) AND (@@ROWCOUNT = 1)
	SET @nRet = 1
ELSE
	SET @nRet = 0
-- end

<
GO

