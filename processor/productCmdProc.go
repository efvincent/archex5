//
// Command Processors accept commands, evaluate their validity, and determine if an event
// should be created. The way to think about it - a command is a request to do something
// that has not yet been done.
//
package processor

import (
	"errors"
	"fmt"
	"log"

	"github.com/efvincent/archex5/commands"
)

// Dispatches the product command to the appropriate command handler
func ProcessProductCommand(cmd interface{}) error {
	switch c := cmd.(type) {
	case *commands.CreateProductCmd:
		return ProcessCreateProduct(c)
	case *commands.HeadCheckCmd:
		return PerformHeadCheck(c)
	case *commands.SetActiveCmd:
		return SetProductActiveState(c)
	case *commands.UpdatePriceCmd:
		return UpdatePrice(c)
	case *commands.UpdateProductAttributesCmd:
		return UpdateAttribs(c)
	case *commands.UpdateProductImagesCmd:
		return UpdateImages(c)
	default:
		return errors.New(fmt.Sprintf("Unknown command type: %v", c))
	}
}

func UpdateImages(cmd *commands.UpdateProductImagesCmd) error {
	log.Printf("Updating images for %s in %s", cmd.SKU, cmd.Namespace)
	return nil
}

func UpdateAttribs(cmd *commands.UpdateProductAttributesCmd) error {
	log.Printf("Updating product attributes for %s in %s", cmd.SKU, cmd.Namespace)
	return nil
}

func UpdatePrice(cmd *commands.UpdatePriceCmd) error {
	log.Printf("Updating price on %s in %s", cmd.SKU, cmd.Namespace)
	return nil
}

func SetProductActiveState(cmd *commands.SetActiveCmd) error {
	log.Printf("Changing Product active state for %s in %s", cmd.SKU, cmd.Namespace)
	return nil
}

func PerformHeadCheck(cmd *commands.HeadCheckCmd) error {
	// In reality, we'd probably want an object for each image and be able to track the headcheck events
	// on a per image basis
	log.Printf("Performing HeadCheck on SKU %s in %s", cmd.SKU, cmd.Namespace)
	return nil
}

func ProcessCreateProduct(cmd *commands.CreateProductCmd) error {
	log.Printf("Checking if SKU %s in namespace %s exists", cmd.Product.SKU, cmd.Namespace)
	log.Printf("Writing product created event for %s in %s", cmd.Product.SKU, cmd.Namespace)
	return nil
}
