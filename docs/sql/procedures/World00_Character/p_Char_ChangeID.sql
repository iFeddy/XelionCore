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
