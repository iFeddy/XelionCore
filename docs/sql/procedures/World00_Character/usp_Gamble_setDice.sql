CREATE PROCEDURE [dbo].[usp_Gamble_setDice] 
	@sDice    nvarchar(3)
,	@nRet     int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	DECLARE @nGameNo AS int

	SET @nRet    = 0
	SET @nGameNo = 0

	IF LEN(ISNULL(@sDice, '')) <= 0 BEGIN
		SET @nRet = -3314
		RETURN
	END

	INSERT INTO tGambleDice (sDice)
	VALUES (@sDice)
	IF @@ERROR <> 0 OR @@IDENTITY = 0 BEGIN
		SET @nRet = -3392
		RETURN
	END
	ELSE BEGIN
		SET @nGameNo = @@IDENTITY
		SELECT @nGameNo
	END

	SET NOCOUNT OFF
END

0
GO
