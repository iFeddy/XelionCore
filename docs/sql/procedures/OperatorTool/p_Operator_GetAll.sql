CREATE PROCEDURE [dbo].[p_Operator_GetAll]
/*
  Operator Get All
  output:
   - Select된 데이타
*/
AS
SET NOCOUNT ON
SELECT nOperNo, sOperID, nOperLevel, dRegDate FROM tOperator
-- end

0鐁羝đ
GO
