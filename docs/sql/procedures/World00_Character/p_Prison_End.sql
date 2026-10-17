CREATE PROCEDURE [dbo].[p_Prison_End] 
@nImprisonRegNum	int

AS
SET NOCOUNT ON

UPDATE tCharacter SET 	nLoginZoneX=5963, 
			nLoginZoneY=9793, 
			nLoginZoneD=90,
			sLoginZone='Rou', 
			nPrisonMin=0 
	WHERE nCharNo=@nImprisonRegNum

RETURN @@error

0樁滌ą
GO
