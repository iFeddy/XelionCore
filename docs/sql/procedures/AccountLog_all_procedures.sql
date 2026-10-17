-- ========================================================================
-- Stored Procedures Oficiales de Fiesta Online - Base de Datos: AccountLog
-- Extraídos directamente de AccountLog.bak
-- Total procedimientos recuperados: 27 de 27
-- ========================================================================

-- ========================================================================
-- Procedure: p_AccountLog_Add
-- ========================================================================
CREATE PROCEDURE [dbo].[p_AccountLog_Add] 
/*
 Account log add
 input:
  로그에 기록될 정보들(See SQL)
 output: 
  nLogNo = 추가된 로그의 번호(에러시 LogNo = 0 )
*/
@nType		tinyint,
@nUserNo	int,
@nWorldNo	tinyint,
@nPlayMin	smallint,
@nErrorCode smallint,
@sIP		varchar(16),
@sUserID	nvarchar(16),
@sUserPassword	nvarchar(16),
@nLogNo bigint OUTPUT
AS
SET NOCOUNT ON
SET @nLogNo = 0

/*-------------------------------------------------------------------------------------
Project: FiestaHeroes - NA2016
Author: xkl

Update the IP associated with the account upon login.
------------------------------------------------------------------------------------*/
IF @nType = 1
BEGIN
    UPDATE Account..tUser SET sUserIP = @sIP Where nUserNo = @nUserNo
END
/*-------------------------------------------------------------------------------------
 Continue as normal ~ 
------------------------------------------------------------------------------------*/

INSERT tAccountLog  WITH (ROWLOCK) (  nType,  nUserNo,  nWorldNo,  nPlayMin,  nErrorCode,  sIP,  sUserID,  sUserPassword )
VALUES             ( @nType, @nUserNo, @nWorldNo, @nPlayMin, @nErrorCode, @sIP, @sUserID, @sUserPassword )
SET @nLogNo = @@IDENTITY
-- end
GO

-- ========================================================================
-- Procedure: p_AccountLog_Get
-- ========================================================================
CREATE PROCEDURE [p_AccountLog_Get] 
/*
  AccountLog Get - 주어진 조건에 맞는 로그를 모두검색해서 얻는다.
  input:
   - 주어진 조건절
  output:
   - Select된 데이타
*/
@WhereExpr nvarchar(256)
AS
SET NOCOUNT ON
DECLARE @SQLCMD nchar(512)
SET  @SQLCMD = N'SELECT nLogNo, dDate, nType, nUserNo, nWorldNo, nPlayMin, nErrorCode, sIP, sUserID, sUserPassword FROM tAccountLog WITH(NOLOCK) WHERE ' + @WhereExpr
EXEC sp_executesql @SQLCMD
-- end

0㴁ㅾă
GO

