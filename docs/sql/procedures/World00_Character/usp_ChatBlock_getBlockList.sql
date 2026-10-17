CREATE PROCEDURE [dbo].[usp_ChatBlock_getBlockList]
	@nUserNo int
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	IF ISNULL(@nUserNo, 0) <= 0 BEGIN
		SET @nRet = -2311
		RETURN
	END

	SELECT T1.nCharNo, T2.sID, T2.nAdminLevel, T2.bDeleted
	  FROM tChatBlock T1, tCharacter T2
	 WHERE T1.nUserNo = @nUserNo AND T1.nCharNo = T2.nCharNo
END

0ȁ罹M
GO
