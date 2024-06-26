package main

func handleMake(arg2 string, arg3 string) error {

	switch arg2 {
	case "migration":
		err := handleMakeMigration(arg3)
		if err != nil {
			exitGracefully(err)
		}

	case "auth":
		err := handleAuth()
		if err != nil {
			exitGracefully(err)
		}

	case "session":
		err := handleSession()
		if err != nil {
			exitGracefully(err)
		}

	case "handler":
		err := handleHandler(arg3)
		if err != nil {
			exitGracefully(err)
		}

	case "model":
		err := handleModel(arg3)
		if err != nil {
			exitGracefully(err)
		}

	case "mail":
		err := handleMail(arg3)
		if err != nil {
			exitGracefully(err)
		}
	}

	return nil
}
