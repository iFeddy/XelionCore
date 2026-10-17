CREATE PROCEDURE [dbo].[p_Guild_GetName]

/*
Get Guild Intro
2006.2 By Blacksky
*/

@nNo int

AS

SET NOCOUNT ON

SELECT sName FROM tGuild WHERE @nNo = nNo

-- end

0訁ģ
GO
