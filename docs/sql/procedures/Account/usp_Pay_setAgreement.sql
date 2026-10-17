CREATE PROCEDURE [dbo].[usp_Pay_setAgreement]
	@userNo    int
,	@name      nvarchar(10)
,	@juminNo   nvarchar(20)
,	@phoneNo   nvarchar(20)
,	@vaildDate int
AS
BEGIN
	SET NOCOUNT ON

	INSERT INTO tPayAgreement (userNo, name, juminNo, phoneNo, vaildDate)
	VALUES (@userNo, @name, @juminNo, @phoneNo, DATEADD(mm, @vaildDate, GETDATE()))

	SET NOCOUNT OFF
END

0
GO
