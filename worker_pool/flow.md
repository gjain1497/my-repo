        +----------------+
        |   Submit Job   |
        +----------------+
                 |
                 v
          +---------------+
          |   Jobs Queue  |   (channel of jobs)
          +---------------+
             ^   ^   ^
             |   |   |
      +------+   |   +------+
      |          |          |
      v          v          v
  +--------+  +--------+  +--------+
  | Worker |  | Worker |  | Worker |
  +--------+  +--------+  +--------+
      |          |          |
      v          v          v
   +-----------------------------+
   |         Results Queue       |
   +-----------------------------+
                 |
                 v
        +----------------+
        |   Get Results  |
        +----------------+

Jobs Queue: buffered/unbuffered channel that stores incoming jobs.

Workers: goroutines that continuously pick jobs, process them, and push results.

Results Queue: another channel to collect outputs from workers.