-- ========================================================================
-- Procedure: p_AccountLog_GetCount
-- ========================================================================
CREATE PROCEDURE [p_AccountLog_GetCount]
/*
  AccountLog Get Count - 주어진 조건에 맞는 로그를 모두검색해서 얻는다.
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
SET @SQLCMD = N'SELECT @RowCount = COUNT(*) FROM tAccountLog  WITH(NOLOCK)  WHERE ' + @WhereExpr
EXEC sp_executesql @SQLCMD, N'@RowCount AS int OUTPUT', @RowCount = @Count OUTPUT
-- end

0瘁▢Ą
GO

-- ========================================================================
-- Procedure: p_AccountLog_GetLog
-- ========================================================================
CREATE PROCEDURE [p_AccountLog_GetLog]
	@QueryExpr nvarchar(256)
AS
BEGIN
	SET NOCOUNT ON
	DECLARE @SQLCMD nvarchar(512)
	SET @SQLCMD = N'SELECT a.nLogNo, a.dDate, a.nUserNo, b.sUserID, a.nType, a.sIP, a.nPlayMin, a.sUserPassword, a.nErrorCode, a.nWorldNo
			  FROM tAccountLog  a INNER JOIN AccountDB.Account.dbo.tUser b ON a.nUserNo = b.nUserNo ' + @QueryExpr
	EXEC sp_executesql @SQLCMD
	SET NOCOUNT OFF
END

0꼁ᧆą
GO

-- ========================================================================
-- Procedure: p_AccountLog_GetLogCount
-- ========================================================================
CREATE PROCEDURE [p_AccountLog_GetLogCount]
	@QueryExpr nvarchar(256),
	@Count int OUTPUT
AS
BEGIN
	SET NOCOUNT ON
	DECLARE @SQLCMD nvarchar(512)
	SET @SQLCMD = N'SELECT @RowCount = COUNT(*)
			  FROM tAccountLog  a INNER JOIN AccountDB.Account.dbo.tUser b ON a.nUserNo = b.nUserNo ' + @QueryExpr
	EXEC sp_executesql @SQLCMD, N'@RowCount AS int OUTPUT', @RowCount = @Count OUTPUT
	SET NOCOUNT OFF
END

0෪Ć
GO

-- ========================================================================
-- Procedure: p_Event_Friend_DelRequest
-- ========================================================================
CREATE PROCEDURE [p_Event_Friend_DelRequest]
@nEventUserNo 	int,
@err 		smallint OUTPUT

AS
SET NOCOUNT ON

DELETE FROM tEvent_Friend WHERE nEventUserNo= @nEventUserNo
IF @@ROWCOUNT = 0
	SET @err = 1
ELSE
	SET @err = 0	

RETURN @@ERROR




0錁Ĉ
GO

-- ========================================================================
-- Procedure: p_Event_Friend_EventComplete
-- ========================================================================
CREATE PROCEDURE [p_Event_Friend_EventComplete]
@nEventUserNo 	int,
@nFriendCharNo	int OUTPUT,
@err 		smallint OUTPUT

AS
SET NOCOUNT ON

SET @nFriendCharNo = 0

SELECT @nFriendCharNo = nFriendCharNo FROM tEvent_Friend WITH(READUNCOMMITTED) WHERE nEventUserNo = @nEventUserNo
IF @@ROWCOUNT = 0
BEGIN
	SET @err = 1
	RETURN @@ERROR
END

UPDATE tEvent_Friend SET  bEventComplete = 1 WHERE nEventUserNo = @nEventUserNo
IF @@ERROR <> 0
	SET @err = @@ERROR
ELSE
	SET @err = 0

RETURN @@ERROR




0찁ĉ
GO

-- ========================================================================
-- Procedure: p_Event_Friend_Get
-- ========================================================================
CREATE PROCEDURE [p_Event_Friend_Get]
	@nFriendUserNo	int

AS
BEGIN

	SET NOCOUNT ON
	 
	SELECT 
		t.nEventUserNo,
		t.nEventCharNo,
		t.nFriendCharNo,
		t.nWorldNo,
		t.bEventComplete,
		t.dDate
	FROM tEvent_Friend t
	WHERE t.nFriendUserNo = @nFriendUserNo

	SET NOCOUNT OFF

	RETURN @@ERROR
END
0
GO

-- ========================================================================
-- Procedure: p_Event_Friend_GetCount
-- ========================================================================
CREATE PROCEDURE [p_Event_Friend_GetCount]
@nUserNo int,
@nCount int OUTPUT 

AS
SET NOCOUNT ON


SELECT @nCount = COUNT(*) FROM tEvent_Friend WITH(READUNCOMMITTED)
WHERE nFriendUserNo = @nUserNo

RETURN @@error



0㸁웄ċ
GO

-- ========================================================================
-- Procedure: p_Event_Friend_Giveup
-- ========================================================================
CREATE PROCEDURE [p_Event_Friend_Giveup]
@nEventUserNo 	int,
@nFriendCharNo 	int OUTPUT,
@err 		smallint OUTPUT

AS
SET NOCOUNT ON

SET @nFriendCharNo = 0

SELECT @nFriendCharNo = nFriendCharNo FROM tEvent_Friend WHERE nEventUserNo= @nEventUserNo
IF @@ROWCOUNT = 0
	BEGIN
		SET @err = 1
	END

ELSE

	BEGIN
		DELETE FROM tEvent_Friend WHERE nEventUserNo= @nEventUserNo
		SET @err = 0
	END

RETURN @@ERROR
GO

-- ========================================================================
-- Procedure: p_Event_Friend_SetHideMsg
-- ========================================================================
CREATE PROCEDURE [p_Event_Friend_SetHideMsg]
@nUserNo 	int,
@err 		smallint OUTPUT

AS
SET NOCOUNT ON

IF EXISTS (SELECT nUserNo FROM tEvent_Friend_HideMsg WITH(READUNCOMMITTED)
		WHERE nUserNo= @nUserNo )
	BEGIN
		SET @err = 1
	END
ELSE
	BEGIN
		INSERT INTO tEvent_Friend_HideMsg(nUserNo) VALUES(@nUserNo)
		IF @@ERROR <> 0
			SET @err = @@ERROR
		ELSE
			SET @err = 0
	END


RETURN




0뀁꼌č
GO

-- ========================================================================
-- Procedure: p_Event_WorldCup_Enter
-- ========================================================================
CREATE PROCEDURE [p_Event_WorldCup_Enter]
@nCharNo	int,
@nUserNo	int,
@nLevel		tinyint,
@nWorldNo	tinyint,
@err 		smallint OUTPUT

AS
SET NOCOUNT ON

IF GETDATE() > '2006-06-07'
BEGIN
	SET @err = 3
	RETURN @@ERROR
END

IF EXISTS ( SELECT nEventCharNo FROM tEventWorldCup WHERE nEventCharNo=@nCharNo AND nLevel=@nLevel AND nWorldNo=@nWorldNo )
	BEGIN
		SET @err = 1
	END 
ELSE
	BEGIN
		INSERT INTO tEventWorldCup(nEventCharNo, nEventUserNo, nLevel, nWorldNo)
			VALUES(@nCharNo, @nUserNo, @nLevel, @nWorldNo )
		IF @@ERROR <> 0 --Error
			SET @err = 2
		ELSE
			SET @err = 0
	END

RETURN @@ERROR



0嬁譹Đ
GO

-- ========================================================================
-- Procedure: p_IPBlock_Add
-- ========================================================================
CREATE PROCEDURE p_IPBlock_Add 
/*
 IP Block Add
 input:
  추가할 블럭IP 대역
 output: 
  nRet,  0 = OK, ?=Error code
*/
@sIPStart	nvarchar(20),
@sIPEnd	nvarchar(20),
@sDesc	nvarchar(30),
@nRet		int OUTPUT
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = -1
	INSERT tIPBlock( sIPStart, sIPEnd, sDesc )
	VALUES             ( @sIPStart, @sIPEnd, @sDesc )
	SET @nRet = @@ERROR

	SET NOCOUNT OFF
