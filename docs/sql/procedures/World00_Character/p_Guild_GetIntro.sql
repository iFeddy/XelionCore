CREATE PROCEDURE [dbo].[p_Guild_GetIntro]

/*
Get Guild Intro
2006.2 By CJC
*/

@nNo int

AS

SET NOCOUNT ON

SELECT sIntro FROM tGuild WHERE @nNo = nNo

-- end

0度賖Ġ
GO
