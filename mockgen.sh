#!/bin/bash
mockgen -source=internal/film/repository.go \
  -destination=internal/film/mock/repository_mock.go \
  -package=mock

mockgen -source=internal/actor/repository.go \
  -destination=internal/actor/mock/repository_mock.go \
  -package=mock