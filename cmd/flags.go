package cmd

var opts = &Options{}

var (
	showVersion  bool
	verboseLevel int
	quiet        bool
	uninstall    bool
	cleanupOnly  bool

	runMode    bool
	configPath string

	installFlag   bool
	customizeFlag bool
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "V", false, "Show version and exit")
	rootCmd.PersistentFlags().CountVarP(&verboseLevel, "verbose", "v", "Increase verbosity (-v, -vv, -vvv). Equivalent to --verbose-level=N")
	rootCmd.PersistentFlags().IntVar(&verboseLevel, "verbose-level", 0, "Set verbosity level explicitly (0-3)")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Suppress all output (overrides verbose)")
	rootCmd.PersistentFlags().BoolVarP(&uninstall, "uninstall", "u", false, "Uninstall KVMage from /usr/local/bin")
	rootCmd.PersistentFlags().BoolVarP(&cleanupOnly, "cleanup", "X", false, "Run cleanup mode to remove orphaned kvmage temp files")

	rootCmd.PersistentFlags().BoolVarP(&runMode, "run", "r", false, "Use CLI args to run KVMage")
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "f", "", "Path to YAML config file")

	rootCmd.Flags().BoolVarP(&installFlag, "install", "i", false, "Run in install mode")
	rootCmd.Flags().BoolVarP(&customizeFlag, "customize", "c", false, "Run in customize mode")

	rootCmd.Flags().StringVarP(&opts.ImageName, "image-name", "n", "", "Name of the image")
	rootCmd.Flags().StringVarP(&opts.OSVariant, "os-var", "o", "", "OS variant")
	rootCmd.Flags().StringVarP(&opts.ImageSize, "image-size", "s", "", "Image size")
	rootCmd.Flags().StringVarP(&opts.ImagePartition, "image-part", "P", "", "Partition to expand inside image (e.g. /dev/sda1)")
	rootCmd.Flags().StringVarP(&opts.InstallMedia, "install-media", "j", "", "Path to local or remote ISO file or install tree")
	rootCmd.Flags().StringVarP(&opts.InstallFile, "install-file", "k", "", "Install file path")
	rootCmd.Flags().StringVarP(&opts.ImageSource, "image-src", "S", "", "Source qcow2 image")
	rootCmd.Flags().StringVarP(&opts.ImageDestination, "image-dest", "D", "", "Destination qcow2 image")
	rootCmd.Flags().StringVarP(&opts.Hostname, "hostname", "H", "", "Hostname (optional)")
	rootCmd.Flags().StringVarP(&opts.CustomScript, "custom-script", "C", "", "Custom script (optional)")
	rootCmd.Flags().StringVarP(&opts.Network, "network", "W", "", "Network interface (optional)")
	rootCmd.Flags().StringVarP(&opts.Console, "console", "", "", "Console type: serial or graphical (optional)")
	rootCmd.Flags().StringVarP(&opts.Firmware, "firmware", "m", "bios", "Firmware type: bios or efi")
}
