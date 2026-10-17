CREATE PROCEDURE p_IPBlock_Get
/*
 IP Block Get
 input:
  블럭이 되어있는 정보 얻기
 output: 
  Selected records
*/
@sIP 		nvarchar(20)
AS
BEGIN
	SET NOCOUNT ON
	SELECT sIPStart, sIPEnd, sDesc, dDate FROM tIPBlock WHERE @sIP >= sIPStart AND @sIP <= sIPEnd
	SET NOCOUNT OFF
END
-- end
0됁Фİ
GO
