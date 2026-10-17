CREATE PROCEDURE [dbo].[usp_Char_chkCharNo]
	@nCharNo int
,	@sID     nvarchar(40)
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	IF EXISTS(SELECT nCharNo FROM tCharacter WHERE nCharNo = @nCharNo AND sID = @sID)
		SET @nRet = 0
	ELSE
		SET @nRet = -1

	SELECT nCharNo, sID FROM tCharacter WHERE nCharNo = @nCharNo

	SET NOCOUNT OFF
END

0䔁䡙$
GO
