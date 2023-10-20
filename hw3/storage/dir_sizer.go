package storage

import (
    "context"
)

// Result represents the Size function result
type Result struct {
    // Total Size of File objects
    Size int64
    // Count is a count of File objects processed
    Count int64
}

type DirSizer interface {
    // Size calculate a size of given Dir, receive a ctx and the root Dir instance
    // will return Result or error if happened
    Size(ctx context.Context, d Dir) (Result, error)
}

// sizer implement the DirSizer interface
type sizer struct {
    // maxWorkersCount number of workers for asynchronous run
    maxWorkersCount int
}

// NewSizer returns new DirSizer instance
func NewSizer() DirSizer {
    return &sizer{8}
}

func dirSize(ctx context.Context, ac chan<- Result, ec chan<- error, wp chan struct{}, dir Dir) {
    defer func() {
        <-wp
    }()

    dirs, files, err := dir.Ls(ctx)
    if err != nil {
        ec <- err
        return
    }

    childrenResultChan := make(chan Result, len(dirs))
    childrenErrorChan := make(chan error, len(dirs))
    for i := range dirs {
        wp <- struct{}{}
        go dirSize(ctx, childrenResultChan, childrenErrorChan, wp, dirs[i])
    }

    result := Result{}
    for i := range files {
        fileSize, fileErr := files[i].Stat(ctx)
        if fileErr != nil {
            ec <- fileErr
            return
        }
        result.Size += fileSize
        result.Count += 1
    }

    for range dirs {
        select {
        case err := <-childrenErrorChan:
            ec <- err
            return
        case <-ctx.Done():
            return
        case childrenResult := <-childrenResultChan:
            result.Size += childrenResult.Size
            result.Count += childrenResult.Count
        }
    }
    ac <- result
}

func (a *sizer) Size(ctx context.Context, d Dir) (Result, error) {
    resultChan := make(chan Result)
    errorChan := make(chan error)
    wp := make(chan struct{}, a.maxWorkersCount)
    go dirSize(ctx, resultChan, errorChan, wp, d)
    select {
    case err := <-errorChan:
        return Result{}, err
    case result := <-resultChan:
        return result, nil
    case <-ctx.Done():
        return Result{}, nil
    }
}
