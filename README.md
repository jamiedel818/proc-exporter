# Mini Node Exporter

For later with prometheus https://signoz.io/guides/why-does-prometheus-consume-so-much-memory/

Note as of 10/25/25
- wrote the meminfo parser and tests. Is to think how I want to take the data and output them as prometheus metrics.
- Make it so the program takes a path to the proc directory and then file names are used by collectors. 
- I want introduce the collection logic. eg, scan every X seconds. Aquire the lock to write, scan, update data, release the lock. Probably in its own goroutine
- Implement HTTP server logic. `/metrics`. When request comes in, grab the lock to read, return the data in prometheus format, release the lock.
- Add in some config values in the form of env vars. eg `SCAN_FREQUENCY`, `PROC_DIRECTORY_PATH`. 
- Add some basic logging.

Step 1.
- learn about the proc files I am interested in and how to understand their contents.

One file I want is `meminfo`. This is in kiB
```
cat /proc/meminfo
MemTotal:        3888568 kB
MemFree:         2838832 kB
MemAvailable:    3551232 kB
Buffers:           72016 kB
Cached:           675900 kB
SwapCached:            0 kB
Active:           714800 kB
Inactive:         185964 kB
Active(anon):     159456 kB
Inactive(anon):        0 kB
Active(file):     555344 kB
Inactive(file):   185964 kB
Unevictable:           0 kB
Mlocked:               0 kB
SwapTotal:        524284 kB
SwapFree:         524284 kB
Zswap:                 0 kB
Zswapped:              0 kB
Dirty:                 8 kB
Writeback:             0 kB
AnonPages:        152932 kB
Mapped:           260464 kB
Shmem:              6608 kB
KReclaimable:      40656 kB
Slab:              75564 kB
SReclaimable:      40656 kB
SUnreclaim:        34908 kB
KernelStack:        6752 kB
PageTables:         4784 kB
SecPageTables:         0 kB
NFS_Unstable:          0 kB
Bounce:                0 kB
WritebackTmp:          0 kB
CommitLimit:     2468568 kB
Committed_AS:    1767120 kB
VmallocTotal:   261087232 kB
VmallocUsed:       28380 kB
VmallocChunk:          0 kB
Percpu:             1440 kB
CmaTotal:         524288 kB
CmaFree:          511596 kB
```

Step 2.
- Open the proc files. I can probably write a general open file function for this.
- What is the best way to monitor a file that changes frequently. eg: if I am going to open this every X seconds. 

Step 3.
- Parse out the information you care about from the proc files.
- Should I do this sequentially or in parallel. eg: Open meminfo, then open cpuinfo, or do them in parallel. latter is more complicated especially if I use a global state.

Step 4.
- Store the information

Step 5.
- Serve the information as a standard prometheus metric. Learn this format.

