CREATE PROCEDURE p_IPBlock_Delete
/*
 IP Block Delete
 input:
  블럭이 되어있는 정보 제거
 output: 
  Selected records
*/
@sIP 		nvarchar(20),
@nRet		int OUTPUT
AS
BEGIN
	SET NOCOUNT ON
	DELETE FROM tIPBlock WHERE @sIP >= sIPStart AND @sIP <= sIPEnd
	SET @nRet = @@ERROR
	SET NOCOUNT OFF
END
-- end
0笁ကį
GO
