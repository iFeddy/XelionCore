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
