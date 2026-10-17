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
