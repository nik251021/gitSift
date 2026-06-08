package handlers

import "fmt"

const helpText = `
goSift — Your AI companion for working with Git repositories.

Usage:
  goSift <command> [arguments]

Commands:
  setAPIkey            Save your Gemini API key to the .env file
  getAPIkey            Retrieve the current Gemini API key from the .env file
  setCurrentLinkToRepo Set the active target repository URL
  getCurrentLinkToRepo View the active target repository URL
  sift                 Fetch and generate a raw context cache from the repository URL
  ask                  Query Gemini about the active cached repository code
  help                 Display this help reference manual

Examples:
  goSift setAPIkey "AIzaSy..."
  goSift getAPIkey
  goSift setCurrentLinkToRepo "https://github.com/user/repo"
  goSift ask "Explain the main architectural entry point"
  goSift help
`

func HelpHandler(args []string) error {
	fmt.Print(helpText)
	return nil
}
