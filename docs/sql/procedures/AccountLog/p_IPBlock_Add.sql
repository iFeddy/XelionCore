CREATE PROCEDURE p_IPBlock_Add 
/*
 IP Block Add
 input:
  추가할 블럭IP 대역
 output: 
  nRet,  0 = OK, ?=Error code
*/
@sIPStart	nvarchar(20),
@sIPEnd	nvarchar(20),
@sDesc	nvarchar(30),
@nRet		int OUTPUT
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = -1
	INSERT tIPBlock( sIPStart, sIPEnd, sDesc )
	VALUES             ( @sIPStart, @sIPEnd, @sDesc )
	SET @nRet = @@ERROR

	SET NOCOUNT OFF
END
-- end
0䈁ᯜĮ
GO
