CREATE PROCEDURE [dbo].[usp_ReturnUser_SetUserItem]
	@nUserNo int
,	@nRet    int = 0 output
AS
BEGIN

	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000
	
	DECLARE @itemNo int
	DECLARE @RET	int

	SET @itemNo = 30001 -- 이벤트 아이템은 변경될 수 있음
	SET @nRet = 0
	
	BEGIN TRAN

	EXEC usp_Charge_ItemInsert @nUserNo, 0, @itemNo, 1, @RET OUTPUT
	IF @RET = 1 BEGIN
		INSERT INTO tReturnUser (dDate, nUserNo, nItemNo)
		VALUES (GETDATE(), @nUserNo, @itemNo)

		IF @@TRANCOUNT > 0 COMMIT TRAN
	END
	ELSE BEGIN
		IF @@TRANCOUNT > 0 ROLLBACK TRAN
		SET @nRet = -1
	END
				
	SET LOCK_TIMEOUT -1
	SET NOCOUNT OFF
	
END

0줁䋢Ť
GO
