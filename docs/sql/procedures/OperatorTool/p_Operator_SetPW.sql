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
