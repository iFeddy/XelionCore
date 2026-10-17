CREATE PROCEDURE [dbo].[usp_ReturnUser_GetDatediffLastlog]
	@nUserNo int
,	@nRet    int = 0 output
AS
BEGIN 	
	
	SET NOCOUNT ON

	SET @nRet = 0
	
	SELECT TOP 1 DATEDIFF(s, dDate, getdate()) AS diffDate
	  FROM tReturnUser
	 WHERE nUserNo = @nUserNo
	 ORDER BY dDate DESC
	 
	SET NOCOUNT OFF
    
END

0送亾ţ
GO
