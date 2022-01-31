# golang code review

This is an informal, WIP list of the things things I ([@gonzaloserrano](https://twitter.com/gonzaloserrano)) do when reviewing Golang code.

## general refs

#### Google
- [go wiki: code review comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [go doc: effective go](https://golang.org/doc/effective_go.html)
- go proverbs: [website](https://go-proverbs.github.io/) - [talk video](https://www.youtube.com/watch?v=PAAkCSZUG1c)

#### 3rd party

- [go best practices](https://peter.bourgon.org/go-best-practices-2016/) by Peter Bourgon
- [ultimate go design guidelines](https://github.com/ardanlabs/gotraining/blob/master/topics/go/README.md) by Ardan Labs
- [idiomatic go](https://about.sourcegraph.com/go/idiomatic-go) by Sourcegraph
- [the zen of go](https://dave.cheney.net/2020/02/23/the-zen-of-go) by Dave Cheney

## code reviewing & github

- the [guidelines for faster PR reviews](https://github.com/kubernetes/community/blob/master/contributors/devel/pull-requests.md#best-practices-for-faster-reviews) from the Kubernetes project are a must. A quick summary:
    - do small commits and, even better, small PRs
    - use separate PRs for fixes not related with your feature
    - add a different commit for review changes so it's easy to review them instead of the whole patch
    - test and document your code
    - don't add features you don't need
- other guidelines:
    - [How to Write a Git Commit Message])https://cbea.ms/git-commit)
    - prefer documentation as code (example tests files) over READMEs
    - separate the vendor updates in a different commit
    - choose a good GitHub merge strategy:
      - choose merge when:
          - multiple commits 
          - deps in different commit
      - choose squash if you just have a single commit to avoid the extra merge commit
    - do Continuous Integration (CI) to ensure the code quality:
        - tests are in green (`go test -race` in TravisCI or CircleCI)
        - new code or changes in functionality has the corresponding tests (e.g `gocovmerge` + codecov.io or coveralls.io)

## code linting

Just use [golangci-lint](https://golangci-lint.run/). The defaults are fine. The more linters you enable, the more you can avoid nitpicks in the PR reviews.

As a formatter, I personally prefer [gofumpt: a stricter gofmt](https://github.com/mvdan/gofumpt).

## software design

- onion / hexagonal / clean architecture
  - use decorator/middleware when possible
    - e.g don't pollute your `http.Handler`s with to instrumentation (logging, metrics, tracing)
  - see [this go-kit presentation example](https://youtu.be/NX0sHF8ZZgw?t=1344) by Peter Bourgon
  - and [Embrace the Interface](https://www.youtube.com/watch?v=xyDkyFjzFVc) video by Tomas Senart ([@tsenart](https://twitter.com/tsenart))
  - [another good blog post with examples](https://medium.com/@matryer/writing-middleware-in-golang-and-how-go-makes-it-so-much-fun-4375c1246e81)
- example: [coupling with logger, metrics etc](https://peter.bourgon.org/go-best-practices-2016/#logging-and-instrumentation):
  - don't log everything, maybe using metrics is better ([blog post by @copyconstruct](https://medium.com/@copyconstruct/logs-and-metrics-6d34d3026e38))
  - [the RED method](https://www.weave.works/blog/the-red-method-key-metrics-for-microservices-architecture/) for choosing metrics.
- [code is communication](https://talks.golang.org/2014/readability.slide#44)
- [api design](https://talks.golang.org/2014/readability.slide#42)
    - [make your APIs sync](https://about.sourcegraph.com/go/idiomatic-go/#asynchronous-apis)
- [always write the simplest code you can
](http://go-talks.appspot.com/github.com/davecheney/presentations/writing-high-performance-go.slide?utm_source=statuscode&utm_medium=medium#43)
- structs with just one field: `has-a` -> `is-a`
- avoid [primitive obsession - C#](http://enterprisecraftsmanship.com/2015/03/07/functional-c-primitive-obsession/)
- [accept interfaces, return concrete types](http://idiomaticgo.com/post/best-practice/accept-interfaces-return-structs/)
- [DDD in go](https://gist.github.com/abdullin/3e3fd199674255e4d206)
- [think about using functional options in your constructors](https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis)

## package design

- learn and follow the stdlib design
- packagitis: this is not PHP or Java, where you put classes wherever and then you refer to them with the namespace.
- watch out import cycles! Go makes you think at compiler time how your types will be grouped.
- are the exported types needed to be exported?
  - maybe they are just used for tests
- reduce the number of package dependencies
    - sometimes [duplication is better than the wrong abstraction](https://www.sandimetz.com/blog/2016/1/20/the-wrong-abstraction?duplication)

#### other refs

- Ben Johnson article on [standard package layout](https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1)
- Bill Kennedy's resources:
  - [blog article](https://www.goinggo.net/2017/02/design-philosophy-on-packaging.html)
  - [GopherCon India '17 video](https://www.youtube.com/watch?v=spKM5CyBwJA)
- [style guides for go pkgs - @rakyll](https://rakyll.org/style-packages):
    - use multiple files
    - keep types closer to where they are used
    - organize by responsibility (e.g avoid package `model`)

## errors

- coverage: at least happy path
- [design for errors - railway oriented programming](https://hackmd.io/JwDgbGDMBmAMBGBaArMA7GRAWATARgFNFgdoRFYwATdXEAQwGNGqg===)
  - https://dave.cheney.net/2015/01/26/errors-and-exceptions-redux
- [indent error flow: early returns](https://github.com/golang/go/wiki/CodeReviewComments#indent-error-flow)
  - [talk: keep the normal code path at a minimal indentation]( https://talks.golang.org/2014/readability.slide#27)
- [errors shouldn't be capitalized](https://github.com/golang/go/wiki/CodeReviewComments#error-strings)
- put tests in a different package
    - else name the file `*_internal_test.go`
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
- don't panic, return errors - specially in libs.

## concurrency

- read [golang concurrency tricks](https://udhos.github.io/golang-concurrency-tricks/)
- go makes concurrency easy enough to be dangerous [source](https://www.youtube.com/watch?v=DJ4d_PZ6Gns&t=1270s)
- shared mutable state is the root of all evil [source](http://henrikeichenhardt.blogspot.com.es/2013/06/why-shared-mutable-state-is-root-of-all.html)
- how protect shared mutable data
  - with Mutex/RWMutex
    - theory in spanish: [el problema de los lectores/escritores](https://github.com/gallir/libro_concurrencia/blob/master/chapters/06-semaforos.adoc#lectores-escritores)
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
  - _good_ concurrent tests are mandatory and run with `-race` flag
  - use parallel subtests when possible to make tests run faster ([official blog post](https://blog.golang.org/subtests))
    - must capture range vars! `tc := tc // capture range variable` (also see next point)
- watch out when:
  - launching goroutines inside loops, see [using goroutines on loop iteration variables](https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables)
  - passing pointers through channels
  - you see the `go` keyword around without much explanation and tests
    - why was it done?
    - is there a benchmark for that?
  - async fire-and-forget logic: does that function return or not a value and/or error that must be handled?
- goroutine lifetime:
  - watch out for leaks: lifetime must be always clear, e.g terminate a `for-select` subscribing to context.Done() to finish return.
  - [more from official core review](https://github.com/golang/go/wiki/CodeReviewComments#goroutine-lifetimes) 
  - [more from Dave Cheney](http://go-talks.appspot.com/github.com/davecheney/presentations/writing-high-performance-go.slide?utm_source=statuscode&utm_medium=medium#35)
- if you write libs, leave concurrency to the callerwhen possible. See [Zen Of Go](https://the-zen-of-go.netlify.app/).
    - your code will be simpler, and the clients will choose the kind of concurrency they want
- [don't expose channels](https://about.sourcegraph.com/go/idiomatic-go/#asynchronous-apis)
- [the channel closing principle](http://www.tapirgames.com/blog/golang-channel-closing): don't close a channel from the receiver side and don't close a channel if the channel has multiple concurrent senders.
- refs:
  - [pipelines and cancellation](https://blog.golang.org/pipelines)
  - [advanced concurrency patterns](https://blog.golang.org/advanced-go-concurrency-patterns)
- performance: 
  - mutex contention [profiling](https://talks.golang.org/2017/state-of-go.slide#23 ) [tips](http://golang.rakyll.org/mutexprofile/) from Jaana B Dogan

## http

- don't use `http.ListenAndServe`. Use `http.Server{}` with good defaults as explained in [how to handle http timeouts in go](https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/).
- an HTTP client [must close the response body](https://golang.org/pkg/net/http/#pkg-overview) when finished reading with it.
  - [closing the body must be done after error checking](https://stackoverflow.com/a/42525774/547956)
- when creating HTTP-based libs, allow injecting an `http.Client`.

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
    - [the right way](https://docs.google.com/presentation/d/e/2PACX-1vQ9aFgICdqCz5pjrVJ4zFZrWtTKbfGYFPCOKsScomkLoE1Kzk3DVd9-u4k_XgZekqJ7nl-YTy4lD8Uq/pub?slide=id.g550f852d27_228_0) by [Daniel Martí](https://twitter.com/mvdan_)
- subtests: see concurrency.
- tags: if you have tagged some tests (e.g _// +build integration_) then you need to run your tests with -tags 
- t.Run and t.Parallel
  - watch out with test cases in a loop, you need probably something like `tc := tc` before
- assert should have `(expected, actual)` not on the contrary
- test doubles naming (from [the little mocker](https://8thlight.com/blog/uncle-bob/2014/05/14/TheLittleMocker.html) by Uncle Bob):
  - **dummy** objects are passed around but never actually used
  - **fake** objects actually have working implementations, but usually take some shortcut which makes them not suitable for production (an InMemoryTestDatabase is a good example).
  - **stubs** provide canned answers to calls made during the test, usually not responding at all to anything outside what's programmed in for the test.
  - **spies** are stubs that also record some information based on how they were called. One form of this might be an email service that records how many messages it was sent.
  - **mocks** are pre-programmed with expectations which form a specification of the calls they are expected to receive. They can throw an exception if they receive a call they don't expect and are checked during verification to ensure they got all the calls they were expecting.

## aws-sdk-go

- Create just a single session in the top level, see [the doc](https://github.com/aws/aws-sdk-go/blob/master/aws/session/doc.go#L4-L11)
- Every service package has an interface that you can use for embedding in test, see an example in the [official blog](https://aws.amazon.com/blogs/developer/mocking-out-then-aws-sdk-for-go-for-unit-testing)

## performance

- [receiver type: value or pointer?](https://github.com/golang/go/wiki/CodeReviewComments#receiver-type)
- [Writing high performance go talk](http://go-talks.appspot.com/github.com/davecheney/presentations/writing-high-performance-go)
- [Escape Analysis and Memory Profiling - Bill Kennedy @ GopherCon SG 2017](https://www.youtube.com/watch?v=2557w0qsDV0)
- from [so you wanna go fast?](https://youtu.be/DJ4d_PZ6Gns?t=2193) talk by by Tyler Treat:
  - the stdlib provides general solutions, and you should generally use them.
  - small, idiomatic changes can have profound performance impact.
  - learn and use tools from the go toolchain to analyze the performance.
  - the performance can change a lot between go versions, review your optimizations.
  - code is marginal, architecture is material: the big wins come from architecutre, do it right first.
  - mechanical sympathy: know how yoour abstractions actually work in your hardware, go makes this possible.
  - optimize for the right trade-off: optimizing for performance means trading something else
