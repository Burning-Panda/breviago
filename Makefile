
clean:
	@echo "Cleaning up..."
	@echo "Cleanup DB files..."
	rm acronyms.db acronyms.db-shm-journal acronyms.db-shm acronyms.db-wal
	@echo "Cleanup build files..."
	rm -rf tmp

