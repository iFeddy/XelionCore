CREATE PROCEDURE [dbo].[p_Item_SetWeaponTitle]
/*

Item Set Weapon Title
2005.8 By CJC

 input:
	see SQL

 output:
	@nRet	0 = OK, 
	        1 = Error Set Weapon mob 
	        2 = Error Set Weapon user
	        3 = Error Delete license item
	        ? = Error DB
*/

@nWeaponItemKey bigint,
@nLicenseItemKey bigint,
@nMobID smallint,

@nTitleCharNo int,
@sTitle nvarchar(24),

-- Output var
@nRet int OUTPUT

AS

SET NOCOUNT ON

BEGIN TRAN
-- 아이템에 몹 라이센스 붙이기
EXEC p_Item_SetMobList @nWeaponItemKey, @nMobID, 0, 0, @nRet OUTPUT
IF @nRet <> 0
BEGIN
	ROLLBACK TRAN
	SET @nRet = 1
	RETURN
END

-- 처음 붙은 아이템의 경우 최초 제작자 정보를 세팅한다.
DECLARE @nRetTemp int
SELECT @nRetTemp = COUNT(*) FROM tItemMobList WHERE nItemKey = @nWeaponItemKey
IF @nRetTemp = 1
BEGIN
	EXEC p_Item_SetTitle @nWeaponItemKey, 0, @nTitleCharNo, @sTitle, @nRet OUTPUT
	IF @nRet <> 0
	BEGIN
		ROLLBACK TRAN
		SET @nRet = 2
		RETURN
	END
END

-- 사용된 라이센스 아이템은 지운다.
SET @nRetTemp = 1
EXEC p_Item_Delete @nLicenseItemKey, @nRetTemp OUTPUT
IF @nRetTemp <> 0
BEGIN
	ROLLBACK TRAN
	SET @nRet = 3
	RETURN
END

-- 성공
COMMIT TRAN
SET @nRet = 0

RETURN

-- end

0萁环Ű
GO
