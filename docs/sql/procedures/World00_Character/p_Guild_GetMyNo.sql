CREATE PROCEDURE [dbo].[p_Guild_GetMyNo]

/*
Get My Guild Number
2006.2 By CJC

Input:
	See SQL
	
Output:
	@nNo = 0, I'm not join guild
	@nNo = ?, My guild No.
*/

@nCharNo int,
@nNo     int output

AS

SET NOCOUNT ON

SET @nNo = 0
SELECT @nNo = nNo FROM tGuildMember WHERE nCharNo = @nCharNo

-- end

0笁㐹ģ
GO
