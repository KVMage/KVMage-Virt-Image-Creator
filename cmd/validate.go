package cmd

import (
	"fmt"
	"reflect"
)

var fieldLabels = map[string]string{
	"ImageName":        "image_name",
	"OSVariant":        "os_var",
	"ImageSize":        "image_size",
	"InstallFile":      "install_file",
	"InstallMedia":     "install_media",
	"ImageSource":      "image_src",
	"ImageDestination": "image_dest",
}

var requiredFields = map[string][]string{
	"install": {
		"ImageName", "OSVariant", "ImageSize",
		"InstallFile", "ImageDestination",
	},
	"customize": {
		"ImageName", "ImageSource", "ImageDestination",
	},
}

func ValidateOptions(opts *Options) error {
	mode := opts.VirtMode
	if mode != "install" && mode != "customize" {
		return fmt.Errorf("invalid or missing virt_mode: must be 'install' or 'customize'")
	}

	missing := []string{}
	v := reflect.ValueOf(opts).Elem()

	for _, field := range requiredFields[mode] {
		val := v.FieldByName(field)
		if !val.IsValid() || val.String() == "" {
			missing = append(missing, fieldLabels[field])
		}
	}

	if mode == "install" {
		hasInstallMedia := opts.InstallMedia != ""

		switch {
		case !hasInstallMedia:
			missing = append(missing, "install_media")
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("%s mode is missing required fields: %v", mode, missing)
	}

	if opts.Console != "" && opts.Console != "serial" && opts.Console != "graphical" {
		return fmt.Errorf("invalid console value: must be 'serial', 'graphical', or unset")
	}

	if opts.Firmware != "" && opts.Firmware != "bios" && opts.Firmware != "efi" {
		return fmt.Errorf("invalid firmware value: must be 'bios', 'efi', or unset")
	}

	return nil
}
