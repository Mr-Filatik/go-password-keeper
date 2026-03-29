// Package gabs describes the implementation of the masker via the github.com/Jeffail/gabs/v2 package.
package gabs

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Jeffail/gabs/v2"
	"github.com/mr-filatik/go-password-keeper/internal/platform/sequrity/mask"
)

// Mask masks data represented as a slice of bytes according to the rules
// described in the IMaskable interface for a specific type.
//
// Implementation of the mask.IMasker interface
// from the github.com/mr-filatik/go-password-keeper/internal/platform/sequrity/mask package.
//
//nolint:cyclop
func Mask(data []byte, rules mask.IMaskable) (string, error) {
	parsed, parsErr := gabs.ParseJSON(data)
	if parsErr != nil {
		return string(data), parsErr
	}

	var errs []error

	opts := rules.Rules()
	for _, opt := range opts {
		switch opt.GetMaskType() {
		case mask.MTypeNone:
			continue

		case mask.MTypeEdit:
			err := edit(parsed, opt.GetMaskFunc(), opt.GetPath())
			if err != nil {
				errs = append(errs, err)
			}

		case mask.MTypeEditInSlice:
			err := editInSlice(parsed, opt.GetMaskFunc(), opt.GetPath(), opt.GetIntPath())
			if err != nil {
				errs = append(errs, err)
			}

		case mask.MTypeRemove:
			continue

		case mask.MTypeRemoveInSlice:
			continue
		}
	}

	if len(errs) > 0 {
		err := fmt.Errorf("mask UserInfo error: %w", errors.Join(errs...))

		return parsed.String(), err
	}

	return parsed.String(), nil
}

func edit(cnt *gabs.Container, fnc mask.MFunc, path []string) error {
	strPath := strings.Join(path, ".")

	item := cnt.Search(path...)
	if item == nil {
		return nil // Not an error for optional fields
	}

	maskedData, err := fnc(item.Data())
	if err != nil {
		return fmt.Errorf("mask in %s failed: %w", strPath, err)
	}

	_, err = cnt.Set(maskedData, path...)
	if err != nil {
		return fmt.Errorf("set value in %s failed: %w", strPath, err)
	}

	return nil
}

//nolint:cyclop
func editInSlice(cnt *gabs.Container, fnc mask.MFunc, path []string, intpath []string) error {
	item := cnt.Search(path...)
	if item == nil {
		return nil // Not an error for optional fields
	}

	childrens := item.Children()

	for idx, child := range childrens {
		strPath := fmt.Sprintf("%s.%d", strings.Join(path, "."), idx)

		if child == nil {
			continue // Not an error for optional fields
		}

		if len(intpath) == 0 {
			data := child.Data()

			maskedData, err := fnc(data)
			if err != nil {
				return fmt.Errorf("mask in %s failed: %w", strPath, err)
			}

			_, err = item.SetIndex(maskedData, idx)
			if err != nil {
				return fmt.Errorf("set value in %s failed: %w", strPath, err)
			}

			continue
		}

		strIntPath := strings.Join(intpath, ".")
		strFullPath := fmt.Sprintf("%s.%s", strPath, strIntPath)

		target := child.Search(intpath...)
		if target == nil {
			continue // Not an error for optional fields
		}

		data := target.Data()

		maskedData, err := fnc(data)
		if err != nil {
			return fmt.Errorf("mask in %s failed: %w", strFullPath, err)
		}

		_, err = child.SetP(maskedData, strIntPath)
		if err != nil {
			return fmt.Errorf("set value in %s failed: %w", strFullPath, err)
		}

		_, err = item.SetIndex(child.Data(), idx)
		if err != nil {
			return fmt.Errorf("set child in %s failed: %w", strPath, err)
		}
	}

	return nil
}
