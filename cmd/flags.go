package cmd

var opts = &Options{}

var (
        showVersion       bool
        checkRequirements bool
        verboseLevel      int
        quiet             bool
        uninstall         bool
        cleanupOnly       bool

        runMode    bool
        configPath string
        envFilePath string

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
        rootCmd.PersistentFlags().BoolVarP(&checkRequirements, "check-requirements", "R", false, "Check system requirements and exit")

        rootCmd.PersistentFlags().BoolVarP(&runMode, "run", "r", false, "Use CLI args to run KVMage")
        rootCmd.PersistentFlags().StringVarP(&configPath, "config", "f", "", "Path to YAML config file (defaults to kvmage.yml if not specified)")
        rootCmd.PersistentFlags().StringVar(&envFilePath, "env-file", "", "Path to env file for variable substitution")
        rootCmd.PersistentFlags().Lookup("config").NoOptDefVal = "AUTO"

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
        rootCmd.Flags().StringSliceVarP(&opts.Upload, "upload", "U", []string{}, "Files or directories to upload (temp)")
        rootCmd.Flags().StringSliceVarP(&opts.Execute, "execute", "E", []string{}, "Files to execute scripts (in order)")
        rootCmd.Flags().StringVarP(&opts.Network, "network", "W", "", "Network interface (optional)")
        rootCmd.Flags().StringVarP(&opts.Console, "console", "", "", "Console type: serial or graphical (optional)")
        rootCmd.Flags().StringVarP(&opts.Firmware, "firmware", "m", "bios", "Firmware type: bios or efi")
}