END
-- end
0䈁ᯜĮ
GO

-- ========================================================================
-- Procedure: p_IPBlock_Delete
-- ========================================================================
CREATE PROCEDURE p_IPBlock_Delete
/*
 IP Block Delete
 input:
  블럭이 되어있는 정보 제거
 output: 
  Selected records
*/
@sIP 		nvarchar(20),
@nRet		int OUTPUT
AS
BEGIN
	SET NOCOUNT ON
	DELETE FROM tIPBlock WHERE @sIP >= sIPStart AND @sIP <= sIPEnd
	SET @nRet = @@ERROR
	SET NOCOUNT OFF
END
-- end
0笁ကį
GO

-- ========================================================================
-- Procedure: p_IPBlock_Get
-- ========================================================================
CREATE PROCEDURE p_IPBlock_Get
/*
 IP Block Get
 input:
  블럭이 되어있는 정보 얻기
 output: 
  Selected records
*/
@sIP 		nvarchar(20)
AS
BEGIN
	SET NOCOUNT ON
	SELECT sIPStart, sIPEnd, sDesc, dDate FROM tIPBlock WHERE @sIP >= sIPStart AND @sIP <= sIPEnd
	SET NOCOUNT OFF
END
-- end
0됁Фİ
GO

-- ========================================================================
-- Procedure: p_IPBlock_GetList
-- ========================================================================
CREATE PROCEDURE p_IPBlock_GetList
/*
 IP Block Get
 input:
  블럭이 되어있는 목록
 output: 
  Selected records
*/
AS
BEGIN
	SET NOCOUNT ON
	SELECT sIPStart, sIPEnd, sDesc, dDate FROM tIPBlock
	SET NOCOUNT OFF
END
-- end
0弁2
GO

-- ========================================================================
-- Procedure: p_IPBlock_IsBlock
-- ========================================================================
CREATE  PROCEDURE p_IPBlock_IsBlock
/*
 Is IP block ?
 input:
  IP
 output: 
  bIsBlock, 0=No block, 1=Block
*/
	@nUserNo    int
,	@sIP 		nvarchar(20)
,	@bIsBlock	tinyint OUTPUT
AS
BEGIN
	SET NOCOUNT ON

	SET @bIsBlock = 0

	IF EXISTS(SELECT sIPStart FROM tIPBlock WHERE @sIP >= sIPStart AND @sIP <= sIPEnd) BEGIN
		IF NOT EXISTS(SELECT TOP 1 nLogNo FROM tAccountLog WHERE nUserNo = @nUserNo AND nType = 1) BEGIN
			SET @bIsBlock = 1
		END
	END

	SET NOCOUNT OFF
