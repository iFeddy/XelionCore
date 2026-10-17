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
