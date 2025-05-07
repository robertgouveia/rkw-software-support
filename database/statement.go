package database

const DisputeChange = `UPDATE [dbo].[DeliveryIssuesHead] SET [Status] = @Status WHERE IssueID = @IssueID`

const ShippingChange = `UPDATE [dbo].[Goods Outward Header] SET [Shipping Agent Service] = '48' WHERE [Sales Order No_] = @OrderNo`