END
-- end
GO

-- ========================================================================
-- Procedure: p_Report_Add
-- ========================================================================
CREATE PROCEDURE [p_Report_Add] 
@nReporter int,
@sReporterCharID nvarchar(16),
@sHarmer nvarchar(16),
@nWorldNo tinyint,
@sReportType nvarchar(16),
@sReason nvarchar(512),
@sChatLog text,
@nReportRegNo int OUTPUT 
AS
SET NOCOUNT ON
DECLARE	@ret	int
SET 		@ret = 0
--1시간에 1회만 가능합니다.
IF EXISTS (SELECT nReporter FROM tReport(NOLOCK) 
		WHERE 	nReporter = @nReporter AND 
			nWorldNo=@nWorldNo AND
			dReportDate > dateadd(hour, -1, getdate()) )
BEGIN
	SET @nReportRegNo = -1
	RETURN @ret
END
--동일 캐릭터에 대해서는 1일 1회만 가능합니다.
IF EXISTS (SELECT nReporter FROM tReport(NOLOCK) 
		WHERE 	nReporter = @nReporter AND 
			nWorldNo=@nWorldNo AND
			dReportDate > dateadd(day, -1, getdate()) AND
			sHarmer = @sHarmer)
BEGIN
	SET @nReportRegNo = -2
	RETURN @ret
END
INSERT tReport (sReportType,nReporter,sReporterCharID,sHarmer,nWorldNo,sReason,sChatLog)
VALUES (@sReportType,@nReporter, @sReporterCharID, @sHarmer, @nWorldNo,@sReason,@sChatLog )
SET @ret = @@error
IF @ret <> 0 
	SET @nReportRegNo = 0
ELSE
	SET @nReportRegNo = @@IDENTITY
RETURN @ret

0촁珁Ē
GO

-- ========================================================================
-- Procedure: p_Report_Deal_OPT
-- ========================================================================
CREATE PROCEDURE [p_Report_Deal_OPT]
	@nRegNumber	int,
	@nResult	tinyint,
	@sRemark	nvarchar(128),
	
	@nRet		int OUTPUT
	
AS
SET NOCOUNT ON
UPDATE tReport SET 
	nResult			= @nResult,
	sRemark			= @sRemark
WHERE nRegNumber = @nRegNumber
IF @@ERROR <> 0 OR @@ROWCOUNT = 0
BEGIN
	SET @nRet = 0
	RETURN
END
SET @nRet = 1

0؁柦ē
GO

-- ========================================================================
-- Procedure: p_Report_Del
-- ========================================================================
CREATE PROCEDURE [p_Report_Del]
@nRegNumber int,
@nRet smallint OUTPUT
AS
SET NOCOUNT ON
DECLARE	
@dReportDate	datetime,
@nResult	tinyint
--NoError(ErrorCode.h와 연동됨)
SET @nRet = 0
SELECT @dReportDate=dReportDate, @nResult=nResult FROM tReport WHERE nRegNumber=@nRegNumber 
--잘못된 등록번호
IF @@ROWCOUNT = 0
BEGIN
	SET @nRet = 1
	RETURN @@error
END
--1시간 지난 신고는 취소할수 없음.
IF @dReportDate < DATEADD(hour, -1, GETDATE())
BEGIN
	SET @nRet = 2
	RETURN @@error
END
--처리된 결과는 취소할 수 없음
IF @nResult <> 0
BEGIN
	SET @nRet = 3
	RETURN @@error
END
DELETE FROM tReport WHERE nRegNumber=@nRegNumber
RETURN @@error

0㼁尊Ĕ
GO

-- ========================================================================
-- Procedure: p_Report_Get
-- ========================================================================
CREATE PROCEDURE [p_Report_Get] 
@nReporter int,
@nWorldNo tinyint
AS
SET NOCOUNT ON
SELECT 	nRegNumber,nResult,sRemark  FROM tReport
WHERE 	nReporter = @nReporter AND nWorldNo=@nWorldNo AND
	dReportDate > dateadd(hour, -1, getdate()) 
RETURN @@error

0㡶
GO

-- ========================================================================
-- Procedure: p_Report_Get_Count_OPT
-- ========================================================================
CREATE PROCEDURE [p_Report_Get_Count_OPT]
	@QueryExpr nvarchar(256),
	@Count int OUTPUT
