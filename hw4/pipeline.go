package executor

import (
    "context"
)

type (
    In  <-chan any
    Out = In
)

type Stage func(in In) (out Out)

func controlPipeline(ctx context.Context, lastIn In, out chan any) {
    defer close(out)
    for {
        select {
        case v, ok := <-lastIn:
            if !ok {
                return
            }
            out <- v
        case <-ctx.Done():
            return
        }
    }
}

func ExecutePipeline(ctx context.Context, in In, stages ...Stage) Out {
    nxtIn := in
    for _, stage := range stages {
        nxtIn = stage(nxtIn)
    }

    out := make(chan any)
    go controlPipeline(ctx, nxtIn, out)
    return out
}
