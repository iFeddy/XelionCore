CREATE  PROCEDURE [dbo].[usp_Guild_SetToken]
/*
Set Guild Token
2008.12 By Lunar
Input:
	See SQL
	
Output:
	nRet = -1,  Error(Not found guild)
	nRet = -2,  Error(Not found guild)
	nRet = 1,  OK
*/
@nNo int,
@nToken bigint,
@nRet int output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	IF EXISTS(SELECT nNo FROM tGuildToken WHERE nNo = @nNo) BEGIN
		UPDATE tGuildToken SET nToken = @nToken WHERE nNo = @nNo
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = -1
			RETURN
		END
	END
	ELSE BEGIN
		INSERT INTO tGuildToken (nNo, nToken)
		VALUES (@nNo, @nToken)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -2
			RETURN
		END
	END			
	
	SET NOCOUNT OFF
END
-- end

0
GO