AS
BEGIN
	SET NOCOUNT ON
	DECLARE @SQLCMD nvarchar(512)
	SET @SQLCMD = N'SELECT @RowCount = COUNT(*) 
		FROM tReport ' + @QueryExpr
	EXEC sp_executesql @SQLCMD, N'@RowCount AS int OUTPUT', @RowCount = @Count OUTPUT
	SET NOCOUNT OFF
END
GO

-- ========================================================================
-- Procedure: p_Report_Get_OPT
-- ========================================================================
CREATE PROCEDURE [p_Report_Get_OPT]
	@QueryExpr nvarchar(256)
AS
BEGIN
	SET NOCOUNT ON
	DECLARE @SQLCMD nvarchar(512)
	SET @SQLCMD = N'SELECT nRegNumber, sReportType, nReporter, sReporterCharID, sHarmer, nWorldNo, sReason, sChatLog, dReportDate, nResult, sRemark 
		FROM tReport ' + @QueryExpr
	EXEC sp_executesql @SQLCMD
	SET NOCOUNT OFF
END 

0威ć
GO

-- ========================================================================
-- Procedure: usp_EventiPad_getList
-- ========================================================================
CREATE PROCEDURE dbo.usp_EventiPad_getList
AS
BEGIN
	SET NOCOUNT ON

	SELECT nUserNo, nTicket, sys_create_date AS [dCreateDate], sys_update_date AS [dLastUpdate]
	  FROM tEventiPad
	 ORDER BY nTicket DESC
	
	SET NOCOUNT OFF
END
0騁ａD
GO

-- ========================================================================
-- Procedure: usp_EventiPad_getTicket
-- ========================================================================
CREATE PROCEDURE dbo.usp_EventiPad_getTicket
	@nUserNo int,
	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	IF ISNULL(@nUserNo, 0) <= 0 BEGIN
		SET @nRet = -1
		RETURN
	END
		
	SELECT nUserNo, nTicket, sys_create_date AS [dCreateDate], sys_update_date AS [dLastUpdate]
	  FROM tEventiPad
	 WHERE nUserNo = @nUserNo
	
	SET NOCOUNT OFF
END
GO

-- ========================================================================
-- Procedure: usp_EventiPad_setTicket
-- ========================================================================
CREATE PROCEDURE dbo.usp_EventiPad_setTicket
	@nUserNo   int,
	@nPlayTime int,
	@nRet      int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	DECLARE @nTicket AS int SET @nTicket = 0
	SET @nRet = 0
	
	IF ISNULL(@nPlayTime, 0) < 60 BEGIN
		SET @nRet = -1
		RETURN
	END
	ELSE BEGIN
		SET @nTicket = @nPlayTime / 60
	END
	
	IF EXISTS(SELECT nUserNo FROM tEventiPad WHERE nUserNo = @nUserNo) BEGIN
		UPDATE tEventiPad SET
			nTicket = nTicket + @nTicket,
			sys_update_date = GETDATE()
		 WHERE nUserNo = @nUserNo
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = -8
			RETURN
		END
	END
	ELSE BEGIN
		INSERT INTO tEventiPad (nUserNo, nTicket)
		VALUES (@nUserNo, @nTicket)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -7
			RETURN
		END
	END
	
	SET NOCOUNT OFF
END
0䌁넢Ķ
GO

-- ========================================================================
-- Procedure: usp_Log_setLogout
-- ========================================================================
CREATE PROCEDURE [usp_Log_setLogout]
	@WORLD_NO TINYINT		-- 서버
,	@USER_NO INT			-- 유저 번호
,	@CHAR_NO INT			-- 캐릭터 번호
,	@CLASS TINYINT			-- 캐릭터 클래스
,	@LEVEL SMALLINT			-- 캐릭터 현재레벨
,	@EXPERIENCE BIGINT		-- 획득 경험치
,	@PLAYTIME SMALLINT		-- 플레이시간
,	@RET INT = 0 OUTPUT		-- 리턴값(0:실패/1:성공)
AS

BEGIN

	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	INSERT INTO tLogout (worldNo, userNo, charNo, class, level, experience, playTime)
	VALUES (@WORLD_NO, @USER_NO, @CHAR_NO, @CLASS, @LEVEL, @EXPERIENCE, @PLAYTIME)
	IF @@ERROR <> 0 OR @@IDENTITY = 0
		SET @RET = 0
	ELSE
		SET @RET = 1

	SET NOCOUNT OFF
	SET LOCK_TIMEOUT -1

END

0∁靕ď
GO

