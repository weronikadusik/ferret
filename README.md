# ferret
A system monitoring tool that tells you why, not just what 🦦

## Sorting

```
# Sort process list by pid
ferret --sort=pid

# Sort process list by cpu usage
ferret --sort=cpu

# Sort process list by memory usage
ferret --sort=memory

# Sort process list by disk usage
ferret --sort=disk
```

> Temporary CLI flag — will move to interactive sorting once the TUI (bubbletea) lands.