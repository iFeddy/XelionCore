CREATE PROCEDURE [dbo].[p_OPTMenu_GetAuth]
@sMenuName	NVARCHAR(20)
AS
SET NOCOUNT ON
SELECT auth FROM tMenuAuth WHERE menu = @sMenuName

0威ć
GO
