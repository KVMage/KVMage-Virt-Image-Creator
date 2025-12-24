package cmd

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

type Options struct {
	Profile          string   `yaml:"-"`                         // Internal use only, not from YAML
	VirtMode         string   `yaml:"virt_mode"`                 // "install" or "customize"
	ImageName        string   `yaml:"image_name"`                // Name of the image
	OSVariant        string   `yaml:"os_var"`                    // OS variant (from osinfo-query)
	ImageSize        string   `yaml:"image_size"`                // e.g., "100G"
	ImagePartition   string   `yaml:"image_part"`                // e.g. "/dev/sda1"
	InstallFile      string   `yaml:"install_file" file:"true"`  // Path to auto install file
	InstallMedia     string   `yaml:"install_media" file:"true"` // Local or remote ISO file or install tree
	ImageSource      string   `yaml:"image_src" file:"true"`     // Source QCOW2 image (customize mode)
	ImageDestination string   `yaml:"image_dest" file:"true"`    // Output QCOW2 image
	Upload			 []string `yaml:"upload" file:"true"`		 // Files or directories to upload (temp)
	Execute			 []string `yaml:"execute": file:"true"`		 // Files to execute scripts (in order)
	Hostname         string   `yaml:"hostname"`                  // Optional
	Network          string   `yaml:"network"`                   // Optional virtual network
	Console          string   `yaml:"console"`                   // "serial", or "graphical"
	Firmware         string   `yaml:"firmware"`                  // "bios" or "efi" (default: bios)
}

func resolvePath(path string) string {
	if path == "" || filepath.IsAbs(path) {
		return path
	}
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path
	}
	cwd, err := os.Getwd()
	if err != nil {
		return path
	}
	return filepath.Join(cwd, path)
}

func (o *Options) ResolvePaths() {
	v := reflect.ValueOf(o).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("file")
		if tag == "true" && field.Type.Kind() == reflect.String {
			val := v.Field(i).String()
			v.Field(i).SetString(resolvePath(val))
        } else if tag == "true" && field.Type.Kind() == reflect.Slice && field.Type.Elem().Kind() == reflect.String {
		    slice := v.Field(i)
	    	for j := 0; j < slice.Len(); j++ {
				val := slice.Index(j).String()
				slice.Index(j).SetString(resolvePath(val))
			}
		}
	}
}
