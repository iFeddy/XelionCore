CREATE PROCEDURE [dbo].[p_Item_CostumWeapon_Endure_Set]
/*
Item Costum Weapon Endure Set
2008.10 By Lunar
 input:
	see SQL
 output:
	@nRet	0 = OK
	        1 = Not found item
	        2 = Missmatch item owner
	        3 = Not enough money to sub
	        4 = Upgrade value setting, upgrade fail count error(option set error)
	        5 = Raw delete fail
	        6 = RawLeft delete fail, if valid item
	        7 = RawRight delete fail, if valid item
	        8 = RawMiddle delete fail, if valid item
	        
 operation:
	1. money sub
	2. upgrade item level
	3. delete raws
*/
@nCostumWeaponKey bigint,
@nCostumWeaponEndureOptionType smallint,
@nCostumWeaponEndureOptionData bigint,
@nRet int OUTPUT
AS
BEGIN
	SET NOCOUNT ON
	
	DECLARE @nTinyintRet tinyint
	DECLARE @nRetTemp int

	BEGIN TRAN

	-- 1 가구 아이템 내구도 상승
	SET @nRetTemp = 1
	EXEC p_Item_SetOption @nCostumWeaponKey, @nCostumWeaponEndureOptionType, @nCostumWeaponEndureOptionData, @nRetTemp OUTPUT
	IF @nRetTemp <> 0 BEGIN
		ROLLBACK TRAN
		SET @nRet = 5
		RETURN
	END
	
	-- 업글 성공
	COMMIT TRAN
	SET @nRet = 0

	SET NOCOUNT OFF
END

0圁ꌌJ
GO
