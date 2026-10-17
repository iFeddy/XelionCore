CREATE PROCEDURE [dbo].[usp_OPT_setUserAuthID]
	@userNo   int
,	@AuthID   tinyint
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	DECLARE @RET AS INT
	
	IF EXISTS(SELECT nUserNo FROM tUser WITH(NOLOCK) WHERE nUserNo = @userNo AND bIsDelete = 0)
		BEGIN
			UPDATE tUser 
			   SET nAuthID = @AuthID
			 WHERE nUserNo = @userNo
			IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
				RETURN (-9)
			ELSE
				RETURN (0)
		END
	ELSE
		RETURN (-1)

	SET LOCK_TIMEOUT -1
	SET NOCOUNT OFF
END

0
GO
