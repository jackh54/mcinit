package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/jackh54/mcinit/internal/cache"
	"github.com/jackh54/mcinit/internal/config"
	"github.com/jackh54/mcinit/internal/java"
	"github.com/jackh54/mcinit/internal/provider"
	"github.com/jackh54/mcinit/internal/scripts"
	"github.com/jackh54/mcinit/internal/utils"
	"github.com/jackh54/mcinit/pkg/jvmflags"
	"github.com/spf13/cobra"
)

var (
	serverType  string
	mcVersion   string
	serverPath  string
	serverName  string
	acceptEula  bool
	ram         string
	xms         string
	xmx         string
	jvmFlags    string
	port        int
	nogui       bool
	gitignore   bool
	javaVersion string
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a new server folder with configuration",
	Long: `Initialize a new Minecraft server with the specified configuration.
Downloads the server jar, generates startup scripts, and creates all necessary files.`,
	Example: `  mcinit init --type vanilla --mc 1.21.4 --accept-eula --ram 4G
  mcinit init --type paper --mc 1.21.4 --accept-eula --ram 4G
  mcinit init --type paper --mc 1.21.4 --path ./test-server --xms 2G --xmx 6G --flags minimal
  mcinit init --type purpur --mc 1.21.4 --java 21 --port 25566 --nogui`,
	RunE: runInit,
}

func init() {
	initCmd.Flags().StringVar(&serverType, "type", "paper", "Server type (vanilla|paper|folia|purpur|velocity|waterfall|bungee)")
	initCmd.Flags().StringVar(&mcVersion, "mc", "", "Minecraft version (required)")
	initCmd.Flags().StringVar(&serverPath, "path", "./server", "Target directory for server")
	initCmd.Flags().StringVar(&serverName, "name", "", "Server name (default: derived from path)")
	initCmd.Flags().BoolVar(&acceptEula, "accept-eula", false, "Accept Minecraft EULA")
	initCmd.Flags().StringVar(&ram, "ram", "", "Total RAM (e.g., 4G) - sets both Xms and Xmx")
	initCmd.Flags().StringVar(&xms, "xms", "", "Initial heap size (e.g., 2G)")
	initCmd.Flags().StringVar(&xmx, "xmx", "", "Maximum heap size (e.g., 4G)")
	initCmd.Flags().StringVar(&jvmFlags, "flags", "aikar", "JVM flags preset (aikar|minimal|custom)")
	initCmd.Flags().IntVar(&port, "port", 25565, "Server port")
	initCmd.Flags().BoolVar(&nogui, "nogui", false, "Disable server GUI")
	initCmd.Flags().BoolVar(&gitignore, "gitignore", false, "Add server path to .gitignore")
	initCmd.Flags().StringVar(&javaVersion, "java", "auto", "Java version or path (auto|17|21|/path/to/java)")

	_ = initCmd.MarkFlagRequired("mc")
}

