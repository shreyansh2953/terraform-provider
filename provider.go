package bashprovider

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider function where we execute Git setup script using schema variables
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"git_username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The global Git username to configure.",
			},
			"git_email": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The global Git email to configure.",
			},
			"git_token": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Git personal access token for authentication.",
			},
		},
		ConfigureFunc: providerConfigure, // This is called to configure the provider.
	}
}

// Function that runs when the provider is initialized.
func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	// Fetch the Git setup variables from the schema
	gitUsername := d.Get("git_username").(string)
	gitEmail := d.Get("git_email").(string)
	gitToken := d.Get("git_token").(string)

	// Execute the script to configure Git globally
	err := setupGitCredentials(gitUsername, gitEmail, gitToken)
	if err != nil {
		fmt.Println("Error")
	}

	return nil, nil
}

// Function to execute Git credential setup script
func setupGitCredentials(username, email, token string) error {
	// Script to set Git credentials globally
	script := fmt.Sprintf(`
		git config --global user.name "%s"
		git config --global user.email "%s"
		git config --global credential.helper store
		echo "https://%s:x-oauth-basic@github.com" > ~/.git-credentials
	`, username, email, token)

	// Execute the bash script
	cmd := exec.Command("bash", "-c", script)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println(`Error executing Git credential setup`, err)
	}

	fmt.Println("Git credentials setup complete.")
	return nil
}
