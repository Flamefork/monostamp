# Monostamp

Wallclock-aligned monotonic timestamp generator.

## Features
- Customizable scale|precision: nanosecond, microsecond, millisecond, second, custom
- Customizable start timestamp for monotonic continuation
- Generated timestamp drift from wall clock reporting

## Usage

```go
import "github.com/Flamefork/monostamp"

m := monostamp.New(
    monostamp.UnixMicro, // microsecond scale
    0,                   // start from wall clock timestamp
    monostamp.NewDriftReporter(
        1000,            // don't warn about drifts less than 1ms (1000 microseconds)
        10000000,        // don't warn more often than once per 10s (10M microseconds)
        func(generatedTS int64, clockTS int64) {
            log.Printf("warn: timestamp drifted by %dms", (generatedTS - clockTS)/1000)
        },
    ).Report,
)

m.Next() // 1686052801000000
m.Next() // 1686052801000001
m.Next() // 1686052801000002
time.Sleep(time.Second)
m.Next() // 1686052802000000
m.Next() // 1686052802000001
```

## License

Copyright 2023 Ilia Ablamonov

Permission is hereby granted, free of charge, to any person obtaining a copy of this
software and associated documentation files (the “Software”), to deal in the Software
without restriction, including without limitation the rights to use, copy, modify,
merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit
persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR
PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE
FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
DEALINGS IN THE SOFTWARE.
