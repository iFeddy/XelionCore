CREATE PROCEDURE [dbo].[usp_Item_GuildToken_BuyAll]
/*
Item buy all
2004.8 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK
	        1 = Not found item owner or Need more owner's money/fame
	        ? = Error code
	        
 operation:
	1. money/fame check and sub
	2. make buy item
	
*/
@nItemKey bigint,
@nOwner int,
@nStorageType tinyint,
@nStorage smallint,
@nItemID int,
@nFlags int,
@nPriceGuildToken bigint,
@nRet int OUTPUT
AS
SET NOCOUNT ON
DECLARE @CreateItemRet int
BEGIN TRAN

-- 아이템 소유자가 충분한 돈을 가지고 있는지 확인 후 돈 감소

EXEC usp_Guild_TokenSub @nOwner, @nPriceGuildToken, @nRet OUTPUT
IF @nRet <> 0
BEGIN
	ROLLBACK TRAN
	SET @nRet = 1
	RETURN
END


-- 아이템 생성 저장프로시져 호출
EXEC p_Item_Create @nItemKey, @nOwner, @nStorageType, @nStorage, @nItemID, @nFlags, @CreateItemRet OUTPUT
IF @CreateItemRet <> 0 OR @@ERROR <> 0
BEGIN
	ROLLBACK TRAN
	SET @nRet = @CreateItemRet
	RETURN
END
-- 완료
COMMIT TRAN
SET @nRet = 0
-- end

0ﰁￔ
GO
