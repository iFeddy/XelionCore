CREATE PROCEDURE [dbo].[usp_User_GetDatediffCreated]
	@nUserNo int
,	@nRet    int = 0 output
AS
BEGIN

	SET NOCOUNT ON

	SET @nRet = 0
	
	SELECT DATEDIFF(s, dDate, getdate()) AS diffDate
	  FROM tUser
	 WHERE nUserNo = @nUserNo
	 
    SET NOCOUNT OFF
    
END

0鼁o
GO
