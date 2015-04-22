# Cuckoo Filter

[![Build Status](https://travis-ci.org/asp2insp/cuckoofilter.svg)](https://travis-ci.org/asp2insp/cuckoofilter)

This is a Golang implementation of [Cuckoo Filters](https://www.cs.cmu.edu/~dga/papers/cuckoo-conext2014.pdf)


# Usage


Basic usage:

```go
// Create a cuckoo table with 10,000 buckets, 4 entries per bucket, using
// 5 bytes per fingerprint representation. Allow 500 iterations of the cuckoo
// displacement before failing on insert.
table := NewCuckooTable(10000, 500, 4, 5)
table.Insert([]byte("Hello World"))
table.Lookup([]byte("Hello World")) // true
table.Size() // 0
table.Stats() // utilization, rebucketRatio, compressionRatio

table.Delete([]byte("Hello World"))
```

See main.go for a benchmark using the `/usr/share/dict/words` file.

The library is "batteries included" which means it contains a hashing function and a fingerprinting function. There's a little big of work that needs to be done to allow the user to provide a hash function and a fingerprinting function of their own.


# Contributing

Pull requests are welcome!


# License

Copyright 2014 Josiah Gaskin

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.