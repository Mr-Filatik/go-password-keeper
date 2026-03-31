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
func Mask(data []byte, rules mask.IMaskable) (string, error) {
	parsed, parsErr := gabs.ParseJSON(data)
	if parsErr != nil {
		return string(data), parsErr
	}

	var errs []error

	opts := rules.Rules()
	for _, opt := range opts {
		editFn := opt.GetMaskFunc()
		deleteFn := opt.GetDeleteFunc()

		if editFn != nil && deleteFn == nil {
			err := editRecursive(parsed, opt.GetMaskFunc(), opt.GetPath())
			if err != nil {
				errs = append(errs, err)
			}
		}

		if editFn == nil && deleteFn != nil {
			err := deleteRecursive(parsed, opt.GetDeleteFunc(), opt.GetPath())
			if err != nil {
				errs = append(errs, err)
			}
		}
	}

	if len(errs) > 0 {
		err := fmt.Errorf("mask UserInfo error: %w", errors.Join(errs...))

		return parsed.String(), err
	}

	return parsed.String(), nil
}

//nolint:cyclop
func editRecursive(node *gabs.Container, fnc mask.EditFunc, path []string) error {
	if node == nil || len(path) == 0 {
		return nil
	}

	part := path[0]
	pathStr := strings.Join(path, ".")

	var errs []error

	if part == mask.Array {
		children := node.Children()
		for idx, child := range children {
			if len(path) == 1 {
				masked, merr := fnc(child.Data())
				if merr != nil {
					errs = append(errs, fmt.Errorf("masked %s failed: %w", pathStr, merr))

					continue
				}

				_, serr := node.SetIndex(masked, idx)
				if serr != nil {
					errs = append(errs, fmt.Errorf("set index %s failed: %w", pathStr, serr))
				}
			} else {
				eerr := editRecursive(child, fnc, path[1:])
				if eerr != nil {
					errs = append(errs, fmt.Errorf("edit %s failed: %w", pathStr, eerr))

					// return eerr
				}

				_, serr := node.SetIndex(child.Data(), idx)
				if serr != nil {
					errs = append(errs, fmt.Errorf("set index %s failed: %w", pathStr, serr))
				}
			}
		}

		if len(errs) > 0 {
			return errors.Join(errs...)
		}

		return nil
	}

	if len(path) == 1 {
		target := node.Search(part)
		if target.Data() == nil {
			return nil
		}

		masked, merr := fnc(target.Data())
		if merr != nil {
			errs = append(errs, fmt.Errorf("masked %s failed: %w", pathStr, merr))

			return errors.Join(errs...)
		}

		_, serr := node.Set(masked, part)
		if serr != nil {
			errs = append(errs, fmt.Errorf("set %s failed: %w", pathStr, serr))
		}

		if len(errs) > 0 {
			return errors.Join(errs...)
		}

		return nil
	}

	nextNode := node.Search(part)
	if nextNode.Data() == nil {
		return nil
	}

	eerr := editRecursive(nextNode, fnc, path[1:])
	if eerr != nil {
		errs = append(errs, fmt.Errorf("edit %s failed: %w", pathStr, eerr))

		return errors.Join(errs...)
	}

	return nil
}

func deleteRecursive(node *gabs.Container, fnc mask.DeleteFunc, path []string) error {
	if node == nil || len(path) == 0 {
		return nil
	}

	part := path[0]
	pathStr := strings.Join(path, ".")

	var errs []error

	// СЛУЧАЙ 1: Работа с массивом
	if part == mask.Array {
		children := node.Children()
		// Идем с конца в начало, чтобы удаление по индексу не ломало порядок
		for idx := len(children) - 1; idx >= 0; idx-- {
			child := children[idx]

			if len(path) == 1 {
				// Удаляем сам элемент массива, если функция вернула true
				ok, perr := fnc(child.Data())
				if perr != nil {
					errs = append(errs, fmt.Errorf("predicate %s failed: %w", pathStr, perr))

					continue // or delete?
				}

				if ok {
					rerr := node.ArrayRemove(idx)
					if rerr != nil {
						// add idx
						errs = append(errs, fmt.Errorf("remove index %s failed: %w", pathStr, rerr))
					}
				}
			} else {
				derr := deleteRecursive(child, fnc, path[1:])
				if derr != nil {
					errs = append(errs, fmt.Errorf("edit %s failed: %w", pathStr, derr))

					//return err
				}

				_, serr := node.SetIndex(child.Data(), idx)
				if serr != nil {
					errs = append(errs, fmt.Errorf("set index %s failed: %w", pathStr, serr))
				}
			}
		}

		if len(errs) > 0 {
			return errors.Join(errs...)
		}

		return nil
	}

	if len(path) == 1 {
		target := node.Search(part)
		if target.Data() == nil {
			return nil
		}

		ok, perr := fnc(target.Data())
		if perr != nil {
			errs = append(errs, fmt.Errorf("predicate %s failed: %w", pathStr, perr))
		}

		if ok {
			derr := node.Delete(part)
			if derr != nil {
				errs = append(errs, fmt.Errorf("delete %s failed: %w", pathStr, derr))
				// return fmt.Errorf("delete field %s failed: %w", part, derr)
			}
		}

		if len(errs) > 0 {
			return errors.Join(errs...)
		}

		return nil
	}

	nextNode := node.Search(part)
	if nextNode.Data() == nil {
		return nil
	}

	derr := deleteRecursive(nextNode, fnc, path[1:])
	if derr != nil {
		errs = append(errs, fmt.Errorf("edit %s failed: %w", pathStr, derr))

		return errors.Join(errs...)
	}

	return nil
}
