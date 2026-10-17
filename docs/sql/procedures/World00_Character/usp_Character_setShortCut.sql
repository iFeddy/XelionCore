CREATE PROCEDURE [dbo].[usp_Character_setShortCut]
	@nCharNo int,
	@nSlotNo tinyint,
	@nCodeNo smallint,
	@nValue  int,
	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	SET @nRet = 0
	
	DECLARE @nCodeType smallint SET @nCodeType = 0
	DECLARE @nDBCharNo int
	DECLARE @nPrevCodeNo smallint SET @nPrevCodeNo = 0
	DECLARE @nPrevValue int SET @nPrevValue = 0

	SELECT @nDBCharNo = nCharNo, @nPrevCodeNo = nCodeNo, @nPrevValue = nValue
	  FROM tCharacterShortCut
	 WHERE nCharNo = @nCharNo 
	   AND nSlotNo = @nSlotNo
	IF @nDBCharNo IS NULL BEGIN
		INSERT INTO tCharacterShortCut (nCharNo, nSlotNo, nCodeType, nCodeNo, nValue)
		VALUES (@nCharNo, @nSlotNo, 1, @nCodeNo, @nValue)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -2791
		END
	END
	ELSE IF @nPrevCodeNo = @nCodeNo AND @nPrevValue = @nValue BEGIN
		SET @nRet = -2722
	END
	ELSE BEGIN
		UPDATE tCharacterShortCut 
		   SET nCodeNo = @nCodeNo
		     , nValue  = @nValue
		     , sys_update_date = GETDATE()
		 WHERE nCharNo = @nCharNo 
		   AND nSlotNo = @nSlotNo
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = -2793
		END
	END

	SET NOCOUNT OFF
END

0@
GO
