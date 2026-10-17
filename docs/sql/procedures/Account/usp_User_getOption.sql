CREATE PROCEDURE [dbo].[usp_User_getOption]
	@userNo int
AS
BEGIN
	SET NOCOUNT ON

	SELECT T1.sEmail, T1.sJuminNo, T1.sPhoneNo, T1.nZipNo, T2.sZIPCODE, T2.sSIDO, T2.sGUGUN, T2.sDONG, T1.sAddress, T1.bIsMail
	  FROM tUserOption T1, tZipcode T2
	 WHERE T1.nZipNo = T2.nZipNo
	   AND T1.nUserNo = @userNo

	SET NOCOUNT OFF
END

0
GO