func runInit(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	// Validate server type
	if !provider.Exists(serverType) {
		return fmt.Errorf("invalid server type: %s (available: %v)", serverType, provider.List())
	}

	// Handle RAM flags
	if ram != "" {
		if xms != "" || xmx != "" {
			return fmt.Errorf("cannot specify both --ram and --xms/--xmx")
		}
		xms = ram
		xmx = ram
	} else {
		if xms == "" {
			xms = "2G"
		}
		if xmx == "" {
			xmx = "4G"
		}
	}

	// Resolve server path
	absPath, err := utils.AbsolutePath(serverPath)
	if err != nil {
		return fmt.Errorf("failed to resolve server path: %w", err)
	}

	// Derive server name from path if not provided
	if serverName == "" {
		serverName = filepath.Base(absPath)
	}

	printf("Initializing %s server (Minecraft %s) at %s\n", serverType, mcVersion, absPath)

	// Dry-run mode
	if dryRun {
		printf("[DRY RUN] Would create server with:\n")
		printf("  Type: %s\n", serverType)
		printf("  Version: %s\n", mcVersion)
		printf("  Path: %s\n", absPath)
		printf("  Name: %s\n", serverName)
		printf("  RAM: Xms=%s Xmx=%s\n", xms, xmx)
		printf("  Flags: %s\n", jvmFlags)
		printf("  Port: %d\n", port)
		printf("  EULA: %v\n", acceptEula)
		return nil
	}

	// Create server directory
	if err := utils.EnsureDir(absPath); err != nil {
		return fmt.Errorf("failed to create server directory: %w", err)
	}

	// Check if config already exists
	cfgPath := filepath.Join(absPath, "mcinit.json")
	if config.Exists(cfgPath) {
		return fmt.Errorf("server already initialized (mcinit.json exists)")
	}

	// Detect and validate Java
	printf("Detecting Java installation...\n")
	javaDetector := java.NewDetector()
	javaValidator := java.NewValidator()

	var javaInst *java.Installation
	if javaVersion != "" && javaVersion != "auto" {
		// Check if it's a path or version number
		if utils.PathExists(javaVersion) {
			inst, err := javaDetector.DetectVersion(javaVersion)
			if err != nil {
				return fmt.Errorf("invalid Java path: %w", err)
			}
			javaInst = inst
		} else {
			// Parse as version number
			var majorVersion int
			_, _ = fmt.Sscanf(javaVersion, "%d", &majorVersion)
			inst, err := javaDetector.FindByMajorVersion(majorVersion)
			if err != nil {
				return fmt.Errorf("java %d not found: %w", majorVersion, err)
			}
			javaInst = inst
		}
	} else {
		// Auto-detect best Java
		inst, err := javaDetector.FindBest()
		if err != nil {
			return fmt.Errorf("no Java installation found: %w", err)
		}
		javaInst = inst
	}

	// Validate Java for Minecraft version
	if err := javaValidator.ValidateForMinecraft(javaInst, mcVersion); err != nil {
		errorLog("Java validation warning: %v\n", err)
		errorLog("Server may not start correctly\n")
	}

	printf("Found Java %s at %s\n", javaInst.Version, javaInst.Path)

	// Get provider and download jar
	printf("Downloading %s server jar for Minecraft %s...\n", serverType, mcVersion)
	prov, err := provider.Get(serverType)
	if err != nil {
		return fmt.Errorf("failed to get provider: %w", err)
	}

	localPath, downloadURL, checksum, err := prov.DownloadJar(ctx, mcVersion, "latest")
	if err != nil {
		return fmt.Errorf("failed to download server jar: %w", err)
	}

	// Copy jar to server directory
	destJarPath := filepath.Join(absPath, "server.jar")
	if err := copyFile(localPath, destJarPath); err != nil {
		return fmt.Errorf("failed to copy server jar: %w", err)
	}

	printf("Server jar downloaded successfully\n")

	// Create configuration
	cfg := config.DefaultConfig()
	cfg.Server.Type = serverType
	cfg.Server.MinecraftVersion = mcVersion
	cfg.Server.JarPath = "server.jar"
	cfg.Server.Name = serverName
	cfg.Server.DownloadURL = downloadURL
	algorithm, _, _ := prov.GetChecksum(ctx, mcVersion, "latest")
	if algorithm == "sha256" {
		cfg.Server.SHA256 = checksum
	} else {
		cfg.Server.SHA1 = checksum
	}

	cfg.Java.Version = javaVersion
	cfg.Java.Path = javaInst.Path
	cfg.Java.DetectedVersion = javaInst.Version

	cfg.JVM.Xms = xms
	cfg.JVM.Xmx = xmx
	cfg.JVM.Flags = jvmFlags

	cfg.ServerConfig.Port = port
	cfg.ServerConfig.NoGUI = nogui

	cfg.EULA.Accepted = acceptEula
	if acceptEula {
		cfg.EULA.AcceptedAt = time.Now().UTC()
	}

	// Get cache directory
	cache, _ := cache.New()
	cfg.Paths.CacheDir = cache.GetBaseDir()
	cfg.Paths.ServerDir = "."

	// Save configuration
	if err := config.Save(cfg, cfgPath); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	printf("Configuration saved to mcinit.json\n")

	// Create EULA file
	eulaPath := filepath.Join(absPath, "eula.txt")
	eulaContent := fmt.Sprintf("# Generated by mcinit\neula=%v\n", acceptEula)
	if err := os.WriteFile(eulaPath, []byte(eulaContent), 0644); err != nil {
		return fmt.Errorf("failed to create eula.txt: %w", err)
	}

	// Generate server.properties
	serverPropsPath := filepath.Join(absPath, "server.properties")
	serverPropsContent := fmt.Sprintf("server-port=%d\n", port)
	if err := os.WriteFile(serverPropsPath, []byte(serverPropsContent), 0644); err != nil {
		return fmt.Errorf("failed to create server.properties: %w", err)
	}

	// Generate startup scripts
	printf("Generating startup scripts...\n")
	scriptGen := scripts.NewGenerator(absPath)

	jvmFlagsArr, err := jvmflags.Get(jvmFlags, xms, xmx, cfg.JVM.CustomFlags)
	if err != nil {
		return fmt.Errorf("failed to build JVM flags: %w", err)
	}

	if err := scriptGen.GenerateAll(javaInst.Path, "server.jar", xms, xmx, jvmFlagsArr); err != nil {
		return fmt.Errorf("failed to generate scripts: %w", err)
	}

	printf("Startup scripts generated\n")

	// Handle .gitignore
	if gitignore {
		if utils.IsInGitRepo(absPath) {
			gitRoot, err := utils.FindGitRoot(absPath)
			if err == nil {
				relPath, err := utils.RelativePath(gitRoot, absPath)
				if err == nil {
					if err := utils.AddToGitignore(gitRoot, relPath); err != nil {
						errorLog("Failed to update .gitignore: %v\n", err)
					} else {
						printf("Added %s to .gitignore\n", relPath)
					}
				}
			}
		} else {
			errorLog("Not in a git repository, skipping .gitignore\n")
		}
	}

	// Create .mcinit directory
	mcinitDir := filepath.Join(absPath, ".mcinit")
	if err := utils.EnsureDir(mcinitDir); err != nil {
		return fmt.Errorf("failed to create .mcinit directory: %w", err)
	}

	printf("\n")
	printf("Server initialized successfully!\n")
	printf("\n")
	printf("Next steps:\n")
	printf("  cd %s\n", serverPath)
	if !acceptEula {
		printf("  # Accept the EULA by editing eula.txt\n")
	}
	printf("  mcinit start\n")
	printf("\n")

	return nil
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() { _ = sourceFile.Close() }()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() { _ = destFile.Close() }()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
