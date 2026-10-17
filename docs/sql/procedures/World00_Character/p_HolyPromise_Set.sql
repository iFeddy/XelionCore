CREATE PROCEDURE [dbo].[p_HolyPromise_Set]
/*
Holy promise set
2008.4 By CJC

 input:
	see SQL
	
 output:
	@nRet
	0 = OK,  Fiest join
	1 = Ring / same char
	2 = i'm over max member count
	3 = OK, ReJoin
	4 = Failed, Already upper exists
	? = Other DB error code
*/
@nMAXPROMISEMEMBER smallint,
@nCharNo int,
@nUpCharNo int,
@nFlags tinyint,
@nRet int OUTPUT

AS

BEGIN
	SET NOCOUNT ON
	
	-- check ring, same char.
	EXEC p_HolyPromise_CheckRing @nCharNo, @nUpCharNo, @nRet OUTPUT
	IF @nRet = 1
		RETURN
	
	-- check over max promise and add
	BEGIN TRAN

	DECLARE @nPromiseCount int
	EXEC p_HolyPromise_GetMemberCount @nUpCharNo, @nPromiseCount OUTPUT
	IF @nPromiseCount >= @nMAXPROMISEMEMBER
	BEGIN
		SET @nRet = 2
		ROLLBACK TRAN
		RETURN
	END

	-- Update/Inset data
	IF EXISTS( SELECT nCharNo FROM  tHolyPromise WHERE nCharNo=@nCharNo)
	BEGIN
		-- Check already upper here
		DECLARE @TempUpCharNo int
		SELECT @TempUpCharNo=nUpCharNo FROM  tHolyPromise WHERE nCharNo=@nCharNo
		IF @TempUpCharNo <> 0
		BEGIN
			SET @nRet = 4
			ROLLBACK TRAN
			RETURN
		END

		UPDATE tHolyPromise SET nUpCharNo=@nUpCharNo, nFlags=@nFlags, dRegDate=GetDate() WHERE nCharNo=@nCharNo
		SET @nRet = @@ERROR
		IF @nRet <> 0
		BEGIN
			ROLLBACK TRAN
			RETURN
		END
		SET @nRet = 3 -- OK. Rejoin		
	END

	ELSE

	BEGIN
		INSERT tHolyPromise (  nCharNo,  nUpCharNo,  nFlags )
		VALUES ( @nCharNo, @nUpCharNo, @nFlags )
		SET @nRet = @@ERROR
		IF @nRet <> 0
		BEGIN
			ROLLBACK TRAN
			RETURN
		END
	END
	
	COMMIT TRAN
	RETURN
	-- end
END

0
GO
