Your link

https://interview.codejudge.io/3a08094a55cb40cb9b3698e02b63c69a?user_type=RECRUITER&int_id=625750

Candidate link

https://interview.codejudge.io/3a08094a55cb40cb9b3698e02b63c69a?user_type=DEVELOPER

Instructions

Introduction

3 Mins

Problem Solving

15 Mins

Guidelines

Ask puzzle-based questions.

1. Should be able to solve complex coding problems.

2. Should be able to think out of box to resolve issues.


Sample questions

Add the asked question here

Type or paste your question here
Evaluation metrics

Problem Understanding & Structured Reasoning

NA12345678910
Solution Design & Translation to Code

NA12345678910
Optimization & Scalability Mindset

NA12345678910
Data Structures and Algorithms

10 Mins

Coding

20 Mins

Database

10 Mins

Conclusion

2 Mins

Overall feedback

Interview complete? Share your feedback


Fill Feedback


Questions (3)

Question 1

Question:

You have 8 balls, 7 are equal, 1 is heavier. Find the heavier one in 2 weighings.
Follow-up: Write a program that generalizes this for N balls and outputs the minimum weighings required.
Assessment Points:

Logical reasoning on first step.
Translate puzzle → function.
Correct formula derivation: ⌈log₃(N)⌉ weighings.
Question 2

Traveling Salesman (Simplified)

Question:

Given a graph of 5 cities with pairwise distances, write an algorithm to compute the shortest route that visits all cities once and returns to the start.
Constraint: Brute force is acceptable since cities ≤ 5.
Assessment Points:

Recognizes NP-hard nature.
Uses permutations/DFS for small N.
Mentions scalability limits + possible heuristics (greedy/DP).
Question 3

Monty Hall Simulation

Question:

Write a program to simulate the Monty Hall problem with 1000 runs and calculate the probability of winning if:
You always stick with your first choice.
You always switch.
Assessment Points:

Ability to translate a probability puzzle into simulation code.
Code structure, randomness, reproducibility.
Insight: Switching ≈ 66% win rate.



Data Structures and Algorithms (Thoery only)
Data Structure and Algorithms

1. Should be able to compare and identify best matched DS for the requirements

2. Should be able to compare and identify best matched algo for the requirements

3. Should be able to optimize for space complexity and time complexity.

Justifies best DS/Algo choice, proactively optimize, handle complexity.



Discuss a DSA problem, do not ask to implement the code
Question 1

Given a 2D grid of 0s (water) and 1s (land), count the number of distinct islands, where:

An island is a group of 1s connected 4-directionally.
Two islands are considered the same if one is a rotation or reflection of the other.
Follow-up Discussion Points:

How would you use DFS to capture the shape of each island?
How do you normalize shapes to handle rotation and reflection?
Time and space complexity concerns for very large grids.
How to hash island shapes canonically.


Question 2

Given a directed graph with red and blue edges, find the longest path (not necessarily simple) that alternates edge colors at each step.

Input:

n nodes.
List of red edges and blue edges.
Follow-up Discussion Points:

DFS with edge color state tracking.
Memoization to avoid redundant paths.
Can cycles be allowed? If yes, how to prevent infinite loops?
Can this be solved in topological order or needs dynamic programming?


Question 3

Given a directed graph, determine the minimum number of edges required to make the graph strongly connected (i.e., there’s a path from every node to every other node).

Follow-up Discussion Points:

