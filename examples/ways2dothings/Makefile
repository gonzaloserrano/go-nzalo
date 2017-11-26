goto:
	$(info [[ REMINDER: GO BACK TO SLIDES ]])
t1:
	go test -v ./... --run TestDoOne
t2:
	go test -v ./... --run TestDoSeveralOK
t3:
	go test -v ./... --run TestDoSeveralConcurrentWrong
t4:
	go test -v ./... --run TestDoSeveralConcurrentClientSide

# ---------------------

t5: goto
	go test -v ./... --run TestDoSeveralConcurrentServerSideStuck
t6:
	go test -v ./... --run TestDoSeveralConcurrentSignalClosesChan
t7:
	go test -v ./... --run TestDoSeveralConcurrentSignalCancelsContext

# ---------------------

t8: goto
	go test -v ./... --run TestDoSeveralConcurrentWithBSK

# ---------------------

t9: goto
	go test -v ./... --run TestDoSeveralConcurrentWithErrgroup

# ---------------------

t10: goto
	go test -v ./... --run TestDoSeveralConcurrentWithOKLogRun
