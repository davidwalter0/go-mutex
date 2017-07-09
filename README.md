---
**mutex**

standalone scoped mutex

https://github.com/davidwalter0/go-mutex.git

Scoped execute auto unlock

Call Monitor with defer for scoped lock/unlock


*Create*

```
    var mtx *Mutex = mutex.NewMutex()
```

*Call*

```
    // Scoped call: acquire lock on entry, and release on scope
    // closure
    { // enter scope lock
        defer mtx.Monitor()()
    /*
     ...
     ...
     ...
    */
    } // exit scope unlock
```
