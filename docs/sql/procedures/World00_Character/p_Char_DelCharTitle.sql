CREATE PROCEDURE [dbo].[p_Char_DelCharTitle]

@nCharNo		int,
@nType			tinyint, 
@nElementNo		tinyint, 
@nRet			tinyint OUTPUT
AS
SET NOCOUNT ON
-- 해당 상태 삭제
DELETE FROM tCharacterTitle WHERE nCharNo = @nCharNo AND nType = @nType AND nElementNo = @nElementNo
SET @nRet = @@ERROR
-- end

0匁ւŀ
GO
