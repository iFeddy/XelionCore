CREATE PROCEDURE [dbo].[usp_Item_getUser]
	@userNo    INT = 0,
	@serialKey INT = 0,
	@USER_ID   NVARCHAR(20) = '' OUTPUT,
	@AUTH_ID   INT = 0 OUTPUT
AS
BEGIN
	SET NOCOUNT ON

	IF @serialKey = 0
		BEGIN
			SET @USER_ID = ''
			SET @AUTH_ID = 0
		END
	ELSE
		BEGIN
			IF EXISTS(SELECT nUserNo FROM tUser WITH (NOLOCK) WHERE nUserNo = @userNo AND bIsDelete = 0)
				BEGIN
					SELECT @USER_ID = sUserID, @AUTH_ID = nAuthID
					  FROM tUser AS T1 INNER JOIN tSerialKey AS T2 ON T1.nUserNo = T2.userNo
					 WHERE T2.userNo = @userNo AND T2.serialKey = @serialKey

					SET @USER_ID = ISNULL(@USER_ID, '')
					IF @USER_ID <> '' OR @USER_ID <> NULL
						DELETE FROM tSerialKey WHERE userNo = @userNo
				END
			ELSE
				BEGIN
					SET @USER_ID = ''
					SET @AUTH_ID = 0
				END
		END

	SET NOCOUNT OFF
END

0
GO
