CREATE PROCEDURE [dbo].[usp_Statistics_getCharacterLevel]
	@dDate nvarchar(10) = ''
AS
BEGIN
	SET NOCOUNT ON

	IF LEN(@dDate) = 0 BEGIN
		SET @dDate = CONVERT(NVARCHAR(10), GETDATE(), 120)
	END

	SELECT nLevel AS 'Level',
		   nTotalBefore + nTotalAfter AS 'Total',
		   nBefore + nAfter AS 'Exp>0',
		   nTotalBefore AS 'Total(Before 2009-04-22)',
		   nBefore AS 'Exp>0(Before 2009-04-22)',
		   nTotalAfter AS 'Total(After 2009-04-22)',
		   nAfter AS 'Exp>0(After 2009-04-22)'
	  FROM tStatCharLevel WITH (NOLOCK)
	 WHERE dDate = @dDate

	SET NOCOUNT OFF
END

0
GO
