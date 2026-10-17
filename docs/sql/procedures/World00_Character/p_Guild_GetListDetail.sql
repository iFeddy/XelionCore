CREATE PROCEDURE [dbo].[p_Guild_GetListDetail]

/*
Char Get All Guild lIst detail
2007.11 By khris
*/

AS

SET NOCOUNT ON

SELECT

 nNo,
 sName,
 nMaxMembers,
 nType,
 nGrade,
 nFame,
 nExp
 
FROM tGuild

-- end

0
GO
