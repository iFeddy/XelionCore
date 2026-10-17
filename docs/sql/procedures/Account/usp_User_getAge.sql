CREATE PROCEDURE [dbo].[usp_User_getAge]
	@userNo int
,	@AGE    int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	DECLARE @JUMIN_NO AS NVARCHAR(14)
	DECLARE @BIRTHDAY AS NVARCHAR(10)

	IF EXISTS(SELECT sJuminNo FROM tUserOption WITH(NOLOCK) WHERE nUserNo = @userNo)
		BEGIN
			SELECT @JUMIN_NO = sJuminNo FROM tUserOption WITH(NOLOCK) WHERE nUserNo = @userNo

			SET @BIRTHDAY = CASE SUBSTRING(@JUMIN_NO, 8, 1)
								WHEN '1' THEN '19'
								WHEN '2' THEN '19'
								WHEN '3' THEN '20'
								WHEN '4' THEN '20'
								ELSE ''
							END
			IF @BIRTHDAY = ''
				SET @AGE = 0
			ELSE
				BEGIN
					SET @BIRTHDAY = @BIRTHDAY + SUBSTRING(@JUMIN_NO, 1, 2) + '-' + SUBSTRING(@JUMIN_NO, 3, 2) + '-' + SUBSTRING(@JUMIN_NO, 5, 2)
					SET @AGE = DATEDIFF(yy, @BIRTHDAY, GETDATE())
				END
		END
	ELSE
		SET @AGE = 0

	SET NOCOUNT OFF
END

0
GO
