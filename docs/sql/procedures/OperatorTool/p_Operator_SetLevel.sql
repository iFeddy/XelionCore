CREATE PROCEDURE [dbo].[p_Operator_SetLevel]
/*
Operator set level
input:
  See SQL
output: 
  nRet = 1 TRUE
	   = 0 FALSE
*/
@nOperNo	int,
@nOperLevel	tinyint,
@nRet		int OUTPUT
AS
SET NOCOUNT ON
UPDATE tOperator SET nOperLevel = @nOperLevel WHERE nOperNo = @nOperNo
IF (@@ERROR = 0) AND (@@ROWCOUNT = 1)
	SET @nRet = 1
ELSE
	SET @nRet = 0
-- end

0
GO
