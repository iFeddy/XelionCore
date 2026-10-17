CREATE PROCEDURE [p_AccountLog_GetLog]
	@QueryExpr nvarchar(256)
AS
BEGIN
	SET NOCOUNT ON
	DECLARE @SQLCMD nvarchar(512)
	SET @SQLCMD = N'SELECT a.nLogNo, a.dDate, a.nUserNo, b.sUserID, a.nType, a.sIP, a.nPlayMin, a.sUserPassword, a.nErrorCode, a.nWorldNo
			  FROM tAccountLog  a INNER JOIN AccountDB.Account.dbo.tUser b ON a.nUserNo = b.nUserNo ' + @QueryExpr
	EXEC sp_executesql @SQLCMD
	SET NOCOUNT OFF
END

0꼁ᧆą
GO
