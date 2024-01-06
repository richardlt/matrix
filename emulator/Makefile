reset:
	rm -rf node_modules

clean:
	rm -rf public

install:
	docker run --rm \
		-v $(PWD):/tmp/workspace -w /tmp/workspace \
		node:20 \
		npm ci

build:
	docker run --rm \
		-v $(PWD):/tmp/workspace -w /tmp/workspace \
		node:20 \
		npm run build
