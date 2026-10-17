CREATE PROCEDURE [dbo].[usp_Collection_setSerial]
	@nCardID      int
,	@nCount       bigint
,	@nRet         int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0
	
	IF ISNULL(@nCardID, 0) = 0 BEGIN
		SET @nRet = -1111
		RETURN
	END
	
	IF EXISTS(SELECT nCardID FROM tCollectionSerialCount WHERE nCardID = @nCardID) BEGIN
		BEGIN TRAN
		
		UPDATE tCollectionSerialCount 
		   SET sys_update_date = GETDATE()
		     , nCount          = @nCount
		 WHERE nCardID = @nCardID 
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = -4444
			ROLLBACK TRAN
			RETURN
		END
		
		COMMIT TRAN
	END
	ELSE BEGIN
		IF @nCount <> 1 BEGIN
			SET @nRet = -5555
			RETURN
		END
		
		INSERT INTO tCollectionSerialCount (nCardID, nCount)
		VALUES (@nCardID, @nCount)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -6666
			RETURN
		END
	END
END
GO
