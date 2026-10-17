CREATE PROCEDURE [dbo].[usp_Character_setStatPoint]
	@nCharNo       int
,	@nMaxStatPoint int
,	@nStrength     tinyint
,	@nConstitute   tinyint
,	@nDexterity    tinyint
,	@nIntelligence tinyint
,	@nMentalPower  tinyint
,	@nItemID       int
,	@nOptionType   smallint
,	@nDeleteCount  int
,	@nRet          int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000
	
	DECLARE @nTotalStat   int
	DECLARE @nDeleteCheck int
	DECLARE @nStr int, @nCon int, @nDex int, @nInt int, @nMen int
	DECLARE @nTotalItemCount int, @nItemKey bigint, @nOptionData bigint, @nItemDeleteRemain int, @nRemainCount int
	
	SET @nRet = 0
	SET @nTotalStat = 0
	SET @nDeleteCheck = 0

	-- invalid nCharNo	
	IF @nCharNo <= 0 BEGIN
		SET @nRet = -1011
		RETURN
	END
	
	-- invalid StatPoint
	IF @nStrength < 0 OR @nConstitute < 0 OR @nDexterity < 0 OR @nIntelligence < 0 OR @nMentalPower < 0 BEGIN
		SET @nRet = -1012
		RETURN
	END
	
	-- have many StatPoint
	SET @nTotalStat = @nStrength + @nConstitute + @nDexterity + @nIntelligence + @nMentalPower
	IF @nMaxStatPoint < @nTotalStat BEGIN
		SET @nRet = -1013
		RETURN
	END
	
	IF EXISTS(SELECT nCharNo FROM tCharacter WHERE nCharNo = @nCharNo AND bDeleted = 0) BEGIN
		SELECT @nStr = nStrength
			,  @nCon = nConstitute 
			,  @nDex = nDexterity 
			,  @nInt = nIntelligence 
			,  @nMen = nMentalPower 
		  FROM tCharacter
		 WHERE nCharNo = @nCharNo
		IF @nStr > @nStrength SET @nDeleteCheck = @nDeleteCheck + (@nStr - @nStrength)
		IF @nCon > @nConstitute SET @nDeleteCheck = @nDeleteCheck + (@nCon - @nConstitute)
		IF @nDex > @nDexterity SET @nDeleteCheck = @nDeleteCheck + (@nDex - @nDexterity)
		IF @nInt > @nIntelligence SET @nDeleteCheck = @nDeleteCheck + (@nInt - @nIntelligence)
		IF @nMen > @nMentalPower SET @nDeleteCheck = @nDeleteCheck + (@nMen - @nMentalPower)
		IF @nDeleteCount <> @nDeleteCheck BEGIN
			SET @nRet = -1021
			RETURN
		END
		
		BEGIN TRAN
		
		UPDATE tCharacter
		   SET nStrength = @nStrength 
		   ,   nConstitute = @nConstitute 
		   ,   nDexterity = @nDexterity 
		   ,   nIntelligence = @nIntelligence 
		   ,   nMentalPower = @nMentalPower 
		 WHERE nCharNo = @nCharNo 
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = -1093
			ROLLBACK TRAN
			RETURN
		END		 

		IF @nDeleteCount > 0  BEGIN
			IF @nDeleteCount = 1 BEGIN
				SELECT TOP 1 @nItemKey = T1.nItemKey, @nOptionData = T2.nOptionData
				  FROM tItem T1, tItemOptions T2
				 WHERE T1.nItemKey = T2.nItemKey 
				   AND T1.nOwner = @nCharNo
				   AND T1.nStorageType = 9
				   AND T1.nItemID = @nItemID 
				 ORDER BY nStorage ASC
				IF @nOptionData > 1 BEGIN
					SET @nOptionData = @nOptionData - 1
					EXEC p_Item_SetOption @nItemKey, @nOptionType, @nOptionData, @nRet output
					IF @nRet <> 0 BEGIN
						SET @nRet = -1022
						ROLLBACK TRAN
						RETURN						
					END					
				END
				ELSE BEGIN
					EXEC p_Item_Delete @nItemKey, @n
GO
