// Copyright 2018 David Walter. All rights reserved.
// Use of this source code is governed by a Apache-style
// license that can be found in the LICENSE file.

/*
---
**mutex**

Note that in the following the _scope_ closure for a deferred call
of a lock is the function exit, therefore scoped lock calls lock at
the defer statement and close at the end of the function not at the
block *{}* level

standalone scoped mutex

https://github.com/davidwalter0/go-mutex.git

```
go get github.com/davidwalter0/go-mutex
```

Scoped execute auto unlock

Call Monitor with defer for scoped lock/unlock


*Create*

```
    import 	"github.com/davidwalter0/go-mutex"

    var mtx *Mutex = mutex.NewMutex()
```

*Call*

```
    // Scoped call: acquire lock on entry, and release on scope
    // closure
    { // enter scope lock
        defer mtx.Monitor()()

    // ...
    // ...
    // ...
    } // exit scope unlock
```


*Alternative Use*

With a shared scoped set of go routines or threads an anonymous monitor created by NewMonitor can be used

```
    import 	"github.com/davidwalter0/go-mutex"
    // Scoped call: acquire lock on entry, and release on scope
    // closure created anonymously in the function closure
    var m := mutex.NewMonitor()
    { // enter scope lock
        defer m()()
    // ...
    // ...
    // ...

    } // exit scope unlock
```

*/
package mutex
