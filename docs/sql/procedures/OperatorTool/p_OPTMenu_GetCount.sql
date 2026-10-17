CREATE PROCEDURE [dbo].[p_OPTMenu_GetCount]
@nCount tinyint OUTPUT
AS
SET NOCOUNT ON
SELECT @nCount = COUNT(*) FROM tMenuAuth

0찁ĉ
GO
