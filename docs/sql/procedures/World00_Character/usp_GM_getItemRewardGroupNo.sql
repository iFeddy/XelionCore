CREATE PROCEDURE [dbo].[usp_GM_getItemRewardGroupNo]
	@sOperID nvarchar(20)
,	@sIP     nvarchar(30)
,	@sCause  nvarchar(1000)
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	SET @nRet = 0
		
	INSERT INTO tGMItemReward (sOperID, sIP, sCause)
	VALUES (@sOperID, @sIP, @sCause)
	IF @@ERROR <> 0 OR @@IDENTITY <= 0 BEGIN
		SET @nRet = -1111
		RETURN
	END
	ELSE BEGIN
		SET @nRet = 0
		
		SELECT @@IDENTITY
		RETURN
	END
END

0靂ŭ
GO
