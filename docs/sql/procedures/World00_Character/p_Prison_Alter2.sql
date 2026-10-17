CREATE PROCEDURE [dbo].[p_Prison_Alter2]
@sGmUserID nvarchar(16),
@sImprison nvarchar(16),
@nMinute smallint,
@sReason nvarchar(16), 
@sRemark nvarchar(64),
@err smallint OUTPUT


AS
SET NOCOUNT ON

DECLARE	
@nRet		int,
@nPrisonMin 	int,
@nImprisonRegNum int 

SET 		@err = 0

	--케릭터ID가 존재하지 않습니다.
SELECT @nImprisonRegNum=nCharNo, @nPrisonMin=nPrisonMin FROM tCharacter(NOLOCK) WHERE sID=@sImprison 
IF @@ROWCOUNT = 0
BEGIN
	SET @err = 1
	RETURN @@error
END

--수감중이 아닐때...
IF @nPrisonMin = 0
BEGIN

	SET @err = 2
	RETURN @@error
END

BEGIN TRAN
IF @nMinute = 0
	BEGIN
		UPDATE tCharacter SET
 		nPrisonMin=@nMinute, 
		nLoginZoneX=5963,
		nLoginZoneY=9793,
		sLoginZone='Rou' 		
		WHERE nCharNo=@nImprisonRegNum
	END
ELSE
	BEGIN
		UPDATE tCharacter SET
 		nPrisonMin=@nMinute, 
		nLoginZoneX=896,
		nLoginZoneY=1122,
		sLoginZone='EldPri' 		
		WHERE nCharNo=@nImprisonRegNum

	END

IF @@error <> 0 GOTO ERR

INSERT tPrison (nCharID,nMinute,sReason,sGmID,sRemark)
VALUES (@nImprisonRegNum,@nMinute,@sReason,@sGmUserID,@sRemark )
IF @@error <> 0 GOTO ERR

COMMIT TRAN
RETURN @@error

ERR:
ROLLBACK TRAN
return @@error

0쾭
GO
