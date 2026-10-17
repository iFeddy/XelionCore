CREATE PROCEDURE [dbo].[usp_Character_getChangeIDList]
	@nUserNo int
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0
	
	SELECT nRowNo, nUserNo, nCharNo, bInit
	  FROM tCharacterChangeID
	 WHERE nUserNo = @nUserNo AND bFlag = 0
	   
	SET NOCOUNT OFF
END

0딁붣Ĭ
GO
