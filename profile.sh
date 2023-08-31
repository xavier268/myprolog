#! /bin/bash

echo "Starting a profiling capture"

go test -coverprofile cover.out ./...
go test -memprofile solver_mem.prof -cpuprofile solver_cpu.prof ./solver
go tool pprof -png solver_mem.prof > solver_mem.png
go tool pprof -png solver_cpu.prof > solver_cpu.png



