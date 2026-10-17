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
