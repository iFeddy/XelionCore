CREATE PROCEDURE [dbo].[usp_User_insert]
	@userID   nvarchar(20)
,	@userPW   nvarchar(32)
,	@nAffiliateID varchar(50) = NULL
,	@szAffilateSubID varchar(50) = NULL
,	@szCountryCode varchar(50) = NULL
,	@userName nvarchar(10)
,	@userIP   nvarchar(30)
,	@eMail    nvarchar(70) = NULL
,	@juminNo  nvarchar(20) = NULL
,	@phoneNo  nvarchar(20) = NULL
,	@zipNo    int = NULL
,	@address  nvarchar(90) = NULL
,	@isMail   bit = NULL
,	@channel  nvarchar(20) = 'ONSON'
,	@isOption bit = 1
,	@userNo   int = 0 OUTPUT

AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	DECLARE @RET AS INT
	DECLARE @itemNo AS INT

	SET @itemNo = 30000 -- 이벤트 아이템은 변경될 수 있음

	EXEC @RET = usp_User_chkUserID @userID
	IF (@RET <> 0)
		RETURN (-2)

	BEGIN TRAN

	INSERT INTO tUser (sUserID, sUserPW, sUserName, sUserIP)
	VALUES (@userID, @userPW, @userName, @userIP)
	IF (@@ERROR <> 0 OR @@IDENTITY = 0)
		BEGIN
			ROLLBACK TRAN
			RETURN (-9)
		END
	ELSE
		BEGIN                        
			SET @userNo = @@IDENTITY

			IF (@isOption = 1)
				BEGIN
					EXEC @RET = usp_User_insertOption @userNo,  @eMail, @juminNo, @phoneNo, @zipNo, @address, @isMail, @channel, @nAffiliateID
					IF (@RET <> 0)
						BEGIN
							ROLLBACK TRAN
							RETURN (-3)
						END
				END
	END

	SET @RET = 0
	EXEC usp_Charge_ItemInsert @userNo, 0, @itemNo, 1, @RET OUTPUT

	COMMIT TRAN
	SET LOCK_TIMEOUT -1
	SET NOCOUNT ON
END

0Ё憓Ŷ
GO
