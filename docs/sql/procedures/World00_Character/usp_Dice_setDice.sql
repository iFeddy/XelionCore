CREATE PROCEDURE [dbo].[usp_Dice_setDice] 
	@nCharNo  int
,	@nItemKey bigint
,	@sDice    nvarchar(3)
,	@nRet     int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	DECLARE @nGameNo AS int

	SET @nRet    = 0
	SET @nGameNo = 0

	IF ISNULL(@nCharNo, 0) <= 0 BEGIN
		SET @nRet = -2411
		RETURN
	END

	IF ISNULL(@nItemKey, 0) <= 0 BEGIN
		SET @nRet = -2412
		RETURN
	END

	IF LEN(ISNULL(@sDice, '')) <= 0 BEGIN
		SET @nRet = -2413
		RETURN
	END

	INSERT INTO tDice (nCharNo, nItemKey, sDice)
	VALUES (@nCharNo, @nItemKey, @sDice)
	IF @@ERROR <> 0 OR @@IDENTITY = 0 BEGIN
		SET @nRet = -2492
		RETURN
	END
	ELSE BEGIN
		SET @nGameNo = @@IDENTITY
		SELECT @nGameNo
	END

	SET NOCOUNT OFF
END
GO
