.PHONY: dist
dist:
	../../../../_dockers/js/env.sh docker_build_pj $(shell pwd)

.PHONY: dist-clean
dist-clean:
	rm -rf dist

.PHONY: clean
clean:
	rm -rf node_modules .pnpm-store

dev:
	../../_env/env.sh pj_pnpm $(shell pwd) run dev