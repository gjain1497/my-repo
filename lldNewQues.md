Aggregating User Activities from JSON

A social media platform needs an activity aggregation service to summarize the count of different activity types per user for analytics dashboards.

Design an API endpoint POST /aggregate-activities that:
Accepts a JSON array of user activity objects.
Aggregates counts of each activity type per user.
Returns a dictionary mapping user IDs to activity counts.
Provide a database schema for storing raw activities and aggregated results.
Mention the design pattern you would use to implement the aggregation logic.


Summing Nested Dictionary Values

A configuration management system stores metrics in nested JSON/dictionary structures. The analytics service must compute the total sum of all numeric values in these structures.

Design a utility function or microservice POST /sum-nested-values that:
Accepts a nested JSON/dictionary as input.
Recursively sums all numeric values.
Returns the total sum.
Provide a database schema if storing nested metric structures.
Specify the design pattern used for recursive processing.



Sales by Region for a Given Year

A retail company wants to build a sales analytics service that allows users to filter sales by year and view aggregated totals per region.

Design an API endpoint GET /sales-summary?year=YYYY that:

Takes a year as input.
Queries the database for total sales per region for that year.
Returns JSON output with region names as keys and total sales as values.
Define the database schema for storing sales data.
Explain which design pattern you would use for separating data access and business logic.
