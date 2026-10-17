CREATE PROCEDURE [dbo].[p_Char_SaveLocation] 
@nCharNo		int,
@sLoginZone		nvarchar(16),
@nLoginZoneX	int,
@nLoginZoneY	int,
@nKQHandle		int,
@sKQMap			nvarchar(16),
@nKQX			int,
@nKQY			int

AS
SET NOCOUNT ON

UPDATE tCharacter SET 
	sLoginZone=@sLoginZone,
	nLoginZoneX=@nLoginZoneX,
	nLoginZoneY=@nLoginZoneY,
	
	nKQHandle=@nKQHandle,
	sKQMap=@sKQMap,
	nKQX=@nKQX,
	nKQY=@nKQY,
	dKQDate=GetDate()
	WHERE nCharNo=@nCharNo

return @@error

0鳝ű
GO
