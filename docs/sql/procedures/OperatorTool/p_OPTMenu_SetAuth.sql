CREATE PROCEDURE [dbo].[p_OPTMenu_SetAuth]
@sMenuName	NVARCHAR(20),
@nAuth		tinyint,
@nRet		int OUTPUT
AS
SET NOCOUNT ON
UPDATE tMenuAuth SET auth = @nAuth WHERE menu = @sMenuName
IF (@@ERROR = 0) AND (@@ROWCOUNT = 1)
	SET @nRet = 1
ELSE
	SET @nRet = 0

0ԁ튠Ċ
GO
