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
