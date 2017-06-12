# golang code review

This is an informal, WIP list of the things things I ([@gonzaloserrano](https://twitter.com/gonzaloserrano)) do when reviewing Golang code.

## general refs

- [go wiki: code review comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [go doc: effective go](https://golang.org/doc/effective_go.html)
- [best practices for PR reviews](https://github.com/kubernetes/community/blob/master/contributors/devel/pull-requests.md#best-practices-for-faster-reviews) by k8s team
- [go best practices](https://peter.bourgon.org/go-best-practices-2016/) by Peter Bourgon
- [ultimate go design guidelines](https://github.com/ardanlabs/gotraining/blob/master/topics/go/README.md) by Ardan Labs

## github

- better multiple PRs with small commits than multiple-commit PRs
- if have deps changes, separate the commit from code updates from the deps 
- squash, rebase or merge
  - deps in different commit -> merge
  - rest of code in same commit
- changes should:
  - be in green in current go version and tip
    - this can be configured in TravisCI for e.g
  - have same or better coverage, or < 1% of less coverage if really needed
    - use e.g coveralls.io

## code linting

The code should pass a linter step in the CI pipeline:

mandatory:
  - gofmt: official tool  to make all the code follow the same syntax guidelines.
  - go vet: official tool to find quircks in your code.
  - go lint: make your code follow the official code review guidelines.

others:
  - errcheck with -blank: check your errors are handled properly
  - interfacer: use the narrowest interface available.
  - gosimple: simplify your code whenever possible.
  - statticcheck: advanced tool to find subtle improvements and bugs in your code.

[gonmetalinter](https://github.com/gonzaloserrano/go-nzalo/blob/master/scripts/gonmetalinter): a dummy shell script I did for launching some linters

## software design

- onion / hexagonal / clean architecture
  - decorators / middlewares: logging, metrics, tracing...
  - see [this go-kit presentation example](https://youtu.be/NX0sHF8ZZgw?t=1344) by Peter Bourgon
- [code is communication](https://talks.golang.org/2014/readability.slide#44)
- [api design](https://talks.golang.org/2014/readability.slide#42)
- [always write the simplest code you can
](http://go-talks.appspot.com/github.com/davecheney/presentations/writing-high-performance-go.slide?utm_source=statuscode&utm_medium=medium#43)
- structs with just one field: `has-a` -> `is-a`
- [coupling with logger, metrics etc](https://peter.bourgon.org/go-best-practices-2016/#logging-and-instrumentation):
  - prefer decorator/middleware pattern
- avoid [primitive obsession - C#](http://enterprisecraftsmanship.com/2015/03/07/functional-c-primitive-obsession/)

## pkg design

- packagitis: this is not PHP or Java, where you put classes wherever and then you refer to them with the namespace.
- watch out import cycles! Go makes you think at compiler time how your types will be grouped.
- are the exported types needed to be exported?
  - maybe they are just used for tests
- Bill Kennedy's resources:
  - [blog article](https://www.goinggo.net/2017/02/design-philosophy-on-packaging.html)
  - [GopherCon India '17 video](https://www.youtube.com/watch?v=spKM5CyBwJA)

## errors

- coverage: at least happy path
- [design for errors - railway oriented programming](https://hackmd.io/JwDgbGDMBmAMBGBaArMA7GRAWATARgFNFgdoRFYwATdXEAQwGNGqg===)
  - https://dave.cheney.net/2015/01/26/errors-and-exceptions-redux
- [indent error flow: early returns](https://github.com/golang/go/wiki/CodeReviewComments#indent-error-flow)
  - [talk: keep the normal code path at a minimal indentation]( https://talks.golang.org/2014/readability.slide#27)
- [errors shouldn't be capitalized](https://github.com/golang/go/wiki/CodeReviewComments#error-strings)
- don't tag error messages in datadog metrics
- modeling:
  - don't couple tests to error messages
  - don't do errors.New(...); model them
  - use assert.Type in the test
- [how to handle errors](https://www.goinggo.net/2017/05/design-philosophy-on-logging.html) by Bill K.
  - wrap with [pkg/errors](https://github.com/pkg/errors) and log with `%v`
  - > Handling an error means:
    > * The error has been logged.
    > * The application is back to 100% integrity.
    > * The current error is not reported any longer.
- an http.Handler shouldn't repeat sending an errored response more than once;  separate the domain logic from the HTTP layer.

## concurrency

- protect shared mutable data
  - with Mutex/RWMutex
    - theory in spanish: [el problema de los lectores/escritores]( https://github.com/gallir/libro_concurrencia/blob/master/chapters/06-semaforos.adoc#lectores-escritores)
    - `mu sync.Mutex` over the fields it protects
    - use `RWMutex` if you have (significantly?) more reads than writes
  - or avoid sharing state and use message passing via channels
    - [channel design](https://github.com/ardanlabs/gotraining/tree/master/topics/go#channel-design)
      - suspect of buffered channels
    - [principles of designin APIs with channels](https://www.youtube.com/watch?v=hFqXgmor74k)
  - or use the `atomic` pkg if your state is a primitive (`float64` for e.g)
    - that's a nice thing for testing using spies in concurrent tests for e.g
  - or use things like `sync.Once` or `sync.Map` (>= go 1.9)
- testing:
  - concurrent tests are needed if your code is going to be run concurrently!
  - are tests run with `tests -race` ?
  - use parallel subtests when possible to make tests run faster ([official blog post](https://blog.golang.org/subtests))
    - must capture range vars! `tc := tc // capture range variable`
    - see also [using goroutines on loop iteration variables](https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables)
- be suspicious of anything similar to `go someFunc(...)`: why was it done?
  - performance: is there a benchmark for that?
  - async fire-and-forget logic: does that function return or not a value and/or error that must be handled?
- what's the lifetime of that goroutine? you could be leaking them: a goroutine lifetime must be always clear, subscribe to context.Done() to finish it.
    - [more from official core review](https://github.com/golang/go/wiki/CodeReviewComments#goroutine-lifetimes) 
    - [more from Dave Cheney](http://go-talks.appspot.com/github.com/davecheney/presentations/writing-high-performance-go.slide?utm_source=statuscode&utm_medium=medium#35)
- refs:
  - [pipelines and cancellation](https://blog.golang.org/pipelines)
  - [advanced concurrency patterns](https://blog.golang.org/advanced-go-concurrency-patterns)
- perf: mutex contention profiling tips from Jaana B Dogan:
  - https://talks.golang.org/2017/state-of-go.slide#23 
  - http://golang.rakyll.org/mutexprofile/

## http

- don't use `http.ListenAndServe`. Use `http.Server{}` with good defaults as explained in [how to handle http timeouts in go](https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/).
- an HTTP client must close the response body when finished with it, see the doc.

## naming

Be careful when using short names. Usually there is no prob with receivers
names, or temporal variables like indexes in loops etc. But if your funcs
violate the SRP or you have many args, then you can end with many short names
and it can make more harm than good. If understanding a var name is not
straightforward, use longer names.

other refs:
  - [avoid GetSomething getter name](https://golang.org/doc/effective_go.html#Getters)
  - [named result parameters](https://github.com/golang/go/wiki/CodeReviewComments#named-result-parameters)
  - [initialisms - ID, HTTP etc](https://github.com/golang/go/wiki/CodeReviewComments#initialisms)
  - [receiver name](https://github.com/golang/go/wiki/CodeReviewComments#receiver-names)
    > Don't name result parameters just to avoid declaring a var inside the function

## gotchas

- [50 Shades of Go](http://devs.cloudimmunity.com/gotchas-and-common-mistakes-in-go-golang/)
- a time.Ticker must be stopped, otherwise a goroutine is leaked
- slices hold a referente to the underlying array as explained here and here. Slices are passed by reference, maybe you think you are returning small data but it's underlying array can be huge.
- interfaces and nil gotcha, see the [understanding nil](https://www.youtube.com/watch?v=ynoY2xz-F8s) talk by Francesc Campoy

## performance

- [receiver type: value or pointer?](https://github.com/golang/go/wiki/CodeReviewComments#receiver-type)
- [Writing high performance go talk](http://go-talks.appspot.com/github.com/davecheney/presentations/writing-high-performance-go)
- [Escape Analysis and Memory Profiling - Bill Kennedy @ GopherCon SG 2017](https://www.youtube.com/watch?v=2557w0qsDV0)

## tests

- test names & scope
- coverage
  - edge cases
  - any change of behaviour should have a test that covers it
  - i.e new code yes, a refactor maybe not
- scope: if tests are complicated think about
  - refactor the code
  - refactor the tests
  - add arrange / act / assert comments
- types:
  - normal
  - internal
  - integration
  - example files
  - benchmarks
- subtests: see concurrency.
- tags: if you have tagged some tests (e.g _// +build integration_) then you need to run your tests with -tags 
- t.Run and t.Parallel
  - watch out with test cases in a loop, you need probably something like `tc := tc` before
- assert should have `(expected, actual)` not on the contrary

#### test doubles naming

Taken from [the little mocker](https://8thlight.com/blog/uncle-bob/2014/05/14/TheLittleMocker.html) by Uncle Bob.

- **dummy** objects are passed around but never actually used. Usually they are just used to fill parameter lists.
- **fake** objects actually have working implementations, but usually take some shortcut which makes them not suitable for production (an InMemoryTestDatabase is a good example).
- **stubs** provide canned answers to calls made during the test, usually not responding at all to anything outside what's programmed in for the test.
- **spies** are stubs that also record some information based on how they were called. One form of this might be an email service that records how many messages it was sent.
- **mocks** are pre-programmed with expectations which form a specification of the calls they are expected to receive. They can throw an exception if they receive a call they don't expect and are checked during verification to ensure they got all the calls they were expecting.
