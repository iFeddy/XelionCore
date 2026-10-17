CREATE PROCEDURE [dbo].[p_Guild_GetCount]

/*
Char Get Guild Count
2006.2 By CJC
*/

@nCount int output

AS

SET NOCOUNT ON

SELECT @nCount = COUNT(*) FROM tGuild

-- end

0댁錗ě
GO
