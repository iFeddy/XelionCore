CREATE PROCEDURE [dbo].[usp_User_chkBlock]
	@userNo INT
,	@RET    INT OUTPUT
AS
BEGIN TRAN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000
	
	DECLARE @BLOCK_NO INT
	DECLARE @TYPE     BIT
	DECLARE @END_DATE DATETIME	
	
	SELECT TOP 1 @BLOCK_NO = nBlockNo, @TYPE = nType, @END_DATE = dEndDate 
	  FROM tUserBlock(NOLOCK) 
	 WHERE nUserNo = @userNo AND bStatus = 1
	 ORDER BY nBlockNo
	IF @TYPE = 1 AND @END_DATE < GETDATE()
		BEGIN
			UPDATE tUserBlock SET bStatus = 0 WHERE nBlockNo = @BLOCK_NO
			IF @@ERROR <> 0 OR @@ROWCOUNT = 0
				BEGIN
					ROLLBACK TRAN
					SET @RET = 0
					RETURN
				END
		
			-- Ãß°¡ °èÁ¤ºí·°»óÅÂ À¯¹«Ã¼Å©
			IF EXISTS(SELECT nBlockNo FROM tUserBlock(NOLOCK) WHERE nUserNo = @userNo AND bStatus = 1)
				SET @RET = 0
			ELSE
				BEGIN
					UPDATE tUser SET bIsBlock = 0 WHERE nUserNo = @userNo
					IF @@ERROR <> 0 OR @@ROWCOUNT = 0
						BEGIN
							ROLLBACK TRAN
							SET @RET = 0
							RETURN
						END
			
					SET @RET = 1
				END	
		END
	SET NOCOUNT OFF
	SET LOCK_TIMEOUT -1
	
COMMIT TRAN

0
GO
