CREATE PROCEDURE [dbo].[usp_Item_setSerialKey]
	@userNo    int = 0,
	@serialKey bigint = 0,
	@RET       int = 0 OUTPUT
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	IF @userNo = 0 OR @serialKey = 0
		RETURN (-1)

	IF EXISTS(SELECT nUserNo FROM tUser WITH(NOLOCK) WHERE nUserNo = @userNo AND bIsDelete = 0)
		BEGIN
			BEGIN TRAN 

			IF EXISTS(SELECT userNo FROM tSerialKey WHERE userNo = @userNo)
				DELETE FROM tSerialKey WHERE userNo = @userNo

			INSERT INTO tSerialKey (userNo, serialKey)
			VALUES (@userNo, @serialKey)
			IF @@ERROR <> 0
				BEGIN 
					ROLLBACK TRAN
					SET @RET = 0
				END
			ELSE
				BEGIN
					COMMIT TRAN
					SET @RET = 1
				END
		END

	SET LOCK_TIMEOUT -1
	SET NOCOUNT OFF
END

0
GO
