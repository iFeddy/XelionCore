CREATE PROCEDURE [dbo].[usp_User_insertOption_Old]
	@userNo  int
,	@eMail   nvarchar(70) = NULL
,	@juminNo nvarchar(20) = NULL
,	@phoneNo nvarchar(20) = NULL
,	@zipNo   int = NULL
,	@address nvarchar(90) = NULL
,	@isMail  bit = NULL
,	@channel nvarchar(20) = NULL
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	INSERT INTO tUserOption (nUserNo, sEmail, sJuminNo, sPhoneNo, nZipNo, sAddress, bIsMail, sChannel)
	VALUES (@userNo, @eMail, @juminNo, @phoneNo, @zipNo, @address, @isMail, @channel)
	IF @@ERROR <> 0
		RETURN (-1)
	ELSE
		RETURN (0)
	
	SET LOCK_TIMEOUT -1
	SET NOCOUNT OFF
END

0
GO
