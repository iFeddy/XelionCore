CREATE PROCEDURE [dbo].[p_Guild_GetNotify]

/*
Get Guild Notify
2006.2 By CJC
*/

@nNo int

AS

SET NOCOUNT ON

SELECT dNotifyDate, sNotifyCharID, sNotify FROM tGuild WHERE @nNo = nNo

-- end

0ᲁĥ
GO
