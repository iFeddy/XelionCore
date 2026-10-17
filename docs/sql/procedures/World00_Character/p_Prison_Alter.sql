CREATE PROCEDURE [dbo].[p_Prison_Alter] 
@sGmUserID nvarchar(16),
@sImprison nvarchar(16),
@nMinute smallint,
@sReason nvarchar(16), 
@sRemark nvarchar(64),
@nImprisonRegNum int,
@err smallint OUTPUT


AS
SET NOCOUNT ON

DECLARE	
@nRet		int,
@nPrisonMin 	int

SET 		@err = 0

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

0蚃ă
GO
