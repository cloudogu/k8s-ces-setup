.PHONY: clean
clean: $(ADDITIONAL_CLEAN)
	rm -rf ${TARGET_DIR}
	rm -rf ${TMP_DIR}
	rm -rf ${UTILITY_BIN_PATH}

.PHONY: dist-clean
dist-clean: clean
	rm -rf node_modules
	rm -rf public/vendor
	rm -rf vendor
	rm -rf npm-cache
	rm -rf bower
