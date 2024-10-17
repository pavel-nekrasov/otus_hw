package main

import (
	"errors"
	"io"
	"os"
	"time"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrSourceFileNotFound    = errors.New("file in 'from' path not found")
	ErrTargetCannotBeCreated = errors.New("file in 'to' path cannot be created")
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrFileCannotBeCopied    = errors.New("file cannot be copied")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	const ChunkSizeDefault = 1024
	const ProgressDelayMs = 50

	source, err := os.Open(fromPath)
	if err != nil {
		if os.IsNotExist(err) {
			return ErrSourceFileNotFound
		}
		return ErrUnsupportedFile
	}
	defer source.Close()

	target, err := os.Create(toPath)
	if err != nil {
		return ErrTargetCannotBeCreated
	}
	defer target.Close()

	sourceStat, err := source.Stat()
	if err != nil {
		return ErrUnsupportedFile
	}

	fileSize := sourceStat.Size()

	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}
	if limit == 0 || limit+offset > fileSize {
		limit = fileSize - offset
	}

	bar := pb.StartNew(int(limit))
	defer bar.Finish()

	var chunkSize, processed int64

	for processed < limit {
		_, err := source.Seek(offset+processed, 0)
		if err != nil {
			return ErrUnsupportedFile
		}

		if limit-processed < ChunkSizeDefault {
			chunkSize = limit - processed
		} else {
			chunkSize = ChunkSizeDefault
		}

		written, err := io.CopyN(target, source, chunkSize)
		if err != nil {
			return ErrFileCannotBeCopied
		}

		bar.Add(int(written))
		time.Sleep(ProgressDelayMs * time.Millisecond)
		processed += written
	}

	return nil
}