Use of Kosaraju’s or Tarjan’s algorithm to identify strongly connected components (SCCs).
Condensation of the graph into a DAG.
Counting sources and sinks in the SCC DAG to calculate the result as max(#sources, #sinks).





Coding : 

1. Should be able to write good readable code with robust error handling and corner cases. We can specify this requirement to the candidate as well.

2. Should be able to specify unit tests for code written.

3. Should be able to identify gaps/issues with the hint.

4. Should verify the code (without asking) before requesting the review.

Write robust, production-grade code

The code is readable, with meaningful variable names and comments. It meets the requirements. It is modularised with input validations. 

The code is further made user friendly. 

We should ask different test scenarios for this code and observe if the candidate is able to identify gaps and improve it further.

Question 1

 Log Parser with Error Handling

Problem:

Write a function that takes a log file (string input, each line = "timestamp level message") and returns the count of each log level (INFO, WARN, ERROR).

Requirements:

Handle invalid/malformed log lines gracefully.
Support extensibility (new log levels).
Unit test for: empty file, only one level present, malformed line in between.
Discussion Follow-up:

How would you scale this for very large files (streaming)?
How would you design unit tests for performance?
Question 2

 JSON Data Flattener

Problem:

Given a nested JSON object, flatten it into a single-level dictionary.

Example:



{"a": {"b": 1, "c": {"d": 2}}, "e": 3}

➡ {"a.b": 1, "a.c.d": 2, "e": 3}

Requirements:

Handle invalid JSON input.
Support different separators (., _) as parameter.
Add unit tests for: deeply nested object, empty object, large object.
Discussion Follow-up:

How would you modify for unflattening back?
How would you optimize memory usage for very deep nesting?
Question 3

 Rate Limiter (API Simulation)

Problem:

Design a simple in-memory rate limiter function allow_request(user_id) that allows max 3 requests per user per 10 seconds.

Requirements:

Handle multiple users.
Handle invalid inputs (null user_id).
Make it production-friendly (clear old data, configurable limits).
Unit test for:
User under limit.
User exceeding limit.
Multiple users at once.
Discussion Follow-up:

How would you implement this in a distributed system?
How do you avoid memory leaks with long-lived users?
Question 4

 File System Path Normalizer

Problem:

Write a function that normalizes a file path string:

Removes . (current dir),
Resolves .. (parent dir),
Collapses multiple / into one.
Example:

Input: "/a//b/./c/../d/"

Output: "/a/b/d"

Requirements:

Handle invalid paths (e.g., null, .. at root).
Ensure cross-platform support (/ vs \).
Unit test for: empty path, root path, long nested path.
Discussion Follow-up:

How would you extend for Windows/Unix compatibility?
How would you handle symbolic links?





Database : 

Should have exposure to RDBMS, Nosql DB.

Should be able to write and optimise queries.

Question 1

Query Optimization & Indexing

You have a table Orders(order_id, customer_id, order_date, amount) with millions of rows.

Write a query to fetch the top 5 customers by total purchase amount in the last 6 months.
Follow-up: How would you optimize this query for performance? (Think indexes, partitioning, materialized views).
What to assess:

Correct use of GROUP BY, ORDER BY, LIMIT.
Understanding of covering indexes, composite indexes.
Trade-off between runtime aggregation vs pre-aggregated tables.
Question 2

Handling Large Joins

You have two large tables:

Users(user_id, name, country)
Transactions(txn_id, user_id, amount, txn_date)
Write a query to find the average transaction amount per country.

Follow-up: How would you handle performance if both tables are very large (hundreds of millions of rows)?
What to assess:

Correct join syntax, GROUP BY.
Indexing strategies (join keys, covering index).
Awareness of partitioning, sharding, or denormalization.
Question 3

Querying Hierarchical Data

You have a table Employees(id, name, manager_id).

Write a query to find the entire reporting hierarchy under a given manager (recursive).
Follow-up: How would you optimize recursive queries in RDBMS? Would you model this differently in a NoSQL DB?
What to assess:

Recursive CTEs in SQL.
Thinking about graph/tree traversal.
Trade-off: RDBMS recursion vs storing hierarchy in NoSQL (adjacency list, materialized path, nested sets).
Question 4

RDBMS vs NoSQL Decision

You are building a real-time recommendation engine that needs to:

Store user activity logs (billions of writes/day).
Query recent activity per user with low latency.
Question: Would you choose RDBMS or NoSQL for this? Why?

Follow-up: What type of NoSQL (Key-Value, Document, Columnar, Graph) would you pick, and why?
What to assess:

Understanding of CAP theorem trade-offs.
Ability to reason about write-heavy vs read-heavy workloads.
Choosing between RDBMS and specific NoSQL models.