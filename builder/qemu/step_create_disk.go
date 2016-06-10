package qemu

import (
	"fmt"
	"path/filepath"

	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
)

// This step creates the virtual disk that will be used as the
// hard drive for the virtual machine.
type stepCreateDisk struct{}

func (s *stepCreateDisk) Run(state multistep.StateBag) multistep.StepAction {
	config := state.Get("config").(*Config)
	driver := state.Get("driver").(Driver)
	ui := state.Get("ui").(packer.Ui)
	name := config.VMName

	if config.DiskImage == true {
		return multistep.ActionContinue
	}

	command := []string{
		"create",
		"-f", config.Format,
		filepath.Join(config.OutputDir, name),
		fmt.Sprintf("%vM", config.DiskSize),
	}

	ui.Say("Creating hard drive...")
	if err := driver.QemuImg(command...); err != nil {
		err := fmt.Errorf("Error creating hard drive: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	state.Put("disk_filename", name)

	if len(config.AdditionalDiskSize) > 0 {
		additionalPaths := make([]string, len(config.AdditionalDiskSize))
		ui.Say("Creating additional hard drives...")

		for i, additionalsize := range config.AdditionalDiskSize {
			additionalPath := filepath.Join(config.OutputDir, fmt.Sprintf("%s-%d", name, i+1))
			command := []string{
				"create",
				"-f", config.Format,
				additionalPath,
				fmt.Sprintf("%vM", additionalsize),
			}

			ui.Say("Creating hard drive...")
			if err := driver.QemuImg(command...); err != nil {
				err := fmt.Errorf("Error creating hard drive: %s", err)
				state.Put("error", err)
				ui.Error(err.Error())
				return multistep.ActionHalt
			}
			additionalPaths[i] = additionalPath
		}
		state.Put("additional_disk_paths", additionalPaths)
	}
	return multistep.ActionContinue
}

func (s *stepCreateDisk) Cleanup(state multistep.StateBag) {}
