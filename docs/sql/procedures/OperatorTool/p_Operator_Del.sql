CREATE PROCEDURE [dbo].[p_Operator_Del]
/*
Operator delete
input:
  See SQL
output: 
  @nOperLevel = 0, Delete failed
                1, Delete successed
*/
@nOperNo	int,
@nOperLevel tinyint OUTPUT
AS
SET NOCOUNT ON
SET @nOperLevel = 0
DELETE FROM tOperator WHERE nOperNo = @nOperNo
IF @@ROWCOUNT > 0
	SET @nOperLevel = 1
	
-- end

0
GO
