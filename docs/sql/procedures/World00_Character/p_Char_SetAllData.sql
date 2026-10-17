CREATE PROCEDURE [dbo].[p_Char_SetAllData]
/*
Character Set all data
2004.6 By CJC
 input:
   All args
 output: 
   1 = update ok
   0 = error at update chatacter data
*/
-- Character columns
@nCharNo int,
@nUserNo int, 
@nSlotNo tinyint,
@sID nvarchar(40), 
@nLevel smallint, 
@nExp bigint, 
@nHP smallint, 
@nSP smallint, 
@nAP smallint, 
@nHPS smallint,
@nSPS smallint,
@nPS int,
@nGS int,
@nMoney bigint,
@nFame int, 
@nFlags int, 
@nAdminLevel tinyint, 
@sLoginZone varchar(16),
@nLoginZoneX int,
@nLoginZoneY int,
@nLoginZoneD tinyint,
@nPrisonMin smallint, 
@nPKCount int, 
@nPKYellowTime tinyint,
@nStrength tinyint,
@nConstitute tinyint,
@nDexterity tinyint,
@nIntelligence tinyint,
@nMentalPower tinyint,
@nRedistributePoint tinyint,
@nPlayMin int, 
@nLoginCount int, 
@dLastLoginDate datetime, 
@dCreateDate datetime, 
@nCreateWorld tinyint, 
@bDeleted tinyint, 
@dDeletedDate datetime,
-- Output var
@nRet int OUTPUT
AS
SET NOCOUNT ON
-- Update character
UPDATE tCharacter SET
	nUserNo		=	@nUserNo , 
	nSlotNo			=	@nSlotNo , 
	sID			=	@sID , 
	nLevel			=	@nLevel , 
	nExp			=	@nExp , 
	nHP			=	@nHP , 
	nSP			=	@nSP , 
	nAP			=	@nAP , 
	nHPS			=	@nHPS , 
	nSPS			=	@nSPS , 
	nPS 			=	@nPS,
	nGS			=	@nGS,
	
	nMoney			= 	@nMoney , 
	nFame			= 	@nFame , 
	nFlags			= 	@nFlags , 
	nAdminLevel		=	@nAdminLevel ,
	
	sLoginZone		=	@sLoginZone , 
	nLoginZoneX		=	@nLoginZoneX , 
	nLoginZoneY		=	@nLoginZoneY , 
	nLoginZoneD		=	@nLoginZoneD ,
	
	nPrisonMin		=	@nPrisonMin , 
	nPKCount		=	@nPKCount , 
	nPKYellowTime		=	@nPKYellowTime ,
	
	nStrength		= 	@nStrength,
	nConstitute		= 	@nConstitute,
	nDexterity		= 	@nDexterity,
	nIntelligence		= 	@nIntelligence,
	nMentalPower		= 	@nMentalPower,
	nRedistributePoint	= 	@nRedistributePoint,
	
	nPlayMin		=	@nPlayMin , 
	nLoginCount		=	@nLoginCount , 
	dLastLoginDate		=	@dLastLoginDate ,
	
	dCreateDate		=	@dCreateDate , 
	nCreateWorld		=	@nCreateWorld ,
	bDeleted		=	@bDeleted, 
	dDeletedDate		=	@dDeletedDate 
WHERE nCharNo = @nCharNo
IF @@ERROR <> 0 OR @@ROWCOUNT = 0
BEGIN
	SET @nRet = 0
	RETURN
END
-- Update success
SET @nRet = 1
-- end

0ꄁ㡤ŵ
GO
