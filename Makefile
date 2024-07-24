.DEFAULT_GOAL := run

IMAGE_PREFIX := simonswine
IMAGE_TAG := latest
BUILD_PLATFORM := linux/amd64,linux/arm64

# This dataset is licensed under a Creative Commons Attribution 4.0 License
# https://www.geonames.org/
data/%.txt:
	mkdir -p data/
	curl -o data/$*.zip https://download.geonames.org/export/dump/$*.zip
	unzip data/$*.zip $*.txt -d data
	rm -rf data/$*.zip

.PHONY: data
data: data/DE.txt data/US.txt data/GB.txt data/ES.txt data/FR.txt data/CA.txt

.PHONY: build
build:
	docker buildx build --load -t $(IMAGE_PREFIX)/find-city-by-population .

.PHONY: push
push:
	docker buildx build --push --platform $(BUILD_PLATFORM) -t $(IMAGE_PREFIX)/find-city-by-population:$(IMAGE_TAG) .

.PHONY: bench
bench:
	go test -run=XXX -bench=BenchmarkFindCityBy -benchtime 10s -cpuprofile cpu.pb.gz -memprofile mem.pb.gz

.PHONY: run
run: data
	go run ./


