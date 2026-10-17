CREATE PROCEDURE [dbo].[p_Item_ClrWeaponTitle]
/*

Item Clear Weapon Title
2006.6 By CJC

 input:
	see SQL

 output:
	@nRet	0 = OK, 
	        1 = Error Not found Weapon mob 
	        2 = Error Clear Weapon mob 
	        3 = Error Delete license item
	        4 = Error Change license item lot
*/

@nWeaponItemKey bigint,

@nMobID smallint,
@nLicenseClearItemKey bigint,
@bLicenseClearItemDelete tinyint,
@nLicenseClearItemOptionType smallint,
@nLicenseClearItemOptionData bigint,

-- Output var
@nRet int OUTPUT

AS

SET NOCOUNT ON

-- 아이템에 몹 라이센스 붙어 있는지 확인 
IF NOT EXISTS(SELECT nMobID FROM tItemMobList WHERE nItemKey=@nWeaponItemKey AND nMobID=@nMobID)
BEGIN
	SET @nRet = 1
	RETURN
END

-- 아이템에 몹 라이센스 지우기
BEGIN TRAN
DELETE FROM tItemMobList WHERE nItemKey=@nWeaponItemKey AND nMobID=@nMobID
IF @@ERROR <> 0
BEGIN
	ROLLBACK TRAN
	SET @nRet = 2
	RETURN
END

-- 라이센스 지우는 아이템 지우기/랏조정
DECLARE @nRetTemp int
IF @bLicenseClearItemDelete = 1
BEGIN
	EXEC p_Item_Delete @nLicenseClearItemKey, @nRetTemp OUTPUT
	IF @nRetTemp <> 0
	BEGIN
		ROLLBACK TRAN
		SET @nRet = 3
		RETURN
	END
END
ELSE
BEGIN
	EXEC p_Item_SetOption @nLicenseClearItemKey, @nLicenseClearItemOptionType, @nLicenseClearItemOptionData, @nRetTemp OUTPUT
	IF @nRetTemp <> 0
	BEGIN
		ROLLBACK TRAN
		SET @nRet = 4
		RETURN
	END
END

-- 성공
COMMIT TRAN
SET @nRet = 0

RETURN

-- end

0렚Ŋ
GO
