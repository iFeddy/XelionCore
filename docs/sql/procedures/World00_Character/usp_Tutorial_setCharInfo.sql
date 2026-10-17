CREATE PROCEDURE [dbo].[usp_Tutorial_setCharInfo]
	@nUserNo int
,	@nCharNo int
,	@nState  tinyint
,	@nStep   tinyint
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	SET @nRet = 0

    IF EXISTS (SELECT nCharNo FROM tTutorial WHERE nUserNo = @nUserNo AND nCharNo = @nCharNo)
    BEGIN
        UPDATE tTutorial SET nState = @nState, nStep = @nStep WHERE nUserNo = @nUserNo AND nCharNo = @nCharNo;
        IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 SET @nRet = -1;
    END
    ELSE
    BEGIN
		/*-------------------------------------------------------------------------------------
		Project: FiestaHeroes - NA2016
		Author: xkl
	
		Skip the character tutorial.
		------------------------------------------------------------------------------------*/
		IF EXISTS (SELECT Enabled FROM Options..Tutorial_Skip WHERE Enabled = 1)
            INSERT INTO tTutorial (nUserNo, nCharNo, nState, nStep) VALUES (@nUserNo, @nCharNo, 2, 0);
        ELSE
            INSERT INTO tTutorial (nUserNo, nCharNo, nState, nStep) VALUES (@nUserNo, @nCharNo, 0, 0);
		/*-------------------------------------------------------------------------------------
		Continue as normal ~ 
		------------------------------------------------------------------------------------*/
		
        IF @@ERROR <> 0 SET @nRet = -1;
    END
END
GO
