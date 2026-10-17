CREATE PROCEDURE [dbo].[p_Prison_UpdateCharPrisonMin] 
@nImprisonRegNum	int, 
@nMinute		smallint

AS
SET NOCOUNT ON

--감옥시간과 로긴존과의 동기화를 위해서 꼭!!! 같이 업데이트 해줘야함.
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
		sLoginZone='EldPri' 		
		WHERE nCharNo=@nImprisonRegNum

	END

RETURN @@error

0ꌁ拰Ć
GO
