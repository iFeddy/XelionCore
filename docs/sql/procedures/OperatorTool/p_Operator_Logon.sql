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
