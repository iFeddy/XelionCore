CREATE PROCEDURE [dbo].[usp_User_chkRegister]
	@juminNo    nvarchar(20) = ''
,	@age        int = 0 OUTPUT
,	@registDate nvarchar(10) = '' OUTPUT
AS
BEGIN
	SET NOCOUNT ON

	DECLARE @USER_NO  AS INT
	DECLARE @BIRTHDAY AS NVARCHAR(10)

	IF @juminNo = ''
		RETURN -1

	SET @BIRTHDAY = CASE SUBSTRING(@juminNo, 8, 1)
						WHEN '1' THEN '19'
						WHEN '2' THEN '19'
						WHEN '3' THEN '20'
						WHEN '4' THEN '20'
						ELSE ''
					END
	IF @BIRTHDAY = ''
		RETURN -1
	ELSE
		SET @BIRTHDAY = @BIRTHDAY + SUBSTRING(@juminNo, 1, 2) + '-' + SUBSTRING(@juminNo, 3, 2) + '-' + SUBSTRING(@juminNo, 5, 2)

	SET @age = DATEDIFF(yy, @BIRTHDAY, GETDATE())
	IF @age < 15
		RETURN -2

	SELECT @USER_NO = T1.nUserNo, @registDate = T1.dDate 
	  FROM tUser T1, tUserOption T2
	 WHERE T1.nUserNo = T2.nUserNo
	   AND T2.sJuminNo = @juminNo
	   AND T1.bIsDelete = 0
	IF @USER_NO	> 100
		BEGIN
			SELECT @registDate = dDate FROM tUser WITH(NOLOCK) WHERE nUserNo = @USER_NO
			RETURN -3
		END
	ELSE
		RETURN 0

	SET NOCOUNT OFF
END
GO
