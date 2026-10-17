create procedure [dbo].[Adam_ListProc] as  select o.name, m.definition   from sys.sql_modules m join sys.objects o  On m.object_id=o.object_id and o.type='P'

0
GO
