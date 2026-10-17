CREATE PROCEDURE [dbo].[usp_User_search]
	@userName nvarchar(10)
,	@juminNo  nvarchar(20)
AS
BEGIN
	SET NOCOUNT ON

	SELECT T1.nUserNo, T1.sUserID, T1.sUserPW, T2.sEmail 
	  FROM tUser T1, tUserOption T2
	 WHERE T1.nUserNo = T2.nUserNo
	   AND T1.sUserName = @userName
	   AND T2.sJuminNo = @juminNo
	   AND T1.bIsDelete = 0

	SET NOCOUNT OFF
END

0
GO
