CREATE PROCEDURE p_IPBlock_GetList
/*
 IP Block Get
 input:
  블럭이 되어있는 목록
 output: 
  Selected records
*/
AS
BEGIN
	SET NOCOUNT ON
	SELECT sIPStart, sIPEnd, sDesc, dDate FROM tIPBlock
	SET NOCOUNT OFF
END
-- end
0弁2
GO
