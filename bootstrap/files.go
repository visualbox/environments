package bootstrap


import (
  "os"
)

var (
	/**
	 *	Integration ID.
	 */
	ID := os.Getenv("ID")
	
	/**
	 *	Integration version.
	 */
	VERSION := os.Getenv("VERSION")
	
	/**
	 *	Initial integration configuration
	 *	model.
	 */
  MODEL := os.Getenv("MODEL")
)
