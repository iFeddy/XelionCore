CREATE PROCEDURE [dbo].[usp_Item_GuildToken_BuyLot]
/*
Item buy - change lot and money/fame
2008.12 By Lunar
 input:
	see SQL
 output:
	@nRet	0 = OK
	        1 = Not found item owner or Need more owner's money
	        ? = Error code
	        
 operation:
	1. money/fame check and sub
	2. make buy item
	
*/
@nItemKey bigint,
@nOwner int,
@nOptionType smallint,
@nOptionData bigint,
@nPriceGuildToken bigint,
@nRet int OUTPUT
AS
SET NOCOUNT ON
BEGIN TRAN


EXEC usp_Guild_TokenSub @nOwner, @nPriceGuildToken, @nRet OUTPUT
IF @nRet <> 0
BEGIN
	ROLLBACK TRAN
	SET @nRet = 1
	RETURN
END

-- 아이템 옵션 저장프로시져 호출(랏 수정)
EXEC p_Item_SetOption @nItemKey, @nOptionType, @nOptionData, @nRet OUTPUT
IF @nRet <> 0
BEGIN
	ROLLBACK TRAN
	RETURN
END
-- 완료
COMMIT TRAN
-- end

0㔁
GO
