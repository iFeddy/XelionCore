CREATE PROCEDURE [dbo].[p_Prison_Add2]
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
	--이미 수감중입니다
IF @nPrisonMin <> 0
BEGIN

	SET @err = 2
	RETURN @@error
END

BEGIN TRAN

UPDATE tCharacter SET
 	nPrisonMin=@nMinute, 
	nLoginZoneX=896,
	nLoginZoneY=1122,
	sLoginZone='EldPri' 		
	WHERE nCharNo=@nImprisonRegNum
IF @@error <> 0 GOTO ERR

INSERT tPrison (nCharID,nMinute,sReason,sGmID,sRemark)
VALUES (@nImprisonRegNum,@nMinute,@sReason,@sGmUserID,@sRemark )
IF @@error <> 0 GOTO ERR

COMMIT TRAN
RETURN @@error

ERR:
ROLLBACK TRAN
return @@error

0稁
GO
