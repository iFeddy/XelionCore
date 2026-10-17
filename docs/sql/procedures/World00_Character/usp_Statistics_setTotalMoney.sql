CREATE PROCEDURE [dbo].[usp_Statistics_setTotalMoney]
AS
BEGIN
	SET NOCOUNT ON
	
	-- 캐릭터가 보유한 총 금액
	IF OBJECT_ID('tCharacter') IS NOT NULL
		AND OBJECTPROPERTY(OBJECT_ID('tCharacter'), 'IsTable') = 1 BEGIN
		INSERT INTO tStatisticsGameMoney (sCode, nMoney)
		SELECT 'Character' AS sCode, ISNULL(SUM(nMoney), 0) AS nMoney
		  FROM tCharacter WITH(NOLOCK)
		 WHERE bDeleted = 0
	END

	-- 계정에 보유한 총 금액
	IF OBJECT_ID('tUserMoney') IS NOT NULL
		AND OBJECTPROPERTY(OBJECT_ID('tUserMoney'), 'IsTable') = 1 BEGIN
		INSERT INTO tStatisticsGameMoney (sCode, nMoney)
		SELECT 'User' AS sCode, ISNULL(SUM(nUserMoney), 0) AS nMoney
		  FROM tUserMoney WITH(NOLOCK)
	END

	-- 길드창고에 보유한 총 금액
	IF OBJECT_ID('tGuildAcademy') IS NOT NULL
		AND OBJECTPROPERTY(OBJECT_ID('tGuildAcademy'), 'IsTable') = 1 BEGIN
		INSERT INTO tStatisticsGameMoney (sCode, nMoney)
		SELECT 'GuildStorage'AS sCode, ISNULL(SUM(nMoney), 0) AS nMoney
		  FROM tGuildAcademy WITH(NOLOCK)
	END

	-- 길드토너먼트 보상 금액중 찾지 않음 금액
	IF OBJECT_ID('tGuild') IS NOT NULL
		AND OBJECTPROPERTY(OBJECT_ID('tGuild'), 'IsTable') = 1 BEGIN
		INSERT INTO tStatisticsGameMoney (sCode, nMoney)
		SELECT 'GuildTournament' AS sCode, ISNULL(SUM(nMoney), 0) AS nMoney
		  FROM tGuild WITH(NOLOCK)
	END

	-- 홀리프라미스 보상 금액중 찾지 않은 금액
	IF OBJECT_ID('tHolyPromiseReward') IS NOT NULL
		AND OBJECTPROPERTY(OBJECT_ID('tHolyPromiseReward'), 'IsTable') = 1 BEGIN
		INSERT INTO tStatisticsGameMoney (sCode, nMoney)
		SELECT 'HolyPromise' AS sCode, ISNULL(SUM(T1.nRemainMoney), 0) AS nMoney
		  FROM tHolyPromiseReward T1 WITH(NOLOCK), (
				SELECT nCharNo, MAX(nSeqNo) AS nMaxSeqNo
				  FROM tHolyPromiseReward WITH(NOLOCK)
				 GROUP BY nCharNo) T2
		 WHERE T1.nCharNo = T2.nCharNo
		   AND T1.nSeqNo = T2.nMaxSeqNo 
		 ORDER BY 1
	END

	-- 경매장 판매 후 찾지 않은 금액
	IF OBJECT_ID('tAuctionInventory') IS NOT NULL
		AND OBJECTPROPERTY(OBJECT_ID('tAuctionInventory'), 'IsTable') = 1 BEGIN
		INSERT INTO tStatisticsGameMoney (sCode, nMoney)
		SELECT 'AuctionInventory' AS sCode, ISNULL(SUM(nMoney), 0) AS nMoney
		  FROM tAuctionInventory WITH(NOLOCK)
		 WHERE bIsDelete = 0
	END
END

0﫱.
GO
