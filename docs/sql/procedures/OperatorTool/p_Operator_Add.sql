CREATE PROCEDURE [dbo].[p_Operator_Add]
/*
Operator add
input:
  See SQL
output: 
  nOperNo = 추가된 Operator의 번호(에러시 = 0 )
*/
@nOperNo	int OUTPUT,
@sOperID	nvarchar(20),
@sOperPW	nvarchar(20),
@nOperLevel	tinyint
AS
SET NOCOUNT ON
SET @nOperNo = 0
INSERT tOperator (  sOperID,  sOperPW,  nOperLevel )
VALUES           ( @sOperID, @sOperPW, @nOperLevel )
IF @@ERROR = 0
	SET @nOperNo = @@IDENTITY
	
-- end

0∁靕ď
GO
