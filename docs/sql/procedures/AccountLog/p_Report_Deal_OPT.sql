CREATE PROCEDURE [p_Report_Deal_OPT]
	@nRegNumber	int,
	@nResult	tinyint,
	@sRemark	nvarchar(128),
	
	@nRet		int OUTPUT
	
AS
SET NOCOUNT ON
UPDATE tReport SET 
	nResult			= @nResult,
	sRemark			= @sRemark
WHERE nRegNumber = @nRegNumber
IF @@ERROR <> 0 OR @@ROWCOUNT = 0
BEGIN
	SET @nRet = 0
	RETURN
END
SET @nRet = 1

0؁柦ē
GO
