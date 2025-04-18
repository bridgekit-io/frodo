package cli

import (
	"fmt"
	"log"
	"strings"

	"github.com/bridgekit-io/frodo/generate"
	"github.com/bridgekit-io/frodo/parser"
	"github.com/spf13/cobra"
)

// GenerateClientRequest contains all of the CLI options used in the "frodo client" command.
type GenerateClientRequest struct {
	templateOption
	// InputFileName is the service definition to parse/process (the "--service" option)
	InputFileName string
	// Language is the programming language for the client to generate (the "--language" option)
	Language string
}

// TemplateName translates the Language option into the name of the template we should use for generation.
func (req GenerateClientRequest) TemplateName() string {
	switch strings.ToLower(req.Language) {
	case "go", "":
		return "client.go"
	case "js", "javascript", "node", "nodejs":
		return "client.js"
	case "dart", "flutter":
		return "client.dart"
	default:
		return ""
	}
}

// GenerateClient handles the registration and execution of the 'frodo client' CLI subcommand.
type GenerateClient struct{}

// Command creates the Cobra struct describing this CLI command and its options.
func (c GenerateClient) Command() *cobra.Command {
	request := &GenerateClientRequest{}
	cmd := &cobra.Command{
		Use:   "client [flags] FILENAME",
		Short: "Process a Go source file with your service interface to generate an RPC client proxy for your service(s).",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			request.InputFileName = args[0]
			crapPants(c.Exec(request))
		},
	}
	cmd.Flags().StringVar(&request.Language, "language", "go", "The file extension of the target language (e.g. 'go' or 'js')")
	cmd.Flags().StringVar(&request.Template, "template", "", "Path to a custom Go template file used to generate this artifact.")
	cmd.Flags().BoolVar(&request.Force, "force", false, "Ignore file modification timestamps and generate the artifact no matter what.")
	return cmd
}

// Exec takes all of the parsed CLI flags and generates the target client artifact.
func (c GenerateClient) Exec(request *GenerateClientRequest) error {
	templateName := request.TemplateName()
	if templateName == "" {
		return fmt.Errorf("unsupported client language")
	}

	if !request.Force && generate.UpToDate(request.InputFileName, templateName) {
		log.Printf("Skipping '%s'. Artifact is up to date '%s'", request.InputFileName, templateName)
		return nil
	}

	return c.generate(request, request.ToFileTemplate(templateName))
}

// generate parses the input service definition file and creates an output client/gateway
// code, writing it to the output gen/ directory.
func (c GenerateClient) generate(request *GenerateClientRequest, artifact generate.FileTemplate) error {
	log.Printf("Parsing service definition: %s", request.InputFileName)
	ctx, err := parser.ParseFile(request.InputFileName)
	if err != nil {
		return err
	}

	log.Printf("Generating '%s'", artifact.Name)
	return generate.File(ctx, artifact)
}
