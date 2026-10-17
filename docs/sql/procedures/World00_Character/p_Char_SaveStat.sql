CREATE PROCEDURE [dbo].[p_Char_SaveStat] 
@nCharNo int,
@nPS int,
@nGS int,
@nHPS smallint,
@nSPS smallint,
@nHP smallint,
@nSP smallint

AS
SET NOCOUNT ON

UPDATE tCharacter 
	SET 
		nPS=@nPS,
		nGS=@nGS,
		nHPS=@nHPS,
		nSPS=@nSPS,
		nHP=@nHP,
		nSP=@nSP 
	WHERE nCharNo=@nCharNo


return @@error

0 鄂Ų
GO
