.PHONY: reset clean install check build

reset:
	rm -rf node_modules

clean:
	rm -rf public/app dist

install:
	npm ci

# Vite transpiles with esbuild and never type-checks, so this is the only thing that
# will reject a type error.
check:
	npm run typecheck

build:
	npm run build
