-- ========================================================================
-- Stored Procedures Oficiales de Fiesta Online - Base de Datos: World00_Character
-- Extraídos directamente de World00_Character.bak
-- Total procedimientos recuperados: 408 de 408
-- ========================================================================

-- ========================================================================
-- Procedure: Adam_ListProc
-- ========================================================================
create procedure [dbo].[Adam_ListProc] as  select o.name, m.definition   from sys.sql_modules m join sys.objects o  On m.object_id=o.object_id and o.type='P'

0
GO

-- ========================================================================
-- Procedure: CommandExecute
-- ========================================================================
CREATE PROCEDURE [dbo].[CommandExecute]

@Command nvarchar(max),
@CommandType nvarchar(max),
@Mode int,
@Comment nvarchar(max) = NULL,
@DatabaseName nvarchar(max) = NULL,
@SchemaName nvarchar(max) = NULL,
@ObjectName nvarchar(max) = NULL,
@ObjectType nvarchar(max) = NULL,
@IndexName nvarchar(max) = NULL,
@IndexType int = NULL,
@StatisticsName nvarchar(max) = NULL,
@PartitionNumber int = NULL,
@ExtendedInfo xml = NULL,
@LogToTable nvarchar(max),
@Execute nvarchar(max)

AS

BEGIN

  ----------------------------------------------------------------------------------------------------
  --// Source: https://ola.hallengren.com                                                          //--
  ----------------------------------------------------------------------------------------------------

  SET NOCOUNT ON

  DECLARE @StartMessage nvarchar(max)
  DECLARE @EndMessage nvarchar(max)
  DECLARE @ErrorMessage nvarchar(max)
  DECLARE @ErrorMessageOriginal nvarchar(max)

  DECLARE @StartTime datetime
  DECLARE @EndTime datetime

  DECLARE @StartTimeSec datetime
  DECLARE @EndTimeSec datetime

  DECLARE @ID int

  DECLARE @Error int
  DECLARE @ReturnCode int

  SET @Error = 0
  SET @ReturnCode = 0

  ----------------------------------------------------------------------------------------------------
  --// Check core requirements                                                                    //--
  ----------------------------------------------------------------------------------------------------

  IF @LogToTable = 'Y' AND NOT EXISTS (SELECT * FROM sys.objects objects INNER JOIN sys.schemas schemas ON objects.[schema_id] = schemas.[schema_id] WHERE objects.[type] = 'U' AND schemas.[name] = 'dbo' AND objects.[name] = 'CommandLog')
  BEGIN
    SET @ErrorMessage = 'The table CommandLog is missing. Download https://ola.hallengren.com/scripts/CommandLog.sql.' + CHAR(13) + CHAR(10) + ' '
    RAISERROR(@ErrorMessage,16,1) WITH NOWAIT
    SET @Error = @@ERROR
  END

  IF @Error <> 0
  BEGIN
    SET @ReturnCode = @Error
    GOTO ReturnCode
  END

  ----------------------------------------------------------------------------------------------------
  --// Check input parameters                                                                     //--
  ----------------------------------------------------------------------------------------------------

  IF @Command IS NULL OR @Command = ''
  BEGIN
    SET @ErrorMessage = 'The value for the parameter @Command is not supported.' + CHAR(13) + CHAR(10) + ' '
    RAISERROR(@ErrorMessage,16,1) WITH NOWAIT
    SET @Error = @@ERROR
  END

  IF @CommandType IS NULL OR @CommandType = '' OR LEN(@CommandType) > 60
  BEGIN
    SET @ErrorMessage = 'The value for the parameter @CommandType is not supported.' + CHAR(13) + CHAR(10) + ' '
    RAISERROR(@ErrorMessage,16,1) WITH NOWAIT
    SET @Error = @@ERROR
  END

  IF @Mode NOT IN(1,2) OR @Mode IS NULL
  BEGIN
    SET @ErrorMessage = 'The value for the parameter @Mode is not supported.' + CHAR(13) + CHAR(10) + ' '
    RAISERROR(@ErrorMessage,16,1) WITH NOWAIT
    SET @Error = @@ERROR
  END

  IF @LogToTable NOT IN('Y','N') OR @LogToTable IS NULL
  BEGIN
    SET @ErrorMessage = 'The value for the parameter @LogToTable is not supported.' + CHAR(13) + CHAR(10) + ' '
    RAISERROR(@ErrorMessage,16,1) WITH NOWAIT
    SET @Error = @@ERROR
  END

  IF @Execute NOT IN('Y','N') OR @Execute IS NULL
  BEGIN
    SET @ErrorMessage = 'The value for the parameter @Execute is not supported.' + CHAR(13) + CHAR(10) + ' '
    RAISERROR(@ErrorMessage,16,1) WITH NOWAIT
    SET @Error = @@ERROR
  END

  IF @Error <> 0
  BEGIN
    SET @ReturnCode = @Error
    GOTO ReturnCode
  END

  ----------------------------------------------------------------------------------------------------
  --// Log initial information                                                                    //--
  ----------------------------------------------------------------------------------------------------

  SET @StartTime = GETDATE()
  SET @StartTimeSec = CONVERT(datetime,CONVERT(nvarchar,@StartTime,120),120)

  SET @StartMessage = 'Date and time: ' + CONVERT(nvarchar,@StartTimeSec,120) + CHAR(13) + CHAR(10)
  SET @StartMessage = @StartMessage + 'Command: ' + @Command
  IF @Comment IS NOT NULL SET @StartMessage = @StartMessage + CHAR(13) + CHAR(10) + 'Comment: ' + @Comment
  SET @StartMessage = REPLACE(@StartMessage,'%','%%')
  RAISERROR(@StartMessage,10,1) WITH NOWAIT

  IF @LogToTable = 'Y'
  BEGIN
    INSERT INTO dbo.CommandLog (DatabaseName, SchemaName, ObjectName, ObjectType, IndexName, IndexType, StatisticsName, PartitionNumber, ExtendedInfo, CommandType, Command, StartTime)
    VALUES (@DatabaseName, @SchemaName, @ObjectName, @ObjectType, @IndexName, @IndexType, @StatisticsName, @PartitionNumber, @ExtendedInfo, @CommandType, @Command, @StartTime)
  END

  SET @ID = SCOPE_IDENTITY()

  ----------------------------------------------------------------------------------------------------
  --// Execute command                                                                            //--
  ----------------------------------------------------------------------------------------------------

  IF @Mode = 1 AND @Execute = 'Y'
  BEGIN
    EXECUTE(@Command)
    SET @Error = @@ERROR
    SET @ReturnCode = @Error
  END

  IF @Mode = 2 AND @Execute = 'Y'
  BEGIN
    BEGIN TRY
      EXECUTE(@Command)
    END TRY
    BEGIN CATCH
      SET @Error = ERROR_NUMBER()
      SET @ReturnCode = @Error
      SET @ErrorMessageOriginal = ERROR_MESSAGE()
      SET @ErrorMessage = 'Msg ' + CAST(@Error AS nvarchar) + ', ' + ISNULL(@ErrorMessageOriginal,'')
      RAISERROR(@ErrorMessage,16,1) WITH NOWAIT
    END CATCH
  END

  ----------------------------------------------------------------------------------------------------
  --// Log completing information                                                                 //--
  ----------------------------------------------------------------------------------------------------

  SET @EndTime = GETDATE()
  SET @EndTimeSec = CONVERT(datetime,CONVERT(varchar,@EndTime,120),120)

  SET @EndMessage = 'Outcome: ' + CASE WHEN @Execute = 'N' THEN 'Not Executed' WHEN @Error = 0 THEN 'Succeeded' ELSE 'Failed' END + CHAR(13) + CHAR(10)
  SET @EndMessage = @EndMessage + 'Duration: ' + CASE WHEN DATEDIFF(ss,@StartTimeSec, @EndTimeSec)/(24*3600) > 0 THEN CAST(DATEDIFF(ss,@StartTimeSec, @EndTimeSec)/(24*3600) AS nvarchar) + '.' ELSE '' END + CONVERT(nvarchar,@EndTimeSec - @StartTimeSec,108) + CHAR(13) + CHAR(10)
  SET @EndMessage = @EndMessage + 'Date and time: ' + CONVERT(nvarchar,@EndTimeSec,120) + CHAR(13) + CHAR(10) + ' '
  SET @EndMessage = REPLACE(@EndMessage,'%','%%')
  RAISERROR(@EndMessage,10,1) WITH NOWAIT

  IF @LogToTable = 'Y'
  BEGIN
    UPDATE dbo.CommandLog
    SET EndTime = @EndTime,
        ErrorNumber = CASE WHEN @Execute = 'N' THEN NULL ELSE @Error END,
        ErrorMessage = @ErrorMessageOriginal
    WHERE ID = @ID
  END

  ReturnCode:
  IF @ReturnCode <> 0
  BEGIN
    RETURN @ReturnCode
  END

  ----------------------------------------------------------------------------------------------------

END

0
GO

-- ========================================================================
-- Procedure: IndexOptimize
-- ========================================================================
CREATE PROCEDURE [dbo].[IndexOptimize]

@Databases nvarchar(max),
@FragmentationLow nvarchar(max) = NULL,
@FragmentationMedium nvarchar(max) = 'INDEX_REORGANIZE,INDEX_REBUILD_ONLINE,INDEX_REBUILD_OFFLINE',
@FragmentationHigh nvarchar(max) = 'INDEX_REBUILD_ONLINE,INDEX_REBUILD_OFFLINE',
@FragmentationLevel1 int = 5,
@FragmentationLevel2 int = 30,
@PageCountLevel int = 1000,
@SortInTempdb nvarchar(max) = 'N',
@MaxDOP int = NULL,
@FillFactor int = NULL,
@PadIndex nvarchar(max) = NULL,
@LOBCompaction nvarchar(max) = 'Y',
@UpdateStatistics nvarchar(max) = NULL,
@OnlyModifiedStatistics nvarchar(max) = 'N',
@StatisticsSample int = NULL,
@StatisticsResample nvarchar(max) = 'N',
@PartitionLevel nvarchar(max) = 'Y',
@MSShippedObjects nvarchar(max) = 'N',
@Indexes nvarchar(max) = NULL,
@TimeLimit int = NULL,
@Delay int = NULL,
@WaitAtLowPriorityMaxDuration int = NULL,
@WaitAtLowPriorityAbortAfterWait nvarchar(max) = NULL,
@LockTimeout int = NULL,
@LogToTable nvarchar(max) = 'N',
@Execute nvarchar(max) = 'Y'

AS

BEGIN

  ----------------------------------------------------------------------------------------------------
  --// Source: https://ola.hallengren.com                                                          //--
  ----------------------------------------------------------------------------------------------------

  SET NOCOUNT ON

  SET ARITHABORT ON

  SET NUMERIC_ROUNDABORT OFF

  DECLARE @StartMessage nvarchar(max)
  DECLARE @EndMessage nvarchar(max)
  DECLARE @DatabaseMessage nvarchar(max)
  DECLARE @ErrorMessage nvarchar(max)

  DECLARE @Version numeric(18,10)
  DECLARE @AmazonRDS bit

  DECLARE @Cluster nvarchar(max)

  DECLARE @StartTime datetime

  DECLARE @CurrentDBID int
  DECLARE @CurrentDatabaseID int
  DECLARE @CurrentDatabaseName nvarchar(max)
  DECLARE @CurrentIsDatabaseAccessible bit
  DECLARE @CurrentAvailabilityGroup nvarchar(max)
  DECLARE @CurrentAvailabilityGroupRole nvarchar(max)
  DECLARE @CurrentDatabaseMirroringRole nvarchar(max)
  DECLARE @CurrentLogShippingRole nvarchar(max)

  DECLARE @CurrentCommand01 nvarchar(max)
  DECLARE @CurrentCommand02 nvarchar(max)
  DECLARE @CurrentCommand03 nvarchar(max)
  DECLARE @CurrentCommand04 nvarchar(max)
  DECLARE @CurrentCommand05 nvarchar(max)
  DECLARE @CurrentCommand06 nvarchar(max)
  DECLARE @CurrentCommand07 nvarchar(max)
  DECLARE @CurrentCommand08 nvarchar(max)
  DECLARE @CurrentCommand09 nvarchar(max)
  DECLARE @CurrentCommand10 nvarchar(max)
  DECLARE @CurrentCommand11 nvarchar(max)
  DECLARE @CurrentCommand12 nvarchar(max)
  DECLARE @CurrentCommand13 nvarchar(max)
  DECLARE @CurrentCommand14 nvarchar(max)

  DECLARE @CurrentCommandOutput13 int
  DECLARE @CurrentCommandOutput14 int

  DECLARE @CurrentCommandType13 nvarchar(max)
  DECLARE @CurrentCommandType14 nvarchar(max)

  DECLARE @CurrentIxID int
  DECLARE @CurrentSchemaID int
  DECLARE @CurrentSchemaName nvarchar(max)
  DECLARE @CurrentObjectID int
  DECLARE @CurrentObjectName nvarchar(max)
  DECLARE @CurrentObjectType nvarchar(max)
  DECLARE @CurrentIsMemoryOptimized bit
  DECLARE @CurrentIndexID int
  DECLARE @CurrentIndexName nvarchar(max)
  DECLARE @CurrentIndexType int
  DECLARE @CurrentStatisticsID int
  DECLARE @CurrentStatisticsName nvarchar(max)
  DECLARE @CurrentPartitionID bigint
  DECLARE @CurrentPartitionNumber int
  DECLARE @CurrentPartitionCount int
  DECLARE @CurrentIsPartition bit
  DECLARE @CurrentIndexExists bit
  DECLARE @CurrentStatisticsExists bit
  DECLARE @CurrentIsImageText bit
  DECLARE @CurrentIsNewLOB bit
  DECLARE @CurrentIsFileStream bit
  DECLARE @CurrentIsColumnStore bit
  DECLARE @CurrentAllowPageLocks bit
  DECLARE @CurrentNoRecompute bit
  DECLARE @CurrentStatisticsModified bit
  DECLARE @CurrentOnReadOnlyFileGroup bit
  DECLARE @CurrentFragmentationLevel float
  DECLARE @CurrentPageCount bigint
  DECLARE @CurrentFragmentationGroup nvarchar(max)
  DECLARE @CurrentAction nvarchar(max)
  DECLARE @CurrentMaxDOP int
  DECLARE @CurrentUpdateStatistics nvarchar(max)
  DECLARE @CurrentComment nvarchar(max)
  DECLARE @CurrentExtendedInfo xml
  DECLARE @CurrentDelay datetime

  DECLARE @tmpDatabases TABLE (ID int IDENTITY,
                               DatabaseName nvarchar(max),
                               DatabaseType nvarchar(max),
                               Selected bit,
                               Completed bit,
                               PRIMARY KEY(Selected, Completed, ID))

  DECLARE @tmpIndexesStatistics TABLE (ID int IDENTITY,
                                       SchemaID int,
                                       SchemaName nvarchar(max),
                                       ObjectID int,
                                       ObjectName nvarchar(max),
                                       ObjectType nvarchar(max),
                                       IsMemoryOptimized bit,
                                       IndexID int,
                                       IndexName nvarchar(max),
                                       IndexType int,
                                       StatisticsID int,
                                       StatisticsName nvarchar(max),
                                       PartitionID bigint,
                                       PartitionNumber int,
                                       PartitionCount int,
                                       Selected bit,
                                       Completed bit,
                                       PRIMARY KEY(Selected, Completed, ID))

  DECLARE @SelectedDatabases TABLE (DatabaseName nvarchar(max),
                                    DatabaseType nvarchar(max),
                                    Selected bit)

  DECLARE @SelectedIndexes TABLE (DatabaseName nvarchar(max),
                                  SchemaName nvarchar(max),
                                  ObjectName nvarchar(max),
                                  IndexName nvarchar(max),
                                  Selected bit)

  DECLARE @Actions TABLE ([Action] nvarchar(max))

  INSERT INTO @Actions([Action]) VALUES('INDEX_REBUILD_ONLINE')
  INSERT INTO @Actions([Action]) VALUES('INDEX_REBUILD_OFFLINE')
  INSERT INTO @Actions([Action]) VALUES('INDEX_REORGANIZE')

  DECLARE @ActionsPreferred TABLE (FragmentationGroup nvarchar(max),
                                   [Priority] int,
                                   [Action] nvarchar(max))

  DECLARE @CurrentActionsAllowed TABLE ([Action] nvarchar(max))

  DECLARE @Error int
  DECLARE @ReturnCode int

  SET @Error = 0
  SET @ReturnCode = 0

  SET @Version = CAST(LEFT(CAST(SERVERPROPERTY('ProductVersion') AS nvarchar(max)),CHARINDEX('.',CAST(SERVERPROPERTY('ProductVersion') AS nvarchar(max))) - 1) + '.' + REPLACE(RIGHT(CAST(SERVERPROPERTY('ProductVersion') AS nvarchar(max)), LEN(CAST(SERVERPROPERTY('ProductVersion') AS nvarchar(max))) - CHARINDEX('.',CAST(SERVERPROPERTY('ProductVersion') AS nvarchar(max)))),'.','') AS numeric(18,10))

  SET @AmazonRDS = CASE WHEN DB_ID('rdsadmin') IS NOT NULL AND SUSER_SNAME(0x01) = 'rdsa' THEN 1 ELSE 0 END

  ----------------------------------------------------------------------------------------------------
  --// Log initial information                                                                    //--
  ----------------------------------------------------------------------------------------------------

  SET @StartTime = CONVERT(datetime,CONVERT(nvarchar,GETDATE(),120),120)

  SET @StartMessage = 'Date and time: ' + CONVERT(nvarchar,@StartTime,120) + CHAR(13) + CHAR(10)
  SET @StartMessage = @StartMessage + 'Server: ' + CAST(SERVERPROPERTY('ServerName') AS nvarchar) + CHAR(13) + CHAR(10)
GO

-- ========================================================================
-- Procedure: p_Char_AddExp
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_AddExp]
@nCharNo int,
@nExp bigint

AS
SET NOCOUNT ON

UPDATE tCharacter 
	SET 
		nExp=nExp+@nExp
	WHERE nCharNo=@nCharNo

return @@error

0趎Ĺ
GO

-- ========================================================================
-- Procedure: p_Char_ChangeClass
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_ChangeClass]
@nCharNo int,
@nClass tinyint

AS
SET NOCOUNT ON

UPDATE tCharacterShape 
	SET 
		nClass=@nClass
	WHERE nCharNo=@nCharNo

return @@error

0쐁墄Ĺ
GO

-- ========================================================================
-- Procedure: p_Char_ChangeID
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_ChangeID]
/*
Character change ID 
2004.9 By CJC
 input:
   See code
 output: 
   1 = ok
   0 = error
*/
@nCharNo int,
@sID nvarchar(40), 
@nRet int OUTPUT
AS
SET NOCOUNT ON
-- Change Char ID
UPDATE tCharacter SET sID = @sID WHERE nCharNo = @nCharNo
IF @@ERROR <> 0 OR @@ROWCOUNT = 0
BEGIN
	SET @nRet = 0
	RETURN
END
-- success
SET @nRet = 1
-- end

0㘁䃍Ļ
GO

-- ========================================================================
-- Procedure: p_Char_CheckOptionsExistance
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_CheckOptionsExistance]
	@nCharNo		int,
	@nRet			int OUTPUT
AS

BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000
	
	IF EXISTS(SELECT * FROM  tCharacterOptions WHERE nCharNo = @nCharNo )
		SET @nRet = 1
	ELSE
		SET @nRet = 0
	
	SET LOCK_TIMEOUT -1
	SET NOCOUNT OFF
END

0漁㓱ļ
GO

-- ========================================================================
-- Procedure: p_Char_Create
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_Create] 
	@nUserNo int,			-- 1	1 ~ ?
	@nCreateWorld tinyint,	-- 2	0 ~ ?
	@nAdminLevel tinyint,	-- 3	0=USER, 1~?=GM Level 1~?
	@nSlotNo tinyint,		-- 4	0 ~ 5
	@sID nvarchar(40),		-- 5  
	@nRace tinyint,			-- 6	0 = Human, 1 = Elf, 2 = Dark elf
	@nClass tinyint,		-- 7	1 = Fighter,4 = Cleric, 7 = Archer, 10 = Mage
	@nGender tinyint,		-- 8	0 = Male, 1 = Female
	@nHairType tinyint,		-- 9	0 ~ ?
	@nHairColor tinyint,	-- 10	0 ~ ?
	@nFaceShape tinyint,	-- 11	0 ~ ?
	@nCharNo int OUTPUT     -- 12
AS
BEGIN
	SET NOCOUNT ON
	-- 에러코드 일단 세팅
	SET @nCharNo = 0

	/*-------------------------------------------------------------------------------------
	Project: FiestaHeroes - NA2016
	Author: xkl

	If the account is admin (nAuthID = 9), then make the character have 100 admin power.
	------------------------------------------------------------------------------------*/
	IF EXISTS (SELECT Enabled FROM Options..Auto_GM WHERE Enabled = 1)
	BEGIN
		DECLARE @nAuthID INT = (SELECT nAuthID FROM Account..tUser WHERE nUserNo = @nUserNo)
	
		IF @nAuthID = 9
		BEGIN
			SET @nAdminLevel = 100
		END
	END
	/*-------------------------------------------------------------------------------------
	 Continue as normal ~ 
	------------------------------------------------------------------------------------*/

	-- 해당유저의 해당슬롯에 캐릭이 있는지 확인
	IF EXISTS (SELECT nUserNo FROM tCharacter(NOLOCK) WHERE nUserNo = @nUserNo AND nSlotNo = @nSlotNo AND bDeleted = 0)
	BEGIN
		SET @nCharNo = -1
		RETURN
	END
	-- 캐릭 생성
	
	IF EXISTS (SELECT sID FROM tCharacter(NOLOCK) WHERE sID = @sID)
		BEGIN
		RETURN
		END
		ELSE
	-- 캐릭 생성
	
	BEGIN TRAN
	IF EXISTS (SELECT sID FROM tCharacter(NOLOCK) WHERE sID = @sID)
		BEGIN
		RETURN
		END
		ELSE
	INSERT tCharacter (  nUserNo,  nSlotNo,  sID,  nCreateWorld,  nAdminLevel )
	VALUES            ( @nUserNo, @nSlotNo, @sID, @nCreateWorld, @nAdminLevel )
	IF @@ERROR <> 0
	BEGIN
		ROLLBACK TRAN
		RETURN
	END
	SET @nCharNo = @@IDENTITY
	-- 캐릭 외모 세팅
	INSERT tCharacterShape (  nCharNo,  nRace,  nClass,  nGender,  nHairType,  nHairColor,  nFaceShape )
	VALUES                 ( @nCharNo, @nRace, @nClass, @nGender, @nHairType, @nHairColor, @nFaceShape ) 
	IF @@ERROR <> 0
	BEGIN
		ROLLBACK TRAN
		SET @nCharNo = -2
		RETURN
	END
	-- 캐릭 옵션 세팅
	INSERT tCharacterOptions (  nCharNo )
	VALUES                   ( @nCharNo )
	IF @@ERROR <> 0
	BEGIN
		ROLLBACK TRAN
		SET @nCharNo = -3
		RETURN
	END
	-- 캐릭 생성 성공

	/*-------------------------------------------------------------------------------------
	Project: FiestaHeroes - NA2016
	Author: xkl

	Auto master for your apprentice character. Holy Promise system.
	------------------------------------------------------------------------------------*/
	IF EXISTS (SELECT Enabled FROM Options..Auto_Master WHERE Enabled = 1)
	BEGIN
		DECLARE @Master_nCharNo INT = (SELECT nCharNo FROM Options..Auto_Master WHERE Enabled = 1)
		INSERT INTO tHolyPromise (nCharNo, nUpCharNo, nFlags) VALUES (@nCharNo, @Master_nCharNo, 0)
	END
	/*-------------------------------------------------------------------------------------
	 Continue as normal ~ 
	------------------------------------------------------------------------------------*/

	-- 2010.10.7
	-- kodong
	-- 새로 추가된 단축키 시스템에 의해 초기값 지정
	DECLARE @nRet int
	SET @nRet = 0
	EXEC usp_Character_initKeyMap @nCharNo, @nRet output
	IF @@ERROR <> 0 OR @nRet <> 0 BEGIN
		ROLLBACK TRAN
		SET @nCharNo = -4
		RETURN
	END
	ELSE BEGIN
		COMMIT TRAN
	END
END
GO

-- ========================================================================
-- Procedure: p_Char_CreateSetDefaultData
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_CreateSetDefaultData]
/*
Character Create set default data
2005.6 By CJC

 input:
   All args
   
 output: 
   0 = error
   1 = ok
*/
-- Character columns
@nCharNo int,
@sLoginZone varchar(16),
@nLoginZoneX int,
@nLoginZoneY int,
@nHP smallint, 
@nSP smallint, 
@nHPS smallint,
@nSPS smallint,
@nMoney bigint,
-- Output var
@nRet int OUTPUT

AS

SET NOCOUNT ON

UPDATE tCharacter SET
	sLoginZone	= @sLoginZone,
	nLoginZoneX	= @nLoginZoneX,
	nLoginZoneY	= @nLoginZoneY,
	nHP			= @nHP,
	nSP			= @nSP,
	nHPS		= @nHPS,
	nSPS		= @nSPS,
	nMoney		= @nMoney
WHERE nCharNo = @nCharNo

IF @@ERROR <> 0 OR @@ROWCOUNT = 0
BEGIN
	SET @nRet = 0
	RETURN
END

SET @nRet = 1
-- end

0ᴹľ
GO

-- ========================================================================
-- Procedure: p_Char_DelCharTitle
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_DelCharTitle]

@nCharNo		int,
@nType			tinyint, 
@nElementNo		tinyint, 
@nRet			tinyint OUTPUT
AS
SET NOCOUNT ON
-- 해당 상태 삭제
DELETE FROM tCharacterTitle WHERE nCharNo = @nCharNo AND nType = @nType AND nElementNo = @nElementNo
SET @nRet = @@ERROR
-- end

0匁ւŀ
GO

-- ========================================================================
-- Procedure: p_Char_DelChargedBuff
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_DelChargedBuff]
/*
Char Delete charged buff
2006.6 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK, 
	        ? = Error code
*/

@nCharNo int,
@nKey int,
@nRet int OUTPUT

AS

SET NOCOUNT ON

DELETE FROM tCharacterChargedBuffer WHERE nCharNo = @nCharNo AND nKey = @nKey
SET @nRet = @@ERROR

RETURN

-- end

0ᨁᅞĿ
GO

-- ========================================================================
-- Procedure: p_Char_DelOptions
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_DelOptions]
	@nCharNo		int,
	@nRet			int OUTPUT
AS

BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000
	
	DELETE FROM tCharacterOptions WHERE nCharNo = @nCharNo 
	SET @nRet = @@ERROR
	
	SET LOCK_TIMEOUT -1
	SET NOCOUNT OFF	
END

0⋔ł
GO

-- ========================================================================
-- Procedure: p_Char_DelState
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_DelState]
/*
Char Delete state
2004.10 By CJC
 input:
	nCharNo, nStateNo = To be delete info.
 output:
	@nRet	0 = OK, 
	        ? = Error code
*/
@nCharNo int,
@nStateNo smallint,
@nRet int OUTPUT
AS
SET NOCOUNT ON
-- 해당 상태 삭제
DELETE FROM tCharacterState WHERE nCharNo = @nCharNo AND nStateNo = @nStateNo
SET @nRet = @@ERROR
-- end

0씁Ł
GO

-- ========================================================================
-- Procedure: p_Char_DelStateAll
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_DelStateAll]
/*
Char Delete state all
2004.10 By CJC
 input:
	nCharNo, nStateNo = To be delete info.
 output:
	@nRet	0 = OK, 
	        ? = Error code
*/
@nCharNo int,
@nRet int OUTPUT
AS
SET NOCOUNT ON
-- 해당 상태 삭제
DELETE FROM tCharacterState WHERE nCharNo = @nCharNo
SET @nRet = @@ERROR
-- end

0㜁혓Ń
GO

-- ========================================================================
-- Procedure: p_Char_Delete
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_Delete]
/*
Character Delete
2004.6 By CJC
  Input:
    nCharNo = Character number to be delete.
  Output:
    0 = Delete OK
    1 = 해당 캐릭이 이미 지워진 캐릭임
    2 = 해당 캐릭을 찾을수 없음
    3 = 해당 캐릭을 지울수 없음
  Desc:
    - 지워진 캐릭의 ID는 '아뒤~날짜#시도횟수' 형태로 바뀐다.
    - 버그 방지로 아뒤에는 '~'는 사용되면 안된다.
*/
@nCharNo int,
@nRet int OUTPUT
AS
SET NOCOUNT ON
-- 이미 지워진 캐릭인지 조사
IF EXISTS(SELECT nCharNo FROM tCharacter(NOLOCK) WHERE nCharNo = @nCharNo AND bDeleted > 0)
BEGIN
	SET @nRet = 1
	RETURN
END
-- 지워질 캐릭의 이름 얻기
DECLARE @sID nvarchar(40)
SET @sID = ''
SELECT @sID = sID FROM tCharacter(NOLOCK) WHERE nCharNo = @nCharNo
IF LEN(@sID) = 0
BEGIN
	SET @nRet = 2 -- 해당 캐릭이 없음
	RETURN
END
-- 아뒤내에 최초 삭제마크 이후의 문자는 모두 지운다.(지운거 또 지울수 있음으로 오류방지용)
DECLARE @DeleteMark nvarchar(2)
SET @DeleteMark = N'~'                       -- 삭제 마크 세팅
DECLARE @StartDeleteMark int
SET @StartDeleteMark = CHARINDEX( @DeleteMark, @sID )
IF @StartDeleteMark > 0
	SET @sID = LEFT( @sID, @StartDeleteMark-1 )
-- 지워질 캐릭의 ID를 변경
DECLARE @CurrentDate datetime
SET @CurrentDate = GETDATE()
DECLARE @nDeleteCount int
SET @nDeleteCount = 0
DECLARE @sChangedID nvarchar(40)
NAME_CHANGE: -- 이름 만들기 '아뒤' + 삭제마크 + '날짜' + '#시도횟수'
SET @sChangedID = 
	RTRIM(@sID) + 
	@DeleteMark + DATENAME(year, @CurrentDate) + DATENAME(month, @CurrentDate) + DATENAME(day, @CurrentDate) +
	N'#' + CAST( @nDeleteCount AS nvarchar(2) )
-- 캐릭 지우기
UPDATE tCharacter SET 
bDeleted = 1,
dDeletedDate = GETDATE(),
sID = @sChangedID
WHERE nCharNo = @nCharNo
-- 실패시 50번 이름을 바꾸면서 시도한다.
IF @@ERROR = 2601 -- 인덱스키 중복인 경우 
BEGIN
	SET @nDeleteCount = @nDeleteCount + 1
	IF @nDeleteCount < 50
		GOTO NAME_CHANGE
		
	-- 실패	
	SET @nRet = 3 
	RETURN
END	
-- 성공
SET @nRet = 0 
-- end

0戁쒜ŀ
GO

-- ========================================================================
-- Procedure: p_Char_FameAdd
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_FameAdd]
/*
Add fame
2004.9 By CJC
 input:
	see SQL
 output:
	@nRet	1 = OK
	        0 = Error
*/
@nCharNo int,
@nAddFame int,
@nRet tinyint OUTPUT
AS
SET NOCOUNT ON
DECLARE @nCharFame int
EXEC p_Char_FameGet @nCharNo, @nCharFame OUTPUT, @nRet OUTPUT
IF @nRet = 0
	RETURN
-- Add fame
SET @nCharFame = @nCharFame + @nAddFame
EXEC p_Char_FameSet @nCharNo, @nCharFame, @nRet OUTPUT
-- End

0瀁쨷ń
GO

-- ========================================================================
-- Procedure: p_Char_FameGet
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_FameGet]
/*
Get fame
2004.9 By CJC
 input:
	see SQL
 output:
	@nRet	1 = OK, Fame return
	        0 = Error
*/
@nCharNo int,
@nCharFame int OUTPUT,
@nRet tinyint OUTPUT
AS
SET NOCOUNT ON
SET @nRet = 0
-- Get fame
SELECT @nCharFame = nFame FROM tCharacter WHERE nCharNo = @nCharNo
IF @@ROWCOUNT = 1
	SET @nRet = 1
	
-- End
RETURN

0꤁빛Ņ
GO

-- ========================================================================
-- Procedure: p_Char_FameSet
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_FameSet]
/*
Set fame
2004.9 By CJC
 input:
	see SQL
 output:
	@nRet	1 = OK
	        0 = Error
*/
@nCharNo int,
@nSetFame int,
@nRet tinyint OUTPUT
AS
SET NOCOUNT ON
SET @nRet = 0
UPDATE tCharacter SET nFame = @nSetFame WHERE nCharNo = @nCharNo
IF @@ERROR = 0 AND @@ROWCOUNT = 1
	SET @nRet = 1
-- end
RETURN

0ఁņ
GO

-- ========================================================================
-- Procedure: p_Char_FameSub
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_FameSub]
/*
Sub fame
2004.9 By CJC
 input:
	see SQL
 output:
	@nRet	1 = OK
	        0 = Error, or Minus value
*/
@nCharNo int,
@nSubFame int,
@nRet tinyint OUTPUT
AS
SET NOCOUNT ON
DECLARE @nCharFame int
EXEC p_Char_FameGet @nCharNo, @nCharFame OUTPUT, @nRet OUTPUT
IF @nRet = 0
	RETURN
-- Sub fame
SET @nCharFame = @nCharFame - @nSubFame
-- Check minus
IF @nCharFame < 0
BEGIN
	SET @nRet = 0
	RETURN	
END
EXEC p_Char_FameSet @nCharNo, @nCharFame, @nRet OUTPUT
-- End

0뉿ņ
GO

-- ========================================================================
-- Procedure: p_Char_Find
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_Find]
/*
Character Find
2004.6 By CJC
  input:
    sID = 찾을 ID
  output:
    0 = 없음
    1~? = 찾은 캐릭번호
*/
@sID nvarchar(40),
@nCharNo int OUTPUT
AS
SET NOCOUNT ON
SET @nCharNo = 0
SELECT @nCharNo = nCharNo FROM tCharacter(NOLOCK) WHERE sID = @sID
-- end

0⨁斾ň
GO

-- ========================================================================
-- Procedure: p_Char_GetAllData
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_GetAllData]
/*
Character Get all data
2004.6 By CJC
 input:
   nCharNo = To be get login data
 output:
	none
 desc:
	fetch success, get data
	fetch fail, not data
*/
@nCharNo int
AS
SET NOCOUNT ON

/*-------------------------------------------------------------------------------------
Project: FiestaHeroes - NA2016
Author: xkl

Rebirth of characters that can be customized through the settings in:
Options.dbo.Character_Rebirth
------------------------------------------------------------------------------------*/

-- Retrieve values from Character_Rebirth.
DECLARE @Trigger_ItemID INT, @Trigger_Level INT, @Reward_Abstate INT, @Unequip_Character_Items INT, @CharacterLevel INT;
SELECT @Trigger_ItemID = Trigger_ItemID, @Trigger_Level = Trigger_Level, @Reward_Abstate = Reward_Abstate, @Unequip_Character_Items = Unequip_Character_Items FROM Options..Character_Rebirth WHERE Enabled = 1;

-- Get character level.
SELECT @CharacterLevel = nLevel FROM tCharacter WHERE nCharNo = @nCharNo;

-- Update the rebirth Abstate time for characters that already have it.
UPDATE tCharacterState SET nStateTime = 100000000, nRestTime = 100000000 WHERE nCharNo = @nCharNo AND nStateNo = @Reward_Abstate;

-- Begin the character rebirth process.
IF EXISTS (SELECT 1 FROM Options..Character_Rebirth WHERE Enabled = 1)
BEGIN
	-- Check if the character has the required item.
	IF EXISTS (SELECT 1 FROM tItem WHERE nOwner = @nCharNo AND nItemID = @Trigger_ItemID)
	BEGIN
	    -- Check if the character meets the required level.
	    IF @CharacterLevel = @Trigger_Level
	    BEGIN
	        -- Apply rebirth.
	        DELETE FROM tItem WHERE nOwner = @nCharNo AND nItemID = @Trigger_ItemID;
	        UPDATE tCharacter SET nLevel = 1, nExp = 0 WHERE nCharNo = @nCharNo;
	
	        -- Check for unequip setting.
	        IF @Unequip_Character_Items = 1
	            UPDATE tItem SET nStorageType = 2 WHERE nOwner = @nCharNo AND nStorageType = 8;
	
	        -- Reward the Abstate.
	        IF NOT EXISTS (SELECT 1 FROM tCharacterState WHERE nCharNo = @nCharNo AND nStateNo = @Reward_Abstate)
	            INSERT INTO tCharacterState (nCharNo, nStateNo, nStateLevel, nStateTime, nRestTime) 
	            VALUES (@nCharNo, @Reward_Abstate, 1, 100000000, 100000000);
	    END
	END
END
/*-------------------------------------------------------------------------------------
 Continue as normal ~ 
------------------------------------------------------------------------------------*/

SELECT
nUserNo,
nSlotNo,
sID,
nLevel,
nExp,
nHP,
nSP,
nAP,
nHPS,
nSPS,
nPS,
nGS,
nMoney,
nFame,
nFlags,
nAdminLevel,
sLoginZone,
nLoginZoneX,
nLoginZoneY,
nLoginZoneD,
nPrisonMin,
nPKCount,
nPKYellowTime,
nStrength,
nConstitute,
nDexterity,
nIntelligence,
nMentalPower,
nRedistributePoint,
nPlayMin,
nLoginCount,
dLastLoginDate,
dCreateDate,
nCreateWorld,
bDeleted,
dDeletedDate
FROM tCharacter
WHERE nCharNo = @nCharNo
-- end
GO

-- ========================================================================
-- Procedure: p_Char_GetCharID
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_GetCharID]
/*
Get Char ID by Char No
input: CharNo
Output:
    @nRet, 0 = Not found char
           1 = @sID, Found char ID
*/
@nCharNo int,
@sID nvarchar(40) OUTPUT,
@nRet int OUTPUT
AS
SET NOCOUNT ON
SELECT @sID = sID FROM tCharacter WHERE nCharNo = @nCharNo
SET @nRet = @@ROWCOUNT
-- END

0㠁歙Ō
GO

-- ========================================================================
-- Procedure: p_Char_GetCharNo
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_GetCharNo]
/*
Get Char No by Char ID
Input: sID
Output:
    @nCharNo, 0 = Not found user
              ? = Found user No
*/
@sID nvarchar(40),
@nCharNo int OUTPUT
AS
SET NOCOUNT ON
SET @nCharNo = 0
SELECT @nCharNo = nCharNo FROM tCharacter WHERE sID = @sID
-- END

0鬁钇ō
GO

-- ========================================================================
-- Procedure: p_Char_GetChargedBufferAll
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_GetChargedBufferAll]
/*
Char Get charaged buffer all
2006.6 By CJC
 input:
	see SQL
 output:
	fetch to take data
*/

@nCharNo int

AS

SET NOCOUNT ON

SELECT nKey, nID, dUseDate, dEndDate
FROM tCharacterChargedBuffer WHERE nCharNo = @nCharNo

-- end

0鰁丆Ŋ
GO

-- ========================================================================
-- Procedure: p_Char_GetChargedBufferCount
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_GetChargedBufferCount]
/*
Char Get charaged buffer all count
2006.6 By CJC
 input:
	see SQL
 output:
	fetch to take data
*/

@nCharNo int,
@nCount int OUTPUT

AS

SET NOCOUNT ON

SELECT @nCount = COUNT(*) FROM tCharacterChargedBuffer WHERE nCharNo = @nCharNo

-- end

0！眴ŋ
GO

-- ========================================================================
-- Procedure: p_Char_GetIsDeleted
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_GetIsDeleted]
@nCharNo int,
@nRet int OUTPUT

AS
SET NOCOUNT ON

IF EXISTS(SELECT nCharNo FROM tCharacter(NOLOCK) WHERE nCharNo = @nCharNo AND bDeleted > 0)

BEGIN
	SET @nRet = 1
		RETURN
END

0焁彽ō
GO

-- ========================================================================
-- Procedure: p_Char_GetKQMap
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_GetKQMap]
/*
Character Get Kingdom quest map
2006.4 By CJC
 input:
   nCharNo = To be get KQ Map data
 output:
	none
 desc:
	fetch success, get data
	fetch fail, not data
*/

@nCharNo int

AS
SET NOCOUNT ON

SELECT

nKQHandle,
sKQMap,
nKQX,
nKQY,
nKQD,
dKQDate

FROM tCharacter

WHERE nCharNo = @nCharNo

-- end

0䟅ŏ
GO

-- ========================================================================
-- Procedure: p_Char_GetListOfUserChar
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_GetListOfUserChar]
/*
Character Get list of user character
2004.6 By CJC
  Input:
    nUserNo = Search user's character
  output: 
	none
 desc:
	fetch success, get data
	fetch fail, not data
*/
@nUserNo int
AS
SET NOCOUNT ON
SELECT nCharNo FROM tCharacter(NOLOCK) WHERE nUserNo = @nUserNo AND bDeleted = 0
-- end

0ᰁ㯪Ő
GO

-- ========================================================================
-- Procedure: p_Char_GetListOfUserCharAll
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_GetListOfUserCharAll]
/*
Character Get list of user character(include deleted character)
2004.6 By CJC
  Input:
    nUserNo = Search user's character
  output: 
	none
 desc:
	fetch success, get data
	fetch fail, not data
*/
@nUserNo int
AS
SET NOCOUNT ON
SELECT nCharNo FROM tCharacter(NOLOCK) WHERE nUserNo = @nUserNo
-- end

0唁『ő
GO

-- ========================================================================
-- Procedure: p_Char_GetLoginData
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_GetLoginData]
/*
Character Get login data
2004.6 By CJC
 input:
	nCharNo = To be get login data
 output:
	none
 desc:
	fetch success, get data
	fetch fail, not data
*/
@nCharNo int
AS
SET NOCOUNT ON
SELECT 
	C.sID, 
	C.nLevel, 
	C.nSlotNo, 
	C.sLoginZone, 
	C.bDeleted, 
	C.dDeletedDate,
	
	S.nRace, 
	S.nClass, 
	S.nGender, 
	S.nHairType, 
	S.nHairColor, 
	S.nFaceShape,
	
	C.nKQHandle,
	C.sKQMap,
	C.nKQX,
	C.nKQY,
	C.nKQD,
	C.dKQDate
FROM tCharacter AS C
INNER JOIN tCharacterShape AS S ON S.nCharNo = @nCharNo
WHERE C.nCharNo = @nCharNo
-- end

0踁␲Œ
GO

-- ========================================================================
-- Procedure: p_Char_GetNumOfUserChar
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_GetNumOfUserChar]
/*
Character number of user character
2004.6 By CJC
  Input:
    nUserNo = Search user's character
  output: 
    0 = none
    ? = number of user's character, execpt deleted.
*/
@nUserNo int,
@nNumOfUserChar int OUTPUT
AS
SET NOCOUNT ON
SET @nNumOfUserChar = 0
SELECT @nNumOfUserChar = COUNT(nCharNo) FROM tCharacter(NOLOCK) WHERE nUserNo = @nUserNo AND bDeleted = 0
-- end

0윁ᡖœ
GO

-- ========================================================================
-- Procedure: p_Char_GetNumOfUserCharAll
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_GetNumOfUserCharAll]
/*
Character number of user character(include deleted character)
2004.6 By CJC
  Input:
    nUserNo = Search user's character
  output: 
    0 = none
    ? = number of user's character
*/
@nUserNo int,
@nNumOfUserChar int OUTPUT
AS
SET NOCOUNT ON
SET @nNumOfUserChar = 0
SELECT @nNumOfUserChar = COUNT(nCharNo) FROM tCharacter(NOLOCK) WHERE nUserNo = @nUserNo
-- end

0⨁䆅Ŕ
GO

-- ========================================================================
-- Procedure: p_Char_GetOptGame
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_GetOptGame]

/*
Character Get options
2005.8 By CJC

 input:
   nCharNo = To be get option data
 output:
	@sData = Option binary data
*/

@nCharNo int,
@nRet tinyint output,
@sData varbinary(64) output

AS

SET NOCOUNT ON

SET @nRet = 0
SELECT @sData = sGame, @nRet = 1 FROM tCharacterOptions WHERE nCharNo = @nCharNo AND sGame IS NOT NULL

-- end

0౻Ŕ
GO

-- ========================================================================
-- Procedure: p_Char_GetOptKeyMapping
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_GetOptKeyMapping]

/*
Character Get options
2005.8 By CJC

 input:
   nCharNo = To be get option data
 output:
	@sData = Option binary data
*/

@nCharNo int,
@nRet tinyint output,
@sData varbinary(308) output

AS

SET NOCOUNT ON

SET @nRet = 0
SELECT @sData = sKeyMapping, @nRet = 1 FROM tCharacterOptions WHERE nCharNo = @nCharNo AND sKeyMapping IS NOT NULL

-- end

0䠁뾹ŕ
GO

-- ========================================================================
-- Procedure: p_Char_GetOptShortCutData
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_GetOptShortCutData]

/*
Character Get options
2005.8 By CJC

 input:
   nCharNo = To be get option data
 output:
	@sData = Option binary data
*/

@nCharNo int,
@nRet tinyint output,
@sData varbinary(1024) output

AS

SET NOCOUNT ON

SET @nRet = 0
SELECT @sData = sShortCutData, @nRet = 1 FROM tCharacterOptions WHERE nCharNo = @nCharNo AND sShortCutData IS NOT NULL

-- end

0脁돝Ŗ
GO

-- ========================================================================
-- Procedure: p_Char_GetOptShortCutSize
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_GetOptShortCutSize]

/*
Character Get options
2005.8 By CJC

 input:
   nCharNo = To be get option data
 output:
	@sData = Option binary data
*/

@nCharNo int,
@nRet tinyint output,
@sData varbinary(24) output

AS

SET NOCOUNT ON

SET @nRet = 0
SELECT @sData = sShortCutSize, @nRet = 1 FROM tCharacterOptions WHERE nCharNo = @nCharNo AND sShortCutSize IS NOT NULL

-- end

0ŗ
GO

-- ========================================================================
-- Procedure: p_Char_GetOptSound
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_GetOptSound]

/*
Character Get options
2005.8 By CJC

 input:
   nCharNo = To be get option data
 output:
	@sData = Option binary data
*/

@nCharNo int,
@nRet tinyint output,
@sData varbinary(1) output

AS

SET NOCOUNT ON

SET @nRet = 0
SELECT @sData = sSound, @nRet = 1 FROM tCharacterOptions WHERE nCharNo = @nCharNo AND sSound IS NOT NULL

-- end

0鰥Ř
GO

-- ========================================================================
-- Procedure: p_Char_GetOptVideo
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_GetOptVideo]

/*
Character Get options
2005.8 By CJC

 input:
   nCharNo = To be get option data
 output:
	@sData = Option binary data
*/

@nCharNo int,
@nRet tinyint output,
@sData varbinary(60) output

AS

SET NOCOUNT ON

SET @nRet = 0
SELECT @sData = sVideo, @nRet = 1 FROM tCharacterOptions WHERE nCharNo = @nCharNo AND sVideo IS NOT NULL

-- end

0Ⰱ遊ř
GO

-- ========================================================================
-- Procedure: p_Char_GetOptWindowsPos
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_GetOptWindowsPos]

/*
Character Get options
2005.8 By CJC

 input:
   nCharNo = To be get option data
 output:
	@sData = Option binary data
*/

@nCharNo int,
@nRet tinyint output,
@sData varbinary(392) output

AS

SET NOCOUNT ON

SET @nRet = 0
SELECT @sData = sWindowsPos, @nRet = 1 FROM tCharacterOptions WHERE nCharNo = @nCharNo AND sWindowsPos IS NOT NULL

-- end

0輁른Ś
GO

-- ========================================================================
-- Procedure: p_Char_GetShape
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_GetShape]
/*
Character Get shape
2004.6 By CJC
 input:
   nCharNo = To be get shape
 output:
	none
 desc:
	fetch success, get data
	fetch fail, not data
*/
@nCharNo int
AS
SET NOCOUNT ON
SELECT 
	nRace, 
	nClass, 
	nGender, 
	nHairType, 
	nHairColor, 
	nFaceShape 
FROM tCharacterShape WHERE nCharNo = @nCharNo
-- end

0적궜ś
GO

-- ========================================================================
-- Procedure: p_Char_GetState
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_GetState]
	@nCharNo int,
	@nStateNo smallint,
	@nStateLevel smallint OUTPUT,
	@nStateTime int OUTPUT,
	@nStateCount smallint OUTPUT,
	@nStatePowerDemage tinyint OUTPUT,
	@nStatePowerSP tinyint OUTPUT,
	@nStatePowerKeepTime tinyint OUTPUT,
	@nStatePowerCoolTime tinyint OUTPUT,
	@nCasterLevel smallint output,
	@nCasterAttackPower int output,		
	@nRet int OUTPUT
AS
BEGIN
	SET NOCOUNT ON

	SELECT 
		@nStateLevel         = nStateLevel, 
		@nStateTime          = nStateTime,
		@nStateCount         = nStateCount,	
		@nStatePowerDemage   = nStatePowerDemage,
		@nStatePowerSP       = nStatePowerSP,
		@nStatePowerKeepTime = nStatePowerKeepTime,
		@nStatePowerCoolTime = nStatePowerCoolTime,
		@nCasterLevel        = nCasterLevel,
		@nCasterAttackPower  = @nCasterAttackPower 
	  FROM tCharacterState 
	 WHERE nCharNo = @nCharNo AND nStateNo = @nStateNo
	SET @nRet = @@ROWCOUNT
END

0휁沶Ŝ
GO

-- ========================================================================
-- Procedure: p_Char_GetStateAll
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_GetStateAll]
	@nCharNo int
AS
BEGIN
	SET NOCOUNT ON
	SELECT 
		nStateNo, 
		nStateLevel, 
		nStateTime,
		nStateCount,
		nStatePowerDemage,
		nStatePowerSP,
		nStatePowerKeepTime,
		nStatePowerCoolTime,
		nCasterLevel,
		nCasterAttackPower 
	 FROM tCharacterState 
	WHERE nCharNo = @nCharNo
END

0㨁闥ŝ
GO

-- ========================================================================
-- Procedure: p_Char_GetStateAllWithRestTime
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_GetStateAllWithRestTime]
	@nCharNo int
AS
BEGIN
	SET NOCOUNT ON
	SELECT 
		nStateNo, 
		nStateLevel, 
		nStateTime,
		nRestTime,
		nStateCount,
		nStatePowerDemage,
		nStatePowerSP,
		nStatePowerKeepTime,
		nStatePowerCoolTime,
		nCasterLevel,
		nCasterAttackPower 		
	 FROM tCharacterState 
	WHERE nCharNo = @nCharNo
END

0䤁哿Ş
GO

-- ========================================================================
-- Procedure: p_Char_GetStateAllWithSetItem
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_GetStateAllWithSetItem]
	@nCharNo int
AS
BEGIN
	SET NOCOUNT ON
	SELECT 
		nStateNo, 
		nStateLevel, 
		nStateTime,
		nStateCount,
		nStatePowerDemage,
		nStatePowerSP,
		nStatePowerKeepTime,
		nStatePowerCoolTime,
		nSetItemIndexA,
		nSetItemValueA,
		nSetItemIndexB,
		nSetItemValueB,
		nSetItemIndexC,
		nSetItemValueC,
		nSetItemIndexD,
		nSetItemValueD,
		nCasterLevel,
		nCasterAttackPower		
	 FROM tCharacterState 
	WHERE nCharNo = @nCharNo
END

0舁䤣ş
GO

-- ========================================================================
-- Procedure: p_Char_GetStateAllWithSetItemAndRestTime
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_GetStateAllWithSetItemAndRestTime]
	@nCharNo int
AS
BEGIN
	SET NOCOUNT ON
	SELECT 
		nStateNo, 
		nStateLevel, 
		nStateTime,
		nRestTime,
		nStateCount,
		nStatePowerDemage,
		nStatePowerSP,
		nStatePowerKeepTime,
		nStatePowerCoolTime,
		nSetItemIndexA,
		nSetItemValueA,
		nSetItemIndexB,
		nSetItemValueB,
		nSetItemIndexC,
		nSetItemValueC,
		nSetItemIndexD,
		nSetItemValueD,
		nCasterLevel,
		nCasterAttackPower
	 FROM tCharacterState 
	WHERE nCharNo = @nCharNo
END

0묁㵇Š
GO

-- ========================================================================
-- Procedure: p_Char_GetStateWithRestTime
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_GetStateWithRestTime]
	@nCharNo int,
	@nStateNo smallint,
	@nStateLevel smallint OUTPUT,
	@nStateTime int OUTPUT,
	@nRestTime int OUTPUT,
	@nStateCount smallint OUTPUT,
	@nStatePowerDemage tinyint OUTPUT,
	@nStatePowerSP tinyint OUTPUT,
	@nStatePowerKeepTime tinyint OUTPUT,
	@nStatePowerCoolTime tinyint OUTPUT,
	@nCasterLevel smallint output,
	@nCasterAttackPower int output,			
	@nRet int OUTPUT
AS
BEGIN
	SET NOCOUNT ON
	SELECT 
		@nStateLevel			= nStateLevel, 
		@nStateTime			= nStateTime,
		@nRestTime			= nRestTime,
		@nStateCount			= nStateCount,	
		@nStatePowerDemage		= nStatePowerDemage,
		@nStatePowerSP			= nStatePowerSP,
		@nStatePowerKeepTime	= nStatePowerKeepTime,
		@nStatePowerCoolTime	= nStatePowerCoolTime,
		@nCasterLevel        = nCasterLevel,
		@nCasterAttackPower  = @nCasterAttackPower 		
	  FROM tCharacterState 
	 WHERE nCharNo = @nCharNo AND nStateNo = @nStateNo

	SET @nRet = @@ROWCOUNT
END

0ḁ晶š
GO

-- ========================================================================
-- Procedure: p_Char_GetStateWithSetItem
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_GetStateWithSetItem]
	@nCharNo int,
	@nStateNo smallint,
	@nStateLevel smallint OUTPUT,
	@nStateTime int OUTPUT,
	@nStateCount smallint OUTPUT,
	@nStatePowerDemage tinyint OUTPUT,
	@nStatePowerSP tinyint OUTPUT,
	@nStatePowerKeepTime tinyint OUTPUT,
	@nStatePowerCoolTime tinyint OUTPUT,
	@nSetItemIndexA smallint OUTPUT,
	@nSetItemValueA smallint OUTPUT,
	@nSetItemIndexB smallint OUTPUT,
	@nSetItemValueB smallint OUTPUT,
	@nSetItemIndexC smallint OUTPUT,
	@nSetItemValueC smallint OUTPUT,
	@nSetItemIndexD smallint OUTPUT,
	@nSetItemValueD smallint OUTPUT,
	@nCasterLevel smallint output,
	@nCasterAttackPower int output,			
	@nRet int OUTPUT
AS
BEGIN
	SET NOCOUNT ON
	SELECT 
		@nStateLevel			= nStateLevel, 
		@nStateTime			= nStateTime,
		@nStateCount			= nStateCount,	
		@nStatePowerDemage		= nStatePowerDemage,
		@nStatePowerSP		= nStatePowerSP,
		@nStatePowerKeepTime	= nStatePowerKeepTime,
		@nStatePowerCoolTime		= nStatePowerCoolTime,
		@nSetItemIndexA		= nSetItemIndexA,
		@nSetItemValueA		= nSetItemValueA,
		@nSetItemIndexB		= nSetItemIndexB,
		@nSetItemValueB		= nSetItemValueB,
		@nSetItemIndexC		= nSetItemIndexC,
		@nSetItemValueC		= nSetItemValueC,
		@nSetItemIndexD		= nSetItemIndexD,
		@nSetItemValueD		= nSetItemValueD,
		@nCasterLevel        = nCasterLevel,
		@nCasterAttackPower  = @nCasterAttackPower 		
	  FROM tCharacterState 
 	 WHERE nCharNo = @nCharNo AND nStateNo = @nStateNo

	SET @nRet = @@ROWCOUNT
END

0圁媚Ţ
GO

-- ========================================================================
-- Procedure: p_Char_GetStateWithSetItemAndRestTime
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_GetStateWithSetItemAndRestTime]
	@nCharNo int,
	@nStateNo smallint,
	@nStateLevel smallint OUTPUT,
	@nStateTime int OUTPUT,
	@nRestTime int OUTPUT,
	@nStateCount smallint OUTPUT,
	@nStatePowerDemage tinyint OUTPUT,
	@nStatePowerSP tinyint OUTPUT,
	@nStatePowerKeepTime tinyint OUTPUT,
	@nStatePowerCoolTime tinyint OUTPUT,
	@nSetItemIndexA smallint OUTPUT,
	@nSetItemValueA smallint OUTPUT,
	@nSetItemIndexB smallint OUTPUT,
	@nSetItemValueB smallint OUTPUT,
	@nSetItemIndexC smallint OUTPUT,
	@nSetItemValueC smallint OUTPUT,
	@nSetItemIndexD smallint OUTPUT,
	@nSetItemValueD smallint OUTPUT,
	@nCasterLevel smallint output,
	@nCasterAttackPower int output,			
	@nRet int OUTPUT
AS
BEGIN
	SET NOCOUNT ON
	SELECT 
		@nStateLevel			= nStateLevel, 
		@nStateTime			= nStateTime,
		@nRestTime			= nRestTime,
		@nStateCount			= nStateCount,	
		@nStatePowerDemage		= nStatePowerDemage,
		@nStatePowerSP		= nStatePowerSP,
		@nStatePowerKeepTime	= nStatePowerKeepTime,
		@nStatePowerCoolTime		= nStatePowerCoolTime,
		@nSetItemIndexA		= nSetItemIndexA,
		@nSetItemValueA		= nSetItemValueA,
		@nSetItemIndexB		= nSetItemIndexB,
		@nSetItemValueB		= nSetItemValueB,
		@nSetItemIndexC		= nSetItemIndexC,
		@nSetItemValueC		= nSetItemValueC,
		@nSetItemIndexD		= nSetItemIndexD,
		@nSetItemValueD		= nSetItemValueD,
		@nCasterLevel        = nCasterLevel,
		@nCasterAttackPower  = @nCasterAttackPower 		
	  FROM tCharacterState 
	 WHERE nCharNo = @nCharNo AND nStateNo = @nStateNo

	SET @nRet = @@ROWCOUNT
END

0送亾ţ
GO

-- ========================================================================
-- Procedure: p_Char_GetTitleAll
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_GetTitleAll]

/*
Character title set
2005.10 By CJC
  Input:
    nCharNo = character number
  output: 
    featch to get data
*/

@nCharNo int

AS

SET NOCOUNT ON

SELECT nType, nElementNo, nElementValue 
  FROM tCharacterTitle WITH(READUNCOMMITTED)
 WHERE nCharNo = @nCharNo ORDER BY nNo  

-- end

0矬Ť
GO

-- ========================================================================
-- Procedure: p_Char_GetTitleCount
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_GetTitleCount]

/*
Character number of title
2005.10 By CJC
  Input:
    nCharNo = character number
  output: 
    0 = none
    ? = number of title on this character
*/

@nCharNo int,
@nCount int OUTPUT

AS

SET NOCOUNT ON

SET @nCount = 0

SELECT @nCount = COUNT(*) FROM tCharacterTitle(NOLOCK) WHERE nCharNo = @nCharNo

-- end

0줁䋢Ť
GO

-- ========================================================================
-- Procedure: p_Char_GetUserNo
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_GetUserNo]
/*
Get User No by Char ID
Input: nCharNo
Output:
    @nUserNo, 0 = Not found user
              ? = Found user No
*/
@nCharNo int,
@nUserNo int OUTPUT
AS
SET NOCOUNT ON
SET @nUserNo = 0
SELECT @nUserNo = nUserNo FROM tCharacter WHERE nCharNo = @nCharNo
-- END

0ᄁť
GO

-- ========================================================================
-- Procedure: p_Char_Login
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_Login]
	@nCharNo	int
,	@nRet		int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	DECLARE @dLastloginDate AS datetime
	
	SET @nRet = 0

	IF EXISTS (SELECT nCharNo FROM tCharacter WHERE nCharNo = @nCharNo) BEGIN
		BEGIN TRAN

		SET @dLastloginDate = (SELECT dLastLoginDate FROM tCharacter WHERE nCharNo = @nCharNo)

		UPDATE tCharacter SET nLoginCount = nLoginCount + 1, dLastLoginDate = GETDATE() WHERE nCharNo = @nCharNo
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = 1
			IF @@TRANCOUNT > 0 ROLLBACK TRAN
			RETURN
		END
		ELSE BEGIN
			SELECT DATEDIFF(d,@dLastLoginDate,GETDATE()) AS nLoginDateDiff
			  FROM tCharacter
			 WHERE nCharNo = @nCharNo

			IF @@TRANCOUNT > 0 COMMIT TRAN
			RETURN
		END
	END
	ELSE BEGIN
		SET @nRet = -1111
	END

	SET LOCK_TIMEOUT -1
	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: p_Char_Logout
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_Logout]
	@nCharNo	int
,	@nPlayMin	int
,	@nRet		int OUTPUT
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	DECLARE @nLastFriendPoint int

	SET @nRet = 1
	SET @nLastFriendPoint = (SELECT nFriendPoint FROM tCharacter WHERE nCharNo = @nCharNo)
	

	IF EXISTS (SELECT nCharNo FROM tCharacter WHERE nCharNo = @nCharNo) BEGIN
		BEGIN TRAN

		UPDATE tCharacter SET nPlayMin = nPlayMin + @nPlayMin WHERE nCharNo = @nCharNo
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = 0
			IF @@TRANCOUNT > 0 ROLLBACK TRAN
			RETURN
		END
		ELSE BEGIN
			UPDATE tCharacter SET nLastFriendPoint = @nLastFriendPoint WHERE nCharNo = @nCharNo
			IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
				SET @nRet = 0
				IF @@TRANCOUNT > 0 ROLLBACK TRAN
				RETURN
			END
			ELSE BEGIN
				IF @@TRANCOUNT > 0 COMMIT TRAN
				RETURN
			END
		END
	END
	ELSE BEGIN
		SET @nRet = 0
	END

	SET LOCK_TIMEOUT -1
	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: p_Char_MoneyAdd
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_MoneyAdd]
/*
Add money
2004.8 By CJC
 input:
	see SQL
 output:
	@nRet	1 = OK
	        0 = Error
*/
@nCharNo int,
@nAddMoney bigint,
@nRet tinyint OUTPUT
AS
SET NOCOUNT ON
DECLARE @nCharMoney bigint
EXEC p_Char_MoneyGet @nCharNo, @nCharMoney OUTPUT, @nRet OUTPUT
IF @nRet = 0
	RETURN
-- Add money
SET @nCharMoney = @nCharMoney + @nAddMoney
EXEC p_Char_MoneySet @nCharNo, @nCharMoney, @nRet OUTPUT
-- End

0ခ㲢ũ
GO

-- ========================================================================
-- Procedure: p_Char_MoneyGet
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_MoneyGet]
/*
Get money
2004.8 By CJC
 input:
	see SQL
 output:
	@nRet	1 = OK, Money return
	        0 = Error
*/
@nCharNo int,
@nCharMoney bigint OUTPUT,
@nRet tinyint OUTPUT
AS
SET NOCOUNT ON
SET @nRet = 0
-- 돈 얻기
SELECT @nCharMoney = nMoney FROM tCharacter WHERE nCharNo = @nCharNo
IF @@ROWCOUNT = 1
	SET @nRet = 1
	
-- End
RETURN

0괁፳Ũ
GO

-- ========================================================================
-- Procedure: p_Char_MoneyGetCharsByUserNo
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_MoneyGetCharsByUserNo]
	@nUserNo		int
AS

BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	SELECT nCharNo, sID, nMoney, nLevel, nSlotNo, bDeleted, nAdminLevel 
	FROM tCharacter WHERE nUserNo =  @nUserNo
	
	SET LOCK_TIMEOUT -1
	SET NOCOUNT OFF	
END

0⸁뫖Ū
GO

-- ========================================================================
-- Procedure: p_Char_MoneyGetGreaterThen
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_MoneyGetGreaterThen]
	@nMoney		bigint
AS

BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	SELECT nUserNo, nCharNo, sID, nMoney, nLevel, nSlotNo, bDeleted, nAdminLevel 
	FROM tCharacter WHERE nMoney >=  @nMoney
	
	SET LOCK_TIMEOUT -1
	SET NOCOUNT OFF	
END

0鄁ū
GO

-- ========================================================================
-- Procedure: p_Char_MoneyMove
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_MoneyMove]
/*
Move money
2008.2 By CJC
 input:
	see SQL
 output:
	@nRet	1 = OK
	        0 = Error
*/
	@nCharNoFrom int,
	@nCharNoTo   int,
	@nMoveMoney  bigint,
	@nRet        tinyint OUTPUT
AS
BEGIN
	SET NOCOUNT ON

	IF @nMoveMoney < 0
		BEGIN
			SET @nRet = 0
			RETURN
		END

	IF @nMoveMoney = 0
		BEGIN
			SET @nRet = 1
			RETURN
		END
		
	DECLARE @money as bigint

	BEGIN TRAN

	IF NOT EXISTS(SELECT nCharNo FROM tCharacter WHERE nCharNo = @nCharNoFrom AND bDeleted = 0)
	BEGIN
		SET @nRet = 0
		ROLLBACK TRAN
		RETURN
	END

	IF NOT EXISTS(SELECT nCharNo FROM tCharacter WHERE nCharNo = @nCharNoTo AND bDeleted = 0)
	BEGIN
		SET @nRet = 0
		ROLLBACK TRAN
		RETURN
	END

	SELECT @money = nMoney FROM tCharacter WHERE nCharNo = @nCharNoFrom AND bDeleted = 0
	IF @money >= @nMoveMoney
		BEGIN
			-- SUB MONEY
			UPDATE tCharacter SET nMoney = @money - @nMoveMoney WHERE nCharNo = @nCharNoFrom AND bDeleted = 0
			IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
			BEGIN
				SET @nRet =0
				ROLLBACK TRAN
				RETURN
			END
			
			-- ADD MONEY
			UPDATE tCharacter SET nMoney = nMoney + @nMoveMoney WHERE nCharNo = @nCharNoTo AND bDeleted = 0
			IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
			BEGIN
				SET @nRet = 0
				ROLLBACK TRAN
				RETURN
			END
	
			SET @nRet = 1
			COMMIT TRAN
		END
	ELSE
		BEGIN
			SET @nRet = 0
			ROLLBACK TRAN
			RETURN
		END

	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: p_Char_MoneySet
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_MoneySet]
/*
Set money
2004.8 By CJC
 input:
	see SQL
 output:
	@nRet	1 = OK
	        0 = Error
*/
@nCharNo int,
@nSetMoney bigint,
@nRet tinyint OUTPUT
AS
SET NOCOUNT ON
SET @nRet = 0
UPDATE tCharacter SET nMoney = @nSetMoney WHERE nCharNo = @nCharNo
IF @@ERROR = 0 AND @@ROWCOUNT = 1
	SET @nRet = 1
-- end
RETURN

0쨁Ŭ
GO

-- ========================================================================
-- Procedure: p_Char_MoneySub
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_MoneySub]
/*
Sub money
2004.8 By CJC
 input:
	see SQL
 output:
	@nRet	1 = OK
	        0 = Error, or Minus value
*/
@nCharNo int,
@nSubMoney bigint,
@nRet tinyint OUTPUT
AS
SET NOCOUNT ON
DECLARE @nCharMoney bigint
EXEC p_Char_MoneyGet @nCharNo, @nCharMoney OUTPUT, @nRet OUTPUT
IF @nRet = 0
	RETURN
-- Sub money
SET @nCharMoney = @nCharMoney - @nSubMoney
-- Check minus
IF @nCharMoney < 0
BEGIN
	SET @nRet = 0
	RETURN	
END
EXEC p_Char_MoneySet @nCharNo, @nCharMoney, @nRet OUTPUT
-- End

0ሁ譧Ů
GO

-- ========================================================================
-- Procedure: p_Char_Rename
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_Rename]
/*
Char Rename
2006.9 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK, 
	        ? = Error code
*/
@nCharNo int,
@sID nvarchar(40),
-- Output var
@nRet int OUTPUT

AS
BEGIN

SET NOCOUNT ON

IF NOT EXISTS(SELECT nCharNo FROM tCharacterChangeID WHERE bFlag = 0 AND nCharNo = @nCharNo) BEGIN
	SET @nRet = @@ERROR
	RETURN
END

UPDATE tCharacter SET sID = @sID WHERE nCharNo = @nCharNo
SET @nRet = @@ERROR

END
-- end
GO

-- ========================================================================
-- Procedure: p_Char_SaveLevelExpFrame
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_SaveLevelExpFrame] 
@nCharNo int,
@nLevel smallint,
@nExp bigint,
@nFame int

AS
SET NOCOUNT ON

UPDATE tCharacter 
	SET 
		nLevel=@nLevel,
		nExp=@nExp,
		nFame=@nFame
	WHERE nCharNo=@nCharNo

return @@error

0긁ꢹŰ
GO

-- ========================================================================
-- Procedure: p_Char_SaveLocation
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_SaveLocation] 
@nCharNo		int,
@sLoginZone		nvarchar(16),
@nLoginZoneX	int,
@nLoginZoneY	int,
@nKQHandle		int,
@sKQMap			nvarchar(16),
@nKQX			int,
@nKQY			int

AS
SET NOCOUNT ON

UPDATE tCharacter SET 
	sLoginZone=@sLoginZone,
	nLoginZoneX=@nLoginZoneX,
	nLoginZoneY=@nLoginZoneY,
	
	nKQHandle=@nKQHandle,
	sKQMap=@sKQMap,
	nKQX=@nKQX,
	nKQY=@nKQY,
	dKQDate=GetDate()
	WHERE nCharNo=@nCharNo

return @@error

0鳝ű
GO

-- ========================================================================
-- Procedure: p_Char_SavePKCount
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_SavePKCount] 
/*
Character Set nPKCount
2007.2 By CJC
 input:
   All args
 output: 
   @@error
*/
@nCharNo int,
@nPKCount int

AS
SET NOCOUNT ON

UPDATE tCharacter SET nPKCount=@nPKCount WHERE nCharNo=@nCharNo

return @@error

0䨁옌Ų
GO

-- ========================================================================
-- Procedure: p_Char_SaveStat
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_SaveStat] 
@nCharNo int,
@nPS int,
@nGS int,
@nHPS smallint,
@nSPS smallint,
@nHP smallint,
@nSP smallint

AS
SET NOCOUNT ON

UPDATE tCharacter 
	SET 
		nPS=@nPS,
		nGS=@nGS,
		nHPS=@nHPS,
		nSPS=@nSPS,
		nHP=@nHP,
		nSP=@nSP 
	WHERE nCharNo=@nCharNo


return @@error

0 鄂Ų
GO

-- ========================================================================
-- Procedure: p_Char_SetAdminLevel
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_SetAdminLevel]
@nCharNo	int,
@nAdmin		Tinyint,

@nRet		int OUTPUT

AS

SET NOCOUNT ON

UPDATE tCharacter SET 
	nAdminLevel			= @nAdmin
WHERE nCharNo = @nCharNo

IF @@ERROR <> 0 OR @@ROWCOUNT = 0
BEGIN
	SET @nRet = 0
	RETURN
END
SET @nRet = 1

0鈁祊Ŵ
GO

-- ========================================================================
-- Procedure: p_Char_SetAllData
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_SetAllData]
/*
Character Set all data
2004.6 By CJC
 input:
   All args
 output: 
   1 = update ok
   0 = error at update chatacter data
*/
-- Character columns
@nCharNo int,
@nUserNo int, 
@nSlotNo tinyint,
@sID nvarchar(40), 
@nLevel smallint, 
@nExp bigint, 
@nHP smallint, 
@nSP smallint, 
@nAP smallint, 
@nHPS smallint,
@nSPS smallint,
@nPS int,
@nGS int,
@nMoney bigint,
@nFame int, 
@nFlags int, 
@nAdminLevel tinyint, 
@sLoginZone varchar(16),
@nLoginZoneX int,
@nLoginZoneY int,
@nLoginZoneD tinyint,
@nPrisonMin smallint, 
@nPKCount int, 
@nPKYellowTime tinyint,
@nStrength tinyint,
@nConstitute tinyint,
@nDexterity tinyint,
@nIntelligence tinyint,
@nMentalPower tinyint,
@nRedistributePoint tinyint,
@nPlayMin int, 
@nLoginCount int, 
@dLastLoginDate datetime, 
@dCreateDate datetime, 
@nCreateWorld tinyint, 
@bDeleted tinyint, 
@dDeletedDate datetime,
-- Output var
@nRet int OUTPUT
AS
SET NOCOUNT ON
-- Update character
UPDATE tCharacter SET
	nUserNo		=	@nUserNo , 
	nSlotNo			=	@nSlotNo , 
	sID			=	@sID , 
	nLevel			=	@nLevel , 
	nExp			=	@nExp , 
	nHP			=	@nHP , 
	nSP			=	@nSP , 
	nAP			=	@nAP , 
	nHPS			=	@nHPS , 
	nSPS			=	@nSPS , 
	nPS 			=	@nPS,
	nGS			=	@nGS,
	
	nMoney			= 	@nMoney , 
	nFame			= 	@nFame , 
	nFlags			= 	@nFlags , 
	nAdminLevel		=	@nAdminLevel ,
	
	sLoginZone		=	@sLoginZone , 
	nLoginZoneX		=	@nLoginZoneX , 
	nLoginZoneY		=	@nLoginZoneY , 
	nLoginZoneD		=	@nLoginZoneD ,
	
	nPrisonMin		=	@nPrisonMin , 
	nPKCount		=	@nPKCount , 
	nPKYellowTime		=	@nPKYellowTime ,
	
	nStrength		= 	@nStrength,
	nConstitute		= 	@nConstitute,
	nDexterity		= 	@nDexterity,
	nIntelligence		= 	@nIntelligence,
	nMentalPower		= 	@nMentalPower,
	nRedistributePoint	= 	@nRedistributePoint,
	
	nPlayMin		=	@nPlayMin , 
	nLoginCount		=	@nLoginCount , 
	dLastLoginDate		=	@dLastLoginDate ,
	
	dCreateDate		=	@dCreateDate , 
	nCreateWorld		=	@nCreateWorld ,
	bDeleted		=	@bDeleted, 
	dDeletedDate		=	@dDeletedDate 
WHERE nCharNo = @nCharNo
IF @@ERROR <> 0 OR @@ROWCOUNT = 0
BEGIN
	SET @nRet = 0
	RETURN
END
-- Update success
SET @nRet = 1
-- end

0ꄁ㡤ŵ
GO

-- ========================================================================
-- Procedure: p_Char_SetChargedBuffer
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_SetChargedBuffer]
/*
Char Set charged buffer
2006.6 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK, 
	        ? = Error code
*/

@nCharNo int,
@nKey int,
@nID smallint,
@dUseDate datetime,
@dEndDate datetime,

@nRet int OUTPUT

AS

SET NOCOUNT ON

IF EXISTS(SELECT nCharNo FROM tCharacterChargedBuffer WHERE nCharNo = @nCharNo AND nKey = @nKey )
BEGIN
	UPDATE tCharacterChargedBuffer SET 
		nID = @nID,
		dUseDate = @dUseDate,
		dEndDate = @dEndDate
	WHERE nCharNo = @nCharNo AND nKey = @nKey
END
ELSE
BEGIN
	INSERT tCharacterChargedBuffer (  nCharNo,  nKey,  nID,  dUseDate,  dEndDate )
	VALUES                         ( @nCharNo, @nKey, @nID, @dUseDate, @dEndDate )
END

SET @nRet = @@ERROR
RETURN

-- end

0ⲈŶ
GO

-- ========================================================================
-- Procedure: p_Char_SetLevelExp
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_SetLevelExp]
@nCharNo	int,
@nLevel		smallint,
@nExp		bigint,

@nRet		int OUTPUT

AS

SET NOCOUNT ON

UPDATE tCharacter SET 
	nLevel			= @nLevel,
	nExp			= @nExp
WHERE nCharNo = @nCharNo

IF @@ERROR <> 0 OR @@ROWCOUNT = 0
BEGIN
	SET @nRet = 0
	RETURN
END
SET @nRet = 1

0朁諁ŷ
GO

-- ========================================================================
-- Procedure: p_Char_SetOptGame
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_SetOptGame]
/*
Character Set options
2004.6 By CJC
 input:
   All args
 output: 
   0 = error
   1 = update ok
*/

@nCharNo int,
@sData varbinary(64),
@nRet int OUTPUT
AS

SET NOCOUNT ON

UPDATE tCharacterOptions SET sGame = @sData WHERE nCharNo = @nCharNo

IF @@ERROR <> 0 OR @@ROWCOUNT = 0
BEGIN
	INSERT tCharacterOptions ( nCharNo, sGame ) VALUES ( @nCharNo, @sData )
	SET @nRet = @@ROWCOUNT
	RETURN
END

SET @nRet = 1

-- end

0㴁喷ŷ
GO

-- ========================================================================
-- Procedure: p_Char_SetOptKeyMapping
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_SetOptKeyMapping]
/*
Character Set options
2004.6 By CJC
 input:
   All args
 output: 
   0 = error
   1 = update ok
*/

@nCharNo int,
@sData varbinary(308),
@nRet int OUTPUT
AS

SET NOCOUNT ON

UPDATE tCharacterOptions SET sKeyMapping = @sData WHERE nCharNo = @nCharNo

IF @@ERROR <> 0 OR @@ROWCOUNT = 0
BEGIN
	INSERT tCharacterOptions (nCharNo, sKeyMapping) VALUES (@nCharNo, @sData)
	SET @nRet = @@ROWCOUNT
	RETURN
END

SET @nRet = 1

-- end

0猉Ź
GO

-- ========================================================================
-- Procedure: p_Char_SetOptShortCutData
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_SetOptShortCutData]
/*
Character Set options
2004.6 By CJC
 input:
   All args
 output: 
   0 = error
   1 = update ok
*/

@nCharNo int,
@sData varbinary(1024),
@nRet int OUTPUT
AS

SET NOCOUNT ON

UPDATE tCharacterOptions SET sShortCutData = @sData WHERE nCharNo = @nCharNo

IF @@ERROR <> 0 OR @@ROWCOUNT = 0
BEGIN
	INSERT tCharacterOptions (nCharNo, sShortCutData) VALUES (@nCharNo, @sData)
	SET @nRet = @@ROWCOUNT
	RETURN
END

SET @nRet = 1

-- end

0꼁㷿Ź
GO

-- ========================================================================
-- Procedure: p_Char_SetOptShortCutSize
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_SetOptShortCutSize]
/*
Character Set options
2004.6 By CJC
 input:
   All args
 output: 
   0 = error
   1 = update ok
*/

@nCharNo int,
@sData varbinary(24),
@nRet int OUTPUT
AS

SET NOCOUNT ON

UPDATE tCharacterOptions SET sShortCutSize = @sData WHERE nCharNo = @nCharNo

IF @@ERROR <> 0 OR @@ROWCOUNT = 0
BEGIN
	INSERT tCharacterOptions (nCharNo, sShortCutSize) VALUES (@nCharNo, @sData)
	SET @nRet = @@ROWCOUNT
	RETURN
END

SET @nRet = 1

-- end

0℁♈Ż
GO

-- ========================================================================
-- Procedure: p_Char_SetOptSound
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_SetOptSound]
/*
Character Set options
2004.6 By CJC
 input:
   All args
 output: 
   0 = error
   1 = update ok
*/

@nCharNo int,
@sData varbinary(1),
@nRet int OUTPUT
AS

SET NOCOUNT ON

UPDATE tCharacterOptions SET sSound = @sData WHERE nCharNo = @nCharNo

IF @@ERROR <> 0 OR @@ROWCOUNT = 0
BEGIN
	INSERT tCharacterOptions (nCharNo, sSound) VALUES (@nCharNo, @sData)
	SET @nRet = @@ROWCOUNT
	RETURN
END

SET @nRet = 1

-- end

0威ᩬż
GO

-- ========================================================================
-- Procedure: p_Char_SetOptVideo
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_SetOptVideo]
/*
Character Set options
2004.6 By CJC
 input:
   All args
 output: 
   0 = error
   1 = update ok
*/

@nCharNo int,
@sData varbinary(60),
@nRet int OUTPUT
AS

SET NOCOUNT ON

UPDATE tCharacterOptions SET sVideo = @sData WHERE nCharNo = @nCharNo

IF @@ERROR <> 0 OR @@ROWCOUNT = 0
BEGIN
	INSERT tCharacterOptions (nCharNo, sVideo) VALUES (@nCharNo, @sData)
	SET @nRet = @@ROWCOUNT
	RETURN
END

SET @nRet = 1

-- end

0봁䎚Ž
GO

-- ========================================================================
-- Procedure: p_Char_SetOptWindowsPos
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_SetOptWindowsPos]
/*
Character Set options
2004.6 By CJC
 input:
   All args
 output: 
   0 = error
   1 = update ok
*/

@nCharNo int,
@sData varbinary(392),
@nRet int OUTPUT
AS

SET NOCOUNT ON

UPDATE tCharacterOptions SET sWindowsPos = @sData WHERE nCharNo = @nCharNo

IF @@ERROR <> 0 OR @@ROWCOUNT = 0
BEGIN
	INSERT tCharacterOptions (nCharNo, sWindowsPos) VALUES (@nCharNo, @sData)
	SET @nRet = @@ROWCOUNT
	RETURN
END

SET @nRet = 1

-- end

0錁ຐŽ
GO

-- ========================================================================
-- Procedure: p_Char_SetShape
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_SetShape]
/*
Character Set shape
2004.6 By CJC
 input:
   All args
 output: 
   0 = error
   1 = update ok
*/
@nCharNo int,
@nRace tinyint, 
@nClass tinyint, 
@nGender tinyint, 
@nHairType tinyint, 
@nHairColor tinyint, 
@nFaceShape tinyint, 
-- Output var
@nRet tinyint OUTPUT
AS
SET NOCOUNT ON
UPDATE tCharacterShape SET
	nRace		= @nRace,
	nClass		= @nClass,
	nGender		= @nGender,
	nHairType	= @nHairType,
	nHairColor	= @nHairColor,
	nFaceShape	= @nFaceShape
WHERE nCharNo = @nCharNo
IF @@ERROR <> 0 OR @@ROWCOUNT = 0
BEGIN
	SET @nRet = 0
	RETURN
END
SET @nRet = 1
-- end

0ԁž
GO

-- ========================================================================
-- Procedure: p_Char_SetState
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_SetState]
	@nCharNo int,
	@nStateNo smallint,
	@nStateLevel smallint,
	@nStateTime int,
	@nStateCount smallint,
	@nStatePowerDemage tinyint,
	@nStatePowerSP tinyint,
	@nStatePowerKeepTime tinyint,
	@nStatePowerCoolTime tinyint,
	@nCasterLevel smallint,
	@nCasterAttackPower int,
	@nRet int OUTPUT
AS
BEGIN
	SET NOCOUNT ON
	IF EXISTS(SELECT nCharNo FROM tCharacterState WHERE nCharNo = @nCharNo AND nStateNo = @nStateNo ) BEGIN
		UPDATE tCharacterState 
		   SET 
			nStateLevel			= @nStateLevel,
			nStateTime			= @nStateTime,
			nStateCount			= @nStateCount,	
			nStatePowerDemage	= @nStatePowerDemage,
			nStatePowerSP		= @nStatePowerSP,
			nStatePowerKeepTime	= @nStatePowerKeepTime,
			nStatePowerCoolTime	= @nStatePowerCoolTime,
		    nCasterLevel         = @nCasterLevel,
		    nCasterAttackPower   = @nCasterAttackPower 				
		WHERE nCharNo = @nCharNo AND nStateNo = @nStateNo
		
		SET @nRet = @@ERROR
	END
	ELSE BEGIN
		INSERT tCharacterState (  nCharNo,  nStateNo,  nStateLevel,  nStateTime,  nStateCount,  nStatePowerDemage,  nStatePowerSP,  nStatePowerKeepTime,  nStatePowerCoolTime, nCasterLevel, nCasterAttackPower )
		VALUES                 ( @nCharNo, @nStateNo, @nStateLevel, @nStateTime, @nStateCount, @nStatePowerDemage, @nStatePowerSP, @nStatePowerKeepTime, @nStatePowerCoolTime, @nCasterLevel, @nCasterAttackPower )

		SET @nRet = @@ERROR
	END
END

0ᐁ뗳ſ
GO

-- ========================================================================
-- Procedure: p_Char_SetStateWithRestTime
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_SetStateWithRestTime]
	@nCharNo int,
	@nStateNo smallint,
	@nStateLevel smallint,
	@nStateTime int,
	@nRestTime int,
	@nStateCount smallint,
	@nStatePowerDemage tinyint,
	@nStatePowerSP tinyint,
	@nStatePowerKeepTime tinyint,
	@nStatePowerCoolTime tinyint,
	@nCasterLevel smallint,
	@nCasterAttackPower int,	
	@nRet int OUTPUT
AS
BEGIN
	SET NOCOUNT ON

	IF EXISTS(SELECT nCharNo FROM tCharacterState WHERE nCharNo = @nCharNo AND nStateNo = @nStateNo ) BEGIN
		UPDATE tCharacterState SET 
			nStateLevel			= @nStateLevel,
			nStateTime			= @nStateTime,
			nRestTime			= @nRestTime,
			nStateCount			= @nStateCount,	
			nStatePowerDemage	= @nStatePowerDemage,
			nStatePowerSP		= @nStatePowerSP,
			nStatePowerKeepTime	= @nStatePowerKeepTime,
			nStatePowerCoolTime	= @nStatePowerCoolTime,
		    nCasterLevel         = @nCasterLevel,
		    nCasterAttackPower   = @nCasterAttackPower 					
		WHERE nCharNo = @nCharNo AND nStateNo = @nStateNo
		
		SET @nRet = @@ERROR
	END
	ELSE BEGIN
		INSERT tCharacterState (  nCharNo,  nStateNo,  nStateLevel,  nStateTime, nRestTime,  nStateCount,  nStatePowerDemage,  nStatePowerSP,  nStatePowerKeepTime,  nStatePowerCoolTime, nCasterLevel, nCasterAttackPower )
		VALUES                 ( @nCharNo, @nStateNo, @nStateLevel, @nStateTime, @nRestTime, @nStateCount, @nStatePowerDemage, @nStatePowerSP, @nStatePowerKeepTime, @nStatePowerCoolTime, @nCasterLevel, @nCasterAttackPower )

		SET @nRet = @@ERROR
	END
END

0䴁ꨗĀ
GO

-- ========================================================================
-- Procedure: p_Char_SetStateWithSetItem
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_SetStateWithSetItem]
	@nCharNo int,
	@nStateNo smallint,
	@nStateLevel smallint,
	@nStateTime int,
	@nStateCount smallint,
	@nStatePowerDemage tinyint,
	@nStatePowerSP tinyint,
	@nStatePowerKeepTime tinyint,
	@nStatePowerCoolTime tinyint,
	@nSetItemIndexA smallint,
	@nSetItemValueA smallint,
	@nSetItemIndexB smallint,
	@nSetItemValueB smallint,
	@nSetItemIndexC smallint,
	@nSetItemValueC smallint,
	@nSetItemIndexd smallint,
	@nSetItemValueD smallint,
	@nCasterLevel smallint,
	@nCasterAttackPower int,
	@nRet int = 0 OUTPUT
AS
BEGIN
	SET NOCOUNT ON
	
	IF EXISTS(SELECT nCharNo FROM tCharacterState WHERE nCharNo = @nCharNo AND nStateNo = @nStateNo) BEGIN
		UPDATE tCharacterState 
		   SET nStateLevel			= @nStateLevel
		     , nStateTime			= @nStateTime
		     , nStateCount			= @nStateCount
		     , nStatePowerDemage	= @nStatePowerDemage
		     , nStatePowerSP		= @nStatePowerSP
		     , nStatePowerKeepTime	= @nStatePowerKeepTime
		     , nStatePowerCoolTime	= @nStatePowerCoolTime
		     , nSetItemIndexA		= @nSetItemIndexA
		     , nSetItemValueA		= @nSetItemValueA
		     , nSetItemIndexB		= @nSetItemIndexB
		     , nSetItemValueB		= @nSetItemValueB
		     , nSetItemIndexC		= @nSetItemIndexC
		     , nSetItemValueC		= @nSetItemValueC
		     , nSetItemIndexD		= @nSetItemIndexD
		     , nSetItemValueD		= @nSetItemValueD
		     , nCasterLevel         = @nCasterLevel 
		     , nCasterAttackPower   = @nCasterAttackPower 
		 WHERE nCharNo = @nCharNo AND nStateNo = @nStateNo

		SET @nRet = @@ERROR
	END
	ELSE BEGIN
		INSERT tCharacterState (nCharNo, nStateNo, nStateLevel, nStateTime, nStateCount, nStatePowerDemage, nStatePowerSP, nStatePowerKeepTime, nStatePowerCoolTime, nSetItemIndexA, nSetItemValueA, nSetItemIndexB, nSetItemValueB, nSetItemIndexC, nSetItemValueC, nSetItemIndexD, nSetItemValueD, nCasterLevel, nCasterAttackPower)
		VALUES (@nCharNo, @nStateNo, @nStateLevel, @nStateTime, @nStateCount, @nStatePowerDemage, @nStatePowerSP, @nStatePowerKeepTime, @nStatePowerCoolTime, @nSetItemIndexA, @nSetItemValueA, @nSetItemIndexB, @nSetItemValueB, @nSetItemIndexC, @nSetItemValueC, @nSetItemIndexD, @nSetItemValueD, @nCasterLevel, @nCasterAttackPower)

		SET @nRet = @@ERROR
	END
END

0뀁퍅ā
GO

-- ========================================================================
-- Procedure: p_Char_SetStateWithSetItemAndRestTime
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_SetStateWithSetItemAndRestTime]
	@nCharNo int,
	@nStateNo smallint,
	@nStateLevel smallint,
	@nStateTime int,
	@nRestTime int,
	@nStateCount smallint,
	@nStatePowerDemage tinyint,
	@nStatePowerSP tinyint,
	@nStatePowerKeepTime tinyint,
	@nStatePowerCoolTime tinyint,
	@nSetItemIndexA smallint,
	@nSetItemValueA smallint,
	@nSetItemIndexB smallint,
	@nSetItemValueB smallint,
	@nSetItemIndexC smallint,
	@nSetItemValueC smallint,
	@nSetItemIndexd smallint,
	@nSetItemValueD smallint,
	@nCasterLevel smallint,
	@nCasterAttackPower int,	
	@nRet int OUTPUT
AS
BEGIN
	SET NOCOUNT ON

	IF EXISTS(SELECT nCharNo FROM tCharacterState WHERE nCharNo = @nCharNo AND nStateNo = @nStateNo )
	BEGIN
		UPDATE tCharacterState SET 
			nStateLevel			= @nStateLevel,
			nStateTime			= @nStateTime,
			nRestTime			= @nRestTime,
			nStateCount			= @nStateCount,	
			nStatePowerDemage	= @nStatePowerDemage,
			nStatePowerSP		= @nStatePowerSP,
			nStatePowerKeepTime	= @nStatePowerKeepTime,
			nStatePowerCoolTime	= @nStatePowerCoolTime,
			nSetItemIndexA		= @nSetItemIndexA,
			nSetItemValueA		= @nSetItemValueA,
			nSetItemIndexB		= @nSetItemIndexB,
			nSetItemValueB		= @nSetItemValueB,
			nSetItemIndexC		= @nSetItemIndexC,
			nSetItemValueC		= @nSetItemValueC,
			nSetItemIndexD		= @nSetItemIndexD,
			nSetItemValueD		= @nSetItemValueD,
		    nCasterLevel        = @nCasterLevel,
		    nCasterAttackPower  = @nCasterAttackPower 			
		WHERE nCharNo = @nCharNo AND nStateNo = @nStateNo
		
		SET @nRet = @@ERROR
	END
	ELSE BEGIN
		INSERT tCharacterState (  nCharNo,  nStateNo,  nStateLevel,  nStateTime, nRestTime,  nStateCount,  nStatePowerDemage,  nStatePowerSP,  nStatePowerKeepTime,  nStatePowerCoolTime, nSetItemIndexA, nSetItemValueA, nSetItemIndexB, nSetItemValueB, nSetItemIndexC, nSetItemValueC, nSetItemIndexD, nSetItemValueD, nCasterLevel, nCasterAttackPower )
		VALUES                 ( @nCharNo, @nStateNo, @nStateLevel, @nStateTime, @nRestTime, @nStateCount, @nStatePowerDemage, @nStatePowerSP, @nStatePowerKeepTime, @nStatePowerCoolTime, @nSetItemIndexA, @nSetItemValueA, @nSetItemIndexB, @nSetItemValueB, @nSetItemIndexC, @nSetItemValueC, @nSetItemIndexD, @nSetItemValueD, @nCasterLevel, @nCasterAttackPower )

		SET @nRet = @@ERROR
	END
END

0ጁﱴĂ
GO

-- ========================================================================
-- Procedure: p_Char_SetTitle
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_SetTitle]
/*
Character Set shape
2004.6 By CJC
 input:
   All args
 output: 
   0 = error
   1 = update/insert ok
*/
@nCharNo		int,
@nType			tinyint, 
@nElementNo		tinyint, 
@nElementValue	bigint,
@nRet			tinyint OUTPUT

AS

SET NOCOUNT ON

-- SET ERROR FIRST
SET @nRet = 0

IF EXISTS( SELECT nCharNo FROM tCharacterTitle WHERE nCharNo = @nCharNo AND nType = @nType AND nElementNo = @nElementNo )
BEGIN
	UPDATE tCharacterTitle SET nElementValue = @nElementValue
	WHERE nCharNo = @nCharNo AND nType = @nType AND nElementNo = @nElementNo

	IF @@ERROR = 0 AND @@ROWCOUNT = 1
	BEGIN
		SET @nRet = 1
		RETURN
	END
END
ELSE
BEGIN
	INSERT tCharacterTitle (  nCharNo,  nType,  nElementNo,  nElementValue )
	VALUES                 ( @nCharNo, @nType, @nElementNo, @nElementValue )

	IF @@ERROR = 0
	BEGIN
		SET @nRet = 1
		RETURN
	END
END

-- end

0欁⡋
GO

-- ========================================================================
-- Procedure: p_Char_StatPoint_Dec
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_StatPoint_Dec]
/*
Character decrement statpoint
2005.11 By CJC

 input:
   See SQL
   
 output: 
   0 = error
   1 = update ok
   
*/

@nCharNo int,
@nStatPointType tinyint,

-- Output var
@nRet int OUTPUT

AS

SET NOCOUNT ON

-- CHECK nRedistributePoint

/* DO NOT USE nRedistributePoint

DECLARE @nRedistributePoint tinyint
SET @nRedistributePoint = 0
SELECT @nRedistributePoint = nRedistributePoint FROM tCharacter WHERE nCharNo = @nCharNo
IF @nRedistributePoint < 1
BEGIN
	SET @nRet = 0
	RETURN
END
*/
-- CHECK Type and Decrement stat-point
DECLARE @nStatPoint tinyint

IF @nStatPointType = 0
BEGIN
	SELECT @nStatPoint = nStrength FROM tCharacter WHERE nCharNo = @nCharNo
	IF @nStatPoint < 1
	BEGIN
		SET @nRet = 0
		RETURN
	END
	
	UPDATE tCharacter	SET nStrength	= nStrength - 1					WHERE nCharNo = @nCharNo AND nStrength > 0
--	UPDATE tCharacter	SET nRedistributePoint = nRedistributePoint - 1 WHERE nCharNo = @nCharNo AND nRedistributePoint > 0
END

IF @nStatPointType = 1
BEGIN
	SELECT @nStatPoint = nConstitute FROM tCharacter WHERE nCharNo = @nCharNo
	IF @nStatPoint < 1
	BEGIN
		SET @nRet = 0
		RETURN
	END
	
	UPDATE tCharacter	SET nConstitute = nConstitute - 1				WHERE nCharNo = @nCharNo AND nConstitute > 0
--	UPDATE tCharacter	SET nRedistributePoint = nRedistributePoint - 1	WHERE nCharNo = @nCharNo AND nRedistributePoint > 0
END
	
IF @nStatPointType = 2
BEGIN
	SELECT @nStatPoint = nDexterity FROM tCharacter WHERE nCharNo = @nCharNo
	IF @nStatPoint < 1
	BEGIN
		SET @nRet = 0
		RETURN
	END
	
	UPDATE tCharacter	SET nDexterity = nDexterity - 1					WHERE nCharNo = @nCharNo AND nDexterity > 0
--	UPDATE tCharacter	SET nRedistributePoint = nRedistributePoint - 1	WHERE nCharNo = @nCharNo AND nRedistributePoint > 0
END

IF @nStatPointType = 3
BEGIN
	SELECT @nStatPoint = nIntelligence FROM tCharacter WHERE nCharNo = @nCharNo
	IF @nStatPoint < 1
	BEGIN
		SET @nRet = 0
		RETURN
	END
	
	UPDATE tCharacter	SET nIntelligence = nIntelligence - 1			WHERE nCharNo = @nCharNo AND nIntelligence > 0
--	UPDATE tCharacter	SET nRedistributePoint = nRedistributePoint - 1 WHERE nCharNo = @nCharNo AND nRedistributePoint > 0
END

IF @nStatPointType = 4 
BEGIN
	SELECT @nStatPoint = nMentalPower FROM tCharacter WHERE nCharNo = @nCharNo
	IF @nStatPoint < 1
	BEGIN
		SET @nRet = 0
		RETURN
	END
	
	UPDATE tCharacter	SET nMentalPower = nMentalPower - 1				WHERE nCharNo = @nCharNo AND nMentalPower > 0
--	UPDATE tCharacter	SET nRedistributePoint = nRedistributePoint - 1 WHERE nCharNo = @nCharNo AND nRedistributePoint > 0
END

-- Check Error
IF @@ERROR <> 0 OR @nStatPointType > 4
BEGIN
	SET @nRet = 0
	RETURN
END
SET @nRet = 1

-- end of file

0蔁Ą
GO

-- ========================================================================
-- Procedure: p_Char_StatPoint_Inc
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_StatPoint_Inc]
/*
Character increment statpoint
2005.11 By CJC

 input:
   See SQL
   
 output: 
   0 = error
   1 = update ok
   
*/

@nCharNo int,
@nStatPointType tinyint,

-- Output var
@nRet int OUTPUT

AS

SET NOCOUNT ON

IF @nStatPointType = 0 UPDATE tCharacter SET nStrength		= nStrength + 1		WHERE nCharNo = @nCharNo
IF @nStatPointType = 1 UPDATE tCharacter SET nConstitute	= nConstitute + 1	WHERE nCharNo = @nCharNo
IF @nStatPointType = 2 UPDATE tCharacter SET nDexterity		= nDexterity + 1	WHERE nCharNo = @nCharNo
IF @nStatPointType = 3 UPDATE tCharacter SET nIntelligence	= nIntelligence + 1 WHERE nCharNo = @nCharNo
IF @nStatPointType = 4 UPDATE tCharacter SET nMentalPower	= nMentalPower + 1	WHERE nCharNo = @nCharNo

-- check error
IF @@ERROR <> 0 OR @nStatPointType > 4
BEGIN
	SET @nRet = 0
	RETURN
END

SET @nRet = 1

-- end of file

0鐁ꏖą
GO

-- ========================================================================
-- Procedure: p_Char_SubExp
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_SubExp]
@nCharNo int,
@nExp bigint

AS
SET NOCOUNT ON

UPDATE tCharacter 
	SET 
		nExp=nExp-@nExp
	WHERE nCharNo=@nCharNo

return @@error

0촄Ć
GO

-- ========================================================================
-- Procedure: p_Char_WeddingStatus_Get
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_WeddingStatus_Get]
	@CharNo INT
AS
BEGIN
	SET NOCOUNT ON

	IF NOT EXISTS(SELECT * FROM tMarriage WHERE nCharNo = @CharNo)
	BEGIN
		INSERT tMarriage(nCharNo, nPartnerCharNo, dWhen, dSubWhen) VALUES(@CharNo, 0, GETDATE(), GETDATE())
	END

	SELECT nPartnerCharNo, nStatus, nSubStatus, dWhen, dSubWhen, nMarriageCount, nDivorceCancel FROM tMarriage WHERE nCharNo = @CharNo

	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: p_Char_WeddingStatus_Set
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_WeddingStatus_Set]
	@CharNo INT,
	@PartnerCharNo INT,
	@Status SMALLINT,
	@SubStatus SMALLINT,
	@When DATETIME,
	@SubWhen DATETIME,
	@MarriageCount TINYINT,
	@DivorceCancel TINYINT
AS

BEGIN

	SET NOCOUNT ON

	IF EXISTS(SELECT nCharNo FROM tMarriage WHERE nCharNo = @CharNo)
	BEGIN
		UPDATE tMarriage
		SET
			nPartnerCharNo = @PartnerCharNo,
			nStatus = @Status,
			nSubStatus = @SubStatus,
			dWhen = @When,
			dSubWhen = @SubWhen,
			nMarriageCount = @MarriageCount,
			nDivorceCancel = @DivorceCancel
		WHERE
			nCharNo = @CharNo
	END
	ELSE
	BEGIN
		INSERT tMarriage(nCharNo, nPartnerCharNo, nStatus, nSubStatus, dWhen, dSubWhen)
		VALUES(@CharNo, @PartnerCharNo, @Status, @SubStatus, @When, @SubWhen)
		
	END

	SET NOCOUNT OFF

END

0
GO

-- ========================================================================
-- Procedure: p_Char_Wedding_Cancel_Propose
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_Wedding_Cancel_Propose]
	@CharANo INT,
	@CharBNo INT,
	@nRet INT OUTPUT
AS

BEGIN

	SET NOCOUNT ON

	SELECT @nRet = 0

-- to Character A

	BEGIN TRAN

	IF EXISTS(SELECT nCharNo FROM tMarriage WHERE nCharNo = @CharANo AND nPartnerCharNo = @CharBNo)
	BEGIN
		UPDATE tMarriage
		SET
			nStatus = 0,
			nSubStatus = 0
		WHERE
			nCharNo = @CharANo

		 IF @@ERROR <> 0 GOTO ERROR_
	END
	ELSE GOTO ERROR_

-- to Character B

	IF EXISTS(SELECT nCharNo FROM tMarriage WHERE nCharNo = @CharBNo AND nPartnerCharNo = @CharANo)
	BEGIN
		UPDATE tMarriage
		SET
			nStatus = 0,
			nSubStatus = 0
		WHERE
			nCharNo = @CharBNo

		 IF @@ERROR <> 0 GOTO ERROR_
	END
	ELSE GOTO ERROR_

	COMMIT TRAN
	GOTO END_

ERROR_:
	ROLLBACK TRAN
	SELECT @nRet = 1

END_:
	SET NOCOUNT OFF

END

0
GO

-- ========================================================================
-- Procedure: p_Char_Wedding_Divorce_Cancel
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_Wedding_Divorce_Cancel]
	@CharANo INT,
	@CharBNo INT,
	@IncreseCount TINYINT,
	@SubStatusA TINYINT OUTPUT,
	@SubStatusB TINYINT OUTPUT,
	@nRet INT OUTPUT
AS

BEGIN

	SET NOCOUNT ON

	SELECT @nRet = 0

-- to Character A : applicant

	BEGIN TRAN

	IF EXISTS(SELECT nCharNo FROM tMarriage WHERE nCharNo = @CharANo AND nPartnerCharNo = @CharBNo)
	BEGIN
		SELECT @SubStatusA = nSubStatus FROM tMarriage WHERE nCharNo = @CharANo
		UPDATE tMarriage
		SET
			nSubStatus =
			CASE nSubStatus
				WHEN 1 THEN 1 + 4
				WHEN 3 THEN 3 + 4
				ELSE nSubStatus
			END,
			nDivorceCancel = nDivorceCancel + @IncreseCount,
			dSubWhen = dWhen
		WHERE
			nCharNo = @CharANo

		 IF @@ERROR <> 0 GOTO ERROR_
	END
	ELSE GOTO ERROR_

-- to Character B

	IF EXISTS(SELECT nCharNo FROM tMarriage WHERE nCharNo = @CharBNo AND nPartnerCharNo = @CharANo)
	BEGIN
		SELECT @SubStatusB = nSubStatus FROM tMarriage WHERE nCharNo = @CharBNo
		UPDATE tMarriage
		SET
			nSubStatus = 0,
			dSubWhen = dWhen
		WHERE
			nCharNo = @CharBNo

		 IF @@ERROR <> 0 GOTO ERROR_
	END
	ELSE GOTO ERROR_

	COMMIT TRAN
	GOTO END_

ERROR_:
	ROLLBACK TRAN
	SELECT @nRet = 1

END_:
	SET NOCOUNT OFF

END

0
GO

-- ========================================================================
-- Procedure: p_Char_Wedding_Divorce_Do
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_Wedding_Divorce_Do]
	@CharANo INT,
	@CharBNo INT,
	@nRet INT OUTPUT
AS

BEGIN

	SET NOCOUNT ON

	SELECT @nRet = 0

-- to Character A : applicant

	BEGIN TRAN

	IF EXISTS(SELECT nCharNo FROM tMarriage WHERE nCharNo = @CharANo AND nPartnerCharNo = @CharBNo)
	BEGIN
		UPDATE tMarriage
		SET
			nStatus = 3,
			nSubStatus = 1,
			dSubWhen = DATEADD(day, 7, GETDATE())
		WHERE
			nCharNo = @CharANo

		 IF @@ERROR <> 0 GOTO ERROR_
	END
	ELSE GOTO ERROR_

-- to Character B

	IF EXISTS(SELECT nCharNo FROM tMarriage WHERE nCharNo = @CharBNo AND nPartnerCharNo = @CharANo)
	BEGIN
		UPDATE tMarriage
		SET
			nStatus = 3,
			nSubStatus = 1,
			dSubWhen = DATEADD(day, 7, GETDATE())
		WHERE
			nCharNo = @CharBNo

		 IF @@ERROR <> 0 GOTO ERROR_
	END
	ELSE GOTO ERROR_

	COMMIT TRAN
	GOTO END_

ERROR_:
	ROLLBACK TRAN
	SELECT @nRet = 1

END_:
	SET NOCOUNT OFF

END

0
GO

-- ========================================================================
-- Procedure: p_Char_Wedding_Divorce_Do_by_Item
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_Wedding_Divorce_Do_by_Item]
	@CharANo INT,
	@CharBNo INT,
	@nRet INT OUTPUT
AS
BEGIN
	SET NOCOUNT ON
	SELECT @nRet = 0
-- to Character A : applicant
	BEGIN TRAN
	IF EXISTS(SELECT nCharNo FROM tMarriage WHERE nCharNo = @CharANo AND nPartnerCharNo = @CharBNo)
	BEGIN
		UPDATE tMarriage
		SET
			nStatus = 3,
			nSubStatus = 1,
			-- 즉시 결혼도 가능하달라고 해서 수정함
			dSubWhen = DATEADD(second, 5, GETDATE())
		WHERE
			nCharNo = @CharANo
		 IF @@ERROR <> 0 GOTO ERROR_
	END
	ELSE GOTO ERROR_
-- to Character B
	IF EXISTS(SELECT nCharNo FROM tMarriage WHERE nCharNo = @CharBNo AND nPartnerCharNo = @CharANo)
	BEGIN
		UPDATE tMarriage
		SET
			nStatus = 3,
			nSubStatus = 2,
			-- 즉시 결혼도 가능하달라고 해서 수정함
			dSubWhen = DATEADD(second, 5, GETDATE())
		WHERE
			nCharNo = @CharBNo
		 IF @@ERROR <> 0 GOTO ERROR_
	END
	ELSE GOTO ERROR_
	COMMIT TRAN
	GOTO END_
ERROR_:
	ROLLBACK TRAN
	SELECT @nRet = 1
END_:
	SET NOCOUNT OFF
END

0⌁僔Č
GO

-- ========================================================================
-- Procedure: p_Char_Wedding_Divorce_by_Agree
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_Wedding_Divorce_by_Agree]
	@CharANo INT,
	@CharBNo INT,
	@nRet INT OUTPUT
AS

BEGIN

	SET NOCOUNT ON

	SELECT @nRet = 0

-- to Character A : applicant

	BEGIN TRAN

	IF EXISTS(SELECT nCharNo FROM tMarriage WHERE nCharNo = @CharANo AND nPartnerCharNo = @CharBNo)
	BEGIN
		UPDATE tMarriage
		SET
			nSubStatus = 1,
			dSubWhen = DATEADD(minute, 10, GETDATE())
		WHERE
			nCharNo = @CharANo

		 IF @@ERROR <> 0 GOTO ERROR_
	END
	ELSE GOTO ERROR_

-- to Character B

	IF EXISTS(SELECT nCharNo FROM tMarriage WHERE nCharNo = @CharBNo AND nPartnerCharNo = @CharANo)
	BEGIN
		UPDATE tMarriage
		SET
			nSubStatus = 2 + 4,
			dSubWhen = DATEADD(minute, 10, GETDATE())
		WHERE
			nCharNo = @CharBNo

		 IF @@ERROR <> 0 GOTO ERROR_
	END
	ELSE GOTO ERROR_

	COMMIT TRAN
	GOTO END_

ERROR_:
	ROLLBACK TRAN
	SELECT @nRet = 1

END_:
	SET NOCOUNT OFF

END

0
GO

-- ========================================================================
-- Procedure: p_Char_Wedding_Divorce_by_Force
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_Wedding_Divorce_by_Force]
	@CharANo INT,
	@CharBNo INT,
	@nRet INT OUTPUT
AS

BEGIN

	SET NOCOUNT ON

	SELECT @nRet = 0

-- to Character A : applicant

	BEGIN TRAN

	IF EXISTS(SELECT nCharNo FROM tMarriage WHERE nCharNo = @CharANo AND nPartnerCharNo = @CharBNo)
	BEGIN
		UPDATE tMarriage
		SET
			nSubStatus = 3,
			dSubWhen = DATEADD(Day, 3, GETDATE())
--			dSubWhen = DATEADD(minute, 10, GETDATE())
		WHERE
			nCharNo = @CharANo

		 IF @@ERROR <> 0 GOTO ERROR_
	END
	ELSE GOTO ERROR_

-- to Character B

	IF EXISTS(SELECT nCharNo FROM tMarriage WHERE nCharNo = @CharBNo AND nPartnerCharNo = @CharANo)
	BEGIN
		UPDATE tMarriage
		SET
			nSubStatus = 4,
			dSubWhen = DATEADD(Day, 3, GETDATE())
--			dSubWhen = DATEADD(minute, 10, GETDATE())
		WHERE
			nCharNo = @CharBNo

		 IF @@ERROR <> 0 GOTO ERROR_
	END
	ELSE GOTO ERROR_

	COMMIT TRAN
	GOTO END_

ERROR_:
	ROLLBACK TRAN
	SELECT @nRet = 1

END_:
	SET NOCOUNT OFF

END

0
GO

-- ========================================================================
-- Procedure: p_Char_Wedding_Do
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_Wedding_Do] 
	@CharANo INT,
	@CharBNo INT,
	@nRet INT OUTPUT
AS

BEGIN

	SET NOCOUNT ON

	SELECT @nRet = 0

-- to Character A

	BEGIN TRAN

	IF EXISTS(SELECT nCharNo FROM tMarriage WHERE nCharNo = @CharANo AND nPartnerCharNo = @CharBNo)
	BEGIN
		UPDATE tMarriage
		SET
			nStatus = 2,
			nSubStatus = 0,
			dWhen = GETDATE(),
			dSubWhen = DATEADD(day, 7, GETDATE()),
--			dSubWhen = DATEADD(minute, 3, GETDATE()),
			nMarriageCount = nMarriageCount + 1,
			nDivorceCancel = 0
		WHERE
			nCharNo = @CharANo

		 IF @@ERROR <> 0 GOTO ERROR_
	END
	ELSE GOTO ERROR_

-- to Character B

	IF EXISTS(SELECT nCharNo FROM tMarriage WHERE nCharNo = @CharBNo AND nPartnerCharNo = @CharANo)
	BEGIN
		UPDATE tMarriage
		SET
			nStatus = 2,
			nSubStatus = 0,
			dWhen = GETDATE(),
			dSubWhen = DATEADD(day, 7, GETDATE()),
--			dSubWhen = DATEADD(minute, 3, GETDATE()),
			nMarriageCount = nMarriageCount + 1,
			nDivorceCancel = 0
		WHERE
			nCharNo = @CharBNo

		 IF @@ERROR <> 0 GOTO ERROR_
	END
	ELSE GOTO ERROR_

	COMMIT TRAN
	GOTO END_

ERROR_:
	ROLLBACK TRAN
	SELECT @nRet = 1

END_:
	SET NOCOUNT OFF

END

0
GO

-- ========================================================================
-- Procedure: p_Char_Wedding_Escape_Divorce_State
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_Wedding_Escape_Divorce_State]
	@CharNo INT,
	@nRet INT OUTPUT
AS

BEGIN

	SET NOCOUNT ON

	IF EXISTS(SELECT nCharNo FROM tMarriage WHERE nCharNo = @CharNo)
	BEGIN
		SELECT @nRet = 0
		UPDATE tMarriage
		SET
			nStatus = 0,
			nSubStatus = 0,
			nDivorceCancel = 0
		WHERE
			nCharNo = @CharNo
	END
	ELSE
	BEGIN
		SELECT @nRet = 1
	END

	SET NOCOUNT OFF

END

0
GO

-- ========================================================================
-- Procedure: p_Char_Wedding_Get_Partner_Info
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_Wedding_Get_Partner_Info]
	@CharNo INT
AS

BEGIN

	SET NOCOUNT ON

	SELECT nLevel, dLastLoginDate FROM tCharacter WHERE nCharNo = @CharNo

	SET NOCOUNT OFF

END

0
GO

-- ========================================================================
-- Procedure: p_Char_Wedding_Propose
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_Wedding_Propose]
	@CharANo INT,
	@CharBNo INT,
	@nRet INT OUTPUT
AS

BEGIN

	SET NOCOUNT ON

	SELECT @nRet = 0

-- to Character A

	BEGIN TRAN

	IF EXISTS(SELECT nCharNo FROM tMarriage WHERE nCharNo = @CharANo)
	BEGIN
		UPDATE tMarriage
		SET
			nPartnerCharNo = @CharBNo,
			nStatus = 1,
			nSubStatus = 0,
			dWhen = GETDATE(),
			dSubWhen = DATEADD(day, 3, GETDATE())

		WHERE
			nCharNo = @CharANo

		 IF @@ERROR <> 0 GOTO ERROR_
	END
	ELSE
	BEGIN
		INSERT tMarriage(nCharNo, nPartnerCharNo, nStatus, nSubStatus, dWhen, dSubWhen)
		VALUES(@CharANo, @CharBNo, 1,  0, GETDATE(), DATEADD(day, 3, GETDATE()))		

		 IF @@ERROR <> 0 GOTO ERROR_
	END

-- to Character B

	IF EXISTS(SELECT nCharNo FROM tMarriage WHERE nCharNo = @CharBNo)
	BEGIN
		UPDATE tMarriage
		SET
			nPartnerCharNo = @CharANo,
			nStatus = 1,
			nSubStatus = 0,
			dWhen = GETDATE(),
			dSubWhen = DATEADD(day, 3, GETDATE())
		WHERE
			nCharNo = @CharBNo

		 IF @@ERROR <> 0 GOTO ERROR_
	END
	ELSE
	BEGIN
		INSERT tMarriage(nCharNo, nPartnerCharNo, nStatus, nSubStatus, dWhen, dSubWhen)
		VALUES(@CharBNo, @CharANo, 1,  0, GETDATE(), DATEADD(day, 3, GETDATE()))		

		 IF @@ERROR <> 0 GOTO ERROR_
	END

	COMMIT TRAN
	GOTO END_

ERROR_:
	ROLLBACK TRAN
	SELECT @nRet = 1

END_:
	SET NOCOUNT OFF

END

0
GO

-- ========================================================================
-- Procedure: p_Char_Wipe
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Char_Wipe]
/*

Character Wipe(Erase)
2006.5 By CJC

 input:
	nCharNo

 output:
	@nRet	0 = OK, 
	        ? = Error code
*/

@nCharNo int,

@nRet int OUTPUT

AS

SET NOCOUNT ON

DELETE FROM tFriend WHERE nCharNo = @nCharNo			-- 관계에 의해 삭제됨
DELETE FROM tFriend WHERE nFriendCharNo = @nCharNo	-- 친구 테이블 순환 참고 없앰용

DELETE FROM tCharacter WHERE nCharNo = @nCharNo

SET @nRet = @@ERROR

-- end

0눁﷑Ē
GO

-- ========================================================================
-- Procedure: p_Event_Friend_Reward
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Event_Friend_Reward]
@nFriendCharNo	int,
@nItemKey	bigint,
@err 		smallint OUTPUT

AS
SET NOCOUNT ON

DECLARE
@nFriendUserNo		int,
@nLevel 			smallint,
@nRewardItemID		int,
@nLevelReward_1_23	int,
@nLevelReward_24_47	int,
@nLevelReward_48_100	int,
@nStorage_Min		smallint,
@nStorage_Max		smallint,
@nStorageType		tinyint,
@nItemID		int,
@nOptionType		smallint,
@nOptionData		bigint,
@nOldItemKey		bigint



-- Hard Fixed...

SET @nStorage_Max 		= 143 	--개정창고 4 맨끝 위치.

SET @nStorageType 		= 6 	--개정창고

SET @nLevelReward_1_23 		= 8004 	-- 축복의 릭스(강도1)
SET @nLevelReward_24_47 	= 8013 	-- 축복의 릭스(강도2)
SET @nLevelReward_48_100 	= 8022 	-- 축복의 릭스(강도3)

SET @nOptionType 		= 1	-- Max = 255


SELECT @nLevel = nLevel, @nFriendUserNo=nUserNo  FROM tCharacter WITH(READUNCOMMITTED) 
	WHERE nCharNo = @nFriendCharNo
IF @@ROWCOUNT = 0
BEGIN
	SET @err = 1
	PRINT('Not Exitist Friend')
	RETURN @@ERROR
END
PRINT('@nLeve=')
PRINT(@nLevel)

IF @nLevel > 0 AND @nLevel < 24
	SET @nRewardItemID = @nLevelReward_1_23
ELSE IF @nLevel > 23 AND @nLevel < 48
	SET @nRewardItemID = @nLevelReward_24_47
ELSE
	SET @nRewardItemID = @nLevelReward_48_100

PRINT('@nRewardItemID=')
PRINT(@nRewardItemID)

SELECT @nStorage_Min = ISNULL(MIN( nStorage ), 0) FROM tItem WHERE nOwner=@nFriendUserNo AND nStorageType =@nStorageType 
	AND ( nItemID= @nLevelReward_1_23 OR nItemID= @nLevelReward_24_47 OR nItemID= @nLevelReward_48_100 )

PRINT('@nStorage_Min=')
PRINT(@nStorage_Min)

-- 개정창고 2 3 4 번 다 썼다.(@nStorage_Min=36)
IF @nStorage_Min > 0 AND @nStorage_Min < 37 
	BEGIN
		PRINT('IF @nStorage_Min > 0 AND @nStorage_Min < 36 ')
		GOTO CRITICAL
	END

IF @nStorage_Min = 0
	BEGIN
		SET	@nStorage_Min = @nStorage_Max
		SET	@nOptionData = 5


	END
ELSE
	BEGIN
		SELECT @nItemID = I.nItemID, @nOptionData=O.nOptionData, @nOldItemKey=I.nItemKey FROM tItem I
			INNER JOIN tItemOptions O ON I.nItemKey = O.nItemKey AND O.nOptionType = @nOptionType
			WHERE I.nOwner=@nFriendUserNo AND I.nStorageType =@nStorageType 
			AND I.nStorage = @nStorage_Min 

		PRINT('@@ROWCOUNT=')
		PRINT(@@ROWCOUNT)

		PRINT('@nItemID=')
		PRINT(@nItemID)

		PRINT('@nOptionData=')
		PRINT(@nOptionData)

		IF @nRewardItemID = @nItemID
			BEGIN
				-- 기존값에 추가한다.
				IF @nOptionData > 0 AND @nOptionData < 6 
				BEGIN
					SET	@nOptionData = @nOptionData+5
					--SET	@nStorage_Min = @nStorage_Min
				END
				ELSE
				BEGIN
					SET 	@nOptionData = 5
					SET	@nStorage_Min = @nStorage_Min-1
					PRINT('SET	@nStorage_Min = @nStorage_Min-1')
					PRINT(@nStorage_Min)

				END

			END
		ELSE
			BEGIN
				SET	@nStorage_Min = @nStorage_Min-1
				SET	@nOptionData = 5
			END
		
	END

-- 넣으려는 자리에 다른 아이템이 있는지 체크.(이벤트 전에 들어가 있던 아이템)
IF @nOptionData = 5
BEGIN
	IF EXISTS (SELECT nItemID  FROM tItem WHERE nOwner=@nFriendUserNo AND nStorageType =@nStorageType 
		AND nStorage=@nStorage_Min )
	BEGIN
		GOTO CRITICAL
	END

END

-- 아템 넣기.
BEGIN TRANSACTION

-- 새롭게 아템을 넣는다.
IF @nOptionData = 5
	BEGIN
		PRINT('INSERT tItem nItemKey')
		PRINT(@nItemKey)

		INSERT INTO tItem(nItemKey, nStorageType, nOwner, nStorage, nItemID, nFlags)
			VALUES(@nItemKey, @nStorageType, @nFriendUserNo, @nStorage_Min, @nRewardItemID, 0 )
		IF @@ERROR <> 0
		BEGIN
			GOTO ERR
		END
		PRINT('INSERT tItemOptions nItemKey')
		PRINT(@nItemKey)
		INSERT INTO tItemOptions(nItemKey, nOptionType, nOptionData)
			VALUES(@nItemKey, @nOptionType, @nOptionData )
		IF @@ERROR <> 0
		BEGIN
			GOTO ERR
		END
	END
-- 기존에 개수를 늘린
GO

-- ========================================================================
-- Procedure: p_Event_WorldCup_Init
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Event_WorldCup_Init]
@nCharNo int,
@bNewJoiner tinyint OUTPUT

AS
SET NOCOUNT ON

DECLARE
@dRegisterDate	datetime,
@dEventStartDate	datetime,
@dEventEndDate	datetime


SET @dEventStartDate = '2006-05-30' --테스트
SET @dEventEndDate =  '2006-06-07' --6월 06일까지.

IF GETDATE() > @dEventEndDate
BEGIN
	SET	@bNewJoiner = 0
	RETURN @@error
END

SELECT @dRegisterDate=dCreateDate   
	FROM tCharacter
	WHERE nCharNo = @nCharNo --AND 

IF @@ROWCOUNT = 0
BEGIN
	SET	@bNewJoiner = 0
	RETURN @@error
END

IF @dRegisterDate > @dEventStartDate AND @dRegisterDate < @dEventEndDate
	BEGIN
		SET	@bNewJoiner = 1
	END
ELSE
	BEGIN
		SET	@bNewJoiner = 0
	END

RETURN @@error

0␁Ĕ
GO

-- ========================================================================
-- Procedure: p_Friend_Del
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Friend_Del]
/*
Friend Delete
2004.12 By CJC
 input:
	nCharNo, nFriendCharNo = To be delete info.
 output:
	@nRet	0 = OK, 
	        ? = Error code
*/

@nCharNo int,
@nFriendCharNo int,
@nRet int OUTPUT

AS

SET NOCOUNT ON

-- cross delete
DELETE FROM tFriend WHERE nCharNo = @nCharNo AND nFriendCharNo = @nFriendCharNo
SET @nRet = @@ERROR

DELETE FROM tFriend WHERE nCharNo = @nFriendCharNo AND nFriendCharNo = @nCharNo 
SET @nRet = @nRet + @@ERROR

-- end

0蜁཈Ė
GO

-- ========================================================================
-- Procedure: p_Friend_Del_All
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Friend_Del_All]
/*
Friend Delete
2007.12 By CJC
 input:
	nCharNo = To be delete info.
 output:
	@nRet	0 = OK, 
	        ? = Error code
*/

@nCharNo int,
@nRet int OUTPUT

AS

SET NOCOUNT ON

-- cross delete
DELETE FROM tFriend WHERE nCharNo = @nCharNo OR nFriendCharNo = @nCharNo
SET @nRet = @@ERROR

-- end

0
GO

-- ========================================================================
-- Procedure: p_Friend_Get
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Friend_Get]

/*
Char Get friend
2004.12 By CJC
 input:
	see SQL
 output:
	fetch to take data
*/

@nCharNo int,
@nFriendNo int

AS

SET NOCOUNT ON SELECT 
F.nFriendCharNo,
F.nFlags, 

C.sID,
S.nClass,
C.nLevel,
C.sLoginZone,
C.dLastLoginDate

FROM tFriend AS F
INNER JOIN tCharacterShape AS S ON S.nCharNo =  @nFriendNo
INNER JOIN tCharacter AS C ON C.nCharNo = @nFriendNo AND C.bDeleted = 0
WHERE F.nCharNo = @nCharNo AND F.nFriendCharNo = @nFriendNo

-- end

0更ė
GO

-- ========================================================================
-- Procedure: p_Friend_GetAll
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Friend_GetAll]

/*
Char Get friend all
2004.12 By CJC
 input:
	see SQL
 output:
	fetch to take data
*/

@nCharNo int

AS

SET NOCOUNT ON

SELECT F.nFriendCharNo, F.nFlags, C.sID
FROM tFriend AS F
INNER JOIN tCharacter AS C ON C.nCharNo =  F.nFriendCharNo AND C.bDeleted = 0
WHERE F.nCharNo = @nCharNo

-- end

0㈁Ę
GO

-- ========================================================================
-- Procedure: p_Friend_GetAll2
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Friend_GetAll2]

/*
Char Get friend all
2007.2.14 By CJC
 input:
	see SQL
 output:
	fetch to take data
*/

@nCharNo int

AS

SET NOCOUNT ON

SELECT 
F.nFriendCharNo,
F.nFlags, 

C.sID,
S.nClass,
C.nLevel,
C.sLoginZone,
C.dLastLoginDate

FROM tFriend AS F
INNER JOIN tCharacterShape AS S ON S.nCharNo = F.nFriendCharNo
INNER JOIN tCharacter AS C ON C.nCharNo = F.nFriendCharNo AND C.bDeleted = 0
WHERE F.nCharNo = @nCharNo

-- end

0켁슆ė
GO

-- ========================================================================
-- Procedure: p_Friend_GetNumOfFriend
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Friend_GetNumOfFriend]

/*
Char Get friend all
2004.12 By CJC
 input:
	see SQL
 output:
	fetch to take data
*/

@nCharNo int,
@nNumOfFriend int OUTPUT

AS

SET NOCOUNT ON

SET @nNumOfFriend = 0

SELECT @nNumOfFriend = COUNT(*) FROM tFriend AS F WITH(READUNCOMMITTED) 
INNER JOIN tCharacter AS C WITH(READUNCOMMITTED) ON C.nCharNo = F.nFriendCharNo AND C.bDeleted = 0
WHERE F.nCharNo = @nCharNo

-- end

0稁黳Ě
GO

-- ========================================================================
-- Procedure: p_Friend_Set
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Friend_Set]
/*
Char Set friend
2004.12 By CJC

 input:
	see SQL
	
 output:
	@nRet	0 = OK, 
			1 = friend and me are same char
			2 = i'm over max friend count
			3 = friend is over max friend count
			547 = DB error, nCharNo or nFriendCharNo is not existed Character(Foreign Key checked)
			2627 = DB error, already friend marked
	        ? = Other DB error code
	        
*/
@nMAXFRIEND smallint,
@nCharNo int,
@nFriendCharNo int,
@nFlags tinyint,
@nRet int OUTPUT

AS

SET NOCOUNT ON

-- check firend and me are same char.
IF @nCharNo = @nFriendCharNo
BEGIN
	SET @nRet = 1
	RETURN
END

-- check over max firend
DECLARE @nFriendCount int
EXEC p_Friend_GetNumOfFriend @nCharNo, @nFriendCount OUTPUT
IF @nFriendCount >= @nMAXFRIEND
BEGIN
	SET @nRet = 2
	RETURN
END

EXEC p_Friend_GetNumOfFriend @nFriendCharNo, @nFriendCount OUTPUT
IF @nFriendCount >= @nMAXFRIEND
BEGIN
	SET @nRet = 3
	RETURN
END

-- add friend
BEGIN TRAN

INSERT tFriend (  nCharNo,  nFriendCharNo,  nFlags )
VALUES         ( @nCharNo, @nFriendCharNo, @nFlags )
SET @nRet = @@ERROR
if( @nRet <> 0 )
BEGIN
	ROLLBACK TRAN
	RETURN
END

-- add friend cross
INSERT tFriend (  nCharNo,  nFriendCharNo,  nFlags )
VALUES         ( @nFriendCharNo, @nCharNo, @nFlags )
SET @nRet = @@ERROR
if( @nRet <> 0 )
BEGIN
	ROLLBACK TRAN
	RETURN
END

COMMIT TRAN
RETURN

-- end

0댁錗ě
GO

-- ========================================================================
-- Procedure: p_GuildTournament_Get
-- ========================================================================
CREATE PROCEDURE [dbo].[p_GuildTournament_Get]
/*
Get Tournament Guild Data
2007.5 By CJC
*/
@nGTNo int
AS
SET NOCOUNT ON
SELECT nType, nData FROM tGuildTournament WHERE @nGTNo = nGTNo
-- end

0
GO

-- ========================================================================
-- Procedure: p_GuildTournament_GetAllTournamentList
-- ========================================================================
CREATE PROCEDURE [dbo].[p_GuildTournament_GetAllTournamentList]
AS
BEGIN
	SELECT nGTNo
	FROM   tGuildTournament
	GROUP BY nGTNo
END

0
GO

-- ========================================================================
-- Procedure: p_GuildTournament_GetLastMatch
-- ========================================================================
CREATE PROCEDURE [dbo].[p_GuildTournament_GetLastMatch]
AS
BEGIN
	SELECT   nGTNo, nType, nData
	FROM      tGuildTournament
	WHERE   (nGTNo IN
                    (SELECT   MAX(nGTNo)
                     FROM      tGuildTournament))
             ORDER BY nType
END

0
GO

-- ========================================================================
-- Procedure: p_GuildTournament_Set
-- ========================================================================
CREATE PROCEDURE [dbo].[p_GuildTournament_Set]
/*
Set/Update Guild Tournament Data
2007.5 By CJC
*/
@nGTNo				int,
@nType				smallint,
@nData				int,
@nRet				int OUTPUT
AS
SET NOCOUNT ON
UPDATE tGuildTournament SET	nData = @nData WHERE nGTNo = @nGTNo AND nType = @nType
IF @@ERROR <> 0 OR @@ROWCOUNT = 0
BEGIN
	INSERT tGuildTournament (  nGTNo,  nType,  nData )
	VALUES					( @nGTNo, @nType, @nData )
END
SET @nRet = @@ERROR
-- end

0
GO

-- ========================================================================
-- Procedure: p_Guild_Delete
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Guild_Delete]

/*
Char Delete Guild
2006.2 By CJC

Input:
	See SQL
	
Output:
	nRet = 0,  Error(Not found guild)
	nRet = 1,  OK
*/

@nNo int,
@nRet tinyint output

AS

SET NOCOUNT ON

SET @nRet = 0

DELETE FROM tGuild WHERE nNo = @nNo
SET @nRet = @@ROWCOUNT

-- end

0ᘁ뱆Ĝ
GO

-- ========================================================================
-- Procedure: p_Guild_Get
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Guild_Get]

/*
Get Guild
2006.2 By CJC
*/

@nNo int

AS

SET NOCOUNT ON

SELECT 

	sName,
	sPassword,
	nMaxMembers,
	nMoney,
	nType,
	nGrade,
	nFame,
	nStoneLevel,
	nExp,
	dCreateDate,

	dWarRequestDate,
	dWarStartDate,
	dWarEndDate,
	nWarEnemy,
	nWarWinCount,
	nWarLoseCount,
	nWarDrawCount,
	
	nDismissStatus,
	dDismissDate,

	dNotifyDate,
	sNotifyCharID,
	sNotify,
	
	sIntro
	
FROM tGuild WHERE @nNo = nNo

-- end

0企끪ĝ
GO

-- ========================================================================
-- Procedure: p_Guild_GetAll
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Guild_GetAll]

/*
Char Get All Guild
2006.2 By CJC
*/

AS

SET NOCOUNT ON

SELECT

	nNo,
	sName,
	sPassword,
	nMaxMembers,
	nMoney,
	nType,
	nGrade,
	nFame,
	nStoneLevel,
	nExp,
	dCreateDate,

	dWarRequestDate,
	dWarStartDate,
	dWarEndDate,
	nWarEnemy,
	nWarWinCount,
	nWarLoseCount,
	nWarDrawCount,

	nDismissStatus,
	dDismissDate,

	dNotifyDate,
	sNotifyCharID,
	sNotify,
	
	sIntro
	
FROM tGuild

-- end

0蠁꒎Ğ
GO

-- ========================================================================
-- Procedure: p_Guild_GetCount
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Guild_GetCount]

/*
Char Get Guild Count
2006.2 By CJC
*/

@nCount int output

AS

SET NOCOUNT ON

SELECT @nCount = COUNT(*) FROM tGuild

-- end

0댁錗ě
GO

-- ========================================================================
-- Procedure: p_Guild_GetIntro
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Guild_GetIntro]

/*
Get Guild Intro
2006.2 By CJC
*/

@nNo int

AS

SET NOCOUNT ON

SELECT sIntro FROM tGuild WHERE @nNo = nNo

-- end

0度賖Ġ
GO

-- ========================================================================
-- Procedure: p_Guild_GetList
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Guild_GetList]

/*
Char Get All Guild lIst
2006.3 By blacksky
*/

AS

SET NOCOUNT ON

SELECT

	nNo,
	sName
	
FROM tGuild

-- end

0퀁埌Ġ
GO

-- ========================================================================
-- Procedure: p_Guild_GetListDetail
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Guild_GetListDetail]

/*
Char Get All Guild lIst detail
2007.11 By khris
*/

AS

SET NOCOUNT ON

SELECT

 nNo,
 sName,
 nMaxMembers,
 nType,
 nGrade,
 nFame,
 nExp
 
FROM tGuild

-- end

0
GO

-- ========================================================================
-- Procedure: p_Guild_GetMyNo
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Guild_GetMyNo]

/*
Get My Guild Number
2006.2 By CJC

Input:
	See SQL
	
Output:
	@nNo = 0, I'm not join guild
	@nNo = ?, My guild No.
*/

@nCharNo int,
@nNo     int output

AS

SET NOCOUNT ON

SET @nNo = 0
SELECT @nNo = nNo FROM tGuildMember WHERE nCharNo = @nCharNo

-- end

0笁㐹ģ
GO

-- ========================================================================
-- Procedure: p_Guild_GetName
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Guild_GetName]

/*
Get Guild Intro
2006.2 By Blacksky
*/

@nNo int

AS

SET NOCOUNT ON

SELECT sName FROM tGuild WHERE @nNo = nNo

-- end

0訁ģ
GO

-- ========================================================================
-- Procedure: p_Guild_GetNotify
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Guild_GetNotify]

/*
Get Guild Notify
2006.2 By CJC
*/

@nNo int

AS

SET NOCOUNT ON

SELECT dNotifyDate, sNotifyCharID, sNotify FROM tGuild WHERE @nNo = nNo

-- end

0ᲁĥ
GO

-- ========================================================================
-- Procedure: p_Guild_Make
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Guild_Make]

/*
Make Guild
2006.2 By CJC

Input:
	See SQL
	
Output:
	nNo = 0,   Error
	nNo = -1,  Error( Same guild name exist )
	nNo = -2,  Error( Not exists need item )
	nNo = -3,  Error( Other guild member )
	nNo = -4,  Error( Not enough maker money )
	nNo = -5,  Error( Make Guild error )
	nNo = -6,  Error( Maker to Join Guild master )
	nNo = 1~?, OK( Maked guild No )
*/

@sName nvarchar(18),
@sPassword nvarchar(9),
@nMaxMembers smallint,
@nType tinyint,

@nMakerCharNo int,
@nMakerPay bigint,
@nMakerNeedItemID int,

@nNo int output

AS

SET NOCOUNT ON

-- 에러코드 일단 세팅
SET @nNo = 0

-- 해당 길드가 있는지 확인
IF EXISTS (SELECT sName FROM tGuild(NOLOCK) WHERE sName = @sName)
BEGIN
	SET @nNo = -1
	RETURN
END

-- 필요아이템을 사용하는경우, 해당 아이템이 있는지 확인
IF @nMakerNeedItemID <> 65535
BEGIN
	IF NOT EXISTS (SELECT nItemID FROM tItem(NOLOCK) WHERE nOwner = @nMakerCharNo AND nItemID = @nMakerNeedItemID)
	BEGIN
		SET @nNo = -2
		RETURN
	END
END

-- 이미 다른 길드 가입했는지 확인
IF EXISTS (SELECT nCharNo FROM tGuildMember(NOLOCK) WHERE nCharNo = @nMakerCharNo)
BEGIN
	SET @nNo = -3
	RETURN
END

-- 길드 창설비용 적용
DECLARE @nRet tinyint

BEGIN TRAN
EXEC p_Char_MoneySub @nMakerCharNo, @nMakerPay, @nRet OUTPUT
IF @nRet = 0 -- 에러 확인
BEGIN
	ROLLBACK TRAN
	SET @nNo = -4
	RETURN
END

-- 길드 만들기
INSERT tGuild (  sName,  sPassword,  nMaxMembers,  nType )
VALUES        ( @sName, @sPassword, @nMaxMembers, @nType )
SET @nNo = @@IDENTITY

-- 길드 만들기 오류 확인
IF @@ERROR <> 0 
BEGIN
	ROLLBACK TRAN
	SET @nNo = -5
	RETURN
END

-- 만든 캐릭을 길마으로 추가
EXEC p_Guild_Member_Join @nNo, @nMakerCharNo, 0, @nRet OUTPUT
IF @nRet = 0
BEGIN 
	ROLLBACK TRAN
	SET @nNo = -6
	RETURN
END

-- 만들기 성공
COMMIT TRAN

RETURN

-- end

0☁ႦĦ
GO

-- ========================================================================
-- Procedure: p_Guild_Make2
-- ========================================================================
CREATE  PROCEDURE [dbo].[p_Guild_Make2]
/*
Make Guild
2006.12.15 By CJC
Input:
	See SQL
	
Output:
	nNo = 0,   Error, Remain char money
	nNo = -1,  Error( Same guild name exist )
	nNo = -2,  Error( Not exists need item )
	nNo = -3,  Error( Other guild member )
	nNo = -4,  Error( Not enough maker money )
	nNo = -5,  Error( Make Guild error )
	nNo = -6,  Error( Maker to Join Guild master )
	nNo = 1~?, OK( Maked guild No )
*/
@sName nvarchar(18),
@sPassword nvarchar(9),
@nMaxMembers smallint,
@nType tinyint,
@nMakerCharNo int,
@nMakerPay bigint,
@nMakerNeedItemID int,
@nNo int output,
@nRemainCharMoney bigint output
AS
SET NOCOUNT ON
-- 에러코드 일단 세팅
SET @nNo = 0
-- 해당 길드가 있는지 확인
IF EXISTS (SELECT sName FROM tGuild(NOLOCK) WHERE sName = @sName)
BEGIN
	SET @nNo = -1
	RETURN
END
-- 필요아이템을 사용하는경우, 해당 아이템이 있는지 확인
IF @nMakerNeedItemID <> 65535
BEGIN
	IF NOT EXISTS (SELECT nItemID FROM tItem(NOLOCK) WHERE nOwner = @nMakerCharNo AND nItemID = @nMakerNeedItemID)
	BEGIN
		SET @nNo = -2
		RETURN
	END
END
-- 이미 다른 길드 가입했는지 확인
IF EXISTS (SELECT nCharNo FROM tGuildMember(NOLOCK) WHERE nCharNo = @nMakerCharNo)
BEGIN
	SET @nNo = -3
	RETURN
END
-- 길드 창설비용 적용
DECLARE @nRet tinyint
BEGIN TRAN
EXEC p_Char_MoneySub @nMakerCharNo, @nMakerPay, @nRet OUTPUT
IF @nRet = 0 -- 에러 확인
BEGIN
	ROLLBACK TRAN
	SET @nNo = -4
	RETURN
END
-- 길드 창설비용을 쓰고 남은돈 얻기
EXEC p_Char_MoneyGet @nMakerCharNo, @nRemainCharMoney OUTPUT, @nRet OUTPUT
IF @nRet = 0 -- 에러 확인
BEGIN
	ROLLBACK TRAN
	SET @nNo = -4
	RETURN
END
-- 길드 만들기
INSERT tGuild (  sName,  sPassword,  nMaxMembers,  nType )
VALUES        ( @sName, @sPassword, @nMaxMembers, @nType )
SET @nNo = @@IDENTITY
-- 길드 만들기 오류 확인
IF @@ERROR <> 0 
BEGIN
	ROLLBACK TRAN
	SET @nNo = -5
	RETURN
END
-- 만든 캐릭을 길마으로 추가
EXEC p_Guild_Member_Join @nNo, @nMakerCharNo, 0, @nRet OUTPUT
IF @nRet = 0
BEGIN 
	ROLLBACK TRAN
	SET @nNo = -6
	RETURN
END
-- 길드아카데미 생성
EXEC usp_GuildAcademy_setAcademy @nNo, @nMakerCharNo, 50, @nRet OUTPUT
IF @nRet <> 0 -- 길드아카데미 생성에서는 정상처리시 0을 return value로 가짐
BEGIN 
	ROLLBACK TRAN
	SET @nNo = @nRet
	RETURN
END
-- 만들기 성공
COMMIT TRAN
RETURN
-- end

0褁㧔ħ
GO

-- ========================================================================
-- Procedure: p_Guild_Member_Get
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Guild_Member_Get]

/*
Get Guild Member
2006.2 By CJC

Input:
	See SQL
	
Output:
	fetch to take data
*/

@nNo int,
@nCharNo int

AS

SET NOCOUNT ON

SELECT G.dJoinDate, G.nGrade, G.nFlags, G.nExpRatio, G.sIntro, C.sID, S.nClass, C.nLevel, C.sLoginZone, C.dLastLoginDate
FROM tGuildMember    AS G INNER JOIN 
     tCharacter      AS C ON C.nCharNo =  G.nCharNo AND C.bDeleted = 0 INNER JOIN 
     tCharacterShape AS S ON S.nCharNo =  G.nCharNo
WHERE G.nCharNo = @nCharNo AND G.nNo = @nNo

-- end

0숁ⷸĨ
GO

-- ========================================================================
-- Procedure: p_Guild_Member_GetAll
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Guild_Member_GetAll]

/*
Get Guild Member All
2006.2 By CJC

Input:
	See SQL
	
Output:
	nRet = 0,  Error(Not found guild)
	nRet = 1,  OK
*/

@nNo int

AS

SET NOCOUNT ON

SELECT G.nCharNo, G.dJoinDate, G.nGrade, G.nFlags, G.nExpRatio, G.sIntro, C.sID, S.nClass, C.nLevel, C.sLoginZone, C.dLastLoginDate
FROM tGuildMember    AS G INNER JOIN 
     tCharacter      AS C ON C.nCharNo =  G.nCharNo AND C.bDeleted = 0 INNER JOIN 
     tCharacterShape AS S ON S.nCharNo =  G.nCharNo
WHERE G.nNo = @nNo

-- end

0弁ӊħ
GO

-- ========================================================================
-- Procedure: p_Guild_Member_GetCount
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Guild_Member_GetCount]

/*
Get Guild Member Count
2006.2 By CJC

Input:
	See SQL
	
Output:
	nCount = Number of member
*/

@nNo int,
@nCount int output
AS

SET NOCOUNT ON

SELECT @nCount = COUNT(*) FROM tGuildMember AS G WITH(READUNCOMMITTED) 
INNER JOIN tCharacter AS C WITH(READUNCOMMITTED) ON C.nCharNo = G.nCharNo AND C.bDeleted = 0
WHERE G.nNo = @nNo 

-- end

0갬ĩ
GO

-- ========================================================================
-- Procedure: p_Guild_Member_Join
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Guild_Member_Join]

/*
Get Guild Member Join
2006.2 By CJC

Input:
	See SQL
	
Output:
	nRet = 0, Error
	nRet = 1, OK
*/

@nNo int,
@nCharNo int,
@nGrade tinyint,

@nRet tinyint output
AS

SET NOCOUNT ON

SET @nRet = 0

-- 만들기
INSERT tGuildMember (  nNo,  nCharNo,  nGrade )
VALUES              ( @nNo, @nCharNo, @nGrade )

-- 오류 확인
IF @@ERROR <> 0 
BEGIN
	RETURN
END

SET @nRet = 1

-- end

RETURN

0䌁핛Ī
GO

-- ========================================================================
-- Procedure: p_Guild_Member_Leave
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Guild_Member_Leave]

/*
Get Guild Member Leave
2006.2 By CJC

Input:
	See SQL
	
Output:
	nRet = 0, Error
	nRet = 1, OK
*/

@nNo int,
@nCharNo int,

@nRet tinyint output
AS

SET NOCOUNT ON

SET @nRet = 0

-- 삭제
DELETE FROM tGuildMember WHERE nNo = @nNo AND nCharNo = @nCharNo
SET @nRet = @@ROWCOUNT

-- end

0刁鑵ī
GO

-- ========================================================================
-- Procedure: p_Guild_Member_SetExpRatio
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Guild_Member_SetExpRatio]

/*
Set Guild Member Exp Ratio
2006.2 By CJC

Input:
	See SQL
	
Output:
	nRet = 0,  Error
	nRet = 1,  OK
*/

@nNo int,
@nCharNo int,
@nExpRatio smallint,

@nRet tinyint output

AS

SET NOCOUNT ON
SET @nRet = 0

UPDATE tGuildMember SET nExpRatio = @nExpRatio WHERE nNo = @nNo AND nCharNo = @nCharNo
SET @nRet = @@ROWCOUNT

IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
BEGIN
	SET @nRet = 0
	RETURN
END

-- end

0謁袙Ĭ
GO

-- ========================================================================
-- Procedure: p_Guild_Member_SetFlags
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Guild_Member_SetFlags]

/*
Set Guild Member Flags
2006.2 By CJC

Input:
	See SQL
	
Output:
	nRet = 0,  Error
	nRet = 1,  OK
*/

@nNo int,
@nCharNo int,
@nFlags int,

@nRet tinyint output

AS

SET NOCOUNT ON
SET @nRet = 0

UPDATE tGuildMember SET nFlags = @nFlags WHERE nNo = @nNo AND nCharNo = @nCharNo
SET @nRet = @@ROWCOUNT

IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
BEGIN
	SET @nRet = 0
	RETURN
END

-- end

0뇇ĭ
GO

-- ========================================================================
-- Procedure: p_Guild_Member_SetGrade
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Guild_Member_SetGrade]
	@nNo     int
,	@nCharNo int
,	@nGrade  tinyint
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	DECLARE @cnt int, @guildMaster int
	DECLARE @lastChangeDate datetime

	SET @nRet = 0
	SET @guildMaster = 0

	IF ISNULL(@nNo, 0) = 0 BEGIN
		SET @nRet = -1311
		RETURN
	END
	IF ISNULL(@nCharNo, 0) = 0 BEGIN
		SET @nRet = -1313
		RETURN
	END
	
	-- Ä³ž¯ÅÍží º¯°æ Œ­ºñœº(ÀÎ°£°ü°è ÃÊ±âÈ­) ŸÆÀÌÅÛ Ãß°¡·Î ÀÎÇÑ ŒöÁ€ by kodong
	SELECT @cnt = ISNULL(COUNT(nCharNo), 0) 
	  FROM tCharacterChangeID 
	 WHERE nCharNo = @nCharNo 
	   AND bFlag = 0 -- ¹Ý¿µ ŸÈµÈ ±â·ÏÀÌ ÀÖŽÂ °æ¿ì
	   AND bInit = 1 -- ÀÎ°£°ü°è ÃÊ±âÈ­žŠ Œ³Á€ÇÑ °æ¿ì
	IF @nGrade = 0 AND @cnt = 1 BEGIN
		SET @nRet = -1323
		RETURN
	END
	ELSE BEGIN
		BEGIN TRAN
		
		-- ±æµåž¶œºÅÍ º¯°æ ¿äÃ»œÃ ±âÁž ±æµåž¶œºÅÍ Ä³ž¯ÅÍ¹øÈ£žŠ °¡Á®¿È
		IF @nGrade = 0 BEGIN
			SELECT @guildMaster = ISNULL(MAX(nCharNo), 0) FROM tGuildMember WHERE nNo = @nNo AND nGrade = 0
			
			-- ±æµåž¶œºÅÍ ÀÌŸç ÄðÅžÀÓ ÃŒÅ© by kodong 2011.9.30
			SELECT @lastChangeDate = MAX(dDate) FROM tLog_GuildMasterChange WITH(HOLDLOCK) WHERE nGuildNo = @nNo
			IF @lastChangeDate IS NOT NULL AND DATEADD(Day, 7, @lastChangeDate) > GETDATE() BEGIN
				ROLLBACK TRAN
				SET @nRet = -1324
				RETURN
			END
		END	
		
		-- ±æµåž¶œºÅÍ º¯°æ ¿äÃ»°ú ±âÁž ±æµåž¶œºÅÍ Ä³ž¯ÅÍ¹øÈ£°¡ ÁžÀçÇÏžé ±âÁž ±æµåž¶œºÅÍ žðµÎžŠ ÀÏ¹Ý ±æµå¿øÀž·Î º¯°æ		
		IF @nGrade = 0 AND @guildMaster <> 0 BEGIN
			UPDATE tGuildMember SET nGrade = 6 WHERE nNo = @nNo AND nGrade = 0
			IF @@ERROR <> 0 BEGIN
				ROLLBACK TRAN
				SET @nRet = -1325
				RETURN
			END
		END		
		
		-- ±æµå¿ø ±ÇÇÑ º¯°æ	
		UPDATE tGuildMember SET nGrade = @nGrade WHERE nNo = @nNo AND nCharNo = @nCharNo
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			ROLLBACK TRAN
			SET @nRet = -1325
		END
		ELSE BEGIN
			-- ±æµåž¶œºÅÍ ÀÌŸç ÄðÅžÀÓ ÃŒÅ© by kodong 2011.9.30				
			IF @nGrade = 0 BEGIN
				INSERT INTO tLog_GuildMasterChange (nGuildNo, nOldGuildMaster, nNewGuildMaster)
				VALUES (@nNo, @guildMaster, @nCharNo)
			END
			
			COMMIT TRAN
			SET @nRet = 0
		END
	END
END

0ﴁ烡Į
GO

-- ========================================================================
-- Procedure: p_Guild_Member_SetIntro
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Guild_Member_SetIntro]

/*
Set Guild Member Intro
2006.2 By CJC

Input:
	See SQL
	
Output:
	nRet = 0,  Error
	nRet = 1,  OK
*/

@nNo int,
@nCharNo int,
@sIntro nvarchar(64),

@nRet tinyint output

AS

SET NOCOUNT ON
SET @nRet = 0

UPDATE tGuildMember SET sIntro = @sIntro WHERE nNo = @nNo AND nCharNo = @nCharNo
SET @nRet = @@ROWCOUNT

IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
BEGIN
	SET @nRet = 0
	RETURN
END

-- end

0㘁攆į
GO

-- ========================================================================
-- Procedure: p_Guild_MoneyAdd
-- ========================================================================
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

-- ========================================================================
-- Procedure: p_Guild_MoneySub
-- ========================================================================
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

-- ========================================================================
-- Procedure: p_Guild_Set
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Guild_Set]

/*
Update Guild Data
2006.3 By blacksky
*/

@nNo				int,

@sName				nvarchar(18),
@sPassword			nvarchar(9),
@nMaxMembers		smallint,
@nMoney				bigint,
@nType				tinyint,
@nGrade				tinyint,	
@nFame				int,
@nStoneLevel		smallint,
@nExp				bigint,
@dCreateDate		datetime,

@dWarRequestDate	datetime,
@dWarStartDate		datetime,
@dWarEndDate		datetime,
@nWarEnemy			int,
@nWarWinCount		int,
@nWarLoseCount		int,
@nWarDrawCount		int,

@nDismissStatus		tinyint,
@dDismissDate		datetime,

@dNotifyDate		datetime,
@sNotifyCharID		nvarchar(16),
@sNotify			nvarchar(512),

@sIntro				nvarchar(128),

@nRet				tinyint output

AS

SET NOCOUNT ON
SET @nRet = 0

UPDATE
	tGuild

SET
	sName			= @sName,
	sPassword		= @sPassword,
	nMaxMembers		= @nMaxMembers,
	nMoney			= @nMoney,
	nType			= @nType,
	nGrade			= @nGrade,
	nFame			= @nFame,
	nStoneLevel		= @nStoneLevel,
	nExp			= @nExp,
	dCreateDate		= @dCreateDate,

	dWarRequestDate	= @dWarRequestDate,
	dWarStartDate	= @dWarStartDate,
	dWarEndDate		= @dWarEndDate,
	nWarEnemy		= @nWarEnemy,
	nWarWinCount	= @nWarWinCount,
	nWarLoseCount	= @nWarLoseCount,
	nWarDrawCount	= @nWarDrawCount,
	
	nDismissStatus	= @nDismissStatus,
	dDismissDate	= @dDismissDate,

	dNotifyDate		= @dNotifyDate,
	sNotifyCharID	= @sNotifyCharID,
	sNotify			= @sNotify,
	
	sIntro			= @sIntro
	
WHERE
	@nNo = nNo

SET @nRet = @@ROWCOUNT

IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
BEGIN
	SET @nRet = 0
	RETURN
END
-- end

0ଁ白Ĳ
GO

-- ========================================================================
-- Procedure: p_Guild_SetDismiss
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Guild_SetDismiss]

/*
Set Guild Type
2006.2 By CJC

Input:
	See SQL
	
Output:
	nRet = 0,  Error(Not found guild)
	nRet = 1,  OK
*/

@nNo int,
@nDismissStatus tinyint,
@dDismissDate datetime,

@nRet tinyint output

AS

SET NOCOUNT ON
SET @nRet = 0

UPDATE tGuild SET nDismissStatus = @nDismissStatus, dDismissDate = @dDismissDate WHERE nNo = @nNo
SET @nRet = @@ROWCOUNT

IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
BEGIN
	SET @nRet = 0
	RETURN
END

-- end

0渁龫ĳ
GO

-- ========================================================================
-- Procedure: p_Guild_SetExp
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Guild_SetExp]

/*
Set Guild Exp
2006.2 By CJC

Input:
	See SQL
	
Output:
	nRet = 0,  Error(Not found guild)
	nRet = 1,  OK
*/

@nNo int,
@nExp bigint,

@nRet tinyint output

AS

SET NOCOUNT ON
SET @nRet = 0

UPDATE tGuild SET nExp = @nExp WHERE nNo = @nNo
SET @nRet = @@ROWCOUNT

IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
BEGIN
	SET @nRet = 0
	RETURN
END

-- end

0꜁鏏Ĵ
GO

-- ========================================================================
-- Procedure: p_Guild_SetFame
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Guild_SetFame]

/*
Set Guild Fame
2006.2 By CJC

Input:
	See SQL
	
Output:
	nRet = 0,  Error(Not found guild)
	nRet = 1,  OK
*/

@nNo int,
@nFame int,

@nRet tinyint output

AS

SET NOCOUNT ON
SET @nRet = 0

UPDATE tGuild SET nFame = @nFame WHERE nNo = @nNo
SET @nRet = @@ROWCOUNT

IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
BEGIN
	SET @nRet = 0
	RETURN
END

-- end

0蟳ĵ
GO

-- ========================================================================
-- Procedure: p_Guild_SetGrade
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Guild_SetGrade]

/*
Set Guild Grade
2006.2 By CJC

Input:
	See SQL
	
Output:
	nRet = 0,  Error(Not found guild)
	nRet = 1,  OK
*/

@nNo int,
@nGrade tinyint,

@nRet tinyint output

AS

SET NOCOUNT ON
SET @nRet = 0

UPDATE tGuild SET nGrade = @nGrade WHERE nNo = @nNo
SET @nRet = @@ROWCOUNT

IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
BEGIN
	SET @nRet = 0
	RETURN
END

-- end

0䐁檡ĳ
GO

-- ========================================================================
-- Procedure: p_Guild_SetIntro
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Guild_SetIntro]

/*
Set Guild Intro
2006.2 By CJC

Input:
	See SQL
	
Output:
	nRet = 0,  Error(Not found guild)
	nRet = 1,  OK
*/

@nNo int,
@sIntro nvarchar(128),

@nRet tinyint output

AS

SET NOCOUNT ON
SET @nRet = 0

UPDATE tGuild SET sIntro = @sIntro WHERE nNo = @nNo
SET @nRet = @@ROWCOUNT

IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
BEGIN
	SET @nRet = 0
	RETURN
END

-- end

<䜍Ķ
GO

-- ========================================================================
-- Procedure: p_Guild_SetNotify
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Guild_SetNotify]

/*
Set Guild Notify
2006.2 By CJC

Input:
	See SQL
	
Output:
	nRet = 0,  Error(Not found guild)
	nRet = 1,  OK
*/

@nNo int,
@sNotifyCharID nvarchar(16),
@sNotify nvarchar(512),

@nRet tinyint output

AS

SET NOCOUNT ON
SET @nRet = 0

UPDATE tGuild SET dNotifyDate = getdate(), sNotifyCharID = @sNotifyCharID, sNotify = @sNotify WHERE nNo = @nNo
SET @nRet = @@ROWCOUNT

IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
BEGIN
	SET @nRet = 0
	RETURN
END

-- end

0愁⽖ĸ
GO

-- ========================================================================
-- Procedure: p_Guild_SetStoneLevel
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Guild_SetStoneLevel]

/*
Set Guild StoneLevel
2006.2 By CJC

Input:
	See SQL
	
Output:
	nRet = 0,  Error(Not found guild)
	nRet = 1,  OK
*/

@nNo int,
@nStoneLevel smallint,

@nRet tinyint output

AS

SET NOCOUNT ON
SET @nRet = 0

UPDATE tGuild SET nStoneLevel = @nStoneLevel WHERE nNo = @nNo
SET @nRet = @@ROWCOUNT

IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
BEGIN
	SET @nRet = 0
	RETURN
END

-- end

0ﴁ䲨ĺ
GO

-- ========================================================================
-- Procedure: p_Guild_SetType
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Guild_SetType]

/*
Set Guild Type
2006.2 By CJC

Input:
	See SQL
	
Output:
	nRet = 0,  Error(Not found guild)
	nRet = 1,  OK
*/

@nNo int,
@nType tinyint,

@nRet tinyint output

AS

SET NOCOUNT ON
SET @nRet = 0

UPDATE tGuild SET nType = @nType WHERE nNo = @nNo
SET @nRet = @@ROWCOUNT

IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
BEGIN
	SET @nRet = 0
	RETURN
END

-- end

0騁⍺Ĺ
GO

-- ========================================================================
-- Procedure: p_Guild_SetWar
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Guild_SetWar]

/*
Set Guild War
2006.2 By CJC

Input:
	See SQL
	
Output:
	nRet = 0,  Error(Not found guild)
	nRet = 1,  OK
*/

@nNo int,
@dWarRequestDate datetime,
@dWarStartDate datetime,
@dWarEndDate datetime,
@nWarEnemy int,

@nRet tinyint output

AS

SET NOCOUNT ON
SET @nRet = 0

UPDATE tGuild SET 
	dWarRequestDate = @dWarRequestDate,
	dWarStartDate = @dWarStartDate,
	dWarEndDate = @dWarEndDate,
	nWarEnemy = @nWarEnemy
WHERE nNo = @nNo
SET @nRet = @@ROWCOUNT

IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
BEGIN
	SET @nRet = 0
	RETURN
END

-- end

0态痗Ļ
GO

-- ========================================================================
-- Procedure: p_Guild_SetWarResult
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Guild_SetWarResult]

/*
Set Guild War Result
2006.2 By CJC

Input:
	See SQL
	
Output:
	nRet = 0,  Error(Not found guild)
	nRet = 1,  OK
*/

@nNo int,
@nWarWinCount int,
@nWarLoseCount int,
@nWarDrawCount int,

@nRet tinyint output

AS

SET NOCOUNT ON
SET @nRet = 0

UPDATE tGuild SET 
	nWarWinCount = @nWarWinCount,
	nWarLoseCount = @nWarLoseCount,
	nWarDrawCount = @nWarDrawCount
WHERE nNo = @nNo
SET @nRet = @@ROWCOUNT

IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
BEGIN
	SET @nRet = 0
	RETURN
END

-- end

0ఁ௃Ļ
GO

-- ========================================================================
-- Procedure: p_Guild_setName
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Guild_setName]
	@nNo   int
,	@sName nvarchar(18)
,	@nRet  int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	IF @nNo <= 0 BEGIN
		SET @nRet = -1311
		RETURN
	END

	IF LEN(ISNULL(@sName, '')) = 0 BEGIN
		SET @nRet = -1312
		RETURN
	END

	IF EXISTS(SELECT nNo FROM tGuild WHERE nNo = @nNo) BEGIN
		IF EXISTS(SELECT sName FROM tGuild WHERE sName = @sName) BEGIN
			SET @nRet = -1322
		END
		ELSE BEGIN

			BEGIN TRAN

			UPDATE tGuild SET sName = @sName WHERE nNo = @nNo
			IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
				ROLLBACK TRAN
				SET @nRet = -1393
			END
			ELSE BEGIN
				COMMIT TRAN
				SET @nRet = 0
			END
		END
	END
	ELSE BEGIN
		SET @nRet = -1321
	END

	SET NOCOUNT OFF
END

0逷7
GO

-- ========================================================================
-- Procedure: p_HolyPromise_CheckRing
-- ========================================================================
CREATE PROCEDURE [dbo].[p_HolyPromise_CheckRing]
/*
Holy promise check ring
2008.4 By CJC

 input:
	see SQL
	
 output:
	@nRet
	0 = No ring 
	1 = Yes ring
	-1 = DB error
*/
@nCharNo int,		-- 현재 캐릭
@nUpCharNo int,	-- 상위로 등록될 캐릭
@nRet int OUTPUT	-- 아미 상위로 등록되어 링구조가 되는지 결과확인값 리턴

AS

BEGIN
	SET NOCOUNT ON

	-- check are same char.
	IF @nCharNo = @nUpCharNo
	BEGIN
		SET @nRet = 1
		RETURN
	END

	DECLARE @nTempCharNo int
	DECLARE @nTempUpCharNo int

	SET @nTempCharNo = @nCharNo

	LOOP_UP:
	SET @nTempUpCharNo = 0
	SELECT @nTempUpCharNo = P.nUpCharNo FROM tHolyPromise AS P 
	INNER JOIN tCharacter AS C ON C.nCharNo =  P.nCharNo AND C.bDeleted = 0
	WHERE P.nCharNo = @nTempCharNo

	IF @nTempUpCharNo > 0
	BEGIN
		-- 링구조가 형성되는지 검사
		IF @nUpCharNo = @nTempUpCharNo
		BEGIN
			SET @nRet = 1
			RETURN
		END

		-- 다음 상위로 이동
		SET @nTempCharNo = @nTempUpCharNo
		GOTO LOOP_UP
	END

	-- 링구조 없음.
	SET @nRet = 0
	RETURN

END

0똁⺰Ł
GO

-- ========================================================================
-- Procedure: p_HolyPromise_Del
-- ========================================================================
CREATE PROCEDURE [dbo].[p_HolyPromise_Del]
/*
Delete holy promise
2008.4 By CJC

 input:
	see SQL
	
 output:
	@nRet
	0 = OK, 
	? = Other DB error code
	        
*/
@nCharNo int,
@nRet int OUTPUT

AS

BEGIN
	SET NOCOUNT ON

	-- Update/Inset data
	IF EXISTS( SELECT nCharNo FROM  tHolyPromise WHERE nCharNo=@nCharNo)
	BEGIN
		UPDATE tHolyPromise SET nUpCharNo=0, nFlags=0, dRegDate=GetDate() WHERE nCharNo=@nCharNo
		SET @nRet = @@ERROR
	END
	ELSE
		SET @nRet = 0

	RETURN
END

0
GO

-- ========================================================================
-- Procedure: p_HolyPromise_DelChar
-- ========================================================================
CREATE PROCEDURE [dbo].[p_HolyPromise_DelChar]
/*
Delete char holy promise
2008.4 By CJC

 input:
	see SQL
	
 output:
	@nRet
	0 = OK, 
	? = Other DB error code
	        
*/
@nDelCharNo int,
@nRet int OUTPUT

AS

BEGIN
	SET NOCOUNT ON

	DELETE FROM tHolyPromise WHERE nUpCharNo = @nDelCharNo
	SET @nRet = @@ERROR

	RETURN
END

0
GO

-- ========================================================================
-- Procedure: p_HolyPromise_Get
-- ========================================================================
CREATE PROCEDURE [dbo].[p_HolyPromise_Get]
/*
Get holy primise
2008.3 By CJC
 input:
	see SQL
 output:
	nRet > 0 fetch to take data
	= 0 No data
*/

@nCharNo int,
@nRet int OUTPUT

AS
BEGIN
	SET NOCOUNT ON
	SELECT @nRet=COUNT(nUpCharNo) FROM tHolyPromise WHERE nCharNo = @nCharNo
	IF @nRet > 0
	BEGIN
		IF EXISTS(
			SELECT P.nUpCharNo, P.nFlags, P.dRegDate, P.dRejoinableDate, P.dMemberAcceptableDate, C.sID, C.nLevel,C.dLastLoginDate FROM tHolyPromise AS P
			INNER JOIN tCharacter AS C ON C.nCharNo =  P.nUpCharNo AND C.bDeleted = 0
			WHERE P.nCharNo = @nCharNo)

			SELECT P.nUpCharNo, P.nFlags, P.dRegDate, P.dRejoinableDate, P.dMemberAcceptableDate, C.sID, C.nLevel,C.dLastLoginDate FROM tHolyPromise AS P
			INNER JOIN tCharacter AS C ON C.nCharNo =  P.nUpCharNo AND C.bDeleted = 0
			WHERE P.nCharNo = @nCharNo

		ELSE
			SELECT P.nUpCharNo, P.nFlags, P.dRegDate, P.dRejoinableDate, P.dMemberAcceptableDate, '', 1, GetDate() FROM tHolyPromise AS P WHERE P.nCharNo = @nCharNo

	END
END

0
GO

-- ========================================================================
-- Procedure: p_HolyPromise_GetMember
-- ========================================================================
CREATE PROCEDURE [dbo].[p_HolyPromise_GetMember]
/*
Get primise member count
2008.3 By CJC
 input:
	see SQL
 output:
	fetch to take data
*/

@nCharNo int,
@nNumOfMember int OUTPUT

AS
BEGIN
	SET NOCOUNT ON

	SELECT @nNumOfMember=COUNT( P.nCharNo) FROM tHolyPromise AS P
	INNER JOIN tCharacter AS C ON C.nCharNo =  P.nCharNo AND C.bDeleted = 0
	WHERE P.nUpCharNo = @nCharNo

	IF @nNumOfMember > 0 
	BEGIN
		SELECT P.nCharNo, P.nFlags, P.dRegDate, P.dRejoinableDate, P.dMemberAcceptableDate, C.sID, C.nLevel,C.dLastLoginDate FROM tHolyPromise AS P
		INNER JOIN tCharacter AS C ON C.nCharNo =  P.nCharNo AND C.bDeleted = 0
		WHERE P.nUpCharNo = @nCharNo
	END
END

0
GO

-- ========================================================================
-- Procedure: p_HolyPromise_GetMemberCount
-- ========================================================================
CREATE PROCEDURE [dbo].[p_HolyPromise_GetMemberCount]
/*
Get primise member count
2008.3 By CJC
 input:
	see SQL
 output:
	fetch to take data
*/

@nCharNo int,
@nNumOfMember int OUTPUT

AS
BEGIN
	SET NOCOUNT ON
	SET @nNumOfMember = 0
	SELECT @nNumOfMember = COUNT(P.nUpCharNo) FROM tHolyPromise AS P 
	INNER JOIN tCharacter AS C ON C.nCharNo =  P.nCharNo AND C.bDeleted = 0
	WHERE P.nUpCharNo = @nCharNo
END

<
GO

-- ========================================================================
-- Procedure: p_HolyPromise_Set
-- ========================================================================
CREATE PROCEDURE [dbo].[p_HolyPromise_Set]
/*
Holy promise set
2008.4 By CJC

 input:
	see SQL
	
 output:
	@nRet
	0 = OK,  Fiest join
	1 = Ring / same char
	2 = i'm over max member count
	3 = OK, ReJoin
	4 = Failed, Already upper exists
	? = Other DB error code
*/
@nMAXPROMISEMEMBER smallint,
@nCharNo int,
@nUpCharNo int,
@nFlags tinyint,
@nRet int OUTPUT

AS

BEGIN
	SET NOCOUNT ON
	
	-- check ring, same char.
	EXEC p_HolyPromise_CheckRing @nCharNo, @nUpCharNo, @nRet OUTPUT
	IF @nRet = 1
		RETURN
	
	-- check over max promise and add
	BEGIN TRAN

	DECLARE @nPromiseCount int
	EXEC p_HolyPromise_GetMemberCount @nUpCharNo, @nPromiseCount OUTPUT
	IF @nPromiseCount >= @nMAXPROMISEMEMBER
	BEGIN
		SET @nRet = 2
		ROLLBACK TRAN
		RETURN
	END

	-- Update/Inset data
	IF EXISTS( SELECT nCharNo FROM  tHolyPromise WHERE nCharNo=@nCharNo)
	BEGIN
		-- Check already upper here
		DECLARE @TempUpCharNo int
		SELECT @TempUpCharNo=nUpCharNo FROM  tHolyPromise WHERE nCharNo=@nCharNo
		IF @TempUpCharNo <> 0
		BEGIN
			SET @nRet = 4
			ROLLBACK TRAN
			RETURN
		END

		UPDATE tHolyPromise SET nUpCharNo=@nUpCharNo, nFlags=@nFlags, dRegDate=GetDate() WHERE nCharNo=@nCharNo
		SET @nRet = @@ERROR
		IF @nRet <> 0
		BEGIN
			ROLLBACK TRAN
			RETURN
		END
		SET @nRet = 3 -- OK. Rejoin		
	END

	ELSE

	BEGIN
		INSERT tHolyPromise (  nCharNo,  nUpCharNo,  nFlags )
		VALUES ( @nCharNo, @nUpCharNo, @nFlags )
		SET @nRet = @@ERROR
		IF @nRet <> 0
		BEGIN
			ROLLBACK TRAN
			RETURN
		END
	END
	
	COMMIT TRAN
	RETURN
	-- end
END

0
GO

-- ========================================================================
-- Procedure: p_HolyPromise_SetDate
-- ========================================================================
CREATE PROCEDURE [dbo].[p_HolyPromise_SetDate]
/*
Holy promise set date
2008.4 By CJC

 input:
	see SQL
	
 output:
	@nRet
	0 = OK, 
	? = Other DB error code
*/
@nCharNo int,
@dRejoinableDate datetime,
@dMemberAcceptableDate datetime,
@nRet int OUTPUT

AS

BEGIN
	SET NOCOUNT ON

	-- Update data
	IF EXISTS( SELECT nCharNo FROM  tHolyPromise WHERE nCharNo=@nCharNo)
	BEGIN
		UPDATE tHolyPromise SET dRejoinableDate=@dRejoinableDate, dMemberAcceptableDate=@dMemberAcceptableDate WHERE nCharNo=@nCharNo
		SET @nRet = @@ERROR
	END
	ELSE	
	BEGIN
		INSERT tHolyPromise (  nCharNo,  nUpCharNo,  nFlags, dRejoinableDate, dMemberAcceptableDate )
		VALUES ( @nCharNo, 0, 0, @dRejoinableDate, @dMemberAcceptableDate )
		SET @nRet = @@ERROR
	END
	
	RETURN
END

0
GO

-- ========================================================================
-- Procedure: p_Item_BuyAll
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Item_BuyAll]
	@nItemKey bigint,
	@nOwner int,
	@nStorageType tinyint,
	@nStorage smallint,
	@nItemID int,
	@nFlags int,
	@nPrice bigint,
	@nFame int,
	@nCoin bigint,
	@nRet int OUTPUT
AS
BEGIN
	SET NOCOUNT ON

	DECLARE @CreateItemRet int

	IF @nPrice <= 0 AND @nFame <= 0 AND @nCoin <= 0 BEGIN
		SET @nRet = 2
		RETURN
	END

	BEGIN TRAN
	
	-- ??? ???? ??? ?? ??? ??? ?? ? ? ??
	IF @nPrice > 0 BEGIN
		EXEC p_Char_MoneySub @nOwner, @nPrice, @nRet OUTPUT
		IF @nRet <> 1
		BEGIN
			ROLLBACK TRAN
			SET @nRet = 1
			RETURN
		END
	END
	
	-- ??? ???? ??? ??? ??? ??? ?? ? ?? ??
	IF @nFame > 0 BEGIN
		EXEC p_Char_FameSub @nOwner, @nFame, @nRet OUTPUT
		IF @nRet <> 1
		BEGIN
			ROLLBACK TRAN
			SET @nRet = 1
			RETURN
		END
	END
	
	-- ??? ???? ??? ??? ??? ??? ?? ? ?? ??
	IF @nCoin > 0 BEGIN
		EXEC usp_Gamble_subCoin @nOwner, @nCoin, @nRet OUTPUT
		IF @nRet <> 0
		BEGIN
			ROLLBACK TRAN
			SET @nRet = 1
			RETURN
		END
	END
	
	-- ??? ?? ?????? ??
	EXEC p_Item_Create @nItemKey, @nOwner, @nStorageType, @nStorage, @nItemID, @nFlags, @CreateItemRet OUTPUT
	IF @CreateItemRet <> 0 OR @@ERROR <> 0
	BEGIN
		ROLLBACK TRAN
		SET @nRet = @CreateItemRet
		RETURN
	END
	-- ??
	COMMIT TRAN
	SET @nRet = 0
	-- end
END

<
GO

-- ========================================================================
-- Procedure: p_Item_BuyLot
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Item_BuyLot]
	@nItemKey bigint,
	@nOwner int,
	@nOptionType smallint,
	@nOptionData bigint,
	@nPrice bigint,
	@nFame int,
	@nCoin bigint,
	@nRet int OUTPUT
AS
BEGIN
	SET NOCOUNT ON

	IF @nPrice <= 0 AND @nFame <= 0 AND @nCoin <= 0 BEGIN
		SET @nRet = 2
		RETURN
	END

	BEGIN TRAN
	
	-- ??? ???? ??? ?? ??? ??? ?? ? ? ??
	IF @nPrice > 0 BEGIN
		EXEC p_Char_MoneySub @nOwner, @nPrice, @nRet OUTPUT
		IF @nRet <> 1
		BEGIN
			ROLLBACK TRAN
			SET @nRet = 1
			RETURN
		END
	END
	
	-- ??? ???? ??? ??? ??? ??? ?? ? ?? ??
	IF @nFame > 0 BEGIN
		EXEC p_Char_FameSub @nOwner, @nFame, @nRet OUTPUT
		IF @nRet <> 1
		BEGIN
			ROLLBACK TRAN
			SET @nRet = 1
			RETURN
		END
	END
	
	-- ??? ???? ??? ??? ??? ??? ?? ? ?? ??
	IF @nCoin > 0 BEGIN
		EXEC usp_Gamble_subCoin @nOwner, @nCoin, @nRet OUTPUT
		IF @nRet <> 0
		BEGIN
			ROLLBACK TRAN
			SET @nRet = 1
			RETURN
		END
	END
	
	-- ??? ?? ?????? ??(? ??)
	EXEC p_Item_SetOption @nItemKey, @nOptionType, @nOptionData, @nRet OUTPUT
	IF @nRet <> 0
	BEGIN
		ROLLBACK TRAN
		RETURN
	END
	-- ??
	COMMIT TRAN
	-- end
END

0
GO

-- ========================================================================
-- Procedure: p_Item_ClrWeaponTitle
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Item_ClrWeaponTitle]
/*

Item Clear Weapon Title
2006.6 By CJC

 input:
	see SQL

 output:
	@nRet	0 = OK, 
	        1 = Error Not found Weapon mob 
	        2 = Error Clear Weapon mob 
	        3 = Error Delete license item
	        4 = Error Change license item lot
*/

@nWeaponItemKey bigint,

@nMobID smallint,
@nLicenseClearItemKey bigint,
@bLicenseClearItemDelete tinyint,
@nLicenseClearItemOptionType smallint,
@nLicenseClearItemOptionData bigint,

-- Output var
@nRet int OUTPUT

AS

SET NOCOUNT ON

-- 아이템에 몹 라이센스 붙어 있는지 확인 
IF NOT EXISTS(SELECT nMobID FROM tItemMobList WHERE nItemKey=@nWeaponItemKey AND nMobID=@nMobID)
BEGIN
	SET @nRet = 1
	RETURN
END

-- 아이템에 몹 라이센스 지우기
BEGIN TRAN
DELETE FROM tItemMobList WHERE nItemKey=@nWeaponItemKey AND nMobID=@nMobID
IF @@ERROR <> 0
BEGIN
	ROLLBACK TRAN
	SET @nRet = 2
	RETURN
END

-- 라이센스 지우는 아이템 지우기/랏조정
DECLARE @nRetTemp int
IF @bLicenseClearItemDelete = 1
BEGIN
	EXEC p_Item_Delete @nLicenseClearItemKey, @nRetTemp OUTPUT
	IF @nRetTemp <> 0
	BEGIN
		ROLLBACK TRAN
		SET @nRet = 3
		RETURN
	END
END
ELSE
BEGIN
	EXEC p_Item_SetOption @nLicenseClearItemKey, @nLicenseClearItemOptionType, @nLicenseClearItemOptionData, @nRetTemp OUTPUT
	IF @nRetTemp <> 0
	BEGIN
		ROLLBACK TRAN
		SET @nRet = 4
		RETURN
	END
END

-- 성공
COMMIT TRAN
SET @nRet = 0

RETURN

-- end

0렚Ŋ
GO

-- ========================================================================
-- Procedure: p_Item_CostumWeapon_Endure_Set
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Item_CostumWeapon_Endure_Set]
/*
Item Costum Weapon Endure Set
2008.10 By Lunar
 input:
	see SQL
 output:
	@nRet	0 = OK
	        1 = Not found item
	        2 = Missmatch item owner
	        3 = Not enough money to sub
	        4 = Upgrade value setting, upgrade fail count error(option set error)
	        5 = Raw delete fail
	        6 = RawLeft delete fail, if valid item
	        7 = RawRight delete fail, if valid item
	        8 = RawMiddle delete fail, if valid item
	        
 operation:
	1. money sub
	2. upgrade item level
	3. delete raws
*/
@nCostumWeaponKey bigint,
@nCostumWeaponEndureOptionType smallint,
@nCostumWeaponEndureOptionData bigint,
@nRet int OUTPUT
AS
BEGIN
	SET NOCOUNT ON
	
	DECLARE @nTinyintRet tinyint
	DECLARE @nRetTemp int

	BEGIN TRAN

	-- 1 가구 아이템 내구도 상승
	SET @nRetTemp = 1
	EXEC p_Item_SetOption @nCostumWeaponKey, @nCostumWeaponEndureOptionType, @nCostumWeaponEndureOptionData, @nRetTemp OUTPUT
	IF @nRetTemp <> 0 BEGIN
		ROLLBACK TRAN
		SET @nRet = 5
		RETURN
	END
	
	-- 업글 성공
	COMMIT TRAN
	SET @nRet = 0

	SET NOCOUNT OFF
END

0圁ꌌJ
GO

-- ========================================================================
-- Procedure: p_Item_CostumWeapon_Endure_Update
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Item_CostumWeapon_Endure_Update]
/*
Item Costum Weapon Endure Update
2008.10 By Lunar
 input:
	see SQL
 output:
	@nRet	0 = OK
	        1 = Not found item
	        2 = Missmatch item owner
	        3 = Not enough money to sub
	        4 = Upgrade value setting, upgrade fail count error(option set error)
	        5 = Raw delete fail
	        6 = RawLeft delete fail, if valid item
	        7 = RawRight delete fail, if valid item
	        8 = RawMiddle delete fail, if valid item
	        
 operation:
	1. money sub
	2. upgrade item level
	3. delete raws
*/
@nOwner int,
@nEndureKitKey bigint,
@nEndureKitOptionType smallint,
@nEndureKitLot smallint,
@nCostumWeaponKey bigint,
@nCostumWeaponEndureOptionType smallint,
@nCostumWeaponEndureOptionData bigint,
@nRet int OUTPUT
AS
BEGIN
	SET NOCOUNT ON
	
	DECLARE @nTinyintRet tinyint
	DECLARE @nRetTemp int

	-- 1.1 내구도 키트 아이템 소유자 확인
	EXEC p_Item_IsValidOwner @nEndureKitKey, @nOwner, @nTinyintRet OUTPUT
	IF @nTinyintRet <> 0 BEGIN
		SET @nRet = @nTinyintRet
		RETURN
	END
	
	-- 1.2 무기 코스튬 아이템 소유자 확인
	EXEC p_Item_IsValidOwner @nCostumWeaponKey, @nOwner, @nTinyintRet OUTPUT
	IF @nTinyintRet <> 0 BEGIN
		SET @nRet = @nTinyintRet
		RETURN
	END

	BEGIN TRAN
	
	-- 2.1 사용된 내구도 키트 아이템 지우기 (사용후 Lot == 0 인 경우)
	IF @nEndureKitLot = 0 BEGIN
		SET @nRetTemp = 1
		EXEC p_Item_Delete @nEndureKitKey, @nRetTemp OUTPUT
		IF @nRetTemp <> 0 BEGIN
			ROLLBACK TRAN
			SET @nRet = 5
			RETURN
		END
	END
	-- 2.2 사용된 내구도 키트 아이템 차감하기 (사용후 Lot > 0 인 경우)
	ELSE BEGIN
		SET @nRetTemp = 1
		EXEC p_Item_SetOption @nEndureKitKey, @nEndureKitOptionType, @nEndureKitLot, @nRetTemp OUTPUT
		IF @nRetTemp <> 0 BEGIN
			ROLLBACK TRAN
			SET @nRet = 5
			RETURN
		END
	END
	
	-- 3 가구 아이템 내구도 상승
	SET @nRetTemp = 1
	EXEC p_Item_SetOption @nCostumWeaponKey, @nCostumWeaponEndureOptionType, @nCostumWeaponEndureOptionData, @nRetTemp OUTPUT
	IF @nRetTemp <> 0 BEGIN
		ROLLBACK TRAN
		SET @nRet = 5
		RETURN
	END
	
	-- 업글 성공
	COMMIT TRAN
	SET @nRet = 0

	SET NOCOUNT OFF
END

0ก㙏Ō
GO

-- ========================================================================
-- Procedure: p_Item_Create
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Item_Create]
/*
Item Create
2004.8 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK, 
	        ? = Error code
*/
@nItemKey bigint,
@nOwner int,
@nStorageType tinyint,
@nStorage smallint,
@nItemID int,
@nFlags int,
@nRet int OUTPUT
AS
SET NOCOUNT ON
-- 아이템 생성
INSERT tItem (  nItemKey,  nOwner,  nStorageType,  nStorage,  nItemID,  nFlags )
VALUES       ( @nItemKey, @nOwner, @nStorageType, @nStorage, @nItemID, @nFlags )
SET @nRet = @@ERROR
-- end
RETURN

0䜁⩳ō
GO

-- ========================================================================
-- Procedure: p_Item_Create_UseEmptyStorage
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Item_Create_UseEmptyStorage]
/*
Item Create
2008.4 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK, 
	        ? = Error code
*/
@nItemKey bigint,
@nOwner int,
@nStorageType tinyint,
@nStorage smallint OUTPUT,
@nItemID int,
@nFlags int,
@nRet int OUTPUT
AS
BEGIN
	SET NOCOUNT ON

	BEGIN TRAN

	-- 빈인벤 위치찾기
	EXEC  p_Item_GetStorageEmpty @nOwner, @nStorageType, @nStorage OUTPUT
	IF @nStorage < 0
	BEGIN
		SET @nRet = -1
		ROLLBACK TRAN
		RETURN
	END

	-- 아이템 생성
	INSERT tItem (  nItemKey,  nOwner,  nStorageType,  nStorage,  nItemID,  nFlags )
	VALUES       ( @nItemKey, @nOwner, @nStorageType, @nStorage, @nItemID, @nFlags )
	SET @nRet = @@ERROR
	IF @nRet = 0
		COMMIT TRAN
	ELSE
		ROLLBACK TRAN
END

0퐁被Ŏ
GO

-- ========================================================================
-- Procedure: p_Item_Create_UseEmptyStorage_tmp
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Item_Create_UseEmptyStorage_tmp]
/*
Item Create
2008.4 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK, 
	        ? = Error code
modified by Adam.
when @nStorage <0,  print the charno + itemkey.
*/
@nItemKey bigint,
@nOwner int,
@nStorageType tinyint,
@nStorage smallint OUTPUT,
@nItemID int,
@nFlags int,
@nRet int OUTPUT
AS
BEGIN
	SET NOCOUNT ON

	BEGIN TRAN

	-- 빈인벤 위치찾기
	EXEC  p_Item_GetStorageEmpty @nOwner, @nStorageType, @nStorage OUTPUT
	IF @nStorage < 0
	BEGIN
		PRINT cast(@nOwner as varchar(20)) + ':' + cast(@nItemKey as varchar(20))
		ROLLBACK TRAN
		RETURN
	END

	-- 아이템 생성
	INSERT tItem (  nItemKey,  nOwner,  nStorageType,  nStorage,  nItemID,  nFlags )
	VALUES       ( @nItemKey, @nOwner, @nStorageType, @nStorage, @nItemID, @nFlags )
	SET @nRet = @@ERROR
	IF @nRet = 0
		BEGIN			
		EXEC p_Item_SetOption @nItemKey, 300,1,0
		COMMIT TRAN
		END
	ELSE
		BEGIN
		PRINT cast(@nOwner as varchar(20)) + ':' + cast(@nItemKey as varchar(20))
		ROLLBACK TRAN
		RETURN
		END
END

0뤁ኻŏ
GO

-- ========================================================================
-- Procedure: p_Item_Delete
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Item_Delete]
/*
Item Delete
2004.8 By CJC
 input:
	nItemKey = To be delete.
 output:
	@nRet	0 = OK, 
	        ? = Error code
*/
@nItemKey bigint,
@nRet int OUTPUT
AS
SET NOCOUNT ON
-- 아이템 삭제 옵션
DELETE FROM tItem WHERE nItemKey = @nItemKey
SET @nRet = @@ERROR
-- end

0۟Ő
GO

-- ========================================================================
-- Procedure: p_Item_ExchangeStorage
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Item_ExchangeStorage]
/*
Exchange item storage
2004.8 By CJC
 input:
	nItemKey1,nStorageType1,nStoage1 = #1 Item
	nItemKey2,nStorageType2,nStoage2 = #2 Item
 output:
	0 = Ok
	1~3 = Error step No.
 operation:
	Set nItemKey1 To nStorageType1,nStoage1
	Set nItemKey2 To nStorageType2,nStoage2
	
	1. #1 item move to temp storage   * #1,#2 is already input that exchanged data.
	2. set #2's storage to #2
	3. set #1's storage to #1
*/
@nItemKey1 bigint,
@nStorageType1 tinyint,
@nStorage1 smallint,
@nItemKey2 bigint,
@nStorageType2 tinyint,
@nStorage2 smallint,
@nRet tinyint OUTPUT
AS
SET NOCOUNT ON
SET @nRet = 0
BEGIN TRAN
-- Step 1
UPDATE tItem SET nStorageType = 11
WHERE nItemKey = @nItemKey1
IF @@ROWCOUNT = 0
BEGIN
	SET @nRet = 1
	ROLLBACK TRAN
	RETURN
END
	
-- Step 2
UPDATE tItem SET nStorageType = @nStorageType2, nStorage = @nStorage2
WHERE nItemKey = @nItemKey2
IF @@ROWCOUNT = 0
BEGIN
	SET @nRet = 2
	ROLLBACK TRAN
	RETURN
END
	
-- Step 3
UPDATE tItem SET nStorageType = @nStorageType1, nStorage = @nStorage1
WHERE nItemKey = @nItemKey1
IF @@ROWCOUNT = 0
BEGIN
	SET @nRet = 3
	ROLLBACK TRAN
	RETURN
END
COMMIT TRAN
-- end

0⬁ﬄŐ
GO

-- ========================================================================
-- Procedure: p_Item_ExchangeStorageAndOwner
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Item_ExchangeStorageAndOwner]
/*
Exchange item storage
2004.8 By CJC
 input:
	nOwner1,nItemKey1,nStorageType1,nStoage1 = #1 Item
	nOwner2,nItemKey2,nStorageType2,nStoage2 = #2 Item
 output:
	0 = Ok
	1~3 = Error step No.
 operation:
	Set nItemKey1 To nOwner1,nStorageType1,nStoage1
	Set nItemKey2 To nOwner2,nStorageType2,nStoage2
	
	1. #1 item move to temp storage   * #1,#2 is already input that exchanged data.
	2. set #2's storage to #2
	3. set #1's storage to #1
*/
@nItemKey1 bigint,
@nOwner1 int,
@nStorageType1 tinyint,
@nStorage1 smallint,
@nItemKey2 bigint,
@nOwner2 int,
@nStorageType2 tinyint,
@nStorage2 smallint,
@nRet tinyint OUTPUT
AS
SET NOCOUNT ON
SET @nRet = 0
BEGIN TRAN
-- Step 1
UPDATE tItem SET nStorageType = 11
WHERE nItemKey = @nItemKey1
IF @@ROWCOUNT = 0
BEGIN
	SET @nRet = 1
	ROLLBACK TRAN
	RETURN
END
	
-- Step 2
UPDATE tItem SET nStorageType = @nStorageType2, nStorage = @nStorage2, nOwner = @nOwner2
WHERE nItemKey = @nItemKey2
IF @@ROWCOUNT = 0
BEGIN
	SET @nRet = 2
	ROLLBACK TRAN
	RETURN
END
	
-- Step 3
UPDATE tItem SET nStorageType = @nStorageType1, nStorage = @nStorage1, nOwner = @nOwner1
WHERE nItemKey = @nItemKey1
IF @@ROWCOUNT = 0
BEGIN
	SET @nRet = 3
	ROLLBACK TRAN
	RETURN
END
COMMIT TRAN
-- end

0搁ő
GO

-- ========================================================================
-- Procedure: p_Item_Furniture_Endure_Update
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Item_Furniture_Endure_Update]
/*
Item Furniture Endure Update
2008.10 By Lunar
 input:
	see SQL
 output:
	@nRet	0 = OK
	        1 = Not found item
	        2 = Missmatch item owner
	        3 = Not enough money to sub
	        4 = Upgrade value setting, upgrade fail count error(option set error)
	        5 = Raw delete fail
	        6 = RawLeft delete fail, if valid item
	        7 = RawRight delete fail, if valid item
	        8 = RawMiddle delete fail, if valid item
	        
 operation:
	1. money sub
	2. upgrade item level
	3. delete raws
*/
@nOwner int,
@nEndureKitKey bigint,
@nEndureKitOptionType smallint,
@nEndureKitLot smallint,
@nFurnitureKey bigint,
@nFurnitureEndureOptionType smallint,
@nFurnitureEndureOptionData bigint,
@nRet int OUTPUT
AS
BEGIN
	SET NOCOUNT ON
	
	DECLARE @nTinyintRet tinyint
	DECLARE @nRetTemp int

	-- 1.1 내구도 키트 아이템 소유자 확인
	EXEC p_Item_IsValidOwner @nEndureKitKey, @nOwner, @nTinyintRet OUTPUT
	IF @nTinyintRet <> 0 BEGIN
		SET @nRet = @nTinyintRet
		RETURN
	END
	
	-- 1.2 가구 아이템 소유자 확인
	EXEC p_Item_IsValidOwner @nFurnitureKey, @nOwner, @nTinyintRet OUTPUT
	IF @nTinyintRet <> 0 BEGIN
		SET @nRet = @nTinyintRet
		RETURN
	END

	BEGIN TRAN
	
	-- 2.1 사용된 내구도 키트 아이템 지우기 (사용후 Lot == 0 인 경우)
	IF @nEndureKitLot = 0 BEGIN
		SET @nRetTemp = 1
		EXEC p_Item_Delete @nEndureKitKey, @nRetTemp OUTPUT
		IF @nRetTemp <> 0 BEGIN
			ROLLBACK TRAN
			SET @nRet = 5
			RETURN
		END
	END
	-- 2.2 사용된 내구도 키트 아이템 차감하기 (사용후 Lot > 0 인 경우)
	ELSE BEGIN
		SET @nRetTemp = 1
		EXEC p_Item_SetOption @nEndureKitKey, @nEndureKitOptionType, @nEndureKitLot, @nRetTemp OUTPUT
		IF @nRetTemp <> 0 BEGIN
			ROLLBACK TRAN
			SET @nRet = 5
			RETURN
		END
	END
	
	-- 3 가구 아이템 내구도 상승
	SET @nRetTemp = 1
	EXEC p_Item_SetOption @nFurnitureKey, @nFurnitureEndureOptionType, @nFurnitureEndureOptionData, @nRetTemp OUTPUT
	IF @nRetTemp <> 0 BEGIN
		ROLLBACK TRAN
		SET @nRet = 5
		RETURN
	END
	
	-- 업글 성공
	COMMIT TRAN
	SET @nRet = 0

	SET NOCOUNT OFF
END

0鴁Œ
GO

-- ========================================================================
-- Procedure: p_Item_GetData
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Item_GetData]
/*
Item Get Data
2004.10 By CJC
 input:
	nItemKey = Item Key
 output:
	@nRet	1 = OK, 
	        0 = Error
*/
@nItemKey bigint,
@nOwner int OUTPUT,
@nStorageType tinyint OUTPUT,
@nStorage smallint OUTPUT,
@nItemID int OUTPUT,
@nFlags int OUTPUT,
@dDate datetime OUTPUT,
@nRet int OUTPUT
AS
SET NOCOUNT ON
SET @nRet = 0
SELECT 
	@nOwner = nOwner,
	@nStorageType = nStorageType,
	@nStorage = nStorage,
	@nItemID = nItemID,
	@nFlags = nFlags,
	@dDate = dDate
FROM tItem
WHERE nItemKey = @nItemKey
SET @nRet = @@ROWCOUNT
-- end

0혁흰œ
GO

-- ========================================================================
-- Procedure: p_Item_GetList
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Item_GetList]
/*
Item Get list
2004.8 By CJC
 input:
	nCharNo	     = Item Owner
 output:
	none
 desc:
	fetch success, get data
	fetch fail, not data
*/
@nCharNo int
AS
SET NOCOUNT ON
SELECT nItemKey, nStorageType, nStorage, nItemID, nFlags, dDate
FROM tItem
WHERE nOwner = @nCharNo
-- end

0持㖩ŕ
GO

-- ========================================================================
-- Procedure: p_Item_GetListType
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Item_GetListType]
/*
Item Get list from StorageType
2004.8 By CJC
 input:
	nCharNo	     = Item Owner
	nStorageType = Item Storage Type
 output:
	none
 desc:
	fetch success, get data
	fetch fail, not data
*/
@nCharNo int,
@nStorageType tinyint
AS
SET NOCOUNT ON
SELECT nItemKey, nStorage, nItemID, nFlags, dDate
FROM tItem WITH(READUNCOMMITTED) 
WHERE nOwner = @nCharNo	AND nStorageType = @nStorageType
-- end

0鰁⧍Ŗ
GO

-- ========================================================================
-- Procedure: p_Item_GetMobList
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Item_GetMobList]
/*
Item Get Mob List
2004.8 By CJC
 input:
	see SQL
 output:
	@nRet	0 = No data fount
	        ? = Number data here(must to be 1)
	        
	@nGrade = grade
	@nKillCount = kill count(progress)
		        
*/
@nItemKey bigint,
@nMobID smallint,
-- Output var
@nGrade tinyint OUTPUT,
@nKillCount int OUTPUT,
@nRet int OUTPUT
AS
SET NOCOUNT ON
SELECT @nGrade = nGrade, @nKillCount = nKillCount 
FROM tItemMobList
WHERE nItemKey = @nItemKey AND nMobID = @nMobID
SET @nRet = @@ROWCOUNT
-- end

0픁ᷱŗ
GO

-- ========================================================================
-- Procedure: p_Item_GetMobListAll
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Item_GetMobListAll]
/*
Item Get Mob List All
2004.8 By CJC
 input:
	see SQL
 output:
	fetch to take data
*/
@nItemKey bigint
AS
SET NOCOUNT ON
SELECT nMobID, nGrade, nKillCount FROM tItemMobList WHERE nItemKey = @nItemKey
-- end

0먁ꠁŗ
GO

-- ========================================================================
-- Procedure: p_Item_GetNumOfItem
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Item_GetNumOfItem]
/*
Get number of item
2004.8 By CJC
 input:
	nCharNo	    = Item Owner
 output:
	nItemCount,	-1 = error
	             ? = Count of item;
*/
@nCharNo int,
@nNumOfItem int OUTPUT
AS
SET NOCOUNT ON
SET @nNumOfItem = -1
SELECT @nNumOfItem = COUNT(*) FROM tItem WHERE nOwner = @nCharNo
-- end

0䜁غř
GO

-- ========================================================================
-- Procedure: p_Item_GetOption
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Item_GetOption]
/*
Item Get Option
2004.8 By CJC
 input:
	see SQL
 output:
	@nRet	0 = No data found
	        ? = Number data here(must to be 1)
	        
	@nOptionData = option data
		        
*/
@nItemKey bigint,
@nOptionType smallint,
-- Output var
@nOptionData bigint OUTPUT,
@nRet int OUTPUT
AS
SET NOCOUNT ON
SELECT @nOptionData = nOptionData FROM tItemOptions WHERE nItemKey = @nItemKey AND nOptionType = @nOptionType
SET @nRet = @@ROWCOUNT
-- end

0輁른Ś
GO

-- ========================================================================
-- Procedure: p_Item_GetOptionAll
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Item_GetOptionAll]
/*
Item Get Option All
2004.8 By CJC
 input:
	see SQL
 output:
	fetch to take data
*/
@nItemKey bigint
AS
SET NOCOUNT ON
SELECT nOptionType, nOptionData FROM tItemOptions WHERE nItemKey = @nItemKey
-- end

0攁葮Ś
GO

-- ========================================================================
-- Procedure: p_Item_GetOptionCount
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Item_GetOptionCount]
/*
Item Get Option Count
2004.8 By CJC
 input:
	see SQL
 output:
	@nOptionCount = Item option count
*/
@nItemKey bigint,
@nOptionCount int OUTPUT
AS
SET NOCOUNT ON
SELECT @nOptionCount = COUNT(nItemKey) FROM tItemOptions WHERE nItemKey = @nItemKey
-- end

0鸁碒ś
GO

-- ========================================================================
-- Procedure: p_Item_GetStorage
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Item_GetStorage]
/*
Get item storage
2004.8 By CJC
 input:
	nItemKey	 = Item Key
 output:
	1, ok nStorageType, nStorage = valid value
	0, Not found item
*/
@nItemKey bigint,
@nStorageType tinyint OUTPUT,
@nStorage smallint OUTPUT,
@nRet tinyint = 1 OUTPUT 
AS
SET NOCOUNT ON
SELECT @nStorageType = nStorageType, @nStorage = nStorage FROM tItem WHERE nItemKey = @nItemKey
-- Check not found item
IF @@ROWCOUNT <> 1
	SET @nRet = 0
-- end

0⬁훋Ŝ
GO

-- ========================================================================
-- Procedure: p_Item_GetStorageAndOwner
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Item_GetStorageAndOwner]
/*
Get item storage
2004.8 By CJC
 input:
	nItemKey	 = Item Key
 output:
	1, ok nOwner, nStorageType, nStorage = valid value
	0, Not found item
*/
@nItemKey bigint,
@nOwner int OUTPUT,
@nStorageType tinyint OUTPUT,
@nStorage smallint OUTPUT,
@nRet tinyint OUTPUT 
AS
SET NOCOUNT ON
SELECT @nOwner = nOwner, @nStorageType = nStorageType, @nStorage = nStorage FROM tItem WHERE nItemKey = @nItemKey
-- Check not found item
IF @@ROWCOUNT <> 1
	SET @nRet = 0
ELSE
	SET @nRet = 1

-- end

0ခ惛ŝ
GO

-- ========================================================================
-- Procedure: p_Item_GetStorageEmpty
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Item_GetStorageEmpty]
	@owner       int = 0
,	@storageType tinyint = 0
,	@ret         int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	DECLARE @err AS bit
	DECLARE @emptyStorage AS smallint
	SET @err = 0

	IF @owner = 0
		SET @err = 1
	IF @storageType = 0
		SET @err = 1

	IF @err <> 1
		BEGIN	
			IF EXISTS(SELECT nItemKey FROM tItem WHERE nOwner = @owner AND nStorageType = @storageType AND nStorage = 0)
				BEGIN
					SELECT @emptyStorage = ISNULL(MIN(nStorage + 1), 0)
					  FROM tItem 
					 WHERE nOwner = @owner 
					   AND nStorageType = @storageType 
					   AND nStorage + 1 NOT IN (
							SELECT nStorage
							  FROM tItem 
							 WHERE nOwner = @owner AND nStorageType = @storageType)
				END
			ELSE
				BEGIN
					SET @emptyStorage = 0
				END

			SET @ret = @emptyStorage
			PRINT @ret
		END
	ELSE
		SET @ret = -1

	SET NOCOUNT OFF
END

0Ёꨅ^
GO

-- ========================================================================
-- Procedure: p_Item_GetStorageEmptyList
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Item_GetStorageEmptyList]
	@nOwner         int
,	@nStorageType   int
,	@nGetStorageCnt int = 1
,	@nRet           int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	IF @nOwner = -1 BEGIN
		SET @nRet = -1000
		RETURN
	END
	
	-- 동적 쿼리 작성 시작
	DECLARE @SQL AS NVARCHAR(3000)
	DECLARE @PARAMS AS NVARCHAR(1000)
	SET @SQL = N''
	SET @PARAMS = N''


	SET @SQL = '
		SELECT TOP ' + CAST(@nGetStorageCnt AS NVARCHAR) + ' A.nStorage
		FROM tItemStorageFindEmpty A LEFT OUTER JOIN (
			SELECT nStorage
			FROM tItem
			WHERE nOwner = ' + CAST(@nOwner AS NVARCHAR) + '
			AND nStorageType = ' + CAST(@nStorageType AS NVARCHAR) + ') B ON A.nStorage = B.nStorage
		WHERE B.nStorage IS NULL
		ORDER BY A.nStorage'
	SET @PARAMS = '@nGetStorageCnt INT, @nOwner INT, @nStorageType INT'

	-- 만들어진 동적 쿼리 실행
	EXEC sp_executesql @SQL, @PARAMS, @nGetStorageCnt, @nOwner, @nStorageType

	SET NOCOUNT OFF
END

0㴁鸩_
GO

-- ========================================================================
-- Procedure: p_Item_GetStorageListType
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Item_GetStorageListType]

/*
Item Get storage list from StorageType
2004.8 By CJC
 input:
	nCharNo	     = Item Owner
	nStorageType = Item Storage Type
 output:
	none
 desc:
	fetch success, get data
	fetch fail, not data
*/

@nCharNo int,
@nStorageType tinyint

AS

SET NOCOUNT ON

SELECT nStorage FROM tItem WHERE nOwner = @nCharNo AND nStorageType = @nStorageType

-- end

0༁ꝜŠ
GO

-- ========================================================================
-- Procedure: p_Item_GetStorageTypeCount
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Item_GetStorageTypeCount]
/*
Item Get number of item count from StorageType
2008.4 By CJC
 input:
	Item Owner
	Item StorageType
 output:
	Number of item count from StorageType
 */

@nOwner int,
@nStorageType tinyint,
@nItemCount int OUTPUT
AS
BEGIN
	SET NOCOUNT ON
	SELECT @nItemCount=COUNT(nOwner) FROM tItem WHERE nOwner = @nOwner AND nStorageType = @nStorageType
END

0
GO

-- ========================================================================
-- Procedure: p_Item_GetTitle
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Item_GetTitle]
/*
Item Get Title
2004.9 By CJC
 input:
	see SQL
 output:
	@nRet	0 = No data fount
	        ? = Number data here(must to be 1)
	        
	@nTitleType,@nTitleCharNo,@sTitle = title data
		        
*/
@nItemKey bigint,
-- Output var
@nTitleType smallint OUTPUT,
@nTitleCharNo int OUTPUT,
@sTitle nvarchar(24) OUTPUT,
@nRet int OUTPUT
AS
SET NOCOUNT ON
SELECT 
@nTitleType = nTitleType,
@nTitleCharNo = nTitleCharNo,
@sTitle = sTitle
FROM tItemTitle WHERE nItemKey = @nItemKey
SET @nRet = @@ROWCOUNT
-- end

0ⴁ▐Ţ
GO

-- ========================================================================
-- Procedure: p_Item_IsValidOwner
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Item_IsValidOwner]
/*
Is valid item owner ?
2004.8 By CJC
 input:
	nItemKey	 = Item Key
	nOwner	     = Item Owner
 output:
	0 = OK
	1 = Not found item
	2 = Missmatch item owner
*/
@nItemKey bigint,
@nOwner int,
@nRet tinyint OUTPUT 
AS
SET NOCOUNT ON
DECLARE @nDBOwner int 
SET @nDBOwner = 0
SET @nRet = 0
SELECT @nDBOwner = nOwner FROM tItem WITH (READUNCOMMITTED) WHERE nItemKey = @nItemKey
-- Check not found item
IF @@ROWCOUNT <> 1
BEGIN
	SET @nRet = 1
	RETURN
END
-- Check missmatch item owner
IF @nDBOwner <> @nOwner
BEGIN
	SET @nRet = 2
	RETURN
END
-- end

0昁ᦴţ
GO

-- ========================================================================
-- Procedure: p_Item_Merge
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Item_Merge]
/*
Item Merge - change lot and delete item
2004.9 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK
	        1 = Not found from item
	        2 = Missmatch from item owner
	        3 = Not found to item
	        4 = Missmatch to item owner
	        ? = Error code(Lot change/Delete item error)
*/
@nOwnerFrom int,
@nOwnerTo int,
@nItemKeyFrom bigint,
@nItemKeyTo bigint,
@nOptionType smallint,
@nOptionData bigint,
@nRet int OUTPUT
AS
SET NOCOUNT ON
-- 삭제될 아이템(From) 소유자 확인
DECLARE @nTinyintRet tinyint
EXEC p_Item_IsValidOwner @nItemKeyFrom, @nOwnerFrom, @nTinyintRet OUTPUT
IF @nTinyintRet <> 0
BEGIN
	SET @nRet = @nTinyintRet
	RETURN
END
-- 랏 변경될 아이템(To) 소유자 확인
EXEC p_Item_IsValidOwner @nItemKeyTo, @nOwnerTo, @nTinyintRet OUTPUT
IF @nTinyintRet <> 0
BEGIN
	SET @nRet = @nTinyintRet + 2
	RETURN
END
BEGIN TRAN
-- 아이템 랏 변경저장프로시져 호출
EXEC p_Item_SetOption @nItemKeyTo, @nOptionType, @nOptionData, @nRet OUTPUT
IF @nRet <> 0
BEGIN
	ROLLBACK TRAN
	RETURN
END
-- 아이템 삭제
EXEC p_Item_Delete @nItemKeyFrom, @nRet OUTPUT
IF @@ERROR <> 0
BEGIN
	ROLLBACK TRAN
	RETURN
END
-- 완료
COMMIT TRAN
-- end

0ȁ㜇ť
GO

-- ========================================================================
-- Procedure: p_Item_Merge_New
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Item_Merge_New]
/*
Item Merge - change lot and delete item
2007.11 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK
	        1 = Not found from item
	        2 = Missmatch from item owner
	        3 = Not found to item
	        4 = Missmatch to item owner
                     5 = Get Sub Lot Failed
                     6 = Get Existed Lot Failed
                     7 = Set new lot failed
                     8 = Delete item failed
*/
@nOwnerFrom int,
@nOwnerTo int,
@nItemKeyFrom bigint,
@nItemKeyTo bigint,
@nOptionType smallint,
@nOptionDataLotSub bigint OUTPUT,
@nOptionDataLotResult bigint OUTPUT,
@nRet int OUTPUT
AS
SET NOCOUNT ON

-- »èÁ¦µÉ ¾ÆÀÌÅÛ(From) ¼ÒÀ¯ÀÚ È®ÀÎ
DECLARE @nTinyintRet tinyint
EXEC p_Item_IsValidOwner @nItemKeyFrom, @nOwnerFrom, @nTinyintRet OUTPUT
IF @nTinyintRet <> 0
BEGIN
	SET @nRet = @nTinyintRet
	RETURN
END

-- ¶ù º¯°æµÉ ¾ÆÀÌÅÛ(To) ¼ÒÀ¯ÀÚ È®ÀÎ
EXEC p_Item_IsValidOwner @nItemKeyTo, @nOwnerTo, @nTinyintRet OUTPUT
IF @nTinyintRet <> 0
BEGIN
	SET @nRet = @nTinyintRet + 2
	RETURN
END

-- »©Áú ¾ÆÀÌÅÛ ¶ù ¾ò±â
EXEC p_Item_GetOption @nItemKeyFrom, @nOptionType, @nOptionDataLotSub OUTPUT, @nRet OUTPUT
IF @nRet <> 1
BEGIN
	SET @nRet = 5
	RETURN
END

-- ´õÇØÁú ¾ÆÀÌÅÛ ¶ù ¾ò±â
DECLARE @nOptionDataLotAdd bigint
EXEC p_Item_GetOption @nItemKeyTo, @nOptionType, @nOptionDataLotAdd OUTPUT, @nRet OUTPUT
IF @nRet <> 1
BEGIN
	SET @nRet = 6
	RETURN
END

-- ¾ÆÀÌÅÛ ¶ù º¯°æÀúÀåÇÁ·Î½ÃÁ® È£Ãâ
BEGIN TRAN
SET @nOptionDataLotResult = @nOptionDataLotSub + @nOptionDataLotAdd
EXEC p_Item_SetOption @nItemKeyTo, @nOptionType, @nOptionDataLotResult, @nRet OUTPUT
IF @nRet <> 0
BEGIN
	ROLLBACK TRAN
	SET @nRet = 7
	RETURN
END

-- ¾ÆÀÌÅÛ »èÁ¦
EXEC p_Item_Delete @nItemKeyFrom, @nRet OUTPUT
IF @@ERROR <> 0
BEGIN
	ROLLBACK TRAN
	SET @nRet = 8
	RETURN
END

-- ¿Ï·á
COMMIT TRAN

-- end

0
GO

-- ========================================================================
-- Procedure: p_Item_SellAll
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Item_SellAll]
/*
Item sell all
2004.9 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK
	        1 = Not found item
	        2 = Missmatch item owner
	        3 = Money add error
	        ? = Error code
	        
 operation:
	1. money add
	2. delete sell item
	
*/
@nItemKey bigint,
@nOwner int,
@nPrice bigint,
@nRet int OUTPUT
AS
SET NOCOUNT ON
-- 아이템 소유자 확인
DECLARE @nTinyintRet tinyint
EXEC p_Item_IsValidOwner @nItemKey, @nOwner, @nTinyintRet OUTPUT
IF @nTinyintRet <> 0
BEGIN
	SET @nRet = @nTinyintRet
	RETURN
END
-- 아이템 소유자의 돈 증가
BEGIN TRAN
EXEC p_Char_MoneyAdd @nOwner, @nPrice, @nTinyintRet OUTPUT
IF @nTinyintRet <> 1
BEGIN
	ROLLBACK TRAN
	SET @nRet = 3
	RETURN
END
-- 아이템 삭제 저장프로시져 호출
EXEC p_Item_Delete @nItemKey, @nRet OUTPUT
IF @@ERROR <> 0
BEGIN
	ROLLBACK TRAN
	RETURN
END
-- 완료
COMMIT TRAN
-- end

0攁怵Ŧ
GO

-- ========================================================================
-- Procedure: p_Item_SellAll_LotCheck
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Item_SellAll_LotCheck]
/*
Item sell all
2004.9 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK
	        1 = Not found item
	        2 = Missmatch item owner
	        3 = Money add error
	        ? = Error code
	        
 operation:
	1. money add
	2. delete sell item
	
*/
@nItemKey bigint,
@nOwner int,
@nOptionType smallint,
@nLot bigint,
@nPrice bigint,
@nRet int OUTPUT
AS
SET NOCOUNT ON
-- 아이템 소유자 확인
DECLARE @nTinyintRet tinyint
EXEC p_Item_IsValidOwner @nItemKey, @nOwner, @nTinyintRet OUTPUT
IF @nTinyintRet <> 0
BEGIN
	SET @nRet = @nTinyintRet
	RETURN
END

-- 1보다 큰 경우에만 Lot 을 체크해 본다.
-- 1일 경우는 어차피 tItem 테이블에 존재하는 것만으로도 옳은 것으로 취급한다.
IF @nLot > 1
BEGIN
	DECLARE @nGetOptionDataLot bigint
	DECLARE @nGetOptionRet int
	EXEC p_Item_GetOption @nItemKey, @nOptionType, @nGetOptionDataLot OUTPUT, @nGetOptionRet OUTPUT
	IF @nGetOptionRet <> 1
	BEGIN
		SET @nRet = 4
		RETURN
	END

	IF @nGetOptionDataLot < @nLot
	BEGIN
		SET @nRet = 5
		RETURN
	END
END

-- 아이템 소유자의 돈 증가
BEGIN TRAN
EXEC p_Char_MoneyAdd @nOwner, @nPrice, @nTinyintRet OUTPUT
IF @nTinyintRet <> 1
BEGIN
	ROLLBACK TRAN
	SET @nRet = 3
	RETURN
END
-- 아이템 삭제 저장프로시져 호출
EXEC p_Item_Delete @nItemKey, @nRet OUTPUT
IF @@ERROR <> 0
BEGIN
	ROLLBACK TRAN
	RETURN
END
-- 완료
COMMIT TRAN
-- end

0䨁Ŧ
GO

-- ========================================================================
-- Procedure: p_Item_SellLot
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Item_SellLot]
/*
Item sell - change lot and add money
2004.9 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK
	        1 = Not found item
	        2 = Missmatch item owner
	        3 = Money add error
	        ? = Error code(Lot change error)
*/
@nItemKey bigint,
@nOwner int,
@nOptionType smallint,
@nOptionData bigint,
@nPrice bigint,
@nRet int OUTPUT
AS
SET NOCOUNT ON
-- 아이템 소유자 확인
DECLARE @nTinyintRet tinyint
EXEC p_Item_IsValidOwner @nItemKey, @nOwner, @nTinyintRet OUTPUT
IF @nTinyintRet <> 0
BEGIN
	SET @nRet = @nTinyintRet
	RETURN
END
BEGIN TRAN
-- 아이템 소유자가 돈 증가
EXEC p_Char_MoneyAdd @nOwner, @nPrice, @nRet OUTPUT
IF @nRet <> 1
BEGIN
	ROLLBACK TRAN
	SET @nRet = 3
	RETURN
END
-- 아이템 랏 저장프로시져 호출
EXEC p_Item_SetOption @nItemKey, @nOptionType, @nOptionData, @nRet OUTPUT
IF @nRet <> 0
BEGIN
	ROLLBACK TRAN
	RETURN
END
-- 완료
COMMIT TRAN
-- end

0茁ŧ
GO

-- ========================================================================
-- Procedure: p_Item_SellLot_LotCheck
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Item_SellLot_LotCheck]
/*
Item sell - change lot and add money
2004.9 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK
	        1 = Not found item
	        2 = Missmatch item owner
	        3 = Money add error
	        ? = Error code(Lot change error)
*/
@nItemKey bigint,
@nOwner int,
@nOptionType smallint,
@nOptionData bigint,
@nPrice bigint,
@nRet int OUTPUT
AS
SET NOCOUNT ON
-- 아이템 소유자 확인
DECLARE @nTinyintRet tinyint
EXEC p_Item_IsValidOwner @nItemKey, @nOwner, @nTinyintRet OUTPUT
IF @nTinyintRet <> 0
BEGIN
	SET @nRet = @nTinyintRet
	RETURN
END

DECLARE @nGetOptionDataLot bigint
DECLARE @nGetOptionRet int
EXEC p_Item_GetOption @nItemKey, @nOptionType, @nGetOptionDataLot OUTPUT, @nGetOptionRet OUTPUT
IF @nGetOptionRet <> 1
BEGIN
	SET @nRet = 4
	RETURN
END

IF @nGetOptionDataLot < @nOptionData
BEGIN
	SET @nRet = 5
	RETURN
END

BEGIN TRAN
-- 아이템 소유자가 돈 증가
EXEC p_Char_MoneyAdd @nOwner, @nPrice, @nRet OUTPUT
IF @nRet <> 1
BEGIN
	ROLLBACK TRAN
	SET @nRet = 3
	RETURN
END
-- 아이템 랏 저장프로시져 호출
EXEC p_Item_SetOption @nItemKey, @nOptionType, @nOptionData, @nRet OUTPUT
IF @nRet <> 0
BEGIN
	ROLLBACK TRAN
	RETURN
END
-- 완료
COMMIT TRAN
-- end

0ခ㲢ũ
GO

-- ========================================================================
-- Procedure: p_Item_SetData
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Item_SetData]
/*
Item Set data
2004.10 By CJC
 input:
	see SQL
 output:
	@nRet	1 = OK, 
	        0 = Error
*/
@nItemKey bigint,
@nOwner int,
@nStorageType tinyint,
@nStorage smallint,
@nItemID int,
@nFlags int,
@dDate datetime,
@nRet int OUTPUT
AS
SET NOCOUNT ON
SET @nRet = 0
UPDATE tItem SET 
nOwner = @nOwner,
nStorageType = @nStorageType,
nStorage = @nStorage,
nItemID = @nItemID,
nFlags = @nFlags,
dDate = @dDate
WHERE nItemKey = @nItemKey
IF @@ERROR = 0 AND @@ROWCOUNT = 1
	SET @nRet = 1
-- end

0省➓i
GO

-- ========================================================================
-- Procedure: p_Item_SetMobList
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Item_SetMobList]
/*
Item Set Mob List
2004.8 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK, 
	        ? = Error code
*/
@nItemKey bigint,
@nMobID smallint,
@nGrade tinyint,
@nKillCount int,
-- Output var
@nRet int OUTPUT
AS
SET NOCOUNT ON
UPDATE tItemMobList SET	nGrade = @nGrade, nKillCount = @nKillCount
WHERE nItemKey = @nItemKey AND nMobID = @nMobID
IF @@ERROR <> 0 OR @@ROWCOUNT = 0
BEGIN
	INSERT tItemMobList (  nItemKey,  nMobID,  nGrade,  nKillCount )
	VALUES				( @nItemKey, @nMobID, @nGrade, @nKillCount )
END
SET @nRet = @@ERROR
-- end

0舁⓪ū
GO

-- ========================================================================
-- Procedure: p_Item_SetOption
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Item_SetOption]
/*
Item Set Option
2004.8 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK, 
	        ? = Error code
*/
@nItemKey bigint,
@nOptionType smallint,
@nOptionData bigint,
-- Output var
@nRet int OUTPUT
AS
SET NOCOUNT ON
UPDATE tItemOptions SET	nOptionData = @nOptionData
WHERE nItemKey = @nItemKey AND nOptionType = @nOptionType
IF @@ERROR <> 0 OR @@ROWCOUNT = 0
BEGIN
	INSERT tItemOptions (  nItemKey,  nOptionType,   nOptionData )
	VALUES				( @nItemKey, @nOptionType,  @nOptionData )
END
SET @nRet = @@ERROR
-- end

RETURN

0朁껺ū
GO

-- ========================================================================
-- Procedure: p_Item_SetStorage
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Item_SetStorage]
/*
Get item storage
2004.8 By CJC
 input:
	nItemKey	 = Item Key
	nStorageType, nStorage = Storage new value
 output:
	1 = OK
	0 = Not found or error
*/
@nItemKey bigint,
@nStorageType tinyint,
@nStorage smallint,
@nRet tinyint OUTPUT
AS
SET NOCOUNT ON
UPDATE tItem 
SET	nStorageType = @nStorageType, nStorage = @nStorage
WHERE nItemKey = @nItemKey
SET @nRet = @@ROWCOUNT
-- end

0ꀁꌞŬ
GO

-- ========================================================================
-- Procedure: p_Item_SetStorageAndOwner
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Item_SetStorageAndOwner]
/*
Get item storage
2004.8 By CJC
 input:
	nItemKey	 = Item Key
	nOwner, nStorageType, nStorage = Storage new value
 output:
	1 = OK
	0 = Not found or error
*/
@nItemKey bigint,
@nOwner int,
@nStorageType tinyint,
@nStorage smallint,
@nRet tinyint OUTPUT
AS
SET NOCOUNT ON
UPDATE tItem 
SET	nOwner = @nOwner, nStorageType = @nStorageType, nStorage = @nStorage
WHERE nItemKey = @nItemKey
SET @nRet = @@ROWCOUNT
-- end

0ⴁŗŮ
GO

-- ========================================================================
-- Procedure: p_Item_SetStorageAndOwner_UseEmptyStorage
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Item_SetStorageAndOwner_UseEmptyStorage]
/*
Get item storage
2004.8 By CJC
 input:
	nItemKey	 = Item Key
	nOwner, nStorageType, nStorage = Storage new value
 output:
	1 = OK
	0 = Not found or error
*/
@nItemKey bigint,
@nOwner int,
@nStorageType tinyint,
@nStorage smallint OUTPUT,
@nRet tinyint OUTPUT
AS
SET NOCOUNT ON


	-- 빈인벤 위치찾기
	EXEC  p_Item_GetStorageEmpty @nOwner, @nStorageType, @nStorage OUTPUT
	IF @nStorage < 0
	BEGIN
		SET @nRet = -1
		ROLLBACK TRAN
		RETURN
	END

	UPDATE tItem 
	SET	nOwner = @nOwner, nStorageType = @nStorageType, nStorage = @nStorage
	WHERE nItemKey = @nItemKey

	SET @nRet = @@ROWCOUNT
-- end

0昁Ů
GO

-- ========================================================================
-- Procedure: p_Item_SetTitle
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Item_SetTitle]
/*
Item Set Title
2004.8 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK, 
	        ? = Error code
*/
@nItemKey bigint,
@nTitleType smallint,
@nTitleCharNo int,
@sTitle nvarchar(24),
-- Output var
@nRet int OUTPUT
AS
SET NOCOUNT ON
UPDATE tItemTitle SET 
nTitleType = @nTitleType,
nTitleCharNo = @nTitleCharNo,
sTitle = @sTitle
WHERE nItemKey = @nItemKey
IF @@ERROR <> 0 OR @@ROWCOUNT = 0
BEGIN
	INSERT tItemTitle (  nItemKey,  nTitleType,  nTitleCharNo,  sTitle )
	VALUES			  ( @nItemKey, @nTitleType, @nTitleCharNo, @sTitle )
END
SET @nRet = @@ERROR
-- end

0䬁羋ů
GO

-- ========================================================================
-- Procedure: p_Item_SetWeaponTitle
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Item_SetWeaponTitle]
/*

Item Set Weapon Title
2005.8 By CJC

 input:
	see SQL

 output:
	@nRet	0 = OK, 
	        1 = Error Set Weapon mob 
	        2 = Error Set Weapon user
	        3 = Error Delete license item
	        ? = Error DB
*/

@nWeaponItemKey bigint,
@nLicenseItemKey bigint,
@nMobID smallint,

@nTitleCharNo int,
@sTitle nvarchar(24),

-- Output var
@nRet int OUTPUT

AS

SET NOCOUNT ON

BEGIN TRAN
-- 아이템에 몹 라이센스 붙이기
EXEC p_Item_SetMobList @nWeaponItemKey, @nMobID, 0, 0, @nRet OUTPUT
IF @nRet <> 0
BEGIN
	ROLLBACK TRAN
	SET @nRet = 1
	RETURN
END

-- 처음 붙은 아이템의 경우 최초 제작자 정보를 세팅한다.
DECLARE @nRetTemp int
SELECT @nRetTemp = COUNT(*) FROM tItemMobList WHERE nItemKey = @nWeaponItemKey
IF @nRetTemp = 1
BEGIN
	EXEC p_Item_SetTitle @nWeaponItemKey, 0, @nTitleCharNo, @sTitle, @nRet OUTPUT
	IF @nRet <> 0
	BEGIN
		ROLLBACK TRAN
		SET @nRet = 2
		RETURN
	END
END

-- 사용된 라이센스 아이템은 지운다.
SET @nRetTemp = 1
EXEC p_Item_Delete @nLicenseItemKey, @nRetTemp OUTPUT
IF @nRetTemp <> 0
BEGIN
	ROLLBACK TRAN
	SET @nRet = 3
	RETURN
END

-- 성공
COMMIT TRAN
SET @nRet = 0

RETURN

-- end

0萁环Ű
GO

-- ========================================================================
-- Procedure: p_Item_Split_N_Merge
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Item_Split_N_Merge]
/*
Item Split and Merge - change lot
2004.9 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK
	        1 = Not found from item
	        2 = Missmatch from item owner
	        3 = Not found to item
	        4 = Missmatch to item owner
	        ? = Error code(Lot change/Delete item error)
*/
@nItemKeyFrom bigint,
@nOwnerFrom int,
@nOptionTypeFrom smallint,
@nOptionDataFrom bigint,
@nItemKeyTo bigint,
@nOwnerTo int,
@nOptionTypeTo smallint,
@nOptionDataTo bigint,
@nRet int OUTPUT
AS
SET NOCOUNT ON
-- 수량 줄어들 아이템(From) 소유자 확인
DECLARE @nTinyintRet tinyint
EXEC p_Item_IsValidOwner @nItemKeyFrom, @nOwnerFrom, @nTinyintRet OUTPUT
IF @nTinyintRet <> 0
BEGIN
	SET @nRet = @nTinyintRet
	RETURN
END
-- 수량 늘어날 아이템(To) 소유자 확인
EXEC p_Item_IsValidOwner @nItemKeyTo, @nOwnerTo, @nTinyintRet OUTPUT
IF @nTinyintRet <> 0
BEGIN
	SET @nRet = @nTinyintRet + 2
	RETURN
END
BEGIN TRAN
-- 수량 줄어들 아이템 랏 변경저장프로시져 호출
EXEC p_Item_SetOption @nItemKeyFrom, @nOptionTypeFrom, @nOptionDataFrom, @nRet OUTPUT
IF @nRet <> 0
BEGIN
	ROLLBACK TRAN
	RETURN
END
-- 수량 늘어날 아이템 랏 변경저장프로시져 호출
EXEC p_Item_SetOption @nItemKeyTo, @nOptionTypeTo, @nOptionDataTo, @nRet OUTPUT
IF @nRet <> 0
BEGIN
	ROLLBACK TRAN
	RETURN
END
-- 완료
COMMIT TRAN
-- end

0봁染ű
GO

-- ========================================================================
-- Procedure: p_Item_Split_N_Merge_New
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Item_Split_N_Merge_New]
/*
Item Split and Merge - change lot
2004.9 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK
	        1 = Not found from item
	        2 = Missmatch from item owner
	        3 = Not found to item
	        4 = Missmatch to item owner
                     5 = From lot read failed
                     6 = From lot less more nSplitLot
                     7 = To lot read failed
                     8 = From lot set failed
                     9 = To lot set failed
*/

@nItemKeyFrom bigint,
@nOwnerFrom int,
@nOptionTypeFrom smallint,
@nOptionDataFrom bigint OUTPUT,

@nItemKeyTo bigint,
@nOwnerTo int,
@nOptionTypeTo smallint,
@nOptionDataTo bigint OUTPUT,

@nSplitLot bigint,
@nRet int OUTPUT

AS

SET NOCOUNT ON

-- ¼ö·® ÁÙ¾îµé ¾ÆÀÌÅÛ(From) ¼ÒÀ¯ÀÚ È®ÀÎ
DECLARE @nTinyintRet tinyint
EXEC p_Item_IsValidOwner @nItemKeyFrom, @nOwnerFrom, @nTinyintRet OUTPUT
IF @nTinyintRet <> 0
BEGIN
	SET @nRet = @nTinyintRet
	RETURN
END

-- ¼ö·® ´Ã¾î³¯ ¾ÆÀÌÅÛ(To) ¼ÒÀ¯ÀÚ È®ÀÎ
EXEC p_Item_IsValidOwner @nItemKeyTo, @nOwnerTo, @nTinyintRet OUTPUT
IF @nTinyintRet <> 0
BEGIN
	SET @nRet = @nTinyintRet + 2
	RETURN
END

-- »©Áú ¾ÆÀÌÅÛ ¶ù ¾ò±â
EXEC p_Item_GetOption @nItemKeyFrom, @nOptionTypeFrom, @nOptionDataFrom OUTPUT, @nRet OUTPUT
IF @nRet <> 1
BEGIN
	SET @nRet = 5
	RETURN
END

-- »©Áú ¾ÆÀÌÅÛ ¶ùÀÌ ÃæºÐÇÑÁö È®ÀÎ
IF @nOptionDataFrom < @nSplitLot
BEGIN
	SET @nRet = 6
	RETURN
END

-- ´õÇØÁú ¾ÆÀÌÅÛ ¶ù ¾ò±â
DECLARE @nOptionDataLotAdd bigint
EXEC p_Item_GetOption @nItemKeyTo, @nOptionTypeTo, @nOptionDataTo OUTPUT, @nRet OUTPUT
IF @nRet <> 1
BEGIN
	SET @nRet = 7
	RETURN
END

-- °á°ú ¶ù °è»ê
SET @nOptionDataFrom = @nOptionDataFrom - @nSplitLot
SET @nOptionDataTo = @nOptionDataTo + @nSplitLot

BEGIN TRAN
-- ¼ö·® ÁÙ¾îµé ¾ÆÀÌÅÛ ¶ù º¯°æÀúÀåÇÁ·Î½ÃÁ® È£Ãâ
EXEC p_Item_SetOption @nItemKeyFrom, @nOptionTypeFrom, @nOptionDataFrom, @nRet OUTPUT
IF @nRet <> 0
BEGIN
	ROLLBACK TRAN
	SET @nRet = 8
	RETURN
END

-- ¼ö·® ´Ã¾î³¯ ¾ÆÀÌÅÛ ¶ù º¯°æÀúÀåÇÁ·Î½ÃÁ® È£Ãâ
EXEC p_Item_SetOption @nItemKeyTo, @nOptionTypeTo, @nOptionDataTo, @nRet OUTPUT
IF @nRet <> 0
BEGIN
	ROLLBACK TRAN
	SET @nRet = 9
	RETURN
END

-- ¿Ï·á
COMMIT TRAN
-- end

0
GO

-- ========================================================================
-- Procedure: p_Item_Upgrade
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Item_Upgrade]
/*
Item upgrade
2005.6 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK
	        1 = Not found item
	        2 = Missmatch item owner
	        3 = Not enough money to sub
	        4 = Upgrade value setting, upgrade fail count error(option set error)
	        5 = Raw delete fail
	        6 = RawLeft delete fail, if valid item
	        7 = RawRight delete fail, if valid item
	        8 = RawMiddle delete fail, if valid item
	        
 operation:
	1. money sub
	2. upgrade item level
	3. delete raws
*/
@nOwner int,
@nItemKey bigint,
@nItemDestroy int,
@nUpgradeOptionType smallint,
@nUpgradeOptionValue bigint,
@nUpgradeFailCountType smallint,
@nUpgradeFailCountValue bigint,
@nRawKey bigint,
@nRawLeftKey bigint,
@nRawRightKey bigint,
@nRawMiddleKey bigint,
@nGiftMoney bigint,
@nRet int OUTPUT
AS
SET NOCOUNT ON
-- 아이템 소유자 확인
DECLARE @nTinyintRet tinyint
EXEC p_Item_IsValidOwner @nItemKey, @nOwner, @nTinyintRet OUTPUT
IF @nTinyintRet <> 0
BEGIN
	SET @nRet = @nTinyintRet
	RETURN
END
-- 아이템 소유자의 돈 감소
BEGIN TRAN
IF @nGiftMoney > 0 -- 수고비가 있는 경우, 금액을 빼준다.
BEGIN
	EXEC p_Char_MoneySub @nOwner, @nGiftMoney, @nTinyintRet OUTPUT
	IF @nTinyintRet <> 1
	BEGIN
		ROLLBACK TRAN
		SET @nRet = 3
		RETURN
	END
END	
-- 아이템 업그레이드 값 조정
DECLARE @nRetTemp int
IF @nItemDestroy = 1 -- 업글 실패로 아이템 파괴해야 되는경우
BEGIN
	SET @nRetTemp = 1
	EXEC p_Item_Delete @nItemKey, @nRetTemp OUTPUT
	IF @nRetTemp <> 0
	BEGIN
		ROLLBACK TRAN
		SET @nRet = 4
		RETURN
	END
END
ELSE
BEGIN
	-- 업글 실패로 파괴가 아닌 경우, 새로운 업글값, 실패 카운트 기록	
	SET @nRetTemp = 1
	EXEC p_Item_SetOption @nItemKey, @nUpgradeOptionType, @nUpgradeOptionValue, @nRetTemp OUTPUT
	IF @nRetTemp <> 0
	BEGIN
		ROLLBACK TRAN
		SET @nRet = 4
		RETURN
	END
	
	SET @nRetTemp = 1
	EXEC p_Item_SetOption @nItemKey, @nUpgradeFailCountType, @nUpgradeFailCountValue, @nRetTemp OUTPUT
	IF @nRetTemp <> 0
	BEGIN
		ROLLBACK TRAN
		SET @nRet = 4
		RETURN
	END
END
-- 사용된 자원 지우기
SET @nRetTemp = 1
EXEC p_Item_Delete @nRawKey, @nRetTemp OUTPUT
IF @nRetTemp <> 0
BEGIN
	ROLLBACK TRAN
	SET @nRet = 5
	RETURN
END
-- 사용된 자원왼족 지우기(입력된 값이 있으면)
IF @nRawLeftKey <> 0
BEGIN
	SET @nRetTemp = 1
	EXEC p_Item_Delete @nRawLeftKey, @nRetTemp OUTPUT
	IF @nRetTemp <> 0
	BEGIN
		ROLLBACK TRAN
		SET @nRet = 6
		RETURN
	END
END
-- 사용된 자원오른쪽 지우기(입력된 값이 있으면)
IF @nRawRightKey <> 0
BEGIN
	SET @nRetTemp = 1
	EXEC p_Item_Delete @nRawRightKey, @nRetTemp OUTPUT
	IF @nRetTemp <> 0
	BEGIN
		ROLLBACK TRAN
		SET @nRet = 7
		RETURN
	END
END
-- 사용된 자원가운데 지우기(입력된 값이 있으면)
IF @nRawMiddleKey <> 0
BEGIN
	SET @nRetTemp = 1
	EXEC p_Item_Delete @nRawMiddleKey, @nRetTemp OUTPUT
	IF @nRetTemp <> 0
	BEGIN
		ROLLBACK TRAN
		SET @nRet = 8
		RETURN
	END
END

-- 업글 성공
COMMIT TRAN
SET @nRet = 0
-- end

0ꔡs
GO

-- ========================================================================
-- Procedure: p_Item_Upgrade_CelllInAll
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Item_Upgrade_CelllInAll]
/*
Item upgrade
2008.2 By Lunar
 input:
	see SQL
 output:
	@nRet	0 = OK
	        1 = Not found item
	        2 = Missmatch item owner
	        3 = Not enough money to sub
	        4 = Upgrade value setting, upgrade fail count error(option set error)
	        5 = Raw delete fail
	        6 = RawLeft delete fail, if valid item
	        7 = RawRight delete fail, if valid item
	        8 = RawMiddle delete fail, if valid item
	        
 operation:
	1. money sub
	2. upgrade item level
	3. delete raws
*/
@nOwner int,
@nItemKey bigint,
@nItemDestroy int,
@nUpgradeOptionType smallint,
@nUpgradeOptionValue bigint,
@nUpgradeFailCountType smallint,
@nUpgradeFailCountValue bigint,
@nRawKey bigint,
@nRawLeftKey bigint,
@nRawRightKey bigint,
@nRawMiddleKey bigint,
@nGiftMoney bigint,
@nRawOptionType smallint,
@nRawLot smallint,
@nRawLeftOptionType smallint,
@nRawLeftLot smallint,
@nRawRightOptionType smallint,
@nRawRightLot smallint,
@nRawMiddleOptionType smallint,
@nRawMiddleLot smallint,
@nRet int OUTPUT
AS
SET NOCOUNT ON
-- 아이템 소유자 확인
DECLARE @nTinyintRet tinyint
EXEC p_Item_IsValidOwner @nItemKey, @nOwner, @nTinyintRet OUTPUT
IF @nTinyintRet <> 0
BEGIN
	SET @nRet = @nTinyintRet
	RETURN
END
-- 아이템 소유자의 돈 감소
BEGIN TRAN
IF @nGiftMoney > 0 -- 수고비가 있는 경우, 금액을 빼준다.
BEGIN
	EXEC p_Char_MoneySub @nOwner, @nGiftMoney, @nTinyintRet OUTPUT
	IF @nTinyintRet <> 1
	BEGIN
		ROLLBACK TRAN
		SET @nRet = 3
		RETURN
	END
END	
-- 아이템 업그레이드 값 조정
DECLARE @nRetTemp int
IF @nItemDestroy = 1 -- 업글 실패로 아이템 파괴해야 되는경우
BEGIN
	SET @nRetTemp = 1
	EXEC p_Item_Delete @nItemKey, @nRetTemp OUTPUT
	IF @nRetTemp <> 0
	BEGIN
		ROLLBACK TRAN
		SET @nRet = 4
		RETURN
	END
END
ELSE
BEGIN
	-- 업글 실패로 파괴가 아닌 경우, 새로운 업글값, 실패 카운트 기록	
	SET @nRetTemp = 1
	EXEC p_Item_SetOption @nItemKey, @nUpgradeOptionType, @nUpgradeOptionValue, @nRetTemp OUTPUT
	IF @nRetTemp <> 0
	BEGIN
		ROLLBACK TRAN
		SET @nRet = 4
		RETURN
	END
	
	SET @nRetTemp = 1
	EXEC p_Item_SetOption @nItemKey, @nUpgradeFailCountType, @nUpgradeFailCountValue, @nRetTemp OUTPUT
	IF @nRetTemp <> 0
	BEGIN
		ROLLBACK TRAN
		SET @nRet = 4
		RETURN
	END
END

-- 사용된 자원 지우기
IF @nRawKey <> 0
BEGIN
	IF @nRawLot = 1
	BEGIN
		SET @nRetTemp = 1
		EXEC p_Item_Delete @nRawKey, @nRetTemp OUTPUT
		IF @nRetTemp <> 0
		BEGIN
			ROLLBACK TRAN
			SET @nRet = 5
			RETURN
		END
	END
	IF @nRawLot > 1
	BEGIN
		SET @nRetTemp = 1
		SET @nRawLot = @nRawLot - 1
		EXEC p_Item_SetOption @nRawKey, @nRawOptionType, @nRawLot, @nRetTemp OUTPUT
		IF @nRetTemp <> 0
		BEGIN
			ROLLBACK TRAN
			SET @nRet = 5
			RETURN
		END
	END
END

-- 사용된 자원왼족 지우기(입력된 값이 있으면)
IF @nRawLeftKey <> 0
BEGIN
	IF @nRawLeftLot = 1
	BEGIN
		SET @nRetTemp = 1
		EXEC p_Item_Delete @nRawLeftKey, @nRetTemp OUTPUT
		IF @nRetTemp <> 0
		BEGIN
			ROLLBACK TRAN
			SET @nRet = 6
			RETURN
		END
	END
	IF @nRawLeftLot > 1
	BEGIN
		SET @nRetTemp = 1
		SET @nRawLeftLot = @nRawLeftLot - 1
		EXEC p_Item_SetOption @nRawLeftKey, @nRawLeftOptionType, @nRawLeftLot, @nRetTemp OUTPUT
		IF @nRetTemp <> 0
		BEGIN
			ROLLBACK TRAN
			SET @nRet = 6
			RETURN
		END
	END
END


-- 사용된 자원오른쪽 지우기(입력된 값이 있으면)
IF @nRawRightKey <> 0
BEGIN
	IF @nRawRightLot = 1
	BEGIN
		SET @nRetTemp = 1
		EXEC p_Item_Delete @nRawRightKey, @nRetTemp OUTPUT
		IF @nRetTemp <> 0
		BEGIN
			ROLLBACK TRAN
			SET @nRet = 7
			RETURN
		END
	END
	IF @nRawRightLot > 1
	BEGIN
		SET @nRetTemp = 1
		SET @nRawRightLot = @nRawRightLot - 1
		EXEC p_Item_SetOption @nRawRightKey, @nRawRightOptionType, @nRawRightLot, @nRetTemp OUTPUT
		IF @nRetTemp <> 0
		BEGIN
			ROLLBACK TRAN
			SET @nRet = 7
			RETURN
		END
	END
END

-- 사용된 자원가운데 지우기(입력된 값이 있으면)
IF @nRawMiddleKey <> 0
BEGIN
	IF @nRawMiddleLot = 1
	BEGIN
		SET @nRetTemp = 1
		EXEC p_Item_Delete @nRawMiddleKey, @nRetTemp OUTPUT
		IF @nRetTemp <> 0
		BEGIN
			ROLLBACK TRAN
			SET @nR
GO

-- ========================================================================
-- Procedure: p_Item_Upgrade_GOLD_NINE
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Item_Upgrade_GOLD_NINE]
/*
Item upgrade
2005.6 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK
	        1 = Not found item
	        2 = Missmatch item owner
	        3 = Not enough money to sub
	        4 = Upgrade value setting, upgrade fail count error(option set error)
	        5 = Raw delete fail
	        6 = RawLeft delete fail, if valid item
	        7 = RawRight delete fail, if valid item
	        8 = RawMiddle delete fail, if valid item
	        
 operation:
	1. money sub
	2. upgrade item level
	3. delete raws
*/
@nOwner int,
@nItemKey bigint,
@nItemDestroy int,
@nUpgradeOptionType smallint,
@nUpgradeOptionValue bigint,
@nUpgradeFailCountType smallint,
@nUpgradeFailCountValue bigint,
@nRawKey bigint,
@nRawLeftKey bigint,
@nRawRightKey bigint,
@nRawMiddleKey bigint,
@nGiftMoney bigint,
@nRet int OUTPUT
AS
SET NOCOUNT ON
-- 아이템 소유자 확인
DECLARE @nTinyintRet tinyint
EXEC p_Item_IsValidOwner @nItemKey, @nOwner, @nTinyintRet OUTPUT
IF @nTinyintRet <> 0
BEGIN
	SET @nRet = @nTinyintRet
	RETURN
END
-- 아이템 소유자의 돈 감소
BEGIN TRAN
IF @nGiftMoney > 0 -- 수고비가 있는 경우, 금액을 빼준다.
BEGIN
	EXEC p_Char_MoneySub @nOwner, @nGiftMoney, @nTinyintRet OUTPUT
	IF @nTinyintRet <> 1
	BEGIN
		ROLLBACK TRAN
		SET @nRet = 3
		RETURN
	END
END	
-- 아이템 업그레이드 값 조정
DECLARE @nRetTemp int
IF @nItemDestroy = 1 -- 업글 실패로 아이템 파괴해야 되는경우
BEGIN
	SET @nRetTemp = 1
	EXEC p_Item_Delete @nItemKey, @nRetTemp OUTPUT
	IF @nRetTemp <> 0
	BEGIN
		ROLLBACK TRAN
		SET @nRet = 4
		RETURN
	END
END
ELSE
BEGIN
	-- 업글 실패로 파괴가 아닌 경우, 새로운 업글값, 실패 카운트 기록	
	SET @nRetTemp = 1
	EXEC p_Item_SetOption @nItemKey, @nUpgradeOptionType, @nUpgradeOptionValue, @nRetTemp OUTPUT
	IF @nRetTemp <> 0
	BEGIN
		ROLLBACK TRAN
		SET @nRet = 4
		RETURN
	END
	
	SET @nRetTemp = 1
	EXEC p_Item_SetOption @nItemKey, @nUpgradeFailCountType, @nUpgradeFailCountValue, @nRetTemp OUTPUT
	IF @nRetTemp <> 0
	BEGIN
		ROLLBACK TRAN
		SET @nRet = 4
		RETURN
	END
END
-- 사용된 자원 지우기
SET @nRetTemp = 1
EXEC p_Item_Delete @nRawKey, @nRetTemp OUTPUT
IF @nRetTemp <> 0
BEGIN
	ROLLBACK TRAN
	SET @nRet = 5
	RETURN
END
-- 사용된 자원왼족 지우기(입력된 값이 있으면)
IF @nRawLeftKey <> 0
BEGIN
	SET @nRetTemp = 1
	EXEC p_Item_Delete @nRawLeftKey, @nRetTemp OUTPUT
	IF @nRetTemp <> 0
	BEGIN
		ROLLBACK TRAN
		SET @nRet = 6
		RETURN
	END
END
-- 사용된 자원오른쪽 지우기(입력된 값이 있으면)
IF @nRawRightKey <> 0
BEGIN
	SET @nRetTemp = 1
	EXEC p_Item_Delete @nRawRightKey, @nRetTemp OUTPUT
	IF @nRetTemp <> 0
	BEGIN
		ROLLBACK TRAN
		SET @nRet = 7
		RETURN
	END
END
-- 사용된 자원가운데 지우기(입력된 값이 있으면)
IF @nRawMiddleKey <> 0
BEGIN
	SET @nRetTemp = 1
	EXEC p_Item_Delete @nRawMiddleKey, @nRetTemp OUTPUT
	IF @nRetTemp <> 0
	BEGIN
		ROLLBACK TRAN
		SET @nRet = 8
		RETURN
	END
END

-- 업글 성공
COMMIT TRAN
SET @nRet = 0
-- end

0封赪u
GO

-- ========================================================================
-- Procedure: p_Item_Upgrade_OLD
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Item_Upgrade_OLD]
/*
Item upgrade
2005.6 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK
	        1 = Not found item
	        2 = Missmatch item owner
	        3 = Not enough money to sub
	        4 = Upgrade value setting, upgrade fail count error(option set error)
	        5 = Raw delete fail
	        6 = RawLeft delete fail, if valid item
	        7 = RawRight delete fail, if valid item
	        
 operation:
	1. money sub
	2. upgrade item level
	3. delete raws
*/
@nOwner int,
@nItemKey bigint,
@nItemDestroy int,
@nUpgradeOptionType smallint,
@nUpgradeOptionValue bigint,
@nUpgradeFailCountType smallint,
@nUpgradeFailCountValue bigint,
@nRawKey bigint,
@nRawLeftKey bigint,
@nRawRightKey bigint,
@nGiftMoney bigint,
@nRet int OUTPUT
AS
SET NOCOUNT ON
-- 아이템 소유자 확인
DECLARE @nTinyintRet tinyint
EXEC p_Item_IsValidOwner @nItemKey, @nOwner, @nTinyintRet OUTPUT
IF @nTinyintRet <> 0
BEGIN
	SET @nRet = @nTinyintRet
	RETURN
END
-- 아이템 소유자의 돈 감소
BEGIN TRAN
IF @nGiftMoney > 0 -- 수고비가 있는 경우, 금액을 빼준다.
BEGIN
	EXEC p_Char_MoneySub @nOwner, @nGiftMoney, @nTinyintRet OUTPUT
	IF @nTinyintRet <> 1
	BEGIN
		ROLLBACK TRAN
		SET @nRet = 3
		RETURN
	END
END	
-- 아이템 업그레이드 값 조정
DECLARE @nRetTemp int
IF @nItemDestroy = 1 -- 업글 실패로 아이템 파괴해야 되는경우
BEGIN
	SET @nRetTemp = 1
	EXEC p_Item_Delete @nItemKey, @nRetTemp OUTPUT
	IF @nRetTemp <> 0
	BEGIN
		ROLLBACK TRAN
		SET @nRet = 4
		RETURN
	END
END
ELSE
BEGIN
	-- 업글 실패로 파괴가 아닌 경우, 새로운 업글값, 실패 카운트 기록	
	SET @nRetTemp = 1
	EXEC p_Item_SetOption @nItemKey, @nUpgradeOptionType, @nUpgradeOptionValue, @nRetTemp OUTPUT
	IF @nRetTemp <> 0
	BEGIN
		ROLLBACK TRAN
		SET @nRet = 4
		RETURN
	END
	
	SET @nRetTemp = 1
	EXEC p_Item_SetOption @nItemKey, @nUpgradeFailCountType, @nUpgradeFailCountValue, @nRetTemp OUTPUT
	IF @nRetTemp <> 0
	BEGIN
		ROLLBACK TRAN
		SET @nRet = 4
		RETURN
	END
END
-- 사용된 자원 지우기
SET @nRetTemp = 1
EXEC p_Item_Delete @nRawKey, @nRetTemp OUTPUT
IF @nRetTemp <> 0
BEGIN
	ROLLBACK TRAN
	SET @nRet = 5
	RETURN
END
-- 사용된 자원왼족 지우기(입력된 값이 있으면)
IF @nRawLeftKey <> 0
BEGIN
	SET @nRetTemp = 1
	EXEC p_Item_Delete @nRawLeftKey, @nRetTemp OUTPUT
	IF @nRetTemp <> 0
	BEGIN
		ROLLBACK TRAN
		SET @nRet = 6
		RETURN
	END
END
-- 사용된 자원오른쪽 지우기(입력된 값이 있으면)
IF @nRawRightKey <> 0
BEGIN
	SET @nRetTemp = 1
	EXEC p_Item_Delete @nRawRightKey, @nRetTemp OUTPUT
	IF @nRetTemp <> 0
	BEGIN
		ROLLBACK TRAN
		SET @nRet = 7
		RETURN
	END
END
-- 업글 성공
COMMIT TRAN
SET @nRet = 0
-- end

0锁膎v
GO

-- ========================================================================
-- Procedure: p_KingdomQuest_Del
-- ========================================================================
CREATE PROCEDURE [dbo].[p_KingdomQuest_Del]
/*
Kingdom Quest Delete
2005.8 By CJC
 input:
	nQuestNo = To be delete info.
 output:
	@nRet	0 = OK, 
	        ? = Error code
*/
@nQuestNo int,
@nRet int OUTPUT

AS

SET NOCOUNT ON

DELETE FROM tKingdomQuest WHERE nQuestNo = @nQuestNo

IF (@@ERROR=0) AND (@@ROWCOUNT=1)
	SET @nRet = 0
ELSE
	SET @nRet = 1
	
-- end

0瘁䧛Ÿ
GO

-- ========================================================================
-- Procedure: p_KingdomQuest_Get
-- ========================================================================
CREATE PROCEDURE [dbo].[p_KingdomQuest_Get]
/*
Kingdom Quest Get
2005.8 By CJC
 input:
	see SQL
 output:
	@nRet	0 = No data found
	        ? = Number data here(must to be 1)
*/

@nQuestNo int,
@sData varchar(256) OUTPUT,
@nRet int OUTPUT

AS

SET NOCOUNT ON

SELECT @sData = sData FROM tKingdomQuest WHERE nQuestNo = @nQuestNo

SET @nRet = @@ROWCOUNT

-- end

0ጁ₭ŷ
GO

-- ========================================================================
-- Procedure: p_KingdomQuest_GetAll
-- ========================================================================
CREATE PROCEDURE [dbo].[p_KingdomQuest_GetAll]
/*
Kingdom Quest Get All
2005.8 By CJC
 input:
	see SQL
 output:
	@nRet	0 = No data found
	        ? = Number data here(must to be 1)
*/

AS

SET NOCOUNT ON

SELECT nQuestNo, sData FROM tKingdomQuest

-- end

0㈣ź
GO

-- ========================================================================
-- Procedure: p_KingdomQuest_GetNumOfKQ
-- ========================================================================
CREATE PROCEDURE [dbo].[p_KingdomQuest_GetNumOfKQ]
/*
Kingdom Quest Get number of KQ
2005.8 By CJC
 input:
	see SQL
 output:
	@nRet	0 = No data found
	        ? = Number data here(must to be 1)
*/

@nNumOfKQ int OUTPUT

AS

SET NOCOUNT ON

SELECT @nNumOfKQ = COUNT(*) FROM tKingdomQuest

-- end

0蔁ࣵŹ
GO

-- ========================================================================
-- Procedure: p_KingdomQuest_Set
-- ========================================================================
CREATE PROCEDURE [dbo].[p_KingdomQuest_Set]
/*
Kingdom Quest Set
2005.8 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK, 
	        ? = Error code
*/
@nQuestNo smallint,
@sData varchar(256),
@nRet int OUTPUT

AS

SET NOCOUNT ON

UPDATE tKingdomQuest SET sData = @sData WHERE nQuestNo = @nQuestNo

IF @@ERROR <> 0 OR @@ROWCOUNT = 0
BEGIN
	INSERT tKingdomQuest ( nQuestNo, sData ) VALUES ( @nQuestNo, @sData )
END

SET @nRet = @@ERROR

-- end

0䬁孒Ż
GO

-- ========================================================================
-- Procedure: p_MasPup_DelAll
-- ========================================================================
CREATE PROCEDURE [dbo].[p_MasPup_DelAll]
/*
Delete My Master and Pupils
2004.12 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK, 
	        ? = Error code
*/
@nCharNo int,
@nRet int OUTPUT
AS
SET NOCOUNT ON
-- 아이템 삭제 옵션
DELETE FROM tMasterPupils WHERE (nMaster = @nCharNo) OR (nPupils = @nCharNo)  
SET @nRet = @@ERROR
-- end

0ź
GO

-- ========================================================================
-- Procedure: p_MasPup_DelMaster
-- ========================================================================
CREATE PROCEDURE [dbo].[p_MasPup_DelMaster]
/*
Delete My Master
2004.12 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK, 
	        ? = Error code
*/
@nCharNo int,   -- As Pupils
@nRet int OUTPUT
AS
SET NOCOUNT ON
DELETE FROM tMasterPupils WHERE nPupils = @nCharNo
SET @nRet = @@ERROR
-- end

0찁ʴž
GO

-- ========================================================================
-- Procedure: p_MasPup_DelPupils
-- ========================================================================
CREATE PROCEDURE [dbo].[p_MasPup_DelPupils]
/*
Delete My Pupils
2004.12 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK, 
			? = Error code
*/
@nCharNo int,   -- As Master
@nPupils int,
@nRet int OUTPUT
AS
SET NOCOUNT ON
DELETE FROM tMasterPupils WHERE (nPupils = @nPupils) AND (nMaster = @nCharNo)
SET @nRet = @@ERROR
-- end

0ꈁ춪Ž
GO

-- ========================================================================
-- Procedure: p_MasPup_GetAll
-- ========================================================================
CREATE PROCEDURE [dbo].[p_MasPup_GetAll]
/*
Get My Master and Pupils
2004.12 By CJC
 input:
	see SQL
 output:
	fetch to take data
*/
@nCharNo int
AS
SET NOCOUNT ON
SELECT nMaster, nPupils, nFlags FROM tMasterPupils WHERE (nMaster = @nCharNo) OR (nPupils = @nCharNo)  
-- end

0㸁ſ
GO

-- ========================================================================
-- Procedure: p_MasPup_SetMaster
-- ========================================================================
CREATE PROCEDURE [dbo].[p_MasPup_SetMaster]
/*
Set My Master
2004.12 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK, 
			1 = 자신과는 스승관계를 맺을수 없다.
			2 = 제자가 4명이상 기록이 안된다.
			? = Error code(또는 스승이 2명이상 기록안된다)
*/
@nCharNo int,	-- As Pupils
@nMaster int,
@nFlags  int,
@nRet int OUTPUT
AS
SET NOCOUNT ON
-- 자신과 사제를 맺는지 확인한다.
IF @nMaster = @nCharNo
BEGIN
	SET @nRet = 1
	RETURN
END
-- 제자가 4명이상 기록이 안되게 확인
DECLARE @PupilsCount int
SELECT @PupilsCount = COUNT(*) FROM tMasterPupils WHERE nMaster = @nMaster
IF @PupilsCount >= 3
BEGIN
	SET @nRet = 2
	RETURN
END
-- 기록된 사항에 수정인 경우 수정처리 한다.
UPDATE tMasterPupils SET 
	nMaster = @nMaster,
	nFlags  = @nFlags
WHERE nPupils = @nCharNo
-- 처음 기록의 경우 처리한다.
IF @@ERROR <> 0 OR @@ROWCOUNT = 0
BEGIN
	INSERT tMasterPupils (  nMaster,  nPupils,  nFlags )
	VALUES				 ( @nMaster, @nCharNo, @nFlags )
END
SET @nRet = @@ERROR
-- end
GO

-- ========================================================================
-- Procedure: p_MasPup_SetPupils
-- ========================================================================
CREATE PROCEDURE [dbo].[p_MasPup_SetPupils]
/*
Set My Pupils
2004.12 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK,또는 이미 기록 되어 있음 
			1 = 자신과 사제를 맺을수 없다.
			2 = 제자가 4명이상 기록이 안된다.
			? = Error code(또는 스승이 2명이상 기록안된다)
*/
@nCharNo int, -- As Master
@nPupils int,
@nFlags  int,
@nRet int OUTPUT
AS
SET NOCOUNT ON
-- 자신과 사제를 맺는지 확인한다.
IF @nPupils = @nCharNo
BEGIN
	SET @nRet = 1
	RETURN
END
-- 이미 기록이 된경우 확인, Flag만 바꾼다.
IF EXISTS( SELECT * FROM tMasterPupils WHERE (nMaster = @nCharNo) AND ( nPupils = @nPupils ) )
BEGIN
	UPDATE tMasterPupils SET nFlags = @nFlags WHERE (nMaster = @nCharNo) AND (nPupils = @nPupils)
	SET @nRet = @@ERROR
	RETURN
END
-- 제자가 4명이상 기록이 안되게 확인
DECLARE @PupilsCount int
SELECT @PupilsCount = COUNT(*) FROM tMasterPupils WHERE nMaster = @nCharNo
IF @PupilsCount >= 3
BEGIN
	SET @nRet = 2
	RETURN
END
-- 새로 추가된 제자 기록(스승이 2명이상 기록안되게 체크는 nPupils가 Alternative key로 확인된다)
INSERT tMasterPupils (  nMaster,  nPupils,  nFlags )
VALUES				 ( @nCharNo, @nPupils, @nFlags )
SET @nRet = @@ERROR
-- end

0ࠁ＝
GO

-- ========================================================================
-- Procedure: p_Prison_Add
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Prison_Add] 
@sGmUserID nvarchar(16),
@sImprison nvarchar(16),
@nMinute smallint,
@sReason nvarchar(16), 
@sRemark nvarchar(64),
@nImprisonRegNum int,
@err smallint OUTPUT


AS
SET NOCOUNT ON

DECLARE	
@nRet		int,
@nPrisonMin 	int

SET 		@err = 0

BEGIN TRAN


UPDATE tCharacter SET 
	nPrisonMin=@nMinute ,
	nLoginZoneX=896,
	nLoginZoneY=1122,
	sLoginZone='EldPri' 	
	WHERE nCharNo=@nImprisonRegNum


IF @@error <> 0 GOTO ERR

INSERT tPrison (nCharID,nMinute,sReason,sGmID,sRemark)
VALUES (@nImprisonRegNum,@nMinute,@sReason,@sGmUserID,@sRemark )
IF @@error <> 0 GOTO ERR

COMMIT TRAN
RETURN @@error

ERR:
ROLLBACK TRAN
return @@error

0蘁鸻ā
GO

-- ========================================================================
-- Procedure: p_Prison_Add2
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Prison_Add2]
@sGmUserID nvarchar(16),
@sImprison nvarchar(16),
@nMinute smallint,
@sReason nvarchar(16), 
@sRemark nvarchar(64),
@err smallint OUTPUT


AS
SET NOCOUNT ON

DECLARE	
@nRet		int,
@nPrisonMin 	int,
@nImprisonRegNum int 

SET 		@err = 0

	--케릭터ID가 존재하지 않습니다.
SELECT @nImprisonRegNum=nCharNo, @nPrisonMin=nPrisonMin FROM tCharacter(NOLOCK) WHERE sID=@sImprison 
IF @@ROWCOUNT = 0
BEGIN
	SET @err = 1
	RETURN @@error
END
	--이미 수감중입니다
IF @nPrisonMin <> 0
BEGIN

	SET @err = 2
	RETURN @@error
END

BEGIN TRAN

UPDATE tCharacter SET
 	nPrisonMin=@nMinute, 
	nLoginZoneX=896,
	nLoginZoneY=1122,
	sLoginZone='EldPri' 		
	WHERE nCharNo=@nImprisonRegNum
IF @@error <> 0 GOTO ERR

INSERT tPrison (nCharID,nMinute,sReason,sGmID,sRemark)
VALUES (@nImprisonRegNum,@nMinute,@sReason,@sGmUserID,@sRemark )
IF @@error <> 0 GOTO ERR

COMMIT TRAN
RETURN @@error

ERR:
ROLLBACK TRAN
return @@error

0稁
GO

-- ========================================================================
-- Procedure: p_Prison_Alter
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Prison_Alter] 
@sGmUserID nvarchar(16),
@sImprison nvarchar(16),
@nMinute smallint,
@sReason nvarchar(16), 
@sRemark nvarchar(64),
@nImprisonRegNum int,
@err smallint OUTPUT


AS
SET NOCOUNT ON

DECLARE	
@nRet		int,
@nPrisonMin 	int

SET 		@err = 0

BEGIN TRAN

IF @nMinute = 0
	BEGIN
		UPDATE tCharacter SET
 		nPrisonMin=@nMinute, 
		nLoginZoneX=5963,
		nLoginZoneY=9793,
		sLoginZone='Rou' 		
		WHERE nCharNo=@nImprisonRegNum
	END
ELSE
	BEGIN
		UPDATE tCharacter SET
 		nPrisonMin=@nMinute, 
		nLoginZoneX=896,
		nLoginZoneY=1122,
		sLoginZone='EldPri' 		
		WHERE nCharNo=@nImprisonRegNum

	END

IF @@error <> 0 GOTO ERR

INSERT tPrison (nCharID,nMinute,sReason,sGmID,sRemark)
VALUES (@nImprisonRegNum,@nMinute,@sReason,@sGmUserID,@sRemark )
IF @@error <> 0 GOTO ERR

COMMIT TRAN
RETURN @@error

ERR:
ROLLBACK TRAN
return @@error

0蚃ă
GO

-- ========================================================================
-- Procedure: p_Prison_Alter2
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Prison_Alter2]
@sGmUserID nvarchar(16),
@sImprison nvarchar(16),
@nMinute smallint,
@sReason nvarchar(16), 
@sRemark nvarchar(64),
@err smallint OUTPUT


AS
SET NOCOUNT ON

DECLARE	
@nRet		int,
@nPrisonMin 	int,
@nImprisonRegNum int 

SET 		@err = 0

	--케릭터ID가 존재하지 않습니다.
SELECT @nImprisonRegNum=nCharNo, @nPrisonMin=nPrisonMin FROM tCharacter(NOLOCK) WHERE sID=@sImprison 
IF @@ROWCOUNT = 0
BEGIN
	SET @err = 1
	RETURN @@error
END

--수감중이 아닐때...
IF @nPrisonMin = 0
BEGIN

	SET @err = 2
	RETURN @@error
END

BEGIN TRAN
IF @nMinute = 0
	BEGIN
		UPDATE tCharacter SET
 		nPrisonMin=@nMinute, 
		nLoginZoneX=5963,
		nLoginZoneY=9793,
		sLoginZone='Rou' 		
		WHERE nCharNo=@nImprisonRegNum
	END
ELSE
	BEGIN
		UPDATE tCharacter SET
 		nPrisonMin=@nMinute, 
		nLoginZoneX=896,
		nLoginZoneY=1122,
		sLoginZone='EldPri' 		
		WHERE nCharNo=@nImprisonRegNum

	END

IF @@error <> 0 GOTO ERR

INSERT tPrison (nCharID,nMinute,sReason,sGmID,sRemark)
VALUES (@nImprisonRegNum,@nMinute,@sReason,@sGmUserID,@sRemark )
IF @@error <> 0 GOTO ERR

COMMIT TRAN
RETURN @@error

ERR:
ROLLBACK TRAN
return @@error

0쾭
GO

-- ========================================================================
-- Procedure: p_Prison_End
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Prison_End] 
@nImprisonRegNum	int

AS
SET NOCOUNT ON

UPDATE tCharacter SET 	nLoginZoneX=5963, 
			nLoginZoneY=9793, 
			nLoginZoneD=90,
			sLoginZone='Rou', 
			nPrisonMin=0 
	WHERE nCharNo=@nImprisonRegNum

RETURN @@error

0樁滌ą
GO

-- ========================================================================
-- Procedure: p_Prison_Get
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Prison_Get] 
@nCharID	int


AS
SET NOCOUNT ON

SELECT TOP 1 sReason, sRemark FROM tPrison
WHERE nCharID = @nCharID 
ORDER BY dPrisonDate DESC

RETURN @@error

0؁谟ć
GO

-- ========================================================================
-- Procedure: p_Prison_UpdateCharPrisonMin
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Prison_UpdateCharPrisonMin] 
@nImprisonRegNum	int, 
@nMinute		smallint

AS
SET NOCOUNT ON

--감옥시간과 로긴존과의 동기화를 위해서 꼭!!! 같이 업데이트 해줘야함.
IF @nMinute = 0
	BEGIN
		UPDATE tCharacter SET
 		nPrisonMin=@nMinute, 
		nLoginZoneX=5963,
		nLoginZoneY=9793,
		sLoginZone='Rou' 		
		WHERE nCharNo=@nImprisonRegNum
	END
ELSE
	BEGIN
		UPDATE tCharacter SET
 		nPrisonMin=@nMinute, 
		sLoginZone='EldPri' 		
		WHERE nCharNo=@nImprisonRegNum

	END

RETURN @@error

0ꌁ拰Ć
GO

-- ========================================================================
-- Procedure: p_Quest_Add
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Quest_Add]
/*
Quest Set
2004.10 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK, 
	        ? = Error code
*/
@nCharNo int,
@nQuestNo int,
@nStatus tinyint,
@sData varbinary(100),
@nRet int OUTPUT

AS

SET NOCOUNT ON

INSERT tQuest (  nCharNo,  nQuestNo,  nStatus,  sData )
VALUES        ( @nCharNo, @nQuestNo, @nStatus, @sData )

SET @nRet = @@ERROR

-- end

0码瑧ĉ
GO

-- ========================================================================
-- Procedure: p_Quest_Del
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Quest_Del]
/*
Quest Delete
2004.10 By CJC
 input:
	nCharNo, nQuestNo = To be delete info.
 output:
	@nRet	0 = OK, 
	        ? = Error code
*/
@nCharNo int,
@nQuestNo int,
@nRet int OUTPUT

AS

SET NOCOUNT ON

DELETE FROM tQuest WHERE nCharNo = @nCharNo AND nQuestNo = @nQuestNo
IF (@@ERROR=0) AND (@@ROWCOUNT=1)
	SET @nRet = 0
ELSE
	SET @nRet = 1
-- end

0ᔁ䬹Ĉ
GO

-- ========================================================================
-- Procedure: p_Quest_Get
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Quest_Get]
/*
Quest Get
2004.10 By CJC
 input:
	see SQL
 output:
	@nRet	0 = No data found
	        ? = Number data here(must to be 1)
*/
@nCharNo int,
@nQuestNo int,
@nStatus tinyint OUTPUT,
@sData varbinary(100) OUTPUT,
@nRet int OUTPUT

AS

SET NOCOUNT ON

SELECT 
	@nStatus = nStatus, 
	@sData = sData
FROM tQuest WHERE nCharNo = @nCharNo AND nQuestNo = @nQuestNo

SET @nRet = @@ROWCOUNT

-- end

0鶕Ċ
GO

-- ========================================================================
-- Procedure: p_Quest_GetAll
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Quest_GetAll]
/*
Quest Get All
2004.10 By CJC
 input:
	see SQL
 output:
	fetch to take data
*/

@nCharNo int

AS

SET NOCOUNT ON

SELECT nQuestNo, nStatus, sData FROM tQuest WHERE nCharNo = @nCharNo

-- end

0蜁㎁Ċ
GO

-- ========================================================================
-- Procedure: p_Quest_GetAllDoing
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Quest_GetAllDoing]
/*
Quest Get All
2004.10 By CJC
 input:
	see SQL
 output:
	fetch to take data
*/
@nCharNo int,
@nStatusDone tinyint
AS
SET NOCOUNT ON
SELECT nQuestNo, nStatus, sData FROM tQuest WHERE nCharNo = @nCharNo AND nStatus = @nStatusDone
-- end

0
GO

-- ========================================================================
-- Procedure: p_Quest_GetAllDone
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Quest_GetAllDone]
/*
Quest Get All
2004.10 By CJC
 input:
	see SQL
 output:
	fetch to take data
*/

@nCharNo int,
@nStatusDone tinyint

AS

SET NOCOUNT ON

SELECT nQuestNo FROM tQuest WHERE nCharNo = @nCharNo AND nStatus = @nStatusDone

-- end

0锁㤜Ď
GO

-- ========================================================================
-- Procedure: p_Quest_Set
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Quest_Set]
/*
Quest Set
2004.10 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK, 
	        ? = Error code
*/
@nCharNo int,
@nQuestNo int,
@nStatus tinyint,
@sData varbinary(100),
@nRet int OUTPUT
AS

SET NOCOUNT ON

--------------------------------------------------------- add start
-- added by CJC for for wedding quest 2008-06-03
DECLARE @nOldStatus tinyint
IF EXISTS(SELECT nStatus FROM tQuest WHERE nCharNo=@nCharNo and nQuestNo=@nQuestNo)
	SELECT @nOldStatus = nStatus FROM tQuest WHERE nCharNo = @nCharNo AND nQuestNo = @nQuestNo
ELSE
	SET @nOldStatus = 0
/*-------------------------------------------------------------------------------------
	Project: FiestaHeroes - NA2016
	Author: xkl
	
	Quest item fix. Simply add the Quest ID and Item ID to the table, and it will
	prevent players from repeatedly turning in the same item for the quest.
------------------------------------------------------------------------------------*/
	DECLARE @QuestID INT, @ItemID INT;
	SELECT @QuestID = QuestID, @ItemID = ItemID FROM QuestItemFix WHERE QuestID = @nQuestNo;
	
	IF @nQuestNo = @QuestID AND @nStatus = 6
	BEGIN
		IF EXISTS (SELECT nItemID FROM tItem WHERE nOwner = @nCharNo AND nItemID = @ItemID)
		RETURN;
	END
/*-------------------------------------------------------------------------------------
	Continue as normal ~ 
------------------------------------------------------------------------------------*/

UPDATE tQuest SET 
    nStatus = @nStatus,
	sData = @sData
WHERE nCharNo = @nCharNo AND nQuestNo = @nQuestNo

IF @@ERROR <> 0 OR @@ROWCOUNT = 0
BEGIN
	INSERT tQuest (  nCharNo,  nQuestNo,  nStatus,  sData )
	VALUES        ( @nCharNo, @nQuestNo, @nStatus, @sData )
END
-- end of event on 2008-06-24 KST
--------------------------------------------------------- add start
--added by Amadeus for wedding quest 2008-05-29, modify by CJC 2008-06-03
-- nStatus 4 = PQS_REPEAT, Repeat quest done, now re-acceptable.
-- nStatus 6 = PQS_IN_PROGRESS 

/*IF @nStatus = 4 and @nQuestNo=20016 and @nOldStatus = 6
BEGIN
	IF EXISTS(SELECT nCharNo FROM tQuestTimes WHERE nCharNo=@nCharNo and nQuestNo=@nQuestNo)
		BEGIN
			UPDATE tQuestTimes SET nTimes=nTimes+1, dLastComplete=getdate()
			WHERE nCharNo=@nCharNo and nQuestNo=@nQuestNo
		END
	ELSE
		BEGIN
			INSERT tQuestTimes 
			VALUES(@nCharNo, @nQuestNo, 1, getdate())
		END
END
--------------------------------------------------------- end of edit by Amadeus
*/
SET @nRet = @@ERROR
--end

<
GO

-- ========================================================================
-- Procedure: p_Skill_Del
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Skill_Del]
/*
Skill Delete
2004.10 By CJC
 input:
	nCharNo, nSkillNo = To be delete info.
 output:
	@nRet	0 = OK, 
	        ? = Error code
*/
@nCharNo int,
@nSkillNo smallint,
@nRet int OUTPUT
AS
SET NOCOUNT ON
-- 아이템 삭제 옵션
DELETE FROM tSkill WHERE nCharNo = @nCharNo AND nSkillNo = @nSkillNo
SET @nRet = @@ERROR
-- end

0ꐁĎ
GO

-- ========================================================================
-- Procedure: p_Skill_Get
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Skill_Get]
/*
Skill Get
2004.10 By CJC
 input:
	see SQL
 output:
	@nRet	0 = No data found
	        ? = Number data here(must to be 1)
*/
@nCharNo int,
@nSkillNo smallint,
@nSkillLevel tinyint OUTPUT,
@nSkillExp int OUTPUT,
@nSkillWriteTime bigint OUTPUT,
@nSkillCoolTime int OUTPUT,
@nSkillPowerDemage tinyint OUTPUT,
@nSkillPowerSP tinyint OUTPUT,
@nSkillPowerKeepTime tinyint OUTPUT,
@nSkillPowerCoolTime tinyint OUTPUT,
@nRet int OUTPUT
AS
SET NOCOUNT ON
SELECT 
	@nSkillLevel		= nSkillLevel,
	@nSkillExp			= nSkillExp,
	@nSkillWriteTime	= nSkillWriteTime,
	@nSkillCoolTime		= nSkillCoolTime,
	@nSkillPowerDemage	= nSkillPowerDemage,
	@nSkillPowerSP		= nSkillPowerSP,
	@nSkillPowerKeepTime = nSkillPowerKeepTime,
	@nSkillPowerCoolTime = nSkillPowerCoolTime
FROM tSkill WHERE nCharNo = @nCharNo AND nSkillNo = @nSkillNo
SET @nRet = @@ROWCOUNT
-- end

0㄁噯Đ
GO

-- ========================================================================
-- Procedure: p_Skill_GetAll
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Skill_GetAll]
/*
Skill Get All
2004.10 By CJC
 input:
	see SQL
 output:
	fetch to take data
*/
@nCharNo int
AS
SET NOCOUNT ON
SELECT	nSkillNo, nSkillLevel, nSkillExp, nSkillWriteTime, nSkillCoolTime, 
		nSkillPowerDemage, nSkillPowerSP, nSkillPowerKeepTime, nSkillPowerCoolTime
FROM tSkill WHERE nCharNo = @nCharNo
-- end

0樁䪓đ
GO

-- ========================================================================
-- Procedure: p_Skill_GetNumOfSkill
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Skill_GetNumOfSkill]
/*
Skill Get Num of skill
2005.5 By CJC
 input:
	see SQL
 output:
	fetch to take data
*/
@nCharNo int,
@nNumOfSkill int OUTPUT
AS
SET NOCOUNT ON
SET @nNumOfSkill = 0
SELECT @nNumOfSkill = COUNT(*) FROM tSkill WHERE nCharNo = @nCharNo
-- end

0企풣đ
GO

-- ========================================================================
-- Procedure: p_Skill_ResetPower
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Skill_ResetPower]
/*
Skill Reset Power
2005.6 By CJC
 input:
	see SQL
 output:
	@nRet	1 = OK, 
	        0 = Error
*/
@nCharNo int,
@nRet int OUTPUT
AS
SET NOCOUNT ON
DECLARE @nFlagsTemp AS int
SET @nRet = 0 -- set error first
SELECT @nFlagsTemp = nFlags FROM tCharacter WHERE nCharNo = @nCharNo
IF @@ROWCOUNT = 0
	RETURN
BEGIN TRAN
SET @nFlagsTemp = @nFlagsTemp & 0x7fffffff -- Clear skillempower_can_reset's Bit
UPDATE tCharacter SET nFlags = @nFlagsTemp WHERE nCharNo = @nCharNo
UPDATE tSkill SET 
	nSkillPowerDemage	= 0,
	nSkillPowerSP		= 0,
	nSkillPowerKeepTime = 0,
	nSkillPowerCoolTime = 0
WHERE nCharNo = @nCharNo
IF @@ERROR <> 0
BEGIN
	ROLLBACK TRAN
	RETURN
END
COMMIT TRAN
SET @nRet = 1 -- no error return
-- end

0蠁죇Ē
GO

-- ========================================================================
-- Procedure: p_Skill_Set
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Skill_Set]
/*
Skill Set
2004.10 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK, 
	        ? = Error code
*/
@nCharNo int,
@nSkillNo smallint,
@nSkillLevel tinyint,
@nSkillExp int,
@nSkillWriteTime bigint,
@nSkillCoolTime int,
@nSkillPowerDemage tinyint,
@nSkillPowerSP tinyint,
@nSkillPowerKeepTime tinyint,
@nSkillPowerCoolTime tinyint,
@nRet int OUTPUT
AS
SET NOCOUNT ON
UPDATE tSkill SET
	nSkillLevel			= @nSkillLevel,
	nSkillExp			= @nSkillExp,
	nSkillWriteTime		= @nSkillWriteTime,
	nSkillCoolTime		= @nSkillCoolTime,
	nSkillPowerDemage	= @nSkillPowerDemage,
	nSkillPowerSP		= @nSkillPowerSP,
	nSkillPowerKeepTime = @nSkillPowerKeepTime,
	nSkillPowerCoolTime = @nSkillPowerCoolTime
WHERE nCharNo = @nCharNo AND nSkillNo = @nSkillNo
IF @@ERROR <> 0 OR @@ROWCOUNT = 0
BEGIN
	INSERT tSkill (  nCharNo,  nSkillNo,  nSkillLevel,  nSkillExp,  nSkillWriteTime,  nSkillCoolTime,  nSkillPowerDemage,  nSkillPowerSP,  nSkillPowerKeepTime,  nSkillPowerCoolTime )
	VALUES        ( @nCharNo, @nSkillNo, @nSkillLevel, @nSkillExp, @nSkillWriteTime, @nSkillCoolTime, @nSkillPowerDemage, @nSkillPowerSP, @nSkillPowerKeepTime, @nSkillPowerCoolTime )
END
SET @nRet = @@ERROR
-- end

0섁볫ē
GO

-- ========================================================================
-- Procedure: p_Skill_SetCoolTime
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Skill_SetCoolTime]
/*
Skill Set Cool Time
2004.10 By CJC
 input:
	see SQL
 output:
	@nRet	1 = OK, 
	        0 = Error
*/
@nCharNo int,
@nSkillNo smallint,
@nSkillWriteTime bigint,
@nSkillCoolTime int,
@nRet int OUTPUT
AS
SET NOCOUNT ON
UPDATE tSkill SET
	nSkillWriteTime = @nSkillWriteTime,
	nSkillCoolTime	= @nSkillCoolTime
WHERE nCharNo = @nCharNo AND nSkillNo = @nSkillNo
IF @@ERROR <> 0 OR @@ROWCOUNT = 0
	SET @nRet = 0
ELSE
	SET @nRet = 1
-- end

0度넏Ĕ
GO

-- ========================================================================
-- Procedure: p_Skill_SetPassive
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Skill_SetPassive]
/*
Skill Set Passive
2005.6 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK, 
	        ? = Error code
*/

@nCharNo int,
@nSkillNewNo smallint,
@bSkillOld tinyint,
@nSkillOldNo smallint,
@nRet int OUTPUT

AS

SET NOCOUNT ON

IF @bSkillOld = 1
BEGIN
	UPDATE tSkill 
	SET nSkillNo = @nSkillNewNo
	WHERE nCharNo = @nCharNo AND nSkillNo = @nSkillOldNo
	
	IF @@ERROR <> 0 OR @@ROWCOUNT = 0
		SET @nRet = 1 -- ERROR, No existed old skill
	ELSE
		SET @nRet = 0
		
	RETURN
END

INSERT tSkill (  nCharNo,     nSkillNo )
VALUES        ( @nCharNo, @nSkillNewNo )

SET @nRet = @@ERROR

-- end

0㌁ꔴĕ
GO

-- ========================================================================
-- Procedure: p_Skill_SetPower
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Skill_SetPower]
/*
Skill Set Power
2005.6 By CJC
 input:
	see SQL
 output:
	@nRet	1 = OK, 
	        0 = Error
*/
@nCharNo int,
@nSkillNo smallint,
@nSkillPowerSlot tinyint,
@nSkillPowerValue tinyint,
@nRet int OUTPUT
AS
BEGIN
	SET NOCOUNT ON

	IF @nSkillPowerValue > 5 BEGIN
		SET @nRet = 0
		RETURN
	END
	
	-- nSkillPowerDemage
	IF @nSkillPowerSlot = 0
	BEGIN
		UPDATE tSkill SET nSkillPowerDemage = @nSkillPowerValue WHERE nCharNo = @nCharNo AND nSkillNo = @nSkillNo
		IF @@ERROR <> 0 OR @@ROWCOUNT = 0
			SET @nRet = 0
		ELSE
			SET @nRet = 1
		RETURN
	END
	-- nSkillPowerSP
	IF @nSkillPowerSlot = 1
	BEGIN
		UPDATE tSkill SET nSkillPowerSP = @nSkillPowerValue WHERE nCharNo = @nCharNo AND nSkillNo = @nSkillNo
		IF @@ERROR <> 0 OR @@ROWCOUNT = 0
			SET @nRet = 0
		ELSE
			SET @nRet = 1
		RETURN
	END
	-- nSkillPowerKeepTime
	IF @nSkillPowerSlot = 2
	BEGIN
		UPDATE tSkill SET nSkillPowerKeepTime = @nSkillPowerValue WHERE nCharNo = @nCharNo AND nSkillNo = @nSkillNo
		IF @@ERROR <> 0 OR @@ROWCOUNT = 0
			SET @nRet = 0
		ELSE
			SET @nRet = 1
		RETURN
	END
	-- nSkillPowerCoolTime
	IF @nSkillPowerSlot = 3
	BEGIN
		UPDATE tSkill SET nSkillPowerCoolTime = @nSkillPowerValue WHERE nCharNo = @nCharNo AND nSkillNo = @nSkillNo
		IF @@ERROR <> 0 OR @@ROWCOUNT = 0
			SET @nRet = 0
		ELSE
			SET @nRet = 1
		RETURN
	END
	-- ERROR
	SET @nRet = 0;
	-- end
END

0
GO

-- ========================================================================
-- Procedure: p_Skill_SetPowerAll
-- ========================================================================
CREATE PROCEDURE [dbo].[p_Skill_SetPowerAll]
/*
Skill Set Power All
2005.6 By CJC
 input:
	see SQL
 output:
	@nRet	1 = OK, 
	        0 = Error
*/
@nCharNo int,
@nSkillNo smallint,
@nSkillPowerDemage tinyint,
@nSkillPowerSP tinyint,
@nSkillPowerKeepTime tinyint,
@nSkillPowerCoolTime tinyint,
@nRet int OUTPUT
AS
BEGIN
	SET NOCOUNT ON
	
	IF @nSkillPowerDemage > 5 OR @nSkillPowerSP > 5 OR @nSkillPowerKeepTime > 5 OR @nSkillPowerCoolTime > 5 BEGIN
		SET @nRet = 0
		RETURN
	END
	
	UPDATE tSkill SET 
		nSkillPowerDemage	= @nSkillPowerDemage,
		nSkillPowerSP		= @nSkillPowerSP,
		nSkillPowerKeepTime	= @nSkillPowerKeepTime,
		nSkillPowerCoolTime	= @nSkillPowerCoolTime
	WHERE nCharNo = @nCharNo AND nSkillNo = @nSkillNo
	IF @@ERROR <> 0 OR @@ROWCOUNT = 0
		SET @nRet = 0
	ELSE
		SET @nRet = 1
	-- end
END

0
GO

-- ========================================================================
-- Procedure: p_User_Deposit
-- ========================================================================
CREATE PROCEDURE [dbo].[p_User_Deposit]
/*
Set User deposit money
2006.5 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK
	        1 = Miss match UserNo and User's CharNo
	        2 = Char Money Sub Failed(Not enough money)
	        3 = Add deposit money faield
	        4 = Get result money faield
*/
@nUserNo int,
@nCharNo int,
@nMoney bigint,
@nResultUserMoney bigint OUTPUT,
@nRet tinyint OUTPUT
AS

SET NOCOUNT ON

SET @nRet = 0

-- Check UserNo and User's CharNo
DECLARE @nDBUserNo int
SET @nDBUserNo = 0

SELECT @nDBUserNo = nUserNo FROM tCharacter WHERE nCharNo = @nCharNo
IF @nDBUserNo <> @nUserNo
BEGIN
	SET @nRet = 1
	RETURN
END

-- Deposit money
BEGIN TRAN

EXEC p_Char_MoneySub @nCharNo, @nMoney, @nRet OUTPUT
IF @nRet <> 1
BEGIN
	SET @nRet = 2
	ROLLBACK TRAN
	RETURN
END

EXEC p_User_MoneyAdd @nUserNo, @nMoney, @nRet OUTPUT
IF @nRet <> 1
BEGIN
	SET @nRet = 3
	ROLLBACK TRAN
	RETURN
END

EXEC p_User_MoneyGet @nUserNo, @nResultUserMoney OUTPUT, @nRet OUTPUT
IF @nRet <> 1
BEGIN
	SET @nRet = 4
	ROLLBACK TRAN
	RETURN
END

-- Done
COMMIT TRAN

SET @nRet = 0
RETURN

0䄁꫏ę
GO

-- ========================================================================
-- Procedure: p_User_MoneyAdd
-- ========================================================================
CREATE PROCEDURE [dbo].[p_User_MoneyAdd]
/*
Add user money
2005.5 By CJC
 input:
	see SQL
 output:
	@nRet	1 = OK
	        0 = Error
*/
@nUserNo int,
@nAddMoney bigint,
@nRet tinyint OUTPUT

AS

SET NOCOUNT ON

IF @nAddMoney < 0
BEGIN
	SET @nRet = 0
	RETURN
END

DECLARE @nUserMoney bigint
EXEC p_User_MoneyGet @nUserNo, @nUserMoney OUTPUT, @nRet OUTPUT
IF @nRet = 0
	RETURN
-- Add money
SET @nUserMoney = @nUserMoney + @nAddMoney
EXEC p_User_MoneySet @nUserNo, @nUserMoney, @nRet OUTPUT
-- End

0ᜁ病ę
GO

-- ========================================================================
-- Procedure: p_User_MoneyGet
-- ========================================================================
CREATE PROCEDURE [dbo].[p_User_MoneyGet]
/*
Get user money
2005.5 By CJC
 input:
	see SQL
 output:
	@nRet	1 = OK, Money return
	        0 = Error
*/
@nUserNo int,
@nUserMoney bigint OUTPUT,
@nRet tinyint OUTPUT

AS

SET NOCOUNT ON

SET @nRet = 0

-- 돈 얻기
IF EXISTS(SELECT nUserMoney FROM tUserMoney WHERE nUserNo = @nUserNo )
BEGIN
	SELECT @nUserMoney = nUserMoney FROM tUserMoney WHERE nUserNo = @nUserNo
	IF @@ROWCOUNT = 1
		SET @nRet = 1
	RETURN 
END

-- First access
SET @nUserMoney = 0
INSERT INTO tUserMoney VALUES (@nUserNo, @nUserMoney)
IF @@ROWCOUNT = 1 AND @@ERROR = 0
	SET @nRet = 1
	
RETURN 

-- End

0倁槩Ě
GO

-- ========================================================================
-- Procedure: p_User_MoneyGetGreaterThen
-- ========================================================================
CREATE PROCEDURE [dbo].[p_User_MoneyGetGreaterThen]
	@nMoney		bigint
AS

BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	SELECT nUserNo, nUserMoney FROM tUserMoney WHERE nUserMoney >=  @nMoney
	
	SET LOCK_TIMEOUT -1
	SET NOCOUNT OFF	
END

0褁帍ě
GO

-- ========================================================================
-- Procedure: p_User_MoneySet
-- ========================================================================
CREATE PROCEDURE [dbo].[p_User_MoneySet]
/*
Set user money
2006.5 By CJC
 input:
	see SQL
 output:
	@nRet	1 = OK
	        0 = Error
*/
@nUserNo int,
@nSetMoney bigint,
@nRet tinyint OUTPUT
AS
SET NOCOUNT ON
SET @nRet = 0

IF EXISTS( SELECT nUserMoney FROM tUserMoney WHERE nUserNo = @nUserNo )
BEGIN
	UPDATE tUserMoney SET nUserMoney = @nSetMoney WHERE nUserNo = @nUserNo
	IF @@ERROR <> 0
		RETURN
END
ELSE
BEGIN
	INSERT INTO tUserMoney VALUES (@nUserNo, @nSetMoney)
	IF @@ERROR <> 0
		RETURN
END

-- DONE
SET @nRet = 1
RETURN

0숁刱Ĝ
GO

-- ========================================================================
-- Procedure: p_User_MoneySub
-- ========================================================================
CREATE PROCEDURE [dbo].[p_User_MoneySub]
/*
Sub user money
2005.5 By CJC
 input:
	see SQL
 output:
	@nRet	1 = OK
	        0 = Error, or Minus value
*/
@nUserNo int,
@nSubMoney bigint,
@nRet tinyint OUTPUT

AS

SET NOCOUNT ON

IF @nSubMoney < 0
BEGIN
	SET @nRet = 0
	RETURN
END

DECLARE @nUserMoney bigint
EXEC p_User_MoneyGet @nUserNo, @nUserMoney OUTPUT, @nRet OUTPUT
IF @nRet = 0
	RETURN
-- Sub money
SET @nUserMoney = @nUserMoney - @nSubMoney
-- Check minus
IF @nUserMoney < 0
BEGIN
	SET @nRet = 0
	RETURN	
END
EXEC p_User_MoneySet @nUserNo, @nUserMoney, @nRet OUTPUT
-- End

0企끪ĝ
GO

-- ========================================================================
-- Procedure: p_User_Withdraw
-- ========================================================================
CREATE PROCEDURE [dbo].[p_User_Withdraw]
/*
Set User withdraw money
2006.5 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK
	        1 = Miss match UserNo and User's CharNo
	        2 = User Money Sub Failed(Not enough money)
	        3 = Withdraw money faield
	        4 = Get result money faield
*/
@nUserNo int,
@nCharNo int,
@nMoney bigint,
@nResultUserMoney bigint OUTPUT,
@nRet tinyint OUTPUT
AS

SET NOCOUNT ON

SET @nRet = 0

-- Check UserNo and User's CharNo
DECLARE @nDBUserNo int
SET @nDBUserNo = 0

SELECT @nDBUserNo = nUserNo FROM tCharacter WHERE nCharNo = @nCharNo
IF @nDBUserNo <> @nUserNo
BEGIN
	SET @nRet = 1
	RETURN
END

-- Deposit money
BEGIN TRAN

EXEC p_User_MoneySub @nUserNo, @nMoney, @nRet OUTPUT
IF @nRet <> 1
BEGIN
	SET @nRet = 2
	ROLLBACK TRAN
	RETURN
END

EXEC p_Char_MoneyAdd @nCharNo, @nMoney, @nRet OUTPUT
IF @nRet <> 1
BEGIN
	SET @nRet = 3
	ROLLBACK TRAN
	RETURN
END

EXEC p_User_MoneyGet @nUserNo, @nResultUserMoney OUTPUT, @nRet OUTPUT
IF @nRet <> 1
BEGIN
	SET @nRet = 4
	ROLLBACK TRAN
	RETURN
END

-- Done
COMMIT TRAN

SET @nRet = 0
RETURN

0蠁꒎Ğ
GO

-- ========================================================================
-- Procedure: p_View_GetData
-- ========================================================================
CREATE PROCEDURE [dbo].[p_View_GetData]
	@QueryExpr nvarchar(512)

AS
BEGIN
	SET NOCOUNT ON
	DECLARE @SQLCMD nvarchar(512)
	SET @SQLCMD = @QueryExpr
	EXEC sp_executesql @SQLCMD
	SET NOCOUNT OFF
END

0轿
GO

-- ========================================================================
-- Procedure: usp_Attendance_getCallDateList
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Attendance_getCallDateList]
	@nCharNo    int
,	@dStartDate datetime
,	@dEndDate   datetime
,	@nRet       int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	SET @nRet = 0
	
	SELECT dAttendanceDate 
	  FROM tAttendance 
	 WHERE nCharNo = @nCharNo
	   AND dAttendanceDate >= @dStartDate 
	   AND dAttendanceDate <= @dEndDate 	
END

0ꘁ⋂Ġ
GO

-- ========================================================================
-- Procedure: usp_Attendance_getCount
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Attendance_getCount]
	@nCharNo    int
,	@dStartDate datetime
,	@dEndDate   datetime
,	@nRet       int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	SET @nRet = 0
	
	SELECT COUNT(dAttendanceDate) AS nAttendanceCount
	  FROM tAttendance 
	 WHERE nCharNo = @nCharNo
	   AND dAttendanceDate >= @dStartDate 
	   AND dAttendanceDate <= @dEndDate 	
END

0㌁胻ġ
GO

-- ========================================================================
-- Procedure: usp_Attendance_getRewardList
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Attendance_getRewardList]
	@nCharNo    int
,	@dDate      datetime
,	@nRet       int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	SET @nRet = 0
	
	SELECT dEventDate, nRewardID  
	  FROM tAttendanceReward  
	 WHERE nCharNo = @nCharNo
	   AND dEventDate > @dDate
END

0氁生Ģ
GO

-- ========================================================================
-- Procedure: usp_Attendance_setCall
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Attendance_setCall]
	@nCharNo        int
,	@dAttendaceDate datetime
,	@nRet       int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	SET @nRet = 0
	
	IF EXISTS(SELECT nCharNo FROM tAttendance WHERE nCharNo = @nCharNo AND dAttendanceDate = @dAttendaceDate) BEGIN
		SET @nRet = -1111
	END
	ELSE BEGIN
		INSERT INTO tAttendance (nCharNo, dAttendanceDate)
		VALUES (@nCharNo, @dAttendaceDate)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -9999
		END
	END
END

0儁ＯĢ
GO

-- ========================================================================
-- Procedure: usp_Attendance_setReward
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Attendance_setReward]
	@nCharNo     int
,	@dEventDate  datetime	
,	@nRewardID   int
,	@nRet        int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	SET @nRet = 0
	
	IF EXISTS(SELECT nCharNo FROM tAttendanceReward WHERE nCharNo = @nCharNo AND dEventDate = CAST(@dEventDate as datetime) AND nRewardID = @nRewardID) BEGIN
		SET @nRet = -1111
	END
	ELSE BEGIN
		INSERT INTO tAttendanceReward (nCharNo, dEventDate, dRewardDate, nRewardID)
		VALUES (@nCharNo, @dEventDate, GETDATE(), @nRewardID)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -9999
		END
	END
END

0嵧Ĥ
GO

-- ========================================================================
-- Procedure: usp_Char_GetAutoPick
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Char_GetAutoPick]
/*
2009.6 By temlost
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

	-- 삭제되지 않은 실존하는 케릭터인지 검색
	IF NOT EXISTS(SELECT nCharNo FROM tCharacter WHERE nCharNo = @nCharNo and bDeleted = 0)
	BEGIN
		SET @nRet = -4444
		RETURN
	END

	IF EXISTS(SELECT nCharNo, IsAutoPick FROM tAutoPick WHERE nCharNo = @nCharNo)
	BEGIN
		SELECT IsAutoPick FROM tAutoPick WHERE nCharNo = @nCharNo
		SET @nRet = 0
		RETURN
	END
	ELSE
	BEGIN
--		SELECT nCharNo, getdate() FROM tAutoPick WHERE nCharNo = @nCharNo
		SET @nRet = 0
		RETURN
	END
END

0ꠁ熇%
GO

-- ========================================================================
-- Procedure: usp_Char_GetRestExpExecTime
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Char_GetRestExpExecTime]
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

	IF EXISTS(SELECT nCharNo, dRestExpLastExcuteTime FROM tRestExp WHERE nCharNo = @nCharNo)
	BEGIN
		SELECT dRestExpLastExcuteTime FROM tRestExp WHERE nCharNo = @nCharNo
		SET @nRet = 0
		RETURN
	END
	ELSE
	BEGIN
--		SELECT nCharNo, getdate() FROM tRestExp WHERE nCharNo = @nCharNo
		SELECT getdate()
		SET @nRet = 0
		RETURN
	END
END

0
GO

-- ========================================================================
-- Procedure: usp_Char_SetAutoPick
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Char_SetAutoPick]
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
@nIsAutoPick		tinyint,
@nRet			int OUTPUT

AS
BEGIN

	SET NOCOUNT ON

	-- 삭제되지 않은 실존하는 케릭터인지 검색
	IF NOT EXISTS(SELECT nCharNo FROM tCharacter WHERE nCharNo = @nCharNo and bDeleted = 0)
	BEGIN
		SET @nRet = -5654;
		RETURN
	END

	IF EXISTS(SELECT nCharNo, IsAutoPick FROM tAutoPick WHERE nCharNo = @nCharNo)
	BEGIN
		UPDATE tAutoPick
		SET IsAutoPick = @nIsAutoPick
		WHERE nCharNo = @nCharNo

		SET @nRet = 0
		RETURN
	END
	ELSE
	BEGIN
		INSERT INTO tAutoPick (nCharNo, IsAutoPick)
		VALUES (@nCharNo, @nIsAutoPick)

		SET @nRet = 0
		RETURN
	END
END

0숁ⷸĨ
GO

-- ========================================================================
-- Procedure: usp_Char_SetRestExpExecTime
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Char_SetRestExpExecTime]
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

	-- 존재하지 않는 케릭터가 요청으로 들어오면 에러
	IF NOT EXISTS(SELECT nCharNo FROM tCharacter WHERE nCharNo = @nCharNo)
	BEGIN
		SET @nRet = -5654;
		RETURN
	END

	IF EXISTS(SELECT nCharNo, dRestExpLastExcuteTime FROM tRestExp WHERE nCharNo = @nCharNo)
	BEGIN
		UPDATE tRestExp
		SET dRestExpLastExcuteTime = getdate()
		WHERE nCharNo = @nCharNo

		SET @nRet = 0
		RETURN
	END
	ELSE
	BEGIN
		INSERT INTO tRestExp (nCharNo, dRestExpLastExcuteTime)
		VALUES (@nCharNo, getdate())

		SET @nRet = 0
		RETURN
	END
END

0꜁레Ĩ
GO

-- ========================================================================
-- Procedure: usp_Char_chkCharNo
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Char_chkCharNo]
	@nCharNo int
,	@sID     nvarchar(40)
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	IF EXISTS(SELECT nCharNo FROM tCharacter WHERE nCharNo = @nCharNo AND sID = @sID)
		SET @nRet = 0
	ELSE
		SET @nRet = -1

	SELECT nCharNo, sID FROM tCharacter WHERE nCharNo = @nCharNo

	SET NOCOUNT OFF
END

0䔁䡙$
GO

-- ========================================================================
-- Procedure: usp_CharacterLoading
-- ========================================================================
CREATE  PROCEDURE [dbo].[usp_CharacterLoading]
/*
char loading save
2009.6 By Luunar
 input:
	see SQL
 output:
	= 0 No Error
	< 0 Error
*/

@nCharNo		int,
@dWorldTime		datetime,
@sLoadingMapName	nvarchar(16),
@nFailType		smallint,
@nLoadingFailCnt	smallint,
@nRet			int OUTPUT

AS
BEGIN

	SET NOCOUNT ON

	DECLARE @nSeqNo int
	SET @nSeqNo       = 0

	BEGIN TRAN

	SELECT TOP 1 @nSeqNo = nSeqNo
	  FROM tCharacterLoading
	 WHERE nCharNo = @nCharNo 
	 ORDER BY nSeqNo DESC

	SET @nSeqNo       = @nSeqNo + 1

	INSERT INTO tCharacterLoading( nCharNo, nSeqNo, dWorldTime, sLoadingMapName, nFailType, nLoadingFailCnt )
	VALUES (@nCharNo, @nSeqNo, @dWorldTime, @sLoadingMapName, @nFailType, @nLoadingFailCnt)
	IF @@ERROR <> 0 BEGIN
		-- 파일 에러처리
		SET @nRet = -1001
		ROLLBACK TRAN
		RETURN
	END

	SET @nRet = 0
	COMMIT TRAN
END

0ᴁᦢA
GO

-- ========================================================================
-- Procedure: usp_Character_GetDiffFriendPoint
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Character_GetDiffFriendPoint]
	@nCharNo int
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	SET @nRet = 0
	
	IF EXISTS (SELECT nCharNo FROM tCharacter WHERE nCharNo = @nCharNo) BEGIN
		SELECT nFriendPoint - ISNULL(nLastFriendPoint, 0) AS Diffpoint
		  FROM tCharacter
		 WHERE @nCharNo = nCharNo
	END
	ELSE BEGIN
		SET @nRet = -1111
	END

	SET LOCK_TIMEOUT -1
    SET NOCOUNT OFF    
END

0
GO

-- ========================================================================
-- Procedure: usp_Character_GetFriendpoint
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Character_GetFriendpoint]
	@nCharNo int
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	SET @nRet = 0
	
	IF EXISTS (SELECT nCharNo FROM tCharacter WHERE nCharNo = @nCharNo) BEGIN
		SELECT ISNULL(nFriendPoint, 0) AS nFriendPoint
		  FROM tCharacter
		 WHERE @nCharNo = nCharNo
	END
	ELSE BEGIN
		SET @nRet = -1111
	END

	SET LOCK_TIMEOUT -1
    SET NOCOUNT OFF    
END

0
GO

-- ========================================================================
-- Procedure: usp_Character_SetFriendpoint
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Character_SetFriendpoint]
	@nCharNo		int
,	@nFriendPoint	int
,	@nRet			int = 0 output

AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	SET @nRet = 0
	
	IF EXISTS (SELECT nCharNo FROM tCharacter WHERE nCharNo = @nCharNo) BEGIN
		BEGIN TRAN

		UPDATE tCharacter SET nFriendPoint = @nFriendPoint WHERE nCharNo = @nCharNo
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = 1
			IF @@TRANCOUNT > 0 ROLLBACK TRAN
			RETURN
		END
		ELSE BEGIN
			IF @@TRANCOUNT > 0 COMMIT TRAN
			SET @nRet = 0
			RETURN
		END
	END
	ELSE BEGIN
		SET @nRet = -1111
	END
	
	SET LOCK_TIMEOUT -1
	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: usp_Character_chkKeyMap
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Character_chkKeyMap]
	@nCharNo      int
,	@nRet         int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000 
	
	SET @nRet = 0

	IF EXISTS(SELECT nCharNo FROM tCharacterKeyMap WHERE nCharNo = @nCharNo) BEGIN
		SELECT 1
	END
	ELSE BEGIN
		SELECT 0
	END
END

0㐁ᙁĪ
GO

-- ========================================================================
-- Procedure: usp_Character_getAutoScroll
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Character_getAutoScroll]
	@nCharNo		int
,	@nRet			int = 0 output

AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000
	
	SET @nRet = 0
	
	SELECT nAutoBuff, nAutoChargedbuff
	  FROM tCharacterSingleOption
	 WHERE nCharNo= @nCharNo
	
	SET LOCK_TIMEOUT -1
	SET NOCOUNT OFF
END

0ᤁꁑĪ
GO

-- ========================================================================
-- Procedure: usp_Character_getChangeIDList
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Character_getChangeIDList]
	@nUserNo int
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0
	
	SELECT nRowNo, nUserNo, nCharNo, bInit
	  FROM tCharacterChangeID
	 WHERE nUserNo = @nUserNo AND bFlag = 0
	   
	SET NOCOUNT OFF
END

0딁붣Ĭ
GO

-- ========================================================================
-- Procedure: usp_Character_getGameOption
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Character_getGameOption]
	@nCharNo int
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	SET @nRet = 0
	
	IF EXISTS(SELECT nCharNo FROM tCharacterGameOption WHERE nCharNo = @nCharNo) BEGIN
		SELECT T1.nOptionNo, ISNULL(T2.nValue, T1.nDefaultValue) AS nValue
		  FROM tGameOptionInit AS T1 LEFT OUTER JOIN (
				SELECT * FROM tCharacterGameOption WHERE nCharNo = @nCharNo) AS T2
					ON T1.nOptionNo = T2.nOptionNo
		 WHERE T1.bIsActive = 1
	END
	ELSE BEGIN
		SELECT nOptionNo, nDefaultValue AS nValue
		  FROM tGameOptionInit 
		 WHERE bIsActive = 1
	END
END

0렁엧.
GO

-- ========================================================================
-- Procedure: usp_Character_getItemActionCoolTime
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Character_getItemActionCoolTime]
	@nCharNo      int
,	@nRet         int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	SET @nRet = 0
	
	IF EXISTS(SELECT nCharNo FROM tItemActionCoolTime WHERE nCharNo = @nCharNo AND dEndTime < GETDATE()) BEGIN
		BEGIN TRAN
		
		DELETE FROM tItemActionCoolTime WHERE nCharNo = @nCharNo AND dEndTime < GETDATE()
		IF @@ERROR <> 0 OR @@ROWCOUNT = 0 BEGIN
			SET @nRet = -3294
			ROLLBACK TRAN
			RETURN
		END	
		
		SELECT nCharNo, nActionGroup, dEndTime FROM tItemActionCoolTime WHERE nCharNo = @nCharNo
		
		COMMIT TRAN
	END
	ELSE BEGIN
		SELECT nCharNo, nActionGroup, dEndTime FROM tItemActionCoolTime WHERE nCharNo = @nCharNo 
	END
	
	SET NOCOUNT OFF
END

0먋/
GO

-- ========================================================================
-- Procedure: usp_Character_getKeyMap
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Character_getKeyMap]
	@nCharNo      int
,	@nRet         int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000 
	
	SET @nRet = 0

	IF EXISTS(SELECT nCharNo FROM tCharacterKeyMap WHERE nCharNo = @nCharNo) BEGIN
		SELECT nFunctionNo, nExtendedKey, nASCIICode
		  FROM tCharacterKeyMap 
		 WHERE nCharNo = @nCharNo 
	END
	ELSE BEGIN
		SET @nRet = -2725
		RETURN
	END
END

0漁太İ
GO

-- ========================================================================
-- Procedure: usp_Character_getShortCutList
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Character_getShortCutList]
	@nCharNo int,
	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	SET @nRet = 0
	
	IF EXISTS(SELECT nCharNo FROM tCharacter WHERE nCharNo = @nCharNo AND bDeleted = 0) BEGIN
		SELECT nSlotNo, nCodeNo, nValue
		  FROM tCharacterShortCut
		 WHERE nCharNo = @nCharNo AND nCodeNo > 0
	END
	ELSE BEGIN
		SET @nRet = -2721
	END
	
	SET NOCOUNT OFF
END

0ꠁ䵎ı
GO

-- ========================================================================
-- Procedure: usp_Character_getSingleOption
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Character_getSingleOption]
	@nCharNo      int
,	@nRet         int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000 
	
	SET @nRet = 0

	IF NOT EXISTS(SELECT nCharNo FROM tCharacterSingleOption WHERE nCharNo = @nCharNo) BEGIN
		INSERT INTO tCharacterSingleOption (nCharNo)
		VALUES (@nCharNo)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -3192
			RETURN
		END
	END

	SELECT bMysteryBolt
		,  bAutoPick 
		,  dRestExpLastExecuteTime
		,  nChatBan
	  FROM tCharacterSingleOption 
	 WHERE nCharNo = @nCharNo 		
END

0䅲Ĳ
GO

-- ========================================================================
-- Procedure: usp_Character_initKeyMap
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Character_initKeyMap]
	@nCharNo      int
,	@nRet         int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000 
	
	SET @nRet = 0

	IF EXISTS(SELECT nCharNo FROM tCharacterKeyMap WHERE nCharNo = @nCharNo) BEGIN
		BEGIN TRAN
		
		DELETE FROM tCharacterKeyMap WHERE nCharNo = @nCharNo
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -2723
			ROLLBACK TRAN
			RETURN
		END
		
		INSERT INTO tCharacterKeyMap (nCharNo, nFunctionNo, sys_create_date, nExtendedKey, nASCIICode)
		SELECT @nCharNo, nFunctionNo, GETDATE(), nExtendedKey, nASCIICode
		  FROM tKeyMapInit
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -2724
			ROLLBACK TRAN
			RETURN
		END
		
		COMMIT TRAN
	END
	ELSE BEGIN
		INSERT INTO tCharacterKeyMap (nCharNo, nFunctionNo, sys_create_date, nExtendedKey, nASCIICode)
		SELECT @nCharNo, nFunctionNo, GETDATE(), nExtendedKey, nASCIICode
		  FROM tKeyMapInit
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -2724
			RETURN
		END
	END
END

0영쮂2
GO

-- ========================================================================
-- Procedure: usp_Character_setAutoPick
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Character_setAutoPick]
	@nCharNo   int
,	@bAutoPick bit
,	@nRet      int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000 
	
	SET @nRet = 0

	IF NOT EXISTS(SELECT nCharNo FROM tCharacterSingleOption WHERE nCharNo = @nCharNo) BEGIN
		INSERT INTO tCharacterSingleOption (nCharNo, bAutoPick)
		VALUES (@nCharNo, @bAutoPick)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -3192
			RETURN
		END		
	END
	ELSE BEGIN
		UPDATE tCharacterSingleOption 
		   SET bAutoPick = @bAutoPick 
		 WHERE nCharNo = @nCharNo
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -3193
			RETURN
		END		
	END
END

0！뾦3
GO

-- ========================================================================
-- Procedure: usp_Character_setBuffScroll
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Character_setBuffScroll]
	@nCharNo		int
,	@nAutoBuff		tinyint
,	@nRet			int = 0 output

AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	DECLARE @nAutoChargedBuff AS tinyint
	
	SET @nRet = 0
	
	BEGIN TRAN
	
	IF EXISTS (SELECT nCharNo FROM tCharacterSingleOption WHERE nCharNo = @nCharNo) BEGIN
		UPDATE tCharacterSingleOption SET nAutoBuff = @nAutoBuff WHERE nCharNo = @nCharNo
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = 1
			IF @@TRANCOUNT > 0 ROLLBACK TRAN
			RETURN
		END
	END
	ELSE BEGIN
		SET @nAutoChargedBuff = 0

		INSERT INTO tCharacterSingleOption (nCharNo, nAutoBuff, nAutoChargedBuff)
		VALUES (@nCharNo, @nAutoBuff, @nAutoChargedBuff)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = 1
			IF @@TRANCOUNT > 0 ROLLBACK TRAN
			RETURN
		END
	END
		
	COMMIT TRAN
	
	SET LOCK_TIMEOUT -1
	SET NOCOUNT OFF
END

0㠁돋4
GO

-- ========================================================================
-- Procedure: usp_Character_setChangeID
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Character_setChangeID]
	@nRowNo  int
,	@nCharNo int
,	@sID     nvarchar(16)
,	@bInit   bit
,	@nUserNo int
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0
	
	DECLARE @sOldID AS nvarchar(50)
	DECLARE @nDBUserNo AS int
	DECLARE @nGrade AS tinyint
	DECLARE @nMyGuildNo AS int
	
	SELECT @sOldID = T1.sID, @nDBUserNo = T1.nUserNo, @nGrade = nGrade
	  FROM tCharacter T1 LEFT OUTER JOIN tGuildMember T2 ON T1.nCharNo = T2.nCharNo
	 WHERE T1.nCharNo = @nCharNo AND T1.bDeleted = 0
	   AND bDeleted = 0
	IF @sOldID = @sID BEGIN
		SET @nRet = -3039
		RETURN
	END	
	IF @nDBUserNo IS NULL BEGIN
		SET @nRet = -3023
		RETURN
	END
	ELSE BEGIN
		-- mismatch userNo
		IF @nDBUserNo <> @nUserNo BEGIN
			SET @nRet = -3024
			RETURN
		END
		
		-- impossible guildmaster
		IF @nGrade IS NOT NULL AND @nGrade = 0 AND @bInit = 1 BEGIN
			SET @nRet = -3025
			RETURN
		END

		IF NOT EXISTS(SELECT nCharNo FROM tCharacterChangeID WHERE bFlag = 0 AND nCharNo = @nCharNo) BEGIN
			SET @nRet = -3024
			RETURN
		END

		BEGIN TRAN
		
		-- change ID
		IF EXISTS(SELECT nCharNo FROM tCharacter WHERE sID = @sID) BEGIN		
			SET @nRet = -3039
			ROLLBACK TRAN
			RETURN
		END
		ELSE BEGIN
			UPDATE tCharacter
			   SET sID = @sID
			 WHERE nCharNo = @nCharNo
			IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
				SET @nRet = -3026
				ROLLBACK TRAN
				RETURN
			END		
		END
		
		-- init Community
		IF @bInit = 1 BEGIN
			-- delete Friend
			IF EXISTS(SELECT nCharNo FROM tFriend WHERE nCharNo = @nCharNo) BEGIN
				DELETE FROM tFriend WHERE nCharNo = @nCharNo
				IF @@ERROR <> 0 OR @@ROWCOUNT < 1 BEGIN
					SET @nRet = -3027
					ROLLBACK TRAN
					RETURN
				END
			END
			IF EXISTS(SELECT nCharNo FROM tFriend WHERE nFriendCharNo = @nCharNo) BEGIN
				DELETE FROM tFriend WHERE nFriendCharNo = @nCharNo
				IF @@ERROR <> 0 OR @@ROWCOUNT < 1 BEGIN
					SET @nRet = -3027
					ROLLBACK TRAN
					RETURN
				END				
			END
			
			-- delete Marriage
			IF EXISTS(SELECT nCharNo FROM tMarriage WHERE nCharNo = @nCharNo) BEGIN
				DELETE FROM tMarriage WHERE nCharNo = @nCharNo
				IF @@ERROR <> 0 OR @@ROWCOUNT < 1 BEGIN
					SET @nRet = -3028
					ROLLBACK TRAN
					RETURN
				END
			END
			IF EXISTS(SELECT nCharNo FROM tMarriage WHERE nPartnerCharNo = @nCharNo) BEGIN				
				DELETE FROM tMarriage WHERE nPartnerCharNo = @nCharNo
				IF @@ERROR <> 0 OR @@ROWCOUNT < 1 BEGIN
					SET @nRet = -3028
					ROLLBACK TRAN
					RETURN
				END				
			END			

			-- delete HolyPromise me
			IF EXISTS(SELECT nCharNo FROM tHolyPromise WHERE nCharNo = @nCharNo) BEGIN
				DELETE FROM tHolyPromise WHERE nCharNo = @nCharNo
				IF @@ERROR <> 0 OR @@ROWCOUNT < 1 BEGIN
					SET @nRet = -3029
					ROLLBACK TRAN
					RETURN
				END
			END
GO

-- ========================================================================
-- Procedure: usp_Character_setChangeIDItemUsed
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Character_setChangeIDItemUsed]
	@nCharNo int
,	@bInit   bit
,	@nUserNo int
,	@nItemKey bigint
,	@nItemOption smallint
,	@nItemLot smallint
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	DECLARE @nGrade AS tinyint
	SET @nRet = 0

	-- 아이템 소유자 확인
	EXEC p_Item_IsValidOwner @nItemKey, @nCharNo, @nRet OUTPUT
	IF @nRet <> 0 BEGIN
		SET @nRet = -3036
		RETURN
	END

	-- impossible guildmaster
	SELECT @nGrade = nGrade
	  FROM tCharacter T1 LEFT OUTER JOIN tGuildMember T2 ON T1.nCharNo = T2.nCharNo
	 WHERE T1.nCharNo = @nCharNo AND T1.bDeleted = 0
	   AND bDeleted = 0
	IF @nGrade IS NOT NULL AND @nGrade = 0 AND @bInit = 1 BEGIN
		SET @nRet = -3025
		RETURN
	END	   
	   
	IF EXISTS(SELECT nRowNo FROM tCharacterChangeID WHERE nCharNo = @nCharNo AND bFlag = 0) BEGIN
		SET @nRet = -3022
		RETURN
	END
	ELSE IF EXISTS(SELECT nCharNo FROM tCharacter WHERE nUserNo = @nUserNo AND bDeleted = 0) BEGIN
		BEGIN TRAN
		
		-- 아이템 사용 기록
		INSERT INTO tCharacterChangeID (nUserNo, nCharNo, bInit)
		VALUES (@nUserNo, @nCharNo, @bInit)
		IF @@ERROR <> 0 OR @@IDENTITY = 0 BEGIN
			SET @nRet = -3092
			ROLLBACK TRAN
			RETURN
		END	
		
		-- 아이템 삭제 (사용후 Lot = 0 인 경우)
		IF @nItemLot = 0 BEGIN
			EXEC p_Item_Delete @nItemKey, @nRet OUTPUT
			IF @nRet <> 0 BEGIN
				ROLLBACK TRAN
				SET @nRet = -3037
				RETURN
			END
		END
		-- 아이템 차감하기 (사용후 Lot > 0 인 경우)
		ELSE BEGIN
			EXEC p_Item_SetOption @nItemKey, @nItemOption, @nItemLot, @nRet OUTPUT
			IF @nRet <> 0 BEGIN
				ROLLBACK TRAN
				SET @nRet = -3038
				RETURN
			END
		END
		
		COMMIT TRAN		
	END
	ELSE BEGIN
		SET @nRet = -3023
		RETURN	
	END
	
	SET NOCOUNT OFF
END

0뤁嬭7
GO

-- ========================================================================
-- Procedure: usp_Character_setChargedbuffScroll
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Character_setChargedbuffScroll]
	@nCharNo			int
,	@nAutoChargedBuff	tinyint
,	@nRet				int = 0 output

AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000
	
	DECLARE @nAutoBuff AS tinyint
	
	SET @nRet = 0
	
	BEGIN TRAN
	
	IF EXISTS (SELECT nCharNo FROM tCharacterSingleOption WHERE nCharNo = @nCharNo) BEGIN
		UPDATE tCharacterSingleOption SET nAutoChargedBuff = @nAutoChargedBuff WHERE nCharNo = @nCharNo
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = 1
			IF @@TRANCOUNT > 0 ROLLBACK TRAN
			RETURN
		END
	END
	ELSE BEGIN
		SET @nAutoBuff = 0
	
		INSERT INTO tCharacterSingleOption (nCharNo, nAutoBuff, nAutoChargedBuff)
		VALUES (@nCharNo, @nAutoBuff, @nAutoChargedBuff)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = 1
			IF @@TRANCOUNT > 0 ROLLBACK TRAN
			RETURN
		END
	END
	
	
	COMMIT TRAN
	
	SET LOCK_TIMEOUT -1
	SET NOCOUNT OFF
END

0ᰁ葜8
GO

-- ========================================================================
-- Procedure: usp_Character_setChatBan
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Character_setChatBan]
	@sID       nvarchar(40)
,	@nChatBan  int
,	@nRet      int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000 
	
	DECLARE @nCharNo AS int
	
	SET @nRet = 0
	
	SELECT @nCharNo = nCharNo FROM tCharacter WHERE sID = @sID AND bDeleted = 0
	IF @nCharNo IS NULL BEGIN
		SET @nRet = -3121
	END
	ELSE BEGIN
		IF NOT EXISTS(SELECT nCharNo FROM tCharacterSingleOption WHERE nCharNo = @nCharNo) BEGIN
			INSERT INTO tCharacterSingleOption (nCharNo, nChatBan)
			VALUES (@nCharNo, @nChatBan)
			IF @@ERROR <> 0 BEGIN
				SET @nRet = -3192
				RETURN
			END		
		END
		ELSE BEGIN
			UPDATE tCharacterSingleOption 
			   SET nChatBan = @nChatBan 
			 WHERE nCharNo = @nCharNo
			IF @@ERROR <> 0 BEGIN
				SET @nRet = -3193
				RETURN
			END		
		END
	END
END

0⬁䍶9
GO

-- ========================================================================
-- Procedure: usp_Character_setGameOption
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Character_setGameOption]
	@nCharNo   int
,	@nOptionNo smallint	
,	@nValue    tinyint
,	@nRet      int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	SET @nRet = 0
	
	IF EXISTS(SELECT nCharNo FROM tCharacterGameOption WHERE nCharNo = @nCharNo AND nOptionNo = @nOptionNo) BEGIN
		UPDATE tCharacterGameOption 
		   SET nValue = @nValue
		 WHERE nCharNo = @nCharNo
		   AND nOptionNo = @nOptionNo 
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = -2792
			RETURN
		END		   
	END
	ELSE BEGIN
		INSERT INTO tCharacterGameOption (nCharNo, nOptionNo, nValue)
		VALUES (@nCharNo, @nOptionNo, @nValue)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -2792
			RETURN
		END		   		
	END		
END

0鴁⮾;
GO

-- ========================================================================
-- Procedure: usp_Character_setGameOptionInit
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Character_setGameOptionInit]
	@nCharNo   int
,	@nRet      int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	SET @nRet = 0
	
	DELETE FROM tCharacterGameOption WHERE nCharNo = @nCharNo
END

0윁惈;
GO

-- ========================================================================
-- Procedure: usp_Character_setItemActionCoolTime
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Character_setItemActionCoolTime]
	@nCharNo      int
,	@nActionGroup smallint
,	@dEndTime     datetime
,	@nRet         int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	SET @nRet = 0

	IF EXISTS(SELECT nCharNo FROM tItemActionCoolTime WHERE nCharNo = @nCharNo AND nActionGroup = @nActionGroup) BEGIN
		BEGIN TRAN
		
		UPDATE tItemActionCoolTime
		   SET dEndTime = @dEndTime 
		 WHERE nCharNo = @nCharNo 
		   AND nActionGroup = @nActionGroup
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = -3291
			ROLLBACK TRAN
			RETURN
		END
		
		COMMIT TRAN
	END
	ELSE BEGIN
		INSERT INTO tItemActionCoolTime (nCharNo, nActionGroup, dEndTime)
		VALUES (@nCharNo, @nActionGroup, @dEndTime)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -3291
			RETURN
		END
	END	
	
	SET NOCOUNT OFF
END

0༁ᐇ=
GO

-- ========================================================================
-- Procedure: usp_Character_setKeyMap
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Character_setKeyMap]
	@nCharNo      int
,	@nFunctionNo  smallint
,	@nExtendedKey tinyint
,	@nASCIICode   tinyint
,	@nRet         int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000 
	
	SET @nRet = 0

	BEGIN TRAN
	
	IF EXISTS(SELECT nCharNo FROM tCharacterKeyMap WHERE nCharNo = @nCharNo AND nFunctionNo = @nFunctionNo) BEGIN
		UPDATE tCharacterKeyMap 
		   SET nExtendedKey = @nExtendedKey 
		     , nASCIICode = @nASCIICode 
		 WHERE nCharNo = @nCharNo 
		   AND nFunctionNo = @nFunctionNo 
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = -2793
			ROLLBACK TRAN
			RETURN
		END
	END
	ELSE BEGIN
		INSERT INTO tCharacterKeyMap (nCharNo, nFunctionNo, nExtendedKey, nASCIICode)
		VALUES (@nCharNo, @nFunctionNo, @nExtendedKey, @nASCIICode)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -2792
			ROLLBACK TRAN
			RETURN
		END
	END
	
	COMMIT TRAN
END

0㤁䤑=
GO

-- ========================================================================
-- Procedure: usp_Character_setMysteryBolt
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Character_setMysteryBolt]
	@nCharNo      int
,	@bMysteryBolt bit
,	@nRet         int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000 
	
	SET @nRet = 0

	IF NOT EXISTS(SELECT nCharNo FROM tCharacterSingleOption WHERE nCharNo = @nCharNo) BEGIN
		INSERT INTO tCharacterSingleOption (nCharNo, bMysteryBolt)
		VALUES (@nCharNo, @bMysteryBolt)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -3192
			RETURN
		END		
	END
	ELSE BEGIN
		UPDATE tCharacterSingleOption 
		   SET bMysteryBolt = @bMysteryBolt 
		 WHERE nCharNo = @nCharNo
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -3193
			RETURN
		END		
	END
END

0爁㴵>
GO

-- ========================================================================
-- Procedure: usp_Character_setRestExpExecTime
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Character_setRestExpExecTime]
	@nCharNo   int
,	@nRet      int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000 
	
	SET @nRet = 0

	IF NOT EXISTS(SELECT nCharNo FROM tCharacterSingleOption WHERE nCharNo = @nCharNo) BEGIN
		INSERT INTO tCharacterSingleOption (nCharNo)
		VALUES (@nCharNo)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -3192
			RETURN
		END		
	END
	ELSE BEGIN
		UPDATE tCharacterSingleOption 
		   SET dRestExpLastExecuteTime = GETDATE() 
		 WHERE nCharNo = @nCharNo
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -3193
			RETURN
		END		
	END
END

0╽@
GO

-- ========================================================================
-- Procedure: usp_Character_setShortCut
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Character_setShortCut]
	@nCharNo int,
	@nSlotNo tinyint,
	@nCodeNo smallint,
	@nValue  int,
	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	SET @nRet = 0
	
	DECLARE @nCodeType smallint SET @nCodeType = 0
	DECLARE @nDBCharNo int
	DECLARE @nPrevCodeNo smallint SET @nPrevCodeNo = 0
	DECLARE @nPrevValue int SET @nPrevValue = 0

	SELECT @nDBCharNo = nCharNo, @nPrevCodeNo = nCodeNo, @nPrevValue = nValue
	  FROM tCharacterShortCut
	 WHERE nCharNo = @nCharNo 
	   AND nSlotNo = @nSlotNo
	IF @nDBCharNo IS NULL BEGIN
		INSERT INTO tCharacterShortCut (nCharNo, nSlotNo, nCodeType, nCodeNo, nValue)
		VALUES (@nCharNo, @nSlotNo, 1, @nCodeNo, @nValue)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -2791
		END
	END
	ELSE IF @nPrevCodeNo = @nCodeNo AND @nPrevValue = @nValue BEGIN
		SET @nRet = -2722
	END
	ELSE BEGIN
		UPDATE tCharacterShortCut 
		   SET nCodeNo = @nCodeNo
		     , nValue  = @nValue
		     , sys_update_date = GETDATE()
		 WHERE nCharNo = @nCharNo 
		   AND nSlotNo = @nSlotNo
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = -2793
		END
	END

	SET NOCOUNT OFF
END

0@
GO

-- ========================================================================
-- Procedure: usp_Character_setStatPoint
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Character_setStatPoint]
	@nCharNo       int
,	@nMaxStatPoint int
,	@nStrength     tinyint
,	@nConstitute   tinyint
,	@nDexterity    tinyint
,	@nIntelligence tinyint
,	@nMentalPower  tinyint
,	@nItemID       int
,	@nOptionType   smallint
,	@nDeleteCount  int
,	@nRet          int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000
	
	DECLARE @nTotalStat   int
	DECLARE @nDeleteCheck int
	DECLARE @nStr int, @nCon int, @nDex int, @nInt int, @nMen int
	DECLARE @nTotalItemCount int, @nItemKey bigint, @nOptionData bigint, @nItemDeleteRemain int, @nRemainCount int
	
	SET @nRet = 0
	SET @nTotalStat = 0
	SET @nDeleteCheck = 0

	-- invalid nCharNo	
	IF @nCharNo <= 0 BEGIN
		SET @nRet = -1011
		RETURN
	END
	
	-- invalid StatPoint
	IF @nStrength < 0 OR @nConstitute < 0 OR @nDexterity < 0 OR @nIntelligence < 0 OR @nMentalPower < 0 BEGIN
		SET @nRet = -1012
		RETURN
	END
	
	-- have many StatPoint
	SET @nTotalStat = @nStrength + @nConstitute + @nDexterity + @nIntelligence + @nMentalPower
	IF @nMaxStatPoint < @nTotalStat BEGIN
		SET @nRet = -1013
		RETURN
	END
	
	IF EXISTS(SELECT nCharNo FROM tCharacter WHERE nCharNo = @nCharNo AND bDeleted = 0) BEGIN
		SELECT @nStr = nStrength
			,  @nCon = nConstitute 
			,  @nDex = nDexterity 
			,  @nInt = nIntelligence 
			,  @nMen = nMentalPower 
		  FROM tCharacter
		 WHERE nCharNo = @nCharNo
		IF @nStr > @nStrength SET @nDeleteCheck = @nDeleteCheck + (@nStr - @nStrength)
		IF @nCon > @nConstitute SET @nDeleteCheck = @nDeleteCheck + (@nCon - @nConstitute)
		IF @nDex > @nDexterity SET @nDeleteCheck = @nDeleteCheck + (@nDex - @nDexterity)
		IF @nInt > @nIntelligence SET @nDeleteCheck = @nDeleteCheck + (@nInt - @nIntelligence)
		IF @nMen > @nMentalPower SET @nDeleteCheck = @nDeleteCheck + (@nMen - @nMentalPower)
		IF @nDeleteCount <> @nDeleteCheck BEGIN
			SET @nRet = -1021
			RETURN
		END
		
		BEGIN TRAN
		
		UPDATE tCharacter
		   SET nStrength = @nStrength 
		   ,   nConstitute = @nConstitute 
		   ,   nDexterity = @nDexterity 
		   ,   nIntelligence = @nIntelligence 
		   ,   nMentalPower = @nMentalPower 
		 WHERE nCharNo = @nCharNo 
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = -1093
			ROLLBACK TRAN
			RETURN
		END		 

		IF @nDeleteCount > 0  BEGIN
			IF @nDeleteCount = 1 BEGIN
				SELECT TOP 1 @nItemKey = T1.nItemKey, @nOptionData = T2.nOptionData
				  FROM tItem T1, tItemOptions T2
				 WHERE T1.nItemKey = T2.nItemKey 
				   AND T1.nOwner = @nCharNo
				   AND T1.nStorageType = 9
				   AND T1.nItemID = @nItemID 
				 ORDER BY nStorage ASC
				IF @nOptionData > 1 BEGIN
					SET @nOptionData = @nOptionData - 1
					EXEC p_Item_SetOption @nItemKey, @nOptionType, @nOptionData, @nRet output
					IF @nRet <> 0 BEGIN
						SET @nRet = -1022
						ROLLBACK TRAN
						RETURN						
					END					
				END
				ELSE BEGIN
					EXEC p_Item_Delete @nItemKey, @n
GO

-- ========================================================================
-- Procedure: usp_ChargeItem_Draw
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_ChargeItem_Draw]
	@ROW_NO INT,
	@USER_NO INT,
	@RET	INT = 0 OUTPUT
AS

BEGIN

	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	EXEC ACCOUNT.Account.dbo.usp_Charge_ItemDraw @ROW_NO, @USER_NO, @RET OUTPUT
	
	SET NOCOUNT OFF
	SET LOCK_TIMEOUT -1

END

0적C
GO

-- ========================================================================
-- Procedure: usp_ChargeItem_GetList
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_ChargeItem_GetList]
	@USER_NO INT
AS

BEGIN

	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	EXEC Account.dbo.usp_Charge_ItemSelectList @USER_NO

	SET NOCOUNT OFF
	SET LOCK_TIMEOUT -1

END
GO

-- ========================================================================
-- Procedure: usp_ChargeItem_UseSkillPowerInitDel
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_ChargeItem_UseSkillPowerInitDel]
/*
Charge item use - clear character all skill power and delete used item
2006.6 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK
	        1 = Error clear character all skill power
	        2 = Error delete used charge item
*/

@nCharNo int,
@nDeleteItemKey bigint,
@nRet tinyint OUTPUT

AS

SET NOCOUNT ON

BEGIN TRAN

-- Clear character all skill power
UPDATE tSkill SET 
	nSkillPowerDemage	= 0,
	nSkillPowerSP		= 0,
	nSkillPowerKeepTime = 0,
	nSkillPowerCoolTime = 0
WHERE nCharNo = @nCharNo

IF @@ERROR <> 0
BEGIN
	ROLLBACK TRAN
	SET @nRet = 1
	RETURN
END

-- Delete used charge item
DECLARE @TempRet int
EXEC p_Item_Delete @nDeleteItemKey, @TempRet OUTPUT
IF @TempRet <> 0
BEGIN
	ROLLBACK TRAN
	SET @nRet = 2
	RETURN
END

-- Done
COMMIT TRAN

SET @nRet = 0
RETURN

0āD
GO

-- ========================================================================
-- Procedure: usp_ChargeItem_UseSkillPowerInitLotChange
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_ChargeItem_UseSkillPowerInitLotChange]
/*
Charge item use - clear character all skill power and use item lot change
2006.6 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK
	        1 = Error clear character all skill power
	        2 = Error used charge item change lot
*/

@nCharNo int,
@nChangeItemKey bigint,
@nChangeItemOptionType smallint,
@nChangeItemOptionData bigint,		
@nRet tinyint OUTPUT

AS

SET NOCOUNT ON

BEGIN TRAN

-- Clear character all skill power
UPDATE tSkill SET 
	nSkillPowerDemage	= 0,
	nSkillPowerSP		= 0,
	nSkillPowerKeepTime = 0,
	nSkillPowerCoolTime = 0
WHERE nCharNo = @nCharNo

IF @@ERROR <> 0
BEGIN
	ROLLBACK TRAN
	SET @nRet = 1
	RETURN
END

-- Used charge item change lot
DECLARE @TempRet int
EXEC p_Item_SetOption @nChangeItemKey, @nChangeItemOptionType, @nChangeItemOptionData, @TempRet OUTPUT
IF @TempRet <> 0
BEGIN
	ROLLBACK TRAN
	SET @nRet = 2
	RETURN
END

-- Done
COMMIT TRAN

SET @nRet = 0
RETURN

0猁퉻F
GO

-- ========================================================================
-- Procedure: usp_ChargeItem_UseStatPointDecDel
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_ChargeItem_UseStatPointDecDel]
/*
Charge item use - Stat point decrement and use item delete
2006.6 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK
	        1 = Error decrement character stat point.
	        2 = Error delete used charge item
*/

@nCharNo int,
@nStatPointType tinyint,
@nDeleteItemKey bigint,
@nRet tinyint OUTPUT

AS

SET NOCOUNT ON

BEGIN TRAN

-- CHECK Type and Decrement stat-point
DECLARE @nStatPoint tinyint

IF @nStatPointType = 0
BEGIN
	SELECT @nStatPoint = nStrength FROM tCharacter WHERE nCharNo = @nCharNo
	IF @nStatPoint < 1
	BEGIN
		SET @nRet = 1
		RETURN
	END
	
	UPDATE tCharacter	SET nStrength	= nStrength - 1					WHERE nCharNo = @nCharNo AND nStrength > 0
END

IF @nStatPointType = 1
BEGIN
	SELECT @nStatPoint = nConstitute FROM tCharacter WHERE nCharNo = @nCharNo
	IF @nStatPoint < 1
	BEGIN
		SET @nRet = 1
		RETURN
	END
	
	UPDATE tCharacter	SET nConstitute = nConstitute - 1				WHERE nCharNo = @nCharNo AND nConstitute > 0
END
	
IF @nStatPointType = 2
BEGIN
	SELECT @nStatPoint = nDexterity FROM tCharacter WHERE nCharNo = @nCharNo
	IF @nStatPoint < 1
	BEGIN
		SET @nRet = 1
		RETURN
	END
	
	UPDATE tCharacter	SET nDexterity = nDexterity - 1					WHERE nCharNo = @nCharNo AND nDexterity > 0
END

IF @nStatPointType = 3
BEGIN
	SELECT @nStatPoint = nIntelligence FROM tCharacter WHERE nCharNo = @nCharNo
	IF @nStatPoint < 1
	BEGIN
		SET @nRet = 1
		RETURN
	END
	
	UPDATE tCharacter	SET nIntelligence = nIntelligence - 1			WHERE nCharNo = @nCharNo AND nIntelligence > 0
END

IF @nStatPointType = 4 
BEGIN
	SELECT @nStatPoint = nMentalPower FROM tCharacter WHERE nCharNo = @nCharNo
	IF @nStatPoint < 1
	BEGIN
		SET @nRet = 1
		RETURN
	END
	
	UPDATE tCharacter	SET nMentalPower = nMentalPower - 1				WHERE nCharNo = @nCharNo AND nMentalPower > 0
END

IF @@ERROR <> 0
BEGIN
	ROLLBACK TRAN
	SET @nRet = 1
	RETURN
END

-- Delete used charge item
DECLARE @TempRet int
EXEC p_Item_Delete @nDeleteItemKey, @TempRet OUTPUT
IF @TempRet <> 0
BEGIN
	ROLLBACK TRAN
	SET @nRet = 2
	RETURN
END

-- Done
COMMIT TRAN

SET @nRet = 0
RETURN

0舁醕G
GO

-- ========================================================================
-- Procedure: usp_ChargeItem_UseStatPointDecLotChange
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_ChargeItem_UseStatPointDecLotChange]
/*
Charge item use - Stat point decrement and use item lot change
2006.6 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK
	        1 = Error decrement character stat point.
	        2 = Error used charge item change lot
*/

@nCharNo int,
@nStatPointType tinyint,
@nChangeItemKey bigint,
@nChangeItemOptionType smallint,
@nChangeItemOptionData bigint,		
@nRet tinyint OUTPUT

AS

SET NOCOUNT ON

BEGIN TRAN

-- CHECK Type and Decrement stat-point
DECLARE @nStatPoint tinyint

IF @nStatPointType = 0
BEGIN
	SELECT @nStatPoint = nStrength FROM tCharacter WHERE nCharNo = @nCharNo
	IF @nStatPoint < 1
	BEGIN
		SET @nRet = 1
		RETURN
	END
	
	UPDATE tCharacter	SET nStrength	= nStrength - 1					WHERE nCharNo = @nCharNo AND nStrength > 0
END

IF @nStatPointType = 1
BEGIN
	SELECT @nStatPoint = nConstitute FROM tCharacter WHERE nCharNo = @nCharNo
	IF @nStatPoint < 1
	BEGIN
		SET @nRet = 1
		RETURN
	END
	
	UPDATE tCharacter	SET nConstitute = nConstitute - 1				WHERE nCharNo = @nCharNo AND nConstitute > 0
END
	
IF @nStatPointType = 2
BEGIN
	SELECT @nStatPoint = nDexterity FROM tCharacter WHERE nCharNo = @nCharNo
	IF @nStatPoint < 1
	BEGIN
		SET @nRet = 1
		RETURN
	END
	
	UPDATE tCharacter	SET nDexterity = nDexterity - 1					WHERE nCharNo = @nCharNo AND nDexterity > 0
END

IF @nStatPointType = 3
BEGIN
	SELECT @nStatPoint = nIntelligence FROM tCharacter WHERE nCharNo = @nCharNo
	IF @nStatPoint < 1
	BEGIN
		SET @nRet = 1
		RETURN
	END
	
	UPDATE tCharacter	SET nIntelligence = nIntelligence - 1			WHERE nCharNo = @nCharNo AND nIntelligence > 0
END

IF @nStatPointType = 4 
BEGIN
	SELECT @nStatPoint = nMentalPower FROM tCharacter WHERE nCharNo = @nCharNo
	IF @nStatPoint < 1
	BEGIN
		SET @nRet = 1
		RETURN
	END
	
	UPDATE tCharacter	SET nMentalPower = nMentalPower - 1				WHERE nCharNo = @nCharNo AND nMentalPower > 0
END

IF @@ERROR <> 0
BEGIN
	ROLLBACK TRAN
	SET @nRet = 1
	RETURN
END

-- Used charge item change lot
DECLARE @TempRet int
EXEC p_Item_SetOption @nChangeItemKey, @nChangeItemOptionType, @nChangeItemOptionData, @TempRet OUTPUT
IF @TempRet <> 0
BEGIN
	ROLLBACK TRAN
	SET @nRet = 2
	RETURN
END

-- Done
COMMIT TRAN

SET @nRet = 0
RETURN

0각욟G
GO

-- ========================================================================
-- Procedure: usp_ChargeItem_UseStatPointInitDel
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_ChargeItem_UseStatPointInitDel]
/*
Charge item use - Stat point init and use item delete
2006.6 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK
	        1 = Error clear character stat point.
	        2 = Error delete used charge item
*/

@nCharNo int,
@nDeleteItemKey bigint,
@nRet tinyint OUTPUT

AS

SET NOCOUNT ON

BEGIN TRAN

-- Clear character stat point.
UPDATE tCharacter SET 
	nStrength = 0, 
	nConstitute = 0, 
	nDexterity = 0, 
	nIntelligence = 0, 
	nMentalPower = 0
WHERE nCharNo = @nCharNo

IF @@ERROR <> 0
BEGIN
	ROLLBACK TRAN
	SET @nRet = 1
	RETURN
END

-- Delete used charge item
DECLARE @TempRet int
EXEC p_Item_Delete @nDeleteItemKey, @TempRet OUTPUT
IF @TempRet <> 0
BEGIN
	ROLLBACK TRAN
	SET @nRet = 2
	RETURN
END

-- Done
COMMIT TRAN

SET @nRet = 0
RETURN

0뫃H
GO

-- ========================================================================
-- Procedure: usp_ChargeItem_UseStatPointInitLotChange
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_ChargeItem_UseStatPointInitLotChange]
/*
Charge item use - Stat point init and use item lot change
2006.6 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK
	        1 = Error clear character stat point.
	        2 = Error used charge item change lot
*/

@nCharNo int,
@nChangeItemKey bigint,
@nChangeItemOptionType smallint,
@nChangeItemOptionData bigint,		
@nRet tinyint OUTPUT

AS

SET NOCOUNT ON

BEGIN TRAN

-- Clear character stat point.
UPDATE tCharacter SET 
	nStrength = 0, 
	nConstitute = 0, 
	nDexterity = 0, 
	nIntelligence = 0, 
	nMentalPower = 0
WHERE nCharNo = @nCharNo

IF @@ERROR <> 0
BEGIN
	ROLLBACK TRAN
	SET @nRet = 1
	RETURN
END

-- Used charge item change lot
DECLARE @TempRet int
EXEC p_Item_SetOption @nChangeItemKey, @nChangeItemOptionType, @nChangeItemOptionData, @TempRet OUTPUT
IF @TempRet <> 0
BEGIN
	ROLLBACK TRAN
	SET @nRet = 2
	RETURN
END

-- Done
COMMIT TRAN

SET @nRet = 0
RETURN

0ⴁ渂J
GO

-- ========================================================================
-- Procedure: usp_ChatBlock_delBlock
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_ChatBlock_delBlock]
	@nUserNo int
,	@nCharNo int
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	IF ISNULL(@nUserNo, 0) <= 0 BEGIN
		SET @nRet = -2311
		RETURN
	END

	IF ISNULL(@nCharNo, 0) <= 0 BEGIN
		SET @nRet = -2312
		RETURN
	END

	IF NOT EXISTS(SELECT nUserNo FROM tChatBlock WHERE nUserNo = @nUserNo AND nCharNo = @nCharNo) BEGIN
		SET @nRet = -2323
		RETURN
	END
	ELSE BEGIN
		DELETE FROM tChatBlock WHERE nUserNo = @nUserNo AND nCharNo = @nCharNo
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -2394
		END
	END
END

0送霰K
GO

-- ========================================================================
-- Procedure: usp_ChatBlock_delBlockALL
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_ChatBlock_delBlockALL]
	@nUserNo int
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	IF ISNULL(@nUserNo, 0) <= 0 BEGIN
		SET @nRet = -2311
		RETURN
	END

	IF NOT EXISTS(SELECT nUserNo FROM tChatBlock WHERE nUserNo = @nUserNo) BEGIN
		SET @nRet = -2321
		RETURN
	END
	ELSE BEGIN
		DELETE FROM tChatBlock WHERE nUserNo = @nUserNo
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -2394
		END
	END
END

0줁譔L
GO

-- ========================================================================
-- Procedure: usp_ChatBlock_getBlockList
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_ChatBlock_getBlockList]
	@nUserNo int
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	IF ISNULL(@nUserNo, 0) <= 0 BEGIN
		SET @nRet = -2311
		RETURN
	END

	SELECT T1.nCharNo, T2.sID, T2.nAdminLevel, T2.bDeleted
	  FROM tChatBlock T1, tCharacter T2
	 WHERE T1.nUserNo = @nUserNo AND T1.nCharNo = T2.nCharNo
END

0ȁ罹M
GO

-- ========================================================================
-- Procedure: usp_ChatBlock_setBlock
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_ChatBlock_setBlock]
	@nUserNo int
,	@sID     nvarchar(40) = ''
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	DECLARE @cnt AS int
	DECLARE @maxCnt AS int
	DECLARE @nCharNo AS int

	SET @nRet = 0
	SET @nCharNo = 0
	SET @maxCnt = 100 -- 최대등록 차단수(차후 변경시 DB 저장프로시저만 수정하면 반영됨)

	IF ISNULL(@nUserNo, 0) <= 0 BEGIN
		SET @nRet = -2311
		RETURN
	END

	IF LEN(ISNULL(@sID, '')) = 0 BEGIN
		SET @nRet = -2312
		RETURN
	END

	SELECT @nCharNo = nCharNo FROM tCharacter WHERE sID = @sID AND bDeleted = 0
	IF @nCharNo <= 0 BEGIN
		SET @nRet = -2323
		RETURN
	END

	IF EXISTS(SELECT nUserNo FROM tChatBlock WHERE nUserNo = @nUserNo AND nCharNo = @nCharNo) BEGIN
		SET @nRet = -2321
		RETURN
	END

	SELECT @cnt = ISNULL(COUNT(*), 0) FROM tChatBlock WHERE nUserNo = @nUserNo
	IF @cnt > @maxCnt BEGIN
		SET @nRet = -2322
		RETURN
	END		
	ELSE BEGIN
		INSERT INTO tChatBlock (nUserNo, nCharNo)
		VALUES (@nUserNo, @nCharNo)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -2392
			RETURN
		END

		SELECT T1.nCharNo, T2.sID, T2.nAdminLevel, T2.bDeleted
		  FROM tChatBlock T1, tCharacter T2
		 WHERE T1.nUserNo = @nUserNo 
		   AND T1.nCharNo = @nCharNo
		   AND T1.nCharNo = T2.nCharNo
	END
END

0ᄁ㺓N
GO

-- ========================================================================
-- Procedure: usp_ChatColor_get
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_ChatColor_get]
	@nCharNo         int
,	@nRet            int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	SET @nRet = 0
	
	IF EXISTS(SELECT nCharNo FROM tChatColor WHERE nCharNo = @nCharNo) BEGIN
		SELECT nFontColorID, nBalloonColorID 
		  FROM tChatColor 
		 WHERE nCharNo = @nCharNo
	END
	ELSE BEGIN
		SELECT 0 AS nFontColorID, 0 AS nBalloonColorID
	END
END

0㬁玝N
GO

-- ========================================================================
-- Procedure: usp_ChatColor_setBalloon
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_ChatColor_setBalloon]
	@nCharNo         int
,	@nBalloonColorID tinyint
,	@nRet            int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	DECLARE @nFontColorID AS tinyint

	SET @nRet = 0
	
	IF EXISTS(SELECT nCharNo FROM tChatColor WHERE nCharNo = @nCharNo) BEGIN
		BEGIN TRAN

		UPDATE tChatColor 
		   SET nBalloonColorID = @nBalloonColorID 
		 WHERE nCharNo = @nCharNo 
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			ROLLBACK TRAN
			SET @nRet = -1
		END
		ELSE BEGIN
			COMMIT TRAN
		END	
	END
	ELSE BEGIN
		SET @nFontColorID = 0
			
		INSERT INTO tChatColor (nCharNo, nFontColorID, nBalloonColorID)
		VALUES (@nCharNo, @nFontColorID, @nBalloonColorID)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -2
		END
	END
END

0괁寥P
GO

-- ========================================================================
-- Procedure: usp_ChatColor_setFont
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_ChatColor_setFont]
	@nCharNo      int
,	@nFontColorID tinyint
,	@nRet         int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	DECLARE @nBalloonColorID AS tinyint

	SET @nRet = 0
	
	IF EXISTS(SELECT nCharNo FROM tChatColor WHERE nCharNo = @nCharNo) BEGIN
		BEGIN TRAN

		UPDATE tChatColor 
		   SET nFontColorID = @nFontColorID 
		 WHERE nCharNo = @nCharNo 
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			ROLLBACK TRAN
			SET @nRet = -1
		END
		ELSE BEGIN
			COMMIT TRAN
		END	
	END
	ELSE BEGIN
		SET @nBalloonColorID = 0
			
		INSERT INTO tChatColor (nCharNo, nFontColorID, nBalloonColorID)
		VALUES (@nCharNo, @nFontColorID, @nBalloonColorID)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -2
		END
	END
END

0倉Q
GO

-- ========================================================================
-- Procedure: usp_Collection_getCardList
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Collection_getCardList]
	@nCharNo      int
,	@nRet         int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0
	
	IF ISNULL(@nCharNo, 0) = 0 BEGIN
		SET @nRet = -1111
		RETURN
	END

	SELECT nCardID, nSerialNo, nStarNumber, dDate 
	  FROM tCollectionCard 
	 WHERE nCharNo = @nCharNo 	
END

0ἁ䐮R
GO

-- ========================================================================
-- Procedure: usp_Collection_getFavoritesList
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Collection_getFavoritesList]
	@nCharNo      int
,	@nRet         int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0
	
	IF ISNULL(@nCharNo, 0) = 0 BEGIN
		SET @nRet = -1111
		RETURN
	END

	SELECT T1.nSlotNo, T1.nCardID, T2.nSerialNo, T2.nStarNumber 
	  FROM tCollectionFavorites T1, tCollectionCard T2 
	 WHERE T1.nCharNo = T2.nCharNo
	   AND T1.nCardID = T2.nCardID 
	   AND T1.nCharNo = @nCharNo 
END

0堁㡒S
GO

-- ========================================================================
-- Procedure: usp_Collection_getRewardList
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Collection_getRewardList]
	@nCharNo      int
,	@nRet         int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0
	
	IF ISNULL(@nCharNo, 0) = 0 BEGIN
		SET @nRet = -1111
		RETURN
	END

	SELECT nReward 
	  FROM tCollectionReward 
	 WHERE nCharNo = @nCharNo 
END

0༁쮕Ŕ
GO

-- ========================================================================
-- Procedure: usp_Collection_getSerialList
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Collection_getSerialList]
	@nRet         int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0
	
	SELECT nCardID, nCount
	  FROM tCollectionSerialCount 
END

0鄁ⱶT
GO

-- ========================================================================
-- Procedure: usp_Collection_setCard
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Collection_setCard]
	@nCharNo      int
,	@nCardID      int
,	@nSerialNo    bigint
,	@nStarNumber  tinyint
,	@nItemKey     bigint
,	@nRet         int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	SET @nRet = 0
	
	IF ISNULL(@nCharNo, 0) = 0 BEGIN
		SET @nRet = -1111
		RETURN
	END
	IF ISNULL(@nCardID, 0) = 0 BEGIN
		SET @nRet = -1111
		RETURN	
	END
	IF ISNULL(@nSerialNo, 0) = 0 BEGIN
		SET @nRet = -1111
		RETURN	
	END
	IF ISNULL(@nStarNumber, 0) = 0 BEGIN
		SET @nRet = -1111
		RETURN	
	END

	IF EXISTS(SELECT nItemKey FROM tItem WHERE nItemKey = @nItemKey) BEGIN
		BEGIN TRAN
		
		-- 1. 컬렉션 카드 등록
		IF EXISTS(SELECT nCharNo FROM tCollectionCard WHERE nCharNo = @nCharNo AND nCardID = @nCardID) BEGIN
			-- 1-1. 이미 존재하는 경우
			UPDATE tCollectionCard 
			   SET nSerialNo   = @nSerialNo 
				 , nStarNumber = @nStarNumber 
			 WHERE nCharNo = @nCharNo 
			   AND nCardID = @nCardID 
			IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
				SET @nRet = -9999
				ROLLBACK TRAN
				RETURN
			END
		END
		ELSE BEGIN
			-- 1-2. 새로 등록하는 경우
			INSERT INTO tCollectionCard (nCharNo, nCardID, nSerialNo, nStarNumber)
			VALUES (@nCharNo, @nCardID, @nSerialNo, @nStarNumber)
			IF @@ERROR <> 0 BEGIN
				SET @nRet = -9999
				ROLLBACK TRAN
				RETURN
			END
		END
		
		-- 2. 컬렉션에 등록된 카드의 경우 tItem 테이블에서 해당 아이템을 삭제
		DELETE FROM tItem WHERE nItemKey = @nItemKey
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = -9999
			ROLLBACK TRAN
			RETURN
		END	
		
		COMMIT TRAN		
	END
	ELSE BEGIN
		SET @nRet = -9999
		RETURN
	END
END

0쨁ₚU
GO

-- ========================================================================
-- Procedure: usp_Collection_setFavorites
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Collection_setFavorites]
	@nCharNo      int
,	@nSlotNo      tinyint	
,	@nCardID      int
,	@bRegist      bit
,	@nRet         int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	SET @nRet = 0
	
	IF ISNULL(@nCharNo, 0) = 0 BEGIN
		SET @nRet = -1111
		RETURN
	END

	IF ISNULL(@nSlotNo, -1) < 0 BEGIN
		SET @nRet = -1111
		RETURN	
	END

	BEGIN TRAN

	IF @bRegist = 1 BEGIN
		IF ISNULL(@nCardID, 0) = 0 BEGIN
			SET @nRet = -1111
			ROLLBACK TRAN
			RETURN	
		END		
		
		-- 1. 마이컬렉션 카드 등록
		IF EXISTS(SELECT nCharNo FROM tCollectionFavorites WHERE nCharNo = @nCharNo AND nSlotNo = @nSlotNo) BEGIN
			-- 1-1. 이미 존재하는 경우
			UPDATE tCollectionFavorites 
			   SET nCardID = @nCardID 
			 WHERE nCharNo = @nCharNo
			   AND nSlotNo = @nSlotNo 
			IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
				SET @nRet = -9999
				ROLLBACK TRAN
				RETURN
			END
		END
		ELSE BEGIN
			-- 1-2. 새로 등록하는 경우
			INSERT INTO tCollectionFavorites (nCharNo, nSlotNo, nCardID)
			VALUES (@nCharNo, @nSlotNo, @nCardID)
			IF @@ERROR <> 0 BEGIN
				SET @nRet = -9999
				ROLLBACK TRAN
				RETURN
			END
		END
	END
	ELSE BEGIN
		IF EXISTS(SELECT nCharNo FROM tCollectionFavorites WHERE nCharNo = @nCharNo AND nSlotNo = @nSlotNo) BEGIN
			DELETE FROM tCollectionFavorites WHERE nCharNo = @nCharNo AND nSlotNo = @nSlotNo
			IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
				SET @nRet = -9999
				ROLLBACK TRAN
				RETURN
			END
		END
		ELSE BEGIN
			SET @nRet = -5555
			ROLLBACK TRAN
			RETURN
		END
	END
	
	COMMIT TRAN
END

0́ᒿV
GO

-- ========================================================================
-- Procedure: usp_Collection_setReward
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Collection_setReward]
	@nCharNo      int
,	@nReward      int
,	@nRet         int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0
	
	IF ISNULL(@nCharNo, 0) = 0 BEGIN
		SET @nRet = -1111
		RETURN
	END

	IF EXISTS(SELECT nCharNo FROM tCollectionReward WHERE nCharNo = @nCharNo AND nReward = @nReward) BEGIN
		SET @nRet = -2222
		RETURN
	END
	ELSE BEGIN
		INSERT INTO tCollectionReward (nCharNo, nReward)
		VALUES (@nCharNo, @nReward)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -3333
			RETURN
		END
	END
END

0㰁ࣣW
GO

-- ========================================================================
-- Procedure: usp_Collection_setSerial
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Collection_setSerial]
	@nCardID      int
,	@nCount       bigint
,	@nRet         int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0
	
	IF ISNULL(@nCardID, 0) = 0 BEGIN
		SET @nRet = -1111
		RETURN
	END
	
	IF EXISTS(SELECT nCardID FROM tCollectionSerialCount WHERE nCardID = @nCardID) BEGIN
		BEGIN TRAN
		
		UPDATE tCollectionSerialCount 
		   SET sys_update_date = GETDATE()
		     , nCount          = @nCount
		 WHERE nCardID = @nCardID 
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = -4444
			ROLLBACK TRAN
			RETURN
		END
		
		COMMIT TRAN
	END
	ELSE BEGIN
		IF @nCount <> 1 BEGIN
			SET @nRet = -5555
			RETURN
		END
		
		INSERT INTO tCollectionSerialCount (nCardID, nCount)
		VALUES (@nCardID, @nCount)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -6666
			RETURN
		END
	END
END
GO

-- ========================================================================
-- Procedure: usp_Dice_chgMachineToMoneycard
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Dice_chgMachineToMoneycard]
	@nItemKey       bigint
,	@nItemIDOld     int
,	@nItemIDNew     int
,	@nOptionType    smallint
,	@nTotalDividend bigint
,	@nRet              int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	DECLARE @nIsPlace AS smallint 
	DECLARE @nIsPlaceValue AS bigint

	SET @nRet          = 0
	SET @nIsPlace      = 904 -- 가구아이템이 거치된 아이템인지 판단하는 OptionType
	SET @nIsPlaceValue = 0   -- 거치된 아이템을 회수하기 위함 Value

	IF ISNULL(@nItemKey, 0) <= 0 BEGIN
		SET @nRet = -2412
		RETURN
	END

	IF ISNULL(@nItemIDOLD, 0) <= 0 BEGIN
		SET @nRet = -2417
		RETURN
	END

	IF ISNULL(@nItemIDNEW, 0) <= 0 BEGIN
		SET @nRet = -2418
		RETURN
	END

	IF ISNULL(@nOptionType, 0) <= 0 BEGIN
		SET @nRet = -2419
		RETURN
	END

	IF ISNULL(@nTotalDividend, 0) <= 0 BEGIN
		SET @nRet = -2416
		RETURN
	END

	IF EXISTS(SELECT nItemKey FROM tItem WHERE nItemKey = @nItemKey AND nItemID = @nItemIDOLD) BEGIN
		BEGIN TRAN

		-- 게임기의 아이템아이디를 머니카드의 아이템아이디로 교체
		UPDATE tItem SET nItemID = @nItemIDNEW
		 WHERE nItemKey = @nItemKey
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = -2493
			ROLLBACK TRAN
			RETURN
		END

		-- 거치된 게임기의 경우 인벤토리로 회수하는 작업을 함.
		UPDATE tItemOptions SET nOptionData = @nIsPlaceValue
		 WHERE nItemKey = @nItemKey AND nOptionType = @nIsPlace
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = -2493
			ROLLBACK TRAN
		END

		-- 머니카드에 배당금을 확정함
		UPDATE tItemOptions SET nOptionData = @nTotalDividend
		 WHERE nItemKey = @nItemKey AND nOptionType = @nOptionType
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = -2493
			ROLLBACK TRAN
		END
		ELSE BEGIN
			SET @nRet = 0
			COMMIT TRAN
		END
	END
	ELSE BEGIN
		SET @nRet = -2422
	END

	SET LOCK_TIMEOUT -1
	SET NOCOUNT OFF
END

0긁X
GO

-- ========================================================================
-- Procedure: usp_Dice_getDiceExpect
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Dice_getDiceExpect]
	@nCharNo  int
,	@nItemKey bigint
,	@nRet     int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	IF ISNULL(@nCharNo, 0) <= 0 BEGIN
		SET @nRet = -2411
		RETURN
	END

	IF ISNULL(@nItemKey, 0) <= 0 BEGIN
		SET @nRet = -2412
		RETURN
	END

	SELECT ISNULL(COUNT(DISTINCT(A.nGameNo)), 0) AS nCount, ISNULL(SUM(nTotalDividend), 0) AS nTotalDividend
	  FROM tDice AS A, tDiceResult AS B
	 WHERE A.nGameNo = B.nGameNo
	   AND A.nCharNo = @nCharNo
	   AND A.nItemKey = @nItemKey

	SET NOCOUNT OFF
END

0Y
GO

-- ========================================================================
-- Procedure: usp_Dice_getDiceRank
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Dice_getDiceRank]
	@nCharNo  int
,	@nItemKey bigint
,	@nRet     int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	IF ISNULL(@nCharNo, 0) <= 0 BEGIN
		SET @nRet = -2411
		RETURN
	END

	IF ISNULL(@nItemKey, 0) <= 0 BEGIN
		SET @nRet = -2412
		RETURN
	END

	SELECT TOP 20 C.sID, COUNT(B.nGameNo) AS nWinCount
	  FROM tDice AS A, tDiceLog AS B, tCharacter AS C
	 WHERE A.nGameNo = B.nGameNo
	   AND B.nCustomCharNo = C.nCharNo
	   AND A.nCharNo = @nCharNo
	   AND A.nItemKey = @nItemKey
	   AND B.bStatus = 1
	   AND C.bDeleted = 0
	 GROUP BY C.sID
	 ORDER BY nWinCount DESC

	SET NOCOUNT OFF
END

0夁춘[
GO

-- ========================================================================
-- Procedure: usp_Dice_setDice
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Dice_setDice] 
	@nCharNo  int
,	@nItemKey bigint
,	@sDice    nvarchar(3)
,	@nRet     int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	DECLARE @nGameNo AS int

	SET @nRet    = 0
	SET @nGameNo = 0

	IF ISNULL(@nCharNo, 0) <= 0 BEGIN
		SET @nRet = -2411
		RETURN
	END

	IF ISNULL(@nItemKey, 0) <= 0 BEGIN
		SET @nRet = -2412
		RETURN
	END

	IF LEN(ISNULL(@sDice, '')) <= 0 BEGIN
		SET @nRet = -2413
		RETURN
	END

	INSERT INTO tDice (nCharNo, nItemKey, sDice)
	VALUES (@nCharNo, @nItemKey, @sDice)
	IF @@ERROR <> 0 OR @@IDENTITY = 0 BEGIN
		SET @nRet = -2492
		RETURN
	END
	ELSE BEGIN
		SET @nGameNo = @@IDENTITY
		SELECT @nGameNo
	END

	SET NOCOUNT OFF
END
GO

-- ========================================================================
-- Procedure: usp_Dice_setDiceLog
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Dice_setDiceLog]
	@nGameNo       int 
,	@nCustomCharNo int
,	@nSlotNo       int
,	@nBattingMoney bigint
,	@nBattingRate  int
,	@bStatus       bit
,	@nTaxRate      int
,	@nDividend     bigint
,	@nRet          int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	INSERT INTO tDiceLog (nGameNo, nCustomCharNo, nSlotNo, nBattingMoney, nBattingRate, bStatus, nTaxRate, nDividend)
	VALUES (@nGameNo, @nCustomCharNo, @nSlotNo, @nBattingMoney, @nBattingRate, @bStatus, @nTaxRate, @nDividend)
	IF @@ERROR <> 0 BEGIN
		SET @nRet = -2492
		RETURN
	END

	SET NOCOUNT OFF
END

0鈁솼\
GO

-- ========================================================================
-- Procedure: usp_Dice_setDicePay
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Dice_setDicePay]
	@nCharNo        int
,	@nItemKey       bigint
,	@nTotalDividend bigint
,	@nRet           int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	SET @nRet = 0

	IF ISNULL(@nCharNo, 0) <= 0 BEGIN
		SET @nRet = -2411
		RETURN
	END

	IF ISNULL(@nItemKey, 0) <= 0 BEGIN
		SET @nRet = -2412
		RETURN
	END

	IF ISNULL(@nTotalDividend, 0) <= 0 BEGIN
		SET @nRet = -2416
		RETURN
	END

	IF EXISTS(SELECT nItemKey FROM tItem WHERE nItemKey = @nItemKey AND nOwner = @nCharNo) BEGIN
		BEGIN TRAN

		-- žÓŽÏÄ«µå ŸÆÀÌÅÛ »èÁŠ
		DELETE FROM tItem WHERE nItemKey = @nItemKey AND nOwner = @nCharNo 
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = -2494
			ROLLBACK TRAN
			RETURN
		END

		-- µ· Áö±Þ
		UPDATE tCharacter SET nMoney = nMoney + @nTotalDividend WHERE nCharNo = @nCharNo
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = -2493
			ROLLBACK TRAN
			RETURN
		END

		-- µ· Áö±Þ ·Î±×
		INSERT INTO tDicePayLog (nCharNo, nItemKey, nTotalDividend)
		VALUES (@nCharNo, @nItemKey, @nTotalDividend)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -2492
			ROLLBACK TRAN
			RETURN
		END
		ELSE BEGIN
			SET @nRet = 0
			COMMIT TRAN
			RETURN
		END
	END
	ELSE BEGIN
		SET @nRet = -2423
		RETURN	
	END
END
GO

-- ========================================================================
-- Procedure: usp_Dice_setDiceResult
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Dice_setDiceResult] 
	@nGameNo        int
,	@nCustomCharNo  int
,	@nTotalDividend bigint
,	@nRet           int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	DECLARE @money bigint, @resultMoney bigint

	SET @nRet  = 0
	SET @money = 0

	IF ISNULL(@nGameNo, 0) <= 0 BEGIN
		SET @nRet = -2414
		RETURN
	END

	IF ISNULL(@nCustomCharNo, 0) <= 0 BEGIN
		SET @nRet = -2415
		RETURN
	END

	BEGIN TRAN

	SELECT @money = ISNULL(nMoney, 0) FROM tCharacter WITH(UPDLOCK, HOLDLOCK) WHERE nCharNo = @nCustomCharNo
	IF ( @money + @nTotalDividend ) > 0 BEGIN
		SET @resultMoney = @money + @nTotalDividend
	END
	ELSE BEGIN
		SET @resultMoney = 0
	END
	
	UPDATE tCharacter SET nMoney = @resultMoney
	WHERE nCharNo = @nCustomCharNo
	IF @@ERROR <> 0 OR @@ROWCOUNT > 1 BEGIN
		SET @nRet = - 2493
		ROLLBACK TRAN
		RETURN
	END

	INSERT INTO tDiceResult (nGameNo, nCustomCharNo, nBeforeMoney, nTotalDividend, nAfterMoney)
	VALUES (@nGameNo, @nCustomCharNo, @money, @nTotalDividend, @resultMoney)
	IF @@ERROR <> 0 BEGIN
		SET @nRet = - 2492
		ROLLBACK TRAN
		RETURN
	END
	ELSE BEGIN
		SET @nRet = 0
		COMMIT TRAN
	END

	SELECT nCharNo AS nCustomCharNo, nMoney FROM tCharacter WHERE nCharNo = @nCustomCharNo
END
GO

-- ========================================================================
-- Procedure: usp_GM_getItemRewardGroupNo
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_GM_getItemRewardGroupNo]
	@sOperID nvarchar(20)
,	@sIP     nvarchar(30)
,	@sCause  nvarchar(1000)
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	SET @nRet = 0
		
	INSERT INTO tGMItemReward (sOperID, sIP, sCause)
	VALUES (@sOperID, @sIP, @sCause)
	IF @@ERROR <> 0 OR @@IDENTITY <= 0 BEGIN
		SET @nRet = -1111
		RETURN
	END
	ELSE BEGIN
		SET @nRet = 0
		
		SELECT @@IDENTITY
		RETURN
	END
END

0靂ŭ
GO

-- ========================================================================
-- Procedure: usp_GM_getItemRewardList
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_GM_getItemRewardList]
	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0
		
	SELECT nGroupNo, sys_create_date AS dDate, sCause 
	  FROM tGMItemReward 
	 ORDER BY nGroupNo DESC
	
END

0嬁l
GO

-- ========================================================================
-- Procedure: usp_GM_setItemReward
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_GM_setItemReward]
	@nGroupNo     bigint
,	@nCharNo      int
,	@nStorageType tinyint
,	@nItemID      int
,	@nItemKey     bigint
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000
	
	DECLARE @nStorage AS smallint

	SET @nStorage = 0 -- nStoragType = 2 번의 nStorage = 0 번에 보상인벤 제일 앞에 인벤토리임.
	SET @nRet = 0
	
	IF NOT EXISTS(SELECT nCharNo FROM tCharacter WHERE nCharNo = @nCharNo AND bDeleted = 0) BEGIN
		SET @nRet = -1111
		RETURN
	END
		
	IF EXISTS(SELECT nGroupNo FROM tGMItemRewardLog WHERE nGroupNo = @nGroupNo AND nCharNo = @nCharNo) BEGIN
		SET @nRet = -2222
		RETURN
	END
	ELSE BEGIN
		BEGIN TRAN
		
		UPDATE tItem 
		   SET nStorage = nStorage + 1 
		 WHERE nOwner = @nCharNo AND nStorageType = @nStorageType 
		IF @@ERROR <> 0 BEGIN
			ROLLBACK TRAN
			SET @nRet = -3333
			RETURN
		END
		ELSE BEGIN
			INSERT INTO tItem (nItemKey, nStorageType, nOwner, nStorage, nItemID)
			VALUES (@nItemKey, @nStorageType, @nCharNo, @nStorage, @nItemID)
			IF @@ERROR <> 0 BEGIN
				ROLLBACK TRAN
				SET @nRet = -4444
				RETURN
			END	
			
			INSERT INTO tGMItemRewardLog (nGroupNo, nCharNo, nStorageType, nItemID, nItemKey)
			VALUES (@nGroupNo, @nCharNo, @nStorageType, @nItemID, @nItemKey)
			IF @@ERROR <> 0 BEGIN
				ROLLBACK TRAN
				SET @nRet = -5555
				RETURN
			END	
			ELSE BEGIN
				COMMIT TRAN
			END			
		END
	END
END

0鐁m
GO

-- ========================================================================
-- Procedure: usp_Gamble_addCoin
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Gamble_addCoin]
	@nCharNo  int
,	@nItemKey bigint
,	@nCoin    bigint
,	@nRet     int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	DECLARE @status tinyint, @coin bigint
	
	SET @status = 2 -- À¯·áŸÆÀÌÅÛ -> ÄÚÀÎ
	SET @nRet = 0

	IF ISNULL(@nCharNo, 0) <= 0 BEGIN
		SET @nRet = -3311
		RETURN
	END

	IF ISNULL(@nCoin, 0) <= 0 BEGIN
		SET @nRet = -3313
		RETURN
	END
		
	IF EXISTS(SELECT nItemKey FROM tItem WHERE nItemKey = @nItemKey AND nOwner = @nCharNo) BEGIN
		BEGIN TRAN

		DELETE FROM tItem WHERE nItemKey = @nItemKey AND nOwner = @nCharNo 
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			ROLLBACK TRAN
			SET @nRet = -3332
			RETURN		
		END
		
		SELECT @coin = nCoin FROM tGambleCoin WITH(UPDLOCK, HOLDLOCK) WHERE nCharNo = @nCharNo
		IF @coin IS NULL BEGIN
			INSERT INTO tGambleCoin (nCharNo, nCoin)
			VALUES (@nCharNo, @nCoin)
			IF @@ERROR <> 0 BEGIN
				ROLLBACK TRAN
				SET @nRet = -3329
				RETURN
			END	
			
			SET @coin = 0	
		END
		ELSE BEGIN
			UPDATE tGambleCoin 
			   SET sys_update_date = GETDATE()
				 , nCoin = nCoin + @nCoin 
			 WHERE nCharNo = @nCharNo 
			IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
				ROLLBACK TRAN
				SET @nRet = -3323
				RETURN
			END
		END
		
		INSERT INTO tGambleCoinExchangeLog (nCharNo, nExchangedCoin, nStatus, nBeforeCoin, nAfterCoin)
		VALUES (@nCharNo, @nCoin, @status, @coin, @coin + @nCoin)
		IF @@ERROR <> 0 BEGIN
			ROLLBACK TRAN
			SET @nRet = -3326
			RETURN
		END
		ELSE BEGIN
			COMMIT TRAN
			RETURN
		END
	END
	ELSE BEGIN
		SET @nRet = -3331
		RETURN	
	END
END
GO

-- ========================================================================
-- Procedure: usp_Gamble_getCoin
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Gamble_getCoin]
	@nCharNo int
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	SET @nRet = 0
	
	IF EXISTS(SELECT nCharNo FROM tGambleCoin WHERE nCharNo = @nCharNo) BEGIN
		SELECT nCoin FROM tGambleCoin WHERE nCharNo = @nCharNo
	END
	ELSE BEGIN
		SELECT 0 AS nCoin
	END
END

0
GO

-- ========================================================================
-- Procedure: usp_Gamble_getDiceRank
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Gamble_getDiceRank]
	@nRet     int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	SELECT TOP 20 C.sID, COUNT(B.nGameNo) AS nWinCount
	  FROM tGambleDice AS A, tGambleDiceLog AS B, tCharacter AS C
	 WHERE A.nGameNo = B.nGameNo
	   AND B.nCharNo = C.nCharNo
	   AND B.bStatus = 1
	   AND C.bDeleted = 0
	 GROUP BY C.sID
	 ORDER BY nWinCount DESC

	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: usp_Gamble_getExchangedCoin
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Gamble_getExchangedCoin]
	@nCharNo int
,	@dDate   datetime
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	SET @nRet = 0
	
	SELECT ISNULL(SUM(nExchangedCoin), 0) AS nExchangedCoin
	  FROM tGambleCoinExchangeLog 
	 WHERE nCharNo = @nCharNo 
	   AND dDate   >= @dDate 
	   AND nStatus = 1
END

0
GO

-- ========================================================================
-- Procedure: usp_Gamble_getJackpotCoin
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Gamble_getJackpotCoin] 
	@sGameName    nvarchar(20)
,	@nRet         int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0
	
	SELECT nBattingUnit, nJackpotCoin 
	  FROM tGambleJackpot 
	 WHERE sGameName = @sGameName 
END
GO

-- ========================================================================
-- Procedure: usp_Gamble_getJackpotRanking
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Gamble_getJackpotRanking] 
	@sGameName    nvarchar(20)
,	@nBattingUnit bigint
,	@nRet         int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0
	
	SELECT TOP 50 sID, nJackpotCoin  
	  FROM tGambleJackpotLog 
	 WHERE sGameName = @sGameName 
	   AND nBattingUnit = @nBattingUnit 
END

0
GO

-- ========================================================================
-- Procedure: usp_Gamble_setCoinToMoney
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Gamble_setCoinToMoney]
	@nCharNo int
,	@nMoney  bigint
,	@nCoin   bigint
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	DECLARE @status tinyint
	DECLARE @money bigint, @coin bigint
	
	SET @status = 0 -- ÄÚÀÎ -> °ÔÀÓžÓŽÏ
	SET @nRet = 0
	
	IF ISNULL(@nCharNo, 0) <= 0 BEGIN
		SET @nRet = -3311
		RETURN
	END

	IF ISNULL(@nMoney, 0) <= 0 BEGIN
		SET @nRet = -3312
		RETURN
	END	
	
	IF ISNULL(@nCoin, 0) <= 0 BEGIN
		SET @nRet = -3313
		RETURN
	END	

	BEGIN TRAN
		
	SELECT @coin = nCoin FROM tGambleCoin WITH(UPDLOCK, HOLDLOCK) WHERE nCharNo = @nCharNo
	IF @coin IS NULL BEGIN
		SET @nRet = -3321
		ROLLBACK TRAN
		RETURN
	END
	
	IF @coin < @nCoin BEGIN
		SET @nRet = -3322
		ROLLBACK TRAN
		RETURN
	END
	
	UPDATE tGambleCoin 
	   SET sys_update_date = GETDATE()
		 , nCoin = nCoin - @nCoin 
	 WHERE nCharNo = @nCharNo 
	   AND nCoin  >= @coin
	IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
		ROLLBACK TRAN
		SET @nRet = -3323
		RETURN
	END
	
	SELECT @money = nMoney FROM tCharacter WITH(UPDLOCK, HOLDLOCK) WHERE nCharNo = @nCharNo AND bDeleted = 0
	IF @money IS NULL BEGIN		
		ROLLBACK TRAN
		SET @nRet = -3324
		RETURN
	END
	ELSE BEGIN
		UPDATE tCharacter
		   SET nMoney   = nMoney + @nMoney
		 WHERE nCharNo  = @nCharNo
		   AND bDeleted = 0
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			ROLLBACK TRAN
			SET @nRet = -3325
			RETURN		
		END	
	END

	INSERT INTO tGambleCoinExchangeLog (nCharNo, nExchangedCoin, nStatus, nBeforeCoin, nAfterCoin, nBeforeMoney, nAfterMoney)
	VALUES (@nCharNo, @nCoin, @status, @coin, @coin - @nCoin, @money, @money + @nMoney)
	IF @@ERROR <> 0 BEGIN
		ROLLBACK TRAN
		SET @nRet = -3326
		RETURN
	END
	ELSE BEGIN
		COMMIT TRAN
		RETURN
	END
END
GO

-- ========================================================================
-- Procedure: usp_Gamble_setDice
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Gamble_setDice] 
	@sDice    nvarchar(3)
,	@nRet     int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	DECLARE @nGameNo AS int

	SET @nRet    = 0
	SET @nGameNo = 0

	IF LEN(ISNULL(@sDice, '')) <= 0 BEGIN
		SET @nRet = -3314
		RETURN
	END

	INSERT INTO tGambleDice (sDice)
	VALUES (@sDice)
	IF @@ERROR <> 0 OR @@IDENTITY = 0 BEGIN
		SET @nRet = -3392
		RETURN
	END
	ELSE BEGIN
		SET @nGameNo = @@IDENTITY
		SELECT @nGameNo
	END

	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: usp_Gamble_setDiceLog
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Gamble_setDiceLog]
	@nGameNo       int 
,	@nCharNo       int
,	@nSlotNo       int
,	@nBattingMoney bigint
,	@nBattingRate  int
,	@bStatus       bit
,	@nTaxRate      int
,	@nDividend     bigint
,	@nRet          int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	INSERT INTO tGambleDiceLog (nGameNo, nCharNo, nSlotNo, nBattingMoney, nBattingRate, bStatus, nTaxRate, nDividend)
	VALUES (@nGameNo, @nCharNo, @nSlotNo, @nBattingMoney, @nBattingRate, @bStatus, @nTaxRate, @nDividend)
	IF @@ERROR <> 0 BEGIN
		SET @nRet = -3392
		RETURN
	END

	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: usp_Gamble_setDiceResult
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Gamble_setDiceResult] 
	@nGameNo        int
,	@nCharNo        int
,	@nTotalDividend bigint
,	@nRet           int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	DECLARE @coin bigint, @resultCoin bigint

	SET @nRet  = 0
	SET @coin = 0
	SET @resultCoin = 0

	IF ISNULL(@nGameNo, 0) <= 0 BEGIN
		SET @nRet = -3315
		RETURN
	END

	IF ISNULL(@nCharNo, 0) <= 0 BEGIN
		SET @nRet = -3311
		RETURN
	END

	BEGIN TRAN

	SELECT @coin = ISNULL(nCoin, 0) FROM tGambleCoin WITH(UPDLOCK, HOLDLOCK) WHERE nCharNo = @nCharNo
	IF ( @coin + @nTotalDividend ) > 0 BEGIN
		SET @resultCoin = @coin + @nTotalDividend
	END
	ELSE BEGIN
		SET @resultCoin = 0
	END
	
	UPDATE tGambleCoin 
	   SET nCoin = @resultCoin
	     , sys_update_date = GETDATE()
	 WHERE nCharNo = @nCharNo
	IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
		SET @nRet = -3323
		ROLLBACK TRAN
		RETURN
	END

	INSERT INTO tGambleDiceResult (nGameNo, nCharNo, nBeforeMoney, nTotalDividend, nAfterMoney)
	VALUES (@nGameNo, @nCharNo, @coin, @nTotalDividend, @resultCoin)
	IF @@ERROR <> 0 BEGIN
		SET @nRet = - 3330
		ROLLBACK TRAN
		RETURN
	END
	ELSE BEGIN
		SET @nRet = 0
		COMMIT TRAN
	END

	SELECT nCharNo, nCoin FROM tGambleCoin WHERE nCharNo = @nCharNo
END
GO

-- ========================================================================
-- Procedure: usp_Gamble_setMoneyToCoin
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Gamble_setMoneyToCoin]
	@nCharNo int
,	@nMoney  bigint
,	@nCoin   bigint
,	@nLimit  bigint
,	@dDate   datetime
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	DECLARE @status tinyint
	DECLARE @money bigint, @coin bigint, @exchangedCoin bigint
	
	SET @status = 1 -- °ÔÀÓžÓŽÏ -> ÄÚÀÎ
	SET @nRet = 0
	
	IF ISNULL(@nCharNo, 0) <= 0 BEGIN
		SET @nRet = -3311
		RETURN
	END

	IF ISNULL(@nMoney, 0) <= 0 BEGIN
		SET @nRet = -3312
		RETURN
	END	
	
	IF ISNULL(@nCoin, 0) <= 0 BEGIN
		SET @nRet = -3313
		RETURN
	END	
	
	-- ÀÏÁ€œÃ°£µ¿ŸÈ ±³È¯ÇÑ ÄÚÀÎÀÌ ³ÑÀ» °æ¿ì ±³È¯ÇÒ Œö ŸøÀœ
	SELECT @exchangedCoin = ISNULL(SUM(nExchangedCoin), 0) 
	  FROM tGambleCoinExchangeLog 
	 WHERE nCharNo = @nCharNo 
	   AND dDate  >= @dDate
	   AND nStatus = @status
	IF @nLimit < @exchangedCoin + @nCoin BEGIN
		SET @nRet = -3327
		RETURN	
	END

	BEGIN TRAN
		
	SELECT @money = ISNULL(nMoney, 0) FROM tCharacter WITH(UPDLOCK, HOLDLOCK) WHERE nCharNo = @nCharNo AND bDeleted = 0
	IF @money IS NULL BEGIN
		SET @nRet = -3324
		ROLLBACK TRAN
		RETURN
	END
	
	IF @money < @nMoney BEGIN
		SET @nRet = -3328
		ROLLBACK TRAN
		RETURN
	END
		
	UPDATE tCharacter
	   SET nMoney   = nMoney - @nMoney
	 WHERE nCharNo  = @nCharNo
	   AND bDeleted = 0
	   AND nMoney  >= @nMoney
	IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
		ROLLBACK TRAN
		SET @nRet = -3325
		RETURN		
	END
	
	SELECT @coin = nCoin FROM tGambleCoin WITH(UPDLOCK, HOLDLOCK) WHERE nCharNo = @nCharNo
	IF @coin IS NULL BEGIN
		INSERT INTO tGambleCoin (nCharNo, nCoin)
		VALUES (@nCharNo, @nCoin)
		IF @@ERROR <> 0 BEGIN
			ROLLBACK TRAN
			SET @nRet = -3329
			RETURN
		END	
		
		SET @coin = 0
	END
	ELSE BEGIN
		UPDATE tGambleCoin 
		   SET sys_update_date = GETDATE()
			 , nCoin = nCoin + @nCoin 
		 WHERE nCharNo = @nCharNo 
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			ROLLBACK TRAN
			SET @nRet = -3323
			RETURN
		END
	END

	INSERT INTO tGambleCoinExchangeLog (nCharNo, nExchangedCoin, nStatus, nBeforeCoin, nAfterCoin, nBeforeMoney, nAfterMoney)
	VALUES (@nCharNo, @nCoin, @status, @coin, @coin + @nCoin, @money, @money - @nMoney)
	IF @@ERROR <> 0 BEGIN
		ROLLBACK TRAN
		SET @nRet = -3326
		RETURN
	END
	ELSE BEGIN
		COMMIT TRAN
		RETURN
	END
END
GO

-- ========================================================================
-- Procedure: usp_Gamble_setSlotMachineResult
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Gamble_setSlotMachineResult] 
	@nCharNo      int
,	@nBattingUnit bigint
,	@nBatting     tinyint
,	@sGameResult  nvarchar(10)
,	@bIsWin       bit
,	@nPrizeCoin   bigint
,	@bIsJackpot   bit
,	@nJackpotRate int
,	@nRet         int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	DECLARE @gameName nvarchar(20), @charID nvarchar(80)
	DECLARE @gameNo bigint, @battingCoin bigint, @coin bigint, @resultCoin bigint, @jackpotCoin bigint
	
	SET @nRet = 0
	SET @gameName = 'SlotMachine'
	SET @battingCoin = 0
	SET @resultCoin  = 0
	SET @battingCoin = @nBattingUnit * @nBatting 
		
	IF ISNULL(@nCharNo, 0) = 0 BEGIN
		SET @nRet = -3311
		RETURN
	END
	IF ISNULL(@nBattingUnit, 0) = 0 BEGIN
		SET @nRet = -3314
		RETURN		
	END
	IF ISNULL(@nBatting, 0) = 0 BEGIN
		SET @nRet = -3315
		RETURN		
	END	
	IF ISNULL(@sGameResult, '') = '' BEGIN
		SET @nRet = -3316
		RETURN		
	END
	
	SELECT @charID = sID FROM tCharacter WHERE nCharNo = @nCharNo AND bDeleted = 0
	IF @charID IS NOT NULL BEGIN
		BEGIN TRAN
		
		-- 1. ???? ??
		INSERT INTO tGambleSlotMachineLog (nCharNo, nBattingUnit, nBatting, sGameResult, bIsWin, nPrizeCoin)
		VALUES (@nCharNo, @nBattingUnit, @nBatting, @sGameResult, @bIsWin, @nPrizeCoin)
		IF @@ERROR <> 0 OR @@IDENTITY = 0 BEGIN
			SET @nRet = -3333
			ROLLBACK TRAN
			RETURN
		END
		SET @gameNo = @@IDENTITY
		
		-- 2. ??( ????(-) / ????(+)
		SELECT @coin = ISNULL(nCoin, 0) FROM tGambleCoin WITH(UPDLOCK, HOLDLOCK) WHERE nCharNo = @nCharNo
		IF (@coin - @battingCoin + @nPrizeCoin) > 0 BEGIN
			SET @resultCoin = @coin - @battingCoin + @nPrizeCoin
		END
		ELSE BEGIN
			SET @resultCoin = 0
		END
		
		UPDATE tGambleCoin 
		   SET nCoin = @resultCoin
			 , sys_update_date = GETDATE()
		 WHERE nCharNo = @nCharNo
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = -3334
			ROLLBACK TRAN
			RETURN
		END		 
		
		-- 3. ??? ??
		IF @bIsJackpot = 1 BEGIN
			-- 3-1. ?? ???? ??
			SELECT @jackpotCoin = ISNULL(nJackpotCoin, 0) 
			  FROM tGambleJackpot WITH(UPDLOCK, HOLDLOCK)
			 WHERE sGameName = @gameName
			   AND nBattingUnit = @nBattingUnit
			IF @jackpotCoin <> 0 BEGIN
				SET @jackpotCoin = @jackpotCoin * @nJackpotRate / 1000
			END
			
			UPDATE tGambleCoin 
			   SET nCoin = nCoin + @jackpotCoin 
			     , sys_update_date = GETDATE()
			 WHERE nCharNo = @nCharNo 
			IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
				SET @nRet = -3335
				ROLLBACK TRAN
				RETURN
			END
			
			-- 3-2. ?? ???? ???
			UPDATE tGambleJackpot 
			   SET nJackpotCoin = 0
			     , sys_update_date = GETDATE()
			 WHERE sGameName = @gameName 
			   AND nBattingUnit = @nBattingUnit 
			IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
				SET @nRet = -3336
				ROLLBACK TRAN
				RETURN
			END
			
			-- 3-3. ???? ??
			INSERT INTO tGambleJackpotLog (sGameName, nBattingUnit, sID, nJackpotCoin, nGameNo)
			VALUES (@gameName, @nBattingUnit, @charID, @jackpotCoin, @gameNo)
			IF @@ERROR <> 0 BEGIN
				SET @nRet = -3337
				ROLLBACK TRAN
				RETURN
			END						   				  
		END
		-- 4. ??? ????
		ELSE BEGIN
			-- 4-1. ???? ??
			IF EXISTS(SELECT sGameName FROM tGambleJackpot WHERE sGameName = @gameName AND nBattingUnit = @nBattingUnit) BEGIN
				UPDATE tGambleJackpot 
				   SET nJackpotCoin = nJackpotCoin + @battingCoin 
					 , sys_update_date = GETDATE()
				 WHERE sGameName = @gameName
				   AND nBattingUnit = @nBattingUnit 
				IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
					SET @nRet = -3338
					ROLLBACK TRAN
					RETURN
				END	
			END
			ELSE BEGIN
				INSERT INTO tGambleJackpot (sGameName, nBattingUnit, nJackpotCoin)
				VALUES (@gameName, @nBattingUnit, @battingCoin)
				IF @@ERROR <> 0 BEGIN
					SET @nRet = -3338
					ROLLBACK TRAN
					RETURN
				END				
			END		   
		END
		
		COMMIT TRAN
	END
	ELSE BEGIN
		SET @nRet = -3339
	END
END
GO

-- ========================================================================
-- Procedure: usp_Gamble_subCoin
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Gamble_subCoin]
	@nCharNo  int
,	@nCoin    bigint
,	@nRet     int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	DECLARE @status tinyint
	DECLARE @coin bigint

	SET @status = 3 -- ??	-> ???
	SET @nRet = 0

	IF ISNULL(@nCharNo, 0) <= 0 BEGIN
		SET @nRet = -3311
		RETURN
	END

	IF ISNULL(@nCoin, 0) < 0 BEGIN
		SET @nRet = -3313
		RETURN
	END
	
	-- ??? ??? 0 ?? ??? ?? ??? ????? ??.
	IF @nCoin = 0 BEGIN
		SET @nRet = 0
		RETURN
	END	
		
	SELECT @coin = ISNULL(nCoin, 0) FROM tGambleCoin WITH(UPDLOCK, HOLDLOCK) WHERE nCharNo = @nCharNo
	IF @coin < @nCoin BEGIN
		SET @nRet = -3322
		RETURN		
	END
	ELSE BEGIN
		BEGIN TRAN
		
		UPDATE tGambleCoin 
		   SET sys_update_date = GETDATE()
			 , nCoin = nCoin - @nCoin 
		 WHERE nCharNo = @nCharNo 
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			ROLLBACK TRAN
			SET @nRet = -3323
			RETURN
		END

		INSERT INTO tGambleCoinExchangeLog (nCharNo, nExchangedCoin, nStatus, nBeforeCoin, nAfterCoin)
		VALUES (@nCharNo, @nCoin, @status, @coin, @coin - @nCoin)
		IF @@ERROR <> 0 BEGIN
			ROLLBACK TRAN
			SET @nRet = -3326
			RETURN
		END
		ELSE BEGIN
			COMMIT TRAN
			RETURN
		END
	END		
END

0
GO

-- ========================================================================
-- Procedure: usp_GuildAcademy_delMember
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_GuildAcademy_delMember]
	@nNo     int
,	@nCharNo int
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	IF @nNo <= 0 BEGIN
		SET @nRet = -1511
		RETURN
	END

	IF @nCharNo <= 0 BEGIN
		SET @nRet = -1512
		RETURN
	END

	DELETE FROM tGuildAcademyMember WHERE nNo = @nNo AND nCharNo = @nCharNo
	IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
		SET @nRet = -1592
	END

	SET NOCOUNT OFF
END

0ꈁq
GO

-- ========================================================================
-- Procedure: usp_GuildAcademy_delReward
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_GuildAcademy_delReward]
	@nNo     int
,	@nLevel  smallint
,	@nTypeID tinyint
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	IF @nNo <= 0 BEGIN
		SET @nRet = -1511
		RETURN
	END

	IF @nTypeID = 1 BEGIN
		DELETE FROM tGuildAcademyReward WHERE nNo = @nNo AND nLevel = @nLevel AND nTypeID = @nTypeID
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = -1596
			RETURN
		END
	END
	ELSE BEGIN
		DELETE FROM tGuildAcademyReward WHERE nNo = @nNo AND nLevel = @nLevel AND nTypeID > 1
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -1596
			RETURN
		END
	END

	SET NOCOUNT OFF
END

0ᐁs
GO

-- ========================================================================
-- Procedure: usp_GuildAcademy_delRewardAll
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_GuildAcademy_delRewardAll]
	@nNo     int
,	@nTypeID tinyint
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	IF @nNo <= 0 BEGIN
		SET @nRet = -1511
		RETURN
	END

	IF @nTypeID = 1 BEGIN
		DELETE FROM tGuildAcademyReward WHERE nNo = @nNo AND nTypeID = @nTypeID
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -1596
			RETURN
		END
	END
	ELSE BEGIN
		DELETE FROM tGuildAcademyReward WHERE nNo = @nNo AND nTypeID > 1
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -1596
			RETURN
		END
	END

	SET NOCOUNT OFF
END

0䴁칐t
GO

-- ========================================================================
-- Procedure: usp_GuildAcademy_depositMoney
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_GuildAcademy_depositMoney]
	@nNo     int
,	@nCharNo int
,	@nMoney  bigint
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	DECLARE @money AS bigint
	SET @nRet = 0
	SET @money = 0

	IF @nNo <= 0 BEGIN
		SET @nRet = -1511
		RETURN
	END

	IF @nCharNo <= 0 BEGIN
		SET @nRet = -1512
		RETURN
	END

	IF @nMoney <= 0 BEGIN
		SET @nRet = -1514
		RETURN
	END

	BEGIN TRAN
	
	IF NOT EXISTS(SELECT nCharNo FROM tCharacter WHERE nCharNo = @nCharNo AND bDeleted = 0) BEGIN
		SET @nRet = -1524
		ROLLBACK TRAN
		RETURN
	END

	IF NOT EXISTS(SELECT nNo FROM tGuildAcademy WHERE nNo = @nNo) BEGIN
		SET @nRet = -1525
		ROLLBACK TRAN
		RETURN
	END

	SELECT @money = nMoney FROM tCharacter WHERE nCharNo = @nCharNo
	IF @money < @nMoney BEGIN		
		SET @nRet = -1526
		ROLLBACK TRAN
		RETURN
	END
	ELSE BEGIN
		UPDATE tCharacter SET nMoney = nMoney - @nMoney WHERE nCharNo = @nCharNo
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = -1593
			ROLLBACK TRAN
			RETURN
		END

		UPDATE tGuildAcademy SET nMoney = nMoney + @nMoney WHERE nNo = @nNo
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = -1593
			ROLLBACK TRAN
			RETURN
		END
		ELSE BEGIN
			SET @nRet = 0
			COMMIT TRAN
		END
	END
			
	SET NOCOUNT OFF
END

0蘁쉴u
GO

-- ========================================================================
-- Procedure: usp_GuildAcademy_getGraduateCount
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_GuildAcademy_getGraduateCount]
	@nNo int
,	@nRet int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	IF @nNo <= 0 BEGIN
		SET @nRet = -1511
		RETURN
	END

	SELECT ISNULL(COUNT(nCharNo), 0) AS cnt
	  FROM tGuildAcademyGraduate
	 WHERE nNo = @nNo
	IF @@ERROR <> 0 BEGIN
		SET @nRet = -1591
	END

	SET NOCOUNT OFF
END

0뼁뚘v
GO

-- ========================================================================
-- Procedure: usp_GuildAcademy_getGraduateList
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_GuildAcademy_getGraduateList]
	@nNo int
,	@nRet int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	IF @nNo <= 0 BEGIN
		SET @nRet = -1511
		RETURN
	END

	SELECT sID, nClass, nJoinLevel, dJoinDate, nGraduateLevel, dGraduateDate
	  FROM tGuildAcademyGraduate
	 WHERE nNo = @nNo
	IF @@ERROR <> 0 BEGIN
		SET @nRet = -1591
	END

	SET NOCOUNT OFF
END

0ꪼw
GO

-- ========================================================================
-- Procedure: usp_GuildAcademy_getGuildAcademyList
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_GuildAcademy_getGuildAcademyList]
	@nRet int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	SELECT GA.nNo, nMasterCharNo, C1.sID, nMaxMember, GA.nMoney, nPoint
		 , T2.nRanking, ISNULL(dBuffEndTime, GETDATE()), ISNULL(dNotifyDate, GETDATE()), ISNULL(sNotify, ''), ISNULL(sIntro, '')
		 , T1.nCharNo, T1.sID, ISNULL(T3.nGraduateCount, 0) AS nGraduateCount
	  FROM tGuildAcademy AS GA INNER JOIN 
		   tCharacter AS C1 ON GA.nMasterCharNo = C1.nCharNo INNER JOIN (
				SELECT G.nNo, GM.nCharNo, C.sID
				  FROM tGuild AS G, tGuildMember AS GM, tCharacter AS C
				 WHERE G.nNo = GM.nNo AND GM.nCharNo = C.nCharNo
				   AND GM.nGrade = 0 AND C.bDeleted = 0
				) AS T1 ON GA.nNo = T1.nNo INNER JOIN (
				-- ranking check
				SELECT T1.nNo, COUNT(T2.nPoint) + 1 AS nRanking
				  FROM tGuildAcademy T1 LEFT OUTER JOIN tGuildAcademy T2 ON T1.nPoint < T2.nPoint
				 GROUP BY T1.nNo, T1.nPoint
				) AS T2 ON GA.nNo = T2.nNo LEFT OUTER JOIN (
				-- graduate counting
				SELECT nNo, ISNULL(COUNT(nCharNo), 0) AS nGraduateCount
				  FROM tGuildAcademyGraduate
				 GROUP BY nNo
				) AS T3 ON GA.nNo = T3.nNo
	 ORDER BY T2.nRanking
	IF @@ERROR <> 0
		BEGIN
			SET @nRet = -1594
		END

	SET NOCOUNT OFF
END

0㄁黡x
GO

-- ========================================================================
-- Procedure: usp_GuildAcademy_getGuildNo
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_GuildAcademy_getGuildNo]
	@nCharNo int
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	IF @nCharNo <= 0 BEGIN
		SET @nRet = -1512
		RETURN
	END

	IF EXISTS(SELECT nNo FROM tGuildAcademyMember WHERE nCharNo = @nCharNo) BEGIN
		SELECT nNo FROM tGuildAcademyMember WHERE nCharNo = @nCharNo
	END
	ELSE BEGIN
		-- 2008.8.28 일 업데이트시 서버프로그램의 수정 실수로 인해
		-- 정상적인 처리가 되지 않는 현상이 발생함.
		-- 2008.9.4 일 업데이트시 정상적으로 진행할 예정임.
		-- SET @nRet = -1521
		SET @nRet = 0
	END

	SET NOCOUNT OFF
END

0樁錅y
GO

-- ========================================================================
-- Procedure: usp_GuildAcademy_getHistoryList
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_GuildAcademy_getHistoryList]

@nAcademyID int,
@nAmount int,
@nRet int output

AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	SELECT TOP (@nAmount) sID, nAction, dDate FROM tGuildAcademyHistory WHERE nAcademyID = @nAcademyID

	SET NOCOUNT OFF
END
GO

-- ========================================================================
-- Procedure: usp_GuildAcademy_getMemberList
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_GuildAcademy_getMemberList]
	@nNo  int
,	@nRet int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	IF @nNo <= 0 BEGIN
		SET @nRet = -1511
		RETURN
	END

	-- add G.nChatBanTime
	SELECT G.nCharNo, G.nJoinLevel, G.dJoinDate, G.nGrade, G.sIntro, G.nChatBanTime
		 , C.sID, C.nLevel, C.sLoginZone, C.dLastLoginDate
		 , S.nClass
	  FROM tGuildAcademyMember AS G INNER JOIN 
		   tCharacter          AS C ON C.nCharNo =  G.nCharNo INNER JOIN 
		   tCharacterShape     AS S ON S.nCharNo =  G.nCharNo
	 WHERE G.nNo = @nNo AND C.bDeleted = 0

	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: usp_GuildAcademy_getMoney
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_GuildAcademy_getMoney]
	@nNo  int
,	@nRet int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	IF EXISTS(SELECT nNo FROM tGuildAcademy WHERE nNo = @nNo) BEGIN
		SELECT nMoney FROM tGuildAcademy WHERE nNo = @nNo
	END
	ELSE BEGIN
		SET @nRet = -1521
	END

	SET NOCOUNT OFF
END

0ꌁ蜩z
GO

-- ========================================================================
-- Procedure: usp_GuildAcademy_getRanking
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_GuildAcademy_getRanking]
	@nCount int = 1
,	@nRet   int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	DECLARE @MAX_COUNT AS INT

	SELECT @MAX_COUNT = ISNULL(MAX(nCount), 0)
	  FROM tGuildAcademyRanking
	
	SELECT dDate, nRanking, nGuildNo, sGuildName, nMasterNo, sMasterID, nPoint
	  FROM tGuildAcademyRanking
	 WHERE nCount = @MAX_COUNT - (@nCount - 1)

	SET NOCOUNT OFF
END

0ᔁ潲|
GO

-- ========================================================================
-- Procedure: usp_GuildAcademy_getReward
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_GuildAcademy_getReward]
	@nNo    int
,	@nLevel tinyint
,	@nRet   int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	IF @nNo <= 0 BEGIN
		SET @nRet = -1511
		RETURN
	END

	IF @nLevel <= 0 BEGIN
		SET @nRet = -1513
		RETURN
	END

	SELECT nTypeID, nValue, nItemID FROM tGuildAcademyReward WHERE nNo = @nNo AND nLevel = @nLevel
	IF @@ERROR <> 0 BEGIN
		SET @nRet = -1591
	END	

	SET NOCOUNT OFF
END

0丁掖}
GO

-- ========================================================================
-- Procedure: usp_GuildAcademy_getRewardAll
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_GuildAcademy_getRewardAll]
	@nNo     int
,	@nTypeID int
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	IF @nNo <= 0 BEGIN
		SET @nRet = -1511
		RETURN
	END

	IF @nTypeID = 1 BEGIN
		SELECT nLevel, nTypeID, nValue, nItemID FROM tGuildAcademyReward WHERE nNo = @nNo AND nTypeID = 1
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -1591
		END	
	END 
	ELSE BEGIN
		SELECT nLevel, nTypeID, nValue, nItemID FROM tGuildAcademyReward WHERE nNo = @nNo AND nTypeID > 1
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -1591
		END	
	END

	SET NOCOUNT OFF
END

0蜁垺~
GO

-- ========================================================================
-- Procedure: usp_GuildAcademy_raiseMember
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_GuildAcademy_raiseMember]
	@nNo     int
,	@nCharNo int
,	@nGrade  tinyint
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	DECLARE @nMax AS INT
	DECLARE @nCnt AS INT

	SET @nRet = 0

	IF EXISTS(SELECT nNo FROM tGuildAcademyMember WHERE nNo = @nNo AND nCharNo = @nCharNo) BEGIN
		SELECT @nMax = nMaxMembers, @nCnt = nCnt
		  FROM tGuild AS T1 INNER JOIN (
				SELECT nNo, ISNULL(COUNT(*), 0) AS nCnt
				  FROM tCharacter AS C, tGuildMember AS G
				 WHERE C.nCharNo = G.nCharNo 
				   AND G.nNo = @nNo
				   AND C.bDeleted = 0
				 GROUP BY nNo) AS T2 ON (T1.nNo = T2.nNo)
		IF @nCnt < @nMax BEGIN
			BEGIN TRAN
	
			-- 길드아카데미 탈퇴(0:ok)
			EXEC usp_GuildAcademy_delMember @nNo, @nCharNo, @nRet output
			IF @nRet <> 0 BEGIN
				ROLLBACK TRAN
				RETURN
			END
	
			-- 길드 가입(1:ok)
			EXEC p_Guild_Member_Join @nNo, @nCharNo, @nGrade, @nRet output
			IF @nRet <> 1 BEGIN
				SET @nRet = -1531
				ROLLBACK TRAN
				RETURN
			END
			ELSE BEGIN
				SET @nRet = 0
				COMMIT TRAN
			END
		END
		ELSE BEGIN
			SET @nRet = -1530
		END
	END
	ELSE BEGIN
		SET @nRet = -1528
	END
END

0쀁䯞
GO

-- ========================================================================
-- Procedure: usp_GuildAcademy_setAcademy
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_GuildAcademy_setAcademy]
	@nNo           int
,	@nMasterCharNo int
,	@nMaxMember    smallint
,	@nRet          int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	IF EXISTS(SELECT nNo FROM tGuildAcademy WHERE nNo = @nNo) BEGIN
		SET @nRet = -1521
		RETURN
	END
	ELSE BEGIN
		INSERT INTO tGuildAcademy (nNo, nMasterCharNo, nMaxMember)
		VALUES (@nNo, @nMasterCharNo, @nMaxMember)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -1592
			RETURN
		END
	END
END

0更䀂
GO

-- ========================================================================
-- Procedure: usp_GuildAcademy_setChatBanTime
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_GuildAcademy_setChatBanTime]
	@nNo          int
,	@nCharNo      int
,	@nChatBanTime int
,	@nRet         int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0
	
	IF EXISTS(SELECT nNo FROM tGuildAcademyMember WHERE nNo = @nNo AND nCharNo = @nCharNo) BEGIN
		UPDATE tGuildAcademyMember 
		   SET nChatBanTime = @nChatBanTime
		 WHERE nNo = @nNo AND nCharNo = @nCharNo
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = -1593
			RETURN
		END
	END 
	ELSE BEGIN
		SET @nRet = -1511
		RETURN
	END

	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: usp_GuildAcademy_setGraduate
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_GuildAcademy_setGraduate]
	@nNo            int
,	@nCharNo        int
,	@sID            nvarchar(40)
,	@nClass         tinyint
,	@nGraduateLevel smallint
,	@nRet           int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	DECLARE @nJoinLevel smallint
	DECLARE @dJoinDate  datetime

	SET @nRet = 0
	SET @nJoinLevel = 0

	IF @nNo <= 0 BEGIN
		SET @nRet = -1511
		RETURN
	END

	IF @nCharNo <= 0 BEGIN
		SET @nRet = -1512
		RETURN
	END

	IF LEN(ISNULL(@sID, '')) = 0 BEGIN
		SET @nRet = -1515
		RETURN
	END

	IF @nClass <= 0 BEGIN
		SET @nRet = -1516
		RETURN
	END

	IF @nGraduateLevel <= 0 BEGIN
		SET @nRet = -1517
		RETURN
	END

	SELECT @nJoinLevel = nJoinLevel, @dJoinDate = dJoinDate
	  FROM tGuildAcademyMember
	 WHERE nNo = @nNo AND nCharNo = @nCharNo
	IF @nJoinLevel <= 0 BEGIN
		SET @nRet = -1528
		RETURN
	END
	ELSE BEGIN
		INSERT INTO tGuildAcademyGraduate (nNo, nCharNo, sID, nClass, nJoinLevel, dJoinDate, nGraduateLevel)
		VALUES (@nNo, @nCharNo, @sID, @nClass, @nJoinLevel, @dJoinDate, @nGraduateLevel)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -1592
			RETURN
		END
	END

	SET NOCOUNT OFF
END
GO

-- ========================================================================
-- Procedure: usp_GuildAcademy_setHistory
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_GuildAcademy_setHistory]
/*
Set Guild Token
2008.12 By Lunar
Input:
    See SQL
    
Output:
    nRet = -1,  Error(Not found guild)
    nRet = -2,  Error(Not found guild)
    nRet = 1,  OK
*/
@nAcademyID int,
@nCharNo int,
@sID nvarchar(50),
@nAction int,
@nRet int output

AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	INSERT INTO tGuildAcademyHistory (nAcademyID, nCharNo, sID, dDate, nAction)
	VALUES (@nAcademyID, @nCharNo, @sID, CURRENT_TIMESTAMP, @nAction)

	IF @@ERROR <> 0
	BEGIN
		SET @nRet = -2
		RETURN
	END

	SET NOCOUNT OFF
END
GO

-- ========================================================================
-- Procedure: usp_GuildAcademy_setMaster
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_GuildAcademy_setMaster]
	@nNo       int
,	@nCharNo   int
,	@nRet      int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0
	
	IF EXISTS(SELECT nNo FROM tGuildMember WHERE nNo = @nNo AND nCharNo = @nCharNo) BEGIN
		UPDATE tGuildAcademy SET nMasterCharNo = @nCharNo
		 WHERE nNo = @nNo
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = -1593
			RETURN
		END
	END 
	ELSE BEGIN
		SET @nRet = -1511
		RETURN
	END

	SET NOCOUNT OFF
END

0ꐁᱯ
GO

-- ========================================================================
-- Procedure: usp_GuildAcademy_setMember
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_GuildAcademy_setMember]
	@nNo       int
,	@nCharNo   int
,	@dJoinDate datetime
,	@nGrade    tinyint = 0
,	@nRet      int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	DECLARE @nLevel AS int

	IF @nNo <= 0 BEGIN
		SET @nRet = -1511
		RETURN
	END

	IF @nCharNo <= 0 BEGIN
		SET @nRet = -1512
		RETURN
	END

	IF EXISTS(SELECT nNo FROM tGuildMember WHERE nNo = @nNo AND nCharNo = @nCharNo) BEGIN
		SET @nRet = -1522
		RETURN
	END

	IF EXISTS(SELECT nCharNo FROM tCharacter WHERE nCharNo = @nCharNo AND bDeleted = 0) BEGIN
		SELECT @nLevel = nLevel FROM tCharacter WHERE nCharNo = @nCharNo
		IF @nLevel <= 0 BEGIN
			SET @nRet = -1523
			RETURN
		END
	END
	ELSE BEGIN
		SET @nRet = -1524
		RETURN
	END
	
	INSERT tGuildAcademyMember (nNo, nCharNo, nJoinLevel, dJoinDate, nGrade)
	VALUES ( @nNo, @nCharNo, @nLevel, @dJoinDate, @nGrade )
	IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
		SET @nRet = -1591
	END

	SET NOCOUNT OFF
END
GO

-- ========================================================================
-- Procedure: usp_GuildAcademy_setMemberLevelup
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_GuildAcademy_setMemberLevelup]
	@nNo       int
,	@nPoint    int
,	@nBuffTime int
,	@nRet      int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	DECLARE @BUFF_END_TIME AS DATETIME

	IF EXISTS(SELECT nNo FROM tGuildAcademy WHERE nNo = @nNo) BEGIN
		BEGIN TRAN

		SELECT @BUFF_END_TIME = ISNULL(dBuffEndTime, GETDATE()) FROM tGuildAcademy WHERE nNo = @nNo
		IF @BUFF_END_TIME > GETDATE() BEGIN
			SET @BUFF_END_TIME = DATEADD ( minute, @nBuffTime, @BUFF_END_TIME ) 
		END
		ELSE BEGIN
			SET @BUFF_END_TIME = DATEADD ( minute, @nBuffTime, GETDATE() ) 
		END

		UPDATE tGuildAcademy SET nPoint = nPoint + @nPoint, dBuffEndTime = @BUFF_END_TIME
		 WHERE nNo = @nNo
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = -1593
			ROLLBACK TRAN
			RETURN
		END
		ELSE BEGIN			
			COMMIT TRAN
		END
	END
	ELSE BEGIN
		SET @nRet = -1521
	END	

	SET NOCOUNT OFF
END

0ᘁҸ
GO

-- ========================================================================
-- Procedure: usp_GuildAcademy_setNotify
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_GuildAcademy_setNotify]
	@nNo     int
,	@sNotify nvarchar(1024)
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	IF @nNo <= 0 BEGIN
		SET @nRet = -1511
		RETURN
	END

	UPDATE tGuildAcademy 
	   SET dNotifyDate = getdate(), sNotify = @sNotify 
	 WHERE nNo = @nNo
	IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
		SET @nRet = -1593
	END
	ELSE BEGIN
		SET @nRet = 0
	END

	SET NOCOUNT OFF
END

0企
GO

-- ========================================================================
-- Procedure: usp_GuildAcademy_setRanking
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_GuildAcademy_setRanking]
	@nRet  int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	DECLARE @MAX_COUNT AS INT

	BEGIN TRAN

	SELECT @MAX_COUNT = ISNULL(MAX(nCount), 0) + 1
	  FROM tGuildAcademyRanking

	INSERT INTO tGuildAcademyRanking (dDate, nRanking, nGuildNo, sGuildName, nMasterNo, sMasterID, nPoint, nCount)
		SELECT GETDATE() AS dDate, R.nRanking, R.nNo, G.sName, GA.nMasterCharNo, C.sID, R.nPoint, @MAX_COUNT
		  FROM tGuild AS G, tGuildAcademy AS GA, tCharacter AS C, (
				SELECT COUNT(T2.nPoint) + 1 AS nRanking, T1.nNo, T1.nPoint
				  FROM tGuildAcademy T1 LEFT OUTER JOIN tGuildAcademy T2 ON T1.nPoint < T2.nPoint
				 GROUP BY T1.nNo, T1.nPoint) AS R
		 WHERE G.nNo = R.nNo AND C.nCharNo = GA.nMasterCharNo AND GA.nNo = R.nNo
	IF @@ERROR <> 0 BEGIN
		SET @nRet = -1592
		ROLLBACK TRAN
		RETURN
	END

	UPDATE tGuildAcademy SET nPoint = 0
	IF @@ERROR <> 0 OR @@ROWCOUNT = 0 BEGIN
		SET @nRet = -1593
		ROLLBACK TRAN
		RETURN
	END
	ELSE BEGIN
		SET @nRet = 0
		COMMIT TRAN
	END

	SET NOCOUNT OFF
END

0蠁
GO

-- ========================================================================
-- Procedure: usp_GuildAcademy_setReward
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_GuildAcademy_setReward]
	@nNo     int
,	@nLevel  smallint
,	@nTypeID tinyint
,	@nValue  bigint
,	@nItemID int = NULL
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	IF @nLevel <= 0 BEGIN
		SET @nRet = -1513
		RETURN
	END

	INSERT INTO tGuildAcademyReward (nNo, nLevel, nTypeID, nValue, nItemID)
	VALUES (@nNo, @nLevel, @nTypeID, @nValue, @nItemID)
	IF @@ERROR <> 0 BEGIN
		SET @nRet = -1595
		RETURN
	END

	SET NOCOUNT OFF
END

0섁
GO

-- ========================================================================
-- Procedure: usp_GuildAcademy_withdrawMoney
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_GuildAcademy_withdrawMoney]
	@nNo     int
,	@nCharNo int
,	@nMoney  bigint
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	DECLARE @money AS bigint
	SET @nRet = 0
	SET @money = 0

	IF @nNo <= 0
		BEGIN
			SET @nRet = -1511
			RETURN
		END

	IF @nCharNo <= 0
		BEGIN
			SET @nRet = -1512
			RETURN
		END

	IF @nMoney <= 0
		BEGIN
			SET @nRet = -1514
			RETURN
		END

	BEGIN TRAN
	
	IF NOT EXISTS(SELECT nCharNo FROM tCharacter WHERE nCharNo = @nCharNo AND bDeleted = 0)
		BEGIN
			SET @nRet = -1524
			ROLLBACK TRAN
			RETURN
		END

	IF NOT EXISTS(SELECT nNo FROM tGuildAcademy WHERE nNo = @nNo)
		BEGIN
			SET @nRet = -1525
			ROLLBACK TRAN
			RETURN
		END

	SELECT @money = ISNULL(nMoney, 0) FROM tGuildAcademy WITH(UPDLOCK, HOLDLOCK) WHERE nNo = @nNo
	IF @money < @nMoney 
		BEGIN		
			SET @nRet = -1527
			ROLLBACK TRAN
			RETURN
		END
	ELSE
		BEGIN
			UPDATE tGuildAcademy SET nMoney = nMoney - @nMoney WHERE nNo = @nNo AND nMoney >= @nMoney
			IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
				BEGIN
					SET @nRet = -1593
					ROLLBACK TRAN
					RETURN
				END

			UPDATE tCharacter SET nMoney = nMoney + @nMoney WHERE nCharNo = @nCharNo
			IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
				BEGIN
					SET @nRet = -1593
					ROLLBACK TRAN
					RETURN
				END
			ELSE
				BEGIN
					SET @nRet = 0
					COMMIT TRAN
				END
		END
			
	SET NOCOUNT OFF
END
GO

-- ========================================================================
-- Procedure: usp_GuildEmblem_GetInfo
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_GuildEmblem_GetInfo]

@nGuildNo INT

AS
BEGIN
	SET NOCOUNT ON
	SELECT nColor, nIcon, dCreated, dUpdated, nUnk FROM tGuildEmblem
	WHERE nGuildNo = @nGuildNo
	SET NOCOUNT OFF
END
GO

-- ========================================================================
-- Procedure: usp_GuildEmblem_GetInfoAll
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_GuildEmblem_GetInfoAll]

AS
BEGIN
	SET NOCOUNT ON
	SELECT nGuildNo, nColor, nIcon, dCreated, dUpdated, nUnk FROM tGuildEmblem
	SET NOCOUNT OFF
END0
GO

-- ========================================================================
-- Procedure: usp_GuildEmblem_SetInfo
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_GuildEmblem_SetInfo]
/*
Set Guild Token
2008.12 By Lunar
Input:
	See SQL
	
Output:
	nRet = -1,  Error(Not found guild)
	nRet = -2,  Error(Not found guild)
	nRet = 1,  OK
*/
@nGuildNo int,
@nColor int,
@nIcon int,
@dCreated datetime,
@dUpdated datetime,
@nUnk int,
@nRet int output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	IF EXISTS(SELECT nGuildNo FROM tGuildEmblem WHERE nGuildNo = @nGuildNo) BEGIN
		UPDATE tGuildEmblem
		SET nColor = @nColor, nIcon = @nIcon, dUpdated = @dUpdated, nUnk = @nUnk
		WHERE nGuildNo = @nGuildNo
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = -1
			RETURN
		END
	END
	ELSE BEGIN
		INSERT INTO tGuildEmblem (nGuildNo, nColor, nIcon, dCreated, dUpdated, nUnk)
		VALUES (@nGuildNo, @nColor, @nIcon, @dCreated, @dUpdated, @nUnk)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -2
			RETURN
		END
	END			
	
	IF @@ERROR <> 0
	BEGIN
		ROLLBACK TRAN
		RETURN
	END
	RETURN @@ERROR	

	SET NOCOUNT OFF
END
-- end
GO

-- ========================================================================
-- Procedure: usp_GuildTournament_getChoiceList
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_GuildTournament_getChoiceList]
	@nGTNo int
,	@nRet  int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0
	IF ISNULL(@nGTNo, 0) <= 0 BEGIN
		SET @nRet = -1411
		RETURN
	END

	SELECT nNo FROM tGuildTournamentRegist WITH(NOLOCK) WHERE nGTNo = @nGTNo AND bChoice = 1

	SET NOCOUNT OFF
END

0쀁➥ċ
GO

-- ========================================================================
-- Procedure: usp_GuildTournament_getMoney
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_GuildTournament_getMoney]
	@nGTNo int
,	@nRet  int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0
	IF ISNULL(@nGTNo, 0) <= 0 BEGIN
		SET @nRet = -1411
		RETURN
	END

	SELECT ISNULL(SUM(nMoney), 0) AS nTotalMoney 
	  FROM tGuildTournamentRegist WITH(NOLOCK) 
	 WHERE nGTNo = @nGTNo

	SET NOCOUNT OFF
END

0㌁쥭
GO

-- ========================================================================
-- Procedure: usp_GuildTournament_getRegistList
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_GuildTournament_getRegistList]
	@nGTNo int
,	@nRet  int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0
	IF ISNULL(@nGTNo, 0) <= 0 BEGIN
		SET @nRet = -1411
		RETURN
	END

	SELECT nNo FROM tGuildTournamentRegist WITH(NOLOCK) WHERE nGTNo = @nGTNo

	SET NOCOUNT OFF
END

0ꔁ놵
GO

-- ========================================================================
-- Procedure: usp_GuildTournament_setChoice
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_GuildTournament_setChoice]
	@nGTNo        int
,	@nLimitLevel  int = 1
,	@nLimitMember int = 1
,	@nRet         int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0
	IF ISNULL(@nGTNo, 0) <= 0 BEGIN
		SET @nRet = -1411
		RETURN
	END

	BEGIN TRAN

	UPDATE tGuildTournamentRegist
	   SET bChoice = 1
	 WHERE nGTNo = @nGTNo AND nNo IN (
			SELECT TOP 16 T4.nNo
			  FROM tGuildTournamentRegist AS T4 INNER JOIN (
					SELECT T1.nNo, ISNULL(COUNT(T2.nCharNo), 0) AS nMember
					  FROM tGuildMember AS T1, tCharacter AS T2 
					 WHERE T1.nCharNo = T2.nCharNo
					   AND T2.nLevel >= @nLimitLevel
					 GROUP BY T1.nNo
					HAVING COUNT(T2.nCharNo) >= @nLimitMember) AS T3 ON T4.nNo = T3.nNo
			 WHERE nGTNo = @nGTNo
			 ORDER BY nSeed DESC)
	IF @@ERROR <> 0 BEGIN
		SET @nRet = -1493
		ROLLBACK TRAN
	END
	ELSE BEGIN
		SET @nRet = 0
		COMMIT TRAN
	END

	SET NOCOUNT OFF
END

0ꗙ
GO

-- ========================================================================
-- Procedure: usp_GuildTournament_setRegist
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_GuildTournament_setRegist]
	@nGTNo        int
,	@nNo          int
,	@nCharNo      int
,	@nMoney       bigint = 0
,	@nLimitLevel  int = 1
,	@nLimitMember int = 1
,	@nRet         int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	DECLARE @prevGTNo int, @prevWinnerNo int
	DECLARE @cnt AS int

	SET @nRet = 0
	SET @cnt  = 0

	IF ISNULL(@nGTNo, 0) <= 0 BEGIN
		SET @nRet = -1411
		RETURN
	END

	IF ISNULL(@nNo, 0) <= 0 BEGIN
		SET @nRet = -1412
		RETURN
	END

	IF ISNULL(@nCharNo, 0) <= 0 BEGIN
		SET @nRet = -1413
		RETURN
	END

	IF EXISTS(SELECT nGTNo FROM tGuildTournamentRegist WHERE nGTNo = @nGTNo AND nNo = @nNo) BEGIN
		SET @nRet = -1421
		RETURN
	END

	-- �������ʸ�Ʈ ��û���� üũ(����, �ο�)
	SELECT @cnt = ISNULL(COUNT(T1.nCharNo), 0)
	  FROM tGuildMember AS T1, tCharacter AS T2 
	 WHERE T1.nCharNo = T2.nCharNo
	   AND T1.nNo = @nNo
	   AND T2.nLevel >= @nLimitLevel
	IF @cnt < @nLimitMember BEGIN
		SET @nRet = -1422
		RETURN
	END

	-- 2011.8.22 by kodong
	-- �������ʸ�Ʈ ��û�� ���� �������ʸ�Ʈ ���±����� ��û�� �� ����.	
	SELECT TOP 1 @prevGTNo = nGTNo 
	  FROM tGuildTournament 
	 WHERE nGTNo < @nGTNo 
	 ORDER BY nGTNo DESC
	IF @prevGTNo IS NOT NULL BEGIN
		SELECT @prevWinnerNo = nGuildNo
		  FROM tGuildTournament 
		 WHERE nGTNo = @prevGTNo
		   AND nType >= 3028 AND nType <= 3029
		   AND nData = 1
		IF @prevWinnerNo = @nNo BEGIN
			SET @nRet = -1424
			RETURN
		END
	END
	
	BEGIN TRAN

	-- �������ʸ�Ʈ ��û�� �
	EXEC p_Char_MoneySub @nCharNo, @nMoney, @nRet OUTPUT
	IF @nRet <> 1 BEGIN
		SET @nRet = -1423
		ROLLBACK TRAN
		RETURN
	END

	INSERT INTO tGuildTournamentRegist (nGTNo, nNo, nMoney, nSeed)
	VALUES (@nGTNo, @nNo, @nMoney, CAST(RAND() * 100000 AS INT))
	IF @@ERROR <> 0 BEGIN
		SET @nRet = -1492
		ROLLBACK TRAN
	END
	ELSE BEGIN
		SET @nRet = 0
		COMMIT TRAN
	END

	SET NOCOUNT OFF
END
GO

-- ========================================================================
-- Procedure: usp_GuildTournament_setResult
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_GuildTournament_setResult]
	@nGTNo    int
,	@nType    smallint
,	@nData    int
,	@nGuildNo int
,	@nRet     int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	SET @nRet = 0
	
	IF EXISTS(SELECT nGTNo FROM tGuildTournament WHERE  nGTNo = @nGTNo AND nType = @nType) BEGIN
		UPDATE tGuildTournament 
		   SET nData = @nData 
		     , nGuildNo = @nGuildNo 
		 WHERE nGTNo = @nGTNo 
		   AND nType = @nType
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = -1493
		END		   
	END
	ELSE BEGIN
		INSERT tGuildTournament (nGTNo, nType, nData, nGuildNo)
		VALUES (@nGTNo, @nType, @nData, @nGuildNo)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -1492
		END			
	END
END

0ď
GO

-- ========================================================================
-- Procedure: usp_Guild_GetToken
-- ========================================================================
CREATE  PROCEDURE [dbo].[usp_Guild_GetToken]
/*
Get Guild Token
2008.12 By Lunar
Input:
	See SQL
	
Output:
	nRet = 0,  Error(Not found guild)
	nRet = 1,  OK
*/
@nNo int,
@nToken bigint output,
@nRet int output
AS
SET NOCOUNT ON
SET @nRet = 0
SELECT @nToken = ISNULL(MAX(nToken), 0) FROM tGuildToken WHERE nNo = @nNo
IF @@ERROR <> 0 --OR @@ROWCOUNT <> 1
BEGIN
	SET @nRet = -1
	RETURN
END
-- end

0
GO

-- ========================================================================
-- Procedure: usp_Guild_GetTokenAll
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Guild_GetTokenAll]
AS
BEGIN
	SET NOCOUNT ON
	SELECT nNo, nToken FROM tGuildToken
	SET NOCOUNT OFF
END

0、ছp
GO

-- ========================================================================
-- Procedure: usp_Guild_SetToken
-- ========================================================================
CREATE  PROCEDURE [dbo].[usp_Guild_SetToken]
/*
Set Guild Token
2008.12 By Lunar
Input:
	See SQL
	
Output:
	nRet = -1,  Error(Not found guild)
	nRet = -2,  Error(Not found guild)
	nRet = 1,  OK
*/
@nNo int,
@nToken bigint,
@nRet int output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	IF EXISTS(SELECT nNo FROM tGuildToken WHERE nNo = @nNo) BEGIN
		UPDATE tGuildToken SET nToken = @nToken WHERE nNo = @nNo
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = -1
			RETURN
		END
	END
	ELSE BEGIN
		INSERT INTO tGuildToken (nNo, nToken)
		VALUES (@nNo, @nToken)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -2
			RETURN
		END
	END			
	
	SET NOCOUNT OFF
END
-- end

0
GO

-- ========================================================================
-- Procedure: usp_Guild_TokenSub
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Guild_TokenSub]
/*
Sub money
2008.12 By Lunar
 input:
	see SQL
 output:
	@nRet	1 = OK
	        0 = Error, or Minus value
*/
@nGuildNo int,
@nSubToken bigint,
@nRet int OUTPUT
AS
SET NOCOUNT ON
DECLARE @nGuildToken bigint
EXEC usp_Guild_GetToken @nGuildNo, @nGuildToken OUTPUT, @nRet OUTPUT
IF @nRet <> 0
	RETURN
-- Sub money
SET @nGuildToken = @nGuildToken - @nSubToken
-- Check minus
IF @nGuildToken < 0
BEGIN
	SET @nRet = -1
	RETURN	
END
EXEC usp_Guild_SetToken @nGuildNo, @nGuildToken, @nRet OUTPUT
-- End

0
GO

-- ========================================================================
-- Procedure: usp_HolyPromise_getRemainMoney
-- ========================================================================
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

-- ========================================================================
-- Procedure: usp_HolyPromise_setMoney
-- ========================================================================
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

-- ========================================================================
-- Procedure: usp_HolyPromise_withDrawRemainMoney
-- ========================================================================
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

-- ========================================================================
-- Procedure: usp_InGameGuide_getViewList
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_InGameGuide_getViewList]
	@nCharNo int
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	SET @nRet = 0

	IF EXISTS(SELECT nCharNo FROM tInGameGuide WITH(NOLOCK) WHERE nCharNo = @nCharNo) BEGIN
		SELECT nGuideNo
		  FROM tInGameGuide 
		 WHERE nCharNo = @nCharNo
	END
	ELSE BEGIN
		SET @nRet = -1111
	END
END

0㐁庳
GO

-- ========================================================================
-- Procedure: usp_InGameGuide_setGuideNo
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_InGameGuide_setGuideNo]
	@nCharNo  int
,	@nGuideNo int
,	@nRet     int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	SET @nRet = 0

	IF EXISTS(SELECT nCharNo FROM tInGameGuide WITH(NOLOCK) WHERE nCharNo = @nCharNo AND nGuideNo = @nGuideNo) BEGIN
		SET @nRet = -1111
	END
	ELSE BEGIN
		INSERT INTO tInGameGuide (nCharNo, nGuideNo)
		VALUES (@nCharNo, @nGuideNo)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -2222
		END
	END
END

0洁勗
GO

-- ========================================================================
-- Procedure: usp_Item_Add_EnchantGem
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Item_Add_EnchantGem] 
	@nCharNo INT,
	@nEquipmentItemKey BIGINT,
	@nGemItemKey BIGINT,
	@nGemSocketIndex TINYINT,
	@nGemID INT,
	@nRestGemCount TINYINT,
	@nRet INT = 0 OUTPUT
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON
	
	SET @nRet = 0
	
    BEGIN TRAN
    
    INSERT tItemEnchant (nItemKey, nSlot, nGemItemID, nRestCount)
		VALUES (@nEquipmentItemKey, @nGemSocketIndex, @nGemID, @nRestGemCount)

	IF @@ERROR <> 0
	BEGIN
		ROLLBACK TRAN
		SET @nRet = -101
		RETURN
	END

	EXEC p_Item_Delete @nGemItemKey, @nRet OUTPUT
	
	IF @@ERROR <> 0 OR @nRet <> 0
	BEGIN
		ROLLBACK TRAN
		SET @nRet = -202
		RETURN
	END

    COMMIT TRAN 
END

0
GO

-- ========================================================================
-- Procedure: usp_Item_Add_EnchantSocket
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Item_Add_EnchantSocket]
	@nCharNo INT,
	@nEquipmentItemKey BIGINT,
	@nDrillItemKey BIGINT,
	@nAllSocketCount TINYINT OUTPUT,
	@nDrillLotType SMALLINT,
	@nDrillLotValue BIGINT OUTPUT,
	@nRet INT = 0 OUTPUT
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON
	
	SET @nRet = 0

    BEGIN TRAN

	EXEC p_Item_Delete @nDrillItemKey, @nRet OUTPUT

	IF @@ERROR <> 0 OR @nRet <> 0
	BEGIN
		ROLLBACK TRAN
		SET @nRet = -101
		RETURN
	END
	
	SET @nDrillLotValue = 0
	
	DECLARE @nCreatedSocketCount TINYINT
	
	EXEC p_Item_GetOption @nEquipmentItemKey, 523, @nAllSocketCount OUTPUT, @nRet OUTPUT

	IF @@ERROR <> 0 OR @nRet = 0
	BEGIN
		ROLLBACK TRAN
		SET @nRet = -201
		RETURN
	END

	EXEC p_Item_GetOption @nEquipmentItemKey, 524, @nCreatedSocketCount OUTPUT, @nRet OUTPUT

	IF @@ERROR <> 0 OR @nRet = 0
	BEGIN
		ROLLBACK TRAN
		SET @nRet = -202
		RETURN
	END
	
	SET @nAllSocketCount = @nAllSocketCount + 1
	SET @nCreatedSocketCount = @nCreatedSocketCount + 1
	
	EXEC p_Item_SetOption @nEquipmentItemKey, 523, @nAllSocketCount, @nRet OUTPUT

	IF @@ERROR <> 0 OR @nRet <> 0
	BEGIN
		ROLLBACK TRAN
		SET @nRet = -301
		RETURN
	END
	
	EXEC p_Item_SetOption @nEquipmentItemKey, 524, @nCreatedSocketCount, @nRet OUTPUT
	
	IF @@ERROR <> 0 OR @nRet <> 0
	BEGIN
		ROLLBACK TRAN
		SET @nRet = -302
		RETURN
	END
	
    COMMIT TRAN  
END

0
GO

-- ========================================================================
-- Procedure: usp_Item_Del_EnchantGem
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Item_Del_EnchantGem] 
	@nCharNo INT,
	@nEquipmentItemKey BIGINT,
	@nRemoverItemKey BIGINT,
	@nResetSocketIndex TINYINT,
	@bNeedItem TINYINT,
	@nRemoverLotType SMALLINT,
	@nRemoverLotValue BIGINT = 0 OUTPUT,
	@nRet INT = 0 OUTPUT
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON
	
	SET @nRet = 0
	SET @nRemoverLotValue = 0
	
    BEGIN TRAN

	IF @bNeedItem = 1
	BEGIN
		EXEC p_Item_Delete @nRemoverItemKey, @nRet OUTPUT

		IF @@ERROR <> 0 OR @nRet <> 0
		BEGIN
			ROLLBACK TRAN
			SET @nRet = -101
			RETURN
		END
	END
	
    DELETE FROM tItemEnchant
		WHERE nItemKey = @nEquipmentItemKey
			AND nSlot = @nResetSocketIndex

	IF @@ERROR <> 0 OR @@ROWCOUNT = 0
	BEGIN
		ROLLBACK TRAN
		SET @nRet = -201
		RETURN
	END

    COMMIT TRAN
END

0
GO

-- ========================================================================
-- Procedure: usp_Item_Get_EnchantAll
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Item_Get_EnchantAll]
	@nEquipmentItemKey BIGINT,
	@nRet INT = 0 OUTPUT
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON

	SET @nRet = 0
	
	SELECT nSlot, nGemItemID, nRestCount FROM tItemEnchant WHERE nItemKey = @nEquipmentItemKey

	IF @@ERROR <> 0
	BEGIN
		SET @nRet = -101
		RETURN
	END 
END

0
GO

-- ========================================================================
-- Procedure: usp_Item_GuildToken_BuyAll
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Item_GuildToken_BuyAll]
/*
Item buy all
2004.8 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK
	        1 = Not found item owner or Need more owner's money/fame
	        ? = Error code
	        
 operation:
	1. money/fame check and sub
	2. make buy item
	
*/
@nItemKey bigint,
@nOwner int,
@nStorageType tinyint,
@nStorage smallint,
@nItemID int,
@nFlags int,
@nPriceGuildToken bigint,
@nRet int OUTPUT
AS
SET NOCOUNT ON
DECLARE @CreateItemRet int
BEGIN TRAN

-- 아이템 소유자가 충분한 돈을 가지고 있는지 확인 후 돈 감소

EXEC usp_Guild_TokenSub @nOwner, @nPriceGuildToken, @nRet OUTPUT
IF @nRet <> 0
BEGIN
	ROLLBACK TRAN
	SET @nRet = 1
	RETURN
END


-- 아이템 생성 저장프로시져 호출
EXEC p_Item_Create @nItemKey, @nOwner, @nStorageType, @nStorage, @nItemID, @nFlags, @CreateItemRet OUTPUT
IF @CreateItemRet <> 0 OR @@ERROR <> 0
BEGIN
	ROLLBACK TRAN
	SET @nRet = @CreateItemRet
	RETURN
END
-- 완료
COMMIT TRAN
SET @nRet = 0
-- end

0ﰁￔ
GO

-- ========================================================================
-- Procedure: usp_Item_GuildToken_BuyLot
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Item_GuildToken_BuyLot]
/*
Item buy - change lot and money/fame
2008.12 By Lunar
 input:
	see SQL
 output:
	@nRet	0 = OK
	        1 = Not found item owner or Need more owner's money
	        ? = Error code
	        
 operation:
	1. money/fame check and sub
	2. make buy item
	
*/
@nItemKey bigint,
@nOwner int,
@nOptionType smallint,
@nOptionData bigint,
@nPriceGuildToken bigint,
@nRet int OUTPUT
AS
SET NOCOUNT ON
BEGIN TRAN


EXEC usp_Guild_TokenSub @nOwner, @nPriceGuildToken, @nRet OUTPUT
IF @nRet <> 0
BEGIN
	ROLLBACK TRAN
	SET @nRet = 1
	RETURN
END

-- 아이템 옵션 저장프로시져 호출(랏 수정)
EXEC p_Item_SetOption @nItemKey, @nOptionType, @nOptionData, @nRet OUTPUT
IF @nRet <> 0
BEGIN
	ROLLBACK TRAN
	RETURN
END
-- 완료
COMMIT TRAN
-- end

0㔁
GO

-- ========================================================================
-- Procedure: usp_Item_Set_EnchantGemLot
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Item_Set_EnchantGemLot]
	@nCharNo INT,
	@nEquipmentItemKey BIGINT,
	@nSlot TINYINT,
	@nRestGemCount TINYINT,
	@nRet INT = 0 OUTPUT
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON
	
	SET @nRet = 0
	
	UPDATE tItemEnchant
		SET nRestCount = @nRestGemCount
		FROM tItemEnchant
		WHERE nItemKey = @nEquipmentItemKey
			AND nSlot = @nSlot

	IF @@ROWCOUNT <> 1 OR @@ERROR <> 0
	BEGIN
		SET @nRet = -101
		RETURN
	END 
END

0
GO

-- ========================================================================
-- Procedure: usp_Item_getAutoMonList
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Item_getAutoMonList]
	@nCharNo int
,	@nType	 int
,	@nRet    int = 0 output

AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000
	
	SET @nRet = 0
	
	SELECT nSlot, nItemID
	  FROM tAutoMon
	 WHERE nCharNo = @nCharNo
	   AND nType = @nType
	   
	SET LOCK_TIMEOUT -1
	SET NOCOUNT OFF
END

0訁ឌ
GO

-- ========================================================================
-- Procedure: usp_Item_getEmptyStorage
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Item_getEmptyStorage]
	@owner       int = 0
,	@storageType tinyint = 0
,	@ret         int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	DECLARE @err AS bit
	DECLARE @emptyStorage AS smallint
	SET @err = 0

	IF @owner = 0
		SET @err = 1
	IF @storageType = 0
		SET @err = 1

	IF @err <> 1
		BEGIN	
			IF EXISTS(SELECT nItemKey FROM tItem WHERE nOwner = @owner AND nStorageType = @storageType AND nStorage = 0)
				BEGIN
					SELECT @emptyStorage = ISNULL(MIN(nStorage + 1), 0)
					  FROM tItem 
					 WHERE nOwner = @owner 
					   AND nStorageType = @storageType 
					   AND nStorage + 1 NOT IN (
							SELECT nStorage
							  FROM tItem 
							 WHERE nOwner = @owner AND nStorageType = @storageType)
				END
			ELSE
				BEGIN
					SET @emptyStorage = 0
				END

			SET @ret = @emptyStorage
			PRINT @ret
		END
	ELSE
		SET @ret = -1

	SET NOCOUNT OFF
END

0쌁ர
GO

-- ========================================================================
-- Procedure: usp_Item_setAutoMon
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Item_setAutoMon]
	@nCharNo int
,	@nType	 int
,	@nSlot	 int
,	@nItemID int
,	@nRet    int = 0 output

AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000
	
	SET @nRet = 0
	
	BEGIN TRAN
				
	IF EXISTS (SELECT nCharNo, nType, nSlot FROM tAutoMon WHERE nCharNo = @nCharNo AND nType = @nType AND nSlot = @nSlot) BEGIN
		UPDATE tAutoMon SET nItemID = @nItemID WHERE nCharNo = @nCharNo AND nType = @nType AND nSlot = @nSlot
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = 1
			IF @@TRANCOUNT > 0 ROLLBACK TRAN
			RETURN
		END
	END
	ELSE BEGIN
		INSERT INTO tAutoMon (nCharNo, nType, nSlot, nItemID)
		VALUES (@nCharNo, @ntype, @nSlot, @nItemID)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = 1
			IF @@TRANCOUNT > 0 ROLLBACK TRAN
			RETURN
		END
	END
			
	COMMIT TRAN
	
	SET LOCK_TIMEOUT -1
	SET NOCOUNT OFF

END

0渁
GO

-- ========================================================================
-- Procedure: usp_Item_setPutOnBelonged
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Item_setPutOnBelonged]
/*
Get item storage
2009.9.23 By kodong
'	output Parameter
'		@nRet	 : check SQL execute(valid/invalid)
'			= 0	 : works ok
'			< 0	 : error
'				-2592 : db error(insert)
'				-2593 : db error(update)
*/
@nItemKey     bigint,
@nStorageType tinyint,
@nStorage     smallint,
@nOptionType  smallint,
@nOptionData  bigint,
@nRet         int OUTPUT
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	BEGIN TRAN
	
	UPDATE tItem 
	SET	nStorageType = @nStorageType, nStorage = @nStorage
	WHERE nItemKey = @nItemKey
	IF @@ERROR <> 0 OR @@ROWCOUNT = 0 BEGIN
		SET @nRet = -2593
		ROLLBACK TRAN
		RETURN
	END

	IF EXISTS(SELECT nItemKey FROM tItemOptions WHERE nItemKey = @nItemKey AND nOptionType = @nOptionType) BEGIN
		UPDATE tItemOptions
		SET nOptionType = @nOptionType,
			nOptionData = @nOptionData
		WHERE nItemKey = @nItemKey
		AND   nOptionType = @nOptionType
		IF @@ERROR <> 0 OR @@ROWCOUNT = 0 BEGIN
			SET @nRet = -2593
			ROLLBACK TRAN
			RETURN
		END
	END
	ELSE BEGIN
		INSERT INTO tItemOptions (nItemKey, nOptionType, nOptionData)
		VALUES (@nItemKey, @nOptionType, @nOptionData)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -2592
			ROLLBACK TRAN
			RETURN
		END
	END

	COMMIT TRAN

	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: usp_Item_setPutOnBelongedCancel
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Item_setPutOnBelongedCancel]
	@nCharNo               int
,	@nCancelItemKey        bigint
,	@nCancelOptionType     smallint
,	@nCancelOptionData     bigint
,	@nTargetItemKey        bigint
,	@nTargetOptionType     smallint
,	@nTargetOptionData     bigint
,	@nRet                  int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	DECLARE @cancelOptionData bigint
	DECLARE @targetOptionData bigint
	
	SET @nRet = 0
	SET @cancelOptionData = 0
	SET @targetOptionData = 0
	
	-- invalid nCharNo	
	IF ISNULL(@nCharNo, 0) <= 0 BEGIN
		SET @nRet = -2511
		RETURN
	END
		
	IF NOT EXISTS(SELECT nItemKey FROM tItem WHERE nItemKey = @nCancelItemKey AND nOwner = @nCharNo) BEGIN
		SET @nRet = -2512
		RETURN
	END

	IF NOT EXISTS(SELECT nItemKey FROM tItem WHERE nItemKey = @nTargetItemKey AND nOwner = @nCharNo) BEGIN
		SET @nRet = -2513
		RETURN
	END

	BEGIN TRAN

	-- ÇØÁŠ À¯·áŸÆÀÌÅÛ Â÷°š ¶ÇŽÂ »èÁŠ	
	SELECT @cancelOptionData = nOptionData
	  FROM tItemOptions WITH (UPDLOCK)
	 WHERE nItemKey    = @nCancelItemKey 
	   AND nOptionType = @nCancelOptionType 
	IF @nCancelOptionData <> @cancelOptionData BEGIN
		SET @nRet = -2521
		ROLLBACK TRAN
		RETURN
	END
	IF @cancelOptionData > 1 BEGIN
		SET @cancelOptionData = @cancelOptionData - 1
		EXEC p_Item_SetOption @nCancelItemKey, @nCancelOptionType, @cancelOptionData, @nRet output
		IF @nRet <> 0 BEGIN
			SET @nRet = -2522
			ROLLBACK TRAN
			RETURN						
		END					
	END
	ELSE BEGIN
		EXEC p_Item_Delete @nCancelItemKey, @nRet output
		IF @nRet <> 0 BEGIN
			SET @nRet = -2523
			ROLLBACK TRAN
			RETURN						
		END					
	END
	
	-- Âø¿ëœÃ ±ÍŒÓ ¿ÉŒÇº¯°æ(1 -> 0)
	SELECT @targetOptionData = nOptionData
	  FROM tItemOptions WITH (UPDLOCK)
	 WHERE nItemKey    = @nTargetItemKey 
	   AND nOptionType = @nTargetOptionType 
	IF @targetOptionData <> 1 BEGIN
		SET @nRet = -2524
		ROLLBACK TRAN
		RETURN
	END
	ELSE BEGIN
		EXEC p_Item_SetOption @nTargetItemKey, @nTargetOptionType, @nTargetOptionData, @nRet output
		IF @nRet <> 0 BEGIN
			SET @nRet = -2525
			ROLLBACK TRAN
			RETURN						
		END			
	END
	
	COMMIT TRAN
END

0ᤁ쒊
GO

-- ========================================================================
-- Procedure: usp_Item_setPutOnBelongedExchange
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Item_setPutOnBelongedExchange]
/*
Exchange item storage
2009.9.23 By kodong
 input:
	nItemKey1,nStorageType1,nStoage1, nOptionType1, nOptionData1 = #1 Item
	nItemKey2,nStorageType2,nStoage2 = #2 Item
'	output Parameter
'		@nRet	 : check SQL execute(valid/invalid)
'			= 0	 : works ok
'			< 0	 : error
'				-2592 : db error(insert)
'				-2593 : db error(update)
 operation:
	Set nItemKey1 To nStorageType1,nStoage1
	Set nItemKey2 To nStorageType2,nStoage2
	
	1. #1 item move to temp storage   * #1,#2 is already input that exchanged data.
	2. set #2's storage to #2
	3. set #1's storage to #1
	4. set #1's option #1
*/
@nItemKey1     bigint,
@nStorageType1 tinyint,
@nStorage1     smallint,
@nItemKey2     bigint,
@nStorageType2 tinyint,
@nStorage2     smallint,
@nOptionType1  smallint,
@nOptionData1  bigint,
@nRet          int OUTPUT
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	BEGIN TRAN

	-- Step 1 (nStorageType = 11 은 swap 용도로 사용하는 슬롯)
	UPDATE tItem SET nStorageType = 11
	WHERE nItemKey = @nItemKey1
	IF @@ERROR <> 0 OR @@ROWCOUNT = 0 BEGIN
		SET @nRet = -2593
		ROLLBACK TRAN
		RETURN
	END
		
	-- Step 2
	UPDATE tItem SET nStorageType = @nStorageType2, nStorage = @nStorage2
	WHERE nItemKey = @nItemKey2
	IF @@ERROR <> 0 OR @@ROWCOUNT = 0 BEGIN
		SET @nRet = -2593
		ROLLBACK TRAN
		RETURN
	END
		
	-- Step 3
	UPDATE tItem SET nStorageType = @nStorageType1, nStorage = @nStorage1
	WHERE nItemKey = @nItemKey1
	IF @@ERROR <> 0 OR @@ROWCOUNT = 0 BEGIN
		SET @nRet = -2593
		ROLLBACK TRAN
		RETURN
	END

	-- Step 4 (착용시 귀속 옵션을 찾아 있으면 업데이트, 없으면 새로 생성)
	IF EXISTS(SELECT nItemKey FROM tItemOptions WHERE nItemKey = @nItemKey1 AND nOptionType = @nOptionType1) BEGIN
		UPDATE tItemOptions
		SET nOptionType = @nOptionType1,
			nOptionData = @nOptionData1
		WHERE nItemKey = @nItemKey1
		AND nOptionType = @nOptionType1
		IF @@ERROR <> 0 OR @@ROWCOUNT = 0 BEGIN
			SET @nRet = -2593
			ROLLBACK TRAN
			RETURN
		END
	END
	ELSE BEGIN
		INSERT INTO tItemOptions (nItemKey, nOptionType, nOptionData)
		VALUES (@nItemKey1, @nOptionType1, @nOptionData1)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -2592
			ROLLBACK TRAN
			RETURN
		END
	END
	
	COMMIT TRAN

	SET NOCOUNT OFF
END

0퀁埌Ġ
GO

-- ========================================================================
-- Procedure: usp_Minihouse_delGate
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Minihouse_delGate]
	@nCharNo  int
,	@nItemID  int
,	@nRet     int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	IF EXISTS(SELECT nCharNo FROM tMinihouseGate WHERE nCharNo = @nCharNo AND nItemID = @nItemID) BEGIN
		DELETE FROM tMinihouseGate WHERE nCharNo = @nCharNo AND nItemID = @nItemID
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = -1694
		END
	END
	ELSE BEGIN
		SET @nRet = -1622
	END

	SET NOCOUNT OFF
END

0謁곒
GO

-- ========================================================================
-- Procedure: usp_Minihouse_getGateList
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Minihouse_getGateList]
	@nCharNo  int
,	@nRet     int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	IF @nCharNo <= 0 BEGIN
		SET @nRet = -1611
		RETURN
	END

	IF NOT EXISTS(SELECT nCharNo FROM tMinihouse WHERE nCharNo = @nCharNo) BEGIN
		INSERT INTO tMinihouse (nCharNo) VALUES (@nCharNo)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -1692
			RETURN
		END
	END

	SELECT nCharNo, nItemID, bType, ISNULL(dEndDate, 0) AS dEndDate
	  FROM tMinihouseGate
	 WHERE nCharNo = @nCharNo AND (dEndDate > GETDATE() OR dEndDate IS NULL)

	SET NOCOUNT OFF
END

0쐁ꃶ!
GO

-- ========================================================================
-- Procedure: usp_Minihouse_getVisitorCount
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Minihouse_getVisitorCount]
	@nCharNo int
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	IF EXISTS(SELECT nCharNo FROM tMinihouse WHERE nCharNo = @nCharNo) BEGIN
		SELECT nVisitor FROM tMinihouse WHERE nCharNo = @nCharNo
	END
	ELSE BEGIN
		INSERT INTO tMinihouse (nCharNo) VALUES (@nCharNo)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -1692
			RETURN
		END

		SELECT nVisitor FROM tMinihouse WHERE nCharNo = @nCharNo
	END

	SET NOCOUNT OFF
END

0ﴁ锚"
GO

-- ========================================================================
-- Procedure: usp_Minihouse_setGate
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Minihouse_setGate]
	@nCharNo  int
,	@nItemID  int
,	@bType    tinyint
,	@dEndDate datetime = NULL
,	@nItemKey bigint
,	@nRet     int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0
	
	IF EXISTS(SELECT nCharNo FROM tMinihouseGate WHERE nCharNo = @nCharNo AND nItemID = @nItemID) BEGIN
		-- 기존에 등록된 게이트가 있을 경우 사용기간이 이미 지났으면 삭제처리를 먼저한다.
		DELETE FROM tMinihouseGate 
		 WHERE nCharNo = @nCharNo 
		   AND nItemID = @nItemID 
		   AND dEndDate < GETDATE()
		   AND dEndDate IS NOT NULL
	END

	IF EXISTS(SELECT nCharNo FROM tMinihouseGate WHERE nCharNo = @nCharNo AND nItemID = @nItemID) BEGIN
		-- 상위에서 삭제했는데도 있는 경우는 기간이 남았거나 영구아이템일 경우임.
		SET @nRet = -1621
	END
	ELSE BEGIN
		IF @bType = 0 BEGIN
			SET @dEndDate = NULL
		END

		BEGIN TRAN
		-- 게이트 스크롤 아이템을 삭제함.
		EXEC p_Item_Delete @nItemKey, @nRet output
		IF @nRet = 0 BEGIN
			INSERT INTO tMinihouseGate (nCharNo, nItemID, bType, dEndDate) 
			VALUES (@nCharNo, @nItemID, @bType, @dEndDate)
			IF @@ERROR <> 0 BEGIN
				SET @nRet = -1692
				ROLLBACK TRAN
			END
			ELSE BEGIN
				COMMIT TRAN
			END
		END
		ELSE BEGIN
			SET @nRet = -1622
			ROLLBACK TRAN
		END		
	END

	SET NOCOUNT OFF
END

0㘁褿#
GO

-- ========================================================================
-- Procedure: usp_Minihouse_setVisitorCount
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Minihouse_setVisitorCount]
	@nCharNo int
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	IF EXISTS(SELECT nCharNo FROM tMinihouse WHERE nCharNo = @nCharNo) BEGIN
		UPDATE tMinihouse SET nVisitor = nVisitor + 1 WHERE nCharNo = @nCharNo
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = -1693
		END
	END
	ELSE BEGIN
		INSERT INTO tMinihouse (nCharNo) VALUES (@nCharNo)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -1692
			RETURN
		END

		UPDATE tMinihouse SET nVisitor = nVisitor + 1 WHERE nCharNo = @nCharNo
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = -1693
		END
	END

	SET NOCOUNT OFF
END

0漁絣$
GO

-- ========================================================================
-- Procedure: usp_Pet_getInfo
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Pet_getInfo]
@nPetNo int,
@nRet int OUTPUT
AS

SELECT
nOwnerNo,
sPetName,
nPetMind,
nPetStress
FROM tPetNew
WHERE nPetNo = @nPetNo

SET @nRet = @@ROWCOUNT

IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
BEGIN
	SET @nRet = -1
	RETURN
END
GO

-- ========================================================================
-- Procedure: usp_Pet_setCreate
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Pet_setCreate]
@nOwerNo int,
@sPetName nvarchar(16),
@nPetMind tinyint,
@nPetStress tinyint,
@nRet int OUTPUT
AS

INSERT INTO tPetNew (nOwnerNo, sPetName, nPetMind, nPetStress)
VALUES (@nOwerNo, @sPetName, @nPetMind, @nPetStress)

SET @nRet = @@ROWCOUNT

IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
BEGIN
	SET @nRet = -1
	RETURN
END

SELECT @@IDENTITY
GO

-- ========================================================================
-- Procedure: usp_Pet_setName
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Pet_setName]
@nPetNo int,
@sPetName nvarchar(16),
@nRet int OUTPUT
AS

UPDATE tPetNew
SET sPetName = @sPetName
WHERE nPetNo = @nPetNo

SET @nRet = @@ROWCOUNT

IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
BEGIN
	SET @nRet = -1
	RETURN
END
GO

-- ========================================================================
-- Procedure: usp_SafeTrade_setAll
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_SafeTrade_setAll]
	@nBuyer       int,
	@nSeller      int,
	@nMoveMoney   bigint,
	@nItemKey     bigint,
	@nStorageType tinyint,
	@nStorage     smallint,	
	@nRet         int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	DECLARE @ret AS int
	DECLARE @nOwner AS int
	DECLARE @nOwnerStorageType AS tinyint
	SET @ret = 0
	SET @nOwner = 0
	SET @nOwnerStorageType = 0

	-- 아이템 소유권 확인(판매자 소유의 아이템인지 확인)
	SELECT @nOwner = nOwner, @nOwnerStorageType = nStorageType FROM tItem WHERE nItemKey = @nItemKey
	IF @nOwner <> @nSeller BEGIN
		SET @nRet = -2621
		RETURN
	END
	IF @nOwnerStorageType <> 9 BEGIN
		SET @nRet = -2622
		RETURN
	END
	ELSE BEGIN
		BEGIN TRAN
		
		-- 구매자의 돈을 판매자에게 이동 시킴
		EXEC p_Char_MoneyMove @nBuyer, @nSeller, @nMoveMoney, @ret output
		IF @ret = 0 BEGIN
			ROLLBACK TRAN
			SET @nRet = -2623
			RETURN
		END	
		
		-- 아이템을 구매자의 인벤토리로 이동시킴
		EXEC p_Item_SetStorageAndOwner @nItemKey, @nBuyer, @nStorageType, @nStorage, @ret output
		IF @ret = 0 BEGIN
			ROLLBACK TRAN
			SET @nRet = -2624
			RETURN
		END
		ELSE BEGIN
			COMMIT TRAN
			SET @nRet = 0
		END
	END

	SET NOCOUNT OFF
END
GO

-- ========================================================================
-- Procedure: usp_SafeTrade_setLot
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_SafeTrade_setLot]
	@nBuyer            int,
	@nBuyerItemKey     bigint,
	@nBuyerStorageType tinyint,
	@nBuyerStorage     smallint,
	@nBuyerItemID      int,
	@nBuyerFlags       int,
	@nBuyerOptionType  smallint,
	@nSeller           int,
	@nSellerItemKey    bigint,
	@nSellerOptionType smallint,
	@nMoveMoney        bigint,
	@nSplitLot         bigint,
	@nRet              int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	DECLARE @ret AS int
	DECLARE @nBuyerOptionData AS bigint
	DECLARE @nSellerOptionData AS bigint
	DECLARE @bIsDelete AS bit
	SET @ret = 0
	SET @nBuyerOptionData = 0
	SET @nSellerOptionData = 0
	SET @bIsDelete = 0

	BEGIN TRAN
	
	-- 구매자에게 판매자의 아이템을 구입할 수 있도록 아이템을 생성시킨다. 단, Lot 은 0개로 초기화
	INSERT INTO tItem (nItemKey,  nOwner,  nStorageType,  nStorage,  nItemID,  nFlags)
	VALUES (@nBuyerItemKey, @nBuyer, @nBuyerStorageType, @nBuyerStorage, @nBuyerItemID, @nBuyerFlags)
	IF @@ERROR <> 0 BEGIN
		ROLLBACK TRAN
		SET @nRet = -2628
		RETURN
	END
	ELSE BEGIN
		INSERT INTO tItemOptions (nItemKey, nOptionType, nOptionData)
		VALUES (@nBuyerItemKey, @nBuyerOptionType, @nBuyerOptionData)
		IF @@ERROR <> 0 BEGIN
			ROLLBACK TRAN
			SET @nRet = -2629
			RETURN
		END	
	END
	
	-- 아이템을 구매자의 인벤토리로 이동시킴
	EXEC usp_SafeTrade_setMerge @nBuyer, @nBuyerItemKey, @nBuyerOptionType, @nSeller, @nSellerItemKey,
								@nSellerOptionType, @nMoveMoney, @nSplitLot, @bIsDelete, @ret output
	IF @ret <> 0 BEGIN
		ROLLBACK TRAN
		SET @nRet = @ret
		RETURN
	END
	ELSE BEGIN
		COMMIT TRAN
		SET @nRet = 0
	END

	SET NOCOUNT OFF
END
GO

-- ========================================================================
-- Procedure: usp_SafeTrade_setMerge
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_SafeTrade_setMerge]
	@nBuyer            int,
	@nBuyerItemKey     bigint,
	@nBuyerOptionType  smallint,
	@nSeller           int,
	@nSellerItemKey    bigint,
	@nSellerOptionType smallint,
	@nMoveMoney        bigint,
	@nSplitLot         bigint,
	@bIsDelete         bit,
	@nRet              int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	DECLARE @ret AS int
	DECLARE @nBuyerOptionData AS bigint
	DECLARE @nSellerOptionData AS bigint
	DECLARE @nSellerLot AS bit
	SET @ret = 0
	SET @nBuyerOptionData = 0
	SET @nSellerOptionData = 0
	SET @nSellerLot = 0

	BEGIN TRAN
	
	-- 구매자의 돈을 판매자에게 이동 시킴
	EXEC p_Char_MoneyMove @nBuyer, @nSeller, @nMoveMoney, @ret output
	IF @ret = 0 BEGIN
		ROLLBACK TRAN
		SET @nRet = -2623
		RETURN
	END	
	
	-- 아이템을 구매자의 인벤토리로 이동시킴
	EXEC p_Item_Split_N_Merge_New @nSellerItemKey, @nSeller, @nSellerOptionType, @nSellerOptionData output,
								  @nBuyerItemKey, @nBuyer, @nBuyerOptionType, @nBuyerOptionData output, @nSplitLot, @ret output
	IF @ret <> 0 BEGIN
		ROLLBACK TRAN
		SET @nRet = -2624
		RETURN
	END

	-- 판매자 아이템의 Lot 개수가 0 이 아닌경우 1 로 만들어 판매자의 아이템이 삭제됨을 방지한다.
	IF @nSellerOptionData > 0 BEGIN
		SET @nSellerLot = 1
	END

	IF @bIsDelete = @nSellerLot BEGIN
		IF @bIsDelete = 0 BEGIN
			-- Lot이 0인데 삭제 요청이 오지 않은 경우
			ROLLBACK TRAN
			SET @nRet = -2625
		END
		ELSE BEGIN
			-- Lot이 0이 아닌데 삭제 요청이 온 경우
			ROLLBACK TRAN
			SET @nRet = -2626
		END
	END
	ELSE BEGIN
		IF @bIsDelete = 1 BEGIN
			-- 판매자의 아이템삭제
			EXEC p_Item_Delete @nSellerItemKey, @ret output
			IF @ret <> 0 BEGIN
				ROLLBACK TRAN
				SET @nRet = -2627
			END
		END

		COMMIT TRAN
		SET @nRet = 0
	END

	SET NOCOUNT OFF
END

0ᨁ姐'
GO

-- ========================================================================
-- Procedure: usp_SpamBlock_getTime
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_SpamBlock_getTime]
	@nCharNo int
,	@nRet int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	SET @nRet = 0
	
	IF EXISTS(SELECT nCharNo FROM tSpamBlockTime WHERE nCharNo = @nCharNo)
	BEGIN
		SELECT nCharNo, nBlockTime, nStandbyTime
		  FROM tSpamBlockTime
		 WHERE nCharNo = @nCharNo
	END
	ELSE BEGIN
		SET @nRet = -2822
	END

	SET NOCOUNT ON
END

0匁䷴(
GO

-- ========================================================================
-- Procedure: usp_SpamBlock_setTime
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_SpamBlock_setTime]
	@nCharNo int
,	@nBlockTime int
,	@nStandbyTime int
,	@nRet int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	SET @nRet = 0
	
	IF EXISTS(SELECT nCharNo FROM tSpamBlockTime WHERE nCharNo = @nCharNo)
	BEGIN
		UPDATE tSpamBlockTime
		   SET sys_update_date = GETDATE()
		     , nBlockTime = @nBlockTime
		     , nStandbyTime = @nStandbyTime
		 WHERE nCharNo = @nCharNo
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
		BEGIN
			SET @nRet = -2893
		END
	END
	ELSE BEGIN
		INSERT INTO tSpamBlockTime (nCharNo, nBlockTime, nStandbyTime)
		VALUES (@nCharNo, @nBlockTime, @nStandbyTime)
		IF @@ERROR <> 0
		BEGIN
			SET @nRet = -2892
		END
	END

	SET NOCOUNT ON
END

0谁䈘)
GO

-- ========================================================================
-- Procedure: usp_SpamFilter_get
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_SpamFilter_get]
	@nUserNo int
,	@nRet int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	SET @nRet = 0
	
	IF EXISTS(SELECT nUserNo FROM tSpamFilter WHERE nUserNo = @nUserNo)
	BEGIN
		SELECT nUserNo, bIsSpamSpread, bIsSpamBlock, nSpamFilterRate
		  FROM tSpamFilter
		 WHERE nUserNo = @nUserNo
	END
	ELSE BEGIN
		SET @nRet = -2821
	END

	SET NOCOUNT ON
END

0씁㘼*
GO

-- ========================================================================
-- Procedure: usp_SpamFilter_set
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_SpamFilter_set]
	@nUserNo int
,	@bIsSpamSpread bit
,	@bIsSpamBlock  bit
,	@nSpamFilterRate tinyint
,	@nRet int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	SET @nRet = 0
	
	IF EXISTS(SELECT nUserNo FROM tSpamFilter WHERE nUserNo = @nUserNo)
	BEGIN
		UPDATE tSpamFilter
		   SET sys_update_date = GETDATE()
		     , bIsSpamSpread = @bIsSpamSpread
		     , bIsSpamBlock = @bIsSpamBlock
		     , nSpamFilterRate = @nSpamFilterRate
		 WHERE nUserNo = @nUserNo
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
		BEGIN
			SET @nRet = -2893
		END
	END
	ELSE BEGIN
		INSERT INTO tSpamFilter (nUserNo, bIsSpamSpread, bIsSpamBlock, nSpamFilterRate)
		VALUES (@nUserNo, @bIsSpamSpread, @bIsSpamBlock, @nSpamFilterRate)
		IF @@ERROR <> 0
		BEGIN
			SET @nRet = -2892
		END
	END

	SET NOCOUNT ON
END

0︁⩠+
GO

-- ========================================================================
-- Procedure: usp_Statistics_getCharacterLevel
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Statistics_getCharacterLevel]
	@dDate nvarchar(10) = ''
AS
BEGIN
	SET NOCOUNT ON

	IF LEN(@dDate) = 0 BEGIN
		SET @dDate = CONVERT(NVARCHAR(10), GETDATE(), 120)
	END

	SELECT nLevel AS 'Level',
		   nTotalBefore + nTotalAfter AS 'Total',
		   nBefore + nAfter AS 'Exp>0',
		   nTotalBefore AS 'Total(Before 2009-04-22)',
		   nBefore AS 'Exp>0(Before 2009-04-22)',
		   nTotalAfter AS 'Total(After 2009-04-22)',
		   nAfter AS 'Exp>0(After 2009-04-22)'
	  FROM tStatCharLevel WITH (NOLOCK)
	 WHERE dDate = @dDate

	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: usp_Statistics_setCharacterLevel
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Statistics_setCharacterLevel]
AS
BEGIN
	SET NOCOUNT ON

	INSERT INTO tStatCharLevel
	SELECT CONVERT(NVARCHAR(10), GETDATE(), 120) AS dDate, 
		   A.nLevel, 
		   ISNULL(A.nCnt, 0) AS 'Total(Before 2009-04-22)', 
		   ISNULL(B.nCnt, 0) AS 'Exp>0(Before 2009-04-22)', 
		   ISNULL(C.nCnt, 0) AS 'Total(After 2009-04-22)', 
		   ISNULL(D.nCnt, 0) AS 'Exp>0(After 2009-04-22)'
	  FROM (
			SELECT nLevel, COUNT(nCharNo) AS nCnt
			  FROM tCharacter WITH(NOLOCK)
			 WHERE bDeleted = 0 
			   AND nAdminLevel = 0 
			   AND dCreateDate < '2009-04-22'
			 GROUP BY nLevel) A LEFT OUTER JOIN
		   (
			SELECT nLevel, COUNT(nCharNo) AS nCnt
			  FROM tCharacter WITH(NOLOCK)
			 WHERE bDeleted = 0 
			   AND nAdminLevel = 0 
			   AND nExp > 0
			   AND dCreateDate < '2009-04-22'
			 GROUP BY nLevel) B ON A.nLevel = B.nLevel LEFT OUTER JOIN
		   (
			SELECT nLevel, COUNT(nCharNo) AS nCnt
			  FROM tCharacter WITH(NOLOCK)
			 WHERE bDeleted = 0 
			   AND nAdminLevel = 0 
			   AND dCreateDate >= '2009-04-22'
			 GROUP BY nLevel) C ON A.nLevel = C.nLevel LEFT OUTER JOIN
		   (
			SELECT nLevel, COUNT(nCharNo) AS nCnt
			  FROM tCharacter WITH(NOLOCK)
			 WHERE bDeleted = 0 
			   AND nAdminLevel = 0 
			   AND nExp > 0
			   AND dCreateDate >= '2009-04-22'
			 GROUP BY nLevel) D ON A.nLevel = D.nLevel
	 ORDER BY A.nLevel

	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: usp_Statistics_setGameMoney
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Statistics_setGameMoney]
AS
BEGIN
	SET NOCOUNT ON

	INSERT INTO tStatGameMoney
	SELECT T1.Date, T1.Lv030, T1.Lv060, T1.Lv110, T1.Lv120, T2.UserMoney, T3.GuildMoney
	  FROM (
			SELECT CONVERT(NVARCHAR(10), GETDATE(), 120) AS [Date],
				SUM(CASE WHEN nLevel >= 1 AND nLevel <= 30 THEN nMoney ELSE 0 END) AS [Lv030], 
				SUM(CASE WHEN nLevel >= 31 AND nLevel <= 60 THEN nMoney ELSE 0 END) AS [Lv060], 
				SUM(CASE WHEN nLevel >= 61 AND nLevel <= 110 THEN nMoney ELSE 0 END) AS [Lv110], 
				SUM(CASE WHEN nLevel >= 111 AND nLevel <= 999 THEN nMoney ELSE 0 END) AS [Lv120]
			  FROM tCharacter WITH(NOLOCK)
			 WHERE nAdminLevel = 0 AND bDeleted = 0) AS T1, (
			SELECT CONVERT(NVARCHAR(10), GETDATE(), 120) AS [Date], SUM(nUserMoney) AS [UserMoney]
			  FROM tUserMoney WITH(NOLOCK)) AS T2, (
			SELECT CONVERT(NVARCHAR(10), GETDATE(), 120) AS [Date], SUM(nMoney) AS [GuildMoney]
			  FROM tGuild WITH(NOLOCK)) AS T3
	 WHERE T1.Date = T2.Date
	   AND T1.Date = T3.Date

	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: usp_Statistics_setTotalMoney
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Statistics_setTotalMoney]
AS
BEGIN
	SET NOCOUNT ON
	
	-- 캐릭터가 보유한 총 금액
	IF OBJECT_ID('tCharacter') IS NOT NULL
		AND OBJECTPROPERTY(OBJECT_ID('tCharacter'), 'IsTable') = 1 BEGIN
		INSERT INTO tStatisticsGameMoney (sCode, nMoney)
		SELECT 'Character' AS sCode, ISNULL(SUM(nMoney), 0) AS nMoney
		  FROM tCharacter WITH(NOLOCK)
		 WHERE bDeleted = 0
	END

	-- 계정에 보유한 총 금액
	IF OBJECT_ID('tUserMoney') IS NOT NULL
		AND OBJECTPROPERTY(OBJECT_ID('tUserMoney'), 'IsTable') = 1 BEGIN
		INSERT INTO tStatisticsGameMoney (sCode, nMoney)
		SELECT 'User' AS sCode, ISNULL(SUM(nUserMoney), 0) AS nMoney
		  FROM tUserMoney WITH(NOLOCK)
	END

	-- 길드창고에 보유한 총 금액
	IF OBJECT_ID('tGuildAcademy') IS NOT NULL
		AND OBJECTPROPERTY(OBJECT_ID('tGuildAcademy'), 'IsTable') = 1 BEGIN
		INSERT INTO tStatisticsGameMoney (sCode, nMoney)
		SELECT 'GuildStorage'AS sCode, ISNULL(SUM(nMoney), 0) AS nMoney
		  FROM tGuildAcademy WITH(NOLOCK)
	END

	-- 길드토너먼트 보상 금액중 찾지 않음 금액
	IF OBJECT_ID('tGuild') IS NOT NULL
		AND OBJECTPROPERTY(OBJECT_ID('tGuild'), 'IsTable') = 1 BEGIN
		INSERT INTO tStatisticsGameMoney (sCode, nMoney)
		SELECT 'GuildTournament' AS sCode, ISNULL(SUM(nMoney), 0) AS nMoney
		  FROM tGuild WITH(NOLOCK)
	END

	-- 홀리프라미스 보상 금액중 찾지 않은 금액
	IF OBJECT_ID('tHolyPromiseReward') IS NOT NULL
		AND OBJECTPROPERTY(OBJECT_ID('tHolyPromiseReward'), 'IsTable') = 1 BEGIN
		INSERT INTO tStatisticsGameMoney (sCode, nMoney)
		SELECT 'HolyPromise' AS sCode, ISNULL(SUM(T1.nRemainMoney), 0) AS nMoney
		  FROM tHolyPromiseReward T1 WITH(NOLOCK), (
				SELECT nCharNo, MAX(nSeqNo) AS nMaxSeqNo
				  FROM tHolyPromiseReward WITH(NOLOCK)
				 GROUP BY nCharNo) T2
		 WHERE T1.nCharNo = T2.nCharNo
		   AND T1.nSeqNo = T2.nMaxSeqNo 
		 ORDER BY 1
	END

	-- 경매장 판매 후 찾지 않은 금액
	IF OBJECT_ID('tAuctionInventory') IS NOT NULL
		AND OBJECTPROPERTY(OBJECT_ID('tAuctionInventory'), 'IsTable') = 1 BEGIN
		INSERT INTO tStatisticsGameMoney (sCode, nMoney)
		SELECT 'AuctionInventory' AS sCode, ISNULL(SUM(nMoney), 0) AS nMoney
		  FROM tAuctionInventory WITH(NOLOCK)
		 WHERE bIsDelete = 0
	END
END

0﫱.
GO

-- ========================================================================
-- Procedure: usp_Tutorial_getCharInfo
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Tutorial_getCharInfo]
	@nCharNo int
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	SET @nRet = 0
	
	SELECT nState, nStep FROM tTutorial WHERE nCharNo = @nCharNo 
END

0ᬁ/
GO

-- ========================================================================
-- Procedure: usp_Tutorial_getIsSkip
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Tutorial_getIsSkip]
	@nUserNo int
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	SET @nRet = 0
	
	SELECT nState FROM tTutorial WHERE nUserNo = @nUserNo AND nState = 1
END

0吁0
GO

-- ========================================================================
-- Procedure: usp_Tutorial_setCharInfo
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Tutorial_setCharInfo]
	@nUserNo int
,	@nCharNo int
,	@nState  tinyint
,	@nStep   tinyint
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	SET @nRet = 0

    IF EXISTS (SELECT nCharNo FROM tTutorial WHERE nUserNo = @nUserNo AND nCharNo = @nCharNo)
    BEGIN
        UPDATE tTutorial SET nState = @nState, nStep = @nStep WHERE nUserNo = @nUserNo AND nCharNo = @nCharNo;
        IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 SET @nRet = -1;
    END
    ELSE
    BEGIN
		/*-------------------------------------------------------------------------------------
		Project: FiestaHeroes - NA2016
		Author: xkl
	
		Skip the character tutorial.
		------------------------------------------------------------------------------------*/
		IF EXISTS (SELECT Enabled FROM Options..Tutorial_Skip WHERE Enabled = 1)
            INSERT INTO tTutorial (nUserNo, nCharNo, nState, nStep) VALUES (@nUserNo, @nCharNo, 2, 0);
        ELSE
            INSERT INTO tTutorial (nUserNo, nCharNo, nState, nStep) VALUES (@nUserNo, @nCharNo, 0, 0);
		/*-------------------------------------------------------------------------------------
		Continue as normal ~ 
		------------------------------------------------------------------------------------*/
		
        IF @@ERROR <> 0 SET @nRet = -1;
    END
END
GO

